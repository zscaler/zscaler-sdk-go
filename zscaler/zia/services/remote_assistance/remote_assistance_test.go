package remote_assistance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestRemoteAssistance(t *testing.T) {
	// Create the client
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Step 1: Retrieve the current Remote Assistance settings
	originalSettings, err := GetRemoteAssistance(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving Remote Assistance settings: %v", err)
	}
	t.Logf("Original Remote Assistance settings: %+v", originalSettings)

	// Step 2: Calculate dates 30 days in the future
	viewOnlyUntil := time.Now().Add(30 * 24 * time.Hour).UnixMilli()   // 30 days from now in milliseconds
	fullAccessUntil := time.Now().Add(60 * 24 * time.Hour).UnixMilli() // 60 days from now in milliseconds (example)

	// Step 3: Modify the settings
	updatedSettings := *originalSettings // Start with the original settings
	updatedSettings.ViewOnlyUntil = viewOnlyUntil
	updatedSettings.FullAccessUntil = fullAccessUntil
	updatedSettings.UsernameObfuscated = !originalSettings.UsernameObfuscated
	updatedSettings.DeviceInfoObfuscate = !originalSettings.DeviceInfoObfuscate

	// Step 4: Update the Remote Assistance settings
	_, _, err = UpdateRemoteAssistance(ctx, service, updatedSettings)
	if err != nil {
		t.Fatalf("Error updating Remote Assistance settings: %v", err)
	}
	t.Logf("Updated Remote Assistance settings sent successfully")

	// Step 5: Retrieve the settings again and verify changes
	newSettings, err := GetRemoteAssistance(ctx, service)
	if err != nil {
		t.Fatalf("Error retrieving updated Remote Assistance settings: %v", err)
	}
	t.Logf("Updated Remote Assistance settings retrieved: %+v", newSettings)

	// Step 6: Assert changes
	assert.Equal(t, updatedSettings.ViewOnlyUntil, newSettings.ViewOnlyUntil, "ViewOnlyUntil value mismatch")
	assert.Equal(t, updatedSettings.FullAccessUntil, newSettings.FullAccessUntil, "FullAccessUntil value mismatch")
	assert.Equal(t, updatedSettings.UsernameObfuscated, newSettings.UsernameObfuscated, "UsernameObfuscated value mismatch")
	assert.Equal(t, updatedSettings.DeviceInfoObfuscate, newSettings.DeviceInfoObfuscate, "DeviceInfoObfuscate value mismatch")

	// Step 7: Revert to original settings to maintain test idempotency
	_, _, err = UpdateRemoteAssistance(ctx, service, *originalSettings)
	if err != nil {
		t.Fatalf("Error reverting Remote Assistance settings: %v", err)
	}
	t.Logf("Reverted Remote Assistance settings to original values")
}
