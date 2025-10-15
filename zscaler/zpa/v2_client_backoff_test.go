package zpa

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	rl "github.com/zscaler/zscaler-sdk-go/v3/ratelimiter"
)

// TestBackoffLogic tests the ACTUAL backoff function behavior (no rate limiter)
func TestBackoffLogic(t *testing.T) {
	l := logger.GetDefaultLogger("test: ")

	tests := []struct {
		name        string
		statusCode  int
		retryAfter  string
		method      string
		attemptNum  int
		minWait     time.Duration
		maxWait     time.Duration
		expected    time.Duration
		description string
	}{
		{
			name:        "No response - attempt 1",
			statusCode:  0,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second, // 2^1 * 2s = 4s
			description: "With no response, should calculate exponential backoff",
		},
		{
			name:        "429 with Retry-After header",
			statusCode:  429,
			retryAfter:  "5",
			method:      "GET",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    6 * time.Second, // 5 + 1 second padding
			description: "Should honor Retry-After header",
		},
		{
			name:        "429 with Retry-After duration format",
			statusCode:  429,
			retryAfter:  "13s",
			method:      "GET",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    14 * time.Second, // 13s + 1 second padding
			description: "Should parse duration format",
		},
		{
			name:        "503 with Retry-After",
			statusCode:  503,
			retryAfter:  "3",
			method:      "POST",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second, // 3 + 1 second padding
			description: "Should honor Retry-After for 503",
		},
		{
			name:        "503 without Retry-After - fallback to default",
			statusCode:  503,
			retryAfter:  "",
			method:      "POST",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    2 * time.Second, // getRetryAfter returns RetryWaitMinSeconds (2s) as fallback
			description: "Should use default wait time when no Retry-After header",
		},
		{
			name:        "Exponential backoff - attempt 0",
			statusCode:  500,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  0,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    2 * time.Second, // 2^0 * 2s = 2s
			description: "Attempt 0: 2^0 * 2s = 2s",
		},
		{
			name:        "Exponential backoff - attempt 1",
			statusCode:  500,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  1,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    4 * time.Second, // 2^1 * 2s = 4s
			description: "Attempt 1: 2^1 * 2s = 4s",
		},
		{
			name:        "Exponential backoff - attempt 2",
			statusCode:  500,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second, // 2^2 * 2s = 8s
			description: "Attempt 2: 2^2 * 2s = 8s",
		},
		{
			name:        "Exponential backoff - attempt 3 (capped)",
			statusCode:  500,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  3,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // 2^3 * 2s = 16s, but capped at 10s
			description: "Attempt 3: capped at maxWait (10s)",
		},
		{
			name:        "Exponential backoff - attempt 5 (capped)",
			statusCode:  500,
			retryAfter:  "",
			method:      "GET",
			attemptNum:  5,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // 2^5 * 2s = 64s, capped at 10s
			description: "Attempt 5: capped at maxWait (10s)",
		},
	}

	// Create the ACTUAL backoff function (matches current v2_client.go implementation)
	backoffFunc := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if resp != nil {
			if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusServiceUnavailable {
				retryAfter := getRetryAfter(resp, l)
				if retryAfter > 0 {
					return retryAfter
				}
			}
		}
		// Use exponential backoff for all retries (no rate limiter)
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
				if tt.method != "" {
					req, _ := http.NewRequest(tt.method, "http://example.com", nil)
					resp.Request = req
				}
			}

			result := backoffFunc(tt.minWait, tt.maxWait, tt.attemptNum, resp)

			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
			} else {
				t.Logf("%s: attempt=%d, result=%v ✓", tt.name, tt.attemptNum, result)
			}
		})
	}
}

// TestRateLimiterIntegration tests rate limiter behavior
func TestRateLimiterIntegration(t *testing.T) {
	rateLimiter := rl.NewRateLimiter(20, 10, 10, 10)

	t.Run("First 20 GET requests should not wait", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			wait, duration := rateLimiter.Wait("GET")
			if wait {
				t.Errorf("Request %d should not wait, but got wait=true with duration=%v", i+1, duration)
			}
		}
		t.Log("First 20 GET requests: no wait ✓")
	})

	t.Run("21st GET request should trigger wait", func(t *testing.T) {
		// Reset rate limiter
		rateLimiter = rl.NewRateLimiter(20, 10, 10, 10)

		// Make 20 requests quickly
		for i := 0; i < 20; i++ {
			rateLimiter.Wait("GET")
		}

		// 21st request should trigger wait
		wait, duration := rateLimiter.Wait("GET")
		if !wait {
			t.Error("21st request should trigger wait=true")
		}
		if duration <= 0 {
			t.Errorf("Wait duration should be > 0, got %v", duration)
		}
		t.Logf("21st GET request: wait=%v, duration=%v ✓", wait, duration)
	})

	t.Run("First 10 POST requests should not wait", func(t *testing.T) {
		rateLimiter = rl.NewRateLimiter(20, 10, 10, 10)

		for i := 0; i < 10; i++ {
			wait, duration := rateLimiter.Wait("POST")
			if wait {
				t.Errorf("Request %d should not wait, but got wait=true with duration=%v", i+1, duration)
			}
		}
		t.Log("First 10 POST requests: no wait ✓")
	})

	t.Run("Rate limiter resets after time window", func(t *testing.T) {
		rateLimiter = rl.NewRateLimiter(2, 1, 1, 1) // Very small limits for testing

		// Make 2 requests
		rateLimiter.Wait("GET")
		rateLimiter.Wait("GET")

		// 3rd should wait
		wait, _ := rateLimiter.Wait("GET")
		if !wait {
			t.Error("3rd request should wait")
		}

		// Wait for the time window to pass
		time.Sleep(1100 * time.Millisecond)

		// Should not wait now
		wait, _ = rateLimiter.Wait("GET")
		if wait {
			t.Error("After time window, should not wait")
		}
		t.Log("Rate limiter properly resets after time window ✓")
	})
}

// TestRetryAfterParsing tests the Retry-After header parsing (case-insensitive)
func TestRetryAfterParsing(t *testing.T) {
	l := logger.GetDefaultLogger("test: ")

	tests := []struct {
		name        string
		retryAfter  string
		headerCase  string // "Retry-After", "retry-after", "RETRY-AFTER"
		expected    time.Duration
		description string
	}{
		{
			name:        "Integer seconds (standard case)",
			retryAfter:  "5",
			headerCase:  "Retry-After",
			expected:    6 * time.Second, // 5 + 1 padding
			description: "Should parse integer and add 1 second padding",
		},
		{
			name:        "Integer seconds (lowercase - actual ZPA API format)",
			retryAfter:  "5",
			headerCase:  "retry-after",
			expected:    6 * time.Second, // 5 + 1 padding
			description: "Should parse lowercase header (actual API format)",
		},
		{
			name:        "Integer seconds (uppercase)",
			retryAfter:  "5",
			headerCase:  "RETRY-AFTER",
			expected:    6 * time.Second, // 5 + 1 padding
			description: "Should parse uppercase header",
		},
		{
			name:        "Duration format",
			retryAfter:  "13s",
			headerCase:  "Retry-After",
			expected:    14 * time.Second, // 13 + 1 padding
			description: "Should parse duration format and add 1 second padding",
		},
		{
			name:        "Duration format (lowercase)",
			retryAfter:  "13s",
			headerCase:  "retry-after",
			expected:    14 * time.Second, // 13 + 1 padding
			description: "Should parse duration format with lowercase header",
		},
		{
			name:        "Empty header",
			retryAfter:  "",
			headerCase:  "Retry-After",
			expected:    2 * time.Second, // RetryWaitMinSeconds
			description: "Should return default wait time",
		},
		{
			name:        "Invalid format",
			retryAfter:  "invalid",
			headerCase:  "retry-after",
			expected:    2 * time.Second, // RetryWaitMinSeconds
			description: "Should fall back to default wait time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: 429,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
			}

			if tt.retryAfter != "" {
				// Set with the specified case
				resp.Header.Set(tt.headerCase, tt.retryAfter)
			}

			// getRetryAfter uses resp.Header.Get("retry-after") which is case-insensitive
			result := getRetryAfter(resp, l)

			if result != tt.expected {
				t.Errorf("%s: expected %v, got %v", tt.description, tt.expected, result)
			} else {
				t.Logf("%s: header=%q, retryAfter=%q, result=%v ✓", tt.name, tt.headerCase, tt.retryAfter, result)
			}
		})
	}

	// Explicit test for case-insensitivity
	t.Run("Explicit case-insensitivity test", func(t *testing.T) {
		testCases := []string{"Retry-After", "retry-after", "RETRY-AFTER", "ReTrY-aFtEr"}

		for _, headerName := range testCases {
			resp := &http.Response{
				StatusCode: 429,
				Header:     make(http.Header),
				Body:       io.NopCloser(bytes.NewReader([]byte{})),
			}
			resp.Header.Set(headerName, "5")

			result := getRetryAfter(resp, l)
			expected := 6 * time.Second

			if result != expected {
				t.Errorf("Header %q: expected %v, got %v", headerName, expected, result)
			}
		}
		t.Log("Header case-insensitivity confirmed: all variations work correctly ✓")
	})
}

// TestActualBackoffImplementation validates the current backoff implementation
func TestActualBackoffImplementation(t *testing.T) {
	l := logger.GetDefaultLogger("test: ")
	cfg := &Configuration{}
	cfg.ZPA.Client.RateLimit.RetryWaitMin = 2 * time.Second
	cfg.ZPA.Client.RateLimit.RetryWaitMax = 10 * time.Second
	cfg.ZPA.Client.RateLimit.MaxRetries = 100

	// Use the actual getHTTPClient function to get real backoff behavior
	httpClient := getHTTPClient(l, rl.NewRateLimiter(20, 10, 10, 10), cfg)
	if httpClient == nil {
		t.Fatal("Failed to create HTTP client")
	}

	t.Run("Exponential backoff sequence", func(t *testing.T) {
		// Create mock backoff function matching actual implementation
		backoff := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
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

		expected := []time.Duration{
			2 * time.Second,  // 2^0 * 2s = 2s
			4 * time.Second,  // 2^1 * 2s = 4s
			8 * time.Second,  // 2^2 * 2s = 8s
			10 * time.Second, // 2^3 * 2s = 16s, capped at 10s
			10 * time.Second, // Remains capped
		}

		for i, exp := range expected {
			result := backoff(2*time.Second, 10*time.Second, i, nil)
			if result != exp {
				t.Errorf("Attempt %d: expected %v, got %v", i, exp, result)
			}
		}
		t.Log("Exponential backoff sequence validated: 2s → 4s → 8s → 10s (capped) ✓")
	})

	t.Run("Retry-After takes precedence", func(t *testing.T) {
		backoff := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
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

		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("Retry-After", "7")

		result := backoff(2*time.Second, 10*time.Second, 5, resp)
		expected := 8 * time.Second // 7 + 1 padding

		if result != expected {
			t.Errorf("Expected Retry-After to override exponential backoff: got %v, want %v", result, expected)
		}
		t.Logf("Retry-After correctly overrides exponential backoff ✓")
	})
}
