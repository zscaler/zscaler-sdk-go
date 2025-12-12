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

func TestPRACredential_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	credName := "Production Credential"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []pracredential.Credential{
			{ID: "cred-001", Name: "Other Credential"},
			{ID: "cred-002", Name: credName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := pracredential.GetByName(context.Background(), service, credName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "cred-002", result.ID)
	assert.Equal(t, credName, result.Name)
}

func TestPRACredential_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential"

	server.On("POST", path, common.SuccessResponse(pracredential.Credential{
		ID:   "new-cred-123",
		Name: "New Credential",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newCred := &pracredential.Credential{
		Name: "New Credential",
	}

	result, _, err := pracredential.Create(context.Background(), service, newCred)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-cred-123", result.ID)
}

func TestPRACredential_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	credID := "cred-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential/" + credID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateCred := &pracredential.Credential{
		ID:   credID,
		Name: "Updated Credential",
	}

	resp, err := pracredential.Update(context.Background(), service, credID, updateCred)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRACredential_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	credID := "cred-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/credential/" + credID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := pracredential.Delete(context.Background(), service, credID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
