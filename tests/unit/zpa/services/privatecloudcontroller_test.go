// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
)

func TestPrivateCloudController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	cloudID := "pc-12345"
	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/privateCloudController/{id}
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudController/" + cloudID

	server.On("GET", path, common.SuccessResponse(private_cloud_controller.PrivateCloudController{
		ID:   cloudID,
		Name: "Test Private Cloud",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud_controller.Get(context.Background(), service, cloudID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, cloudID, result.ID)
}

func TestPrivateCloudController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/privateCloudController
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/privateCloudController"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []private_cloud_controller.PrivateCloudController{{ID: "pc-001"}, {ID: "pc-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := private_cloud_controller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
