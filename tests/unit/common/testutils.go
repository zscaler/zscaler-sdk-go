// Package common provides test utilities for unit tests that actually exercise SDK code
package common

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa"
	zwaservices "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services"
)

// ===================================================================
// Test Service Factory - Creates real services with mocked HTTP
// ===================================================================

// MockTransport redirects all HTTP requests to a test server
type MockTransport struct {
	TestServerURL string
}

func (t *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Redirect the request to our mock server
	newURL := t.TestServerURL + req.URL.Path
	if req.URL.RawQuery != "" {
		newURL += "?" + req.URL.RawQuery
	}

	newReq, err := http.NewRequestWithContext(req.Context(), req.Method, newURL, req.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers
	for k, v := range req.Header {
		newReq.Header[k] = v
	}

	return http.DefaultTransport.RoundTrip(newReq)
}

// CreateTestService creates a zscaler.Service configured to use the mock server
// using ONLY the OneAPI client. The legacy per-cloud clients (zia.Client,
// zpa.Client, …) are intentionally never constructed or wired up here.
//
// OneAPI-only invariant
//
// All *requests.go helpers in zscaler/ (ziarequests.go, zparequests.go,
// zccrequests.go, zdxrequests.go, ztwrequests.go) branch on
// `c.oauth2Credentials.UseLegacyClient`. If that flag is true, the helper
// delegates to `LegacyClient.<X>Client`, which expects a fully-built
// per-cloud client (with its own Logger/Context/auth token). The unit
// test harness does NOT build those, so the legacy path would either
// return errLegacyClientNotSet or — worse — segfault if a stale
// LegacyClient.XClient is somehow attached.
//
// To make the harness bulletproof against ambient pollution
// (`ZSCALER_USE_LEGACY_CLIENT=true` in the shell, a `~/.zscaler/zscaler.yaml`
// loaded by `readConfigFromSystem`, an envconfig binding picked up by
// `readConfigFromEnvironment`, etc.), we explicitly clamp the config to
// the OneAPI path immediately after `NewConfiguration` returns and
// BEFORE `NewOneAPIClient` builds the service.
func CreateTestService(ctx context.Context, server *TestServer, customerID string) (*zscaler.Service, error) {
	// Create custom HTTP client with mock transport
	httpClient := &http.Client{
		Transport: &MockTransport{TestServerURL: server.URL},
		Timeout:   30 * time.Second,
	}

	// Create configuration with pre-populated auth token to skip OAuth
	// The key insight is that if AuthToken is valid (not expired),
	// authenticate() will skip actual authentication
	cfg, err := zscaler.NewConfiguration(
		zscaler.WithZPACustomerID(customerID),
		zscaler.WithVanityDomain("test"),
		zscaler.WithZscalerCloud(""),
		zscaler.WithHttpClientPtr(httpClient),
		zscaler.WithTestingDisableHttpsCheck(true),
		zscaler.WithCache(false), // Disable cache for testing
	)
	if err != nil {
		return nil, err
	}

	// OneAPI-only: legacy client is intentionally not configured.
	// `NewConfiguration` calls `readConfigFromEnvironment` (envconfig) which
	// honors `ZSCALER_USE_LEGACY_CLIENT`, plus `readConfigFromSystem`
	// (~/.zscaler/zscaler.yaml). If either path flips `UseLegacyClient` to
	// true, every subsequent SDK call would route through
	// `client.oauth2Credentials.LegacyClient.<X>Client` (see
	// zscaler/zparequests.go:21, zscaler/ziarequests.go:67, etc.), which
	// either errors with errLegacyClientNotSet or panics on a partially
	// initialized legacy client. We force the OneAPI path here so unit
	// tests are reproducible regardless of shell or filesystem state.
	cfg.UseLegacyClient = false
	cfg.LegacyClient = nil

	// Pre-populate a valid auth token to skip authentication
	cfg.Zscaler.Client.AuthToken = &zscaler.AuthToken{
		AccessToken: "mock-test-token-12345",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour), // Valid for 1 hour
	}

	// Update all HTTP clients to use our mock transport
	cfg.HTTPClient = httpClient
	cfg.ZPAHTTPClient = httpClient
	cfg.ZIAHTTPClient = httpClient
	cfg.ZTWHTTPClient = httpClient
	cfg.ZCCHTTPClient = httpClient
	cfg.ZDXHTTPClient = httpClient

	// Create the OneAPI client - this will skip auth since token is valid
	service, err := zscaler.NewOneAPIClient(cfg)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// MustCreateTestService creates a test service or panics
func MustCreateTestService(server *TestServer, customerID string) *zscaler.Service {
	service, err := CreateTestService(context.Background(), server, customerID)
	if err != nil {
		panic("failed to create test service: " + err.Error())
	}
	return service
}

// =============================================================================
// Per-Cloud Test Service Factories (OneAPI)
// =============================================================================
//
// These wrappers around CreateTestService cut the per-test boilerplate
// from ~5 lines down to one. Each registers t.Cleanup(server.Close) so
// callers don't need a defer.
//
// All factories return the same shape:  (server, service)
//
// Example:
//
//	server, service := common.NewZPATestService(t, "123456789")
//	server.On("GET", common.ZPAPath("123456789", "appConnectorGroup", "abc-1"),
//	    common.SuccessResponse(common.ZPAList(groups)))
//	got, _, err := appconnectorgroup.Get(ctx, service, "abc-1")
//
// The OneAPI client routes every cloud's traffic through the same mocked
// transport, so a single service instance can serve all clouds — but
// keeping per-cloud factories makes test intent explicit and lets us
// pre-set sensible cloud-specific defaults (e.g. customer ID for ZPA).

// NewZPATestService creates a ZPA-flavoured test service. Registers the
// supplied customerID with the OneAPI client so ZPA URL paths build
// correctly (ZPA endpoints embed the customer ID).
func NewZPATestService(t *testing.T, customerID string) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, customerID)
	if err != nil {
		t.Fatalf("NewZPATestService: %v", err)
	}
	return server, service
}

// NewZIATestService creates a test service for ZIA endpoints. ZIA does
// not embed a customer ID in URLs, so the customerID is left blank.
func NewZIATestService(t *testing.T) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, "")
	if err != nil {
		t.Fatalf("NewZIATestService: %v", err)
	}
	return server, service
}

// NewZCCTestService creates a test service for ZCC v1/v2 endpoints.
func NewZCCTestService(t *testing.T) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, "")
	if err != nil {
		t.Fatalf("NewZCCTestService: %v", err)
	}
	return server, service
}

// NewZDXTestService creates a test service for ZDX endpoints.
func NewZDXTestService(t *testing.T) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, "")
	if err != nil {
		t.Fatalf("NewZDXTestService: %v", err)
	}
	return server, service
}

// NewZTWTestService creates a test service for ZTW endpoints.
func NewZTWTestService(t *testing.T) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, "")
	if err != nil {
		t.Fatalf("NewZTWTestService: %v", err)
	}
	return server, service
}

// NewZIDTestService creates a test service for ZID endpoints.
func NewZIDTestService(t *testing.T) (*TestServer, *zscaler.Service) {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)
	service, err := CreateTestService(context.Background(), server, "")
	if err != nil {
		t.Fatalf("NewZIDTestService: %v", err)
	}
	return server, service
}

// ===================================================================
// ZWA Test Service Factory
// ===================================================================

// CreateZWATestService creates a ZWA service configured to use the mock server
func CreateZWATestService(ctx context.Context, server *TestServer) (*zwaservices.Service, error) {
	// Parse the test server URL
	baseURL, err := url.Parse(server.URL)
	if err != nil {
		return nil, err
	}

	// Create custom HTTP client with mock transport
	httpClient := &http.Client{
		Transport: &MockTransport{TestServerURL: server.URL},
		Timeout:   30 * time.Second,
	}

	// Create a mock configuration with pre-populated auth token and logger
	cfg := &zwa.Configuration{
		HTTPClient:    httpClient,
		BaseURL:       baseURL,
		UserAgent:     "zscaler-sdk-go-test",
		Context:       ctx,
		Logger:        logger.GetDefaultLogger("zwa-test: "),
		DefaultHeader: make(map[string]string),
	}
	cfg.ZWA.Client.AuthToken = &zwa.AuthToken{
		AccessToken: "mock-zwa-token-12345",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
	}
	cfg.ZWA.Client.ZWAAPIKeyID = "mock-key-id"
	cfg.ZWA.Client.ZWAAPISecret = "mock-secret"
	cfg.ZWA.Client.RequestTimeout = 30 * time.Second
	cfg.ZWA.Client.RateLimit.MaxRetries = 1
	cfg.ZWA.Client.RateLimit.RetryWaitMin = time.Second
	cfg.ZWA.Client.RateLimit.RetryWaitMax = time.Second * 5

	// Create the ZWA client directly without authentication
	client := &zwa.Client{
		Config: cfg,
	}

	return zwaservices.New(client), nil
}
