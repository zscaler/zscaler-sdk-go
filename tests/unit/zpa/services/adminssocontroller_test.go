// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/admin_sso_controller"
)

func TestAdminSSOController_GetSSOLoginController_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/v2/ssoLoginOptions"

	server.On("GET", path, common.SuccessResponse(admin_sso_controller.AdminSSOLoginOptions{
		SSOLoginOnly: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := admin_sso_controller.GetSSOLoginController(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.SSOLoginOnly)
}

// Note: UpdateSSOLoginController test omitted as it uses query params that are difficult to mock
