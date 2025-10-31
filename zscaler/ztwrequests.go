package zscaler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Create sends a POST request to create an object.
func (c *Client) CreateResource(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.Create(ctx, removeOneApiEndpointPrefix(endpoint), o)
	}

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

	// Adjusting to handle the extra return value from ExecuteRequest
	respBody, response, _, err := c.ExecuteRequest(ctx, "POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	if len(respBody) > 0 && strings.EqualFold(response.Header.Get("Content-Type"), "application/json") {
		responseObject := reflect.New(t).Interface()
		err = json.Unmarshal(respBody, &responseObject)
		if err != nil {
			return nil, err
		}
		id := reflect.Indirect(reflect.ValueOf(responseObject)).FieldByName("ID")
		c.oauth2Credentials.Logger.Printf("Created Object with ID %v", id)
		return responseObject, nil
	} else {
		if len(respBody) > 0 {
			response.Body = io.NopCloser(bytes.NewReader(respBody))
		}
		return response, nil
	}
}

// Read ...
func (c *Client) ReadResource(ctx context.Context, endpoint string, o interface{}) error {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.Read(ctx, removeOneApiEndpointPrefix(endpoint), o)
	}

	resp, _, _, err := c.ExecuteRequest(ctx, "GET", endpoint, nil, nil, contentTypeJSON)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}
	return nil
}

// ReadTextResource reads a resource that returns plain text instead of JSON
func (c *Client) ReadTextResource(ctx context.Context, endpoint string, o interface{}) error {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.Read(ctx, removeOneApiEndpointPrefix(endpoint), o)
	}

	resp, _, _, err := c.ExecuteRequest(ctx, "GET", endpoint, nil, nil, contentTypeJSON)
	if err != nil {
		return err
	}

	// Handle string pointer - convert response to string and assign to the pointer
	if strPtr, ok := o.(*string); ok {
		*strPtr = strings.TrimSpace(string(resp))
		return nil
	}

	// Fallback to JSON unmarshaling for other types
	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}
	return nil
}

// UpdateWithPut sends an update (PUT request) with the given object.
func (c *Client) UpdateWithPutResource(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.UpdateWithPut(ctx, removeOneApiEndpointPrefix(endpoint), o)
	}
	return c.updateGeneric(ctx, endpoint, o, "PUT", contentTypeJSON)
}

// Update sends an update (PATCH request) with the given object.
func (c *Client) UpdateResource(ctx context.Context, endpoint string, o interface{}) (interface{}, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.Update(ctx, removeOneApiEndpointPrefix(endpoint), o)
	}
	return c.updateGeneric(ctx, endpoint, o, "PATCH", "application/merge-patch+json")
}

// General method to update an object using the specified HTTP method.
func (c *Client) updateGenericResource(ctx context.Context, endpoint string, o interface{}, method, contentType string) (interface{}, error) {
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

	resp, _, _, err := c.ExecuteRequest(ctx, method, endpoint, bytes.NewReader(data), nil, contentType)
	if err != nil {
		return nil, err
	}

	// Check for an empty response body (e.g., 204 No Content)
	if len(resp) == 0 {
		return nil, nil // Return nil for responseObject and no error
	}

	// Unmarshal response body into the provided struct type
	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return responseObject, nil
}

// Delete sends a DELETE request to the specified endpoint.
func (c *Client) DeleteResource(ctx context.Context, endpoint string) error {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.Delete(ctx, removeOneApiEndpointPrefix(endpoint))
	}
	_, _, _, err := c.ExecuteRequest(ctx, "DELETE", endpoint, nil, nil, contentTypeJSON)
	if err != nil {
		return err
	}
	return nil
}

// BulkDelete sends a POST request for bulk deletion.
func (c *Client) BulkDeleteResource(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.BulkDelete(ctx, removeOneApiEndpointPrefix(endpoint), payload)
	}

	if payload == nil {
		return nil, errors.New("tried to delete with a nil payload, expected a struct")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, _, _, err := c.ExecuteRequest(ctx, "POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		c.oauth2Credentials.Logger.Printf("[DEBUG] Bulk delete successful with 204 No Content")
		return &http.Response{StatusCode: 204}, nil
	}

	return &http.Response{StatusCode: 200}, fmt.Errorf("unexpected response: %s", string(resp))
}

func (c *Client) CreateWithSlicePayloadResource(ctx context.Context, endpoint string, slice interface{}) ([]byte, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.CreateWithSlicePayload(ctx, removeOneApiEndpointPrefix(endpoint), slice)
	}

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

	// Explicitly set the contentType as "application/json"
	resp, _, _, err := c.ExecuteRequest(ctx, "POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
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

func (c *Client) UpdateWithSlicePayloadResource(ctx context.Context, endpoint string, slice interface{}) ([]byte, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.UpdateWithSlicePayload(ctx, removeOneApiEndpointPrefix(endpoint), slice)
	}

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

	// Explicitly set the contentType as "application/json"
	resp, _, _, err := c.ExecuteRequest(ctx, "PUT", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CreateWithRawPayload sends an HTTP POST request with a raw string payload.
func (c *Client) CreateWithRawPayloadResource(ctx context.Context, endpoint string, payload string) ([]byte, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}
		return c.oauth2Credentials.LegacyClient.ZtwClient.CreateWithRawPayload(ctx, removeOneApiEndpointPrefix(endpoint), payload)
	}

	if payload == "" {
		return nil, errors.New("tried to create with an empty string payload")
	}

	// Convert the string payload to []byte
	data := []byte(payload)

	// Send the raw string as a POST request
	resp, _, _, err := c.ExecuteRequest(ctx, "POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
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

// CreateWithNoContent handles POST requests that return a 204 No Content response.
func (c *Client) CreateWithNoContentResource(ctx context.Context, endpoint string, o interface{}) (*http.Response, error) {
	if c.oauth2Credentials.UseLegacyClient {
		if c.oauth2Credentials.LegacyClient == nil || c.oauth2Credentials.LegacyClient.ZtwClient == nil {
			return nil, errLegacyClientNotSet
		}

		// Type assertion for legacy client's response
		resp, err := c.oauth2Credentials.LegacyClient.ZtwClient.CreateWithNoContent(ctx, removeOneApiEndpointPrefix(endpoint), o)
		if err != nil {
			return nil, err
		}

		// Ensure the returned value is of type *http.Response
		httpResp, ok := resp.(*http.Response)
		if !ok {
			return nil, fmt.Errorf("unexpected response type: %T, expected *http.Response", resp)
		}

		return httpResp, nil
	}

	// Validate the payload
	if o == nil {
		return nil, errors.New("tried to create with a nil payload, expected a Struct")
	}

	t := reflect.TypeOf(o)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("tried to create with a %s, expected a Struct", t.Kind().String())
	}

	// Marshal the payload
	data, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	// Use the existing ExecuteRequest method
	_, response, _, err := c.ExecuteRequest(ctx, "POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	// Handle the 204 No Content scenario
	if response.StatusCode == http.StatusNoContent {
		c.oauth2Credentials.Logger.Printf("Successfully created object at endpoint: %s (204 No Content)", endpoint)
		return response, nil
	}

	// Check for unexpected response codes
	// if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
	// 	return response, fmt.Errorf("unexpected response code: %d", response.StatusCode)
	// }

	c.oauth2Credentials.Logger.Printf("Successfully created object at endpoint: %s", endpoint)
	return response, nil
}
