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

func TestProvisioningKey_Get_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	keyType := "CONNECTOR_GRP"
	path := common.ZPAPath(api.CustomerID, "associationType", keyType, "provisioningKey", "missing")
	api.On("GET", path, common.NotFoundResponse())

	got, _, err := provisioningkey.Get(context.Background(), api.Service, keyType, "missing")
	require.Error(t, err)
	assert.Nil(t, got)
}

func TestProvisioningKey_GetAll_MultiAssociationTypes_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	types := []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"}
	keysPerType := map[string][]provisioningkey.ProvisioningKey{
		"CONNECTOR_GRP":    {{ID: "key-c1", Name: "Connector Key"}},
		"SERVICE_EDGE_GRP": {{ID: "key-se1", Name: "Edge Key"}},
		"NP_ASSISTANT_GRP": {},
	}
	for _, assocType := range types {
		path := common.ZPAPath(api.CustomerID, "associationType", assocType, "provisioningKey")
		items := keysPerType[assocType]
		api.On("GET", path, common.SuccessResponse(common.ZPAList(items)))
	}

	result, err := provisioningkey.GetAll(context.Background(), api.Service)
	require.NoError(t, err)
	require.Len(t, result, 2)
}

func TestProvisioningKey_GetAllByZComponentID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyType := "CONNECTOR_GRP"
	zcomponentID := "acg-123"
	path := common.ZPAPath(testCustomerID, "associationType", keyType, "zcomponent", zcomponentID, "provisioningKey")

	server.On("GET", path, common.SuccessResponse([]provisioningkey.ProvisioningKey{
		{ID: "key-001"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, err := provisioningkey.GetAllByZComponentID(context.Background(), service, keyType, zcomponentID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestProvisioningKey_GetByNameAllAssociations_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyName := "Shared Key"
	for _, assocType := range []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"} {
		path := common.ZPAPath(testCustomerID, "associationType", assocType, "provisioningKey")
		if assocType == "SERVICE_EDGE_GRP" {
			server.On("GET", path, common.SuccessResponse(map[string]interface{}{
				"list": []provisioningkey.ProvisioningKey{
					{ID: "key-001", Name: keyName},
				},
				"totalPages": 1,
			}))
		} else {
			server.On("GET", path, common.SuccessResponse(map[string]interface{}{
				"list":       []provisioningkey.ProvisioningKey{},
				"totalPages": 1,
			}))
		}
	}

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, assocType, _, err := provisioningkey.GetByNameAllAssociations(context.Background(), service, keyName)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "SERVICE_EDGE_GRP", assocType)
}

func TestProvisioningKey_GetByIDAllAssociations_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyID := "key-123"
	for _, assocType := range []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"} {
		path := common.ZPAPath(testCustomerID, "associationType", assocType, "provisioningKey", keyID)
		if assocType == "CONNECTOR_GRP" {
			server.On("GET", path, common.SuccessResponse(provisioningkey.ProvisioningKey{
				ID: keyID, Name: "My Key",
			}))
		} else {
			server.On("GET", path, common.NotFoundResponse())
		}
	}

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, assocType, _, err := provisioningkey.GetByIDAllAssociations(context.Background(), service, keyID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "CONNECTOR_GRP", assocType)
}

func TestProvisioningKey_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	keyType := "CONNECTOR_GRP"
	path := common.ZPAPath(testCustomerID, "associationType", keyType, "provisioningKey")
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []provisioningkey.ProvisioningKey{},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := provisioningkey.GetByName(context.Background(), service, keyType, "missing")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestProvisioningKey_GetByIDAllAssociations_AllNotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	keyID := "missing-key"
	for _, assocType := range []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"} {
		path := common.ZPAPath(api.CustomerID, "associationType", assocType, "provisioningKey", keyID)
		api.On("GET", path, common.NotFoundResponse())
	}

	got, assoc, _, err := provisioningkey.GetByIDAllAssociations(context.Background(), api.Service, keyID)
	require.Error(t, err)
	assert.Nil(t, got)
	assert.Empty(t, assoc)
}

func TestProvisioningKey_GetByNameAllAssociations_AllNotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	for _, assocType := range []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"} {
		path := common.ZPAPath(api.CustomerID, "associationType", assocType, "provisioningKey")
		api.On("GET", path, common.SuccessResponse(common.ZPAList([]provisioningkey.ProvisioningKey{})))
	}

	got, assoc, _, err := provisioningkey.GetByNameAllAssociations(context.Background(), api.Service, "no-such-key")
	require.Error(t, err)
	assert.Nil(t, got)
	assert.Empty(t, assoc)
}

func TestProvisioningKey_GetAllByZComponentID_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	keyType := "CONNECTOR_GRP"
	zid := "z-missing"
	path := common.ZPAPath(api.CustomerID, "associationType", keyType, "zcomponent", zid, "provisioningKey")
	api.On("GET", path, common.NotFoundResponse())

	got, err := provisioningkey.GetAllByZComponentID(context.Background(), api.Service, keyType, zid)
	require.Error(t, err)
	assert.Nil(t, got)
}
