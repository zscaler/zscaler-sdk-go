// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func TestC2CIPRanges_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/c2cIPRanges"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []c2c_ip_ranges.IPRanges{{ID: "range-001"}, {ID: "range-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := c2c_ip_ranges.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
