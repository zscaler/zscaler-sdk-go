// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
)

func TestSamlAttribute_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	attrID := "saml-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/samlAttribute/" + attrID

	server.On("GET", path, common.SuccessResponse(samlattribute.SamlAttribute{
		ID:               attrID,
		Name:             "Test SAML Attribute",
		IdpName:          "Test IDP",
		SamlName:         "email",
		UserAttribute:    true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := samlattribute.Get(context.Background(), service, attrID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, attrID, result.ID)
	assert.Equal(t, "Test SAML Attribute", result.Name)
}

func TestSamlAttribute_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	attrName := "email"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/samlAttribute"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []samlattribute.SamlAttribute{
			{ID: "saml-001", Name: "Other Attribute"},
			{ID: "saml-002", Name: attrName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := samlattribute.GetByName(context.Background(), service, attrName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "saml-002", result.ID)
	assert.Equal(t, attrName, result.Name)
}

func TestSamlAttribute_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/samlAttribute"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []samlattribute.SamlAttribute{
			{ID: "saml-001", Name: "Attribute 1"},
			{ID: "saml-002", Name: "Attribute 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := samlattribute.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestSamlAttribute_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/samlAttribute"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []samlattribute.SamlAttribute{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := samlattribute.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no saml attribute named")
}

func TestSamlAttribute_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	attrID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/samlAttribute/" + attrID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := samlattribute.Get(context.Background(), service, attrID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
