// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/aup"
)

func TestUserPortalAUP_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	aupID := "aup-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup/" + aupID

	server.On("GET", path, common.SuccessResponse(aup.UserPortalAup{
		ID:   aupID,
		Name: "Test AUP",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := aup.Get(context.Background(), service, aupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, aupID, result.ID)
}

func TestUserPortalAUP_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []aup.UserPortalAup{{ID: "aup-001"}, {ID: "aup-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := aup.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUserPortalAUP_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	aupName := "Production AUP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []aup.UserPortalAup{
			{ID: "aup-001", Name: "Other AUP"},
			{ID: "aup-002", Name: aupName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := aup.GetByName(context.Background(), service, aupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "aup-002", result.ID)
	assert.Equal(t, aupName, result.Name)
}

func TestUserPortalAUP_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup"

	server.On("POST", path, common.SuccessResponse(aup.UserPortalAup{
		ID:   "new-aup-123",
		Name: "New AUP",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newAUP := &aup.UserPortalAup{
		Name: "New AUP",
	}

	result, _, err := aup.Create(context.Background(), service, newAUP)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-aup-123", result.ID)
}

func TestUserPortalAUP_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	aupID := "aup-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup/" + aupID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateAUP := &aup.UserPortalAup{
		ID:   aupID,
		Name: "Updated AUP",
	}

	resp, err := aup.Update(context.Background(), service, aupID, updateAUP)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUserPortalAUP_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	aupID := "aup-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userportal/aup/" + aupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := aup.Delete(context.Background(), service, aupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
