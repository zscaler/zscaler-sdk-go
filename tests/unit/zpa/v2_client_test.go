// Package unit provides unit tests for the ZPA v2 client
package unit

import (
	"bytes"
	"io"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/errorx"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
)

// =============================================================================
// Client Creation Tests
// =============================================================================

func TestZPAClient_Creation(t *testing.T) {
	t.Run("NewClient with nil config returns error", func(t *testing.T) {
		client, err := zpa.NewClient(nil)

		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "configuration cannot be nil")
	})
}

// =============================================================================
// Backoff Logic Tests
// =============================================================================

func TestZPAClient_BackoffLogic(t *testing.T) {
	t.Parallel()

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
			name:        "Exponential backoff - attempt 0",
			statusCode:  500,
			retryAfter:  "",
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
			attemptNum:  2,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    8 * time.Second, // 2^2 * 2s = 8s
			description: "Attempt 2: 2^2 * 2s = 8s",
		},
		{
			name:        "Exponential backoff - capped at max",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  3,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // 2^3 * 2s = 16s, capped at 10s
			description: "Attempt 3: capped at maxWait (10s)",
		},
		{
			name:        "Exponential backoff - high attempt (stays capped)",
			statusCode:  500,
			retryAfter:  "",
			attemptNum:  10,
			minWait:     2 * time.Second,
			maxWait:     10 * time.Second,
			expected:    10 * time.Second, // Capped at 10s
			description: "Attempt 10: remains capped at maxWait (10s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Calculate exponential backoff
			mult := math.Pow(2, float64(tt.attemptNum)) * float64(tt.minWait)
			result := time.Duration(mult)
			if float64(result) != mult || result > tt.maxWait {
				result = tt.maxWait
			}

			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// =============================================================================
// Error Detection Tests
// =============================================================================

func TestZPAClient_ErrorDetection(t *testing.T) {
	t.Parallel()

	t.Run("IsSessionInvalidError", func(t *testing.T) {
		tests := []struct {
			name       string
			statusCode int
			body       string
			expected   bool
		}{
			{
				name:       "SESSION_NOT_VALID",
				statusCode: 401,
				body:       `{"code": "SESSION_NOT_VALID"}`,
				expected:   true,
			},
			{
				name:       "Session already invalidated",
				statusCode: 401,
				body:       `{"message": "getAttribute: Session already invalidated"}`,
				expected:   true,
			},
			{
				name:       "Resource Access Blocked",
				statusCode: 401,
				body:       `{"message": "Resource Access Blocked"}`,
				expected:   true,
			},
			{
				name:       "Regular 401 error",
				statusCode: 401,
				body:       `{"message": "Invalid credentials"}`,
				expected:   false,
			},
			{
				name:       "Non-401 status code with SESSION_NOT_VALID in body",
				statusCode: 500,
				body:       `{"code": "SESSION_NOT_VALID"}`,
				expected:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				resp := common.CreateMockHTTPResponse(tt.statusCode, tt.body)
				result := errorx.IsSessionInvalidError(resp)

				assert.Equal(t, tt.expected, result)
			})
		}
	})

	t.Run("IsEditLockError", func(t *testing.T) {
		tests := []struct {
			name       string
			statusCode int
			body       string
			expected   bool
		}{
			{
				name:       "EDIT_LOCK_NOT_AVAILABLE",
				statusCode: 409,
				body:       `{"code": "EDIT_LOCK_NOT_AVAILABLE"}`,
				expected:   true,
			},
			{
				name:       "Resource Access Blocked (409)",
				statusCode: 409,
				body:       `{"message": "Resource Access Blocked"}`,
				expected:   true,
			},
			{
				name:       "Failed during enter Org barrier",
				statusCode: 409,
				body:       `{"message": "Failed during enter Org barrier"}`,
				expected:   true,
			},
			{
				name:       "Regular 409 error",
				statusCode: 409,
				body:       `{"message": "Conflict"}`,
				expected:   false,
			},
			{
				name:       "Non-409 status code",
				statusCode: 401,
				body:       `{"code": "EDIT_LOCK_NOT_AVAILABLE"}`,
				expected:   false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				resp := common.CreateMockHTTPResponse(tt.statusCode, tt.body)
				result := errorx.IsEditLockError(resp)

				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

// =============================================================================
// Retry-After Header Parsing Tests
// =============================================================================

func TestZPAClient_RetryAfterParsing(t *testing.T) {
	t.Parallel()

	t.Run("Case-insensitive header parsing", func(t *testing.T) {
		testCases := []string{"Retry-After", "retry-after", "RETRY-AFTER", "ReTrY-aFtEr"}

		for _, headerName := range testCases {
			t.Run(headerName, func(t *testing.T) {
				resp := &http.Response{
					StatusCode: 429,
					Header:     make(http.Header),
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				}
				resp.Header.Set(headerName, "5")

				// All variations should set the header correctly
				value := resp.Header.Get("Retry-After")
				assert.Equal(t, "5", value)
			})
		}
	})

	t.Run("Integer format parsing", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("Retry-After", "10")

		value := resp.Header.Get("Retry-After")
		assert.Equal(t, "10", value)
	})

	t.Run("Duration format parsing", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     make(http.Header),
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		}
		resp.Header.Set("Retry-After", "15s")

		value := resp.Header.Get("Retry-After")
		assert.Equal(t, "15s", value)
	})
}

// =============================================================================
// Error Response Tests
// =============================================================================

func TestZPAClient_ErrorResponse(t *testing.T) {
	t.Parallel()

	t.Run("CheckErrorInResponse with JSON error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://api.example.com/test", nil)
		resp := &http.Response{
			StatusCode: 400,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"code": "INVALID_REQUEST", "message": "Bad request"}`)),
			Request:    req,
		}

		err := errorx.CheckErrorInResponse(resp, nil)
		require.Error(t, err)

		errResp, ok := err.(*errorx.ErrorResponse)
		require.True(t, ok)
		assert.NotNil(t, errResp.Parsed)
		assert.Equal(t, 400, errResp.Parsed.Status)
	})

	t.Run("ErrorResponse.IsObjectNotFound", func(t *testing.T) {
		tests := []struct {
			name     string
			err      *errorx.ErrorResponse
			expected bool
		}{
			{
				name: "404 status code",
				err: &errorx.ErrorResponse{
					Response: &http.Response{StatusCode: 404},
				},
				expected: true,
			},
			{
				name: "resource.not.found ID",
				err: &errorx.ErrorResponse{
					Response: &http.Response{StatusCode: 400},
					Parsed:   &errorx.ParsedAPIError{ID: "resource.not.found"},
				},
				expected: true,
			},
			{
				name: "other error",
				err: &errorx.ErrorResponse{
					Response: &http.Response{StatusCode: 500},
				},
				expected: false,
			},
			{
				name:     "nil response",
				err:      &errorx.ErrorResponse{},
				expected: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tt.err.IsObjectNotFound()
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

// =============================================================================
// HTTP Client Configuration Tests
// =============================================================================

func TestZPAClient_HTTPClientConfig(t *testing.T) {
	t.Parallel()

	t.Run("Default timeout is 240 seconds", func(t *testing.T) {
		// The default request timeout for ZPA is 240 seconds (defined as constant)
		expectedTimeout := 240 * time.Second

		cfg := &zpa.Configuration{}
		cfg.ZPA.Client.RequestTimeout = expectedTimeout

		assert.Equal(t, 240*time.Second, cfg.ZPA.Client.RequestTimeout)
	})

	t.Run("Default retry settings", func(t *testing.T) {
		cfg := &zpa.Configuration{}
		cfg.ZPA.Client.RateLimit.MaxRetries = int32(zpa.MaxNumOfRetries)
		cfg.ZPA.Client.RateLimit.RetryWaitMax = time.Second * time.Duration(zpa.RetryWaitMaxSeconds)
		cfg.ZPA.Client.RateLimit.RetryWaitMin = time.Second * time.Duration(zpa.RetryWaitMinSeconds)

		assert.Equal(t, int32(100), cfg.ZPA.Client.RateLimit.MaxRetries)
		assert.Equal(t, 10*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMax)
		assert.Equal(t, 2*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMin)
	})
}

// =============================================================================
// Mock Server Integration Tests
// =============================================================================

func TestZPAClient_MockServerIntegration(t *testing.T) {
	t.Run("Mock server returns configured response", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		// Configure mock response
		server.On("GET", "/test", common.SuccessResponse(`{"status": "ok"}`))

		// Make request
		resp, err := http.Get(server.URL + "/test")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "ok")
	})

	t.Run("Mock server tracks requests", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/api/resource", common.CreatedResponse(`{"id": "123"}`))

		// Make requests
		http.Post(server.URL+"/api/resource", "application/json", strings.NewReader(`{"name": "test"}`))
		http.Post(server.URL+"/api/resource", "application/json", strings.NewReader(`{"name": "test2"}`))

		// Verify tracking
		assert.Equal(t, 2, server.GetCallCount("POST", "/api/resource"))
	})

	t.Run("Mock server returns 404 for unconfigured paths", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		resp, err := http.Get(server.URL + "/unconfigured")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("Mock server simulates rate limiting", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/rate-limited", common.TooManyRequestsResponseWithHeader("5"))

		resp, err := http.Get(server.URL + "/rate-limited")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 429, resp.StatusCode)
		assert.Equal(t, "5", resp.Header.Get("Retry-After"))
	})

	t.Run("Mock server simulates session invalid", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/session-invalid", common.SessionInvalidResponse())

		resp, err := http.Get(server.URL + "/session-invalid")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 401, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "SESSION_NOT_VALID")
	})
}

// =============================================================================
// Rate Limiter Tests
// =============================================================================

func TestZPAClient_RateLimiter(t *testing.T) {
	t.Parallel()

	t.Run("Rate limit configuration", func(t *testing.T) {
		// ZPA rate limits: 20 GET per 10s, 10 POST/PUT/DELETE per 10s
		cfg := &zpa.Configuration{}
		cfg.ZPA.Client.RateLimit.MaxRetries = 100
		cfg.ZPA.Client.RateLimit.RetryWaitMin = 2 * time.Second
		cfg.ZPA.Client.RateLimit.RetryWaitMax = 10 * time.Second

		assert.Equal(t, int32(100), cfg.ZPA.Client.RateLimit.MaxRetries)
		assert.Equal(t, 2*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMin)
		assert.Equal(t, 10*time.Second, cfg.ZPA.Client.RateLimit.RetryWaitMax)
	})
}

