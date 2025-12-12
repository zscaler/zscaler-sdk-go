// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

func TestServiceEdgeGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "seg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(serviceedgegroup.ServiceEdgeGroup{
		ID:          groupID,
		Name:        "Test Service Edge Group",
		Description: "Test description",
		Enabled:     true,
		Latitude:    "37.7749",
		Longitude:   "-122.4194",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := serviceedgegroup.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Test Service Edge Group", result.Name)
}

func TestServiceEdgeGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Service Edge Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []serviceedgegroup.ServiceEdgeGroup{
			{ID: "seg-001", Name: "Other Group", Enabled: true},
			{ID: "seg-002", Name: groupName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgegroup.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "seg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

func TestServiceEdgeGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup"

	server.On("POST", path, common.SuccessResponse(serviceedgegroup.ServiceEdgeGroup{
		ID:          "new-seg-123",
		Name:        "New Service Edge Group",
		Description: "Created via unit test",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newGroup := serviceedgegroup.ServiceEdgeGroup{
		Name:        "New Service Edge Group",
		Description: "Created via unit test",
		Enabled:     true,
	}

	createdGroup, _, err := serviceedgegroup.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, createdGroup)
	assert.Equal(t, "new-seg-123", createdGroup.ID)
}

func TestServiceEdgeGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "seg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup/" + groupID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateGroup := &serviceedgegroup.ServiceEdgeGroup{
		ID:          groupID,
		Name:        "Updated Service Edge Group",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := serviceedgegroup.Update(context.Background(), service, groupID, updateGroup)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServiceEdgeGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "seg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup/" + groupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := serviceedgegroup.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestServiceEdgeGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []serviceedgegroup.ServiceEdgeGroup{
			{ID: "seg-001", Name: "Group 1", Enabled: true},
			{ID: "seg-002", Name: "Group 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgegroup.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestServiceEdgeGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []serviceedgegroup.ServiceEdgeGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgegroup.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no service edge group named")
}

func TestServiceEdgeGroup_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeGroup/" + groupID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgegroup.Get(context.Background(), service, groupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
