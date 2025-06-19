package browser_control_settings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	browserControlSettingsEndpoint = "/zia/api/v1/browserControlSettings"
)

type BrowserControlSettings struct {
	// Specifies how frequently the service checks browsers and relevant applications to warn users regarding outdated or vulnerable browsers, plugins, and applications. If not set, the warnings are disabled
	// Supported Values: DAILY, WEEKLY, MONTHLY, EVERY_2_HOURS, EVERY_4_HOURS, EVERY_6_HOURS, EVERY_8_HOURS, EVERY_12_HOURS
	PluginCheckFrequency string `json:"pluginCheckFrequency,omitempty"`

	// List of plugins that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable plugins are warned.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BypassPlugins []string `json:"bypassPlugins,omitempty"`

	// List of applications that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable applications are warned.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BypassApplications []string `json:"bypassApplications,omitempty"`

	// Versions of Microsoft browser that need to be blocked. If not set, all Microsoft browser versions are allowed.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BlockedInternetExplorerVersions []string `json:"blockedInternetExplorerVersions,omitempty"`

	// Versions of Google Chrome browser that need to be blocked. If not set, all Google Chrome versions are allowed.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BlockedChromeVersions []string `json:"blockedChromeVersions,omitempty"`

	// Versions of Mozilla Firefox browser that need to be blocked. If not set, all Mozilla Firefox versions are allowed.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BlockedFirefoxVersions []string `json:"blockedFirefoxVersions,omitempty"`

	// Versions of Apple Safari browser that need to be blocked. If not set, all Apple Safari versions are allowed.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BlockedSafariVersions []string `json:"blockedSafariVersions,omitempty"`

	// Versions of Opera browser that need to be blocked. If not set, all Opera versions are allowed.
	// See Browser Control API Reference for supported values: https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
	BlockedOperaVersions []string `json:"blockedOperaVersions,omitempty"`

	// Name-ID pairs of users for which the rule is applied
	SmartIsolationUsers []common.IDNameExtensions `json:"smartIsolationUsers,omitempty"`

	// Name-ID pairs of groups for which the rule is applied
	SmartIsolationGroups []common.IDNameExtensions `json:"smartIsolationGroups,omitempty"`

	// The isolation profile.
	// See https://help.zscaler.com/zia/creating-isolation-profiles-zia
	SmartIsolationProfile SmartIsolationProfile `json:"smartIsolationProfile,omitempty"`

	// If set to true, all the browsers are bypassed for warnings.
	// If not set, all vulnerable browsers are warned. This attribute has effect only if the 'enableWarnings' is set to true.
	BypassAllBrowsers bool `json:"bypassAllBrowsers,omitempty"`

	// A Boolean value that specifies whether or not to allow all the browsers and their respective versions access to the internet
	AllowAllBrowsers bool `json:"allowAllBrowsers,omitempty"`

	// A Boolean value that specifies if the warnings are enabled
	EnableWarnings bool `json:"enableWarnings,omitempty"`

	// A Boolean value that specifies if Smart Browser Isolation is enabled
	EnableSmartBrowserIsolation bool `json:"enableSmartBrowserIsolation,omitempty"`

	// The isolation profile ID
	SmartIsolationProfileID int `json:"smartIsolationProfileId,omitempty"`
}

type SmartIsolationProfile struct {
	// The universally unique identifier (UUID) for the browser isolation profile
	ID string `json:"id"`

	// Name of the browser isolation profile
	Name string `json:"name"`

	// The browser isolation profile URL
	URL string `json:"url"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile"`
}

func GetBrowserControlSettings(ctx context.Context, service *zscaler.Service) (*BrowserControlSettings, error) {
	var browserSettings BrowserControlSettings
	err := service.Client.Read(ctx, browserControlSettingsEndpoint, &browserSettings)
	if err != nil {
		return nil, err
	}
	return &browserSettings, nil
}

func UpdateBrowserControlSettings(ctx context.Context, service *zscaler.Service, settings BrowserControlSettings) (*BrowserControlSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, browserControlSettingsEndpoint, settings)
	if err != nil {
		return nil, nil, err
	}

	browserSettings, ok := resp.(*BrowserControlSettings)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated Browser Control Settings : %+v", browserSettings)
	return browserSettings, nil, nil
}
