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

func TestC2CIPRanges_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	rangeID := "range-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges/" + rangeID

	server.On("GET", path, common.SuccessResponse(c2c_ip_ranges.IPRanges{
		ID:   rangeID,
		Name: "Test IP Range",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := c2c_ip_ranges.Get(context.Background(), service, rangeID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, rangeID, result.ID)
}

func TestC2CIPRanges_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	rangeName := "Production Range"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges"

	// GetByName uses GetAllPagesGenericWithCustomFilters which expects paginated response
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []c2c_ip_ranges.IPRanges{
			{ID: "range-001", Name: "Other Range"},
			{ID: "range-002", Name: rangeName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := c2c_ip_ranges.GetByName(context.Background(), service, rangeName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "range-002", result.ID)
	assert.Equal(t, rangeName, result.Name)
}

func TestC2CIPRanges_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges"

	server.On("POST", path, common.SuccessResponse(c2c_ip_ranges.IPRanges{
		ID:   "new-range-123",
		Name: "New IP Range",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newRange := c2c_ip_ranges.IPRanges{
		Name: "New IP Range",
	}

	result, _, err := c2c_ip_ranges.Create(context.Background(), service, &newRange)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-range-123", result.ID)
}

func TestC2CIPRanges_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	rangeID := "range-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges/" + rangeID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateRange := &c2c_ip_ranges.IPRanges{
		ID:   rangeID,
		Name: "Updated IP Range",
	}

	resp, err := c2c_ip_ranges.Update(context.Background(), service, rangeID, updateRange)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestC2CIPRanges_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	rangeID := "range-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ipRanges/" + rangeID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := c2c_ip_ranges.Delete(context.Background(), service, rangeID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
