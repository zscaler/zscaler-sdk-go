// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/admin_users"
)

func TestAdminUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminUser JSON marshaling", func(t *testing.T) {
		user := admin_users.AdminUser{
			ID:             123,
			UserName:       "admin@example.com",
			CompanyID:      "company-456",
			AccountEnabled: "true",
			EditEnabled:    "true",
			IsDefaultAdmin: "false",
			ServiceType:    "ZCC",
			CompanyRole: admin_users.Role{
				ID:              "role-789",
				RoleName:        "Administrator",
				AdminManagement: "FULL",
				Dashboard:       "FULL",
				DeviceOverview:  "FULL",
				AuditLogs:       "READ",
				IsEditable:      true,
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"userName":"admin@example.com"`)
		assert.Contains(t, string(data), `"accountEnabled":"true"`)
		assert.Contains(t, string(data), `"roleName":"Administrator"`)
	})

	t.Run("AdminUser JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 123,
			"userName": "admin@example.com",
			"companyId": "company-456",
			"accountEnabled": "true",
			"editEnabled": "true",
			"isDefaultAdmin": "false",
			"serviceType": "ZCC",
			"companyRole": {
				"id": "role-789",
				"roleName": "Administrator",
				"adminManagement": "FULL",
				"dashboard": "FULL",
				"deviceOverview": "FULL",
				"auditLogs": "READ",
				"isEditable": true
			}
		}`

		var user admin_users.AdminUser
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 123, user.ID)
		assert.Equal(t, "admin@example.com", user.UserName)
		assert.Equal(t, "company-456", user.CompanyID)
		assert.Equal(t, "true", user.AccountEnabled)
		assert.Equal(t, "Administrator", user.CompanyRole.RoleName)
		assert.True(t, user.CompanyRole.IsEditable)
	})

	t.Run("Role JSON marshaling", func(t *testing.T) {
		role := admin_users.Role{
			ID:                           "role-123",
			RoleName:                     "Custom Admin",
			AdminManagement:              "FULL",
			Dashboard:                    "FULL",
			DeviceOverview:               "READ",
			DeviceGroups:                 "FULL",
			AuditLogs:                    "READ",
			AuthSetting:                  "NONE",
			TrustedNetwork:               "FULL",
			ForwardingProfile:            "READ",
			ClientConnectorAppStore:      "NONE",
			ClientConnectorIDP:           "READ",
			ClientConnectorSupport:       "FULL",
			ClientConnectorNotifications: "FULL",
			IsEditable:                   true,
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"role-123"`)
		assert.Contains(t, string(data), `"roleName":"Custom Admin"`)
		assert.Contains(t, string(data), `"adminManagement":"FULL"`)
		assert.Contains(t, string(data), `"isEditable":true`)
	})

	t.Run("SyncZiaZdxZpaAdminUsers JSON marshaling", func(t *testing.T) {
		sync := admin_users.SyncZiaZdxZpaAdminUsers{
			CompanyIDs:   []int{1, 2, 3},
			ErrorCode:    "",
			ErrorMessage: "",
			Success:      "true",
		}

		data, err := json.Marshal(sync)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"companyIds":[1,2,3]`)
		assert.Contains(t, string(data), `"success":"true"`)
	})
}

func TestAdminUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin users list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"userName": "admin1@example.com",
				"companyId": "company-1",
				"accountEnabled": "true",
				"companyRole": {
					"id": "role-1",
					"roleName": "Super Admin"
				}
			},
			{
				"id": 2,
				"userName": "admin2@example.com",
				"companyId": "company-1",
				"accountEnabled": "true",
				"companyRole": {
					"id": "role-2",
					"roleName": "Read Only"
				}
			}
		]`

		var users []admin_users.AdminUser
		err := json.Unmarshal([]byte(jsonResponse), &users)
		require.NoError(t, err)

		assert.Len(t, users, 2)
		assert.Equal(t, "admin1@example.com", users[0].UserName)
		assert.Equal(t, "Super Admin", users[0].CompanyRole.RoleName)
		assert.Equal(t, "admin2@example.com", users[1].UserName)
		assert.Equal(t, "Read Only", users[1].CompanyRole.RoleName)
	})

	t.Run("Parse sync response", func(t *testing.T) {
		jsonResponse := `{
			"companyIds": [100, 200],
			"errorCode": "",
			"errorMessage": "",
			"success": "true",
			"responseData": {}
		}`

		var sync admin_users.SyncZiaZdxZpaAdminUsers
		err := json.Unmarshal([]byte(jsonResponse), &sync)
		require.NoError(t, err)

		assert.Equal(t, []int{100, 200}, sync.CompanyIDs)
		assert.Equal(t, "true", sync.Success)
		assert.Empty(t, sync.ErrorCode)
	})
}

