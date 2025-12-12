// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
)

func TestCloudConnectorGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "ccg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/cloudConnectorGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(cloud_connector_group.CloudConnectorGroup{
		ID:   groupID,
		Name: "Test Cloud Connector Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cloud_connector_group.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
}

func TestCloudConnectorGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/cloudConnectorGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []cloud_connector_group.CloudConnectorGroup{{ID: "ccg-001"}, {ID: "ccg-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cloud_connector_group.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
