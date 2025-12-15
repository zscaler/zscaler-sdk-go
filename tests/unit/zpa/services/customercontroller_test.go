// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	customercontroller "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/customer_controller"
)

func TestCustomerController_GetAllAuthDomains_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/authDomains"

	server.On("GET", path, common.SuccessResponse(customercontroller.AuthDomain{
		AuthDomains: []string{"example.com", "test.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := customercontroller.GetAllAuthDomains(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.AuthDomains, 2)
	assert.Contains(t, result.AuthDomains, "example.com")
}

func TestCustomerController_GetAncestorPolicy_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/ancestorPolicy"

	server.On("GET", path, common.SuccessResponse(customercontroller.AncestorPolicy{
		AccessType: "FULL_ACCESS",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := customercontroller.GetAncestorPolicy(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "FULL_ACCESS", result.AccessType)
}

// Note: Create test omitted as it uses query params that are difficult to mock accurately
