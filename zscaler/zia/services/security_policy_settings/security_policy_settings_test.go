package security_policy_settings

/*
import (
	"context"
	"log"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

const (
	maxConflictRetries    = 10
	conflictRetryInterval = 3 * time.Second
)

func randomURL() string {
	// Use deterministic name for VCR testing
	// URLs must start with a dot for the ZIA API whitelist/blacklist
	return "." + tests.GetTestName("tests-site") + ".example.com"
}

func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		errStr := lastErr.Error()
		// Retry on edit lock or operation in progress errors
		if strings.Contains(errStr, `"code":"EDIT_LOCK_NOT_AVAILABLE"`) ||
			strings.Contains(errStr, `"code":"INVALID_OPERATION"`) ||
			strings.Contains(errStr, "Another custom url operation is in progress") ||
			strings.Contains(errStr, "operation is in progress") {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestSecurityPolicySettings(t *testing.T) {
	tests.ResetTestNameCounter()
	client, err := tests.NewVCRTestClient(t, "security_policy_settings", "zia")
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}
	defer client.Stop()
	service := client.Service

	// Backup initial settings
	initialSettings, err := GetListUrls(context.Background(), service)
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

	err = retryOnConflict(func() error {
		_, err := UpdateListUrls(context.Background(), service, newSettings)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating settings: %v", err)
	}

	updatedSettings, err := GetListUrls(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching updated settings: %v", err)
	}

	// Verify if the settings were updated
	if !areSlicesEqual(updatedSettings.White, testWhiteListURLs) {
		t.Errorf("Whitelist URLs were not updated correctly. Expected: %v, Got: %v", testWhiteListURLs, updatedSettings.White)
	}

	if !areSlicesEqual(updatedSettings.Black, testBlackListURLs) {
		t.Errorf("Blacklist URLs were not updated correctly. Expected: %v, Got: %v", testBlackListURLs, updatedSettings.Black)
	}

	// Restore initial settings
	err = retryOnConflict(func() error {
		_, err := UpdateListUrls(context.Background(), service, *initialSettings)
		return err
	})
	if err != nil {
		t.Fatalf("Error restoring initial settings: %v", err)
	}

	// Verify if the settings were restored
	finalSettings, err := GetListUrls(context.Background(), service)
	if err != nil {
		t.Fatalf("Error fetching final settings: %v", err)
	}

	if !areSlicesEqual(finalSettings.White, initialSettings.White) || !areSlicesEqual(finalSettings.Black, initialSettings.Black) {
		t.Errorf("Settings were not restored correctly. Expected: %+v, Got: %+v", initialSettings, finalSettings)
	}
}

func areSlicesEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}

	sort.Strings(s1)
	sort.Strings(s2)

	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
*/
