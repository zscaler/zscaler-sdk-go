// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

func TestMicrotenants_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	mtID := "mt-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants/" + mtID

	server.On("GET", path, common.SuccessResponse(microtenants.MicroTenant{
		ID:          mtID,
		Name:        "Test Microtenant",
		Description: "Test description",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := microtenants.Get(context.Background(), service, mtID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, mtID, result.ID)
	assert.Equal(t, "Test Microtenant", result.Name)
}

func TestMicrotenants_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	mtName := "Production Microtenant"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []microtenants.MicroTenant{
			{ID: "mt-001", Name: "Other MT", Enabled: true},
			{ID: "mt-002", Name: mtName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := microtenants.GetByName(context.Background(), service, mtName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "mt-002", result.ID)
	assert.Equal(t, mtName, result.Name)
}

func TestMicrotenants_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants"

	server.On("POST", path, common.SuccessResponse(microtenants.MicroTenant{
		ID:          "new-mt-123",
		Name:        "New Microtenant",
		Description: "Created via unit test",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newMT := microtenants.MicroTenant{
		Name:        "New Microtenant",
		Description: "Created via unit test",
		Enabled:     true,
	}

	result, _, err := microtenants.Create(context.Background(), service, newMT)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-mt-123", result.ID)
}

func TestMicrotenants_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	mtID := "mt-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants/" + mtID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateMT := &microtenants.MicroTenant{
		ID:          mtID,
		Name:        "Updated Microtenant",
		Description: "Updated description",
		Enabled:     false,
	}

	resp, err := microtenants.Update(context.Background(), service, mtID, updateMT)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestMicrotenants_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	mtID := "mt-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants/" + mtID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := microtenants.Delete(context.Background(), service, mtID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestMicrotenants_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []microtenants.MicroTenant{
			{ID: "mt-001", Name: "MT 1", Enabled: true},
			{ID: "mt-002", Name: "MT 2", Enabled: false},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := microtenants.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestMicrotenants_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []microtenants.MicroTenant{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := microtenants.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no microtenant named")
}

func TestMicrotenants_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	mtID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/microtenants/" + mtID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := microtenants.Get(context.Background(), service, mtID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
