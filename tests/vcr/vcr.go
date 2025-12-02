// Package vcr provides VCR (Video Cassette Recorder) testing utilities
// for recording and replaying HTTP interactions in tests.
package vcr

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// RecordMode determines how VCR handles requests
type RecordMode int

const (
	// ModeRecording records new interactions (MOCK_TESTS=false)
	ModeRecording RecordMode = iota
	// ModePlayback replays recorded interactions (MOCK_TESTS=true)
	ModePlayback
)

// GetRecordMode returns the current mode based on MOCK_TESTS env var
func GetRecordMode() RecordMode {
	if os.Getenv("MOCK_TESTS") == "true" {
		return ModePlayback
	}
	return ModeRecording
}

// IsPlaybackMode returns true if we're in VCR playback mode
func IsPlaybackMode() bool {
	return GetRecordMode() == ModePlayback
}

// VCRConfig holds VCR configuration
type VCRConfig struct {
	CassetteName string     // Name of cassette file (without path)
	Service      string     // Service name: zia, zpa, zcc, zdx, zidentity, ztw
	Mode         RecordMode // Recording or playback mode
	TestDir      string     // Optional: override test directory detection
}

// getProjectRoot finds the project root directory
func getProjectRoot() string {
	// Try to find go.mod file by walking up directories
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}

	dir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "."
}

// NewVCRRecorder creates a new VCR recorder with sanitization
func NewVCRRecorder(t *testing.T, config VCRConfig) (*recorder.Recorder, error) {
	// Build cassette path
	projectRoot := getProjectRoot()
	cassettePath := filepath.Join(projectRoot, "tests", "cassettes", config.Service, config.CassetteName)

	// Ensure directory exists
	cassetteDir := filepath.Dir(cassettePath)
	if err := os.MkdirAll(cassetteDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cassette directory: %w", err)
	}

	// Determine recorder mode
	var mode recorder.Mode
	if config.Mode == ModePlayback {
		mode = recorder.ModeReplayOnly
	} else {
		// Use ModeRecordOnly to allow recording new interactions
		// This will pass unmatched requests through to the real API and record them
		mode = recorder.ModeRecordOnly
	}

	// Create recorder with options
	opts := []recorder.Option{
		recorder.WithMode(mode),
		recorder.WithSkipRequestLatency(true),
		recorder.WithHook(sanitizeInteraction, recorder.BeforeSaveHook),
		recorder.WithMatcher(customMatcher),
	}

	// In recording mode, set a real transport so unmatched requests go through
	if config.Mode == ModeRecording {
		opts = append(opts, recorder.WithRealTransport(http.DefaultTransport))
	}

	r, err := recorder.New(cassettePath, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create recorder: %w", err)
	}

	t.Logf("VCR: Using cassette %s (mode: %s)", cassettePath, modeString(config.Mode))

	return r, nil
}

func modeString(mode RecordMode) string {
	if mode == ModePlayback {
		return "playback"
	}
	return "recording"
}

// GetVCRTransport returns the http.RoundTripper from the recorder
func GetVCRTransport(r *recorder.Recorder) http.RoundTripper {
	return r
}

// sanitizeInteraction removes sensitive data from recorded interactions
func sanitizeInteraction(i *cassette.Interaction) error {
	// Sanitize request
	sanitizeRequest(i)

	// Sanitize response
	sanitizeResponse(i)

	return nil
}

// Regex patterns for sanitization
var (
	// Token patterns
	bearerTokenRegex  = regexp.MustCompile(`Bearer\s+[A-Za-z0-9\-_\.]+`)
	accessTokenRegex  = regexp.MustCompile(`"access_token"\s*:\s*"[^"]*"`)
	tokenRegex        = regexp.MustCompile(`"token"\s*:\s*"[^"]*"`)
	jwtRegex          = regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`)

	// Credential patterns (URL-encoded form data)
	clientIDRegex     = regexp.MustCompile(`client_id=[^&\s]*`)
	clientSecretRegex = regexp.MustCompile(`client_secret=[^&\s]*`)

	// JSON field patterns
	emailRegex        = regexp.MustCompile(`"[^"]*@[^"]*"`)
	privateKeyRegex   = regexp.MustCompile(`-----BEGIN[A-Z ]*PRIVATE KEY-----[\s\S]*?-----END[A-Z ]*PRIVATE KEY-----`)
	preSharedKeyRegex = regexp.MustCompile(`"preSharedKey"\s*:\s*"[^"]*"`)
	pskRegex          = regexp.MustCompile(`"psk"\s*:\s*"[^"]*"`)
	passwordRegex     = regexp.MustCompile(`"password"\s*:\s*"[^"]*"`)
	apiKeyRegex       = regexp.MustCompile(`"apiKey"\s*:\s*"[^"]*"`)

	// URL patterns for customer IDs (in paths like /customers/216196257331281920/)
	customerIDInURLRegex = regexp.MustCompile(`/customers/\d+`)
)

func sanitizeRequest(i *cassette.Interaction) {
	// Sanitize Authorization header
	for key, values := range i.Request.Headers {
		if strings.ToLower(key) == "authorization" {
			for idx := range values {
				i.Request.Headers[key][idx] = "Bearer REDACTED_TOKEN"
			}
		}
	}

	// Sanitize form values (parsed from request body)
	sensitiveFormFields := []string{
		"client_id", "client_secret", "password", "api_key", "apiKey",
		"secret", "token", "vanity_domain", "customer_id",
	}
	for _, field := range sensitiveFormFields {
		if _, exists := i.Request.Form[field]; exists {
			i.Request.Form[field] = []string{"REDACTED"}
		}
	}

	// Sanitize request body (URL-encoded form data)
	body := i.Request.Body
	body = clientIDRegex.ReplaceAllString(body, "client_id=REDACTED")
	body = clientSecretRegex.ReplaceAllString(body, "client_secret=REDACTED")
	body = passwordRegex.ReplaceAllString(body, `"password":"REDACTED"`)
	body = preSharedKeyRegex.ReplaceAllString(body, `"preSharedKey":"REDACTED"`)
	body = pskRegex.ReplaceAllString(body, `"psk":"REDACTED"`)
	body = apiKeyRegex.ReplaceAllString(body, `"apiKey":"REDACTED"`)
	body = privateKeyRegex.ReplaceAllString(body, "-----BEGIN PRIVATE KEY-----REDACTED-----END PRIVATE KEY-----")
	i.Request.Body = body

	// Sanitize Host header (contains vanity domain)
	i.Request.Host = "REDACTED.zslogin.net"

	// Normalize URL (replace real domains with test domains)
	// Note: We keep customer_id in URLs for VCR matching to work
	i.Request.URL = normalizeURL(i.Request.URL)
}

func sanitizeResponse(i *cassette.Interaction) {
	body := i.Response.Body

	// Redact tokens
	body = accessTokenRegex.ReplaceAllString(body, `"access_token":"REDACTED_TOKEN"`)
	body = tokenRegex.ReplaceAllString(body, `"token":"REDACTED_TOKEN"`)
	body = bearerTokenRegex.ReplaceAllString(body, "Bearer REDACTED_TOKEN")
	body = jwtRegex.ReplaceAllString(body, "REDACTED_JWT_TOKEN")

	// Redact emails
	body = emailRegex.ReplaceAllString(body, `"REDACTED"`)

	// Redact sensitive fields
	body = passwordRegex.ReplaceAllString(body, `"password":"REDACTED"`)
	body = preSharedKeyRegex.ReplaceAllString(body, `"preSharedKey":"REDACTED"`)
	body = pskRegex.ReplaceAllString(body, `"psk":"REDACTED"`)
	body = apiKeyRegex.ReplaceAllString(body, `"apiKey":"REDACTED"`)
	body = privateKeyRegex.ReplaceAllString(body, "-----BEGIN PRIVATE KEY-----REDACTED-----END PRIVATE KEY-----")

	i.Response.Body = body

	// Sanitize Content-Security-Policy header (contains vanity domain)
	if csp, exists := i.Response.Headers["Content-Security-Policy"]; exists {
		for idx, val := range csp {
			// Redact vanity domains in CSP header
			val = regexp.MustCompile(`https://[a-z0-9-]+\.zslogin\.net`).ReplaceAllString(val, "https://REDACTED.zslogin.net")
			val = regexp.MustCompile(`https://[a-z0-9-]+-admin\.zslogin\.net`).ReplaceAllString(val, "https://REDACTED-admin.zslogin.net")
			i.Response.Headers["Content-Security-Policy"][idx] = val
		}
	}

	// Remove sensitive headers
	delete(i.Response.Headers, "Set-Cookie")
	delete(i.Response.Headers, "set-cookie")
	delete(i.Response.Headers, "X-Zscloud-Customer-Id")
	delete(i.Response.Headers, "x-zscloud-customer-id")
}

// URL normalization patterns - redact vanity domains and normalize URLs
var urlPatterns = map[*regexp.Regexp]string{
	regexp.MustCompile(`https://[a-z0-9-]+\.zsapi\.net`):                    "https://api.zsapi.net",
	regexp.MustCompile(`https://[a-z0-9-]+\.zscaler[a-z]*\.net`):            "https://api.zscaler.net",
	regexp.MustCompile(`https://[a-z0-9-]+\.zslogin\.net`):                  "https://REDACTED.zslogin.net",
	regexp.MustCompile(`https://[a-z0-9-]+-admin\.zslogin\.net`):            "https://REDACTED-admin.zslogin.net",
	regexp.MustCompile(`connect-src\s+https://[a-z0-9-]+\.zslogin\.net`):    "connect-src https://REDACTED.zslogin.net",
	regexp.MustCompile(`frame-src\s+https://[a-z0-9-]+\.zslogin\.net`):      "frame-src https://REDACTED.zslogin.net",
	regexp.MustCompile(`frame-src\s+https://[a-z0-9-]+-admin\.zslogin\.net`):"frame-src https://REDACTED-admin.zslogin.net",
}

func normalizeURL(rawURL string) string {
	for pattern, replacement := range urlPatterns {
		rawURL = pattern.ReplaceAllString(rawURL, replacement)
	}
	return rawURL
}

// customMatcher matches requests by method, path, and query params (ignoring host)
func customMatcher(r *http.Request, i cassette.Request) bool {
	// Match HTTP method
	if r.Method != i.Method {
		return false
	}

	// Parse URLs
	rURL := r.URL
	iURL, err := url.Parse(i.URL)
	if err != nil {
		return false
	}

	// Match path (ignore host differences)
	if rURL.Path != iURL.Path {
		return false
	}

	// Match query parameters (order-independent)
	rQuery := rURL.Query()
	iQuery := iURL.Query()

	// Check all query params from request exist in recorded
	for key, values := range rQuery {
		iValues, ok := iQuery[key]
		if !ok {
			return false
		}
		if !stringSliceEqual(values, iValues) {
			return false
		}
	}

	return true
}

func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ============================================================================
// Deterministic Name Generator
// ============================================================================

var (
	nameCounter int
	nameMutex   sync.Mutex
)

// ResetNameCounter resets the deterministic name counter
// Call this at the beginning of each test
func ResetNameCounter() {
	nameMutex.Lock()
	defer nameMutex.Unlock()
	nameCounter = 0
}

// GenerateName generates a deterministic name for VCR testing
func GenerateName(prefix string) string {
	nameMutex.Lock()
	defer nameMutex.Unlock()
	nameCounter++
	return fmt.Sprintf("%s-vcr%04d", prefix, nameCounter)
}

// GenerateNameWithLength generates a deterministic name with max length
func GenerateNameWithLength(prefix string, maxLen int) string {
	name := GenerateName(prefix)
	if len(name) > maxLen {
		return name[:maxLen]
	}
	return name
}

