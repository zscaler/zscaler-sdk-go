// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/browser_control_settings"
)

const browserControlSettingsPath = "/zia/api/v1/browserControlSettings"

// =====================================================
// SDK Function Tests
// =====================================================

func TestBrowserControlSettings_GetBrowserControlSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", browserControlSettingsPath, common.SuccessResponse(browser_control_settings.BrowserControlSettings{
		PluginCheckFrequency:        "DAILY",
		BypassPlugins:               []string{"FLASH", "JAVA"},
		BypassApplications:          []string{"APP1"},
		BlockedChromeVersions:       []string{"OLD"},
		EnableWarnings:              true,
		EnableSmartBrowserIsolation: true,
		SmartIsolationProfileID:     42,
		SmartIsolationProfile: browser_control_settings.SmartIsolationProfile{
			ID:             "profile-uuid",
			Name:           "Default Isolation",
			URL:            "https://isolation.example.com",
			DefaultProfile: true,
		},
		SmartIsolationUsers: []ziacommon.IDNameExtensions{
			{ID: 1, Name: "user@example.com"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_control_settings.GetBrowserControlSettings(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "DAILY", result.PluginCheckFrequency)
	assert.True(t, result.EnableWarnings)
	assert.True(t, result.EnableSmartBrowserIsolation)
	assert.Equal(t, 42, result.SmartIsolationProfileID)
	assert.Equal(t, "profile-uuid", result.SmartIsolationProfile.ID)
	assert.Len(t, result.SmartIsolationUsers, 1)
}

func TestBrowserControlSettings_GetBrowserControlSettings_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", browserControlSettingsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := browser_control_settings.GetBrowserControlSettings(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestBrowserControlSettings_UpdateBrowserControlSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", browserControlSettingsPath, common.SuccessResponse(browser_control_settings.BrowserControlSettings{
		PluginCheckFrequency: "WEEKLY",
		EnableWarnings:       true,
		AllowAllBrowsers:     false,
		BypassAllBrowsers:    false,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := browser_control_settings.BrowserControlSettings{
		PluginCheckFrequency: "WEEKLY",
		EnableWarnings:       true,
	}

	result, _, err := browser_control_settings.UpdateBrowserControlSettings(context.Background(), service, settings)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "WEEKLY", result.PluginCheckFrequency)
	assert.True(t, result.EnableWarnings)
}

func TestBrowserControlSettings_UpdateBrowserControlSettings_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", browserControlSettingsPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := browser_control_settings.BrowserControlSettings{
		EnableWarnings: true,
	}

	result, _, err := browser_control_settings.UpdateBrowserControlSettings(context.Background(), service, settings)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestBrowserControlSettings_UpdateBrowserControlSettings_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Empty 204 body makes UpdateWithPut return (nil, nil); the service
	// type-asserts the response and surfaces "unexpected response type".
	server.On("PUT", browserControlSettingsPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := browser_control_settings.BrowserControlSettings{
		EnableWarnings: true,
	}

	result, _, err := browser_control_settings.UpdateBrowserControlSettings(context.Background(), service, settings)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unexpected response type")
}

// =====================================================
// Structure Tests
// =====================================================

func TestBrowserControlSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BrowserControlSettings JSON marshaling", func(t *testing.T) {
		settings := browser_control_settings.BrowserControlSettings{
			PluginCheckFrequency:            "EVERY_4_HOURS",
			BypassPlugins:                   []string{"FLASH"},
			BypassApplications:              []string{"APP1", "APP2"},
			BlockedInternetExplorerVersions: []string{"IE11"},
			BlockedChromeVersions:             []string{"CHROME_OLD"},
			BlockedFirefoxVersions:          []string{"FF_OLD"},
			BlockedSafariVersions:           []string{"SAFARI_OLD"},
			BlockedOperaVersions:            []string{"OPERA_OLD"},
			SmartIsolationUsers: []ziacommon.IDNameExtensions{
				{ID: 10, Name: "user@example.com", Extensions: map[string]interface{}{"key": "val"}},
			},
			SmartIsolationGroups: []ziacommon.IDNameExtensions{
				{ID: 20, Name: "Engineering"},
			},
			SmartIsolationProfile: browser_control_settings.SmartIsolationProfile{
				ID:             "uuid-123",
				Name:           "Corp Isolation",
				URL:            "https://iso.example.com",
				DefaultProfile: false,
			},
			BypassAllBrowsers:           true,
			AllowAllBrowsers:            false,
			EnableWarnings:              true,
			EnableSmartBrowserIsolation: true,
			SmartIsolationProfileID:     99,
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"pluginCheckFrequency":"EVERY_4_HOURS"`)
		assert.Contains(t, string(data), `"bypassAllBrowsers":true`)
		assert.Contains(t, string(data), `"enableWarnings":true`)
		assert.Contains(t, string(data), `"enableSmartBrowserIsolation":true`)
		assert.Contains(t, string(data), `"smartIsolationProfileId":99`)
		assert.Contains(t, string(data), `"blockedChromeVersions"`)
	})

	t.Run("BrowserControlSettings JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"pluginCheckFrequency": "MONTHLY",
			"bypassPlugins": ["JAVA"],
			"blockedFirefoxVersions": ["LEGACY"],
			"smartIsolationProfile": {
				"id": "abc",
				"name": "Profile A",
				"url": "https://a.example.com",
				"defaultProfile": true
			},
			"allowAllBrowsers": false,
			"enableWarnings": false,
			"enableSmartBrowserIsolation": false,
			"smartIsolationProfileId": 7
		}`

		var settings browser_control_settings.BrowserControlSettings
		err := json.Unmarshal([]byte(jsonData), &settings)
		require.NoError(t, err)

		assert.Equal(t, "MONTHLY", settings.PluginCheckFrequency)
		assert.Equal(t, "abc", settings.SmartIsolationProfile.ID)
		assert.True(t, settings.SmartIsolationProfile.DefaultProfile)
		assert.False(t, settings.EnableWarnings)
		assert.Equal(t, 7, settings.SmartIsolationProfileID)
	})

	t.Run("SmartIsolationProfile JSON marshaling", func(t *testing.T) {
		profile := browser_control_settings.SmartIsolationProfile{
			ID:             "iso-id",
			Name:           "Isolation Profile",
			URL:            "https://isolation.zscaler.com",
			DefaultProfile: true,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"iso-id"`)
		assert.Contains(t, string(data), `"defaultProfile":true`)
	})

	t.Run("SmartIsolationProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "profile-1",
			"name": "Default",
			"url": "https://default.example.com",
			"defaultProfile": false
		}`

		var profile browser_control_settings.SmartIsolationProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, "profile-1", profile.ID)
		assert.Equal(t, "Default", profile.Name)
		assert.False(t, profile.DefaultProfile)
	})
}

func TestBrowserControlSettings_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse full browser control settings response", func(t *testing.T) {
		jsonResponse := `{
			"pluginCheckFrequency": "DAILY",
			"bypassPlugins": ["FLASH"],
			"bypassApplications": ["APP"],
			"blockedChromeVersions": ["OLD_CHROME"],
			"smartIsolationUsers": [{"id": 1, "name": "admin@example.com"}],
			"smartIsolationGroups": [{"id": 2, "name": "Admins"}],
			"smartIsolationProfile": {
				"id": "p1",
				"name": "Profile",
				"url": "https://p.example.com",
				"defaultProfile": true
			},
			"bypassAllBrowsers": false,
			"allowAllBrowsers": true,
			"enableWarnings": true,
			"enableSmartBrowserIsolation": true,
			"smartIsolationProfileId": 55
		}`

		var settings browser_control_settings.BrowserControlSettings
		err := json.Unmarshal([]byte(jsonResponse), &settings)
		require.NoError(t, err)

		assert.Equal(t, "DAILY", settings.PluginCheckFrequency)
		assert.True(t, settings.AllowAllBrowsers)
		assert.Len(t, settings.SmartIsolationUsers, 1)
		assert.Equal(t, 55, settings.SmartIsolationProfileID)
	})
}
