package zdx

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

// TestZDXBackoffLogic tests the ZDX backoff function behavior
func TestZDXBackoffLogic(t *testing.T) {
	l := logger.GetDefaultLogger("zdx-test: ")

	tests := []struct {
		name        string
		statusCode  int
		remaining   string
		attemptNum  int
		minWait     time.Duration
		maxWait     time.Duration
		expectedMin time.Duration
		expectedMax time.Duration
		description string
	}{
		{
			name:        "429 returns 2s fallback",
			statusCode:  429,
			remaining:   "",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expectedMin: 2 * time.Second,
			expectedMax: 2 * time.Second,
			description: "429 should return 2s fallback",
		},
		{
			name:        "Exponential backoff - attempt 1",
			statusCode:  500,
			remaining:   "",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expectedMin: 4 * time.Second,
			expectedMax: 4 * time.Second,
			description: "Should use exponential backoff: 2^1 * 2s = 4s",
		},
		{
			name:        "Exponential backoff - attempt 2",
			statusCode:  500,
			remaining:   "",
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expectedMin: 8 * time.Second,
			expectedMax: 8 * time.Second,
			description: "Should use exponential backoff: 2^2 * 2s = 8s",
		},
		{
			name:        "Backoff capped at max",
			statusCode:  502,
			remaining:   "",
			attemptNum:  10,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expectedMin: 10 * time.Second,
			expectedMax: 10 * time.Second,
			description: "Should cap at maxWait: 10s",
		},
	}

	backoffFunc := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			retryAfter := getRetryAfter(resp, l, 0)
			if retryAfter > 0 {
				return retryAfter
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
				if tt.remaining != "" {
					resp.Header.Set("X-Ratelimit-Remaining-Second", tt.remaining)
				}
			}

			result := backoffFunc(tt.minWait, tt.maxWait, tt.attemptNum, resp)

			// ZDX may include jitter, so use range-based assertion
			require.GreaterOrEqual(t, result, tt.expectedMin, tt.description)
			require.LessOrEqual(t, result, tt.expectedMax, tt.description)
		})
	}

	t.Log("ZDX backoff logic validated ✓")
}

// TestZDXGlobalRateLimiter tests ZDX global rate limiting (10 requests per second)
func TestZDXGlobalRateLimiter(t *testing.T) {
	globalLimiter := rl.NewGlobalRateLimiter(10, 1)

	t.Run("ZDX global rate limit (10 per second)", func(t *testing.T) {
		// First 10 requests should not wait
		for i := 0; i < 10; i++ {
			wait, duration := globalLimiter.Wait()
			require.False(t, wait, "Request %d should not wait", i+1)
			require.Equal(t, time.Duration(0), duration)
		}

		// 11th request should trigger wait
		wait, duration := globalLimiter.Wait()
		require.True(t, wait, "11th request should trigger wait")
		require.Greater(t, duration, time.Duration(0), "Wait duration should be > 0")

		t.Log("ZDX global rate limit validated: 10 per second ✓")
	})

	t.Run("Rate limit resets after time window", func(t *testing.T) {
		limiter := rl.NewGlobalRateLimiter(2, 1) // Very small limit for testing

		// Make 2 requests
		limiter.Wait()
		limiter.Wait()

		// 3rd should wait
		wait, _ := limiter.Wait()
		require.True(t, wait, "3rd request should wait")

		// Wait for window to pass
		time.Sleep(1100 * time.Millisecond)

		// Should not wait now
		wait, _ = limiter.Wait()
		require.False(t, wait, "After time window, should not wait")

		t.Log("ZDX rate limiter properly resets after time window ✓")
	})
}

// TestZDXRateLimitHeaders validates ZDX's X-Ratelimit header handling
func TestZDXRateLimitHeaders(t *testing.T) {
	l := logger.GetDefaultLogger("zdx-test: ")

	t.Run("429 returns 2s fallback", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}

		result := getRetryAfter(resp, l, 0)
		require.Equal(t, 2*time.Second, result, "429 should return 2s fallback")
	})

	t.Run("Default returns 500ms", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}

		result := getRetryAfter(resp, l, 0)
		require.Equal(t, 500*time.Millisecond, result, "Default should return 500ms")
	})

	t.Run("Preemptive backoff when approaching limit", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("X-Ratelimit-Remaining-Second", "2")
		resp.Header.Set("X-Ratelimit-Limit-Second", "10")

		result := getRetryAfter(resp, l, 5) // Threshold of 5
		// Should be 1s + jitter (0-500ms) = 1s to 1.5s
		require.GreaterOrEqual(t, result, 1*time.Second, "Preemptive backoff should be >= 1s")
		require.LessOrEqual(t, result, 2*time.Second, "Preemptive backoff should be <= 2s")
	})

	t.Log("ZDX rate limit header handling validated ✓")
}
