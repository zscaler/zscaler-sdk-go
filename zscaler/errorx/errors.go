package errorx

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ParsedAPIError struct {
	Code      interface{} `json:"code"`
	Message   string      `json:"message"`
	ID        string      `json:"id,omitempty"`
	Reason    string      `json:"reason,omitempty"`
	Exception string      `json:"exception,omitempty"`
	URL       string      `json:"url"`
	Status    int         `json:"status"`
}

type ErrorResponse struct {
	Response *http.Response
	Err      error
	Parsed   *ParsedAPIError
	Message  string
}

func (r *ErrorResponse) Error() string {
	if r.Parsed != nil {
		out, _ := json.MarshalIndent(r.Parsed, "", "  ")
		return fmt.Sprintf("Error: %s", string(out))
	}
	if r.Response != nil {
		return fmt.Sprintf("FAILED: %v %v -> %d %v\n%v",
			r.Response.Request.Method,
			r.Response.Request.URL,
			r.Response.StatusCode,
			r.Response.Status,
			r.Message,
		)
	}
	return fmt.Sprintf("FAILED: %v", r.Err)
}

func CheckErrorInResponse(res *http.Response, respErr error) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return respErr
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")
	isJSON := strings.Contains(contentType, "application/json")
	msg := strings.TrimSpace(string(bodyBytes))

	// ✅ Only fallback if it's non-JSON and the message matches known OneAPI error
	if !isJSON && strings.Contains(strings.ToLower(msg), "only through zscaler oneapi") {
		return NewOneAPIFallbackError(bodyBytes, res.Request.Method, res.Request.URL.Path, getBaseURL(res.Request.URL))
	}

	// Continue with normal JSON or generic error handling
	parsed := &ParsedAPIError{
		Status: res.StatusCode,
		URL:    res.Request.URL.String(),
	}

	if isJSON {
		var jsonBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
			if code, ok := jsonBody["code"]; ok {
				parsed.Code = code
			}
			if msg, ok := jsonBody["message"].(string); ok {
				parsed.Message = msg
			}
			if id, ok := jsonBody["id"].(string); ok {
				parsed.ID = id
			}
			if reason, ok := jsonBody["reason"].(string); ok {
				parsed.Reason = reason
			}
			if ex, ok := jsonBody["exception"].(string); ok {
				parsed.Exception = ex
			}
		} else {
			parsed.Message = fmt.Sprintf("Failed to parse JSON error body: %s", err.Error())
		}
	} else {
		parsed.Message = msg
	}

	return &ErrorResponse{
		Response: res,
		Err:      respErr,
		Parsed:   parsed,
		Message:  string(bodyBytes),
	}
}

func getBaseURL(u *url.URL) string {
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host)
}

func NewOneAPIFallbackError(respBody []byte, method, endpoint, baseURL string) *ErrorResponse {
	fullURL := fmt.Sprintf("%s%s", baseURL, endpoint)

	return &ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Status:     "401 Unauthorized",
			Request: &http.Request{
				Method: method,
				URL:    mustParseURL(fullURL), // ✅ use helper here
			},
		},
		Parsed: &ParsedAPIError{
			Status:  http.StatusUnauthorized,
			Message: strings.TrimSpace(string(respBody)),
			URL:     fullURL,
			Code:    "ONLY_ONEAPI_SUPPORTED",
		},
		Message: strings.TrimSpace(string(respBody)),
		Err:     fmt.Errorf("unexpected non-JSON error response"),
	}
}

func mustParseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		return &url.URL{Path: raw}
	}
	return u
}

func (r *ErrorResponse) IsObjectNotFound() bool {
	if r == nil || r.Response == nil {
		return false
	}
	if r.Response.StatusCode == http.StatusNotFound {
		return true
	}
	if r.Parsed != nil && r.Parsed.ID == "resource.not.found" {
		return true
	}
	return false
}

// IsSessionInvalidError checks if the response indicates a session invalidation error
// that requires token refresh. Only checks for known error messages returned by the API.
func IsSessionInvalidError(res *http.Response) bool {
	if res.StatusCode != http.StatusUnauthorized {
		return false
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	// Rewind the response body for potential reuse
	res.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	bodyStr := string(bodyBytes)

	// Only check for error messages we know for certain are returned by the API
	// The API may return different formats depending on the service/endpoint
	knownSessionInvalidMessages := []string{
		"SESSION_NOT_VALID",                         // Legacy/direct error code
		"getAttribute: Session already invalidated", // Java exception message format
		"Resource Access Blocked",                   // Occurs under high concurrency/load - API returns 401 instead of 429
	}

	for _, msg := range knownSessionInvalidMessages {
		if strings.Contains(bodyStr, msg) {
			return true
		}
	}

	return false
}

// IsEditLockError checks if the response indicates an edit lock conflict error
// that should be retried. This occurs when another operation is in progress.
func IsEditLockError(res *http.Response) bool {
	if res.StatusCode != http.StatusConflict {
		return false
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	defer res.Body.Close()

	// Rewind the response body for potential reuse
	res.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	bodyStr := string(bodyBytes)

	// Check for known edit lock error messages
	editLockErrorMessages := []string{
		"EDIT_LOCK_NOT_AVAILABLE",
		"Resource Access Blocked",
		"Failed during enter Org barrier",
	}

	for _, msg := range editLockErrorMessages {
		if strings.Contains(bodyStr, msg) {
			return true
		}
	}

	return false
}
