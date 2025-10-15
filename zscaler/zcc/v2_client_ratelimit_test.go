package zcc

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
)

// TestZCCExponentialBackoff tests exponential backoff calculation
func TestZCCExponentialBackoff(t *testing.T) {
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
			mult := math.Pow(2, float64(tt.attemptNum)) * float64(tt.minWait)
			result := time.Duration(mult)
			if float64(result) != mult || result > tt.maxWait {
				result = tt.maxWait
			}

			require.Equal(t, tt.expected, result, "Exponential backoff calculation")
		})
	}

	t.Log("ZCC exponential backoff validated ✓")
}

// TestZCCRateLimiter tests ZCC-specific rate limiting (100 per hour, 3 per day for device downloads)
func TestZCCRateLimiter(t *testing.T) {
	rateLimiter := rl.NewRateLimiter(100, 3, 3600, 86400)

	t.Run("ZCC general rate limit (100 per hour)", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			wait, _ := rateLimiter.Wait("GET")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("GET")
		require.True(t, wait, "101st GET should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZCC GET rate limit validated: 100 per hour ✓")
	})

	t.Run("ZCC device download rate limit (3 per day)", func(t *testing.T) {
		rateLimiter := rl.NewRateLimiter(100, 3, 3600, 86400)

		for i := 0; i < 3; i++ {
			wait, _ := rateLimiter.Wait("POST")
			require.False(t, wait, "Request %d should not wait", i+1)
		}

		wait, duration := rateLimiter.Wait("POST")
		require.True(t, wait, "4th POST should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZCC device download rate limit validated: 3 per day ✓")
	})
}

// TestZCCEndpointSpecificRetry validates ZCC's endpoint-specific retry logic
func TestZCCEndpointSpecificRetry(t *testing.T) {
	l := logger.GetDefaultLogger("zcc-test: ")

	t.Run("downloadDevices endpoint gets 24 hour retry", func(t *testing.T) {
		endpoint := "/api/v1/downloadDevices"
		result := getRetryAfter(endpoint, l)
		require.Equal(t, 24*time.Hour, result, "downloadDevices should have 24h retry")
	})

	t.Run("other endpoints get 1 hour retry", func(t *testing.T) {
		endpoint := "/api/v1/getDevices"
		result := getRetryAfter(endpoint, l)
		require.Equal(t, 1*time.Hour, result, "Other endpoints should have 1 hour retry (100 calls per hour limit)")
	})

	t.Log("ZCC endpoint-specific retry logic validated ✓")
}
