// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/location_controller"
)

func TestLocationController_GetSummary_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/location/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "loc-001"}, {ID: "loc-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := location_controller.GetLocationSummary(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLocationController_GetLocationSummaryByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	locName := "Production Location"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/location/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []zpacommon.CommonSummary{
			{ID: "loc-001", Name: "Other Location"},
			{ID: "loc-002", Name: locName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := location_controller.GetLocationSummaryByName(context.Background(), service, locName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "loc-002", result.ID)
	assert.Equal(t, locName, result.Name)
}
