// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

func TestPRAApproval_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	approvalID := "approval-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval/" + approvalID

	server.On("GET", path, common.SuccessResponse(praapproval.PrivilegedApproval{
		ID: approvalID,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praapproval.Get(context.Background(), service, approvalID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, approvalID, result.ID)
}

func TestPRAApproval_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []praapproval.PrivilegedApproval{{ID: "approval-001"}, {ID: "approval-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praapproval.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestPRAApproval_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval"

	server.On("POST", path, common.SuccessResponse(praapproval.PrivilegedApproval{
		ID: "new-approval-123",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newApproval := &praapproval.PrivilegedApproval{}

	result, _, err := praapproval.Create(context.Background(), service, newApproval)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-approval-123", result.ID)
}

func TestPRAApproval_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	approvalID := "approval-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval/" + approvalID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateApproval := &praapproval.PrivilegedApproval{
		ID: approvalID,
	}

	resp, err := praapproval.Update(context.Background(), service, approvalID, updateApproval)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRAApproval_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	approvalID := "approval-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval/" + approvalID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := praapproval.Delete(context.Background(), service, approvalID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRAApproval_GetByEmailID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	emailID := "user@example.com"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []praapproval.PrivilegedApproval{
			{ID: "approval-001", EmailIDs: []string{"other@example.com"}},
			{ID: "approval-002", EmailIDs: []string{emailID}},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praapproval.GetByEmailID(context.Background(), service, emailID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "approval-002", result.ID)
}

func TestPRAApproval_DeleteExpired_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/approval/expired"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := praapproval.DeleteExpired(context.Background(), service)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
