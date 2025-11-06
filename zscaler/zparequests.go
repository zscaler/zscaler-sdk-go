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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/common"
)

func (client *Client) NewRequestDo(ctx context.Context, method, endpoint string, options, body, v interface{}) (*http.Response, error) {
	if client.oauth2Credentials.UseLegacyClient {
		if client.oauth2Credentials.LegacyClient == nil || client.oauth2Credentials.LegacyClient.ZpaClient == nil {
			return nil, errLegacyClientNotSet
		}
		return client.oauth2Credentials.LegacyClient.ZpaClient.NewRequestDo(method, removeOneApiEndpointPrefix(endpoint), options, body, v)
	}
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
	
	// Parse query string to work with parameters
	q, err := url.ParseQuery(query)
	if err != nil {
		return nil, err
	}
	
	// Inject microtenant ID if needed
	q = common.InjectMicrotentantID(body, q, client.oauth2Credentials.Zscaler.Client.MicrotenantID)
	
	// For ZPA endpoints, use custom encoding that preserves %20 for spaces
	// Standard url.Values.Encode() uses + for spaces, but ZPA API requires %20
	isZPAEndpoint := strings.Contains(endpoint, "/zpa/") || strings.Contains(endpoint, "/mgmtconfig/")
	if isZPAEndpoint {
		query = encodeQueryWithSpacesAsPercent20(q)
	} else {
		query = q.Encode()
	}
	
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
	respBody, _, _, err := client.ExecuteRequest(ctx, method, endpoint, bodyReader, nil, contentTypeJSON)
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

// encodeQueryWithSpacesAsPercent20 encodes url.Values similar to url.Values.Encode()
// but uses %20 for spaces instead of + to match ZPA API requirements
func encodeQueryWithSpacesAsPercent20(v url.Values) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	// Sort keys for consistent output
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			// Use QueryEscape which encodes spaces as %20
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}

func (c *Client) GetCustomerID() string {
	if c.oauth2Credentials.UseLegacyClient && c.oauth2Credentials.LegacyClient != nil && c.oauth2Credentials.LegacyClient.ZpaClient != nil && c.oauth2Credentials.LegacyClient.ZpaClient.Config.ZPA.Client.ZPACustomerID != "" {
		return c.oauth2Credentials.LegacyClient.ZpaClient.Config.ZPA.Client.ZPACustomerID
	}
	return c.oauth2Credentials.Zscaler.Client.CustomerID
}
