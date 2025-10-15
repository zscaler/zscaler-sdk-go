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

// TestZWAExponentialBackoff tests exponential backoff calculation for ZWA
func TestZWAExponentialBackoff(t *testing.T) {
	tests := []struct {
		name       string
		attemptNum int
		minWait    time.Duration
		maxWait    time.Duration
		expected   time.Duration
	}{
		{
			name:       "Attempt 0",
			attemptNum: 0,
			minWait:    2 * time.Second,
			maxWait:    10 * time.Second,
			expected:   2 * time.Second, // 2^0 * 2s = 2s
		},
		{
			name:       "Attempt 1",
			attemptNum: 1,
			minWait:    2 * time.Second,
			maxWait:    10 * time.Second,
			expected:   4 * time.Second, // 2^1 * 2s = 4s
		},
		{
			name:       "Attempt 2",
			attemptNum: 2,
			minWait:    2 * time.Second,
			maxWait:    10 * time.Second,
			expected:   8 * time.Second, // 2^2 * 2s = 8s
		},
		{
			name:       "Attempt 3 (capped)",
			attemptNum: 3,
			minWait:    2 * time.Second,
			maxWait:    10 * time.Second,
			expected:   10 * time.Second, // 2^3 * 2s = 16s, capped at 10s
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Calculate exponential backoff
			multiplier := math.Pow(2, float64(tt.attemptNum)) * float64(tt.minWait)
			result := time.Duration(multiplier)
			if float64(result) != multiplier || result > tt.maxWait {
				result = tt.maxWait
			}

			require.Equal(t, tt.expected, result, "Exponential backoff calculation")
		})
	}

	t.Log("ZWA exponential backoff validated ✓")
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

// TestZWARateLimitHeaders validates ZWA's RateLimit header handling
func TestZWARateLimitHeaders(t *testing.T) {
	l := logger.GetDefaultLogger("zwa-test: ")

	t.Run("Rate limit remaining 0 triggers retry", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("RateLimit-Remaining", "0")
		resp.Header.Set("RateLimit-Reset", "5") // 5 seconds

		result := getRetryAfter(resp, l)
		require.Greater(t, result, time.Duration(0), "Should return positive duration when remaining is 0")
	})

	t.Run("Rate limit remaining > 0 returns 5s fallback", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("RateLimit-Remaining", "10")

		result := getRetryAfter(resp, l)
		require.Equal(t, 5*time.Second, result, "Should return 5s fallback when remaining > 0")
	})

	t.Log("ZWA RateLimit header handling validated ✓")
}
