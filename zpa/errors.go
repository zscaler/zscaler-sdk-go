package zpa

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response
	Message  string
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("FAILED: %v, %v, %d, %v, %v", r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Response.Status, r.Message)
}

func checkErrorInResponse(res *http.Response) error {
	if c := res.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: res}
	errorMessage, err := io.ReadAll(res.Body)
	if err == nil && len(errorMessage) > 0 {
		errorResponse.Message = string(errorMessage)
	}
	return errorResponse
}

type apiErrorResponse struct {
	ID string `json:"id"`
}

// isResourceNotFoundError returns true on missing object error (400).
func (r ErrorResponse) isResourceNotFoundError() bool {
	resp := apiErrorResponse{}
	err := json.Unmarshal([]byte(r.Message), &resp)
	if err != nil {
		return false
	}
	return resp.ID == "resource.not.found"
}

// IsObjectNotFound returns true on missing object error (404 & 400 with response  "id": "resource.not.found",).
func (r ErrorResponse) IsObjectNotFound() bool {

	return r.Response.StatusCode == 404 || r.Response.StatusCode == 400 && r.isResourceNotFoundError()
}
