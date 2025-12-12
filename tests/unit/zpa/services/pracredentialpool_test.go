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
	// Correct path: /zpa/waap-pra-config/v1/admin/customers/{customerId}/credential-pool/{id}
	path := "/zpa/waap-pra-config/v1/admin/customers/" + testCustomerID + "/credential-pool/" + poolID

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

	// Correct path: /zpa/waap-pra-config/v1/admin/customers/{customerId}/credential-pool
	path := "/zpa/waap-pra-config/v1/admin/customers/" + testCustomerID + "/credential-pool"

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

func TestPRACredentialPool_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	poolName := "Production Pool"
	path := "/zpa/waap-pra-config/v1/admin/customers/" + testCustomerID + "/credential-pool"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []pracredentialpool.CredentialPool{
			{ID: "pool-001", Name: "Other Pool"},
			{ID: "pool-002", Name: poolName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredentialpool.GetByName(context.Background(), service, poolName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pool-002", result.ID)
	assert.Equal(t, poolName, result.Name)
}

func TestPRACredentialPool_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	poolID := "pool-12345"
	path := "/zpa/waap-pra-config/v1/admin/customers/" + testCustomerID + "/credential-pool/" + poolID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updatePool := &pracredentialpool.CredentialPool{
		ID:   poolID,
		Name: "Updated Pool",
	}

	resp, err := pracredentialpool.Update(context.Background(), service, poolID, updatePool)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRACredentialPool_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	poolID := "pool-12345"
	path := "/zpa/waap-pra-config/v1/admin/customers/" + testCustomerID + "/credential-pool/" + poolID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := pracredentialpool.Delete(context.Background(), service, poolID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
