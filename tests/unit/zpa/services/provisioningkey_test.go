// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
)

func TestProvisioningKey_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := "key-12345"
	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey/" + keyID

	server.On("GET", path, common.SuccessResponse(provisioningkey.ProvisioningKey{
		ID:   keyID,
		Name: "Test Key",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := provisioningkey.Get(context.Background(), service, keyType, keyID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, keyID, result.ID)
}

func TestProvisioningKey_GetAllByAssociationType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []provisioningkey.ProvisioningKey{{ID: "key-001"}, {ID: "key-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, err := provisioningkey.GetAllByAssociationType(context.Background(), service, keyType)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestProvisioningKey_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyName := "Production Key"
	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []provisioningkey.ProvisioningKey{
			{ID: "key-001", Name: "Other Key"},
			{ID: "key-002", Name: keyName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := provisioningkey.GetByName(context.Background(), service, keyType, keyName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "key-002", result.ID)
	assert.Equal(t, keyName, result.Name)
}

func TestProvisioningKey_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey"

	server.On("POST", path, common.SuccessResponse(provisioningkey.ProvisioningKey{
		ID:   "new-key-123",
		Name: "New Key",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newKey := &provisioningkey.ProvisioningKey{
		Name: "New Key",
	}

	result, _, err := provisioningkey.Create(context.Background(), service, keyType, newKey)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-key-123", result.ID)
}

func TestProvisioningKey_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := "key-12345"
	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey/" + keyID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateKey := &provisioningkey.ProvisioningKey{
		ID:   keyID,
		Name: "Updated Key",
	}

	resp, err := provisioningkey.Update(context.Background(), service, keyType, keyID, updateKey)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestProvisioningKey_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := "key-12345"
	keyType := "CONNECTOR_GRP"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/associationType/" + keyType + "/provisioningKey/" + keyID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := provisioningkey.Delete(context.Background(), service, keyType, keyID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
