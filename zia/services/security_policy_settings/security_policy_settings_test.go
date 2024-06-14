package security_policy_settings

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/tests"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

const (
	maxConflictRetries    = 5
	conflictRetryInterval = 1 * time.Second
)

func randomURL() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(".example%d.com", rand.Intn(10000))
}

func retryOnConflict(operation func() error) error {
	var lastErr error
	for i := 0; i < maxConflictRetries; i++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}

		if strings.Contains(lastErr.Error(), `"code":"EDIT_LOCK_NOT_AVAILABLE"`) {
			log.Printf("Conflict error detected, retrying in %v... (Attempt %d/%d)", conflictRetryInterval, i+1, maxConflictRetries)
			time.Sleep(conflictRetryInterval)
			continue
		}

		return lastErr
	}
	return lastErr
}

func TestSecurityPolicySettings(t *testing.T) {
	client, err := tests.NewZiaClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	service := &services.Service{Client: client}

	// Backup initial settings
	initialSettings, err := GetListUrls(service)
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
		_, err := UpdateListUrls(service, newSettings)
		return err
	})
	if err != nil {
		t.Fatalf("Error updating settings: %v", err)
	}

	updatedSettings, err := GetListUrls(service)
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
		_, err := UpdateListUrls(service, *initialSettings)
		return err
	})
	if err != nil {
		t.Fatalf("Error restoring initial settings: %v", err)
	}

	// Verify if the settings were restored
	finalSettings, err := GetListUrls(service)
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
