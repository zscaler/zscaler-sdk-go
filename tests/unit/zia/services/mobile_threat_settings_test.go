package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/mobile_threat_settings"
)

const mobileThreatSettingsPath = "/zia/api/v1/mobileAdvanceThreatSettings"

func sampleMobileThreatSettings() mobile_threat_settings.MobileAdvanceThreatSettings {
	return mobile_threat_settings.MobileAdvanceThreatSettings{
		BlockAppsWithMaliciousActivity:                 true,
		BlockAppsWithKnownVulnerabilities:              true,
		BlockAppsSendingUnencryptedUserCredentials:     true,
		BlockAppsSendingLocationInfo:                   true,
		BlockAppsSendingPersonallyIdentifiableInfo:       true,
		BlockAppsSendingDeviceIdentifier:               true,
		BlockAppsCommunicatingWithAdWebsites:           true,
		BlockAppsCommunicatingWithRemoteUnknownServers: true,
	}
}

func TestMobileThreatSettings_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", mobileThreatSettingsPath, common.SuccessResponse(sampleMobileThreatSettings()))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := mobile_threat_settings.GetMobileThreatSettings(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.BlockAppsWithMaliciousActivity)
	assert.True(t, result.BlockAppsCommunicatingWithRemoteUnknownServers)
}

func TestMobileThreatSettings_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	updated := sampleMobileThreatSettings()
	updated.BlockAppsWithMaliciousActivity = false

	server.On("PUT", mobileThreatSettingsPath, common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleMobileThreatSettings()
	input.BlockAppsWithMaliciousActivity = false

	result, _, err := mobile_threat_settings.UpdateMobileThreatSettings(context.Background(), service, input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.BlockAppsWithMaliciousActivity)
	assert.True(t, result.BlockAppsWithKnownVulnerabilities)
}

func TestMobileThreatSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		settings := sampleMobileThreatSettings()

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"blockAppsWithMaliciousActivity":true`)
		assert.Contains(t, string(data), `"blockAppsCommunicatingWithAdWebsites":true`)
	})
}
