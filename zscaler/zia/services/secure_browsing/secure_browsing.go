package securebrowsing

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const browserControlSettingsEndpoint = "/zia/api/v1/browserControlSettings"

type SupportedBrowserVersion struct {
	BrowserType   string   `json:"browserType,omitempty"`
	Versions      []string `json:"versions,omitempty"`
	OlderVersions []string `json:"olderVersions,omitempty"`
}

type SmartIsolation struct {
	PluginCheckFrequency            string                    `json:"pluginCheckFrequency,omitempty"`
	BypassPlugins                   []string                  `json:"bypassPlugins,omitempty"`
	BypassApplications              []string                  `json:"bypassApplications,omitempty"`
	BypassAllBrowsers               bool                      `json:"bypassAllBrowsers,omitempty"`
	BlockedInternetExplorerVersions []string                  `json:"blockedInternetExplorerVersions,omitempty"`
	BlockedChromeVersions           []string                  `json:"blockedChromeVersions,omitempty"`
	BlockedFirefoxVersions          []string                  `json:"blockedFirefoxVersions,omitempty"`
	BlockedSafariVersions           []string                  `json:"blockedSafariVersions,omitempty"`
	BlockedOperaVersions            []string                  `json:"blockedOperaVersions,omitempty"`
	AllowAllBrowsers                bool                      `json:"allowAllBrowsers,omitempty"`
	EnableWarnings                  bool                      `json:"enableWarnings,omitempty"`
	EnableSmartBrowserIsolation     bool                      `json:"enableSmartBrowserIsolation,omitempty"`
	SmartIsolationUsers             []common.IDNameExtensions `json:"smartIsolationUsers,omitempty"`
	SmartIsolationGroups            []common.IDNameExtensions `json:"smartIsolationGroups,omitempty"`
	SmartIsolationProfile           *SmartIsolationProfile    `json:"smartIsolationProfile,omitempty"`
	SmartIsolationProfileID         int                       `json:"smartIsolationProfileId,omitempty"`
}

type SmartIsolationProfile struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	URL            string `json:"url,omitempty"`
	DefaultProfile bool   `json:"defaultProfile,omitempty"`
}

type BrowserControlSettings struct {
	EnableSmartBrowserIsolation     bool               `json:"enableSmartBrowserIsolation,omitempty"`
	EnableWarnings                  bool               `json:"enableWarnings,omitempty"`
	PluginCheckFrequency            string             `json:"pluginCheckFrequency,omitempty"`
	BypassPlugins                   []string           `json:"bypassPlugins,omitempty"`
	BypassApplications              []string           `json:"bypassApplications,omitempty"`
	AllowAllBrowsers                bool               `json:"allowAllBrowsers,omitempty"`
	BypassAllBrowsers               bool               `json:"bypassAllBrowsers,omitempty"`
	BlockedChromeVersions           []string           `json:"blockedChromeVersions,omitempty"`
	BlockedFirefoxVersions          []string           `json:"blockedFirefoxVersions,omitempty"`
	BlockedInternetExplorerVersions []string           `json:"blockedInternetExplorerVersions,omitempty"`
	BlockedOperaVersions            []string           `json:"blockedOperaVersions,omitempty"`
	BlockedSafariVersions           []string           `json:"blockedSafariVersions,omitempty"`
	SmartIsolationProfile           *common.CBIProfile `json:"smartIsolationProfile,omitempty"`
	SmartIsolationProfileID         int                `json:"smartIsolationProfileId,omitempty"`
}

// GetSupportedBrowserVersions returns the per-browser supported/older
// version list. The API responds with a JSON array, one element per
// browser type, so this function returns []SupportedBrowserVersion.
func GetSupportedBrowserVersions(ctx context.Context, service *zscaler.Service) ([]SupportedBrowserVersion, error) {
	var supportedBrowserVersions []SupportedBrowserVersion
	err := service.Client.Read(ctx, browserControlSettingsEndpoint+"/supportedBrowserVersions", &supportedBrowserVersions)
	if err != nil {
		return nil, err
	}
	return supportedBrowserVersions, nil
}

func UpdateSmartIsolation(ctx context.Context, service *zscaler.Service, settings SmartIsolation) (*SmartIsolation, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, browserControlSettingsEndpoint+"/smartIsolation", settings)
	if err != nil {
		return nil, nil, err
	}

	smartIsolationSettings, ok := resp.(*SmartIsolation)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated Smart Isolation : %+v", smartIsolationSettings)
	return smartIsolationSettings, nil, nil
}

func GetBrowserControlSettings(ctx context.Context, service *zscaler.Service) (*BrowserControlSettings, error) {
	var browserControlSettings BrowserControlSettings
	err := service.Client.Read(ctx, browserControlSettingsEndpoint, &browserControlSettings)
	if err != nil {
		return nil, err
	}
	return &browserControlSettings, nil
}

func UpdateBrowserControlSettings(ctx context.Context, service *zscaler.Service, settings BrowserControlSettings) (*BrowserControlSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, browserControlSettingsEndpoint, settings)
	if err != nil {
		return nil, nil, err
	}

	browserControlSettings, ok := resp.(*BrowserControlSettings)
	if !ok {
		return nil, nil, fmt.Errorf("unexpected response type")
	}
	service.Client.GetLogger().Printf("[DEBUG] Updated Browser Control Settings : %+v", browserControlSettings)
	return browserControlSettings, nil, nil
}
