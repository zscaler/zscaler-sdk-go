// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/inspectioncontrol/inspection_predefined_controls"
)

func TestInspectionPredefinedControls_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	controlID := "predefined-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/predefined/" + controlID

	server.On("GET", path, common.SuccessResponse(inspection_predefined_controls.PredefinedControls{
		ID:   controlID,
		Name: "Test Predefined Control",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := inspection_predefined_controls.Get(context.Background(), service, controlID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, controlID, result.ID)
}

func TestInspectionPredefinedControls_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// GetAll returns []ControlGroupItem, not []PredefinedControls
	// The path includes a version query parameter
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/predefined"

	// Return ControlGroupItem array (the actual response structure)
	server.On("GET", path, common.SuccessResponse([]inspection_predefined_controls.ControlGroupItem{
		{
			ControlGroup: "Protocol Issues",
			PredefinedInspectionControls: []inspection_predefined_controls.PredefinedControls{
				{ID: "predefined-001", Name: "Control 1"},
				{ID: "predefined-002", Name: "Control 2"},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, err := inspection_predefined_controls.GetAll(context.Background(), service, "OWASP_CRS/3.3.0")

	require.NoError(t, err)
	// GetAll flattens the ControlGroupItem into individual PredefinedControls
	assert.Len(t, result, 2)
}
