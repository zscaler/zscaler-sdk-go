package errorx

import (
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response
	Err      error
	Message  string
}

func (r *ErrorResponse) Error() string {
	if r.Response != nil {
		return fmt.Sprintf("FAILED: %v, %v, %d, %v, %v, %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Response.Status, r.Message, r.Err)
	}
	return fmt.Sprintf("FAILED: %v", r.Err)
}

func CheckErrorInResponse(res *http.Response, respErr error) error {
	if c := res.StatusCode; c >= 200 && c <= 299 {
		return respErr
	}
	errorResponse := &ErrorResponse{Response: res, Err: respErr}
	errorMessage, err := io.ReadAll(res.Body)
	if err == nil && len(errorMessage) > 0 {
		errorResponse.Message = string(errorMessage)
	}
	return errorResponse
}

// IsObjectNotFound returns true on missing object error (404).
func (r ErrorResponse) IsObjectNotFound() bool {
	return r.Response.StatusCode == 404
}
