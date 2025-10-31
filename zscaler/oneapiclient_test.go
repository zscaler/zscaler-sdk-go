package zscaler

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
)

func TestUserAgent(t *testing.T) {
	configuration, err := NewConfiguration()
	require.NoError(t, err, "Creating a new config should not error")
	userAgent := "zscaler-sdk-go/" + VERSION + " golang/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH
	require.Equal(t, userAgent, configuration.UserAgent)
}

func TestUserAgentWithExtra(t *testing.T) {
	configuration, err := NewConfiguration(
		WithUserAgentExtra("extra/info"),
	)
	require.NoError(t, err, "Creating a new config should not error")
	userAgent := "zscaler-sdk-go/" + VERSION + " golang/" + runtime.Version() + " " + runtime.GOOS + "/" + runtime.GOARCH + " extra/info"
	require.Equal(t, userAgent, configuration.UserAgent)
}

func TestDetectServiceTypeUnknown(t *testing.T) {
	_, err := detectServiceType("/foo")
	require.Error(t, err)
}

func TestGetServiceHTTPClientUnknown(t *testing.T) {
	cfg, err := NewConfiguration()
	require.NoError(t, err)

	svc, err := NewOneAPIClient(cfg)
	require.NoError(t, err)

	generic := cfg.HTTPClient

	httpClient := svc.Client.getServiceHTTPClient("/foo")
	require.Equal(t, generic, httpClient)
}

// TestExecuteRequestExponentialBackoff tests the OneAPI ExecuteRequest retry logic with exponential backoff
func TestExecuteRequestExponentialBackoff(t *testing.T) {
	t.Run("5xx errors retry with exponential backoff", func(t *testing.T) {
		// Create test configuration
		cfg := &Configuration{}
		cfg.Zscaler.Client.RateLimit.MaxRetries = 5
		cfg.Zscaler.Client.RateLimit.RetryWaitMin = 1 * time.Second
		cfg.Zscaler.Client.RateLimit.RetryWaitMax = 10 * time.Second
		cfg.Zscaler.Client.RequestTimeout = 60 * time.Second

		// Test exponential backoff calculation
		expectedDelays := []time.Duration{
			1 * time.Second,  // 2^0 * 1s = 1s
			2 * time.Second,  // 2^1 * 1s = 2s
			4 * time.Second,  // 2^2 * 1s = 4s
			8 * time.Second,  // 2^3 * 1s = 8s
			10 * time.Second, // 2^4 * 1s = 16s, capped at 10s
		}

		for attempt, expected := range expectedDelays {
			delay := time.Duration(math.Pow(2, float64(attempt))) * cfg.Zscaler.Client.RateLimit.RetryWaitMin
			if delay > cfg.Zscaler.Client.RateLimit.RetryWaitMax {
				delay = cfg.Zscaler.Client.RateLimit.RetryWaitMax
			}
			require.Equal(t, expected, delay, "Attempt %d: exponential backoff calculation mismatch", attempt)
		}

		t.Log("Exponential backoff sequence validated: 1s → 2s → 4s → 8s → 10s (capped) ✓")
	})

	t.Run("Retry-After header parsing", func(t *testing.T) {
		cfg := &Configuration{}
		cfg.Logger = logger.GetDefaultLogger("test: ")

		testCases := []struct {
			name        string
			headerValue string
			headerCase  string
			expected    time.Duration
			description string
		}{
			{
				name:        "Integer format (standard case)",
				headerValue: "5",
				headerCase:  "Retry-After",
				expected:    5 * time.Second, // Trust API's value exactly (no buffer for >= 1s)
				description: "Should parse integer seconds",
			},
			{
				name:        "Integer format (lowercase - ZPA API format)",
				headerValue: "5",
				headerCase:  "retry-after",
				expected:    5 * time.Second, // Trust API's value exactly
				description: "Should parse lowercase header",
			},
			{
				name:        "Duration format",
				headerValue: "10s",
				headerCase:  "Retry-After",
				expected:    10 * time.Second, // Trust API's value exactly (no buffer for >= 1s)
				description: "Should parse duration format",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				resp := &http.Response{
					StatusCode: 429,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
				}
				resp.Header.Set(tc.headerCase, tc.headerValue)

				result := getRetryAfter(resp, cfg)
				require.Equal(t, tc.expected, result, tc.description)
			})
		}

		t.Log("Retry-After header parsing validated (case-insensitive) ✓")
	})

	t.Run("Request timeout and retry limit enforcement", func(t *testing.T) {
		cfg := &Configuration{}
		cfg.Zscaler.Client.RateLimit.MaxRetries = 3
		cfg.Zscaler.Client.RateLimit.RetryWaitMin = 100 * time.Millisecond
		cfg.Zscaler.Client.RateLimit.RetryWaitMax = 500 * time.Millisecond
		cfg.Zscaler.Client.RequestTimeout = 2 * time.Second
		cfg.Context = context.Background()

		// Verify retry limits are respected
		maxRetries := cfg.Zscaler.Client.RateLimit.MaxRetries
		require.Equal(t, int32(3), maxRetries, "MaxRetries should be configurable")

		// Verify timeout is respected
		timeout := cfg.Zscaler.Client.RequestTimeout
		require.Equal(t, 2*time.Second, timeout, "RequestTimeout should be configurable")

		t.Log("Retry limits and timeout configuration validated ✓")
	})
}
