package ztw

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

// These tests lock in the retry and error-surfacing contract established in
// response to issue #449: a 5xx carrying a deterministic API verdict must not be
// retried, and an exhausted retry budget must not replace the API's error
// payload with a bare *url.Error.

func retryTestResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Request:    &http.Request{Method: http.MethodPost},
	}
}

func TestZTWCheckRetryOnServerErrors(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		wantRetry bool
	}{
		{
			name:      "deterministic 500 UNEXPECTED_ERROR is not retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred"}`,
			wantRetry: false,
		},
		{
			name:      "500 with org barrier marker is retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`,
			wantRetry: true,
		},
		{
			name:      "500 with empty body is retried",
			status:    http.StatusInternalServerError,
			body:      "",
			wantRetry: true,
		},
		{
			name:      "502 gateway HTML is retried",
			status:    http.StatusBadGateway,
			body:      "<html>502</html>",
			wantRetry: true,
		},
		{
			// Gateway and overload statuses never carry an API verdict, and 503
			// in particular is wired into Backoff as a Retry-After rate-limit
			// signal, so a JSON code in the body must not stop the retry.
			name:      "502 with a deterministic code is still retried",
			status:    http.StatusBadGateway,
			body:      `{"code":"SOME_PERMANENT_CODE","message":"nope"}`,
			wantRetry: true,
		},
		{
			name:      "503 with a deterministic code is still retried",
			status:    http.StatusServiceUnavailable,
			body:      `{"code":"SERVICE_UNAVAILABLE","message":"try later"}`,
			wantRetry: true,
		},
		{
			name:      "504 with a deterministic code is still retried",
			status:    http.StatusGatewayTimeout,
			body:      `{"code":"GATEWAY_TIMEOUT","message":"upstream timed out"}`,
			wantRetry: true,
		},
		{
			name:      "501 is not retried",
			status:    http.StatusNotImplemented,
			body:      "",
			wantRetry: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := retryTestResponse(tt.status, tt.body)
			shouldRetry, err := checkRetry(context.Background(), resp, nil)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRetry, shouldRetry)

			remaining, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr)
			assert.Equal(t, tt.body, string(remaining), "body must be rewound")
		})
	}
}

// TestZTWCheckRetryPreservesExistingBehaviour guards the retry paths that
// predate the #449 fix.
func TestZTWCheckRetryPreservesExistingBehaviour(t *testing.T) {
	t.Run("429 is retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), retryTestResponse(http.StatusTooManyRequests, ""), nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("409 edit lock is retried", func(t *testing.T) {
		resp := retryTestResponse(http.StatusConflict, `{"code":"EDIT_LOCK_NOT_AVAILABLE","message":"locked"}`)
		shouldRetry, err := checkRetry(context.Background(), resp, nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("412 org barrier is retried", func(t *testing.T) {
		resp := retryTestResponse(http.StatusPreconditionFailed, `{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`)
		shouldRetry, err := checkRetry(context.Background(), resp, nil)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("400 is not retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), retryTestResponse(http.StatusBadRequest, `{"code":"INVALID"}`), nil)
		require.NoError(t, err)
		assert.False(t, shouldRetry)
	})

	t.Run("transport errors are still retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), nil, io.ErrUnexpectedEOF)
		require.NoError(t, err)
		assert.True(t, shouldRetry)
	})

	t.Run("cancelled context is not retried", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		shouldRetry, err := checkRetry(ctx, retryTestResponse(http.StatusInternalServerError, ""), nil)
		assert.False(t, shouldRetry)
		assert.Error(t, err)
	})
}

func TestZTWErrorHandlerIsWired(t *testing.T) {
	cfg := &Configuration{}
	cfg.ZTW.Client.RateLimit.MaxRetries = 2
	cfg.ZTW.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZTW.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	httpClient := getHTTPClient(logger.GetDefaultLogger("ztw-test: "), nil, cfg)
	require.NotNil(t, httpClient)

	rt, ok := httpClient.Transport.(*retryablehttp.RoundTripper)
	require.True(t, ok, "expected a retryablehttp RoundTripper")
	require.NotNil(t, rt.Client.ErrorHandler, "ErrorHandler must be installed (issue #449)")

	// retryablehttp signals an exhausted budget by calling the handler with a
	// nil error; any non-nil error is a real failure that must not be masked.
	t.Run("preserves the response when the retry budget is exhausted", func(t *testing.T) {
		resp := retryTestResponse(http.StatusInternalServerError, `{"code":"UNEXPECTED_ERROR"}`)
		got, err := rt.Client.ErrorHandler(resp, nil, 3)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)
	})

	t.Run("does not mask a genuine error that arrives with a response", func(t *testing.T) {
		resp := retryTestResponse(http.StatusInternalServerError, `{"code":"UNEXPECTED_ERROR"}`)
		got, err := rt.Client.ErrorHandler(resp, context.Canceled, 3)
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

func TestZTWDeterministic500NotRetriedEndToEnd(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"UNEXPECTED_ERROR","message":"permanent failure"}`))
	}))
	defer srv.Close()

	cfg := &Configuration{}
	cfg.ZTW.Client.RateLimit.MaxRetries = 5
	cfg.ZTW.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZTW.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	httpClient := getHTTPClient(logger.GetDefaultLogger("ztw-test: "), nil, cfg)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/api/v1/x", bytes.NewBufferString("{}"))
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts), "deterministic 500 must not be retried")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
