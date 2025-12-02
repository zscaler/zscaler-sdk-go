package ftp_control_policy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestFTPControlPolicy(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "ftp_control_policy", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	ctx := context.Background()

	// Step 1: Fetch current settings
	currentSettings, err := GetFTPControlPolicy(ctx, service)
	if err != nil {
		t.Fatalf("Error fetching initial FTP Control Policy settings: %v", err)
	}

	// Step 2: Backup original settings
	originalSettings := *currentSettings

	// Step 3: Prepare test settings
	testSettings := *currentSettings
	testSettings.FtpOverHttpEnabled = true
	testSettings.FtpEnabled = true
	testSettings.UrlCategories = []string{"ADULT_THEMES", "ADULT_SEX_EDUCATION"}
	testSettings.Urls = []string{"zscaler.com", "zscaler.net"}

	// Step 4: Update the settings
	updatedSettings, _, err := UpdateFTPControlPolicy(ctx, service, &testSettings)
	if err != nil {
		t.Fatalf("Error updating FTP Control Policy settings: %v", err)
	}

	// Step 5: Verify the update
	assert.Equal(t, testSettings.FtpOverHttpEnabled, updatedSettings.FtpOverHttpEnabled, "FtpOverHttpEnabled should be updated")
	assert.Equal(t, testSettings.FtpEnabled, updatedSettings.FtpEnabled, "FtpEnabled should be updated")
	assert.ElementsMatch(t, testSettings.UrlCategories, updatedSettings.UrlCategories, "UrlCategories should be updated")
	assert.ElementsMatch(t, testSettings.Urls, updatedSettings.Urls, "Urls should be updated")

	// Step 6: Restore original settings
	_, _, err = UpdateFTPControlPolicy(ctx, service, &originalSettings)
	if err != nil {
		t.Fatalf("Error restoring original FTP Control Policy settings: %v", err)
	}
}
