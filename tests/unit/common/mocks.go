// Package common provides shared test utilities and mocks for unit testing
package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// =============================================================================
// Mock HTTP Server
// =============================================================================

// MockResponse represents a configured mock response
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
	Delay      time.Duration
}

// MockHandler is a configurable HTTP handler for testing
type MockHandler struct {
	mu        sync.Mutex
	Responses map[string]map[string]MockResponse // path -> method -> response
	Requests  []RecordedRequest
	// CallCount tracks how many times each path/method combination was called
	CallCount map[string]int
}

// RecordedRequest stores information about a received request
type RecordedRequest struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	Body    []byte
}

// NewMockHandler creates a new mock handler
func NewMockHandler() *MockHandler {
	return &MockHandler{
		Responses: make(map[string]map[string]MockResponse),
		Requests:  []RecordedRequest{},
		CallCount: make(map[string]int),
	}
}

// On registers a mock response for a path and method
func (m *MockHandler) On(method, path string, response MockResponse) *MockHandler {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Responses[path] == nil {
		m.Responses[path] = make(map[string]MockResponse)
	}
	m.Responses[path][method] = response
	return m
}

// ServeHTTP implements http.Handler
func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()

	// Record the request
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	m.Requests = append(m.Requests, RecordedRequest{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.RawQuery,
		Headers: r.Header,
		Body:    body,
	})

	// Track call count
	key := r.Method + ":" + r.URL.Path
	m.CallCount[key]++

	// Find matching response
	pathResponses, ok := m.Responses[r.URL.Path]
	if !ok {
		// Try prefix matching for paths with IDs
		for path, responses := range m.Responses {
			if strings.HasPrefix(r.URL.Path, path) || matchPath(path, r.URL.Path) {
				pathResponses = responses
				ok = true
				break
			}
		}
	}

	m.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "no mock configured for path: ` + r.URL.Path + `"}`))
		return
	}

	response, ok := pathResponses[r.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Apply delay if configured
	if response.Delay > 0 {
		time.Sleep(response.Delay)
	}

	// Set headers
	for k, v := range response.Headers {
		w.Header().Set(k, v)
	}

	// Set default content type
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(response.StatusCode)

	// Write body
	if response.Body != nil {
		switch v := response.Body.(type) {
		case string:
			w.Write([]byte(v))
		case []byte:
			w.Write(v)
		default:
			json.NewEncoder(w).Encode(v)
		}
	}
}

// matchPath checks if a pattern matches a path (simple wildcard support)
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if part == "*" || part == "{id}" || strings.HasPrefix(part, ":") {
			continue
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}

// =============================================================================
// Test Server Factory
// =============================================================================

// TestServer wraps httptest.Server with additional utilities
type TestServer struct {
	*httptest.Server
	Handler *MockHandler
}

// NewTestServer creates a new test server with the mock handler
func NewTestServer() *TestServer {
	handler := NewMockHandler()
	server := httptest.NewServer(handler)
	return &TestServer{
		Server:  server,
		Handler: handler,
	}
}

// On is a convenience method to add mock responses
func (ts *TestServer) On(method, path string, response MockResponse) *TestServer {
	ts.Handler.On(method, path, response)
	return ts
}

// AssertRequestCount verifies the number of requests made
func (ts *TestServer) AssertRequestCount(t *testing.T, expected int) {
	require.Equal(t, expected, len(ts.Handler.Requests), "unexpected request count")
}

// GetCallCount returns the number of times a specific path/method was called
func (ts *TestServer) GetCallCount(method, path string) int {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()
	return ts.Handler.CallCount[method+":"+path]
}

// LastRequest returns the most recent request
func (ts *TestServer) LastRequest() *RecordedRequest {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()

	if len(ts.Handler.Requests) == 0 {
		return nil
	}
	return &ts.Handler.Requests[len(ts.Handler.Requests)-1]
}

// =============================================================================
// Response Builders
// =============================================================================

// SuccessResponse creates a 200 OK response
func SuccessResponse(body interface{}) MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

// SuccessResponseWithStatus creates a response with a custom status code
func SuccessResponseWithStatus(statusCode int, body interface{}) MockResponse {
	return MockResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}

// CreatedResponse creates a 201 Created response
func CreatedResponse(body interface{}) MockResponse {
	return MockResponse{
		StatusCode: http.StatusCreated,
		Body:       body,
	}
}

// NoContentResponse creates a 204 No Content response
func NoContentResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusNoContent,
	}
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	}
}

// UnauthorizedResponse creates a 401 Unauthorized response
func UnauthorizedResponse(body string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       body,
	}
}

// SessionInvalidResponse creates a 401 SESSION_NOT_VALID response
func SessionInvalidResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       `{"code": "SESSION_NOT_VALID", "message": "Session is not valid"}`,
	}
}

// ResourceAccessBlockedResponse creates a 401 Resource Access Blocked response
func ResourceAccessBlockedResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       `{"message": "Resource Access Blocked"}`,
	}
}

// TooManyRequestsResponse creates a 429 Too Many Requests response
func TooManyRequestsResponse(retryAfterSeconds int) MockResponse {
	return MockResponse{
		StatusCode: http.StatusTooManyRequests,
		Body:       `{"message": "Rate limit exceeded"}`,
		Headers: map[string]string{
			"Retry-After": string(rune('0' + retryAfterSeconds)),
		},
	}
}

// TooManyRequestsResponseWithHeader creates a 429 response with custom Retry-After
func TooManyRequestsResponseWithHeader(retryAfter string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusTooManyRequests,
		Body:       `{"message": "Rate limit exceeded"}`,
		Headers: map[string]string{
			"Retry-After": retryAfter,
		},
	}
}

// ConflictResponse creates a 409 Conflict response
func ConflictResponse(body string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusConflict,
		Body:       body,
	}
}

// EditLockResponse creates a 409 EDIT_LOCK_NOT_AVAILABLE response
func EditLockResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusConflict,
		Body:       `{"code": "EDIT_LOCK_NOT_AVAILABLE", "message": "Edit lock not available"}`,
	}
}

// ServerErrorResponse creates a 500 Internal Server Error response
func ServerErrorResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"message": "Internal server error"}`,
	}
}

// =============================================================================
// HTTP Response Helpers
// =============================================================================

// CreateMockHTTPResponse creates an http.Response for testing
func CreateMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// CreateMockHTTPResponseWithHeaders creates an http.Response with custom headers
func CreateMockHTTPResponseWithHeaders(statusCode int, body string, headers map[string]string) *http.Response {
	resp := CreateMockHTTPResponse(statusCode, body)
	for k, v := range headers {
		resp.Header.Set(k, v)
	}
	return resp
}

// =============================================================================
// OAuth Token Mock
// =============================================================================

// MockOAuthResponse returns a valid OAuth token response
func MockOAuthResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"access_token": "mock-access-token-12345",
			"token_type":   "Bearer",
			"expires_in":   "3600",
		},
	}
}

// MockZPAOAuthResponse returns a valid ZPA OAuth token response
func MockZPAOAuthResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body: map[string]interface{}{
			"access_token": "mock-zpa-access-token-12345",
			"token_type":   "Bearer",
			"expires_in":   "3600",
		},
	}
}

// RawResponse creates a response with raw bytes and custom content type
func RawResponse(body []byte, statusCode int, headers map[string]string) MockResponse {
	return MockResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	}
}

// CSVResponse creates a CSV response
func CSVResponse(csvData string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body:       csvData,
		Headers: map[string]string{
			"Content-Type": "text/csv",
		},
	}
}

// =============================================================================
// ZPA Specific Helpers
// =============================================================================

// MockSegmentGroupResponse returns a mock segment group response
func MockSegmentGroupResponse(id, name string) map[string]interface{} {
	return map[string]interface{}{
		"id":                  id,
		"name":                name,
		"description":         "Test segment group",
		"enabled":             true,
		"configSpace":         "DEFAULT",
		"policyMigrated":      false,
		"tcpKeepAliveEnabled": "0",
		"applications":        []interface{}{},
	}
}

// MockSegmentGroupListResponse returns a mock list of segment groups
func MockSegmentGroupListResponse(groups ...map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"list":       groups,
		"totalPages": 1,
	}
}

