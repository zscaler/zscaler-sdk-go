package mobile_threat_settings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestMobileThreatSettings(t *testing.T) {
	// Initialize the API client using the testing setup (not shown)
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	// Retrieve current Mobile Threat Settings
	currentSettings, err := GetMobileThreatSettings(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching initial Mobile Threat Settings: %v", err)
	}

	// Store original settings to restore later
	originalSettings := *currentSettings

	// Modify the settings for the test
	testSettings := *currentSettings
	testSettings.BlockAppsWithMaliciousActivity = !currentSettings.BlockAppsWithMaliciousActivity
	testSettings.BlockAppsWithKnownVulnerabilities = !currentSettings.BlockAppsWithKnownVulnerabilities
	testSettings.BlockAppsSendingUnencryptedUserCredentials = !currentSettings.BlockAppsSendingUnencryptedUserCredentials
	testSettings.BlockAppsSendingLocationInfo = !currentSettings.BlockAppsSendingLocationInfo
	testSettings.BlockAppsSendingPersonallyIdentifiableInfo = !currentSettings.BlockAppsSendingPersonallyIdentifiableInfo
	testSettings.BlockAppsSendingDeviceIdentifier = !currentSettings.BlockAppsSendingDeviceIdentifier
	testSettings.BlockAppsCommunicatingWithAdWebsites = !currentSettings.BlockAppsCommunicatingWithAdWebsites
	testSettings.BlockAppsCommunicatingWithRemoteUnknownServers = !currentSettings.BlockAppsCommunicatingWithRemoteUnknownServers
	// Update with the new settings
	updatedSettings, _, err := UpdateMobileThreatSettings(context.Background(), service, testSettings)
	if err != nil {
		t.Fatalf("Error updating Mobile Threat Settings: %v", err)
	}

	// Check if the settings were updated correctly
	assert.Equal(t, testSettings.BlockAppsWithMaliciousActivity, updatedSettings.BlockAppsWithMaliciousActivity, "BlockAppsWithMaliciousActivity should be updated")
	assert.Equal(t, testSettings.BlockAppsWithKnownVulnerabilities, updatedSettings.BlockAppsWithKnownVulnerabilities, "BlockAppsWithKnownVulnerabilities should be updated")
	assert.Equal(t, testSettings.BlockAppsSendingUnencryptedUserCredentials, updatedSettings.BlockAppsSendingUnencryptedUserCredentials, "BlockAppsSendingUnencryptedUserCredentials should be updated")
	assert.Equal(t, testSettings.BlockAppsSendingLocationInfo, updatedSettings.BlockAppsSendingLocationInfo, "BlockAppsSendingLocationInfo should be updated")
	assert.Equal(t, testSettings.BlockAppsSendingPersonallyIdentifiableInfo, updatedSettings.BlockAppsSendingPersonallyIdentifiableInfo, "BlockAppsSendingPersonallyIdentifiableInfo should be updated")
	assert.Equal(t, testSettings.BlockAppsSendingDeviceIdentifier, updatedSettings.BlockAppsSendingDeviceIdentifier, "BlockAppsSendingDeviceIdentifier should be updated")
	assert.Equal(t, testSettings.BlockAppsCommunicatingWithAdWebsites, updatedSettings.BlockAppsCommunicatingWithAdWebsites, "BlockAppsCommunicatingWithAdWebsites should be updated")
	assert.Equal(t, testSettings.BlockAppsCommunicatingWithRemoteUnknownServers, updatedSettings.BlockAppsCommunicatingWithRemoteUnknownServers, "BlockAppsCommunicatingWithRemoteUnknownServers should be updated")

	// Restore original settings
	_, _, err = UpdateMobileThreatSettings(context.Background(), service, originalSettings)
	if err != nil {
		t.Fatalf("Error restoring original Mobile Threat Settings: %v", err)
	}
}
