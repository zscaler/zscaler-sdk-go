// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/machinegroup"
)

func TestMachineGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "mg-12345"
	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/machineGroup/{id}
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/machineGroup/" + groupID

	server.On("GET", path, common.SuccessResponse(machinegroup.MachineGroup{
		ID:          groupID,
		Name:        "Test Machine Group",
		Description: "Test description",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := machinegroup.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Test Machine Group", result.Name)
}

func TestMachineGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Production Machine Group"
	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/machineGroup (v1, not v2)
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/machineGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []machinegroup.MachineGroup{
			{ID: "mg-001", Name: "Other Group", Enabled: true},
			{ID: "mg-002", Name: groupName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := machinegroup.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "mg-002", result.ID)
	assert.Equal(t, groupName, result.Name)
}

func TestMachineGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/machineGroup (v1, not v2)
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/machineGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []machinegroup.MachineGroup{
			{ID: "mg-001", Name: "Group 1", Enabled: true},
			{ID: "mg-002", Name: "Group 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := machinegroup.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestMachineGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/machineGroup (v1, not v2)
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/machineGroup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []machinegroup.MachineGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := machinegroup.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no machine group named")
}

func TestMachineGroup_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/machineGroup/" + groupID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := machinegroup.Get(context.Background(), service, groupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
