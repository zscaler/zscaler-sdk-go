// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/application_profiles"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestApplicationProfiles_GetList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfilesResponse{
		TotalCount: 2,
		Policies: []application_profiles.ApplicationProfile{
			{ID: 1, Name: "Default", Active: 1, DeviceType: "DEVICE_TYPE_IOS"},
			{ID: 2, Name: "Custom Policy", Active: 1, DeviceType: "DEVICE_TYPE_WINDOWS"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.Policies, 2)
	assert.Equal(t, "Default", result.Policies[0].Name)
}

func TestApplicationProfiles_GetList_WithDeviceTypeName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfilesResponse{
		TotalCount: 1,
		Policies: []application_profiles.ApplicationProfile{
			{ID: 3, Name: "Windows Policy", Active: 1, DeviceType: "DEVICE_TYPE_WINDOWS"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "windows", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.TotalCount)
	assert.Equal(t, "Windows Policy", result.Policies[0].Name)
}

func TestApplicationProfiles_GetList_InvalidDeviceType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = application_profiles.GetApplicationProfiles(context.Background(), service, "", "", "invalid_os", nil, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid deviceType")
}

func TestApplicationProfiles_GetByProfileID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles/42"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfile{
		ID:          42,
		Name:        "Enterprise Policy",
		Description: "Main enterprise policy",
		Active:      1,
		DeviceType:  "DEVICE_TYPE_WINDOWS",
		LogMode:     3,
		PolicyExtension: application_profiles.PolicyExtension{
			FollowRoutingTable:  "1",
			UseV8JsEngine:       "1",
			AdvanceZpaReauth:    false,
			DropQuicTraffic:     0,
			EnableAntiTampering: "0",
		},
		DisasterRecovery: application_profiles.DisasterRecovery{
			EnableZiaDR:  false,
			EnableZpaDR:  false,
			AllowZiaTest: false,
			AllowZpaTest: false,
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := application_profiles.GetByProfileID(context.Background(), service, "42")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 42, result.ID)
	assert.Equal(t, "Enterprise Policy", result.Name)
	assert.Equal(t, "1", result.PolicyExtension.FollowRoutingTable)
	assert.False(t, result.DisasterRecovery.EnableZiaDR)
}

func TestApplicationProfiles_GetByProfileID_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = application_profiles.GetByProfileID(context.Background(), service, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "profileId is required")
}

func TestApplicationProfiles_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfilesResponse{
		TotalCount: 3,
		Policies: []application_profiles.ApplicationProfile{
			{ID: 1, Name: "Default", Active: 1},
			{ID: 2, Name: "Enterprise Policy", Active: 1},
			{ID: 3, Name: "Branch Policy", Active: 0},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := application_profiles.GetByName(context.Background(), service, "Enterprise Policy")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, "Enterprise Policy", result.Name)
}

func TestApplicationProfiles_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfilesResponse{
		TotalCount: 1,
		Policies: []application_profiles.ApplicationProfile{
			{ID: 5, Name: "Enterprise Policy", Active: 1},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := application_profiles.GetByName(context.Background(), service, "enterprise policy")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 5, result.ID)
}

func TestApplicationProfiles_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles"

	server.On("GET", path, common.SuccessResponse(application_profiles.ApplicationProfilesResponse{
		TotalCount: 1,
		Policies: []application_profiles.ApplicationProfile{
			{ID: 1, Name: "Default", Active: 1},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = application_profiles.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no application profile found with name: NonExistent")
}

func TestApplicationProfiles_Patch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/application-profiles/42"

	server.On("PATCH", path, common.SuccessResponse(application_profiles.ApplicationProfile{
		ID:          42,
		Name:        "Updated Policy",
		Description: "Updated description",
		Active:      1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	patch := &application_profiles.ApplicationProfile{
		Description: "Updated description",
	}

	result, _, err := application_profiles.PatchApplicationProfile(context.Background(), service, "42", patch)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 42, result.ID)
	assert.Equal(t, "Updated description", result.Description)
}

func TestApplicationProfiles_Patch_EmptyProfileID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	patch := &application_profiles.ApplicationProfile{Description: "test"}
	_, _, err = application_profiles.PatchApplicationProfile(context.Background(), service, "", patch)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "profileId is required")
}

func TestApplicationProfiles_Patch_NilBody_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = application_profiles.PatchApplicationProfile(context.Background(), service, "42", nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "patch body is required")
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestApplicationProfiles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationProfile JSON marshaling", func(t *testing.T) {
		profile := application_profiles.ApplicationProfile{
			ID:          1,
			Name:        "Default",
			Description: "Default Policy",
			Active:      1,
			DeviceType:  "DEVICE_TYPE_IOS",
			RuleOrder:   1,
			LogMode:     3,
			LogFileSize: 100,
			GroupAll:    1,
			Ipv6Mode:    4,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1`)
		assert.Contains(t, string(data), `"name":"Default"`)
		assert.Contains(t, string(data), `"active":1`)
		assert.Contains(t, string(data), `"deviceType":"DEVICE_TYPE_IOS"`)
		assert.Contains(t, string(data), `"ipv6Mode":4`)
	})

	t.Run("ApplicationProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 42,
			"name": "Enterprise",
			"description": "Enterprise Policy",
			"active": 1,
			"deviceType": "DEVICE_TYPE_WINDOWS",
			"ruleOrder": 2,
			"logMode": 3,
			"logLevel": 0,
			"logFileSize": 100,
			"reactivateWebSecurityMinutes": "0",
			"tunnelZappTraffic": 0,
			"groupAll": 1,
			"passcode": "",
			"showVPNTunNotification": 0,
			"ipv6Mode": 4
		}`

		var profile application_profiles.ApplicationProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, 42, profile.ID)
		assert.Equal(t, "Enterprise", profile.Name)
		assert.Equal(t, 1, profile.Active)
		assert.Equal(t, "DEVICE_TYPE_WINDOWS", profile.DeviceType)
		assert.Equal(t, "0", profile.ReactivateWebSecurityMinutes)
		assert.Equal(t, 4, profile.Ipv6Mode)
	})

	t.Run("PolicyExtension JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"sourcePortBasedBypasses": "3389:*",
			"packetTunnelExcludeList": "10.0.0.0/8",
			"packetTunnelIncludeList": "0.0.0.0/0",
			"useV8JsEngine": "1",
			"followRoutingTable": "1",
			"advanceZpaReauth": false,
			"enableSetProxyOnVPNAdapters": 1,
			"dropQuicTraffic": 0,
			"enableAntiTampering": "0",
			"browserAuthType": "-1",
			"useDefaultBrowser": "0"
		}`

		var ext application_profiles.PolicyExtension
		err := json.Unmarshal([]byte(jsonData), &ext)
		require.NoError(t, err)

		assert.Equal(t, "3389:*", ext.SourcePortBasedBypasses)
		assert.Equal(t, "1", ext.UseV8JsEngine)
		assert.Equal(t, "1", ext.FollowRoutingTable)
		assert.False(t, ext.AdvanceZpaReauth)
		assert.Equal(t, 1, ext.EnableSetProxyOnVPNAdapters)
		assert.Equal(t, "-1", ext.BrowserAuthType)
	})

	t.Run("DisasterRecovery JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"enableZiaDR": false,
			"enableZpaDR": false,
			"ziaDRMethod": 0,
			"ziaCustomDbUrl": "",
			"useZiaGlobalDb": false,
			"ziaGlobalDbUrl": "",
			"ziaGlobalDbUrlv2": "",
			"ziaDomainName": "",
			"ziaRSAPubKeyName": "",
			"ziaRSAPubKey": "",
			"zpaDomainName": "",
			"zpaRSAPubKeyName": "",
			"zpaRSAPubKey": "",
			"allowZiaTest": false,
			"allowZpaTest": false
		}`

		var dr application_profiles.DisasterRecovery
		err := json.Unmarshal([]byte(jsonData), &dr)
		require.NoError(t, err)

		assert.False(t, dr.EnableZiaDR)
		assert.False(t, dr.EnableZpaDR)
		assert.False(t, dr.UseZiaGlobalDb)
		assert.False(t, dr.AllowZiaTest)
		assert.False(t, dr.AllowZpaTest)
		assert.Equal(t, 0, dr.ZiaDRMethod)
	})

	t.Run("GenerateCliPasswordContract JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"enableCli": false,
			"allowZpaDisableWithoutPassword": false,
			"allowZiaDisableWithoutPassword": false,
			"allowZdxDisableWithoutPassword": false
		}`

		var contract application_profiles.GenerateCliPasswordContract
		err := json.Unmarshal([]byte(jsonData), &contract)
		require.NoError(t, err)

		assert.False(t, contract.EnableCli)
		assert.False(t, contract.AllowZpaDisableWithoutPassword)
		assert.False(t, contract.AllowZiaDisableWithoutPassword)
		assert.False(t, contract.AllowZdxDisableWithoutPassword)
	})

	t.Run("LocationRulesetPolicies JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"offTrusted": {"id": 0, "name": ""},
			"trusted": {"id": 1, "name": "Trust Policy"},
			"vpnTrusted": {"id": 0, "name": ""},
			"splitVpnTrusted": {"id": 0, "name": ""}
		}`

		var lrp application_profiles.LocationRulesetPolicies
		err := json.Unmarshal([]byte(jsonData), &lrp)
		require.NoError(t, err)

		assert.Equal(t, 1, lrp.Trusted.ID)
		assert.Equal(t, "Trust Policy", lrp.Trusted.Name)
	})

	t.Run("ApplicationProfilesResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 2,
			"policies": [
				{"id": 1, "name": "Default", "active": 1, "deviceType": "DEVICE_TYPE_IOS"},
				{"id": 2, "name": "Custom", "active": 0, "deviceType": "DEVICE_TYPE_WINDOWS"}
			]
		}`

		var response application_profiles.ApplicationProfilesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.TotalCount)
		assert.Len(t, response.Policies, 2)
		assert.Equal(t, "Default", response.Policies[0].Name)
		assert.Equal(t, 1, response.Policies[0].Active)
		assert.Equal(t, "Custom", response.Policies[1].Name)
		assert.Equal(t, 0, response.Policies[1].Active)
	})

	t.Run("API list response with groups, users, string groupIds", func(t *testing.T) {
		jsonData := `{
			"totalCount": 1,
			"policies": [{
				"deviceType": "DEVICE_TYPE_ANDROID",
				"id": 171803,
				"name": "AndroidPolicy01",
				"active": 0,
				"groups": [
					{"id": 62718389, "name": "A001", "authType": "SAFECHANNEL_DIR", "active": 1, "lastModification": "1691828147"}
				],
				"users": [
					{"id": "5807211", "loginName": "user@example.com", "lastModification": "2026-02-19 00:49:32.0", "active": 1, "companyId": "4543"}
				],
				"groupIds": ["62718389"],
				"userIds": ["5807211", "35129345"],
				"forwardingProfileId": null,
				"reauth_period": "12",
				"uninstall_password": null
			}]
		}`

		var response application_profiles.ApplicationProfilesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		require.Len(t, response.Policies, 1)
		p := response.Policies[0]
		require.Len(t, p.Groups, 1)
		assert.Equal(t, int64(62718389), p.Groups[0].ID)
		assert.Equal(t, "A001", p.Groups[0].Name)
		require.Len(t, p.Users, 1)
		assert.Equal(t, "5807211", p.Users[0].ID)
		assert.Equal(t, "user@example.com", p.Users[0].LoginName)
		assert.Equal(t, []string{"62718389"}, p.GroupIds)
		assert.Equal(t, []string{"5807211", "35129345"}, p.UserIds)
		assert.Nil(t, p.ForwardingProfileId)
		require.NotNil(t, p.ReauthPeriod)
		assert.Equal(t, "12", *p.ReauthPeriod)
		assert.Nil(t, p.UninstallPassword)
	})

	t.Run("Full payload JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 1,
			"policies": [{
				"deviceType": "DEVICE_TYPE_IOS",
				"id": 10,
				"name": "Default",
				"description": "Default Policy",
				"pac_url": "",
				"active": 1,
				"ruleOrder": 1,
				"logMode": 3,
				"logLevel": 0,
				"logFileSize": 100,
				"reactivateWebSecurityMinutes": "0",
				"highlightActiveControl": 0,
				"sendDisableServiceReason": 0,
				"refreshKerberosToken": 0,
				"enableDeviceGroups": 0,
				"notificationTemplateId": 733,
				"tunnelZappTraffic": 0,
				"groupAll": 1,
				"passcode": "",
				"logout_password": "",
				"disable_password": "",
				"showVPNTunNotification": 0,
				"useTunnelSDK4_3": 0,
				"ipv6Mode": 4,
				"policyExtension": {
					"sourcePortBasedBypasses": "3389:*",
					"packetTunnelExcludeList": "10.0.0.0/8",
					"packetTunnelIncludeList": "0.0.0.0/0",
					"useV8JsEngine": "1",
					"followRoutingTable": "1",
					"advanceZpaReauth": false,
					"enableSetProxyOnVPNAdapters": 1,
					"browserAuthType": "-1",
					"generateCliPasswordContract": {
						"enableCli": false,
						"allowZpaDisableWithoutPassword": false,
						"allowZiaDisableWithoutPassword": false,
						"allowZdxDisableWithoutPassword": false
					},
					"locationRulesetPolicies": {
						"offTrusted": {"id": 0, "name": ""},
						"trusted": {"id": 0, "name": ""},
						"vpnTrusted": {"id": 0, "name": ""},
						"splitVpnTrusted": {"id": 0, "name": ""}
					}
				},
				"disasterRecovery": {
					"enableZiaDR": false,
					"enableZpaDR": false,
					"ziaDRMethod": 0,
					"useZiaGlobalDb": false,
					"allowZiaTest": false,
					"allowZpaTest": false
				}
			}]
		}`

		var response application_profiles.ApplicationProfilesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.TotalCount)
		require.Len(t, response.Policies, 1)

		p := response.Policies[0]
		assert.Equal(t, 10, p.ID)
		assert.Equal(t, "Default", p.Name)
		assert.Equal(t, "DEVICE_TYPE_IOS", p.DeviceType)
		assert.Equal(t, 1, p.Active)
		assert.Equal(t, 733, p.NotificationTemplateId)
		assert.Equal(t, 4, p.Ipv6Mode)
		assert.Equal(t, "3389:*", p.PolicyExtension.SourcePortBasedBypasses)
		assert.Equal(t, "1", p.PolicyExtension.UseV8JsEngine)
		assert.False(t, p.PolicyExtension.AdvanceZpaReauth)
		assert.Equal(t, "-1", p.PolicyExtension.BrowserAuthType)
		assert.False(t, p.PolicyExtension.GenerateCliPasswordContract.EnableCli)
		assert.False(t, p.DisasterRecovery.EnableZiaDR)
		assert.False(t, p.DisasterRecovery.AllowZpaTest)
	})
}
