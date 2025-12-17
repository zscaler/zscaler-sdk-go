// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_engines"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDLPEngines_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	engineID := 12345
	path := "/zia/api/v1/dlpEngines/12345"

	server.On("GET", path, common.SuccessResponse(dlp_engines.DLPEngines{
		ID:              engineID,
		Name:            "Custom DLP Engine",
		Description:     "Custom engine for PII detection",
		EngineExpression: "((D63.S > 1))",
		CustomDlpEngine: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.Get(context.Background(), service, engineID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, engineID, result.ID)
	assert.Equal(t, "Custom DLP Engine", result.Name)
	assert.True(t, result.CustomDlpEngine)
}

func TestDLPEngines_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	engineName := "Custom DLP Engine"
	path := "/zia/api/v1/dlpEngines"

	server.On("GET", path, common.SuccessResponse([]dlp_engines.DLPEngines{
		{ID: 1, Name: "Other Engine", CustomDlpEngine: true},
		{ID: 2, Name: engineName, CustomDlpEngine: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.GetByName(context.Background(), service, engineName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, engineName, result.Name)
}

func TestDLPEngines_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dlpEngines"

	server.On("GET", path, common.SuccessResponse([]dlp_engines.DLPEngines{
		{ID: 1, Name: "Engine 1", CustomDlpEngine: true},
		{ID: 2, Name: "Engine 2", CustomDlpEngine: false},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPEngines_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dlpEngines"

	server.On("POST", path, common.SuccessResponse(dlp_engines.DLPEngines{
		ID:              99999,
		Name:            "New DLP Engine",
		CustomDlpEngine: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newEngine := &dlp_engines.DLPEngines{
		Name:            "New DLP Engine",
		CustomDlpEngine: true,
	}

	result, _, err := dlp_engines.Create(context.Background(), service, newEngine)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDLPEngines_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	engineID := 12345
	path := "/zia/api/v1/dlpEngines/12345"

	server.On("PUT", path, common.SuccessResponse(dlp_engines.DLPEngines{
		ID:   engineID,
		Name: "Updated DLP Engine",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateEngine := &dlp_engines.DLPEngines{
		ID:   engineID,
		Name: "Updated DLP Engine",
	}

	result, _, err := dlp_engines.Update(context.Background(), service, engineID, updateEngine)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated DLP Engine", result.Name)
}

func TestDLPEngines_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	engineID := 12345
	path := "/zia/api/v1/dlpEngines/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dlp_engines.Delete(context.Background(), service, engineID)

	require.NoError(t, err)
}

func TestDLPEngines_GetEngineLiteID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	engineID := 12345
	path := "/zia/api/v1/dlpEngines/lite"

	server.On("GET", path, common.SuccessResponse([]dlp_engines.DLPEngines{
		{ID: 1, Name: "Other Lite Engine"},
		{ID: engineID, Name: "Target Engine"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.GetEngineLiteID(context.Background(), service, engineID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, engineID, result.ID)
}

func TestDLPEngines_GetByPredefinedEngine_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	predefinedEngineName := "HIPAA"
	path := "/zia/api/v1/dlpEngines/lite"

	server.On("GET", path, common.SuccessResponse([]dlp_engines.DLPEngines{
		{ID: 1, Name: "Other Engine", PredefinedEngineName: "OTHER"},
		{ID: 2, Name: "HIPAA Engine", PredefinedEngineName: predefinedEngineName},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.GetByPredefinedEngine(context.Background(), service, predefinedEngineName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, predefinedEngineName, result.PredefinedEngineName)
}

func TestDLPEngines_GetAllEngineLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dlpEngines/lite"

	server.On("GET", path, common.SuccessResponse([]dlp_engines.DLPEngines{
		{ID: 1, Name: "Lite Engine 1"},
		{ID: 2, Name: "Lite Engine 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_engines.GetAllEngineLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDLPEngines_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DLPEngines JSON marshaling", func(t *testing.T) {
		engine := dlp_engines.DLPEngines{
			ID:               12345,
			Name:             "Custom DLP Engine",
			Description:      "Custom engine for detection",
			EngineExpression: "((D63.S > 1))",
			CustomDlpEngine:  true,
		}

		data, err := json.Marshal(engine)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom DLP Engine"`)
		assert.Contains(t, string(data), `"customDlpEngine":true`)
	})

	t.Run("DLPEngines JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Predefined Engine",
			"description": "System engine",
			"engineExpression": "((D1.S > 0))",
			"customDlpEngine": false,
			"predefinedEngineName": "SSN_ENGINE"
		}`

		var engine dlp_engines.DLPEngines
		err := json.Unmarshal([]byte(jsonData), &engine)
		require.NoError(t, err)

		assert.Equal(t, 54321, engine.ID)
		assert.Equal(t, "Predefined Engine", engine.Name)
		assert.False(t, engine.CustomDlpEngine)
	})
}

