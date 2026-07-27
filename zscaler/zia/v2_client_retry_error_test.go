package zia

import (
	"bytes"
	"context"
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
			name:      "502 with a deterministic code is not retried",
			status:    http.StatusBadGateway,
			body:      `{"code":"SOME_PERMANENT_CODE","message":"nope"}`,
			wantRetry: false,
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

// newExhaustingClient mirrors the ErrorHandler wiring from getHTTPClient with a
// small retry budget, so the exhaustion path can be exercised quickly.
func newExhaustingClient(t *testing.T, retryMax int) *http.Client {
	t.Helper()
	l := logger.GetDefaultLogger("zia-test: ")

	rc := retryablehttp.NewClient()
	rc.RetryMax = retryMax
	rc.RetryWaitMin = time.Millisecond
	rc.RetryWaitMax = 2 * time.Millisecond
	rc.CheckRetry = checkRetry
	rc.Logger = nil
	rc.ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
		if resp != nil {
			l.Printf("[WARN] giving up after %d attempt(s); surfacing API response status %d", numTries, resp.StatusCode)
			return resp, nil
		}
		return nil, err
	}
	return rc.StandardClient()
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
