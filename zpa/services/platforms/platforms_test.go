package platforms

import (
	"reflect"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
)

func TestGetAllPlatforms(t *testing.T) {
	client, err := tests.NewZpaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	service := New(client)

	platforms, resp, err := service.GetAllPlatforms()
	if err != nil {
		t.Fatalf("Error getting all platforms: %v", err)
	}

	if resp.StatusCode >= 400 {
		t.Errorf("Received an HTTP error: %d", resp.StatusCode)
	}

	if platforms == nil {
		t.Error("Expected platforms to be non-nil, even for empty response")
		return
	}

	tests := map[string]string{
		"linux":   "Linux platform",
		"android": "Android platform",
		"windows": "Windows platform",
		"ios":     "IOS platform",
		"mac":     "Mac platform", // adjusted this line
	}

	for tag, platform := range tests {
		if value := getValueByTag(platforms, tag); value == "" {
			t.Errorf("%s not found", platform)
		}
	}
}

func getValueByTag(platforms *Platforms, tag string) string {
	r := reflect.ValueOf(platforms).Elem()
	for i := 0; i < r.NumField(); i++ {
		fieldTag := r.Type().Field(i).Tag.Get("json")
		if fieldTag == tag {
			return r.Field(i).String()
		}
	}
	return ""
}
