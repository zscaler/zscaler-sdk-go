// Package common provides test utilities for unit tests that actually exercise SDK code
package common

import (
	"context"
	"net/http"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
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
