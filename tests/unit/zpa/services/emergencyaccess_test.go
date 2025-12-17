// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/emergencyaccess"
)

func TestEmergencyAccess_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user/" + userID

	server.On("GET", path, common.SuccessResponse(emergencyaccess.EmergencyAccess{
		UserID:    userID,
		EmailID:   "user@example.com",
		FirstName: "Test",
		LastName:  "User",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := emergencyaccess.Get(context.Background(), service, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
}

func TestEmergencyAccess_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user"

	server.On("POST", path, common.SuccessResponse(emergencyaccess.EmergencyAccess{
		UserID:    "new-user-123",
		EmailID:   "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newUser := &emergencyaccess.EmergencyAccess{
		EmailID:   "newuser@example.com",
		FirstName: "New",
		LastName:  "User",
	}

	result, _, err := emergencyaccess.Create(context.Background(), service, newUser)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-user-123", result.UserID)
}

func TestEmergencyAccess_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user/" + userID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateUser := &emergencyaccess.EmergencyAccess{
		UserID:    userID,
		FirstName: "Updated",
		LastName:  "User",
	}

	resp, err := emergencyaccess.Update(context.Background(), service, userID, updateUser)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEmergencyAccess_Deactivate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user/" + userID + "/deactivate"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := emergencyaccess.Deactivate(context.Background(), service, userID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEmergencyAccess_Activate_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/user/" + userID + "/activate"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := emergencyaccess.Activate(context.Background(), service, userID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestEmergencyAccess_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// GetAll uses the /emergencyAccess/users endpoint (plural) with pagination
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/users"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"items": []emergencyaccess.EmergencyAccess{
			{UserID: "user-001", EmailID: "user1@example.com"},
			{UserID: "user-002", EmailID: "user2@example.com"},
		},
		"nextPage": "",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := emergencyaccess.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestEmergencyAccess_GetByEmailID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	emailID := "user@example.com"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/users"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"items": []emergencyaccess.EmergencyAccess{
			{UserID: "user-001", EmailID: "other@example.com"},
			{UserID: "user-002", EmailID: emailID},
		},
		"nextPage": "",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := emergencyaccess.GetByEmailID(context.Background(), service, emailID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "user-002", result.UserID)
	assert.Equal(t, emailID, result.EmailID)
}

func TestEmergencyAccess_GetByEmailID_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/emergencyAccess/users"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"items": []emergencyaccess.EmergencyAccess{
			{UserID: "user-001", EmailID: "other@example.com"},
		},
		"nextPage": "",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := emergencyaccess.GetByEmailID(context.Background(), service, "notfound@example.com")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no emergency access record found")
}
