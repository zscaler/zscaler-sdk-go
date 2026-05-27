package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_idm_profile_lite"
)

const idmProfileLitePath = "/zia/api/v1/idmprofile/lite"

func TestDLPIDMProfileLite_GetDLPProfileLiteID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 1
	server.On("GET", idmProfileLitePath, common.SuccessResponse([]dlp_idm_profile_lite.DLPIDMProfileLite{
		{ProfileID: 2, TemplateName: "Other Template"},
		{ProfileID: profileID, TemplateName: "BD_IDM_TEMPLATE01", NumDocuments: 100},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profile_lite.GetDLPProfileLiteID(context.Background(), service, profileID, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ProfileID)
	assert.Equal(t, "BD_IDM_TEMPLATE01", result.TemplateName)
}

func TestDLPIDMProfileLite_GetDLPProfileLiteByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	templateName := "BD_IDM_TEMPLATE01"
	server.On("GET", idmProfileLitePath, common.SuccessResponse([]dlp_idm_profile_lite.DLPIDMProfileLite{
		{ProfileID: 1, TemplateName: templateName, NumDocuments: 100},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profile_lite.GetDLPProfileLiteByName(context.Background(), service, templateName, true)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, templateName, result.TemplateName)
}

func TestDLPIDMProfileLite_GetDLPProfileLiteByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", idmProfileLitePath, common.SuccessResponse([]dlp_idm_profile_lite.DLPIDMProfileLite{
		{ProfileID: 1, TemplateName: "BD_IDM_TEMPLATE01"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profile_lite.GetDLPProfileLiteByName(context.Background(), service, "ThisIDMTemplateDoesNotExist", true)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPIDMProfileLite_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", idmProfileLitePath, common.SuccessResponse([]dlp_idm_profile_lite.DLPIDMProfileLite{
		{ProfileID: 1, TemplateName: "BD_IDM_TEMPLATE01", NumDocuments: 100},
		{ProfileID: 2, TemplateName: "BD_IDM_TEMPLATE02", NumDocuments: 50},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profile_lite.GetAll(context.Background(), service, true)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPIDMProfileLite_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		profile := dlp_idm_profile_lite.DLPIDMProfileLite{
			ProfileID:    1,
			TemplateName: "BD_IDM_TEMPLATE01",
			NumDocuments: 100,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"templateName":"BD_IDM_TEMPLATE01"`)
		assert.Contains(t, string(data), `"numDocuments":100`)
	})
}
