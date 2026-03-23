package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/custom_ip_apps"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestCustomIPApps_GetList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPAppsResponse{
		TotalCount: 2,
		CustomAppContracts: []custom_ip_apps.CustomIPApp{
			{ID: 6089, AppName: "App01", Active: false, CreatedBy: "140371"},
			{ID: 6090, AppName: "App02", Active: true, CreatedBy: "140371"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.CustomAppContracts, 2)
	assert.Equal(t, "App01", result.CustomAppContracts[0].AppName)
}

func TestCustomIPApps_GetList_WithSearch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPAppsResponse{
		TotalCount: 1,
		CustomAppContracts: []custom_ip_apps.CustomIPApp{
			{ID: 6089, AppName: "App01", Active: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := custom_ip_apps.GetCustomIPApps(context.Background(), service, "App01", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.TotalCount)
}

func TestCustomIPApps_GetByAppID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps/6089"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPApp{
		ID:      6089,
		AppName: "App01",
		Active:  false,
		AppDataBlob: []custom_ip_apps.AppDataBlob{
			{Proto: "TCP", Port: "8080", Ipaddr: "10.10.10.1"},
		},
		AppDataBlobV6: []custom_ip_apps.AppDataBlob{
			{Proto: "", Port: "", Ipaddr: ""},
		},
		CreatedBy:      "140371",
		EditedBy:       "140371",
		ZappDataBlob:   "10.10.10.1",
		ZappDataBlobV6: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := custom_ip_apps.GetByAppID(context.Background(), service, "6089")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 6089, result.ID)
	assert.Equal(t, "App01", result.AppName)
	assert.False(t, result.Active)
	require.Len(t, result.AppDataBlob, 1)
	assert.Equal(t, "TCP", result.AppDataBlob[0].Proto)
	assert.Equal(t, "8080", result.AppDataBlob[0].Port)
	assert.Equal(t, "10.10.10.1", result.AppDataBlob[0].Ipaddr)
}

func TestCustomIPApps_GetByAppID_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = custom_ip_apps.GetByAppID(context.Background(), service, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "appId is required")
}

func TestCustomIPApps_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPAppsResponse{
		TotalCount: 2,
		CustomAppContracts: []custom_ip_apps.CustomIPApp{
			{ID: 6089, AppName: "App01", Active: false},
			{ID: 6090, AppName: "App02", Active: true},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := custom_ip_apps.GetByName(context.Background(), service, "App02")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 6090, result.ID)
	assert.Equal(t, "App02", result.AppName)
}

func TestCustomIPApps_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPAppsResponse{
		TotalCount: 1,
		CustomAppContracts: []custom_ip_apps.CustomIPApp{
			{ID: 6089, AppName: "App01", Active: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := custom_ip_apps.GetByName(context.Background(), service, "app01")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 6089, result.ID)
}

func TestCustomIPApps_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/custom-ip-based-apps"

	server.On("GET", path, common.SuccessResponse(custom_ip_apps.CustomIPAppsResponse{
		TotalCount: 1,
		CustomAppContracts: []custom_ip_apps.CustomIPApp{
			{ID: 6089, AppName: "App01", Active: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = custom_ip_apps.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no custom IP-based app found with name: NonExistent")
}

// =====================================================
// Structure Tests
// =====================================================

func TestCustomIPApps_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CustomIPApp JSON marshaling", func(t *testing.T) {
		app := custom_ip_apps.CustomIPApp{
			ID:      6089,
			AppName: "App01",
			Active:  false,
			AppDataBlob: []custom_ip_apps.AppDataBlob{
				{Proto: "TCP", Port: "8080", Ipaddr: "10.10.10.1"},
			},
			CreatedBy:    "140371",
			ZappDataBlob: "10.10.10.1",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":6089`)
		assert.Contains(t, string(data), `"appName":"App01"`)
		assert.Contains(t, string(data), `"active":false`)
		assert.Contains(t, string(data), `"proto":"TCP"`)
	})

	t.Run("CustomIPApp JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 6089,
			"appName": "App01",
			"active": false,
			"uid": null,
			"appDataBlob": [
				{"proto": "TCP", "port": "8080", "ipaddr": "10.10.10.1", "fqdn": null}
			],
			"appDataBlobV6": [
				{"proto": "", "port": "", "ipaddr": "", "fqdn": null}
			],
			"createdBy": "140371",
			"editedBy": "140371",
			"editedTimestamp": "1773897128",
			"zappDataBlob": "10.10.10.1",
			"zappDataBlobV6": ""
		}`

		var app custom_ip_apps.CustomIPApp
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, 6089, app.ID)
		assert.Equal(t, "App01", app.AppName)
		assert.False(t, app.Active)
		require.Len(t, app.AppDataBlob, 1)
		assert.Equal(t, "TCP", app.AppDataBlob[0].Proto)
		assert.Equal(t, "8080", app.AppDataBlob[0].Port)
		assert.Equal(t, "10.10.10.1", app.AppDataBlob[0].Ipaddr)
		assert.Equal(t, "140371", app.CreatedBy)
		assert.Equal(t, "10.10.10.1", app.ZappDataBlob)
	})

	t.Run("CustomIPAppsResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 1,
			"customAppContracts": [
				{
					"id": 6089,
					"appName": "App01",
					"active": false,
					"appDataBlob": [{"proto": "TCP", "port": "8080", "ipaddr": "10.10.10.1"}],
					"createdBy": "140371"
				}
			]
		}`

		var response custom_ip_apps.CustomIPAppsResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 1, response.TotalCount)
		require.Len(t, response.CustomAppContracts, 1)
		assert.Equal(t, "App01", response.CustomAppContracts[0].AppName)
	})
}
