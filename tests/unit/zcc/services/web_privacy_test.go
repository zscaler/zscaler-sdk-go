// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/web_privacy"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestWebPrivacy_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getWebPrivacyInfo"

	server.On("GET", path, common.SuccessResponse(web_privacy.WebPrivacyInfo{
		ID:                            "privacy-001",
		Active:                        "true",
		CollectMachineHostname:        "true",
		CollectUserInfo:               "true",
		CollectZdxLocation:            "true",
		DisableCrashlytics:            "false",
		EnablePacketCapture:           "true",
		ExportLogsForNonAdmin:         "false",
		GrantAccessToZscalerLogFolder: "true",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_privacy.GetWebPrivacyInfo(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "privacy-001", result.ID)
	assert.Equal(t, "true", result.Active)
}

func TestWebPrivacy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/setWebPrivacyInfo"

	server.On("PUT", path, common.SuccessResponse(web_privacy.WebPrivacyInfo{
		ID:                     "privacy-001",
		Active:                 "false",
		CollectMachineHostname: "false",
		CollectUserInfo:        "false",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateInfo := &web_privacy.WebPrivacyInfo{
		ID:                     "privacy-001",
		Active:                 "false",
		CollectMachineHostname: "false",
	}

	result, err := web_privacy.UpdatePrivacyInfo(context.Background(), service, updateInfo)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "false", result.Active)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestWebPrivacy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebPrivacyInfo JSON marshaling", func(t *testing.T) {
		privacy := web_privacy.WebPrivacyInfo{
			ID:                            "privacy-123",
			Active:                        "true",
			CollectMachineHostname:        "true",
			CollectUserInfo:               "true",
			CollectZdxLocation:            "true",
			DisableCrashlytics:            "false",
			EnablePacketCapture:           "true",
			ExportLogsForNonAdmin:         "false",
			GrantAccessToZscalerLogFolder: "true",
			OverrideT2ProtocolSetting:     "false",
			RestrictRemotePacketCapture:   "true",
		}

		data, err := json.Marshal(privacy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"privacy-123"`)
		assert.Contains(t, string(data), `"active":"true"`)
		assert.Contains(t, string(data), `"collectMachineHostname":"true"`)
		assert.Contains(t, string(data), `"collectUserInfo":"true"`)
	})

	t.Run("WebPrivacyInfo JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "privacy-789",
			"active": "false",
			"collectMachineHostname": "false",
			"collectUserInfo": "false",
			"collectZdxLocation": "true",
			"disableCrashlytics": "true",
			"enablePacketCapture": "false"
		}`

		var privacy web_privacy.WebPrivacyInfo
		err := json.Unmarshal([]byte(jsonData), &privacy)
		require.NoError(t, err)

		assert.Equal(t, "privacy-789", privacy.ID)
		assert.Equal(t, "false", privacy.Active)
		assert.Equal(t, "false", privacy.CollectMachineHostname)
		assert.Equal(t, "true", privacy.DisableCrashlytics)
	})
}
