// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
)

func TestAppConnectorGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "acg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(appconnectorgroup.AppConnectorGroup{
		ID:          groupID,
		Name:        "Test Connector Group",
		Description: "Test description",
		Enabled:     true,
		Latitude:    "37.7749",
		Longitude:   "-122.4194",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := appconnectorgroup.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Test Connector Group", result.Name)
}

func TestAppConnectorGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appconnectorgroup.AppConnectorGroup{
			{ID: "acg-001", Name: "Other Group", Enabled: true},
			{ID: "acg-002", Name: groupName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorgroup.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "acg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

func TestAppConnectorGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup"

	server.On("POST", path, common.SuccessResponse(appconnectorgroup.AppConnectorGroup{
		ID:          "new-acg-123",
		Name:        "New Connector Group",
		Description: "Created via unit test",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newGroup := appconnectorgroup.AppConnectorGroup{
		Name:        "New Connector Group",
		Description: "Created via unit test",
		Enabled:     true,
	}

	result, _, err := appconnectorgroup.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-acg-123", result.ID)
	assert.Equal(t, "New Connector Group", result.Name)
}

func TestAppConnectorGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "acg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/" + groupID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateGroup := &appconnectorgroup.AppConnectorGroup{
		ID:          groupID,
		Name:        "Updated Group",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := appconnectorgroup.Update(context.Background(), service, groupID, updateGroup)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppConnectorGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "acg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/" + groupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := appconnectorgroup.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAppConnectorGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []appconnectorgroup.AppConnectorGroup{
			{ID: "acg-001", Name: "Group 1", Enabled: true},
			{ID: "acg-002", Name: "Group 2", Enabled: false},
			{ID: "acg-003", Name: "Group 3", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorgroup.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestAppConnectorGroup_GetSummary_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []map[string]interface{}{
			{"id": "acg-001", "name": "Group 1"},
			{"id": "acg-002", "name": "Group 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorgroup.GetAppconnectorGroupSummary(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestAppConnectorGroup_GetSG_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "acg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/" + groupID + "/sg"

	server.On("GET", path, common.SuccessResponse(appconnectorgroup.AppConnectorGroup{
		ID:      groupID,
		Name:    "Test Connector Group",
		Enabled: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := appconnectorgroup.GetAppConnectorGroupSG(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, result.ID)
}

func TestAppConnectorGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []appconnectorgroup.AppConnectorGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorgroup.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no app connector group named")
}

func TestAppConnectorGroup_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/appConnectorGroup/" + groupID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorgroup.Get(context.Background(), service, groupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
