package errorx

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// newResponse builds an *http.Response with the given status, content type, and
// body, plus an associated request so error helpers can read Request fields.
func newResponse(status int, contentType, body, method, rawURL string) *http.Response {
	u, _ := url.Parse(rawURL)
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d", status),
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request: &http.Request{
			Method: method,
			URL:    u,
		},
	}
}

func TestIsLimitExceeded_True(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
		Parsed:   &ParsedAPIError{Code: "LIMIT_EXCEEDED", Message: "Maximum 100 static IPs are allowed. Limit has exceeded."},
	}
	require.True(t, errResp.IsLimitExceeded(), "should detect LIMIT_EXCEEDED on 403")
}

func TestIsLimitExceeded_WrongCode(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
		Parsed:   &ParsedAPIError{Code: "ACCESS_DENIED", Message: "Access denied"},
	}
	require.False(t, errResp.IsLimitExceeded(), "should not match ACCESS_DENIED")
}

func TestIsLimitExceeded_WrongStatus(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusBadRequest},
		Parsed:   &ParsedAPIError{Code: "LIMIT_EXCEEDED", Message: "limit exceeded"},
	}
	require.False(t, errResp.IsLimitExceeded(), "LIMIT_EXCEEDED on non-403 should not match")
}

func TestIsLimitExceeded_NilParsed(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
		Parsed:   nil,
	}
	require.False(t, errResp.IsLimitExceeded(), "nil Parsed should return false")
}

func TestIsLimitExceeded_NilResponse(t *testing.T) {
	errResp := &ErrorResponse{Response: nil}
	require.False(t, errResp.IsLimitExceeded(), "nil Response should return false")
}

func TestIsLimitExceeded_NilErrorResponse(t *testing.T) {
	var errResp *ErrorResponse
	require.False(t, errResp.IsLimitExceeded(), "nil ErrorResponse should return false")
}

func TestIsLimitExceeded_NumericCode(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
		Parsed:   &ParsedAPIError{Code: float64(403), Message: "limit exceeded"},
	}
	require.False(t, errResp.IsLimitExceeded(), "numeric code should not match string assertion")
}

func TestIsObjectNotFound_404(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusNotFound},
	}
	require.True(t, errResp.IsObjectNotFound(), "404 should be object not found")
}

func TestIsObjectNotFound_ResourceNotFoundID(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusBadRequest},
		Parsed:   &ParsedAPIError{ID: "resource.not.found"},
	}
	require.True(t, errResp.IsObjectNotFound(), "resource.not.found ID should match")
}

func TestIsObjectNotFound_Nil(t *testing.T) {
	var errResp *ErrorResponse
	require.False(t, errResp.IsObjectNotFound(), "nil ErrorResponse should return false")
}

func TestPackageIsObjectNotFound_PlainError(t *testing.T) {
	// A plain *errors.errorString (e.g. cancelled/timed-out request) must not
	// panic and must report false, rather than crashing on a type assertion.
	require.False(t, IsObjectNotFound(errors.New("Request cancelled")), "plain error should not match and must not panic")
}

func TestPackageIsObjectNotFound_Nil(t *testing.T) {
	require.False(t, IsObjectNotFound(nil), "nil error should return false")
}

func TestPackageIsObjectNotFound_Match(t *testing.T) {
	var err error = &ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	require.True(t, IsObjectNotFound(err), "404 ErrorResponse should match")
}

func TestPackageIsObjectNotFound_Wrapped(t *testing.T) {
	var inner error = &ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	wrapped := fmt.Errorf("failed to read resource: %w", inner)
	require.True(t, IsObjectNotFound(wrapped), "wrapped ErrorResponse should be detected via errors.As")
}

func TestPackageIsLimitExceeded_PlainError(t *testing.T) {
	require.False(t, IsLimitExceeded(errors.New("boom")), "plain error should not match and must not panic")
}

func TestPackageIsLimitExceeded_Wrapped(t *testing.T) {
	var inner error = &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusForbidden},
		Parsed:   &ParsedAPIError{Code: "LIMIT_EXCEEDED"},
	}
	wrapped := fmt.Errorf("create failed: %w", inner)
	require.True(t, IsLimitExceeded(wrapped), "wrapped LIMIT_EXCEEDED should be detected")
}

func TestAsErrorResponse(t *testing.T) {
	var inner error = &ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}
	got, ok := AsErrorResponse(fmt.Errorf("wrap: %w", inner))
	require.True(t, ok)
	require.NotNil(t, got)

	got, ok = AsErrorResponse(errors.New("plain"))
	require.False(t, ok)
	require.Nil(t, got)

	got, ok = AsErrorResponse(nil)
	require.False(t, ok)
	require.Nil(t, got)
}

// =====================================================
// Error()
// =====================================================

func TestError_Parsed(t *testing.T) {
	errResp := &ErrorResponse{
		Parsed: &ParsedAPIError{Code: "X", Message: "boom", Status: 400},
	}
	out := errResp.Error()
	require.Contains(t, out, "Error:")
	require.Contains(t, out, "boom")
}

func TestError_ResponseOnly(t *testing.T) {
	errResp := &ErrorResponse{
		Response: newResponse(http.StatusInternalServerError, "text/plain", "", http.MethodGet, "https://api.example.com/x"),
		Message:  "server exploded",
	}
	out := errResp.Error()
	require.Contains(t, out, "FAILED")
	require.Contains(t, out, "server exploded")
	require.Contains(t, out, "GET")
}

func TestError_ErrOnly(t *testing.T) {
	errResp := &ErrorResponse{Err: errors.New("network down")}
	require.Contains(t, errResp.Error(), "network down")
}

// =====================================================
// Unwrap()
// =====================================================

func TestUnwrap(t *testing.T) {
	inner := errors.New("root cause")
	errResp := &ErrorResponse{Err: inner}
	require.Equal(t, inner, errResp.Unwrap())
	require.True(t, errors.Is(fmt.Errorf("wrap: %w", errResp), inner))

	var nilResp *ErrorResponse
	require.Nil(t, nilResp.Unwrap())
}

// =====================================================
// CheckErrorInResponse()
// =====================================================

func TestCheckErrorInResponse_Success(t *testing.T) {
	sentinel := errors.New("passthrough")
	res := newResponse(http.StatusOK, "application/json", "{}", http.MethodGet, "https://api.example.com/x")
	require.Equal(t, sentinel, CheckErrorInResponse(res, sentinel))
}

func TestCheckErrorInResponse_JSONError(t *testing.T) {
	body := `{"code":"BAD","message":"invalid","id":"resource.not.found","reason":"nope","exception":"ex"}`
	res := newResponse(http.StatusNotFound, "application/json", body, http.MethodPost, "https://api.example.com/x")

	err := CheckErrorInResponse(res, nil)
	respErr, ok := AsErrorResponse(err)
	require.True(t, ok)
	require.Equal(t, "invalid", respErr.Parsed.Message)
	require.Equal(t, "resource.not.found", respErr.Parsed.ID)
	require.Equal(t, "nope", respErr.Parsed.Reason)
	require.Equal(t, "ex", respErr.Parsed.Exception)
	require.True(t, respErr.IsObjectNotFound())
}

func TestCheckErrorInResponse_InvalidJSON(t *testing.T) {
	res := newResponse(http.StatusBadRequest, "application/json", "{not-json", http.MethodGet, "https://api.example.com/x")

	err := CheckErrorInResponse(res, nil)
	respErr, ok := AsErrorResponse(err)
	require.True(t, ok)
	require.Contains(t, respErr.Parsed.Message, "Failed to parse JSON error body")
}

func TestCheckErrorInResponse_NonJSONGeneric(t *testing.T) {
	res := newResponse(http.StatusInternalServerError, "text/plain", "boom", http.MethodGet, "https://api.example.com/x")

	err := CheckErrorInResponse(res, nil)
	respErr, ok := AsErrorResponse(err)
	require.True(t, ok)
	require.Equal(t, "boom", respErr.Parsed.Message)
}

func TestCheckErrorInResponse_OneAPIFallback(t *testing.T) {
	body := "This resource is accessible only through Zscaler OneAPI"
	res := newResponse(http.StatusForbidden, "text/html", body, http.MethodGet, "https://api.example.com/zia/api/v1/x")

	err := CheckErrorInResponse(res, nil)
	respErr, ok := AsErrorResponse(err)
	require.True(t, ok)
	require.Equal(t, "ONLY_ONEAPI_SUPPORTED", respErr.Parsed.Code)
	require.Equal(t, http.StatusUnauthorized, respErr.Response.StatusCode)
}

// =====================================================
// getBaseURL / NewOneAPIFallbackError / mustParseURL
// =====================================================

func TestGetBaseURL(t *testing.T) {
	u, _ := url.Parse("https://api.example.com/zia/api/v1/x?y=1")
	require.Equal(t, "https://api.example.com", getBaseURL(u))
}

func TestNewOneAPIFallbackError(t *testing.T) {
	err := NewOneAPIFallbackError([]byte("  only oneapi  "), http.MethodPost, "/zia/api/v1/x", "https://api.example.com")
	require.Equal(t, http.StatusUnauthorized, err.Response.StatusCode)
	require.Equal(t, "https://api.example.com/zia/api/v1/x", err.Parsed.URL)
	require.Equal(t, "only oneapi", err.Parsed.Message)
	require.Equal(t, "ONLY_ONEAPI_SUPPORTED", err.Parsed.Code)
}

func TestMustParseURL(t *testing.T) {
	u := mustParseURL("https://api.example.com/x")
	require.Equal(t, "api.example.com", u.Host)

	// A control character forces url.Parse to fail; the helper falls back to a
	// URL whose Path is the raw string.
	bad := mustParseURL("http://\x7f")
	require.Equal(t, "http://\x7f", bad.Path)
}

// =====================================================
// IsObjectNotFound remaining branches
// =====================================================

func TestIsObjectNotFound_NilResponse(t *testing.T) {
	errResp := &ErrorResponse{Response: nil}
	require.False(t, errResp.IsObjectNotFound())
}

func TestIsObjectNotFound_NoMatch(t *testing.T) {
	errResp := &ErrorResponse{
		Response: &http.Response{StatusCode: http.StatusBadRequest},
		Parsed:   &ParsedAPIError{ID: "some.other.error"},
	}
	require.False(t, errResp.IsObjectNotFound())
}

// =====================================================
// IsSessionInvalidError
// =====================================================

func TestIsSessionInvalidError_NonUnauthorized(t *testing.T) {
	res := newResponse(http.StatusOK, "text/plain", "", http.MethodGet, "https://api.example.com/x")
	require.False(t, IsSessionInvalidError(res))
}

func TestIsSessionInvalidError_Match(t *testing.T) {
	res := newResponse(http.StatusUnauthorized, "text/plain", "SESSION_NOT_VALID", http.MethodGet, "https://api.example.com/x")
	require.True(t, IsSessionInvalidError(res))
	// Body must be rewound for reuse.
	body, _ := io.ReadAll(res.Body)
	require.Equal(t, "SESSION_NOT_VALID", string(body))
}

func TestIsSessionInvalidError_NoMatch(t *testing.T) {
	res := newResponse(http.StatusUnauthorized, "text/plain", "some other 401", http.MethodGet, "https://api.example.com/x")
	require.False(t, IsSessionInvalidError(res))
}

// =====================================================
// IsEditLockError
// =====================================================

func TestIsEditLockError_WrongStatus(t *testing.T) {
	res := newResponse(http.StatusOK, "text/plain", "EDIT_LOCK_NOT_AVAILABLE", http.MethodGet, "https://api.example.com/x")
	require.False(t, IsEditLockError(res))
}

func TestIsEditLockError_Conflict(t *testing.T) {
	res := newResponse(http.StatusConflict, "text/plain", "EDIT_LOCK_NOT_AVAILABLE", http.MethodGet, "https://api.example.com/x")
	require.True(t, IsEditLockError(res))
}

func TestIsEditLockError_PreconditionFailed(t *testing.T) {
	res := newResponse(http.StatusPreconditionFailed, "text/plain", "Failed during enter Org barrier", http.MethodGet, "https://api.example.com/x")
	require.True(t, IsEditLockError(res))
}

func TestIsEditLockError_NoMatch(t *testing.T) {
	res := newResponse(http.StatusConflict, "text/plain", "some other conflict", http.MethodGet, "https://api.example.com/x")
	require.False(t, IsEditLockError(res))
}

// =====================================================
// IsRetryableServerError
// =====================================================

func TestIsRetryableServerError(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{
			name:   "deterministic 500 with API code is not retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred, please contact Zscaler's support"}`,
			want:   false,
		},
		{
			name:   "500 with org barrier marker is retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":"UNEXPECTED_ERROR","message":"Failed during enter Org barrier"}`,
			want:   true,
		},
		{
			name:   "500 with edit lock marker is retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":"EDIT_LOCK_NOT_AVAILABLE"}`,
			want:   true,
		},
		{
			name:   "500 with resource access blocked marker is retryable",
			status: http.StatusInternalServerError,
			body:   `{"message":"Resource Access Blocked"}`,
			want:   true,
		},
		{
			name:   "500 with precondition marker is retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":"UNEXPECTED_ERROR","message":"Request processing failed, possibly because an expected precondition was not met"}`,
			want:   true,
		},
		{
			name:   "empty body is retryable",
			status: http.StatusInternalServerError,
			body:   "",
			want:   true,
		},
		{
			name:   "non JSON gateway body is retryable",
			status: http.StatusBadGateway,
			body:   "<html>502 Bad Gateway</html>",
			want:   true,
		},
		{
			name:   "JSON without a code is retryable",
			status: http.StatusServiceUnavailable,
			body:   `{"message":"temporarily unavailable"}`,
			want:   true,
		},
		{
			name:   "non string code is retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":500,"message":"numeric code"}`,
			want:   true,
		},
		{
			name:   "empty string code is retryable",
			status: http.StatusInternalServerError,
			body:   `{"code":"","message":"blank code"}`,
			want:   true,
		},
		{
			name:   "501 Not Implemented is not retryable",
			status: http.StatusNotImplemented,
			body:   "",
			want:   false,
		},
		{
			name:   "non 5xx status is not retryable",
			status: http.StatusBadRequest,
			body:   `{"code":"INVALID_INPUT"}`,
			want:   false,
		},
		{
			name:   "success status is not retryable",
			status: http.StatusOK,
			body:   "{}",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := newResponse(tt.status, "application/json", tt.body, http.MethodPost, "https://api.example.com/x")
			require.Equal(t, tt.want, IsRetryableServerError(res))

			// The body must be rewound so downstream parsing still works.
			remaining, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.body, string(remaining))
		})
	}
}

func TestIsRetryableServerError_NilResponse(t *testing.T) {
	require.False(t, IsRetryableServerError(nil))
}

// A preserved 5xx response must still convert into a structured error, which is
// the contract the ZIA legacy client relies on after exhausting its retries.
func TestIsRetryableServerError_ResponseStillParsable(t *testing.T) {
	body := `{"code":"UNEXPECTED_ERROR","message":"An unexpected error has occurred"}`
	res := newResponse(http.StatusInternalServerError, "application/json", body, http.MethodPost, "https://api.example.com/urlFilteringRules")

	require.False(t, IsRetryableServerError(res))

	err := CheckErrorInResponse(res, nil)
	respErr, ok := AsErrorResponse(err)
	require.True(t, ok)
	require.NotNil(t, respErr.Parsed)
	require.Equal(t, "UNEXPECTED_ERROR", respErr.Parsed.Code)
	require.Equal(t, http.StatusInternalServerError, respErr.Parsed.Status)
}
