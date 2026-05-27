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
