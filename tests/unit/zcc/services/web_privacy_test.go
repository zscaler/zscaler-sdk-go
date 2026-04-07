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

func TestWebPrivacy_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getWebPrivacyInfo"

	server.On("GET", path, common.SuccessResponse(web_privacy.WebPrivacyInfo{
		ID:                            "7777",
		Active:                        "1",
		CollectUserInfo:               "1",
		CollectMachineHostname:        "1",
		CollectZdxLocation:            "1",
		EnablePacketCapture:           "1",
		DisableCrashlytics:            "1",
		OverrideT2ProtocolSetting:     "1",
		RestrictRemotePacketCapture:   "1",
		GrantAccessToZscalerLogFolder: "1",
		ExportLogsForNonAdmin:         "1",
		EnableAutoLogSnippet:          "1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_privacy.GetWebPrivacyInfo(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "7777", result.ID)
	assert.Equal(t, "1", result.Active)
	assert.Equal(t, "1", result.DisableCrashlytics)
}

func TestWebPrivacy_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	putPath := "/zcc/papi/public/v1/setWebPrivacyInfo"
	getPath := "/zcc/papi/public/v1/getWebPrivacyInfo"

	server.On("PUT", putPath, common.SuccessResponse(nil))
	server.On("GET", getPath, common.SuccessResponse(web_privacy.WebPrivacyInfo{
		ID:                            "7777",
		Active:                        "0",
		CollectMachineHostname:        "0",
		CollectUserInfo:               "0",
		DisableCrashlytics:            "0",
		EnableFQDNMatchForVpnBypasses: "1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateInfo := &web_privacy.WebPrivacyInfo{
		ID:                            "7777",
		Active:                        "0",
		CollectMachineHostname:        "0",
		DisableCrashlytics:            "0",
		EnableFQDNMatchForVpnBypasses: "1",
	}

	result, err := web_privacy.UpdateWebPrivacyInfo(context.Background(), service, updateInfo)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "0", result.Active)
	assert.Equal(t, "0", result.DisableCrashlytics)
	assert.Equal(t, "1", result.EnableFQDNMatchForVpnBypasses)
}

func TestWebPrivacy_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebPrivacyInfo JSON marshaling", func(t *testing.T) {
		privacy := web_privacy.WebPrivacyInfo{
			ID:                            "7777",
			Active:                        "1",
			CollectUserInfo:               "1",
			CollectMachineHostname:        "1",
			CollectZdxLocation:            "1",
			EnablePacketCapture:           "1",
			DisableCrashlytics:            "1",
			OverrideT2ProtocolSetting:     "1",
			RestrictRemotePacketCapture:   "1",
			GrantAccessToZscalerLogFolder: "1",
			ExportLogsForNonAdmin:         "1",
			EnableAutoLogSnippet:          "1",
			EnforceSecurePacUrls:          "1",
			EnableFQDNMatchForVpnBypasses: "0",
		}

		data, err := json.Marshal(privacy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"7777"`)
		assert.Contains(t, string(data), `"active":"1"`)
		assert.Contains(t, string(data), `"collectMachineHostname":"1"`)
		assert.Contains(t, string(data), `"disableCrashlytics":"1"`)
		assert.Contains(t, string(data), `"enableAutoLogSnippet":"1"`)
		assert.Contains(t, string(data), `"enforceSecurePacUrls":"1"`)
		assert.Contains(t, string(data), `"enableFQDNMatchForVpnBypasses":"0"`)
	})

	t.Run("WebPrivacyInfo JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "7777",
			"active": "1",
			"collectUserInfo": "1",
			"collectMachineHostname": "1",
			"collectZdxLocation": "1",
			"enablePacketCapture": "1",
			"disableCrashlytics": "1",
			"overrideT2ProtocolSetting": "1",
			"restrictRemotePacketCapture": "1",
			"grantAccessToZscalerLogFolder": "1",
			"exportLogsForNonAdmin": "1",
			"enableAutoLogSnippet": "1"
		}`

		var privacy web_privacy.WebPrivacyInfo
		err := json.Unmarshal([]byte(jsonData), &privacy)
		require.NoError(t, err)

		assert.Equal(t, "7777", privacy.ID)
		assert.Equal(t, "1", privacy.Active)
		assert.Equal(t, "1", privacy.CollectUserInfo)
		assert.Equal(t, "1", privacy.DisableCrashlytics)
		assert.Equal(t, "1", privacy.OverrideT2ProtocolSetting)
		assert.Equal(t, "1", privacy.RestrictRemotePacketCapture)
		assert.Equal(t, "1", privacy.GrantAccessToZscalerLogFolder)
		assert.Equal(t, "1", privacy.EnableAutoLogSnippet)
	})
}
