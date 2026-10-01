// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/scimgroup"
)

func TestScimGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := "12345"
	path := "/zpa/userconfig/v1/customers/" + testCustomerID + "/scimgroup/" + groupID

	server.On("GET", path, common.SuccessResponse(scimgroup.ScimGroup{
		ID:   12345,
		Name: "Test SCIM Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := scimgroup.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(12345), result.ID)
}

func TestScimGroup_GetAllByIdpId_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	idpID := "idp-12345"
	path := "/zpa/userconfig/v1/customers/" + testCustomerID + "/scimgroup/idpId/" + idpID

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []scimgroup.ScimGroup{{ID: 1, Name: "Group 1"}, {ID: 2, Name: "Group 2"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := scimgroup.GetAllByIdpId(context.Background(), service, idpID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestScimGroup_GetByName_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-777"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	wantName := "Engineering SCIM"
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Other"},
		{ID: 2, Name: wantName},
	})))

	got, _, err := scimgroup.GetByName(context.Background(), api.Service, wantName, idpID)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(2), got.ID)
	assert.Equal(t, wantName, got.Name)
}

func TestScimGroup_GetByName_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-888"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Only Group"},
	})))

	got, _, err := scimgroup.GetByName(context.Background(), api.Service, "missing-group", idpID)
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "missing-group")
}

func TestScimGroup_GetByNameAndIamIdp_ByIamIdpID_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-901"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering", IamIdpID: "iam-a", IamIdpName: "Okta"},
		{ID: 2, Name: "Engineering", IamIdpID: "iam-b", IamIdpName: "Entra"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "Engineering", idpID,
		scimgroup.IamIdpCriteria{IamIdpID: "iam-b"})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(2), got.ID)
	assert.Equal(t, "Entra", got.IamIdpName)
}

func TestScimGroup_GetByNameAndIamIdp_ByIamIdpName_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-902"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering", IamIdpID: "iam-a", IamIdpName: "Okta"},
		{ID: 2, Name: "Engineering", IamIdpID: "iam-b", IamIdpName: "Entra"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "engineering", idpID,
		scimgroup.IamIdpCriteria{IamIdpName: "okta"})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.ID)
}

func TestScimGroup_GetByNameAndIamIdp_DecoratedName_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-903"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering(Okta)"},
		{ID: 2, Name: "Engineering(Entra)"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "Engineering", idpID,
		scimgroup.IamIdpCriteria{IamIdpName: "Entra"})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(2), got.ID)
}

func TestScimGroup_GetByNameAndIamIdp_Ambiguous_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-904"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering", IamIdpID: "iam-a", IamIdpName: "Okta"},
		{ID: 2, Name: "Engineering", IamIdpID: "iam-a", IamIdpName: "Okta"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "Engineering", idpID,
		scimgroup.IamIdpCriteria{IamIdpName: "Okta"})
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "multiple SCIM groups")
	assert.Contains(t, err.Error(), "id=1")
	assert.Contains(t, err.Error(), "id=2")
}

func TestScimGroup_GetByNameAndIamIdp_NoIamIdpMatch_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-905"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering", IamIdpID: "iam-a", IamIdpName: "Okta"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "Engineering", idpID,
		scimgroup.IamIdpCriteria{IamIdpID: "iam-z"})
	require.Error(t, err)
	require.Nil(t, got)
	assert.Contains(t, err.Error(), "IAM IdP ID 'iam-z'")
	assert.Contains(t, err.Error(), `iamIdpName="Okta"`)
}

func TestScimGroup_GetByNameAndIamIdp_EmptyCriteria_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	idpID := "idp-906"
	path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup", "idpId", idpID)
	api.On("GET", path, common.SuccessResponse(common.ZPAList([]scimgroup.ScimGroup{
		{ID: 1, Name: "Engineering", IamIdpID: "iam-a"},
		{ID: 2, Name: "Engineering", IamIdpID: "iam-b"},
	})))

	got, _, err := scimgroup.GetByNameAndIamIdp(context.Background(), api.Service, "Engineering", idpID,
		scimgroup.IamIdpCriteria{})
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.ID, "empty criteria must keep GetByName's first-match behavior")
}
