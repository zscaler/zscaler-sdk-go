package errorx

import (
	"encoding/json"
	"errors"
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

// Unwrap exposes the underlying error so that *ErrorResponse participates in
// Go's error chain. This allows callers to use errors.Is / errors.As even when
// an *ErrorResponse has been wrapped (e.g. via fmt.Errorf("...: %w", err)).
func (r *ErrorResponse) Unwrap() error {
	if r == nil {
		return nil
	}
	return r.Err
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

// AsErrorResponse safely extracts an *ErrorResponse from an arbitrary error.
// It returns (nil, false) when the error is not (and does not wrap) an
// *ErrorResponse — for example a plain *errors.errorString produced by a
// cancelled or timed-out request. Callers should use this instead of an
// unguarded type assertion (err.(*ErrorResponse)), which panics when the
// dynamic type does not match.
func AsErrorResponse(err error) (*ErrorResponse, bool) {
	if err == nil {
		return nil, false
	}
	var respErr *ErrorResponse
	if errors.As(err, &respErr) {
		return respErr, true
	}
	return nil, false
}

// IsObjectNotFound safely reports whether err represents an object-not-found
// condition. It never panics regardless of the concrete error type, and it is
// aware of wrapped errors. This is the preferred way for callers to check for
// not-found errors.
func IsObjectNotFound(err error) bool {
	respErr, ok := AsErrorResponse(err)
	return ok && respErr.IsObjectNotFound()
}

// IsLimitExceeded safely reports whether err represents a tenant resource limit
// exceeded condition. It never panics regardless of the concrete error type,
// and it is aware of wrapped errors.
func IsLimitExceeded(err error) bool {
	respErr, ok := AsErrorResponse(err)
	return ok && respErr.IsLimitExceeded()
}

// IsLimitExceeded checks if the response indicates a tenant resource limit has been
// reached (e.g., "Maximum 100 static IPs are allowed. Limit has exceeded.").
// The API returns HTTP 403 with code "LIMIT_EXCEEDED" when the tenant's maximum
// capacity for a given resource type is exhausted. This error is non-retryable.
func (r *ErrorResponse) IsLimitExceeded() bool {
	if r == nil || r.Response == nil {
		return false
	}
	if r.Response.StatusCode == http.StatusForbidden && r.Parsed != nil {
		code, ok := r.Parsed.Code.(string)
		if ok && code == "LIMIT_EXCEEDED" {
			return true
		}
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
		"SESSION_NOT_VALID",                              // Legacy/direct error code
		"getAttribute: Session already invalidated",      // Java exception message format (legacy)
		"getAttributeNames: Session already invalidated", // Java exception message format (newer API variant)
		"Session already invalidated",                    // Broad match for any future method name variants
		"Resource Access Blocked",                        // Occurs under high concurrency/load - API returns 401 instead of 429
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
// The API may return either 409 Conflict or 412 Precondition Failed for these
// transient "org barrier" / edit lock conditions.
func IsEditLockError(res *http.Response) bool {
	if res.StatusCode != http.StatusConflict && res.StatusCode != http.StatusPreconditionFailed {
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

// transientServerErrorMessages are API codes/messages that accompany a 5xx and
// indicate a genuinely transient server-side condition, which a retry can
// plausibly resolve. These are distinct from deterministic failures that reuse
// the same 5xx status.
var transientServerErrorMessages = []string{
	"EDIT_LOCK_NOT_AVAILABLE",
	"Resource Access Blocked",
	"Failed during enter Org barrier",
	"Request processing failed, possibly because an expected precondition was not met",
}

// alwaysRetryableServerStatuses are 5xx codes that never carry an application
// verdict. They are produced by gateways, load balancers and overload
// protection rather than by the API's business logic, so a body that happens to
// parse as a JSON error must not be mistaken for a deterministic failure.
//
// 503 matters most: every client's Backoff closure treats it exactly like a 429
// and honours Retry-After, so it has to stay retryable for that path to be
// reachable at all.
var alwaysRetryableServerStatuses = map[int]bool{
	http.StatusBadGateway:         true, // 502
	http.StatusServiceUnavailable: true, // 503
	http.StatusGatewayTimeout:     true, // 504
}

// IsRetryableServerError reports whether a 5xx response is worth retrying.
//
// The ZIA API reuses HTTP 500 with code UNEXPECTED_ERROR for permanent request
// validation failures — for example submitting a URL filtering rule that names a
// category which does not exist. Retrying those is futile, and because the
// legacy retry budget is large it turns a single bad request into a long stall
// that ends with the API's own explanation discarded (see issue #449).
//
// The decision is deliberately conservative: a retry is refused only when the
// body is a well-formed API error carrying a code, which is a deterministic
// verdict from the application itself. An empty body, an HTML error page from a
// load balancer, or any payload that does not parse is still treated as a
// transient infrastructure fault and retried, preserving the previous behaviour
// for real outages. A recognised transient marker always wins over both rules,
// and the body is not inspected at all for the infrastructure statuses in
// alwaysRetryableServerStatuses.
func IsRetryableServerError(res *http.Response) bool {
	if res == nil {
		return false
	}
	// 501 Not Implemented is a permanent condition, matching the upstream
	// retryablehttp policy.
	if res.StatusCode < 500 || res.StatusCode == http.StatusNotImplemented {
		return false
	}
	// 502 / 503 / 504 are infrastructure signals, never API verdicts: retry them
	// unconditionally, exactly as the upstream policy did before this function
	// existed.
	if alwaysRetryableServerStatuses[res.StatusCode] {
		return true
	}
	// A 5xx with no body at all carries no verdict, so it stays retryable. This
	// only happens for hand-built responses; the transport always sets a body.
	if res.Body == nil {
		return true
	}

	bodyBytes, _ := io.ReadAll(res.Body)
	res.Body.Close()

	// Rewind the response body so downstream error parsing still sees it.
	res.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	bodyStr := string(bodyBytes)

	for _, msg := range transientServerErrorMessages {
		if strings.Contains(bodyStr, msg) {
			return true
		}
	}

	var parsed struct {
		Code interface{} `json:"code"`
	}
	if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
		if code, ok := parsed.Code.(string); ok && code != "" {
			return false
		}
	}

	return true
}
