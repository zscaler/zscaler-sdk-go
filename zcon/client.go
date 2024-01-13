package zcon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zscaler/zscaler-sdk-go/v2/cache"
	"github.com/zscaler/zscaler-sdk-go/v2/logger"
)

func (c *Client) do(req *http.Request, start time.Time, reqID string) (*http.Response, error) {
	key := cache.CreateCacheKey(req)
	if c.cacheEnabled {
		if req.Method != http.MethodGet {
			// this will allow to remove resource from cache when PUT/DELETE/PATCH requests are called, which modifies the resource
			c.cache.Delete(key)
			// to avoid resources that GET url is not the same as DELETE/PUT/PATCH url, because of different query params.
			// example delete app segment has key url/<id>?forceDelete=true but GET has url/<id>, in this case we clean the whole cache entries with key prefix url/<id>
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

	resp, err := c.HTTPClient.Do(req)
	logger.LogResponse(c.Logger, resp, start, reqID)
	if err != nil {
		return resp, err
	}
	if c.cacheEnabled && resp.StatusCode >= 200 && resp.StatusCode <= 299 && req.Method == http.MethodGet {
		c.Logger.Printf("[INFO] saving to cache, key:%s\n", key)
		c.cache.Set(key, cache.CopyResponse(resp))
	}
	return resp, nil
}

// Request ... // Needs to review this function.
func (c *Client) Request(endpoint, method string, data []byte, contentType string) ([]byte, error) {
	if contentType == "" {
		contentType = contentTypeJSON
	}

	var req *http.Request
	var resp *http.Response
	var err error

	req, err = http.NewRequest(method, c.URL+endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.UserAgent != "" {
		req.Header.Add("User-Agent", c.UserAgent)
	}
	err = c.checkSession()
	if err != nil {
		return nil, err
	}
	reqID := uuid.New().String()
	start := time.Now()
	logRequestBody := true // or false, based on your requirements

	logger.LogRequest(c.Logger, req, reqID, map[string]string{"JSessionID": c.session.JSessionID}, logRequestBody)

	for retry := 1; retry <= 5; retry++ {
		err = c.checkSession()
		if err != nil {
			return nil, err
		}

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
			return nil, checkErrorInResponse(resp, fmt.Errorf("api responded with code: %d", resp.StatusCode))
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *Client) WithFreshCache() {
	client.freshCache = true
}

// Create send HTTP Post request.
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
