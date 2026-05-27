package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_exact_data_match_lite"
)

const edmSchemaLitePath = "/zia/api/v1/dlpExactDataMatchSchemas/lite"

func TestDLPExactDataMatchLite_GetBySchemaName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	schemaName := "BD_EDM_TEMPLATE01"
	server.On("GET", edmSchemaLitePath, common.SuccessResponse([]dlp_exact_data_match_lite.DLPEDMLite{
		{
			Schema: dlp_exact_data_match_lite.SchemaIDNameExtension{
				ID:   1,
				Name: schemaName,
			},
			TokenList: []dlp_exact_data_match_lite.TokenList{
				{Name: "SSN", Type: "STRING", PrimaryKey: true, OriginalColumn: 1},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match_lite.GetBySchemaName(context.Background(), service, schemaName, true, true)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, schemaName, result[0].Schema.Name)
	assert.True(t, result[0].TokenList[0].PrimaryKey)
}

func TestDLPExactDataMatchLite_GetAllEDMSchema_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", edmSchemaLitePath, common.SuccessResponse([]dlp_exact_data_match_lite.DLPEDMLite{
		{
			Schema: dlp_exact_data_match_lite.SchemaIDNameExtension{
				ID:   1,
				Name: "BD_EDM_TEMPLATE01",
			},
		},
		{
			Schema: dlp_exact_data_match_lite.SchemaIDNameExtension{
				ID:   2,
				Name: "BD_EDM_TEMPLATE02",
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_exact_data_match_lite.GetAllEDMSchema(context.Background(), service, true, false)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPExactDataMatchLite_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		edm := dlp_exact_data_match_lite.DLPEDMLite{
			Schema: dlp_exact_data_match_lite.SchemaIDNameExtension{
				ID:   1,
				Name: "BD_EDM_TEMPLATE01",
			},
			TokenList: []dlp_exact_data_match_lite.TokenList{
				{Name: "SSN", Type: "STRING", PrimaryKey: true},
			},
		}

		data, err := json.Marshal(edm)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"BD_EDM_TEMPLATE01"`)
		assert.Contains(t, string(data), `"primaryKey":true`)
	})
}
