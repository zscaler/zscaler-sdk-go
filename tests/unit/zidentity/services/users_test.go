// Package services provides unit tests for ZIdentity services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/users"
)

func TestUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Users JSON marshaling", func(t *testing.T) {
		user := users.Users{
			ID:             "user-123",
			Source:         "OKTA",
			LoginName:      "john.doe@example.com",
			DisplayName:    "John Doe",
			FirstName:      "John",
			LastName:       "Doe",
			PrimaryEmail:   "john.doe@example.com",
			SecondaryEmail: "jdoe@personal.com",
			Status:         true,
			Department: &common.IDNameDisplayName{
				ID:          "dept-001",
				Name:        "Engineering",
				DisplayName: "Engineering Department",
			},
			IDP: &common.IDNameDisplayName{
				ID:          "idp-001",
				Name:        "Okta",
				DisplayName: "Okta Identity Provider",
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"user-123"`)
		assert.Contains(t, string(data), `"loginName":"john.doe@example.com"`)
		assert.Contains(t, string(data), `"displayName":"John Doe"`)
		assert.Contains(t, string(data), `"firstName":"John"`)
		assert.Contains(t, string(data), `"status":true`)
	})

	t.Run("Users JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "user-456",
			"source": "AZURE_AD",
			"loginName": "jane.smith@example.com",
			"displayName": "Jane Smith",
			"firstName": "Jane",
			"lastName": "Smith",
			"primaryEmail": "jane.smith@example.com",
			"status": true,
			"department": {
				"id": "dept-002",
				"name": "Sales",
				"displayName": "Sales Department"
			},
			"idp": {
				"id": "idp-002",
				"name": "Azure AD",
				"displayName": "Azure Active Directory"
			},
			"customAttrsInfo": {
				"employeeId": "12345",
				"costCenter": "CC-100"
			}
		}`

		var user users.Users
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, "user-456", user.ID)
		assert.Equal(t, "AZURE_AD", user.Source)
		assert.Equal(t, "jane.smith@example.com", user.LoginName)
		assert.Equal(t, "Jane Smith", user.DisplayName)
		assert.True(t, user.Status)
		assert.NotNil(t, user.Department)
		assert.Equal(t, "Sales", user.Department.Name)
		assert.NotNil(t, user.IDP)
		assert.Equal(t, "Azure AD", user.IDP.Name)
		assert.NotNil(t, user.CustomAttrsInfo)
		assert.Equal(t, "12345", user.CustomAttrsInfo["employeeId"])
	})

	t.Run("Users without optional fields", func(t *testing.T) {
		user := users.Users{
			ID:           "user-minimal",
			LoginName:    "minimal@example.com",
			DisplayName:  "Minimal User",
			PrimaryEmail: "minimal@example.com",
			Status:       true,
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"user-minimal"`)
		assert.NotContains(t, string(data), `"department"`)
		assert.NotContains(t, string(data), `"idp"`)
		assert.NotContains(t, string(data), `"secondaryEmail"`)
	})
}

func TestUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse users list response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 3,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "",
			"records": [
				{
					"id": "user-001",
					"loginName": "admin@example.com",
					"displayName": "Admin User",
					"status": true
				},
				{
					"id": "user-002",
					"loginName": "john.doe@example.com",
					"displayName": "John Doe",
					"status": true
				},
				{
					"id": "user-003",
					"loginName": "jane.smith@example.com",
					"displayName": "Jane Smith",
					"status": false
				}
			]
		}`

		var response common.PaginationResponse[users.Users]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.ResultsTotal)
		assert.Len(t, response.Records, 3)
		assert.Equal(t, "Admin User", response.Records[0].DisplayName)
		assert.True(t, response.Records[0].Status)
		assert.False(t, response.Records[2].Status)
	})

	t.Run("Parse single user with full details", func(t *testing.T) {
		jsonResponse := `{
			"id": "user-full",
			"source": "OKTA",
			"loginName": "full.user@example.com",
			"displayName": "Full User",
			"firstName": "Full",
			"lastName": "User",
			"primaryEmail": "full.user@example.com",
			"secondaryEmail": "fulluser@personal.com",
			"status": true,
			"department": {
				"id": "dept-eng",
				"name": "Engineering",
				"displayName": "Engineering Department"
			},
			"idp": {
				"id": "idp-okta",
				"name": "Okta",
				"displayName": "Corporate Okta"
			},
			"customAttrsInfo": {
				"manager": "manager@example.com",
				"location": "San Jose",
				"team": "Platform"
			}
		}`

		var user users.Users
		err := json.Unmarshal([]byte(jsonResponse), &user)
		require.NoError(t, err)

		assert.Equal(t, "user-full", user.ID)
		assert.Equal(t, "Full User", user.DisplayName)
		assert.Equal(t, "fulluser@personal.com", user.SecondaryEmail)
		assert.NotNil(t, user.Department)
		assert.Equal(t, "Engineering", user.Department.Name)
		assert.NotNil(t, user.CustomAttrsInfo)
		assert.Equal(t, "San Jose", user.CustomAttrsInfo["location"])
	})

	t.Run("Parse paginated users response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 500,
			"pageOffset": 200,
			"pageSize": 100,
			"next_link": "/admin/api/v1/users?offset=300&limit=100",
			"prev_link": "/admin/api/v1/users?offset=100&limit=100",
			"records": [
				{"id": "user-201", "displayName": "User 201"},
				{"id": "user-202", "displayName": "User 202"}
			]
		}`

		var response common.PaginationResponse[users.Users]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 500, response.ResultsTotal)
		assert.Equal(t, 200, response.PageOffset)
		assert.NotEmpty(t, response.NextLink)
		assert.NotEmpty(t, response.PrevLink)
		assert.Len(t, response.Records, 2)
	})
}

