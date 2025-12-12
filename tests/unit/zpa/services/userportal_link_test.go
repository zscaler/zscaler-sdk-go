// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

func TestUserPortalLink_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	linkID := "link-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink/" + linkID

	server.On("GET", path, common.SuccessResponse(portal_link.UserPortalLink{
		ID:   linkID,
		Name: "Test Link",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_link.Get(context.Background(), service, linkID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, linkID, result.ID)
}

func TestUserPortalLink_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []portal_link.UserPortalLink{{ID: "link-001"}, {ID: "link-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_link.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUserPortalLink_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	linkName := "Production Link"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []portal_link.UserPortalLink{
			{ID: "link-001", Name: "Other Link"},
			{ID: "link-002", Name: linkName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := portal_link.GetByName(context.Background(), service, linkName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "link-002", result.ID)
	assert.Equal(t, linkName, result.Name)
}

func TestUserPortalLink_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink"

	server.On("POST", path, common.SuccessResponse(portal_link.UserPortalLink{
		ID:   "new-link-123",
		Name: "New Link",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newLink := portal_link.UserPortalLink{
		Name: "New Link",
	}

	result, _, err := portal_link.Create(context.Background(), service, newLink)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-link-123", result.ID)
}

func TestUserPortalLink_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	linkID := "link-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink/" + linkID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateLink := &portal_link.UserPortalLink{
		ID:   linkID,
		Name: "Updated Link",
	}

	resp, err := portal_link.Update(context.Background(), service, linkID, updateLink)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestUserPortalLink_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	linkID := "link-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/userPortalLink/" + linkID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := portal_link.Delete(context.Background(), service, linkID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
