package zwa

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

// TestZWABackoffLogic tests the ZWA backoff function behavior
func TestZWABackoffLogic(t *testing.T) {
	l := logger.GetDefaultLogger("zwa-test: ")

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
			retryAfter:  "4",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    5 * time.Second, // 4 + 1 padding
			description: "Should honor Retry-After header",
		},
		{
			name:        "503 with duration format",
			statusCode:  503,
			retryAfter:  "15s",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    16 * time.Second, // 15 + 1 padding
			description: "Should parse duration format",
		},
		{
			name:        "Exponential backoff - attempt 1",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second, // 2^1 * 2s = 4s
			description: "Should use exponential backoff",
		},
		{
			name:        "Exponential backoff - attempt 2",
			statusCode:  502,
			retryAfter:  "",
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second, // 2^2 * 2s = 8s
			description: "Should grow exponentially",
		},
		{
			name:        "Backoff capped at max",
			statusCode:  504,
			retryAfter:  "",
			attemptNum:  5,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // Capped at 10s
			description: "Should cap at maxWait",
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

	t.Log("ZWA backoff logic validated ✓")
}

// TestZWARateLimiter tests ZWA-specific rate limiting (20 GET per 10s, 10 POST/PUT/DELETE per 10s)
func TestZWARateLimiter(t *testing.T) {
	rateLimiter := rl.NewRateLimiter(20, 10, 10, 10)

	t.Run("ZWA GET rate limit (20 per 10s)", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			wait, _ := rateLimiter.Wait("GET")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("GET")
		require.True(t, wait, "21st GET should trigger wait")
		require.Greater(t, duration, time.Duration(0))

		t.Log("ZWA GET rate limit validated: 20 per 10s ✓")
	})

	t.Run("ZWA POST/PUT/DELETE rate limit (10 per 10s)", func(t *testing.T) {
		rateLimiter := rl.NewRateLimiter(20, 10, 10, 10)

		for i := 0; i < 10; i++ {
			wait, _ := rateLimiter.Wait("POST")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("POST")
		require.True(t, wait, "11th POST should trigger wait")
		require.Greater(t, duration, time.Duration(0))

		t.Log("ZWA POST/PUT/DELETE rate limit validated: 10 per 10s ✓")
	})
}

// TestZWAHeaderCaseInsensitive validates case-insensitive header parsing for ZWA
func TestZWAHeaderCaseInsensitive(t *testing.T) {
	l := logger.GetDefaultLogger("zwa-test: ")
	testCases := []string{"Retry-After", "retry-after", "RETRY-AFTER"}

	for _, headerName := range testCases {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set(headerName, "6")

		result := getRetryAfter(resp, l)
		require.Equal(t, 7*time.Second, result, "Header %q should work", headerName)
	}

	t.Log("ZWA Retry-After header case-insensitivity confirmed ✓")
}
