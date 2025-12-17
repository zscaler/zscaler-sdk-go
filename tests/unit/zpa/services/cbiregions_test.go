// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
)

func TestCBIRegions_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/regions"

	server.On("GET", path, common.SuccessResponse([]cbiregions.CBIRegions{{ID: "region-001", Name: "US West"}, {ID: "region-002", Name: "US East"}}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbiregions.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCBIRegions_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	regionName := "US West"
	path := "/zpa/cbiconfig/cbi/api/customers/" + testCustomerID + "/regions"

	server.On("GET", path, common.SuccessResponse([]cbiregions.CBIRegions{
		{ID: "region-001", Name: regionName},
		{ID: "region-002", Name: "US East"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := cbiregions.GetByName(context.Background(), service, regionName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "region-001", result.ID)
	assert.Equal(t, regionName, result.Name)
}
