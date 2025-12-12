// Package common provides test utilities for unit tests that actually exercise SDK code
package common

import (
	"context"
	"net/http"
	"net/url"
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
// This allows tests to call actual SDK functions and generate real coverage
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
