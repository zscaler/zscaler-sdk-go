package security_policy_settings

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/tests"
)

func randomURL() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(".example%d.com", rand.Intn(10000))
}

func TestSecurityPolicySettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &Service{Client: client}

	// Backup initial settings
	initialSettings, err := service.GetListUrls()
	if err != nil {
		t.Fatalf("Error fetching initial settings: %v", err)
	}

	// Generate test URLs for whitelisting and blacklisting
	testWhiteListURLs := []string{randomURL(), randomURL(), randomURL()}
	testBlackListURLs := []string{randomURL(), randomURL(), randomURL()}

	// Update white and black list URLs
	newSettings := ListUrls{
		White: testWhiteListURLs,
		Black: testBlackListURLs,
	}

	updatedSettings, err := service.UpdateListUrls(newSettings)
	if err != nil {
		t.Fatalf("Error updating settings: %v", err)
	}

	// Verify if the settings were updated
	if !reflect.DeepEqual(updatedSettings.White, testWhiteListURLs) {
		t.Errorf("Whitelist URLs were not updated correctly. Expected: %v, Got: %v", testWhiteListURLs, updatedSettings.White)
	}

	if !reflect.DeepEqual(updatedSettings.Black, testBlackListURLs) {
		t.Errorf("Blacklist URLs were not updated correctly. Expected: %v, Got: %v", testBlackListURLs, updatedSettings.Black)
	}

	// Restore initial settings
	_, err = service.UpdateListUrls(*initialSettings)
	if err != nil {
		t.Fatalf("Error restoring initial settings: %v", err)
	}

	// Verify if the settings were restored
	finalSettings, err := service.GetListUrls()
	if err != nil {
		t.Fatalf("Error fetching final settings: %v", err)
	}

	if !reflect.DeepEqual(finalSettings, initialSettings) {
		t.Errorf("Settings were not restored correctly. Expected: %+v, Got: %+v", initialSettings, finalSettings)
	}
}
