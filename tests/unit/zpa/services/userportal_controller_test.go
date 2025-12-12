// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

func TestUserPortalController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal/" + portalID

	server.On("GET", path, common.SuccessResponse(portal_controller.UserPortalController{
		ID:   portalID,
		Name: "Test Portal",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_controller.Get(context.Background(), service, portalID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, portalID, result.ID)
}

func TestUserPortalController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []portal_controller.UserPortalController{{ID: "portal-001"}, {ID: "portal-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_controller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUserPortalController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalName := "Production Portal"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []portal_controller.UserPortalController{
			{ID: "portal-001", Name: "Other Portal"},
			{ID: "portal-002", Name: portalName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_controller.GetByName(context.Background(), service, portalName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "portal-002", result.ID)
	assert.Equal(t, portalName, result.Name)
}

func TestUserPortalController_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal"

	server.On("POST", path, common.SuccessResponse(portal_controller.UserPortalController{
		ID:   "new-portal-123",
		Name: "New Portal",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newPortal := portal_controller.UserPortalController{
		Name: "New Portal",
	}

	result, _, err := portal_controller.Create(context.Background(), service, newPortal)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-portal-123", result.ID)
}

func TestUserPortalController_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal/" + portalID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updatePortal := &portal_controller.UserPortalController{
		ID:   portalID,
		Name: "Updated Portal",
	}

	resp, err := portal_controller.Update(context.Background(), service, portalID, updatePortal)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUserPortalController_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortal/" + portalID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := portal_controller.Delete(context.Background(), service, portalID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
