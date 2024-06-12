package platforms

import (
	"reflect"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

func TestGetAllPlatforms(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Failed to create ZPA client: %v", err)
	}
	service := services.New(client)

	// Test case: Normal scenario
	t.Run("TestGetAllPlatformsNormal", func(t *testing.T) {
		platforms, resp, err := GetAllPlatforms(service)
		if err != nil {
			t.Fatalf("Failed to fetch platforms: %v", err)
		}

		if resp.StatusCode >= 400 {
			t.Errorf("Received an HTTP error %d when fetching platforms", resp.StatusCode)
		}

		if platforms == nil {
			t.Fatal("Platforms nil, expected a valid response")
		}

		tests := map[string]string{
			"linux":   "Linux",
			"android": "Android",
			"windows": "Windows",
			"ios":     "iOS",
			"mac":     "Mac", // adjusted this line
		}

		platformValues := getValuesByTags(platforms)
		for jsonTag, expectedValue := range tests {
			actualValue, found := platformValues[jsonTag]
			if !found || actualValue != expectedValue {
				t.Errorf("Expected %s but got %s for json tag %s", expectedValue, actualValue, jsonTag)
			}
		}
	})

	// Test case: Error scenario
	t.Run("TestGetAllPlatformsError", func(t *testing.T) {
		// Temporarily change the client configuration to trigger an error
		service.Client.Config.CustomerID = "invalid_customer_id"
		_, _, err := GetAllPlatforms(service)
		if err == nil {
			t.Errorf("Expected error while fetching platforms with invalid customer ID, got nil")
		}
		// Reset the customer ID to avoid affecting other tests
		service.Client.Config.CustomerID = client.Config.CustomerID
	})
}

func getValuesByTags(types *Platforms) map[string]string {
	values := make(map[string]string)
	r := reflect.ValueOf(types).Elem()
	for i := 0; i < r.NumField(); i++ {
		fieldTag := r.Type().Field(i).Tag.Get("json")
		values[fieldTag] = r.Field(i).String()
	}
	return values
}
