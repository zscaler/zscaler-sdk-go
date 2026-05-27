package secure_browsing

import (
	"context"
	"testing"

	"github.com/zscaler/zscaler-sdk-go/v3/tests"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_isolation"
)

func TestSecureBrowsing(t *testing.T) {
	service, err := tests.NewOneAPIClient()
	if err != nil {
		t.Fatalf("Error creating client: %v", err)
	}

	ctx := context.Background()

	// Step 1: read-only listing — verify the supportedBrowserVersions
	// sub-resource decodes into a slice (the API returns one entry per
	// browser type: IE, Chrome, Firefox, Safari, Opera).
	t.Run("GetSupportedBrowserVersions", func(t *testing.T) {
		versions, err := GetSupportedBrowserVersions(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching supported browser versions: %v", err)
		}
		if len(versions) == 0 {
			t.Fatal("Expected at least one SupportedBrowserVersion entry, got empty slice")
		}
		for _, v := range versions {
			t.Logf("SupportedBrowserVersion: browserType=%s versions=%d olderVersions=%d",
				v.BrowserType, len(v.Versions), len(v.OlderVersions))
		}
	})

	// Step 2: GET → tweak conservative knobs → PUT → log. Values match
	// the documented payload (browser_control_settings.json) so we never
	// push something the tenant would reject.
	t.Run("UpdateBrowserControlSettings", func(t *testing.T) {
		settings, err := GetBrowserControlSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching browser control settings: %v", err)
		}

		updated := *settings
		updated.PluginCheckFrequency = "DAILY"
		updated.EnableWarnings = true
		updated.BypassPlugins = []string{"DOTNET", "ACROBAT"}
		updated.BypassApplications = []string{"MSOFFICE", "OUTLOOKEXP"}
		updated.BlockedChromeVersions = []string{"NONE"}
		updated.BlockedFirefoxVersions = []string{"NONE"}
		updated.BlockedInternetExplorerVersions = []string{"NONE"}
		updated.BlockedOperaVersions = []string{"NONE"}
		updated.BlockedSafariVersions = []string{"NONE"}

		result, _, err := UpdateBrowserControlSettings(ctx, service, updated)
		if err != nil {
			t.Fatalf("Error updating browser control settings: %v", err)
		}
		t.Logf("Successfully updated BrowserControlSettings: %+v", result)
	})

	// Step 3: smart-isolation PUT. The package does not expose a
	// dedicated GetSmartIsolation, so we seed the payload from the
	// parent BrowserControlSettings GET — the smartIsolation endpoint
	// expects a full replacement body and the parent record carries the
	// authoritative current state for these fields.
	//
	// EnableSmartBrowserIsolation is intentionally NOT toggled here: a
	// `true` value requires a valid smartIsolationProfileId, which is
	// tenant-specific and not safe to fabricate in an integration test.

	cbiProfileList, err := browser_isolation.GetAll(context.Background(), service)
	if err != nil {
		t.Errorf("Error getting cbi profile: %v", err)
		return
	}
	if len(cbiProfileList) == 0 {
		t.Error("Expected retrieved cbi profile to be non-empty, but got empty slice")
	}

	t.Run("UpdateSmartIsolation", func(t *testing.T) {
		parent, err := GetBrowserControlSettings(ctx, service)
		if err != nil {
			t.Fatalf("Error fetching browser control settings: %v", err)
		}

		smart := SmartIsolation{
			PluginCheckFrequency:            "DAILY",
			BypassPlugins:                   parent.BypassPlugins,
			BypassApplications:              parent.BypassApplications,
			BypassAllBrowsers:               parent.BypassAllBrowsers,
			BlockedInternetExplorerVersions: parent.BlockedInternetExplorerVersions,
			BlockedChromeVersions:           parent.BlockedChromeVersions,
			BlockedFirefoxVersions:          parent.BlockedFirefoxVersions,
			BlockedSafariVersions:           parent.BlockedSafariVersions,
			BlockedOperaVersions:            parent.BlockedOperaVersions,
			AllowAllBrowsers:                parent.AllowAllBrowsers,
			EnableWarnings:                  parent.EnableWarnings,
			EnableSmartBrowserIsolation:     parent.EnableSmartBrowserIsolation,
			SmartIsolationProfileID:         parent.SmartIsolationProfileID,
			SmartIsolationProfile: &SmartIsolationProfile{
				ID:   cbiProfileList[0].ID,
				Name: cbiProfileList[0].Name,
				URL:  cbiProfileList[0].URL,
			},
		}

		result, _, err := UpdateSmartIsolation(ctx, service, smart)
		if err != nil {
			t.Fatalf("Error updating smart isolation: %v", err)
		}
		t.Logf("Successfully updated SmartIsolation: %+v", result)
	})
}
