// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customer_dr_tool"
)

func TestCustomerDRTool_GetCustomerDRTool_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/customerDRToolVersion"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []customer_dr_tool.CustomerDrTool{
			{
				ID:       "dr-001",
				Name:     "DR Tool 1",
				Platform: "WINDOWS",
				Version:  "1.0.0",
				Latest:   true,
			},
			{
				ID:       "dr-002",
				Name:     "DR Tool 2",
				Platform: "MACOS",
				Version:  "2.0.0",
				Latest:   false,
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := customer_dr_tool.GetCustomerDRTool(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "dr-001", result[0].ID)
	assert.Equal(t, "WINDOWS", result[0].Platform)
}
