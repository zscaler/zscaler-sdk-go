package zdx

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/zscaler/zscaler-sdk-go/logger"
)

type Client struct {
	Config *Config
}

type AuthRequest struct {
	APIKeyID     string `json:"key_id"`
	APIKeySecret string `json:"key_secret"`
	Timestamp    int64  `json:"timestamp"`
}

// NewClient returns a new client for the specified apiKey.
func NewClient(config *Config) (c *Client) {
	if config == nil {
		config, _ = NewConfig("", "", "")
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
		if client.Config.APIKeyID == "" || client.Config.APISecret == "" {
			client.Config.Logger.Printf("[ERROR] No client credentials were provided. Please set %s, %s environment variables.\n", ZDX_API_KEY_ID, ZDX_API_SECRET)
			return nil, errors.New("no client credentials were provided")
		}
		client.Config.Logger.Printf("[TRACE] Getting access token for %s=%s\n", ZDX_API_KEY_ID, client.Config.APIKeyID)
		currTimestamp := time.Now().Unix()
		authReq := AuthRequest{
			Timestamp:    currTimestamp,
			APIKeyID:     client.Config.APIKeyID,
			APIKeySecret: generateHash(client.Config.APISecret, currTimestamp),
		}
		data, _ := json.Marshal(authReq)
		url := client.Config.BaseURL.String() + "/v1/oauth/token"
		req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
		}

		req.Header.Add("Content-Type", "application/json")
		if client.Config.UserAgent != "" {
			req.Header.Add("User-Agent", client.Config.UserAgent)
		}
		resp, err := client.Config.GetHTTPClient().Do(req)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
		}
		if resp.StatusCode >= 300 {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s, url:%s\n", ZDX_API_KEY_ID, client.Config.APIKeyID, resp.StatusCode, respBody, url)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s, url:%s", ZDX_API_KEY_ID, client.Config.APIKeyID, resp.StatusCode, respBody, url)
		}
		var a AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
			return nil, fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZDX_API_KEY_ID, client.Config.APIKeyID, err)
		}
		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		return nil, err
	}
	reqID := uuid.NewString()
	start := time.Now()
	logger.LogRequest(client.Config.Logger, req, reqID)
	return client.do(req, v, start, reqID)
}

func generateHash(apiSecret string, currTimestamp int64) string {
	currTimestampStr := strconv.FormatInt(currTimestamp, 10)
	hash := sha256.New()
	hash.Write([]byte(apiSecret + ":" + currTimestampStr))
	return hex.EncodeToString(hash.Sum(nil))
}

// Generating the Http request
func (client *Client) newRequest(method, urlPath string, options, body interface{}) (*http.Request, error) {
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" {
		client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s\n", ZDX_API_KEY_ID, client.Config.APIKeyID)
		return nil, fmt.Errorf("failed to signin the user %s=%s", ZDX_API_KEY_ID, client.Config.APIKeyID)
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

	if err := checkErrorInResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		if err := decodeJSON(resp, v); err != nil {
			return resp, err
		}
	}
	logger.LogResponse(client.Config.Logger, resp, start, reqID)
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
	json.Unmarshal(data, entity)
}
