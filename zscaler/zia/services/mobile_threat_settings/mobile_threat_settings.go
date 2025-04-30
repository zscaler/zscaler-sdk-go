package mobile_threat_settings

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	MobileThreatSettingsEndpoint = "/zia/api/v1/mobileAdvanceThreatSettings"
)

type MobileAdvanceThreatSettings struct {
	// Blocks applications that are known to be malicious, compromised, or perform activities unknown to or hidden from the user
	BlockAppsWithMaliciousActivity bool `json:"blockAppsWithMaliciousActivity,omitempty"`

	//Blocks applications that contain vulnerabilities or that use insecure features, modules, or protocols
	BlockAppsWithKnownVulnerabilities bool `json:"blockAppsWithKnownVulnerabilities,omitempty"`

	// Blocks an application from leaking a user's credentials in an unencrypted format
	BlockAppsSendingUnencryptedUserCredentials bool `json:"blockAppsSendingUnencryptedUserCredentials,omitempty"`

	// Blocks an application from leaking device location details via communication in an unencrypted format or for an unknown purpose
	BlockAppsSendingLocationInfo bool `json:"blockAppsSendingLocationInfo,omitempty"`

	// Blocks an application from leaking a user's personally identifiable information (PII) via communication in an unencrypted format or for an unknown purpose
	BlockAppsSendingPersonallyIdentifiableInfo bool `json:"blockAppsSendingPersonallyIdentifiableInfo,omitempty"`

	// Blocks an application from leaking device identifiers via communication in an unencrypted format or for an unknown purpose
	BlockAppsSendingDeviceIdentifier bool `json:"blockAppsSendingDeviceIdentifier,omitempty"`

	// Blocks an application from communicating with known advertisement websites
	BlockAppsCommunicatingWithAdWebsites bool `json:"blockAppsCommunicatingWithAdWebsites,omitempty"`

	// Blocks an application from communicating with unknown servers (i.e., servers not normally or historically associated with the application)
	BlockAppsCommunicatingWithRemoteUnknownServers bool `json:"blockAppsCommunicatingWithRemoteUnknownServers,omitempty"`
}

func GetMobileThreatSettings(ctx context.Context, service *zscaler.Service) (*MobileAdvanceThreatSettings, error) {
	var advSettings MobileAdvanceThreatSettings
	err := service.Client.Read(ctx, MobileThreatSettingsEndpoint, &advSettings)
	if err != nil {
		return nil, err
	}
	return &advSettings, nil
}

func UpdateMobileThreatSettings(ctx context.Context, service *zscaler.Service, advancedSettings MobileAdvanceThreatSettings) (*MobileAdvanceThreatSettings, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, (MobileThreatSettingsEndpoint), advancedSettings)
	if err != nil {
		return nil, nil, err
	}
	updatedMobileThreatSettings, _ := resp.(*MobileAdvanceThreatSettings)

	service.Client.GetLogger().Printf("[DEBUG]returning updates mobile threat settings from update: %d", updatedMobileThreatSettings)
	return updatedMobileThreatSettings, nil, nil
}
