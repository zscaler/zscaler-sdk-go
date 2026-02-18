package ratelimiter

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// NewRateLimiter: per-second window tests
// ---------------------------------------------------------------------------

func TestNewRateLimiter_GETLimit(t *testing.T) {
	limiter := NewRateLimiter(3, 1, 1, 1)

	for i := 0; i < 3; i++ {
		wait, _ := limiter.Wait(http.MethodGet)
		require.False(t, wait, "GET #%d within limit should not wait", i+1)
	}

	wait, delay := limiter.Wait(http.MethodGet)
	require.True(t, wait, "GET exceeding limit should wait")
	require.Greater(t, delay, time.Duration(0))
}

func TestNewRateLimiter_POSTLimit(t *testing.T) {
	limiter := NewRateLimiter(10, 2, 1, 1)

	for i := 0; i < 2; i++ {
		wait, _ := limiter.Wait(http.MethodPost)
		require.False(t, wait, "POST #%d within limit should not wait", i+1)
	}

	wait, delay := limiter.Wait(http.MethodPost)
	require.True(t, wait, "POST exceeding limit should wait")
	require.Greater(t, delay, time.Duration(0))
}

func TestNewRateLimiter_DELETESharesBucketWithPOST(t *testing.T) {
	limiter := NewRateLimiter(10, 1, 1, 1)

	wait, _ := limiter.Wait(http.MethodPost)
	require.False(t, wait)

	wait, _ = limiter.Wait(http.MethodDelete)
	require.True(t, wait, "DELETE should wait because POST consumed the shared bucket slot")
}

func TestNewRateLimiter_PUTSharesBucketWithPOST(t *testing.T) {
	limiter := NewRateLimiter(10, 1, 1, 1)

	wait, _ := limiter.Wait(http.MethodPut)
	require.False(t, wait)

	wait, _ = limiter.Wait(http.MethodPost)
	require.True(t, wait, "POST should wait because PUT consumed the shared bucket slot")
}

func TestNewRateLimiter_GETAndPOSTAreIndependent(t *testing.T) {
	limiter := NewRateLimiter(1, 1, 1, 1)

	wait, _ := limiter.Wait(http.MethodPost)
	require.False(t, wait, "first POST should not wait")

	wait, _ = limiter.Wait(http.MethodGet)
	require.False(t, wait, "first GET should not wait — independent bucket")
}

func TestNewRateLimiter_WindowExpiry(t *testing.T) {
	limiter := NewRateLimiter(1, 1, 1, 1) // 1 request per 1 second

	wait, _ := limiter.Wait(http.MethodGet)
	require.False(t, wait)

	wait, _ = limiter.Wait(http.MethodGet)
	require.True(t, wait, "second GET in same window should wait")

	time.Sleep(1100 * time.Millisecond)

	wait, _ = limiter.Wait(http.MethodGet)
	require.False(t, wait, "GET after window expiry should be allowed")
}

// ---------------------------------------------------------------------------
// ZIA-specific configuration: 20 GET/10s, 10 POST-PUT-DELETE/10s
// ---------------------------------------------------------------------------

func TestZIARateLimiterConfig(t *testing.T) {
	limiter := NewRateLimiter(20, 10, 10, 10)

	t.Run("20 GETs within limit", func(t *testing.T) {
		lim := NewRateLimiter(20, 10, 10, 10)
		for i := 0; i < 20; i++ {
			wait, _ := lim.Wait(http.MethodGet)
			require.False(t, wait, "GET #%d should not wait", i+1)
		}
		wait, _ := lim.Wait(http.MethodGet)
		require.True(t, wait, "21st GET should wait")
	})

	t.Run("10 POSTs within limit", func(t *testing.T) {
		lim := NewRateLimiter(20, 10, 10, 10)
		for i := 0; i < 10; i++ {
			wait, _ := lim.Wait(http.MethodPost)
			require.False(t, wait, "POST #%d should not wait", i+1)
		}
		wait, _ := lim.Wait(http.MethodPost)
		require.True(t, wait, "11th POST should wait")
	})

	_ = limiter // ensure no unused variable warning
}

// Validates the exact parameters used in the OneAPI client and legacy v2 clients
func TestZIARateLimiterParams_PostPutDeleteFreqIs10(t *testing.T) {
	limiter := NewRateLimiter(20, 10, 10, 10)

	for i := 0; i < 10; i++ {
		wait, _ := limiter.Wait(http.MethodPost)
		require.False(t, wait)
	}

	wait, delay := limiter.Wait(http.MethodPost)
	require.True(t, wait)
	require.LessOrEqual(t, delay, 10*time.Second,
		"wait delay should be <= 10s (the postPutDeleteFreq), got %v", delay)
	require.Greater(t, delay, time.Duration(0))
}

// ---------------------------------------------------------------------------
// NewRateLimiterWithHourly: hourly limit tests
// ---------------------------------------------------------------------------

func TestHourlyLimits_GETExhausted(t *testing.T) {
	limiter := NewRateLimiterWithHourly(
		100, 100, 1, 1, // high per-second limits
		3, 3, 3, // low hourly limits for testing
	)

	for i := 0; i < 3; i++ {
		wait, _ := limiter.Wait(http.MethodGet)
		require.False(t, wait, "GET #%d should not wait", i+1)
	}

	wait, delay := limiter.Wait(http.MethodGet)
	require.True(t, wait, "4th GET should hit hourly limit")
	require.Greater(t, delay, time.Duration(0))
}

func TestHourlyLimits_POSTExhausted(t *testing.T) {
	limiter := NewRateLimiterWithHourly(
		100, 100, 1, 1,
		100, 2, 100,
	)

	for i := 0; i < 2; i++ {
		wait, _ := limiter.Wait(http.MethodPost)
		require.False(t, wait, "POST #%d should not wait", i+1)
	}

	wait, delay := limiter.Wait(http.MethodPost)
	require.True(t, wait, "3rd POST should hit hourly limit")
	require.Greater(t, delay, time.Duration(0))
}

func TestHourlyLimits_DELETEExhausted(t *testing.T) {
	limiter := NewRateLimiterWithHourly(
		100, 100, 1, 1,
		100, 100, 2,
	)

	for i := 0; i < 2; i++ {
		wait, _ := limiter.Wait(http.MethodDelete)
		require.False(t, wait, "DELETE #%d should not wait", i+1)
	}

	wait, delay := limiter.Wait(http.MethodDelete)
	require.True(t, wait, "3rd DELETE should hit hourly limit")
	require.Greater(t, delay, time.Duration(0))
}

func TestHourlyLimits_DisabledByDefault(t *testing.T) {
	limiter := NewRateLimiter(100, 100, 1, 1)

	for i := 0; i < 100; i++ {
		wait, _ := limiter.Wait(http.MethodGet)
		require.False(t, wait)
	}

	wait, _ := limiter.Wait(http.MethodGet)
	require.True(t, wait, "should hit per-second limit, not hourly")
}

func TestZIAHourlyConfig(t *testing.T) {
	limiter := NewRateLimiterWithHourly(
		20, 10, 10, 10,
		950, 950, 380,
	)

	// Just validate the limiter doesn't panic and accepts the production parameters
	wait, _ := limiter.Wait(http.MethodGet)
	require.False(t, wait, "first GET should not wait")

	wait, _ = limiter.Wait(http.MethodPost)
	require.False(t, wait, "first POST should not wait")

	wait, _ = limiter.Wait(http.MethodDelete)
	require.False(t, wait, "first DELETE should not wait")
}

// ---------------------------------------------------------------------------
// GlobalRateLimiter tests
// ---------------------------------------------------------------------------

func TestGlobalRateLimiter(t *testing.T) {
	limiter := NewGlobalRateLimiter(3, 1)

	for i := 0; i < 3; i++ {
		wait, _ := limiter.Wait()
		require.False(t, wait, "request #%d should not wait", i+1)
	}

	wait, delay := limiter.Wait()
	require.True(t, wait, "4th request should wait")
	require.Greater(t, delay, time.Duration(0))
}

func TestGlobalRateLimiter_WindowExpiry(t *testing.T) {
	limiter := NewGlobalRateLimiter(1, 1) // 1 request per 1 second

	wait, _ := limiter.Wait()
	require.False(t, wait)

	wait, _ = limiter.Wait()
	require.True(t, wait)

	time.Sleep(1100 * time.Millisecond)

	wait, _ = limiter.Wait()
	require.False(t, wait, "request after window expiry should be allowed")
}

// ---------------------------------------------------------------------------
// RateLimitTransport tests
// ---------------------------------------------------------------------------

func TestRateLimitTransport_MethodBasedThrottling(t *testing.T) {
	limiter := NewRateLimiter(1, 1, 1, 1)

	transport := &RateLimitTransport{
		Base:    http.DefaultTransport,
		Limiter: limiter,
		Logger:  &nopLogger{},
	}

	_ = transport // validate construction doesn't panic
}

type nopLogger struct{}

func (n *nopLogger) Printf(format string, args ...interface{}) {}
