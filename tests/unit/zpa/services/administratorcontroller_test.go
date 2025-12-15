// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/administrator_controller"
)

func TestAdministratorController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := "admin-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators/" + adminID

	server.On("GET", path, common.SuccessResponse(administrator_controller.AdministratorController{
		ID:          adminID,
		Email:       "admin@example.com",
		DisplayName: "Test Admin",
		IsEnabled:   true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := administrator_controller.Get(context.Background(), service, adminID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, adminID, result.ID)
	assert.Equal(t, "admin@example.com", result.Email)
}

func TestAdministratorController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminName := "Test Admin"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []administrator_controller.AdministratorController{
			{ID: "admin-001", Username: "Other Admin"},
			{ID: "admin-002", Username: adminName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := administrator_controller.GetByName(context.Background(), service, adminName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "admin-002", result.ID)
	assert.Equal(t, adminName, result.Username)
}

func TestAdministratorController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []administrator_controller.AdministratorController{
			{ID: "admin-001", Username: "Admin 1"},
			{ID: "admin-002", Username: "Admin 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := administrator_controller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestAdministratorController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators"

	server.On("POST", path, common.SuccessResponse(administrator_controller.AdministratorController{
		ID:          "new-admin-123",
		Email:       "newadmin@example.com",
		DisplayName: "New Admin",
		IsEnabled:   true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newAdmin := &administrator_controller.AdministratorController{
		Email:       "newadmin@example.com",
		DisplayName: "New Admin",
		IsEnabled:   true,
	}

	result, _, err := administrator_controller.Create(context.Background(), service, newAdmin)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-admin-123", result.ID)
}

func TestAdministratorController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := "admin-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators/" + adminID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateAdmin := &administrator_controller.AdministratorController{
		ID:          adminID,
		DisplayName: "Updated Admin",
		IsEnabled:   false,
	}

	resp, err := administrator_controller.Update(context.Background(), service, adminID, updateAdmin)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAdministratorController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	adminID := "admin-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/administrators/" + adminID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := administrator_controller.Delete(context.Background(), service, adminID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
