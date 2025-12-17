// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
)

const testCustomerID = "123456789"

func TestAppConnectorController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	connectorID := "connector-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector/" + connectorID

	server.On("GET", path, common.SuccessResponse(appconnectorcontroller.AppConnector{
		ID:          connectorID,
		Name:        "Test Connector",
		Description: "Test description",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := appconnectorcontroller.Get(context.Background(), service, connectorID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, connectorID, result.ID)
	assert.Equal(t, "Test Connector", result.Name)
}

func TestAppConnectorController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	connectorName := "Production Connector"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appconnectorcontroller.AppConnector{
			{ID: "conn-001", Name: "Other Connector", Enabled: true},
			{ID: "conn-002", Name: connectorName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorcontroller.GetByName(context.Background(), service, connectorName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "conn-002", result.ID)
	assert.Equal(t, connectorName, result.Name)
}

func TestAppConnectorController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appconnectorcontroller.AppConnector{
			{ID: "conn-001", Name: "Connector 1", Enabled: true},
			{ID: "conn-002", Name: "Connector 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorcontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestAppConnectorController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	connectorID := "connector-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector/" + connectorID

	server.On("PUT", path, common.SuccessResponse(appconnectorcontroller.AppConnector{
		ID:      connectorID,
		Name:    "Updated Connector",
		Enabled: true,
	}))
	server.On("GET", path, common.SuccessResponse(appconnectorcontroller.AppConnector{
		ID:      connectorID,
		Name:    "Updated Connector",
		Enabled: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateConnector := appconnectorcontroller.AppConnector{
		ID:          connectorID,
		Name:        "Updated Connector",
		Description: "Updated description",
		Enabled:     true,
	}

	result, _, err := appconnectorcontroller.Update(context.Background(), service, connectorID, updateConnector)

	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestAppConnectorController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	connectorID := "connector-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector/" + connectorID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := appconnectorcontroller.Delete(context.Background(), service, connectorID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppConnectorController_BulkDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector/bulkDelete"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	ids := []string{"conn-001", "conn-002", "conn-003"}
	resp, err := appconnectorcontroller.BulkDelete(context.Background(), service, ids)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppConnectorController_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []appconnectorcontroller.AppConnector{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorcontroller.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no app connector named")
}

func TestAppConnectorController_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	connectorID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connector/" + connectorID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorcontroller.Get(context.Background(), service, connectorID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
