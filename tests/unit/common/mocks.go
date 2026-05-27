// Package common provides shared test utilities and mocks for unit testing.
//
// The mock server in this file is the single source of truth used by every
// OneAPI-routed test (ZPA, ZIA, ZCC, ZDX, ZTW, ZID, ZWA). It works in
// concert with CreateTestService in testutils.go: every HTTP request the
// SDK would issue is redirected to a httptest.Server backed by MockHandler.
package common

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// =============================================================================
// Mock HTTP Server
// =============================================================================

// MockResponse is a single static response: status code, body, optional
// headers, optional artificial delay.
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
	Delay      time.Duration
}

// MockResponseFunc is a dynamic response handler. It receives the request
// and the (already-buffered) body, and returns the response to send back.
// Useful for assert-on-body or vary-by-query-param tests.
type MockResponseFunc func(req *http.Request, body []byte) MockResponse

// responder is the internal storage shape for a (path, method) entry.
// Exactly one of static / sequence / fn is populated; the dispatcher in
// ServeHTTP decides which to use.
type responder struct {
	static   *MockResponse    // single sticky response
	sequence []MockResponse   // popped one-at-a-time; last item sticks
	fn       MockResponseFunc // dynamic handler
}

// MockHandler is a configurable HTTP handler for testing. It records every
// request it receives and serves responses registered via On / OnSequence /
// OnFunc.
type MockHandler struct {
	mu sync.Mutex
	// Responses is keyed by path → method → responder.
	// Exposed (capitalised) for backward compatibility with older tests
	// that read it directly; new code should use the On* registration
	// methods instead.
	Responses map[string]map[string]*responder
	Requests  []RecordedRequest
	// CallCount tracks how many times each "METHOD:/path" was invoked.
	CallCount map[string]int
}

// RecordedRequest stores information about a received request so tests can
// assert on what the SDK actually sent.
type RecordedRequest struct {
	Method  string
	Path    string
	Query   string
	Headers http.Header
	Body    []byte
}

// NewMockHandler creates a new mock handler.
func NewMockHandler() *MockHandler {
	return &MockHandler{
		Responses: make(map[string]map[string]*responder),
		Requests:  []RecordedRequest{},
		CallCount: make(map[string]int),
	}
}

// On registers a single sticky response for a (method, path) pair. Every
// request matching the pair returns this same response. Re-registering on
// the same pair overwrites the previous registration.
//
// Path matching: exact match preferred, with template fallback for
// segments that look like ":id", "*", or "{id}". Crucially this does NOT
// do prefix matching — a registration for `/foo` will not steal requests
// to `/foo/123` (which was a bug in earlier versions of this file).
func (m *MockHandler) On(method, path string, response MockResponse) *MockHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	r := response
	m.entry(path)[method] = &responder{static: &r}
	return m
}

// OnSequence registers a queue of responses for a (method, path) pair.
// Each request consumes the next response in the queue; once the queue is
// drained, the last response in the sequence becomes the sticky response
// for any further calls.
//
// Use this for paginated GetAll-style tests:
//
//	server.OnSequence("GET", path,
//	    common.SuccessResponse(common.ZPAListPaged(page1, 2)),
//	    common.SuccessResponse(common.ZPAListPaged(page2, 2)),
//	)
func (m *MockHandler) OnSequence(method, path string, responses ...MockResponse) *MockHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(responses) == 0 {
		return m
	}
	seq := make([]MockResponse, len(responses))
	copy(seq, responses)
	m.entry(path)[method] = &responder{sequence: seq}
	return m
}

// OnFunc registers a dynamic handler that builds the response per-request.
// Useful for assert-on-body or response-varies-by-query tests.
//
//	server.OnFunc("PUT", path, func(r *http.Request, body []byte) common.MockResponse {
//	    var got mypkg.Resource
//	    require.NoError(t, json.Unmarshal(body, &got))
//	    require.Equal(t, "expected", got.Name)
//	    return common.NoContentResponse()
//	})
func (m *MockHandler) OnFunc(method, path string, fn MockResponseFunc) *MockHandler {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entry(path)[method] = &responder{fn: fn}
	return m
}

// entry returns (and lazily creates) the inner method→responder map for a path.
// Caller must hold m.mu.
func (m *MockHandler) entry(path string) map[string]*responder {
	if m.Responses[path] == nil {
		m.Responses[path] = make(map[string]*responder)
	}
	return m.Responses[path]
}

// ServeHTTP implements http.Handler. It records the request, finds the
// matching responder, and writes the response.
func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	m.mu.Lock()
	m.Requests = append(m.Requests, RecordedRequest{
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.RawQuery,
		Headers: r.Header.Clone(),
		Body:    append([]byte(nil), body...),
	})
	m.CallCount[r.Method+":"+r.URL.Path]++

	resp, found := m.lookup(r.Method, r.URL.Path)
	m.mu.Unlock()

	if !found {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "no mock configured for path: ` + r.URL.Path + `"}`))
		return
	}

	mr := resp
	if respFunc := m.fnFor(r.Method, r.URL.Path); respFunc != nil {
		mr = respFunc(r, body)
	}

	if mr.Delay > 0 {
		time.Sleep(mr.Delay)
	}
	for k, v := range mr.Headers {
		w.Header().Set(k, v)
	}
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(mr.StatusCode)

	if mr.Body == nil {
		return
	}
	switch v := mr.Body.(type) {
	case string:
		_, _ = w.Write([]byte(v))
	case []byte:
		_, _ = w.Write(v)
	default:
		_ = json.NewEncoder(w).Encode(v)
	}
}

// fnFor checks if the matched (method, path) registration is a dynamic
// handler. Returns nil for static / sequence registrations.
func (m *MockHandler) fnFor(method, path string) MockResponseFunc {
	m.mu.Lock()
	defer m.mu.Unlock()
	rsp, ok := m.findResponder(method, path)
	if !ok {
		return nil
	}
	return rsp.fn
}

// lookup resolves a (method, path) request to a concrete MockResponse,
// consuming a sequence step if the matched registration is a queue.
// Caller must hold m.mu.
//
// Returns (resp, true) when a registration matches, (zero, false) otherwise.
// Note: dynamic handlers (fn) return a placeholder here and are dispatched
// in ServeHTTP via fnFor.
func (m *MockHandler) lookup(method, path string) (MockResponse, bool) {
	rsp, ok := m.findResponder(method, path)
	if !ok {
		return MockResponse{}, false
	}
	switch {
	case rsp.fn != nil:
		// real dispatch happens in ServeHTTP; placeholder OK
		return MockResponse{StatusCode: http.StatusOK}, true
	case len(rsp.sequence) > 0:
		next := rsp.sequence[0]
		// Pop unless this is the last element — keep sticky behaviour
		// so downstream calls after the sequence drains still respond.
		if len(rsp.sequence) > 1 {
			rsp.sequence = rsp.sequence[1:]
		}
		return next, true
	case rsp.static != nil:
		return *rsp.static, true
	}
	return MockResponse{}, false
}

// findResponder finds the responder for a (method, path) request, trying
// exact match first then template/wildcard match. Caller must hold m.mu.
func (m *MockHandler) findResponder(method, path string) (*responder, bool) {
	if methods, ok := m.Responses[path]; ok {
		if rsp, ok := methods[method]; ok {
			return rsp, true
		}
	}
	for pattern, methods := range m.Responses {
		if pattern == path {
			continue
		}
		if !matchPath(pattern, path) {
			continue
		}
		if rsp, ok := methods[method]; ok {
			return rsp, true
		}
	}
	return nil, false
}

// matchPath returns true when `pattern` matches `path`, treating any segment
// equal to "*", "{...}", or starting with ":" as a wildcard for that segment.
// Both pattern and path must have the same number of segments (no prefix
// matching — that was the bug in the previous implementation).
func matchPath(pattern, path string) bool {
	pp := strings.Split(pattern, "/")
	xp := strings.Split(path, "/")
	if len(pp) != len(xp) {
		return false
	}
	for i, seg := range pp {
		if seg == "*" || seg == "" || strings.HasPrefix(seg, ":") ||
			(strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}")) {
			continue
		}
		if seg != xp[i] {
			return false
		}
	}
	return true
}

// =============================================================================
// Test Server Factory
// =============================================================================

// TestServer wraps httptest.Server with the MockHandler attached.
type TestServer struct {
	*httptest.Server
	Handler *MockHandler
}

// NewTestServer creates a new in-memory test server with a fresh MockHandler.
// Callers should `defer ts.Close()` (or, when using NewZPATestService etc.,
// the cleanup is registered automatically).
func NewTestServer() *TestServer {
	handler := NewMockHandler()
	server := httptest.NewServer(handler)
	return &TestServer{
		Server:  server,
		Handler: handler,
	}
}

// On is a convenience wrapper around MockHandler.On that returns the
// TestServer for chaining.
func (ts *TestServer) On(method, path string, response MockResponse) *TestServer {
	ts.Handler.On(method, path, response)
	return ts
}

// OnSequence registers a multi-step response queue for paginated tests.
func (ts *TestServer) OnSequence(method, path string, responses ...MockResponse) *TestServer {
	ts.Handler.OnSequence(method, path, responses...)
	return ts
}

// OnFunc registers a dynamic response handler.
func (ts *TestServer) OnFunc(method, path string, fn MockResponseFunc) *TestServer {
	ts.Handler.OnFunc(method, path, fn)
	return ts
}

// AssertRequestCount fails the test if the recorded request count differs.
func (ts *TestServer) AssertRequestCount(t *testing.T, expected int) {
	t.Helper()
	require.Equal(t, expected, len(ts.Handler.Requests), "unexpected request count")
}

// GetCallCount returns how many times a particular (method, path) was called.
func (ts *TestServer) GetCallCount(method, path string) int {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()
	return ts.Handler.CallCount[method+":"+path]
}

// LastRequest returns the most recent request received, or nil.
func (ts *TestServer) LastRequest() *RecordedRequest {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()
	if len(ts.Handler.Requests) == 0 {
		return nil
	}
	r := ts.Handler.Requests[len(ts.Handler.Requests)-1]
	return &r
}

// AllRequests returns a snapshot copy of every recorded request so tests
// can iterate without racing the server's record loop.
func (ts *TestServer) AllRequests() []RecordedRequest {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()
	out := make([]RecordedRequest, len(ts.Handler.Requests))
	copy(out, ts.Handler.Requests)
	return out
}

// Reset clears all recorded requests, call counts, and registered
// responders. Useful between sub-tests when reusing a server.
func (ts *TestServer) Reset() {
	ts.Handler.mu.Lock()
	defer ts.Handler.mu.Unlock()
	ts.Handler.Responses = make(map[string]map[string]*responder)
	ts.Handler.Requests = nil
	ts.Handler.CallCount = make(map[string]int)
}

// =============================================================================
// Response Builders
// =============================================================================

// SuccessResponse creates a 200 OK response with the given body.
func SuccessResponse(body interface{}) MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body:       body,
	}
}

// SuccessResponseWithStatus creates a response with a custom status code.
func SuccessResponseWithStatus(statusCode int, body interface{}) MockResponse {
	return MockResponse{StatusCode: statusCode, Body: body}
}

// CreatedResponse creates a 201 Created response.
func CreatedResponse(body interface{}) MockResponse {
	return MockResponse{StatusCode: http.StatusCreated, Body: body}
}

// NoContentResponse creates a 204 No Content response.
func NoContentResponse() MockResponse {
	return MockResponse{StatusCode: http.StatusNoContent}
}

// NotFoundResponse creates a 404 Not Found response with a generic body.
func NotFoundResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	}
}

// UnauthorizedResponse creates a 401 Unauthorized response with the given body.
func UnauthorizedResponse(body string) MockResponse {
	return MockResponse{StatusCode: http.StatusUnauthorized, Body: body}
}

// SessionInvalidResponse mocks the OneAPI 401 SESSION_NOT_VALID payload.
func SessionInvalidResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       `{"code": "SESSION_NOT_VALID", "message": "Session is not valid"}`,
	}
}

// ResourceAccessBlockedResponse mocks a 401 with the access-blocked message.
func ResourceAccessBlockedResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusUnauthorized,
		Body:       `{"message": "Resource Access Blocked"}`,
	}
}

// TooManyRequestsResponse mocks a 429 with a numeric Retry-After.
// (The previous version rendered the int via `string(rune('0'+n))` which
// would produce garbage for n>9. Fixed.)
func TooManyRequestsResponse(retryAfterSeconds int) MockResponse {
	return MockResponse{
		StatusCode: http.StatusTooManyRequests,
		Body:       `{"message": "Rate limit exceeded"}`,
		Headers: map[string]string{
			"Retry-After": strconv.Itoa(retryAfterSeconds),
		},
	}
}

// TooManyRequestsResponseWithHeader mocks a 429 with a verbatim Retry-After
// header value (e.g. "120s" for ZIA, an HTTP-date for some clouds).
func TooManyRequestsResponseWithHeader(retryAfter string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusTooManyRequests,
		Body:       `{"message": "Rate limit exceeded"}`,
		Headers:    map[string]string{"Retry-After": retryAfter},
	}
}

// ConflictResponse creates a 409 Conflict response.
func ConflictResponse(body string) MockResponse {
	return MockResponse{StatusCode: http.StatusConflict, Body: body}
}

// EditLockResponse mocks ZIA's 409 EDIT_LOCK_NOT_AVAILABLE payload.
func EditLockResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusConflict,
		Body:       `{"code": "EDIT_LOCK_NOT_AVAILABLE", "message": "Edit lock not available"}`,
	}
}

// ServerErrorResponse creates a 500 Internal Server Error response.
func ServerErrorResponse() MockResponse {
	return MockResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       `{"message": "Internal server error"}`,
	}
}

// =============================================================================
// HTTP Response Helpers (low-level)
// =============================================================================

// CreateMockHTTPResponse creates an http.Response for tests that synthesise
// responses directly without going through the mock server.
func CreateMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// CreateMockHTTPResponseWithHeaders is CreateMockHTTPResponse plus headers.
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

// MockOAuthResponse returns a stock OAuth2 token response.
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

// MockZPAOAuthResponse is a ZPA-flavoured OAuth2 token response.
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

// RawResponse builds a response with raw bytes and explicit headers.
func RawResponse(body []byte, statusCode int, headers map[string]string) MockResponse {
	return MockResponse{StatusCode: statusCode, Body: body, Headers: headers}
}

// CSVResponse builds a CSV response (Content-Type: text/csv).
func CSVResponse(csvData string) MockResponse {
	return MockResponse{
		StatusCode: http.StatusOK,
		Body:       csvData,
		Headers:    map[string]string{"Content-Type": "text/csv"},
	}
}

// =============================================================================
// ZPA-Specific Convenience (kept for backward compatibility)
// =============================================================================

// MockSegmentGroupResponse returns a mock segment-group payload. Prefer
// constructing the real struct in your test for type safety.
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

// MockSegmentGroupListResponse wraps groups in the ZPA list envelope.
// New tests should use ZPAList(...) directly instead.
func MockSegmentGroupListResponse(groups ...map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"list":       groups,
		"totalPages": 1,
	}
}
