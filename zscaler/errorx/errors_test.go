package errorx

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

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
