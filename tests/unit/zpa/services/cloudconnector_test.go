// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector"
)

func TestCloudConnector_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/cloudConnector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []cloud_connector.CloudConnector{{ID: "cc-001"}, {ID: "cc-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cloud_connector.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCloudConnector_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/cloudConnector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []cloud_connector.CloudConnector{{ID: "cc-001", Name: "Test CC"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cloud_connector.GetByName(context.Background(), service, "Test CC")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "cc-001", result.ID)
}
