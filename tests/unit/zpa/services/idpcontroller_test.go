// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/idpcontroller"
)

func TestIdpController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "idp-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/idp/" + idpID

	server.On("GET", path, common.SuccessResponse(idpcontroller.IdpController{
		ID:          idpID,
		Name:        "Test IDP",
		Description: "Test description",
		Enabled:     true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := idpcontroller.Get(context.Background(), service, idpID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, idpID, result.ID)
	assert.Equal(t, "Test IDP", result.Name)
}

func TestIdpController_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpName := "Production IDP"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/idp"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []idpcontroller.IdpController{
			{ID: "idp-001", Name: "Other IDP", Enabled: true},
			{ID: "idp-002", Name: idpName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := idpcontroller.GetByName(context.Background(), service, idpName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "idp-002", result.ID)
	assert.Equal(t, idpName, result.Name)
}

func TestIdpController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/idp"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []idpcontroller.IdpController{
			{ID: "idp-001", Name: "IDP 1", Enabled: true},
			{ID: "idp-002", Name: "IDP 2", Enabled: false},
			{ID: "idp-003", Name: "IDP 3", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := idpcontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestIdpController_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/idp"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []idpcontroller.IdpController{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := idpcontroller.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no Idp-Controller named")
}

func TestIdpController_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/idp/" + idpID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := idpcontroller.Get(context.Background(), service, idpID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
