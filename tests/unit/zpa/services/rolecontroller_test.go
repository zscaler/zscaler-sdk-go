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
