package zia

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/cache"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/logger"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/google/uuid"
)

func (c *Client) do(req *http.Request, start time.Time, reqID string) (*http.Response, error) {
	key := cache.CreateCacheKey(req)
	if c.cacheEnabled {
		if req.Method != http.MethodGet {
			c.cache.Delete(key)
			c.cache.ClearAllKeysWithPrefix(strings.Split(key, "?")[0])
		}
		resp := c.cache.Get(key)
		inCache := resp != nil
		if c.freshCache {
			c.cache.Delete(key)
			inCache = false
			c.freshCache = false
		}
		if inCache {
			c.Logger.Printf("[INFO] served from cache, key:%s\n", key)
			return resp, nil
		}
	}

	// Ensure the session is valid before making the request
	err := c.checkSession()
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	logger.LogResponse(c.Logger, resp, start, reqID)
	if err != nil {
		return resp, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Fallback check for SESSION_NOT_VALID
	if resp.StatusCode == http.StatusUnauthorized || strings.Contains(string(body), "SESSION_NOT_VALID") {
		// Refresh session and retry
		err := c.refreshSession()
		if err != nil {
			return nil, err
		}
		req.Header.Set("JSessionID", c.session.JSessionID)
		resp, err = c.HTTPClient.Do(req)
		logger.LogResponse(c.Logger, resp, start, reqID)
		if err != nil {
			return resp, err
		}
	}

	if c.cacheEnabled && resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet {
		c.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
		c.cache.Set(key, cache.CopyResponse(resp))
	}

	return resp, nil
}

// Request ... // Needs to review this function.
func (c *Client) GenericRequest(baseUrl, endpoint, method string, body io.Reader, urlParams url.Values, contentType string) ([]byte, error) {
	if contentType == "" {
		contentType = contentTypeJSON
	}

	var req *http.Request
	var resp *http.Response
	var err error
	params := ""
	if urlParams != nil {
		params = urlParams.Encode()
	}
	if strings.Contains(endpoint, "?") && params != "" {
		endpoint += "&" + params
	} else if params != "" {
		endpoint += "?" + params
	}
	fullURL := fmt.Sprintf("%s%s", baseUrl, endpoint)
	isSandboxRequest := baseUrl == c.GetSandboxURL()
	req, err = http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	var otherHeaders map[string]string
	if !isSandboxRequest {
		err = c.checkSession()
		if err != nil {
			return nil, err
		}
		otherHeaders = map[string]string{"JSessionID": c.session.JSessionID}
	}
	reqID := uuid.New().String()
	start := time.Now()
	logger.LogRequest(c.Logger, req, reqID, otherHeaders, !isSandboxRequest)
	for retry := 1; retry <= 5; retry++ {
		resp, err = c.do(req, start, reqID)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode <= 299 {
			defer resp.Body.Close()
			break
		}

		resp.Body.Close()
		if resp.StatusCode > 299 && resp.StatusCode != http.StatusUnauthorized {
			return nil, errorx.CheckErrorInResponse(resp, fmt.Errorf("api responded with code: %d", resp.StatusCode))
		}
	}

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyResp, nil
}

// Request ... // Needs to review this function.
func (c *Client) Request(endpoint, method string, data []byte, contentType string) ([]byte, error) {
	return c.GenericRequest(c.URL, endpoint, method, bytes.NewReader(data), nil, contentType)
}

func (client *Client) WithFreshCache() {
	client.freshCache = true
}

// Create sends an HTTP POST request.
func (c *Client) Create(endpoint string, o interface{}) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to create with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to create with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		// Check if the response is an array of strings
		var stringArrayResponse []string
		if json.Unmarshal(resp, &stringArrayResponse) == nil {
			return stringArrayResponse, nil
		}

		// Otherwise, handle as usual
		responseObject := reflect.New(t).Interface()
		err = json.Unmarshal(resp, &responseObject)
		if err != nil {
			return nil, err
		}
		id := reflect.Indirect(reflect.ValueOf(responseObject)).FieldByName("ID")

		c.Logger.Printf("Created Object with ID %v", id)
		return responseObject, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

func (c *Client) CreateWithSlicePayload(endpoint string, slice interface{}) ([]byte, error) {
	if slice == nil {
		return nil, errors.New("tried to create with a nil payload not a Slice")
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("tried to create with a " + v.Kind().String() + " not a Slice")
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}
	if len(resp) > 0 {
		return resp, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

func (c *Client) UpdateWithSlicePayload(endpoint string, slice interface{}) ([]byte, error) {
	if slice == nil {
		return nil, errors.New("tried to update with a nil payload not a Slice")
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, errors.New("tried to update with a " + v.Kind().String() + " not a Slice")
	}

	data, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, "PUT", data, "application/json")
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateWithRawPayload sends an HTTP POST request with a raw string payload.
func (c *Client) CreateWithRawPayload(endpoint string, payload string) ([]byte, error) {
	if payload == "" {
		return nil, errors.New("tried to create with an empty string payload")
	}

	// Convert the string payload to []byte
	data := []byte(payload)

	// Send the raw string as a POST request
	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}

	// Handle the response
	if len(resp) > 0 {
		return resp, nil
	} else {
		// in case of 204 no content
		return nil, nil
	}
}

// Read ...
func (c *Client) Read(endpoint string, o interface{}) error {
	contentType := c.GetContentType()
	resp, err := c.Request(endpoint, "GET", nil, contentType)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}

	return nil
}

// Update ...
func (c *Client) UpdateWithPut(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PUT", "application/json")
}

// Update ...
func (c *Client) Update(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PATCH", "application/merge-patch+json")
}

// Update ...
func (c *Client) updateGeneric(endpoint string, o interface{}, method, contentType string) (interface{}, error) {
	if o == nil {
		return nil, errors.New("tried to update with a nil payload not a Struct")
	}
	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("tried to update with a " + t.Kind().String() + " not a Struct")
	}
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	resp, err := c.Request(endpoint, method, data, contentType)
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	return responseObject, err
}

// Delete ...
func (c *Client) Delete(endpoint string) error {
	_, err := c.Request(endpoint, "DELETE", nil, "application/json")
	if err != nil {
		return err
	}
	return nil
}

// BulkDelete sends an HTTP POST request for bulk deletion and expects a 204 No Content response.
func (c *Client) BulkDelete(endpoint string, payload interface{}) (*http.Response, error) {
	if payload == nil {
		return nil, errors.New("tried to delete with a nil payload, expected a struct")
	}

	// Marshal the payload into JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Send the POST request
	resp, err := c.Request(endpoint, "POST", data, "application/json")
	if err != nil {
		return nil, err
	}

	// Check the status code (204 No Content expected)
	if len(resp) == 0 {
		c.Logger.Printf("[DEBUG] Bulk delete successful with 204 No Content")
		return &http.Response{StatusCode: 204}, nil
	}

	// If the response is not empty, this might indicate an error or unexpected behavior
	return &http.Response{StatusCode: 200}, fmt.Errorf("unexpected response: %s", string(resp))
}
