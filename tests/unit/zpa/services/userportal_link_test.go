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
