package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/web_app_service"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestWebAppService_GetList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true, AppSvcId: 123456, UID: "abc-123"},
		{ID: 5, AppName: "Zoom", Active: true, AppSvcId: 789, UID: "def-456"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_app_service.GetWebAppServices(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Microsoft Teams", result[0].AppName)
	assert.Equal(t, 3, result[0].ID)
	assert.True(t, result[0].Active)
}

func TestWebAppService_GetList_WithSearch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_app_service.GetWebAppServices(context.Background(), service, "Microsoft", nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Microsoft Teams", result[0].AppName)
}

func TestWebAppService_GetByAppID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true, AppSvcId: 123456},
		{ID: 5, AppName: "Zoom", Active: true, AppSvcId: 789},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_app_service.GetByAppID(context.Background(), service, "5")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 5, result.ID)
	assert.Equal(t, "Zoom", result.AppName)
}

func TestWebAppService_GetByAppID_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = web_app_service.GetByAppID(context.Background(), service, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "appId is required")
}

func TestWebAppService_GetByAppID_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = web_app_service.GetByAppID(context.Background(), service, "999")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestWebAppService_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true},
		{ID: 5, AppName: "Zoom", Active: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_app_service.GetByName(context.Background(), service, "Zoom")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 5, result.ID)
	assert.Equal(t, "Zoom", result.AppName)
}

func TestWebAppService_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{
		{ID: 3, AppName: "Microsoft Teams", Active: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := web_app_service.GetByName(context.Background(), service, "microsoft teams")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Microsoft Teams", result.AppName)
}

func TestWebAppService_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webAppService/listByCompany"

	server.On("GET", path, common.SuccessResponse([]web_app_service.WebAppService{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = web_app_service.GetByName(context.Background(), service, "NonExistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// =====================================================
// JSON Structure Tests
// =====================================================

func TestWebAppService_JSONMarshalUnmarshal(t *testing.T) {
	original := web_app_service.WebAppService{
		ID:         3,
		AppVersion: 0,
		AppSvcId:   123456,
		AppName:    "Microsoft Teams",
		Active:     true,
		UID:        "33f85c91-598b-4db6-8f12-b74116ea3560",
		AppDataBlob: []web_app_service.AppDataBlob{
			{Proto: "UDP", Port: "3478,3479,3480,3481", Ipaddr: "52.112.0.0/14", Fqdn: ""},
			{Proto: "TCP", Port: "443,80", Ipaddr: "52.112.0.0/14", Fqdn: ""},
		},
		AppDataBlobV6: []web_app_service.AppDataBlob{
			{Proto: "UDP", Port: "3478,3479,3480,3481", Ipaddr: "2603:1063::/38"},
		},
		CreatedBy:      "2315",
		EditedBy:       "0",
		EditedTimestamp: "1733811797",
		ZappDataBlob:   "52.112.0.0/14",
		ZappDataBlobV6: "2603:1063::/38",
		Version:        0,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded web_app_service.WebAppService
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.AppVersion, decoded.AppVersion)
	assert.Equal(t, original.AppSvcId, decoded.AppSvcId)
	assert.Equal(t, original.AppName, decoded.AppName)
	assert.Equal(t, original.Active, decoded.Active)
	assert.Equal(t, original.UID, decoded.UID)
	assert.Len(t, decoded.AppDataBlob, 2)
	assert.Equal(t, "UDP", decoded.AppDataBlob[0].Proto)
	assert.Len(t, decoded.AppDataBlobV6, 1)
	assert.Equal(t, original.CreatedBy, decoded.CreatedBy)
	assert.Equal(t, original.EditedBy, decoded.EditedBy)
	assert.Equal(t, original.EditedTimestamp, decoded.EditedTimestamp)
	assert.Equal(t, original.ZappDataBlob, decoded.ZappDataBlob)
	assert.Equal(t, original.ZappDataBlobV6, decoded.ZappDataBlobV6)
	assert.Equal(t, original.Version, decoded.Version)
}

func TestWebAppService_FullPayloadUnmarshal(t *testing.T) {
	payload := `[{
		"id": 3,
		"appVersion": 0,
		"appSvcId": 123456,
		"appName": "Microsoft Teams",
		"active": true,
		"uid": "33f85c91-598b-4db6-8f12-b74116ea3560",
		"appDataBlob": [
			{"proto": "UDP", "port": "3478,3479", "ipaddr": "52.112.0.0/14", "fqdn": ""},
			{"proto": "TCP", "port": "443,80", "ipaddr": "52.112.0.0/14", "fqdn": ""}
		],
		"appDataBlobV6": [
			{"proto": "UDP", "port": "3478,3479", "ipaddr": "2603:1063::/38", "fqdn": ""}
		],
		"createdBy": "2315",
		"editedBy": "0",
		"editedTimestamp": "1733811797",
		"zappDataBlob": "52.112.0.0/14",
		"zappDataBlobV6": "2603:1063::/38",
		"version": 0
	}]`

	var apps []web_app_service.WebAppService
	err := json.Unmarshal([]byte(payload), &apps)
	require.NoError(t, err)
	require.Len(t, apps, 1)

	app := apps[0]
	assert.Equal(t, 3, app.ID)
	assert.Equal(t, 0, app.AppVersion)
	assert.Equal(t, 123456, app.AppSvcId)
	assert.Equal(t, "Microsoft Teams", app.AppName)
	assert.True(t, app.Active)
	assert.Equal(t, "33f85c91-598b-4db6-8f12-b74116ea3560", app.UID)
	assert.Len(t, app.AppDataBlob, 2)
	assert.Equal(t, "UDP", app.AppDataBlob[0].Proto)
	assert.Equal(t, "3478,3479", app.AppDataBlob[0].Port)
	assert.Len(t, app.AppDataBlobV6, 1)
	assert.Equal(t, "2315", app.CreatedBy)
	assert.Equal(t, "0", app.EditedBy)
	assert.Equal(t, "1733811797", app.EditedTimestamp)
	assert.Equal(t, "52.112.0.0/14", app.ZappDataBlob)
	assert.Equal(t, "2603:1063::/38", app.ZappDataBlobV6)
	assert.Equal(t, 0, app.Version)
}
