package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/filetypecontrol/custom_file_types"
)

const customFileTypesPath = "/zia/api/v1/customFileTypes"
const customFileTypeCountPath = "/zia/api/v1/customFileTypes/count"

func TestCustomFileTypes_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	fileID := 100
	server.On("GET", customFileTypesPath+"/100", common.SuccessResponse(custom_file_types.CustomFileTypes{
		ID:         fileID,
		Name:       "Custom Archive",
		Extension:  "xyz",
		FileTypeID: 500,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := custom_file_types.Get(context.Background(), service, fileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, fileID, result.ID)
	assert.Equal(t, "xyz", result.Extension)
}

func TestCustomFileTypes_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	fileName := "Custom Archive"
	server.On("GET", customFileTypesPath, common.SuccessResponse([]custom_file_types.CustomFileTypes{
		{ID: 1, Name: "Other Type", Extension: "abc"},
		{ID: 100, Name: fileName, Extension: "xyz"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := custom_file_types.GetByName(context.Background(), service, fileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, fileName, result.Name)
}

func TestCustomFileTypes_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", customFileTypesPath, common.SuccessResponse(custom_file_types.CustomFileTypes{
		ID:        99999,
		Name:      "Custom Archive",
		Extension: "xyz",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newType := &custom_file_types.CustomFileTypes{
		Name:      "Custom Archive",
		Extension: "xyz",
	}

	result, err := custom_file_types.Create(context.Background(), service, newType)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestCustomFileTypes_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	fileID := 100
	server.On("PUT", customFileTypesPath+"/100", common.SuccessResponse(custom_file_types.CustomFileTypes{
		ID:        fileID,
		Name:      "Updated Custom Archive",
		Extension: "xyz",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateType := &custom_file_types.CustomFileTypes{
		ID:        fileID,
		Name:      "Updated Custom Archive",
		Extension: "xyz",
	}

	result, err := custom_file_types.Update(context.Background(), service, fileID, updateType)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Custom Archive", result.Name)
}

func TestCustomFileTypes_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", customFileTypesPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = custom_file_types.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestCustomFileTypes_GetCustomFileTypes_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", customFileTypesPath, common.SuccessResponse([]custom_file_types.CustomFileTypes{
		{ID: 100, Name: "Custom Archive", Extension: "xyz"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := custom_file_types.GetCustomFileTypes(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCustomFileTypes_GetCustomFileTypeCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", customFileTypeCountPath, common.SuccessResponse(5))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := custom_file_types.GetCustomFileTypeCount(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestCustomFileTypes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		ft := custom_file_types.CustomFileTypes{
			ID:        100,
			Name:      "Custom Archive",
			Extension: "xyz",
		}

		data, err := json.Marshal(ft)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"extension":"xyz"`)
	})
}
