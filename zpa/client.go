package zpa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"

	"github.com/zscaler/zscaler-sdk-go/v2/cache"
	"github.com/zscaler/zscaler-sdk-go/v2/logger"
	"github.com/zscaler/zscaler-sdk-go/v2/utils"
)

type Client struct {
	Config *Config
	cache  cache.Cache
}

// NewClient returns a new client for the specified apiKey.
func NewClient(config *Config) (c *Client) {
	if config == nil {
		config, _ = NewConfig("", "", "", "", "")
	}
	cche, err := cache.NewCache(config.cacheTtl, config.cacheCleanwindow, config.cacheMaxSizeMB)
	if err != nil {
		cche = cache.NewNopCache()
	}
	c = &Client{Config: config, cache: cche}
	return
}

func (client *Client) WithFreshCache() {
	client.Config.freshCache = true
}

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	req, err := client.getRequest(method, url, options, body)
	if err != nil {
		return nil, err
	}
	key := cache.CreateCacheKey(req)
	if client.Config.cacheEnabled {
		if req.Method != http.MethodGet {
			// this will allow to remove resource from cache when PUT/DELETE/PATCH requests are called, which modifies the resource
			client.cache.Delete(key)
			// to avoid resources that GET url is not the same as DELETE/PUT/PATCH url, because of different query params.
			// example delete app segment has key url/<id>?forceDelete=true but GET has url/<id>, in this case we clean the whole cache entries with key prefix url/<id>
			client.cache.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
		}
		resp := client.cache.Get(key)
		inCache := resp != nil
		if client.Config.freshCache {
			client.cache.Delete(key)
			inCache = false
			client.Config.freshCache = false
		}
		if inCache {
			if v != nil {
				respData, err := io.ReadAll(resp.Body)
				if err == nil {
					resp.Body = io.NopCloser(bytes.NewBuffer(respData))
				}
				if err := decodeJSON(respData, v); err != nil {
					return resp, err
				}
			}
			unescapeHTML(v)
			client.Config.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return resp, nil
		}
	}
	resp, err := client.newRequestDoCustom(method, url, options, body, v)
	if err != nil {
		return resp, err
	}
	if client.Config.cacheEnabled && resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet && v != nil && reflect.TypeOf(v).Kind() != reflect.Slice {
		d, err := json.Marshal(v)
		if err == nil {
			resp.Body = io.NopCloser(bytes.NewReader(d))
			client.Config.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
			client.cache.Set(key, cache.CopyResponse(resp))
		} else {
			client.Config.Logger.Printf("[ERROR] saving to cache error:%s, key:%s\n", err, key)
		}
	}
	return resp, nil
}

func (client *Client) authenticate() error {
	client.Config.Lock()
	defer client.Config.Unlock()
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" || utils.IsTokenExpired(client.Config.AuthToken.AccessToken) {
		if client.Config.ClientID == "" || client.Config.ClientSecret == "" {
			client.Config.Logger.Printf("[ERROR] No client credentials were provided. Please set %s, %s and %s environment variables.\n", ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID)
			return errors.New("no client credentials were provided")
		}
		client.Config.Logger.Printf("[TRACE] Getting access token for %s=%s\n", ZPA_CLIENT_ID, client.Config.ClientID)
		data := url.Values{}
		data.Set("client_id", client.Config.ClientID)
		data.Set("client_secret", client.Config.ClientSecret)
		authUrl := client.Config.BaseURL.String() + "/signin"
		if client.Config.Cloud == "DEV" {
			authUrl = devAuthUrl
		}
		req, err := http.NewRequest("POST", authUrl, strings.NewReader(data.Encode()))
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZPA_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZPA_CLIENT_ID, client.Config.ClientID, err)
		}

		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		if client.Config.UserAgent != "" {
			req.Header.Add("User-Agent", client.Config.UserAgent)
		}
		resp, err := client.Config.GetHTTPClient().Do(req)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZPA_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZPA_CLIENT_ID, client.Config.ClientID, err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZPA_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZPA_CLIENT_ID, client.Config.ClientID, err)
		}
		if resp.StatusCode >= 300 {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, got http status:%dn response body:%s\n", ZPA_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s", ZPA_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
		}
		var a AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZPA_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZPA_CLIENT_ID, client.Config.ClientID, err)
		}
		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	return nil
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	err := client.authenticate()
	if err != nil {
		return nil, err

	}
	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		return nil, err
	}
	reqID := uuid.NewString()
	start := time.Now()
	logger.LogRequest(client.Config.Logger, req, reqID, nil, true)
	resp, err := client.do(req, v, start, reqID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		err := client.authenticate()
		if err != nil {
			return nil, err
		}

		resp, err := client.do(req, v, start, reqID)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
		return resp, nil
	}
	return resp, err
}

func getMicrotenantIDFromBody(body interface{}) string {
	if body == nil {
		return ""
	}

	d, err := json.Marshal(body)
	if err != nil {
		return ""
	}
	dataMap := map[string]interface{}{}
	err = json.Unmarshal(d, &dataMap)
	if err != nil {
		return ""
	}
	if microTenantID, ok := dataMap["microtenantId"]; ok && microTenantID != nil && microTenantID != "" {
		return fmt.Sprintf("%v", microTenantID)
	}
	return ""
}

func getMicrotenantIDFromEnvVar(body interface{}) string {
	return os.Getenv("ZPA_MICROTENANT_ID")
}

func (client *Client) injectMicrotentantID(body interface{}, q url.Values) url.Values {
	if q.Has("microtenantId") && q.Get("microtenantId") != "" {
		return q
	}

	microTenantID := getMicrotenantIDFromBody(body)
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}

	microTenantID = getMicrotenantIDFromEnvVar(body)
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}
	return q
}

func (client *Client) getRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	// Join the path to the base-url
	u := *client.Config.BaseURL
	unescaped, err := url.PathUnescape(urlPath)
	if err != nil {
		return nil, err
	}

	// Set the encoded path data
	u.RawPath = u.Path + urlPath
	u.Path = u.Path + unescaped

	// Set the query parameters
	if options == nil {
		options = struct{}{}
	}
	// Set the query parameters

	q, err := query.Values(options)
	if err != nil {
		return nil, err
	}
	q = client.injectMicrotentantID(body, q)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s\n", ZPA_CLIENT_ID, client.Config.ClientID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", ZPA_CLIENT_ID, client.Config.ClientID)
	}
	req, err := client.getRequest(method, urlPath, options, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Config.AuthToken.AccessToken))
	req.Header.Add("Content-Type", "application/json")

	if client.Config.UserAgent != "" {
		req.Header.Add("User-Agent", client.Config.UserAgent)
	}

	return req, nil
}

func (client *Client) do(req *http.Request, v interface{}, start time.Time, reqID string) (*http.Response, error) {
	resp, err := client.Config.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	respData, err := io.ReadAll(resp.Body)
	if err == nil {
		resp.Body = io.NopCloser(bytes.NewBuffer(respData))
	}
	if err := checkErrorInResponse(resp, respData); err != nil {
		return resp, err
	}

	if v != nil {
		if err := decodeJSON(respData, v); err != nil {
			return resp, err
		}
	}
	logger.LogResponse(client.Config.Logger, resp, start, reqID)
	unescapeHTML(v)
	return resp, nil
}

func decodeJSON(respData []byte, v interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(respData)).Decode(&v)
}

func unescapeHTML(entity interface{}) {
	if entity == nil {
		return
	}
	data, err := json.Marshal(entity)
	if err != nil {
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return
	}
	for _, field := range []string{"name", "description"} {
		if v, ok := mapData[field]; ok && v != nil {
			str, ok := v.(string)
			if ok {
				mapData[field] = html.UnescapeString(html.UnescapeString(str))
			}
		}
	}
	data, err = json.Marshal(mapData)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, entity)
}
