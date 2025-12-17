// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestUsers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := 12345
	path := "/zia/api/v1/users/12345"

	server.On("GET", path, common.SuccessResponse(users.Users{
		ID:       userID,
		Name:     "John Doe",
		Email:    "john.doe@company.com",
		Comments: "Test user",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := users.Get(context.Background(), service, userID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.ID)
	assert.Equal(t, "John Doe", result.Name)
}

func TestUsers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users"

	server.On("GET", path, common.SuccessResponse([]users.Users{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := users.GetAllUsers(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUsers_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users"

	server.On("POST", path, common.SuccessResponse(users.Users{
		ID:    99999,
		Name:  "New User",
		Email: "new@example.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newUser := users.Users{
		Name:  "New User",
		Email: "new@example.com",
	}

	result, err := users.Create(context.Background(), service, &newUser)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestUsers_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := 12345
	path := "/zia/api/v1/users/12345"

	server.On("PUT", path, common.SuccessResponse(users.Users{
		ID:   userID,
		Name: "Updated User",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateUser := users.Users{
		ID:   userID,
		Name: "Updated User",
	}

	result, _, err := users.Update(context.Background(), service, userID, &updateUser)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated User", result.Name)
}

func TestUsers_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := 12345
	path := "/zia/api/v1/users/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = users.Delete(context.Background(), service, userID)

	require.NoError(t, err)
}

func TestUsers_GetUserByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userName := "John Doe"
	path := "/zia/api/v1/users"

	server.On("GET", path, common.SuccessResponse([]users.Users{
		{ID: 1, Name: "Other User", Email: "other@example.com"},
		{ID: 2, Name: userName, Email: "john.doe@company.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := users.GetUserByName(context.Background(), service, userName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userName, result.Name)
}

func TestUsers_EnrollUser_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userID := 12345
	path := "/zia/api/v1/users/12345/zia/api/v1/enroll"

	server.On("POST", path, common.SuccessResponse(users.EnrollResult{
		UserID: userID,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	request := users.EnrollUserRequest{
		AuthMethods: []string{"BASIC"},
		Password:    "SecureP@ss123",
	}

	result, err := users.EnrollUser(context.Background(), service, userID, request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
}

func TestUsers_EnrollUser_InvalidAuthMethod(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	request := users.EnrollUserRequest{
		AuthMethods: []string{"INVALID_METHOD"},
		Password:    "SecureP@ss123",
	}

	_, err = users.EnrollUser(context.Background(), service, 12345, request)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "authMethods must be one of the following")
}

func TestUsers_BulkDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users/bulkDelete"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ids := []int{1, 2, 3}
	_, err = users.BulkDelete(context.Background(), service, ids)

	require.NoError(t, err)
}

func TestUsers_GetAllAuditors_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users/auditors"

	server.On("GET", path, common.SuccessResponse([]users.Users{
		{ID: 1, Name: "Auditor 1", Email: "auditor1@example.com", Type: "AUDITOR"},
		{ID: 2, Name: "Auditor 2", Email: "auditor2@example.com", Type: "AUDITOR"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := users.GetAllAuditors(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUsers_GetUserReferences_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users/references"

	server.On("GET", path, common.SuccessResponse([]struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		ExternalID string `json:"externalId"`
	}{
		{ID: 1, Name: "User 1", ExternalID: "ext-1"},
		{ID: 2, Name: "User 2", ExternalID: "ext-2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := users.GetUserReferences(context.Background(), service, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUsers_GetAllWithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/users"

	server.On("GET", path, common.SuccessResponse([]users.Users{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &users.GetAllUsersFilterOptions{
		Name:  "John",
		Dept:  "Engineering",
		Group: "Admins",
	}

	result, err := users.GetAllUsers(context.Background(), service, opts)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Users JSON marshaling", func(t *testing.T) {
		user := users.Users{
			ID:       12345,
			Name:     "John Doe",
			Email:    "john.doe@company.com",
			Comments: "Engineering team lead",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"John Doe"`)
	})

	t.Run("Users JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Jane Smith",
			"email": "jane.smith@company.com",
			"groups": [{"id": 100, "name": "Engineering"}],
			"department": {"id": 200, "name": "Product"}
		}`

		var user users.Users
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 54321, user.ID)
		assert.Equal(t, "Jane Smith", user.Name)
	})
}
