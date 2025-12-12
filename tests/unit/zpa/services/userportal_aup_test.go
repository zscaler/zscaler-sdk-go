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
