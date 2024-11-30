package advancedthreatsettings

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	cleanResources()
}

func teardown() {
	cleanResources()
}

func shouldClean() bool {
	val, present := os.LookupEnv("ZSCALER_SDK_TEST_SWEEP")
	return !present || (present && (val == "" || val == "true"))
}

func cleanResources() {
	if !shouldClean() {
		return
	}

	// Define the context here
	ctx := context.Background()

	service, err := tests.NewOneAPIClient()
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Clean malicious URLs
	maliciousResources, err := GetMaliciousURLs(ctx, service)
	if err != nil {
		log.Printf("Error retrieving malicious URLs during cleanup: %v", err)
		return
	}

	var maliciousUrlsToRemove []string
	for _, url := range maliciousResources.MaliciousUrls {
		if strings.HasPrefix(url, "site") {
			maliciousUrlsToRemove = append(maliciousUrlsToRemove, url)
		}
	}
	if len(maliciousUrlsToRemove) > 0 {
		_, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", maliciousUrlsEndpoint), MaliciousURLs{MaliciousUrls: maliciousUrlsToRemove})
		if err != nil {
			log.Printf("Error removing malicious URLs during cleanup: %v", err)
		}
	}

	// Clean security exceptions
	securityResources, err := GetSecurityExceptions(ctx, service)
	if err != nil {
		log.Printf("Error retrieving security exceptions during cleanup: %v", err)
		return
	}

	var securityUrlsToRemove []string
	for _, url := range securityResources.BypassUrls {
		if strings.HasPrefix(url, "site") {
			securityUrlsToRemove = append(securityUrlsToRemove, url)
		}
	}
	if len(securityUrlsToRemove) > 0 {
		_, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", securityExceptionsEndpoint), SecurityExceptions{BypassUrls: securityUrlsToRemove})
		if err != nil {
			log.Printf("Error removing security exceptions during cleanup: %v", err)
		}
	}
}

func TestAdvancedThreatSettings(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Errorf("Error creating client: %v", err)
		return
	}

	ctx := context.Background()

	// Step 1: Retrieve and remove malicious URLs
	t.Run("RemoveMaliciousURLs", func(t *testing.T) {
		initialUrls, err := GetMaliciousURLs(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching initial malicious URLs: %v", err)
		}

		// Generate new random URLs and update
		newUrls := generateRandomUrls(3)
		allUrls := append(initialUrls.MaliciousUrls, newUrls...)
		_, err = UpdateMaliciousURLs(ctx, service, MaliciousURLs{MaliciousUrls: allUrls})
		if err != nil {
			t.Fatalf("Error updating malicious URLs: %v", err)
		}

		// Remove added URLs
		_, err = service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", maliciousUrlsEndpoint), MaliciousURLs{MaliciousUrls: newUrls})
		if err != nil {
			t.Fatalf("Error removing malicious URLs: %v", err)
		}
		t.Logf("Successfully removed malicious URLs: %v", newUrls)
	})

	// Step 2: Retrieve and remove security exceptions
	t.Run("RemoveSecurityExceptions", func(t *testing.T) {
		initialExceptions, err := GetSecurityExceptions(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching initial security exceptions: %v", err)
		}

		// Generate new random URLs and update
		newExceptions := generateRandomUrls(3)
		allExceptions := append(initialExceptions.BypassUrls, newExceptions...)
		_, err = UpdateSecurityExceptions(ctx, service, SecurityExceptions{BypassUrls: allExceptions})
		if err != nil {
			t.Fatalf("Error updating security exceptions: %v", err)
		}

		// Remove added URLs
		_, err = service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s?action=REMOVE_FROM_LIST", securityExceptionsEndpoint), SecurityExceptions{BypassUrls: newExceptions})
		if err != nil {
			t.Fatalf("Error removing security exceptions: %v", err)
		}
		t.Logf("Successfully removed security exceptions: %v", newExceptions)
	})

	// Step 3: Retrieve and update advanced threat settings
	t.Run("UpdateAdvancedThreatSettings", func(t *testing.T) {
		settings, err := GetAdvancedThreatSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching advanced threat settings: %v", err)
		}

		updatedSettings := *settings
		updatedSettings.RiskTolerance = 50
		updatedSettings.RiskToleranceCapture = true
		updatedSettings.CmdCtlServerBlocked = true
		updatedSettings.CmdCtlServerCapture = true
		updatedSettings.CmdCtlTrafficBlocked = true
		updatedSettings.CmdCtlTrafficCapture = true
		updatedSettings.MalwareSitesBlocked = true
		updatedSettings.MalwareSitesCapture = true
		updatedSettings.ActiveXBlocked = true
		updatedSettings.ActiveXCapture = true
		updatedSettings.BrowserExploitsBlocked = true
		updatedSettings.BrowserExploitsCapture = true
		updatedSettings.FileFormatVulnerabilitiesBlocked = true
		updatedSettings.FileFormatVulnerabilitiesCapture = true
		updatedSettings.KnownPhishingSitesBlocked = true
		updatedSettings.KnownPhishingSitesCapture = true
		updatedSettings.SuspectedPhishingSitesBlocked = true
		updatedSettings.SuspectedPhishingSitesCapture = true
		updatedSettings.SuspectAdwareSpywareSitesBlocked = true
		updatedSettings.SuspectAdwareSpywareSitesCapture = true
		updatedSettings.WebspamBlocked = true
		updatedSettings.WebspamCapture = true
		updatedSettings.IrcTunnellingBlocked = true
		updatedSettings.IrcTunnellingCapture = true
		updatedSettings.AnonymizerBlocked = true
		updatedSettings.AnonymizerCapture = true
		updatedSettings.CookieStealingBlocked = true
		updatedSettings.CookieStealingPCAPEnabled = true
		updatedSettings.PotentialMaliciousRequestsBlocked = true
		updatedSettings.PotentialMaliciousRequestsCapture = true
		updatedSettings.BlockedCountries = []string{
			"COUNTRY_CA",
			"COUNTRY_US",
			"COUNTRY_MX",
			"COUNTRY_AU",
			"COUNTRY_GB",
		}
		updatedSettings.BlockCountriesCapture = true
		updatedSettings.BitTorrentBlocked = true
		updatedSettings.BitTorrentCapture = true
		updatedSettings.TorBlocked = true
		updatedSettings.TorCapture = true
		updatedSettings.GoogleTalkBlocked = true
		updatedSettings.GoogleTalkCapture = true
		updatedSettings.SshTunnellingBlocked = true
		updatedSettings.SshTunnellingCapture = true
		updatedSettings.CryptoMiningBlocked = true
		updatedSettings.CryptoMiningCapture = true
		updatedSettings.AdSpywareSitesBlocked = true
		updatedSettings.AdSpywareSitesCapture = true
		updatedSettings.DgaDomainsBlocked = true
		updatedSettings.AlertForUnknownOrSuspiciousC2Traffic = true
		updatedSettings.DgaDomainsCapture = true
		updatedSettings.MaliciousUrlsCapture = true

		result, _, err := UpdateAdvancedThreatSettings(ctx, service, updatedSettings)
		if err != nil {
			t.Fatalf("Error updating advanced threat settings: %v", err)
		}
		t.Logf("Successfully updated advanced threat settings: %+v", result)
	})
}

func generateRandomUrls(count int) []string {
	var urls []string
	domains := []string{".example.com", ".test.com"}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		url := fmt.Sprintf("%s%d%s", "site", rand.Intn(1000), domains[rand.Intn(len(domains))])
		urls = append(urls, url)
	}
	return urls
}

// func contains(slice []string, item string) bool {
// 	for _, s := range slice {
// 		if s == item {
// 			return true
// 		}
// 	}
// 	return false
// }
