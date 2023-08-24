package platforms

import (
	"testing"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func TestGetAllPlatforms(t *testing.T) {
	// Setup the client
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	// Basic check for all platforms
	platforms, _, err := service.GetAllPlatforms()
	if err != nil {
		t.Fatalf("Error getting all platforms: %v", err)
	}

	// Basic tests on the received platforms data
	if platforms.Linux == "" {
		t.Errorf("Linux platform not found")
	}
	if platforms.Android == "" {
		t.Errorf("Android platform not found")
	}
	if platforms.Windows == "" {
		t.Errorf("Windows platform not found")
	}
	if platforms.IOS == "" {
		t.Errorf("IOS platform not found")
	}
	if platforms.MacOS == "" {
		t.Errorf("Mac platform not found")
	}

	// 1. Test with specific customer ID
	client.Config.CustomerID = "216196257331281920"
	_, _, err = service.GetAllPlatforms()
	if err != nil {
		t.Errorf("Error fetching platforms for specific customer ID: %v", err)
	}

	// 2. Ensure response has no HTTP errors
	_, resp, err := service.GetAllPlatforms()
	if err != nil {
		t.Fatalf("Error fetching platforms: %v", err)
	}
	if resp.StatusCode >= 400 {
		t.Errorf("Received an HTTP error: %d", resp.StatusCode)
	}

	// 3. Test empty platform response - this assumes that the API might sometimes return a valid response with empty data
	platforms, _, err = service.GetAllPlatforms()
	if err != nil {
		t.Fatalf("Error fetching platforms: %v", err)
	}
	if platforms == nil {
		t.Error("Expected platforms to be non-nil, even for empty response")
	}
}
