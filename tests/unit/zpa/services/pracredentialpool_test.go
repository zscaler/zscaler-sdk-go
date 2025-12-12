// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func TestPRACredentialPool_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	poolID := "pool-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credentialPool/" + poolID

	server.On("GET", path, common.SuccessResponse(pracredentialpool.CredentialPool{
		ID:   poolID,
		Name: "Test Credential Pool",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredentialpool.Get(context.Background(), service, poolID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, poolID, result.ID)
}

func TestPRACredentialPool_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credentialPool"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []pracredentialpool.CredentialPool{{ID: "pool-001"}, {ID: "pool-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredentialpool.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
