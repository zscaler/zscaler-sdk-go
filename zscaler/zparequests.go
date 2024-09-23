package zscaler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
)

func (client *Client) NewRequestDo(method, endpoint string, options, body, v interface{}) (*http.Response, error) {
	// Call the custom request handler
	// Handle query parameters from options and any additional logic
	if options == nil {
		options = struct{}{}
	}
	var params string
	if options != nil {
		switch opt := options.(type) {
		case url.Values:
			params = opt.Encode()
		default:
			q, err := query.Values(options)
			if err != nil {
				return nil, err
			}
			params = q.Encode()
		}
	}

	if strings.Contains(endpoint, "?") && params != "" {
		endpoint += "&" + params
	} else if params != "" {
		endpoint += "?" + params
	}

	parts := strings.Split(endpoint, "?")
	path := parts[0]
	query := ""
	if len(parts) > 1 {
		query = parts[1]
	}
	q, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	q = client.injectMicrotentantID(body, q)
	query = q.Encode()
	endpoint = path
	if query != "" {
		endpoint += "?" + query
	}
	// Use ExecuteRequest to handle the request
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(bodyBytes)
	}

	// Capture the three return values from ExecuteRequest
	respBody, _, err := client.ExecuteRequest(method, endpoint, bodyReader, nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(respBody)),
	}

	if v != nil {
		if err := decodeJSON(respBody, v); err != nil {
			return resp, err
		}
	}
	unescapeHTML(v)

	return resp, nil
}

func (c *Client) GetCustomerID() string {
	return c.oauth2Credentials.Zscaler.Client.CustomerID
}

func (client *Client) GetFullPath(endpoint string) (string, error) {
	customerID := client.GetCustomerID()
	if customerID == "" {
		return "", fmt.Errorf("CustomerID is not set")
	}
	// Construct the full path with mgmtConfig and CustomerID
	return fmt.Sprintf("%s%s%s", mgmtConfig, customerID, endpoint), nil
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

	microTenantID = client.oauth2Credentials.Zscaler.Client.MicrotenantID
	if microTenantID != "" {
		q.Add("microtenantId", microTenantID)
		return q
	}

	return q
}
