// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_custom_controls"
)

func TestInspectionCustomControls_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	controlID := "ctrl-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/custom/" + controlID

	// The Get function parses ControlRuleJson, so include a valid one
	server.On("GET", path, common.SuccessResponse(inspection_custom_controls.InspectionCustomControl{
		ID:              controlID,
		Name:            "Test Custom Control",
		ControlRuleJson: "[]", // Empty valid JSON array
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_custom_controls.Get(context.Background(), service, controlID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, controlID, result.ID)
}

func TestInspectionCustomControls_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/custom"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []inspection_custom_controls.InspectionCustomControl{{ID: "ctrl-001"}, {ID: "ctrl-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_custom_controls.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestInspectionCustomControls_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/custom"

	server.On("POST", path, common.SuccessResponse(inspection_custom_controls.InspectionCustomControl{
		ID:   "new-ctrl-123",
		Name: "New Control",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newCtrl := inspection_custom_controls.InspectionCustomControl{
		Name: "New Control",
	}

	result, _, err := inspection_custom_controls.Create(context.Background(), service, newCtrl)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-ctrl-123", result.ID)
}

func TestInspectionCustomControls_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	controlID := "ctrl-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/custom/" + controlID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateCtrl := &inspection_custom_controls.InspectionCustomControl{
		ID:   controlID,
		Name: "Updated Control",
	}

	resp, err := inspection_custom_controls.Update(context.Background(), service, controlID, updateCtrl)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestInspectionCustomControls_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	controlID := "ctrl-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/custom/" + controlID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := inspection_custom_controls.Delete(context.Background(), service, controlID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
