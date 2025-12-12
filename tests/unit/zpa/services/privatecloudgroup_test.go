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
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/zpccGroup/" + groupID

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

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/zpccGroup"

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
