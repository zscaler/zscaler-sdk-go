// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func TestPrivateCloudGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "pcg-12345"
	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/privateCloudControllerGroup/{id}
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudControllerGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(private_cloud_group.PrivateCloudGroup{
		ID:   groupID,
		Name: "Test Private Cloud Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud_group.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
}

func TestPrivateCloudGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/privateCloudControllerGroup
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudControllerGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []private_cloud_group.PrivateCloudGroup{{ID: "pcg-001"}, {ID: "pcg-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud_group.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPrivateCloudGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudControllerGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []private_cloud_group.PrivateCloudGroup{
			{ID: "pcg-001", Name: "Other Group"},
			{ID: "pcg-002", Name: groupName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud_group.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pcg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

func TestPrivateCloudGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudControllerGroup"

	server.On("POST", path, common.SuccessResponse(private_cloud_group.PrivateCloudGroup{
		ID:   "new-pcg-123",
		Name: "New Private Cloud Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newGroup := private_cloud_group.PrivateCloudGroup{
		Name: "New Private Cloud Group",
	}

	result, _, err := private_cloud_group.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-pcg-123", result.ID)
}

func TestPrivateCloudGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "pcg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudControllerGroup/" + groupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := private_cloud_group.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
