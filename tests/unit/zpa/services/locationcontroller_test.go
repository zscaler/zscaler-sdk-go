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

func TestLocationController_GetLocationExtranetResource_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	zpnErID := "er-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/location/extranetResource/" + zpnErID

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "loc-001", Name: "Location 1"}, {ID: "loc-002", Name: "Location 2"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := location_controller.GetLocationExtranetResource(context.Background(), service, zpnErID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLocationController_GetLocationSummary_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/location/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []zpacommon.CommonSummary{{ID: "loc-001", Name: "Location 1"}, {ID: "loc-002", Name: "Location 2"}},
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

	locationName := "My Location"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/location/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []zpacommon.CommonSummary{
			{ID: "loc-001", Name: "Other Location"},
			{ID: "loc-002", Name: locationName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := location_controller.GetLocationSummaryByName(context.Background(), service, locationName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "loc-002", result.ID)
	assert.Equal(t, locationName, result.Name)
}

func TestLocationController_GetLocationGroupExtranetResource_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	zpnErID := "er-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/locationGroup/extranetResource/" + zpnErID

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []zpacommon.LocationGroupDTO{
			{ID: "locgrp-001", Name: "Location Group 1"},
			{ID: "locgrp-002", Name: "Location Group 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := location_controller.GetLocationGroupExtranetResource(context.Background(), service, zpnErID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
