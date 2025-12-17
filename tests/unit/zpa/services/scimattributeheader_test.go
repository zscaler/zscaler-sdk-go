// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
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
