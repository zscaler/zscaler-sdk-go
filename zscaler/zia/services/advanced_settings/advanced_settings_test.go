package advanced_settings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestAdvancedSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Step 1: Retrieve current Advanced Settings
	currentSettings, err := GetAdvancedSettings(context.Background(), service)
	if err != nil {
		t.Fatalf("Error retrieving Advanced Settings: %v", err)
	}

	// Store original settings to restore later
	originalSettings := *currentSettings

	// Modify settings for the test
	testSettings := *currentSettings
	testSettings.EnableOffice365 = !currentSettings.EnableOffice365 // Toggle the current setting for testing

	// Step 2: Update Advanced Settings
	updatedSettings, _, err := UpdateAdvancedSettings(context.Background(), service, &testSettings)
	if err != nil {
		t.Fatalf("Error updating Advanced Settings: %v", err)
	}

	// Check if the settings were updated correctly
	assert.Equal(t, testSettings.EnableOffice365, updatedSettings.EnableOffice365, "EnableOffice365 setting should be updated")

	// Log for debugging
	t.Logf("Updated Advanced Settings: %+v", updatedSettings)

	// Optional: Restore original settings to clean up after test
	_, _, err = UpdateAdvancedSettings(context.Background(), service, &originalSettings)
	if err != nil {
		t.Logf("Error restoring original Advanced Settings: %v", err)
	}
}
