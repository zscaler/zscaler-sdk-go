// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/workload_tag_group"
)

func TestWorkloadTagGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/workloadTagGroup/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "tag-001"}, {ID: "tag-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := workload_tag_group.GetWorkloadTagGroup(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestWorkloadTagGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/workloadTagGroup/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "tag-001", Name: "Test Tag Group"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := workload_tag_group.GetByName(context.Background(), service, "Test Tag Group")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tag-001", result.ID)
}
