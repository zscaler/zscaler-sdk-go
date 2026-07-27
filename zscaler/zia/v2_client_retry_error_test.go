package zia

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
)

// These tests lock in the retry and error-surfacing contract established in
// response to issue #449, where a deterministic HTTP 500 was retried 100 times
// and the API's own error payload was replaced by a bare *url.Error.

func newTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request:    &http.Request{Method: http.MethodPost},
	}
}

// TestCheckRetryOnServerErrors asserts that a 5xx carrying a deterministic API
// verdict is not retried, while genuinely transient conditions still are.
func TestCheckRetryOnServerErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		body        string
		wantRetry   bool
		description string
	}{
		{
			name:        "deterministic 500 UNEXPECTED_ERROR is not retried",
			status:      http.StatusInternalServerError,
			body:        `{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred, please contact Zscaler's support"}`,
			wantRetry:   false,
			description: "the exact payload reported in issue #449",
		},
		{
			name:        "500 with org barrier marker is retried",
			status:      http.StatusInternalServerError,
			body:        `{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`,
			wantRetry:   true,
			description: "transient markers must win over the deterministic rule",
		},
		{
			name:      "500 with edit lock marker is retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"EDIT_LOCK_NOT_AVAILABLE","message":"edit lock unavailable"}`,
			wantRetry: true,
		},
		{
			name:      "500 with precondition marker is retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"UNEXPECTED_ERROR","message":"Request processing failed, possibly because an expected precondition was not met"}`,
			wantRetry: true,
		},
		{
			name:        "500 with empty body is retried",
			status:      http.StatusInternalServerError,
			body:        "",
			wantRetry:   true,
			description: "no API verdict means a possible infrastructure fault",
		},
		{
			name:        "500 with HTML gateway body is retried",
			status:      http.StatusInternalServerError,
			body:        "<html><body>502 Bad Gateway</body></html>",
			wantRetry:   true,
			description: "load balancer errors must retain the old behaviour",
		},
		{
			name:        "503 without a code is retried",
			status:      http.StatusServiceUnavailable,
			body:        `{"message":"service unavailable"}`,
			wantRetry:   true,
			description: "a message with no code is not a deterministic verdict",
		},
		{
			name:        "502 with a deterministic code is still retried",
			status:      http.StatusBadGateway,
			body:        `{"code":"SOME_PERMANENT_CODE","message":"nope"}`,
			wantRetry:   true,
			description: "gateway statuses never carry an API verdict",
		},
		{
			name:        "503 with a deterministic code is still retried",
			status:      http.StatusServiceUnavailable,
			body:        `{"code":"SERVICE_UNAVAILABLE","message":"try later"}`,
			wantRetry:   true,
			description: "503 is wired into Backoff as a Retry-After rate-limit signal",
		},
		{
			name:        "504 with a deterministic code is still retried",
			status:      http.StatusGatewayTimeout,
			body:        `{"code":"GATEWAY_TIMEOUT","message":"upstream timed out"}`,
			wantRetry:   true,
			description: "gateway statuses never carry an API verdict",
		},
		{
			name:      "501 Not Implemented is never retried",
			status:    http.StatusNotImplemented,
			body:      "",
			wantRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := newTestResponse(tt.status, tt.body)

			shouldRetry, err := checkRetry(context.Background(), resp, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRetry, shouldRetry, tt.description)

			// The body must survive the retry decision so that downstream error
			// parsing can still read it.
			remaining, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr)
			assert.Equal(t, tt.body, string(remaining), "response body must be rewound")
		})
	}
}

// TestCheckRetryPreservesExistingBehaviour guards the retry paths that predate
// the #449 fix, so narrowing the 5xx rule did not regress them.
func TestCheckRetryPreservesExistingBehaviour(t *testing.T) {
	t.Run("429 is retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), newTestResponse(http.StatusTooManyRequests, ""), nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("409 edit lock is retried", func(t *testing.T) {
		body := `{"code":"EDIT_LOCK_NOT_AVAILABLE","message":"edit lock"}`
		shouldRetry, err := checkRetry(context.Background(), newTestResponse(http.StatusConflict, body), nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("412 org barrier is retried", func(t *testing.T) {
		body := `{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`
		shouldRetry, err := checkRetry(context.Background(), newTestResponse(http.StatusPreconditionFailed, body), nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("400 is not retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), newTestResponse(http.StatusBadRequest, `{"code":"INVALID_INPUT"}`), nil)
		require.NoError(t, err)
		assert.False(t, shouldRetry)
	})

	t.Run("transport errors are still retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), nil, io.ErrUnexpectedEOF)
		require.NoError(t, err)
		assert.True(t, shouldRetry, "connection level failures remain retryable")
	})

	t.Run("cancelled context is not retried", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		shouldRetry, err := checkRetry(ctx, newTestResponse(http.StatusInternalServerError, ""), nil)
		assert.False(t, shouldRetry)
		assert.Error(t, err)
	})
}

// TestMaxNumOfRetriesDefault pins the retry budget so a future change cannot
// silently restore the 100-attempt behaviour reported in issue #449.
func TestMaxNumOfRetriesDefault(t *testing.T) {
	assert.Equal(t, 10, MaxNumOfRetries,
		"raise this assertion deliberately if changing the SDK retry contract")
}

// newExhaustingClient builds the real production client through getHTTPClient
// with a small retry budget, so the exhaustion path can be exercised quickly.
// It must not re-implement the ErrorHandler: a copy would keep passing if the
// production wiring regressed.
func newExhaustingClient(t *testing.T, retryMax int) *http.Client {
	t.Helper()

	cfg := &Configuration{}
	cfg.ZIA.Client.RateLimit.MaxRetries = int32(retryMax)
	cfg.ZIA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZIA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	client := getHTTPClient(logger.GetDefaultLogger("zia-test: "), nil, cfg)
	require.NotNil(t, client)

	rt, ok := client.Transport.(*retryablehttp.RoundTripper)
	require.True(t, ok, "expected a retryablehttp RoundTripper")
	require.NotNil(t, rt.Client.ErrorHandler, "ErrorHandler must be installed (issue #449)")

	return client
}

// TestZIAErrorHandlerSemantics pins the two branches of the installed
// ErrorHandler, which retryablehttp distinguishes by whether err is nil.
func TestZIAErrorHandlerSemantics(t *testing.T) {
	cfg := &Configuration{}
	cfg.ZIA.Client.RateLimit.MaxRetries = 2
	cfg.ZIA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZIA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	rt, ok := getHTTPClient(logger.GetDefaultLogger("zia-test: "), nil, cfg).Transport.(*retryablehttp.RoundTripper)
	require.True(t, ok)
	require.NotNil(t, rt.Client.ErrorHandler)

	t.Run("preserves the response when the retry budget is exhausted", func(t *testing.T) {
		got, err := rt.Client.ErrorHandler(newTestResponse(http.StatusInternalServerError, `{"code":"UNEXPECTED_ERROR"}`), nil, 3)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)
	})

	t.Run("does not mask a genuine error that arrives with a response", func(t *testing.T) {
		got, err := rt.Client.ErrorHandler(newTestResponse(http.StatusInternalServerError, `{"code":"UNEXPECTED_ERROR"}`), context.Canceled, 3)
		assert.Nil(t, got, "a cancelled context must not be reported as a response")
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("reports the transport error when no response exists", func(t *testing.T) {
		sentinel := errors.New("connection refused")
		got, err := rt.Client.ErrorHandler(nil, sentinel, 3)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, sentinel)
	})
}

// TestErrorHandlerPreservesResponseOnExhaustion is the core regression test for
// issue #449: when the retry budget runs out, the caller must still receive the
// API response rather than a sanitized *url.Error.
func TestErrorHandlerPreservesResponseOnExhaustion(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		// Force exhaustion regardless of the checkRetry verdict by using a
		// transient marker, isolating the ErrorHandler behaviour under test.
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`))
	}))
	defer srv.Close()

	client := newExhaustingClient(t, 2)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/urlFilteringRules", bytes.NewBufferString("{}"))
	require.NoError(t, err)

	resp, err := client.Do(req)

	require.NoError(t, err, "exhausted retries must not surface a transport error")
	require.NotNil(t, resp, "the final API response must be preserved")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, int32(3), atomic.LoadInt32(&attempts), "RetryMax 2 means 3 attempts")

	// The preserved response must convert into a structured SDK error.
	apiErr := errorx.CheckErrorInResponse(resp, nil)
	respErr, ok := errorx.AsErrorResponse(apiErr)
	require.True(t, ok, "caller must be able to recover an *errorx.ErrorResponse")
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, "UNEXPECTED_ERROR", respErr.Parsed.Code)
	assert.Equal(t, http.StatusInternalServerError, respErr.Parsed.Status)
	assert.NotEmpty(t, respErr.Parsed.Message)
}

// TestErrorHandlerReturnsTransportErrorWhenNoResponse ensures the fix does not
// swallow genuine connection failures, where there is no response to surface.
func TestErrorHandlerReturnsTransportErrorWhenNoResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	unreachable := srv.URL
	srv.Close() // nothing is listening now

	client := newExhaustingClient(t, 1)
	req, err := http.NewRequest(http.MethodGet, unreachable+"/api/v1/status", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)

	require.Error(t, err, "a connection failure must still be reported")
	assert.Nil(t, resp)

	var urlErr *url.Error
	assert.ErrorAs(t, err, &urlErr)
}

// TestDeterministic500StopsAfterOneAttempt ties the two fixes together: the
// scenario from issue #449 must now fail fast instead of consuming the budget.
func TestDeterministic500StopsAfterOneAttempt(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred, please contact Zscaler's support"}`))
	}))
	defer srv.Close()

	client := newExhaustingClient(t, MaxNumOfRetries)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/urlFilteringRules", bytes.NewBufferString("{}"))
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts),
		"a deterministic API verdict must not be retried at all")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// newGenericRequestTestClient builds a ZIA Client wired to srv with a
// pre-established session, so GenericRequest can be exercised without the
// authentication round trip.
func newGenericRequestTestClient(t *testing.T, srvURL string, retryMax int) *Client {
	t.Helper()

	cfg := &Configuration{}
	cfg.ZIA.Client.RateLimit.MaxRetries = int32(retryMax)
	cfg.ZIA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZIA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	return &Client{
		URL:              srvURL,
		HTTPClient:       getHTTPClient(logger.GetDefaultLogger("zia-test: "), nil, cfg),
		Logger:           logger.GetDefaultLogger("zia-test: "),
		session:          &Session{JSessionID: "test-session"},
		sessionRefreshed: time.Now(),
		sessionTimeout:   time.Hour,
	}
}

// TestGenericRequestSurfacesDeterministic500 reproduces issue #449 end to end
// through the exact call path a caller uses: the API's code and message must
// reach the caller as an *errorx.ErrorResponse, not an opaque *url.Error.
func TestGenericRequestSurfacesDeterministic500(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred, please contact Zscaler's support"}`))
	}))
	defer srv.Close()

	c := newGenericRequestTestClient(t, srv.URL, MaxNumOfRetries)

	body, err := c.GenericRequest(context.Background(), srv.URL, "/api/v1/urlFilteringRules",
		http.MethodPost, bytes.NewBufferString("{}"), nil, contentTypeJSON)

	require.Error(t, err)
	assert.Nil(t, body)
	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts),
		"a deterministic API verdict must not be retried at all")

	var urlErr *url.Error
	assert.False(t, errors.As(err, &urlErr), "the sanitized *url.Error from issue #449 must be gone")

	respErr, ok := errorx.AsErrorResponse(err)
	require.True(t, ok, "caller must receive a structured *errorx.ErrorResponse")
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, "UNEXPECTED_ERROR", respErr.Parsed.Code)
	assert.Equal(t, http.StatusInternalServerError, respErr.Parsed.Status)
	assert.Contains(t, respErr.Parsed.Message, "An unexpected error has occurred")
}

// TestGenericRequestSurfacesPersistent401 covers the post-loop guard: a 401 that
// survives every attempt must not be handed back as a successful body.
func TestGenericRequestSurfacesPersistent401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code":"AUTHENTICATION_FAILED","message":"not authorized"}`))
	}))
	defer srv.Close()

	c := newGenericRequestTestClient(t, srv.URL, 1)

	body, err := c.GenericRequest(context.Background(), srv.URL, "/api/v1/urlCategories",
		http.MethodGet, nil, nil, contentTypeJSON)

	require.Error(t, err, "a persistent 401 must not be reported as success")
	assert.Nil(t, body)

	respErr, ok := errorx.AsErrorResponse(err)
	require.True(t, ok)
	require.NotNil(t, respErr.Parsed)
	assert.Equal(t, http.StatusUnauthorized, respErr.Parsed.Status)
	assert.Equal(t, "AUTHENTICATION_FAILED", respErr.Parsed.Code)
}

// TestGenericRequestSucceeds guards the happy path against the new post-loop
// guard: a 2xx must still return its body.
func TestGenericRequestSucceeds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":42}`))
	}))
	defer srv.Close()

	c := newGenericRequestTestClient(t, srv.URL, 1)

	body, err := c.GenericRequest(context.Background(), srv.URL, "/api/v1/urlCategories",
		http.MethodGet, nil, nil, contentTypeJSON)

	require.NoError(t, err)
	assert.JSONEq(t, `{"id":42}`, string(body))
}
