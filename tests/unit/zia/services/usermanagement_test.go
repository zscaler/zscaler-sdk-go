// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/usermanagement/users"
)

func TestUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Users JSON marshaling", func(t *testing.T) {
		user := users.Users{
			ID:            12345,
			Name:          "John Doe",
			Email:         "john.doe@company.com",
			Comments:      "Engineering team lead",
			TempAuthEmail: "temp@company.com",
			AuthMethods:   []string{"BASIC", "DIGEST"},
			AdminUser:     false,
			Type:          "STANDARD",
			Deleted:       false,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"John Doe"`)
		assert.Contains(t, string(data), `"email":"john.doe@company.com"`)
	})

	t.Run("Users JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Jane Smith",
			"email": "jane.smith@company.com",
			"groups": [
				{"id": 100, "name": "Engineering", "idp_id": 1}
			],
			"department": {
				"id": 200,
				"name": "Product",
				"idp_id": 1
			},
			"comments": "Product manager",
			"authMethods": ["BASIC"],
			"adminUser": false,
			"type": "STANDARD",
			"deleted": false
		}`

		var user users.Users
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 54321, user.ID)
		assert.Equal(t, "Jane Smith", user.Name)
		assert.Len(t, user.Groups, 1)
		assert.NotNil(t, user.Department)
		assert.Equal(t, "Product", user.Department.Name)
	})

	t.Run("EnrollResult JSON marshaling", func(t *testing.T) {
		result := users.EnrollResult{
			UserID: 12345,
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"userId":12345`)
	})

	t.Run("EnrollUserRequest JSON marshaling", func(t *testing.T) {
		request := users.EnrollUserRequest{
			AuthMethods: []string{"BASIC", "DIGEST"},
			Password:    "securePassword123",
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"authMethods"`)
		assert.Contains(t, string(data), `"password":"securePassword123"`)
	})

	t.Run("GetAllUsersFilterOptions structure", func(t *testing.T) {
		opts := users.GetAllUsersFilterOptions{
			Name:  "John",
			Dept:  "Engineering",
			Group: "Developers",
		}

		assert.Equal(t, "John", opts.Name)
		assert.Equal(t, "Engineering", opts.Dept)
		assert.Equal(t, "Developers", opts.Group)
	})
}

func TestUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse users list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "User 1", "email": "user1@company.com", "adminUser": false},
			{"id": 2, "name": "User 2", "email": "user2@company.com", "adminUser": false},
			{"id": 3, "name": "Admin User", "email": "admin@company.com", "adminUser": true}
		]`

		var usersList []users.Users
		err := json.Unmarshal([]byte(jsonResponse), &usersList)
		require.NoError(t, err)

		assert.Len(t, usersList, 3)
		assert.True(t, usersList[2].AdminUser)
	})

	t.Run("Parse user with groups and department", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Full User",
			"email": "full@company.com",
			"groups": [
				{"id": 1, "name": "Group 1"},
				{"id": 2, "name": "Group 2"}
			],
			"department": {
				"id": 10,
				"name": "Engineering"
			}
		}`

		var user users.Users
		err := json.Unmarshal([]byte(jsonResponse), &user)
		require.NoError(t, err)

		assert.Len(t, user.Groups, 2)
		assert.NotNil(t, user.Department)
	})

	t.Run("Parse auditors list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Auditor 1", "adminUser": false, "type": "AUDITOR"},
			{"id": 2, "name": "Auditor 2", "adminUser": false, "type": "AUDITOR"}
		]`

		var auditors []users.Users
		err := json.Unmarshal([]byte(jsonResponse), &auditors)
		require.NoError(t, err)

		assert.Len(t, auditors, 2)
		assert.Equal(t, "AUDITOR", auditors[0].Type)
	})
}

