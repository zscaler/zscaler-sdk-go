// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

func TestPRACredential_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	credID := "cred-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential/" + credID

	server.On("GET", path, common.SuccessResponse(pracredential.Credential{
		ID:   credID,
		Name: "Test Credential",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredential.Get(context.Background(), service, credID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, credID, result.ID)
}

func TestPRACredential_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []pracredential.Credential{{ID: "cred-001"}, {ID: "cred-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredential.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
