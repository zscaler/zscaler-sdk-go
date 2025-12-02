package sandbox_settings

import (
	"context"
	"regexp"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

// isValidMD5 checks if the given string is a valid MD5 hash.
func isValidMD5(hash string) bool {
	matched, _ := regexp.MatchString(`^[a-fA-F0-9]{32}$`, hash)
	return matched
}

func TestUpdateBaAdvancedSettings(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "sandbox_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Define the desired settings for the update using new V2 structure
	desiredSettings := BaAdvancedSettingsV2{
		MD5HashValueList: []MD5HashValue{
			{
				URL:        "42914d6d213a20a2684064be5c80ffa9",
				URLComment: "Test hash 1",
				Type:       MD5HashTypeCustomFilehashDeny,
			},
			{
				URL:        "c0202cf6aeab8437c638533d14563d35",
				URLComment: "Test hash 2",
				Type:       MD5HashTypeCustomFilehashDeny,
			},
			{
				URL:        "1ca31319721740ecb79f4b9ee74cd9b0",
				URLComment: "Test hash 3",
				Type:       MD5HashTypeCustomFilehashDeny,
			},
			{
				URL:        "2c373a7e86d0f3469849971e053bcc82",
				URLComment: "Test hash 4",
				Type:       MD5HashTypeCustomFilehashDeny,
			},
			{
				URL:        "40858748e03a544f6b562a687777397a",
				URLComment: "Test hash 5",
				Type:       MD5HashTypeCustomFilehashDeny,
			},
		},
	}

	updatedSettings, err := UpdateV2(context.Background(), service, desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings V2: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings V2, got nil")
	}
}

func TestValidateMD5Hashes(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "sandbox_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Define the desired settings for the update using new V2 structure
	hashEntries := []MD5HashValue{
		{URL: "42914d6d213a20a2684064be5c80ffa9", URLComment: "Valid hash 1", Type: MD5HashTypeCustomFilehashDeny},
		{URL: "c0202cf6aeab8437c638533d14563d35", URLComment: "Valid hash 2", Type: MD5HashTypeCustomFilehashDeny},
		{URL: "1ca31319721740ecb79f4b9ee74cd9b0", URLComment: "Valid hash 3", Type: MD5HashTypeCustomFilehashDeny},
		{URL: "2c373a7e86d0f3469849971e053bcc82", URLComment: "Valid hash 4", Type: MD5HashTypeCustomFilehashDeny},
		{URL: "40858748e03a544f6b562a687777397a", URLComment: "Valid hash 5", Type: MD5HashTypeCustomFilehashDeny},
		{URL: "invalidhash", URLComment: "Invalid hash", Type: MD5HashTypeCustomFilehashDeny}, // invalid
	}

	// Validate each hash
	validHashEntries := make([]MD5HashValue, 0, len(hashEntries))
	for _, entry := range hashEntries {
		if isValidMD5(entry.URL) {
			validHashEntries = append(validHashEntries, entry)
		} else {
			t.Logf("Hash '%s' is not a valid MD5 hash", entry.URL)
		}
	}

	// Proceed only if all hashes are valid
	if len(validHashEntries) == len(hashEntries) {
		desiredSettings := BaAdvancedSettingsV2{
			MD5HashValueList: validHashEntries,
		}

		updatedSettings, err := UpdateV2(context.Background(), service, desiredSettings)
		if err != nil {
			t.Errorf("Error updating BA Advanced Settings V2: %v", err)
		}
		if updatedSettings == nil {
			t.Errorf("Expected updated BA Advanced Settings V2, got nil")
		}
	} else {
		t.Log("Update skipped due to invalid MD5 hash in the list")
	}
}

func TestGetBaAdvancedSettings(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "sandbox_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	settings, err := GetV2(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting BA Advanced Settings V2: %v", err)
	}
	if settings == nil {
		t.Error("Expected BA Advanced Settings V2, got nil")
	}

	// Log the retrieved hash entries for verification
	if settings != nil && len(settings.MD5HashValueList) > 0 {
		for _, entry := range settings.MD5HashValueList {
			t.Logf("Hash: %s, Comment: %s, Type: %s", entry.URL, entry.URLComment, entry.Type)
		}
	}
}

func TestGetFileHashCount(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "sandbox_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	hashCount, err := GetFileHashCount(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting file hash count: %v", err)
	}
	if hashCount == nil {
		t.Error("Expected file hash count, got nil")
	}
}

// TestEmptyHashList tests clearing all MD5 hashes by sending an empty list.
// The API expects {"md5HashValueList": []} to clear all hashes (not an empty object {}).
func TestEmptyHashList(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "sandbox_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Define empty list to clear all hashes
	// This sends {"md5HashValueList": []} which is required by the API
	desiredSettings := BaAdvancedSettingsV2{
		MD5HashValueList: []MD5HashValue{},
	}

	updatedSettings, err := UpdateV2(context.Background(), service, desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings V2 with empty list: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings V2, got nil")
	}
}
