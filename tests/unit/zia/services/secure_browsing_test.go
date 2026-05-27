// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	securebrowsing "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/secure_browsing"
)

const secureBrowsingPath = "/zia/api/v1/browserControlSettings"

func TestSecureBrowsing_GetSupportedBrowserVersions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := secureBrowsingPath + "/supportedBrowserVersions"
	server.On("GET", path, common.SuccessResponse([]securebrowsing.SupportedBrowserVersion{
		{BrowserType: "CHROME", Versions: []string{"120", "121"}, OlderVersions: []string{"119"}},
		{BrowserType: "FIREFOX", Versions: []string{"115"}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := securebrowsing.GetSupportedBrowserVersions(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "CHROME", result[0].BrowserType)
}

func TestSecureBrowsing_GetBrowserControlSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", secureBrowsingPath, common.SuccessResponse(securebrowsing.BrowserControlSettings{
		PluginCheckFrequency:            "DAILY",
		EnableWarnings:                    true,
		BypassPlugins:                     []string{"DOTNET", "ACROBAT"},
		BypassApplications:                []string{"MSOFFICE", "OUTLOOKEXP"},
		BlockedChromeVersions:             []string{"NONE"},
		BlockedFirefoxVersions:            []string{"NONE"},
		BlockedInternetExplorerVersions:   []string{"NONE"},
		BlockedOperaVersions:              []string{"NONE"},
		BlockedSafariVersions:             []string{"NONE"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := securebrowsing.GetBrowserControlSettings(context.Background(), service)
	require.NoError(t, err)
	assert.Equal(t, "DAILY", result.PluginCheckFrequency)
	assert.Contains(t, result.BypassPlugins, "DOTNET")
}

func TestSecureBrowsing_UpdateBrowserControlSettings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", secureBrowsingPath, common.SuccessResponse(securebrowsing.BrowserControlSettings{
		PluginCheckFrequency: "DAILY",
		EnableWarnings:       true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := securebrowsing.BrowserControlSettings{
		PluginCheckFrequency: "DAILY",
		EnableWarnings:       true,
		BypassPlugins:        []string{"DOTNET", "ACROBAT"},
	}

	result, _, err := securebrowsing.UpdateBrowserControlSettings(context.Background(), service, settings)
	require.NoError(t, err)
	assert.True(t, result.EnableWarnings)
}

func TestSecureBrowsing_UpdateSmartIsolation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := secureBrowsingPath + "/smartIsolation"
	server.On("PUT", path, common.SuccessResponse(securebrowsing.SmartIsolation{
		PluginCheckFrequency: "DAILY",
		EnableWarnings:       true,
		BypassPlugins:        []string{"DOTNET"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := securebrowsing.SmartIsolation{
		PluginCheckFrequency: "DAILY",
		EnableWarnings:       true,
		BypassPlugins:        []string{"DOTNET"},
	}

	result, _, err := securebrowsing.UpdateSmartIsolation(context.Background(), service, settings)
	require.NoError(t, err)
	assert.Equal(t, "DAILY", result.PluginCheckFrequency)
}

func TestSecureBrowsing_GetBrowserControlSettings_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", secureBrowsingPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := securebrowsing.GetBrowserControlSettings(context.Background(), service)
	require.Error(t, err)
	assert.Nil(t, result)
}
