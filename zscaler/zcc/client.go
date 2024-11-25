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
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/uuid"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/logger"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/utils"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/errorx"
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
	client.Config.Logger.Printf("[DEBUG] Creating new request: method=%s, url=%s", method, url)
	return client.newRequestDoCustom(method, url, options, body, v)
}

func (client *Client) authenticate() error {
	client.Config.Lock()
	defer client.Config.Unlock()
	client.Config.Logger.Printf("[DEBUG] Authenticating client: clientID=%s", client.Config.ClientID)
	if client.Config.AuthToken == nil || client.Config.AuthToken.AccessToken == "" || utils.IsTokenExpired(client.Config.AuthToken.AccessToken) {
		client.Config.Logger.Printf("[DEBUG] No valid auth token found, performing authentication")
		if client.Config.ClientID == "" || client.Config.ClientSecret == "" {
			client.Config.Logger.Printf("[ERROR] No client credentials were provided. Please set %s and %s environment variables.\n", ZCC_CLIENT_ID, ZCC_CLIENT_SECRET)
			return errors.New("no client credentials were provided")
		}
		client.Config.Logger.Printf("[TRACE] Getting access token for %s=%s\n", ZCC_CLIENT_ID, client.Config.ClientID)
		authReq := AuthRequest{}
		authReq.APIKey = client.Config.ClientID
		authReq.SecretKey = client.Config.ClientSecret
		data, err := json.Marshal(authReq)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		req, err := http.NewRequest("POST", client.Config.BaseURL.String()+"/auth/v1/login", bytes.NewBuffer(data))
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}

		req.Header.Add("Content-Type", "application/json")
		if client.Config.UserAgent != "" {
			req.Header.Add("User-Agent", client.Config.UserAgent)
		}
		resp, err := client.Config.GetHTTPClient().Do(req)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		if resp.StatusCode >= 300 {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, got http status:%d response body:%s\n", ZCC_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, got http status:%d, response body:%s", ZCC_CLIENT_ID, client.Config.ClientID, resp.StatusCode, respBody)
		}
		var a AuthToken
		err = json.Unmarshal(respBody, &a)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to signin the user %s=%s, err: %v\n", ZCC_CLIENT_ID, client.Config.ClientID, err)
			return fmt.Errorf("[ERROR] Failed to signin the user %s=%s, err: %v", ZCC_CLIENT_ID, client.Config.ClientID, err)
		}
		client.Config.Logger.Printf("[DEBUG] Authentication successful, token received")
		// we need keep auth token for future http request
		client.Config.AuthToken = &a
	}
	return nil
}

func (client *Client) newRequestDoCustom(method, urlStr string, options, body, v interface{}) (*http.Response, error) {
	client.Config.Logger.Printf("[DEBUG] newRequestDoCustom called with method=%s, urlStr=%s", method, urlStr)
	err := client.authenticate()
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Authentication failed: %v", err)
		return nil, err
	}

	req, err := client.newRequest(method, urlStr, options, body)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Failed to create new request: %v", err)
		return nil, err
	}
	start := time.Now()
	reqID := uuid.NewString()
	logger.LogRequest(client.Config.Logger, req, reqID, nil, true)
	resp, err := client.do(req, v, start, reqID)
	if err != nil {
		client.Config.Logger.Printf("[ERROR] Request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		client.Config.Logger.Printf("[WARN] Unauthorized or forbidden response, retrying authentication")
		err := client.authenticate()
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Re-authentication failed: %v", err)
			return nil, err
		}

		resp, err := client.do(req, v, start, reqID)
		if err != nil {
			client.Config.Logger.Printf("[ERROR] Request failed after re-authentication: %v", err)
			return nil, err
		}
		resp.Body.Close()
		return resp, nil
	}
	return resp, err
}

// Generating the Http request.
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
	u.Path += unescaped

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

func (client *Client) do(req *http.Request, v interface{}, start time.Time, reqID string) (*http.Response, error) {
	resp, err := client.Config.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}

	if err := errorx.CheckErrorInResponse(resp, err); err != nil {
		return resp, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	// Log the raw response body for debugging
	client.Config.Logger.Printf("[DEBUG] Raw response body: %s", string(bodyBytes))

	if v != nil {
		// Decode JSON from the raw response body
		if err := json.Unmarshal(bodyBytes, v); err != nil {
			client.Config.Logger.Printf("[ERROR] Failed to parse JSON response: %v", err)
			return resp, err
		}
	}
	logger.LogResponse(client.Config.Logger, resp, start, reqID)
	unescapeHTML(v)
	return resp, nil
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
