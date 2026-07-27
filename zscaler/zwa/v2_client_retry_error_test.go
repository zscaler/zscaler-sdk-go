package zwa

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
		Request:    &http.Request{Method: http.MethodGet},
	}
}

func TestZWACheckRetryOnServerErrors(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		wantRetry bool
	}{
		{
			name:      "deterministic 500 with API code is not retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred"}`,
			wantRetry: false,
		},
		{
			name:      "500 with edit lock marker is retried",
			status:    http.StatusInternalServerError,
			body:      `{"code":"EDIT_LOCK_NOT_AVAILABLE"}`,
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

func TestZWACheckRetryPreservesExistingBehaviour(t *testing.T) {
	t.Run("429 is retried", func(t *testing.T) {
		shouldRetry, err := checkRetry(context.Background(), retryTestResponse(http.StatusTooManyRequests, ""), nil)
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

func TestZWAErrorHandlerIsWired(t *testing.T) {
	cfg := &Configuration{}
	cfg.ZWA.Client.RateLimit.MaxRetries = 2
	cfg.ZWA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZWA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	httpClient := getHTTPClient(logger.GetDefaultLogger("zwa-test: "), nil, cfg)
	require.NotNil(t, httpClient)

	rt, ok := httpClient.Transport.(*retryablehttp.RoundTripper)
	require.True(t, ok, "expected a retryablehttp RoundTripper")
	require.NotNil(t, rt.Client.ErrorHandler, "ErrorHandler must be installed (issue #449)")

	t.Run("preserves the response when one exists", func(t *testing.T) {
		resp := retryTestResponse(http.StatusInternalServerError, `{"code":"UNEXPECTED_ERROR"}`)
		got, err := rt.Client.ErrorHandler(resp, errors.New("giving up"), 3)
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, http.StatusInternalServerError, got.StatusCode)
	})

	t.Run("reports the transport error when no response exists", func(t *testing.T) {
		sentinel := errors.New("connection refused")
		got, err := rt.Client.ErrorHandler(nil, sentinel, 3)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, sentinel)
	})
}

func TestZWADeterministic500NotRetriedEndToEnd(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"code":"UNEXPECTED_ERROR","message":"permanent failure"}`))
	}))
	defer srv.Close()

	cfg := &Configuration{}
	cfg.ZWA.Client.RateLimit.MaxRetries = 5
	cfg.ZWA.Client.RateLimit.RetryWaitMin = time.Millisecond
	cfg.ZWA.Client.RateLimit.RetryWaitMax = 2 * time.Millisecond

	httpClient := getHTTPClient(logger.GetDefaultLogger("zwa-test: "), nil, cfg)
	req, err := http.NewRequest(http.MethodGet, srv.URL+"/zwa/api/v1/x", nil)
	require.NoError(t, err)

	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	assert.Equal(t, int32(1), atomic.LoadInt32(&attempts), "deterministic 500 must not be retried")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
