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

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/v2/ipRanges
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges"

	// GetAll returns []*IPRanges (slice of pointers) and expects a raw array (not paginated)
	server.On("GET", path, common.SuccessResponse([]*c2c_ip_ranges.IPRanges{
		{ID: "range-001", Name: "Range 1"},
		{ID: "range-002", Name: "Range 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := c2c_ip_ranges.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
