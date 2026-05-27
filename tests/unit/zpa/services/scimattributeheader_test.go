// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimattributeheader"
)

func TestScimAttributeHeader_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "idp-12345"
	attrID := "attr-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/idp/" + idpID + "/scimattribute/" + attrID

	server.On("GET", path, common.SuccessResponse(scimattributeheader.ScimAttributeHeader{
		ID:   attrID,
		Name: "Test SCIM Attribute",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := scimattributeheader.Get(context.Background(), service, idpID, attrID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, attrID, result.ID)
}

func TestScimAttributeHeader_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "idp-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/idp/" + idpID + "/scimattribute"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []scimattributeheader.ScimAttributeHeader{{ID: "attr-001"}, {ID: "attr-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := scimattributeheader.GetAllByIdpId(context.Background(), service, idpID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestScimAttributeHeader_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "idp-12345"
	attrName := "email"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/idp/" + idpID + "/scimattribute"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []scimattributeheader.ScimAttributeHeader{
			{ID: "attr-001", Name: "otherAttr"},
			{ID: "attr-002", Name: attrName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := scimattributeheader.GetByName(context.Background(), service, attrName, idpID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "attr-002", result.ID)
	assert.Equal(t, attrName, result.Name)
}

func TestScimAttributeHeader_GetValues_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-values"
	attrID := "attr-collect"
	path := fmt.Sprintf("/zpa/userconfig/v1/customers/%s/scimattribute/idpId/%s/attributeId/%s", api.CustomerID, idpID, attrID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]string{"alice", "bob"})))

	got, err := scimattributeheader.GetValues(context.Background(), api.Service, idpID, attrID)
	require.NoError(t, err)
	assert.Equal(t, []string{"alice", "bob"}, got)
}

func TestScimAttributeHeader_SearchValues_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-search"
	attrID := "attr-search"
	path := fmt.Sprintf("/zpa/userconfig/v1/customers/%s/scimattribute/idpId/%s/attributeId/%s", api.CustomerID, idpID, attrID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]string{"jdoe", "jane"})))

	got, err := scimattributeheader.SearchValues(context.Background(), api.Service, idpID, attrID, "jdoe@example.com")
	require.NoError(t, err)
	assert.Equal(t, []string{"jdoe", "jane"}, got)
}

func TestScimAttributeHeader_GetByName_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-404-attrs"
	path := common.ZPAPath(api.CustomerID, "idp", idpID, "scimattribute")
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimattributeheader.ScimAttributeHeader{
		{ID: "a1", Name: "department"},
	})))

	got, _, err := scimattributeheader.GetByName(context.Background(), api.Service, "costCenter", idpID)
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "costCenter")
}
