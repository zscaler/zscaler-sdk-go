package zia

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

// TestZIABackoffLogic tests the ZIA backoff function behavior
func TestZIABackoffLogic(t *testing.T) {
	l := logger.GetDefaultLogger("zia-test: ")

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
			name:        "429 with Retry-After header",
			statusCode:  429,
			retryAfter:  "5",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    5 * time.Second, // ZIA doesn't add padding
			description: "Should honor Retry-After header",
		},
		{
			name:        "503 with Retry-After duration",
			statusCode:  503,
			retryAfter:  "8s",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second, // ZIA doesn't add padding
			description: "Should parse duration format",
		},
		{
			name:        "Exponential backoff - attempt 1",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second, // 2^1 * 2s
			description: "Should use exponential backoff",
		},
		{
			name:        "Exponential backoff - attempt 2",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second, // 2^2 * 2s
			description: "Should grow exponentially",
		},
		{
			name:        "Exponential backoff - capped at max",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  5,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // 2^5 * 2s = 64s, capped at 10s
			description: "Should cap at maxWait",
		},
	}

	// Test the backoff calculation logic
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

	t.Log("ZIA backoff logic validated ✓")
}

// TestZIARateLimiter tests the ZIA-specific rate limiter behavior
func TestZIARateLimiter(t *testing.T) {
	// ZIA rate limits: 100 GET per 10s, 1 POST/PUT/DELETE per 1s
	rateLimiter := rl.NewRateLimiter(100, 1, 10, 1)

	t.Run("ZIA GET rate limit (100 per 10s)", func(t *testing.T) {
		// First 100 GET requests should not wait
		for i := 0; i < 100; i++ {
			wait, duration := rateLimiter.Wait("GET")
			require.False(t, wait, "Request %d should not wait", i+1)
			require.Equal(t, time.Duration(0), duration, "Duration should be 0")
		}

		// 101st GET should trigger wait
		wait, duration := rateLimiter.Wait("GET")
		require.True(t, wait, "101st GET should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZIA GET rate limit validated: 100 per 10s ✓")
	})

	t.Run("ZIA POST/PUT/DELETE rate limit (1 per 1s)", func(t *testing.T) {
		rateLimiter := rl.NewRateLimiter(100, 1, 10, 1)

		// First POST should not wait
		wait, _ := rateLimiter.Wait("POST")
		require.False(t, wait, "First POST should not wait")

		// Second POST should trigger wait (only 1 allowed per 1s)
		wait, duration := rateLimiter.Wait("POST")
		require.True(t, wait, "Second POST should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZIA POST/PUT/DELETE rate limit validated: 1 per 1s ✓")
	})
}

// TestZIARateLimiterProductionParams validates the exact parameters used in v2_config.go
// and v2_client.go after the rate limiter fix (20 GET/10s, 10 POST-PUT-DELETE/10s).
// Previously the legacy v2 clients incorrectly passed retry config values as rate limiter
// parameters, and the postPutDeleteFreq was 61s instead of 10s.
func TestZIARateLimiterProductionParams(t *testing.T) {
	limiter := rl.NewRateLimiter(20, 10, 10, 10)

	t.Run("20 GETs allowed then throttled", func(t *testing.T) {
		lim := rl.NewRateLimiter(20, 10, 10, 10)
		for i := 0; i < 20; i++ {
			wait, _ := lim.Wait("GET")
			require.False(t, wait, "GET #%d should not wait", i+1)
		}
		wait, delay := lim.Wait("GET")
		require.True(t, wait, "21st GET should trigger wait")
		require.LessOrEqual(t, delay, 10*time.Second, "delay should be within 10s window")
	})

	t.Run("10 POSTs allowed then throttled within 10s window", func(t *testing.T) {
		lim := rl.NewRateLimiter(20, 10, 10, 10)
		for i := 0; i < 10; i++ {
			wait, _ := lim.Wait("POST")
			require.False(t, wait, "POST #%d should not wait", i+1)
		}
		wait, delay := lim.Wait("POST")
		require.True(t, wait, "11th POST should trigger wait")
		require.LessOrEqual(t, delay, 10*time.Second,
			"delay should be <= 10s (not 61s as in the old config)")
		require.Greater(t, delay, time.Duration(0))
	})

	_ = limiter
}

// TestZIARetryAfterHeaderCaseInsensitive validates case-insensitive header parsing
func TestZIARetryAfterHeaderCaseInsensitive(t *testing.T) {
	l := logger.GetDefaultLogger("zia-test: ")

	testCases := []string{"Retry-After", "retry-after", "RETRY-AFTER", "ReTrY-aFtEr"}

	for _, headerName := range testCases {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set(headerName, "5")

		result := getRetryAfter(resp, l)
		expected := 5 * time.Second // ZIA doesn't add padding

		require.Equal(t, expected, result, "Header %q should work", headerName)
	}

	t.Log("ZIA Retry-After header case-insensitivity confirmed ✓")
}
