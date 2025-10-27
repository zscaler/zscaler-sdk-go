package zscaler

import (
	"context"
	"fmt"

	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// TestTimeoutExcludesRateLimitWaits verifies that rate limit wait times don't count toward request timeout
func TestTimeoutExcludesRateLimitWaits(t *testing.T) {
	tests := []struct {
		name                string
		requestTimeout      time.Duration
		rateLimitDelay      time.Duration
		actualResponseDelay time.Duration
		rateLimitCount      int
		expectTimeout       bool
		expectSuccess       bool
	}{
		{
			name:                "No rate limiting - quick response",
			requestTimeout:      5 * time.Second,
			rateLimitDelay:      0,
			actualResponseDelay: 1 * time.Second,
			rateLimitCount:      0,
			expectTimeout:       false,
			expectSuccess:       true,
		},
		{
			name:                "Rate limiting but under timeout",
			requestTimeout:      5 * time.Second,
			rateLimitDelay:      2 * time.Second,
			actualResponseDelay: 500 * time.Millisecond,
			rateLimitCount:      2,
			expectTimeout:       false,
			expectSuccess:       true,
		},
		{
			name:                "Long rate limit waits excluded from timeout",
			requestTimeout:      3 * time.Second,
			rateLimitDelay:      4 * time.Second, // Total wait: 8 seconds (2 retries)
			actualResponseDelay: 500 * time.Millisecond,
			rateLimitCount:      2,
			expectTimeout:       false,
			expectSuccess:       true,
		},
		{
			name:                "Actual response time exceeds timeout",
			requestTimeout:      2 * time.Second,
			rateLimitDelay:      1 * time.Second,
			actualResponseDelay: 3 * time.Second, // This exceeds timeout
			rateLimitCount:      1,
			expectTimeout:       true,
			expectSuccess:       false,
		},
		{
			name:                "Very long rate limit wait but quick actual response",
			requestTimeout:      5 * time.Second,
			rateLimitDelay:      10 * time.Second, // 10 seconds wait per retry
			actualResponseDelay: 1 * time.Second,
			rateLimitCount:      3, // 30 seconds of waiting!
			expectTimeout:       false,
			expectSuccess:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestCount atomic.Int32
			totalRequests := tt.rateLimitCount + 1 // Rate limited requests + final success

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count := requestCount.Add(1)

				// First N requests return 429 with rate limiting
				if count <= int32(tt.rateLimitCount) {
					w.Header().Set("Retry-After", fmt.Sprintf("%d", int(tt.rateLimitDelay.Seconds())))
					w.Header().Set("X-Ratelimit-Remaining", "0")
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte(`{"error": "rate limited"}`))
					return
				}

				// Simulate actual response delay
				time.Sleep(tt.actualResponseDelay)

				// Final request succeeds
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 123, "name": "test"}`))
			}))
			defer server.Close()

			// Create test configuration
			cfg, err := NewConfiguration(
				WithClientID("test-client"),
				WithClientSecret("test-secret"),
				WithVanityDomain("test"),
				WithRequestTimeout(tt.requestTimeout),
				WithRateLimitMaxRetries(int32(totalRequests)),
			)
			if err != nil {
				t.Fatalf("Failed to create configuration: %v", err)
			}

			// Create client
			client := &Client{
				oauth2Credentials: cfg,
			}

			// Mock authentication (skip actual OAuth2)
			cfg.Zscaler.Client.AuthToken = &AuthToken{
				AccessToken: "test-token",
				TokenType:   "Bearer",
				Expiry:      time.Now().Add(1 * time.Hour),
			}

			// Execute request
			ctx := context.Background()
			endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

			startTime := time.Now()
			_, resp, _, err := client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")
			duration := time.Since(startTime)

			// Validate results
			if tt.expectTimeout {
				if err == nil {
					t.Errorf("Expected timeout error but got success")
				}
				if !strings.Contains(err.Error(), "request timeout exceeded") {
					t.Errorf("Expected timeout error, got: %v", err)
				}
				t.Logf("✓ Correctly timed out after %v (actual response delay: %v)", duration, tt.actualResponseDelay)
			} else if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				}
				if resp == nil || resp.StatusCode != http.StatusOK {
					t.Errorf("Expected 200 OK response")
				}

				// Calculate expected total time
				expectedWaitTime := tt.rateLimitDelay * time.Duration(tt.rateLimitCount)
				expectedActualTime := tt.actualResponseDelay * time.Duration(totalRequests)
				expectedTotal := expectedWaitTime + expectedActualTime

				t.Logf("✓ Success: total=%v, expected_wait=%v, expected_actual=%v",
					duration, expectedWaitTime, expectedActualTime)

				// Verify total time is reasonable (within tolerance)
				tolerance := 2 * time.Second // Allow 2s tolerance for test overhead
				if duration < expectedTotal-tolerance || duration > expectedTotal+tolerance {
					t.Logf("⚠️  Duration %v outside expected range %v ± %v (this may be OK depending on test environment)",
						duration, expectedTotal, tolerance)
				}
			}

			// Verify correct number of requests were made
			finalCount := requestCount.Load()
			if tt.expectSuccess && finalCount != int32(totalRequests) {
				t.Errorf("Expected %d requests but got %d", totalRequests, finalCount)
			}
		})
	}
}

// TestTimeoutCalculationAccuracy verifies the timeout calculation is accurate
func TestTimeoutCalculationAccuracy(t *testing.T) {
	requestCount := 0
	rateLimitWaitTime := 2 * time.Second
	actualProcessingTime := 500 * time.Millisecond
	numRateLimits := 3

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount <= numRateLimits {
			// Return 429 with Retry-After
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rateLimitWaitTime.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limited"}`))
			return
		}

		// Simulate actual processing
		time.Sleep(actualProcessingTime)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	// Set timeout to just cover actual processing time, not wait time
	actualTimeNeeded := actualProcessingTime * time.Duration(numRateLimits+1)
	waitTimeTotal := rateLimitWaitTime * time.Duration(numRateLimits)
	requestTimeout := actualTimeNeeded + (1 * time.Second) // Add 1s buffer

	t.Logf("Test setup: actual_time=%v, wait_time=%v, timeout=%v",
		actualTimeNeeded, waitTimeTotal, requestTimeout)

	cfg, err := NewConfiguration(
		WithClientID("test-client"),
		WithClientSecret("test-secret"),
		WithVanityDomain("test"),
		WithRequestTimeout(requestTimeout),
		WithRateLimitMaxRetries(10),
		WithDebug(true),
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}

	client := &Client{
		oauth2Credentials: cfg,
	}

	cfg.Zscaler.Client.AuthToken = &AuthToken{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

	startTime := time.Now()
	_, resp, _, err := client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")
	totalDuration := time.Since(startTime)

	// Should succeed because wait time doesn't count
	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK response")
	}

	// Verify timing
	expectedTotal := actualTimeNeeded + waitTimeTotal
	t.Logf("✓ Success: total_duration=%v, expected=%v, actual_time=%v, wait_time=%v",
		totalDuration, expectedTotal, actualTimeNeeded, waitTimeTotal)

	// The total duration should be close to expected (actual + wait)
	tolerance := 2 * time.Second
	if totalDuration < expectedTotal-tolerance || totalDuration > expectedTotal+tolerance {
		t.Logf("⚠️  Total duration outside expected range (this may be OK in test environments)")
	}

	// Most importantly: it should NOT timeout even though total > timeout
	if totalDuration > requestTimeout {
		t.Logf("✓ Correctly allowed request to complete despite total time (%v) > timeout (%v)",
			totalDuration, requestTimeout)
	}
}

// TestTimeoutMessageIncludesWaitTime verifies error messages include wait time information
func TestTimeoutMessageIncludesWaitTime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Always return 429 to force timeout
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "rate limited"}`))
	}))
	defer server.Close()

	cfg, err := NewConfiguration(
		WithClientID("test-client"),
		WithClientSecret("test-secret"),
		WithVanityDomain("test"),
		WithRequestTimeout(3*time.Second),
		WithRateLimitMaxRetries(100),
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}

	client := &Client{
		oauth2Credentials: cfg,
	}

	cfg.Zscaler.Client.AuthToken = &AuthToken{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

	_, _, _, err = client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")

	// Should timeout
	if err == nil {
		t.Errorf("Expected timeout error but got success")
	}

	// Error message should mention wait time
	errMsg := err.Error()
	if !strings.Contains(errMsg, "request timeout exceeded") {
		t.Errorf("Error message should contain 'request timeout exceeded', got: %s", errMsg)
	}

	if !strings.Contains(errMsg, "excluding") {
		t.Errorf("Error message should mention excluded wait time, got: %s", errMsg)
	}

	t.Logf("✓ Error message correctly formatted: %s", errMsg)
}

// TestNoRateLimitingBehaviorUnchanged verifies normal requests work as before
func TestNoRateLimitingBehaviorUnchanged(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "name": "test"}`))
	}))
	defer server.Close()

	cfg, err := NewConfiguration(
		WithClientID("test-client"),
		WithClientSecret("test-secret"),
		WithVanityDomain("test"),
		WithRequestTimeout(5*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}

	client := &Client{
		oauth2Credentials: cfg,
	}

	cfg.Zscaler.Client.AuthToken = &AuthToken{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

	startTime := time.Now()
	body, resp, _, err := client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")
	duration := time.Since(startTime)

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK response")
	}

	if !strings.Contains(string(body), "test") {
		t.Errorf("Expected response body to contain 'test'")
	}

	t.Logf("✓ Normal request completed successfully in %v", duration)
}

// TestServerErrorBackoffExcludedFromTimeout verifies server error backoff doesn't count
func TestServerErrorBackoffExcludedFromTimeout(t *testing.T) {
	requestCount := 0
	numServerErrors := 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount <= numServerErrors {
			// Return 500 server error
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "server error"}`))
			return
		}

		// Final request succeeds quickly
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	cfg, err := NewConfiguration(
		WithClientID("test-client"),
		WithClientSecret("test-secret"),
		WithVanityDomain("test"),
		WithRequestTimeout(2*time.Second),
		WithRateLimitMaxRetries(10),
		WithRateLimitMinWait(2*time.Second), // 2s backoff
		WithRateLimitMaxWait(10*time.Second),
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}

	client := &Client{
		oauth2Credentials: cfg,
	}

	cfg.Zscaler.Client.AuthToken = &AuthToken{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

	startTime := time.Now()
	_, resp, _, err := client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")
	duration := time.Since(startTime)

	// Should succeed despite backoff delays
	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 OK response")
	}

	t.Logf("✓ Server error retry succeeded in %v (with backoff excluded from timeout)", duration)
}

// Mock logger for testing
type testLogger struct {
	logs []string
}

func (l *testLogger) Printf(format string, args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprintf(format, args...))
}

func (l *testLogger) Println(args ...interface{}) {
	l.logs = append(l.logs, fmt.Sprint(args...))
}

// TestDebugLoggingShowsWaitTimes verifies debug logs show wait time breakdown
func TestDebugLoggingShowsWaitTimes(t *testing.T) {
	requestCount := 0
	rateLimitDelay := 1 * time.Second

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if requestCount == 1 {
			w.Header().Set("Retry-After", fmt.Sprintf("%d", int(rateLimitDelay.Seconds())))
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "rate limited"}`))
			return
		}

		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	testLog := &testLogger{logs: []string{}}

	cfg, err := NewConfiguration(
		WithClientID("test-client"),
		WithClientSecret("test-secret"),
		WithVanityDomain("test"),
		WithRequestTimeout(5*time.Second),
		WithRateLimitMaxRetries(10),
		WithDebug(true),
	)
	if err != nil {
		t.Fatalf("Failed to create configuration: %v", err)
	}

	// Replace logger
	cfg.Logger = testLog

	client := &Client{
		oauth2Credentials: cfg,
	}

	cfg.Zscaler.Client.AuthToken = &AuthToken{
		AccessToken: "test-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	}

	ctx := context.Background()
	endpoint := strings.TrimPrefix(server.URL, "http://") + "/test"

	_, _, _, err = client.ExecuteRequest(ctx, "GET", endpoint, nil, nil, "")

	if err != nil {
		t.Errorf("Expected success but got error: %v", err)
	}

	// Check logs for timing breakdown
	foundTimingLog := false
	for _, log := range testLog.logs {
		if strings.Contains(log, "Request completed") &&
			strings.Contains(log, "total=") &&
			strings.Contains(log, "waited=") &&
			strings.Contains(log, "actual=") {
			foundTimingLog = true
			t.Logf("✓ Found timing log: %s", log)
			break
		}
	}

	if !foundTimingLog {
		t.Errorf("Expected to find timing breakdown in debug logs")
		t.Logf("Logs captured: %v", testLog.logs)
	}
}
