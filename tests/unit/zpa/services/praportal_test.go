// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func TestPRAPortal_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal/" + portalID

	server.On("GET", path, common.SuccessResponse(praportal.PRAPortal{
		ID:   portalID,
		Name: "Test Portal",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praportal.Get(context.Background(), service, portalID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, portalID, result.ID)
}

func TestPRAPortal_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []praportal.PRAPortal{{ID: "portal-001"}, {ID: "portal-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praportal.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPRAPortal_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalName := "Production Portal"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []praportal.PRAPortal{
			{ID: "portal-001", Name: "Other Portal"},
			{ID: "portal-002", Name: portalName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praportal.GetByName(context.Background(), service, portalName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "portal-002", result.ID)
	assert.Equal(t, portalName, result.Name)
}

func TestPRAPortal_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal"

	server.On("POST", path, common.SuccessResponse(praportal.PRAPortal{
		ID:   "new-portal-123",
		Name: "New Portal",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newPortal := &praportal.PRAPortal{
		Name: "New Portal",
	}

	result, _, err := praportal.Create(context.Background(), service, newPortal)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-portal-123", result.ID)
}

func TestPRAPortal_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal/" + portalID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updatePortal := &praportal.PRAPortal{
		ID:   portalID,
		Name: "Updated Portal",
	}

	resp, err := praportal.Update(context.Background(), service, portalID, updatePortal)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRAPortal_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	portalID := "portal-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praPortal/" + portalID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := praportal.Delete(context.Background(), service, portalID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
