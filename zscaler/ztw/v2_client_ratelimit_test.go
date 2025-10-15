package ztw

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
)

// TestZTWBackoffLogic tests the ZTW backoff function behavior
func TestZTWBackoffLogic(t *testing.T) {
	l := logger.GetDefaultLogger("ztw-test: ")

	tests := []struct {
		name        string
		statusCode  int
		retryAfter  string
		attemptNum  int
		minWait     time.Duration
		maxWait     time.Duration
		expected    time.Duration
		description string
	}{
		{
			name:        "429 with Retry-After",
			statusCode:  429,
			retryAfter:  "5",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    5 * time.Second, // ZTW doesn't add padding
			description: "Should honor Retry-After header",
		},
		{
			name:        "Exponential backoff - attempt 1",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second,
			description: "Should use exponential backoff: 2^1 * 2s = 4s",
		},
		{
			name:        "Exponential backoff - attempt 2",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second,
			description: "Should use exponential backoff: 2^2 * 2s = 8s",
		},
		{
			name:        "Exponential backoff - capped",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  5,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second,
			description: "Should cap at maxWait: 10s",
		},
	}

	backoffFunc := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				retryAfter := getRetryAfter(resp, l)
				if retryAfter > 0 {
					return retryAfter
				}
			}
		}
		mult := math.Pow(2, float64(attemptNum)) * float64(min)
		sleep := time.Duration(mult)
		if float64(sleep) != mult || sleep > max {
			sleep = max
		}
		return sleep
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			if tt.statusCode > 0 {
				resp = &http.Response{
					StatusCode: tt.statusCode,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
				if tt.retryAfter != "" {
					resp.Header.Set("Retry-After", tt.retryAfter)
				}
			}

			result := backoffFunc(tt.minWait, tt.maxWait, tt.attemptNum, resp)
			require.Equal(t, tt.expected, result, tt.description)
		})
	}

	t.Log("ZTW backoff logic validated ✓")
}

// TestZTWRateLimiter tests ZTW-specific rate limiting (20 GET per 10s, 10 POST/PUT/DELETE per 61s)
func TestZTWRateLimiter(t *testing.T) {
	rateLimiter := rl.NewRateLimiter(20, 10, 10, 61)

	t.Run("ZTW GET rate limit (20 per 10s)", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			wait, _ := rateLimiter.Wait("GET")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("GET")
		require.True(t, wait, "21st GET should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZTW GET rate limit validated: 20 per 10s ✓")
	})

	t.Run("ZTW POST/PUT/DELETE rate limit (10 per 61s)", func(t *testing.T) {
		rateLimiter := rl.NewRateLimiter(20, 10, 10, 61)

		for i := 0; i < 10; i++ {
			wait, _ := rateLimiter.Wait("POST")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("POST")
		require.True(t, wait, "11th POST should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZTW POST/PUT/DELETE rate limit validated: 10 per 61s ✓")
	})
}

// TestZTWHeaderCaseInsensitive validates case-insensitive header parsing for ZTW
func TestZTWHeaderCaseInsensitive(t *testing.T) {
	l := logger.GetDefaultLogger("ztw-test: ")
	testCases := []string{"Retry-After", "retry-after", "RETRY-AFTER"}

	for _, headerName := range testCases {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set(headerName, "5")

		result := getRetryAfter(resp, l)
		require.Equal(t, 5*time.Second, result, "Header %q should work (ZTW doesn't add padding)", headerName)
	}

	t.Log("ZTW Retry-After header case-insensitivity confirmed ✓")
}
