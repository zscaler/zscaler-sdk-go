// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

func TestServerGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(servergroup.ServerGroup{
		ID:          groupID,
		Name:        "Test Server Group",
		Description: "Test description",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := servergroup.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Test Server Group", result.Name)
}

func TestServerGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Server Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []servergroup.ServerGroup{
			{ID: "sg-001", Name: "Other Group", Enabled: true},
			{ID: "sg-002", Name: groupName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := servergroup.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "sg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

func TestServerGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup"

	server.On("POST", path, common.SuccessResponse(servergroup.ServerGroup{
		ID:          "new-sg-123",
		Name:        "New Server Group",
		Description: "Created via unit test",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newGroup := &servergroup.ServerGroup{
		Name:        "New Server Group",
		Description: "Created via unit test",
		Enabled:     true,
	}

	result, _, err := servergroup.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-sg-123", result.ID)
	assert.Equal(t, "New Server Group", result.Name)
}

func TestServerGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup/" + groupID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateGroup := &servergroup.ServerGroup{
		ID:          groupID,
		Name:        "Updated Server Group",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := servergroup.Update(context.Background(), service, groupID, updateGroup)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServerGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "sg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup/" + groupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := servergroup.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServerGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []servergroup.ServerGroup{
			{ID: "sg-001", Name: "Group 1", Enabled: true},
			{ID: "sg-002", Name: "Group 2", Enabled: false},
			{ID: "sg-003", Name: "Group 3", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := servergroup.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestServerGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []servergroup.ServerGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := servergroup.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no server group named")
}

func TestServerGroup_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serverGroup/" + groupID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := servergroup.Get(context.Background(), service, groupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
