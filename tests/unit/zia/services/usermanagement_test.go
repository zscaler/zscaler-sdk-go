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
