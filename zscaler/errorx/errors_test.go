package errorx

import (
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
