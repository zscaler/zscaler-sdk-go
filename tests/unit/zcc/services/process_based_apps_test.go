package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/process_based_apps"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestProcessBasedApps_GetList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/process-based-apps"

	server.On("GET", path, common.SuccessResponse(process_based_apps.ProcessBasedAppsResponse{
		TotalCount: 2,
		AppIdentities: []process_based_apps.ProcessBasedApp{
			{ID: 7119, AppName: "App01", MatchingCriteria: 2},
			{ID: 7120, AppName: "App02", MatchingCriteria: 1},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := process_based_apps.GetProcessBasedApps(context.Background(), service, "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.AppIdentities, 2)
	assert.Equal(t, "App01", result.AppIdentities[0].AppName)
}

func TestProcessBasedApps_GetByAppID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/process-based-apps/7119"

	server.On("GET", path, common.SuccessResponse(process_based_apps.ProcessBasedApp{
		ID:                 7119,
		AppName:            "App01",
		FileNames:          []string{""},
		FilePaths:          []string{`*\Program\Apps\MSApps\ms-teams.exe`},
		MatchingCriteria:   2,
		SignaturePayload:   `{"signatures":[]}`,
		CertificatePayload: `{"certificates":[{"certName":"cert01","thumbprint":"304B4346F337146144839209C6C79146D05E4764"}]}`,
		CreatedBy:          "140371",
		EditedBy:           "140371",
		EditedTimestamp:    "1773897787",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := process_based_apps.GetByAppID(context.Background(), service, "7119")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 7119, result.ID)
	assert.Equal(t, "App01", result.AppName)
	assert.Equal(t, 2, result.MatchingCriteria)
	require.Len(t, result.FilePaths, 1)
	assert.Contains(t, result.FilePaths[0], "ms-teams.exe")
	assert.Contains(t, result.CertificatePayload, "cert01")
	assert.Equal(t, "140371", result.CreatedBy)
}

func TestProcessBasedApps_GetByAppID_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = process_based_apps.GetByAppID(context.Background(), service, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "appId is required")
}

func TestProcessBasedApps_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/process-based-apps"

	server.On("GET", path, common.SuccessResponse(process_based_apps.ProcessBasedAppsResponse{
		TotalCount: 2,
		AppIdentities: []process_based_apps.ProcessBasedApp{
			{ID: 7119, AppName: "App01", MatchingCriteria: 2},
			{ID: 7120, AppName: "App02", MatchingCriteria: 1},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := process_based_apps.GetByName(context.Background(), service, "App02")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 7120, result.ID)
	assert.Equal(t, "App02", result.AppName)
}

func TestProcessBasedApps_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/process-based-apps"

	server.On("GET", path, common.SuccessResponse(process_based_apps.ProcessBasedAppsResponse{
		TotalCount: 1,
		AppIdentities: []process_based_apps.ProcessBasedApp{
			{ID: 7119, AppName: "App01", MatchingCriteria: 2},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := process_based_apps.GetByName(context.Background(), service, "app01")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 7119, result.ID)
}

func TestProcessBasedApps_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/process-based-apps"

	server.On("GET", path, common.SuccessResponse(process_based_apps.ProcessBasedAppsResponse{
		TotalCount: 1,
		AppIdentities: []process_based_apps.ProcessBasedApp{
			{ID: 7119, AppName: "App01", MatchingCriteria: 2},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = process_based_apps.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no process-based app found with name: NonExistent")
}

// =====================================================
// Structure Tests
// =====================================================

func TestProcessBasedApps_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ProcessBasedApp JSON marshaling", func(t *testing.T) {
		app := process_based_apps.ProcessBasedApp{
			ID:               7119,
			AppName:          "App01",
			FileNames:        []string{""},
			FilePaths:        []string{`*\Program\Apps\MSApps\ms-teams.exe`},
			MatchingCriteria: 2,
			CreatedBy:        "140371",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":7119`)
		assert.Contains(t, string(data), `"appName":"App01"`)
		assert.Contains(t, string(data), `"matchingCriteria":2`)
		assert.Contains(t, string(data), `"filePaths"`)
	})

	t.Run("ProcessBasedApp JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 7119,
			"appName": "App01",
			"fileNames": [""],
			"filePaths": ["*\\Program\\Apps\\MSApps\\ms-teams.exe"],
			"matchingCriteria": 2,
			"signaturePayload": "{\"signatures\":[]}",
			"certificatePayload": "{\"certificates\":[{\"certName\":\"cert01\",\"thumbprint\":\"304B4346F337146144839209C6C79146D05E4764\"}]}",
			"createdBy": "140371",
			"editedBy": "140371",
			"editedTimestamp": "1773897787"
		}`

		var app process_based_apps.ProcessBasedApp
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, 7119, app.ID)
		assert.Equal(t, "App01", app.AppName)
		assert.Equal(t, 2, app.MatchingCriteria)
		require.Len(t, app.FilePaths, 1)
		assert.Contains(t, app.FilePaths[0], "ms-teams.exe")
		assert.Contains(t, app.SignaturePayload, "signatures")
		assert.Contains(t, app.CertificatePayload, "cert01")
		assert.Equal(t, "140371", app.CreatedBy)
		assert.Equal(t, "1773897787", app.EditedTimestamp)
	})

	t.Run("ProcessBasedAppsResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 2,
			"appIdentities": [
				{"id": 7119, "appName": "App01", "matchingCriteria": 2},
				{"id": 7120, "appName": "App02", "matchingCriteria": 1}
			]
		}`

		var response process_based_apps.ProcessBasedAppsResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.TotalCount)
		require.Len(t, response.AppIdentities, 2)
		assert.Equal(t, "App01", response.AppIdentities[0].AppName)
		assert.Equal(t, 2, response.AppIdentities[0].MatchingCriteria)
		assert.Equal(t, "App02", response.AppIdentities[1].AppName)
	})
}
