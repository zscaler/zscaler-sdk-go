package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/predefined_ip_apps"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestPredefinedIPApps_GetList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/predefined-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(predefined_ip_apps.PredefinedIPAppsResponse{
		TotalCount: 2,
		AppServiceContracts: []predefined_ip_apps.PredefinedIPApp{
			{ID: 3, AppName: "Microsoft Teams", Active: true, AppSvcId: 123456},
			{ID: 7, AppName: "SharePoint", Active: true, AppSvcId: 310},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := predefined_ip_apps.GetPredefinedIPApps(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.AppServiceContracts, 2)
	assert.Equal(t, "Microsoft Teams", result.AppServiceContracts[0].AppName)
}

func TestPredefinedIPApps_GetByAppID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/predefined-ip-based-apps/3"

	server.On("GET", path, common.SuccessResponse(predefined_ip_apps.PredefinedIPApp{
		ID:       3,
		AppSvcId: 123456,
		AppName:  "Microsoft Teams",
		Active:   true,
		UID:      "33f85c91-598b-4db6-8f12-b74116ea3560",
		AppDataBlob: []predefined_ip_apps.AppDataBlob{
			{Proto: "UDP", Port: "3478,3479,3480,3481", Ipaddr: "52.112.0.0/14"},
			{Proto: "TCP", Port: "443,80", Ipaddr: "52.112.0.0/14"},
		},
		AppDataBlobV6: []predefined_ip_apps.AppDataBlob{
			{Proto: "UDP", Port: "3478,3479,3480,3481", Ipaddr: "2603:1063::/38"},
		},
		CreatedBy:    "2315",
		ZappDataBlob: "52.112.0.0/14",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := predefined_ip_apps.GetByAppID(context.Background(), service, "3")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.ID)
	assert.Equal(t, "Microsoft Teams", result.AppName)
	assert.True(t, result.Active)
	assert.Equal(t, 123456, result.AppSvcId)
	require.Len(t, result.AppDataBlob, 2)
	assert.Equal(t, "UDP", result.AppDataBlob[0].Proto)
	assert.Equal(t, "TCP", result.AppDataBlob[1].Proto)
}

func TestPredefinedIPApps_GetByAppID_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = predefined_ip_apps.GetByAppID(context.Background(), service, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "appId is required")
}

func TestPredefinedIPApps_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/predefined-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(predefined_ip_apps.PredefinedIPAppsResponse{
		TotalCount: 3,
		AppServiceContracts: []predefined_ip_apps.PredefinedIPApp{
			{ID: 3, AppName: "Microsoft Teams", Active: true},
			{ID: 7, AppName: "SharePoint", Active: true},
			{ID: 9, AppName: "Microsoft Exchange", Active: true},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := predefined_ip_apps.GetByName(context.Background(), service, "SharePoint")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 7, result.ID)
	assert.Equal(t, "SharePoint", result.AppName)
}

func TestPredefinedIPApps_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/predefined-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(predefined_ip_apps.PredefinedIPAppsResponse{
		TotalCount: 1,
		AppServiceContracts: []predefined_ip_apps.PredefinedIPApp{
			{ID: 3, AppName: "Microsoft Teams", Active: true},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := predefined_ip_apps.GetByName(context.Background(), service, "microsoft teams")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 3, result.ID)
}

func TestPredefinedIPApps_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/predefined-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(predefined_ip_apps.PredefinedIPAppsResponse{
		TotalCount: 1,
		AppServiceContracts: []predefined_ip_apps.PredefinedIPApp{
			{ID: 3, AppName: "Microsoft Teams", Active: true},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = predefined_ip_apps.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no predefined IP-based app found with name: NonExistent")
}

// =====================================================
// Structure Tests
// =====================================================

func TestPredefinedIPApps_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PredefinedIPApp JSON marshaling", func(t *testing.T) {
		app := predefined_ip_apps.PredefinedIPApp{
			ID:       3,
			AppSvcId: 123456,
			AppName:  "Microsoft Teams",
			Active:   true,
			UID:      "33f85c91-598b-4db6-8f12-b74116ea3560",
			AppDataBlob: []predefined_ip_apps.AppDataBlob{
				{Proto: "UDP", Port: "3478,3479,3480,3481", Ipaddr: "52.112.0.0/14"},
			},
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":3`)
		assert.Contains(t, string(data), `"appSvcId":123456`)
		assert.Contains(t, string(data), `"appName":"Microsoft Teams"`)
		assert.Contains(t, string(data), `"active":true`)
		assert.Contains(t, string(data), `"proto":"UDP"`)
	})

	t.Run("PredefinedIPApp JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 7,
			"appVersion": 0,
			"appSvcId": 310,
			"appName": "SharePoint",
			"active": true,
			"uid": "CC84C4F3-DEB3-497B-ABBA-A9F36FF690C5",
			"appDataBlob": [
				{"proto": "TCP", "port": "443,80", "ipaddr": "13.107.136.0/22", "fqdn": ""}
			],
			"appDataBlobV6": [
				{"proto": "TCP", "port": "443,80", "ipaddr": "2603:1061:1300::/40", "fqdn": ""}
			],
			"createdBy": "0",
			"editedBy": "0",
			"editedTimestamp": "1729263594",
			"zappDataBlob": "13.107.136.0/22",
			"zappDataBlobV6": "2603:1061:1300::/40"
		}`

		var app predefined_ip_apps.PredefinedIPApp
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, 7, app.ID)
		assert.Equal(t, 310, app.AppSvcId)
		assert.Equal(t, "SharePoint", app.AppName)
		assert.True(t, app.Active)
		assert.Equal(t, "CC84C4F3-DEB3-497B-ABBA-A9F36FF690C5", app.UID)
		require.Len(t, app.AppDataBlob, 1)
		assert.Equal(t, "TCP", app.AppDataBlob[0].Proto)
		assert.Equal(t, "443,80", app.AppDataBlob[0].Port)
		require.Len(t, app.AppDataBlobV6, 1)
		assert.Equal(t, "13.107.136.0/22", app.ZappDataBlob)
		assert.Equal(t, "2603:1061:1300::/40", app.ZappDataBlobV6)
	})

	t.Run("PredefinedIPAppsResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 2,
			"appServiceContracts": [
				{"id": 3, "appName": "Microsoft Teams", "active": true, "appSvcId": 123456},
				{"id": 7, "appName": "SharePoint", "active": true, "appSvcId": 310}
			]
		}`

		var response predefined_ip_apps.PredefinedIPAppsResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.TotalCount)
		require.Len(t, response.AppServiceContracts, 2)
		assert.Equal(t, "Microsoft Teams", response.AppServiceContracts[0].AppName)
		assert.Equal(t, 123456, response.AppServiceContracts[0].AppSvcId)
		assert.Equal(t, "SharePoint", response.AppServiceContracts[1].AppName)
	})
}
