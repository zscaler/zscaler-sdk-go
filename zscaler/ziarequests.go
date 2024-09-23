package zscaler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

// Create sends a POST request to create an object.
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

	// Adjusting to handle the extra return value from ExecuteRequest
	resp, _, err := c.ExecuteRequest("POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
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
		return nil, nil // for 204 No Content
	}
}

// Read ...
func (c *Client) Read(endpoint string, o interface{}) error {
	resp, _, err := c.ExecuteRequest("GET", endpoint, nil, nil, contentTypeJSON)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp, o)
	if err != nil {
		return err
	}
	return nil
}

// UpdateWithPut sends an update (PUT request) with the given object.
func (c *Client) UpdateWithPut(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PUT", contentTypeJSON)
}

// Update sends an update (PATCH request) with the given object.
func (c *Client) Update(endpoint string, o interface{}) (interface{}, error) {
	return c.updateGeneric(endpoint, o, "PATCH", "application/merge-patch+json")
}

// General method to update an object using the specified HTTP method.
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

	resp, _, err := c.ExecuteRequest(method, endpoint, bytes.NewReader(data), nil, contentType)
	if err != nil {
		return nil, err
	}

	responseObject := reflect.New(t).Interface()
	err = json.Unmarshal(resp, &responseObject)
	return responseObject, err
}

// Delete sends a DELETE request to the specified endpoint.
func (c *Client) Delete(endpoint string) error {
	_, _, err := c.ExecuteRequest("DELETE", endpoint, nil, nil, contentTypeJSON)
	if err != nil {
		return err
	}
	return nil
}

// BulkDelete sends a POST request for bulk deletion.
func (c *Client) BulkDelete(endpoint string, payload interface{}) (*http.Response, error) {
	if payload == nil {
		return nil, errors.New("tried to delete with a nil payload, expected a struct")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, _, err := c.ExecuteRequest("POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		c.Logger.Printf("[DEBUG] Bulk delete successful with 204 No Content")
		return &http.Response{StatusCode: 204}, nil
	}

	return &http.Response{StatusCode: 200}, fmt.Errorf("unexpected response: %s", string(resp))
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

	// Explicitly set the contentType as "application/json"
	resp, _, err := c.ExecuteRequest("POST", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
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

	// Explicitly set the contentType as "application/json"
	resp, _, err := c.ExecuteRequest("PUT", endpoint, bytes.NewReader(data), nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	return resp, nil
}