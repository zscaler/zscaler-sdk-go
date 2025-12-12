// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector"
)

func TestBranchConnector_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/branchConnector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []branch_connector.BranchConnector{{ID: "bc-001"}, {ID: "bc-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := branch_connector.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestBranchConnector_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/branchConnector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []branch_connector.BranchConnector{{ID: "bc-001", Name: "Test BC"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := branch_connector.GetByName(context.Background(), service, "Test BC")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "bc-001", result.ID)
}
