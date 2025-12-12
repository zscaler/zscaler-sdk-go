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

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/inspectionControls/predefined"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []inspection_predefined_controls.PredefinedControls{{ID: "predefined-001"}, {ID: "predefined-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, err := inspection_predefined_controls.GetAll(context.Background(), service, "OWASP_CRS/3.3.0")

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
