// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

func TestAppServerController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverID := "server-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server/" + serverID

	server.On("GET", path, common.SuccessResponse(appservercontroller.ApplicationServer{
		ID:          serverID,
		Name:        "Test Server",
		Description: "Test description",
		Address:     "192.168.1.100",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := appservercontroller.Get(context.Background(), service, serverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, serverID, result.ID)
	assert.Equal(t, "Test Server", result.Name)
}

func TestAppServerController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverName := "Production Server"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appservercontroller.ApplicationServer{
			{ID: "server-001", Name: "Other Server", Enabled: true},
			{ID: "server-002", Name: serverName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appservercontroller.GetByName(context.Background(), service, serverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "server-002", result.ID)
	assert.Equal(t, serverName, result.Name)
}

func TestAppServerController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server"

	server.On("POST", path, common.SuccessResponse(appservercontroller.ApplicationServer{
		ID:          "new-server-123",
		Name:        "New Server",
		Description: "Created via unit test",
		Address:     "192.168.1.101",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newServer := appservercontroller.ApplicationServer{
		Name:        "New Server",
		Description: "Created via unit test",
		Address:     "192.168.1.101",
		Enabled:     true,
	}

	result, _, err := appservercontroller.Create(context.Background(), service, newServer)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-server-123", result.ID)
}

func TestAppServerController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverID := "server-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server/" + serverID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateServer := appservercontroller.ApplicationServer{
		ID:          serverID,
		Name:        "Updated Server",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := appservercontroller.Update(context.Background(), service, serverID, updateServer)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppServerController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverID := "server-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server/" + serverID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := appservercontroller.Delete(context.Background(), service, serverID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppServerController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appservercontroller.ApplicationServer{
			{ID: "server-001", Name: "Server 1", Enabled: true},
			{ID: "server-002", Name: "Server 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appservercontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestAppServerController_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []appservercontroller.ApplicationServer{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appservercontroller.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no application server named")
}

func TestAppServerController_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/server/" + serverID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appservercontroller.Get(context.Background(), service, serverID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
