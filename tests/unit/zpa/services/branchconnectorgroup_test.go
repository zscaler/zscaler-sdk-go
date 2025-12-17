// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector_group"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestBranchConnectorGroup_GetSummary_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/branchConnectorGroup/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "bcg-001", Name: "Group 1"}, {ID: "bcg-002", Name: "Group 2"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := branch_connector_group.GetBranchConnectorGroupSummary(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
