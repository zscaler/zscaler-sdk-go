package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_exact_data_match"
)

const edmSchemaPath = "/zia/api/v1/dlpExactDataMatchSchemas"

func TestDLPExactDataMatch_GetDLPEDMSchemaID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	schemaID := 1
	server.On("GET", edmSchemaPath+"/1", common.SuccessResponse(dlp_exact_data_match.DLPEDMSchema{
		SchemaID:     schemaID,
		ProjectName:  "BD_EDM_TEMPLATE01",
		SchemaActive: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match.GetDLPEDMSchemaID(context.Background(), service, schemaID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, schemaID, result.SchemaID)
	assert.Equal(t, "BD_EDM_TEMPLATE01", result.ProjectName)
	assert.True(t, result.SchemaActive)
}

func TestDLPExactDataMatch_GetDLPEDMByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	schemaName := "BD_EDM_TEMPLATE01"
	server.On("GET", edmSchemaPath, common.SuccessResponse([]dlp_exact_data_match.DLPEDMSchema{
		{SchemaID: 2, ProjectName: "Other Template"},
		{SchemaID: 1, ProjectName: schemaName, SchemaActive: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match.GetDLPEDMByName(context.Background(), service, schemaName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, schemaName, result.ProjectName)
}

func TestDLPExactDataMatch_GetDLPEDMByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", edmSchemaPath, common.SuccessResponse([]dlp_exact_data_match.DLPEDMSchema{
		{SchemaID: 1, ProjectName: "BD_EDM_TEMPLATE01", SchemaActive: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match.GetDLPEDMByName(context.Background(), service, "ThisEdmDoesNotExist")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPExactDataMatch_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", edmSchemaPath, common.SuccessResponse([]dlp_exact_data_match.DLPEDMSchema{
		{SchemaID: 1, ProjectName: "BD_EDM_TEMPLATE01", SchemaActive: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDLPExactDataMatch_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		schema := dlp_exact_data_match.DLPEDMSchema{
			SchemaID:     1,
			ProjectName:  "BD_EDM_TEMPLATE01",
			SchemaActive: true,
			Revision:     1,
		}

		data, err := json.Marshal(schema)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"projectName":"BD_EDM_TEMPLATE01"`)
		assert.Contains(t, string(data), `"schemaActive":true`)
	})
}
