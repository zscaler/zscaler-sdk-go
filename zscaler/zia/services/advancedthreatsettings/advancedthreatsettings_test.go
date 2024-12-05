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
		if strings.HasPrefix(url, "site") { // Adjust prefix condition as needed
			maliciousUrlsToRemove = append(maliciousUrlsToRemove, url)
		}
	}
	if len(maliciousUrlsToRemove) > 0 {
		log.Printf("Cleaning malicious URLs: %v", maliciousUrlsToRemove)
		_, err := UpdateMaliciousURLs(ctx, service, []string{}) // Clear the list
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
		if strings.HasPrefix(url, "site") { // Adjust prefix condition as needed
			securityUrlsToRemove = append(securityUrlsToRemove, url)
		}
	}
	if len(securityUrlsToRemove) > 0 {
		log.Printf("Cleaning security exception URLs: %v", securityUrlsToRemove)
		_, err := UpdateSecurityExceptions(ctx, service, []string{}) // Clear the list
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

	// Step 1: Retrieve and update malicious URLs
	t.Run("UpdateMaliciousURLs", func(t *testing.T) {
		initialUrls, err := GetMaliciousURLs(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching initial malicious URLs: %v", err)
		}

		// Generate new random URLs and update
		newUrls := generateRandomUrls(3)
		allUrls := append(initialUrls.MaliciousUrls, newUrls...)
		_, err = UpdateMaliciousURLs(ctx, service, allUrls) // Passing the []string directly
		if err != nil {
			t.Fatalf("Error updating malicious URLs: %v", err)
		}

		// Verify the update
		updatedUrls, err := GetMaliciousURLs(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching updated malicious URLs: %v", err)
		}
		t.Logf("Updated malicious URLs: %v", updatedUrls.MaliciousUrls)
	})

	// Step 2: Retrieve and update security exceptions
	t.Run("UpdateSecurityExceptions", func(t *testing.T) {
		initialExceptions, err := GetSecurityExceptions(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching initial security exceptions: %v", err)
		}

		// Generate new random URLs and update
		newExceptions := generateRandomUrls(3)
		allExceptions := append(initialExceptions.BypassUrls, newExceptions...)
		_, err = UpdateSecurityExceptions(ctx, service, allExceptions) // Passing the []string directly
		if err != nil {
			t.Fatalf("Error updating security exceptions: %v", err)
		}

		// Verify the update
		updatedExceptions, err := GetSecurityExceptions(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching updated security exceptions: %v", err)
		}
		t.Logf("Updated security exceptions: %v", updatedExceptions.BypassUrls)
	})

	// Step 3: Retrieve and update advanced threat settings
	t.Run("UpdateAdvancedThreatSettings", func(t *testing.T) {
		settings, err := GetAdvancedThreatSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching advanced threat settings: %v", err)
		}

		updatedSettings := *settings
		updatedSettings.RiskTolerance = 50
		updatedSettings.RiskToleranceCapture = false
		updatedSettings.CmdCtlServerBlocked = true
		updatedSettings.CmdCtlServerCapture = false
		updatedSettings.CmdCtlTrafficBlocked = true
		updatedSettings.CmdCtlTrafficCapture = false
		updatedSettings.MalwareSitesBlocked = false
		updatedSettings.MalwareSitesCapture = false
		updatedSettings.ActiveXBlocked = true
		updatedSettings.ActiveXCapture = false
		updatedSettings.BrowserExploitsBlocked = true
		updatedSettings.BrowserExploitsCapture = false
		updatedSettings.FileFormatVulnerabilitiesBlocked = true
		updatedSettings.FileFormatVulnerabilitiesCapture = false
		updatedSettings.KnownPhishingSitesBlocked = true
		updatedSettings.KnownPhishingSitesCapture = false
		updatedSettings.SuspectedPhishingSitesBlocked = true
		updatedSettings.SuspectedPhishingSitesCapture = false
		updatedSettings.SuspectAdwareSpywareSitesBlocked = true
		updatedSettings.SuspectAdwareSpywareSitesCapture = false
		updatedSettings.WebspamBlocked = true
		updatedSettings.WebspamCapture = false
		updatedSettings.IrcTunnellingBlocked = true
		updatedSettings.IrcTunnellingCapture = false
		updatedSettings.AnonymizerBlocked = true
		updatedSettings.AnonymizerCapture = false
		updatedSettings.CookieStealingBlocked = true
		updatedSettings.CookieStealingPCAPEnabled = false
		updatedSettings.PotentialMaliciousRequestsBlocked = true
		updatedSettings.PotentialMaliciousRequestsCapture = false
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