package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_idm_profiles"
)

const idmProfilePath = "/zia/api/v1/idmprofile"

func TestDLPIDMProfiles_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 1
	server.On("GET", idmProfilePath+"/1", common.SuccessResponse(dlp_idm_profiles.DLPIDMProfile{
		ProfileID:   profileID,
		ProfileName: "BD_IDM_TEMPLATE01",
		ProfileDesc: "IDM template for testing",
		ProfileType: "LOCAL",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profiles.Get(context.Background(), service, profileID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileID, result.ProfileID)
	assert.Equal(t, "BD_IDM_TEMPLATE01", result.ProfileName)
}

func TestDLPIDMProfiles_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "BD_IDM_TEMPLATE01"
	server.On("GET", idmProfilePath, common.SuccessResponse([]dlp_idm_profiles.DLPIDMProfile{
		{ProfileID: 2, ProfileName: "Other Template"},
		{ProfileID: 1, ProfileName: profileName, ProfileType: "LOCAL"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profiles.GetByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, profileName, result.ProfileName)
}

func TestDLPIDMProfiles_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", idmProfilePath, common.SuccessResponse([]dlp_idm_profiles.DLPIDMProfile{
		{ProfileID: 1, ProfileName: "BD_IDM_TEMPLATE01"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profiles.GetByName(context.Background(), service, "ThisIDMTemplateDoesNotExist")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPIDMProfiles_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", idmProfilePath, common.SuccessResponse([]dlp_idm_profiles.DLPIDMProfile{
		{ProfileID: 1, ProfileName: "BD_IDM_TEMPLATE01", ProfileType: "LOCAL"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_idm_profiles.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDLPIDMProfiles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		profile := dlp_idm_profiles.DLPIDMProfile{
			ProfileID:      1,
			ProfileName:    "BD_IDM_TEMPLATE01",
			ProfileDesc:    "IDM template for testing",
			ProfileType:    "LOCAL",
			ScheduleType:   "DAILY",
			ScheduleTime:   180,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"profileName":"BD_IDM_TEMPLATE01"`)
		assert.Contains(t, string(data), `"profileType":"LOCAL"`)
	})
}
