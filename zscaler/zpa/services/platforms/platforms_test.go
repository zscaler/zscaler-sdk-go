package platforms

import (
	"context"
	"reflect"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestGetAllPlatforms(t *testing.T) {
	client, err := tests.NewVCRTestClient(t, "platforms", "zpa")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// service, err := tests.NewZPAClient()
	// if err != nil {
	// 	t.Fatalf("Error creating client: %v", err)
	// }

	// Test case: Normal scenario
	t.Run("TestGetAllPlatformsNormal", func(t *testing.T) {
		platforms, resp, err := GetAllPlatforms(context.Background(), service)
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
