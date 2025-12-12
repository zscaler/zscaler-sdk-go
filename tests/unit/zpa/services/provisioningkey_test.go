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

// Note: GetAll iterates through all ProvisioningKeyAssociationTypes
// (CONNECTOR_GRP, SERVICE_EDGE_GRP, SITE_CONTROLLER_GRP, NP_ASSISTANT_GRP)
// Testing GetAll would require mocking multiple paths.
// For comprehensive GetAll testing, use integration tests.
