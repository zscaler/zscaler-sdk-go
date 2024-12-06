package zscaler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// NewRequestDo for ZCC with OAuth2 authentication and centralized request handling.
// This function is consistent with the ZPA request handler.
func (client *Client) NewZccRequestDo(ctx context.Context, method, endpoint string, options, body, v interface{}) (*http.Response, error) {
	if client.oauth2Credentials.UseLegacyClient {
		if client.oauth2Credentials.LegacyClient == nil || client.oauth2Credentials.LegacyClient.ZccClient == nil {
			return nil, errLegacyClientNotSet
		}
		return client.oauth2Credentials.LegacyClient.ZccClient.NewRequestDo(method, removeOneApiEndpointPrefix(endpoint), options, body, v)
	}
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

	// Make the request and get the response
	respBody, _, _, err := client.ExecuteRequest(ctx, method, endpoint, bodyReader, nil, contentTypeJSON)
	if err != nil {
		return nil, err
	}

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBuffer(respBody)),
	}

	// Handle response parsing
	if v != nil {
		if err := decodeJSON(respBody, v); err != nil {
			return resp, err
		}
	}
	unescapeHTML(v)

	return resp, nil
}
