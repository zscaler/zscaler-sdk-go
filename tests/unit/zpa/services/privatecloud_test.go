// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud"
)

// privateCloudPath builds the base collection path for the private_cloud service.
func privateCloudPath() string {
	return "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloud"
}

func TestPrivateCloud_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	path := privateCloudPath() + "/" + cloudID

	server.On("GET", path, common.SuccessResponse(private_cloud.PrivateCloudController{
		ID:          cloudID,
		Name:        "Test Private Cloud",
		Description: "desc",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.Get(context.Background(), service, cloudID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, cloudID, result.ID)
	assert.Equal(t, "Test Private Cloud", result.Name)
	assert.True(t, result.Enabled)
}

func TestPrivateCloud_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-does-not-exist"
	path := privateCloudPath() + "/" + cloudID

	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.Get(context.Background(), service, cloudID)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPrivateCloud_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", privateCloudPath(), common.SuccessResponse(map[string]interface{}{
		"list": []private_cloud.PrivateCloudController{
			{ID: "pc-001", Name: "Cloud One"},
			{ID: "pc-002", Name: "Cloud Two"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "pc-001", result[0].ID)
}

func TestPrivateCloud_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", privateCloudPath(), common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.GetAll(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPrivateCloud_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudName := "Production Cloud"

	server.On("GET", privateCloudPath(), common.SuccessResponse(map[string]interface{}{
		"list": []private_cloud.PrivateCloudController{
			{ID: "pc-001", Name: "Other Cloud"},
			{ID: "pc-002", Name: cloudName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.GetByName(context.Background(), service, cloudName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pc-002", result.ID)
	assert.Equal(t, cloudName, result.Name)
}

func TestPrivateCloud_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", privateCloudPath(), common.SuccessResponse(map[string]interface{}{
		"list": []private_cloud.PrivateCloudController{
			{ID: "pc-001", Name: "Production Cloud"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.GetByName(context.Background(), service, "production cloud")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pc-001", result.ID)
}

func TestPrivateCloud_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", privateCloudPath(), common.SuccessResponse(map[string]interface{}{
		"list": []private_cloud.PrivateCloudController{
			{ID: "pc-001", Name: "Existing Cloud"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPrivateCloud_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", privateCloudPath(), common.SuccessResponse(private_cloud.PrivateCloudController{
		ID:             "pc-new",
		Name:           "New Private Cloud",
		Description:    "New Private Cloud",
		Enabled:        true,
		ReEnrollPeriod: "90",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newCloud := private_cloud.PrivateCloudController{
		Name:           "New Private Cloud",
		Description:    "New Private Cloud",
		Enabled:        true,
		ReEnrollPeriod: "90",
	}

	result, _, err := private_cloud.Create(context.Background(), service, newCloud)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pc-new", result.ID)
	assert.Equal(t, "New Private Cloud", result.Name)
}

func TestPrivateCloud_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", privateCloudPath(), common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud.Create(context.Background(), service, private_cloud.PrivateCloudController{
		Name: "New Private Cloud",
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestPrivateCloud_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	path := privateCloudPath() + "/" + cloudID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateCloud := private_cloud.PrivateCloudController{
		ID:   cloudID,
		Name: "Updated Private Cloud",
	}

	resp, err := private_cloud.Update(context.Background(), service, cloudID, &updateCloud)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPrivateCloud_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	path := privateCloudPath() + "/" + cloudID

	server.On("PUT", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := private_cloud.Update(context.Background(), service, cloudID, &private_cloud.PrivateCloudController{
		ID: cloudID,
	})

	require.Error(t, err)
	_ = resp
}

func TestPrivateCloud_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	path := privateCloudPath() + "/" + cloudID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := private_cloud.Delete(context.Background(), service, cloudID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPrivateCloud_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	path := privateCloudPath() + "/" + cloudID

	server.On("DELETE", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	_, err = private_cloud.Delete(context.Background(), service, cloudID)

	require.Error(t, err)
}
