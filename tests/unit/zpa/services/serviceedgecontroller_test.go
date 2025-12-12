// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
)

func TestServiceEdgeController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	edgeID := "se-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge/" + edgeID

	server.On("GET", path, common.SuccessResponse(serviceedgecontroller.ServiceEdgeController{
		ID:   edgeID,
		Name: "Test Service Edge",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgecontroller.Get(context.Background(), service, edgeID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, edgeID, result.ID)
}

func TestServiceEdgeController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []serviceedgecontroller.ServiceEdgeController{{ID: "se-001"}, {ID: "se-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgecontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestServiceEdgeController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	edgeName := "Production Edge"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []serviceedgecontroller.ServiceEdgeController{
			{ID: "se-001", Name: "Other Edge"},
			{ID: "se-002", Name: edgeName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgecontroller.GetByName(context.Background(), service, edgeName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "se-002", result.ID)
	assert.Equal(t, edgeName, result.Name)
}

func TestServiceEdgeController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	edgeID := "se-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge/" + edgeID

	// Update calls PUT first, then GET to retrieve the updated resource
	server.On("PUT", path, common.NoContentResponse())
	server.On("GET", path, common.SuccessResponse(serviceedgecontroller.ServiceEdgeController{
		ID:   edgeID,
		Name: "Updated Edge",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateEdge := serviceedgecontroller.ServiceEdgeController{
		ID:   edgeID,
		Name: "Updated Edge",
	}

	result, _, err := serviceedgecontroller.Update(context.Background(), service, edgeID, updateEdge)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestServiceEdgeController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	edgeID := "se-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge/" + edgeID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := serviceedgecontroller.Delete(context.Background(), service, edgeID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServiceEdgeController_BulkDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdge/bulkDelete"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	ids := []string{"se-001", "se-002", "se-003"}
	resp, err := serviceedgecontroller.BulkDelete(context.Background(), service, ids)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
