// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/role_controller"
)

func TestRoleController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := "role-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles/" + roleID

	server.On("GET", path, common.SuccessResponse(role_controller.RoleController{
		ID:   roleID,
		Name: "Test Role",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := role_controller.Get(context.Background(), service, roleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, roleID, result.ID)
}

func TestRoleController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []role_controller.RoleController{{ID: "role-001"}, {ID: "role-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := role_controller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRoleController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleName := "Admin Role"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []role_controller.RoleController{
			{ID: "role-001", Name: "Other Role"},
			{ID: "role-002", Name: roleName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := role_controller.GetByName(context.Background(), service, roleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "role-002", result.ID)
	assert.Equal(t, roleName, result.Name)
}

func TestRoleController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles"

	server.On("POST", path, common.SuccessResponse(role_controller.RoleController{
		ID:   "new-role-123",
		Name: "New Role",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newRole := &role_controller.RoleController{
		Name: "New Role",
	}

	result, _, err := role_controller.Create(context.Background(), service, newRole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-role-123", result.ID)
}

func TestRoleController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := "role-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles/" + roleID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := role_controller.Delete(context.Background(), service, roleID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRoleController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	roleID := "role-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/roles/" + roleID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateRole := &role_controller.RoleController{
		ID:   roleID,
		Name: "Updated Role",
	}

	resp, err := role_controller.Update(context.Background(), service, roleID, updateRole)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRoleController_GetPermissionGroups_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/permissionGroups"

	// This endpoint returns a raw array, not a paginated response
	server.On("GET", path, common.SuccessResponse([]role_controller.ClassPermissionGroup{
		{ID: "group-001", Name: "Admin Permissions"},
		{ID: "group-002", Name: "Read-Only Permissions"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := role_controller.GetPermissionGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
