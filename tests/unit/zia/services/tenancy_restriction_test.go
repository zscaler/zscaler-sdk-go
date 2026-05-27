// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/tenancy_restriction"
)

func TestTenancyRestriction_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 12345
	path := "/zia/api/v1/tenancyRestrictionProfile/12345"
	server.On("GET", path, common.SuccessResponse(tenancy_restriction.TenancyRestrictionProfile{
		ID: profileID, Name: "O365 Restriction", AppType: "O365", ItemTypePrimary: "DOMAIN",
		ItemDataPrimary: []string{"company.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := tenancy_restriction.Get(context.Background(), service, profileID)
	require.NoError(t, err)
	assert.Equal(t, "O365", result.AppType)
}

func TestTenancyRestriction_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/tenancyRestrictionProfile"
	server.On("GET", path, common.SuccessResponse([]tenancy_restriction.TenancyRestrictionProfile{
		{ID: 1, Name: "Profile 1", AppType: "O365"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := tenancy_restriction.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestTenancyRestriction_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/tenancyRestrictionProfile"
	server.On("POST", path, common.SuccessResponse(tenancy_restriction.TenancyRestrictionProfile{
		ID: 99999, Name: "New Profile", AppType: "GOOGLE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newProfile := &tenancy_restriction.TenancyRestrictionProfile{
		Name: "New Profile", AppType: "GOOGLE", ItemTypePrimary: "DOMAIN",
		ItemDataPrimary: []string{"corp.example.com"},
	}

	result, _, err := tenancy_restriction.Create(context.Background(), service, newProfile)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestTenancyRestriction_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileID := 12345
	path := "/zia/api/v1/tenancyRestrictionProfile/12345"
	server.On("PUT", path, common.SuccessResponse(tenancy_restriction.TenancyRestrictionProfile{
		ID: profileID, Name: "Updated Profile",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &tenancy_restriction.TenancyRestrictionProfile{ID: profileID, Name: "Updated Profile"}
	result, _, err := tenancy_restriction.Update(context.Background(), service, profileID, update)
	require.NoError(t, err)
	assert.Equal(t, "Updated Profile", result.Name)
}

func TestTenancyRestriction_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/tenancyRestrictionProfile/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = tenancy_restriction.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestTenancyRestriction_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "O365 Restriction"
	path := "/zia/api/v1/tenancyRestrictionProfile"
	server.On("GET", path, common.SuccessResponse([]tenancy_restriction.TenancyRestrictionProfile{
		{ID: 1, Name: name, AppType: "O365"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := tenancy_restriction.GetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestTenancyRestriction_GetAppItemCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/tenancyRestrictionProfile/app-item-count/O365/DOMAIN"
	server.On("GET", path, common.SuccessResponse(map[string]int{
		"company.com": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := tenancy_restriction.GetAppItemCount(context.Background(), service, "O365", "DOMAIN")
	require.NoError(t, err)
	assert.Equal(t, 1, result["company.com"])
}

func TestTenancyRestriction_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/tenancyRestrictionProfile/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := tenancy_restriction.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}
