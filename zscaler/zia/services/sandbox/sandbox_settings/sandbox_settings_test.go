package sandbox_settings

import (
	"context"
	"regexp"
	"testing"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/tests"
)

// isValidMD5 checks if the given string is a valid MD5 hash.
func isValidMD5(hash string) bool {
	matched, _ := regexp.MatchString(`^[a-fA-F0-9]{32}$`, hash)
	return matched
}

func TestUpdateBaAdvancedSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: []string{
			"42914d6d213a20a2684064be5c80ffa9",
			"c0202cf6aeab8437c638533d14563d35",
			"1ca31319721740ecb79f4b9ee74cd9b0",
			"2c373a7e86d0f3469849971e053bcc82",
			"40858748e03a544f6b562a687777397a",
		},
	}

	updatedSettings, err := Update(context.Background(), service, desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}

func TestValidateMD5Hashes(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Define the desired settings for the update
	hashes := []string{
		"42914d6d213a20a2684064be5c80ffa9",
		"c0202cf6aeab8437c638533d14563d35",
		"1ca31319721740ecb79f4b9ee74cd9b0",
		"2c373a7e86d0f3469849971e053bcc82",
		"40858748e03a544f6b562a687777397a",
		"invalidhash", // invalid
	}

	// Validate each hash
	validHashes := make([]string, 0, len(hashes))
	for _, hash := range hashes {
		if isValidMD5(hash) {
			validHashes = append(validHashes, hash)
		} else {
			t.Logf("Hash '%s' is not a valid MD5 hash", hash)
		}
	}

	// Proceed only if all hashes are valid
	if len(validHashes) == len(hashes) {
		desiredSettings := BaAdvancedSettings{
			FileHashesToBeBlocked: validHashes,
		}

		updatedSettings, err := Update(context.Background(), service, desiredSettings)
		if err != nil {
			t.Errorf("Error updating BA Advanced Settings: %v", err)
		}
		if updatedSettings == nil {
			t.Errorf("Expected updated BA Advanced Settings, got nil")
		}
	} else {
		t.Log("Update skipped due to invalid MD5 hash in the list")
	}
}

func TestGetBaAdvancedSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	settings, err := Get(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting BA Advanced Settings: %v", err)
	}
	if settings == nil {
		t.Error("Expected BA Advanced Settings, got nil")
	}
}

func TestGetFileHashCount(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	hashCount, err := GetFileHashCount(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting file hash count: %v", err)
	}
	if hashCount == nil {
		t.Error("Expected file hash count, got nil")
	}
}

func TestEmptyHashList(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	// Define the desired settings for the update
	desiredSettings := BaAdvancedSettings{
		FileHashesToBeBlocked: []string{},
	}

	updatedSettings, err := Update(context.Background(), service, desiredSettings)
	if err != nil {
		t.Errorf("Error updating BA Advanced Settings: %v", err)
	}
	if updatedSettings == nil {
		t.Error("Expected updated BA Advanced Settings, got nil")
	}
}
