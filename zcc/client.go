package zcc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/zscaler/zscaler-sdk-go/logger"
)

type Client struct {
	Config *Config
}

// NewClient returns a new client for the specified apiKey.
func NewClient(config *Config) (c *Client) {
	if config == nil {
		config, _ = NewConfig("", "", "", "")
	}
	c = &Client{Config: config}
	return
}

func (client *Client) NewRequestDo(method, url string, options, body, v interface{}) (*http.Response, error) {
	return client.newRequestDoCustom(method, url, options, body, v)
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	client.Config.Lock()
	defer client.Config.Unlock()
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		if client.Config.ClientID == "" || client.Config.ClientSecret == "" {
			client.Config.Logger.Printf("[ERROR] No client credentials were provided. Please set %s and %s environment variables.\n", ZCC_CLIENT_ID, ZCC_CLIENT_SECRET)
			return nil, errors.New("no client credentials were provided")
		}
		client.Config.Logger.Printf("[TRACE] Getting access token for %s=%s\n", ZCC_CLIENT_ID, client.Config.ClientID)
		authReq := AuthRequest{}
		authReq.APIKey = client.Config.ClientID
		authReq.SecretKey = client.Config.ClientSecret
		data, err := json.Marshal(authReq)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		req, err := http.NewRequest("POST", client.Config.BaseURL.String()+"/auth/v1/login", bytes.NewBuffer(data))
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}

		req.Header.Add("Content-Type", "application/json")
		if client.Config.UserAgent != "" {
			req.Header.Add("User-Agent", client.Config.UserAgent)
		}
		resp, err := client.Config.GetHTTPClient().Do(req)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		if resp.StatusCode >= 300 {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, got http status:%dn response body:%s\n", ZCC_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s", ZCC_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
		}
		var a AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		return nil, err
	}
	logger.LogRequest(client.Config.Logger, req)
	return client.do(req, v)
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s\n", ZCC_CLIENT_ID, client.Config.ClientID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", ZCC_CLIENT_ID, client.Config.ClientID)
	}
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
	if options != nil {
		q, err := query.Values(options)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("auth-token", client.Config.AuthToken.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	if client.Config.UserAgent != "" {
		req.Header.Add("User-Agent", client.Config.UserAgent)
	}

	return req, nil
}

func (client *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := client.Config.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	if err := checkErrorInResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err := decodeJSON(resp, v); err != nil {
			return resp, err
		}
	}
	logger.LogResponse(client.Config.Logger, resp)
	unescapeHTML(v)
	return resp, nil
}

func decodeJSON(res *http.Response, v interface{}) error {
	return json.NewDecoder(res.Body).Decode(&v)
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
