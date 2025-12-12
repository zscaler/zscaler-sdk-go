// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/adminuserrolemgmt/adminusers"
)

func TestAdminUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminUsers JSON marshaling", func(t *testing.T) {
		user := adminusers.AdminUsers{
			ID:                          1234,
			LoginName:                   "admin@company.com",
			UserName:                    "Admin User",
			Email:                       "admin@company.com",
			Comments:                    "Primary administrator",
			Disabled:                    false,
			IsNonEditable:               true,
			IsPasswordLoginAllowed:      true,
			IsPasswordExpired:           false,
			IsAuditor:                   false,
			IsSecurityReportCommEnabled: true,
			IsServiceUpdateCommEnabled:  true,
			IsProductUpdateCommEnabled:  true,
			IsExecMobileAppEnabled:      true,
			AdminScopeType:              "ORGANIZATION",
			Role: &adminusers.Role{
				ID:   1,
				Name: "Super Admin",
			},
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1234`)
		assert.Contains(t, string(data), `"loginName":"admin@company.com"`)
		assert.Contains(t, string(data), `"adminScopeType":"ORGANIZATION"`)
	})

	t.Run("AdminUsers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 5678,
			"loginName": "security@company.com",
			"userName": "Security Admin",
			"email": "security@company.com",
			"comments": "Security team administrator",
			"disabled": false,
			"isNonEditable": false,
			"isPasswordLoginAllowed": true,
			"isAuditor": false,
			"adminScopeType": "DEPARTMENT",
			"role": {
				"id": 2,
				"name": "Security Admin",
				"isNameL10nTag": false
			},
			"adminScopeScopeEntities": [
				{"id": 100, "name": "Engineering"},
				{"id": 101, "name": "Sales"}
			]
		}`

		var user adminusers.AdminUsers
		err := json.Unmarshal([]byte(jsonData), &user)
		require.NoError(t, err)

		assert.Equal(t, 5678, user.ID)
		assert.Equal(t, "security@company.com", user.LoginName)
		assert.Equal(t, "DEPARTMENT", user.AdminScopeType)
		assert.NotNil(t, user.Role)
		assert.Equal(t, "Security Admin", user.Role.Name)
		assert.Len(t, user.AdminScopeEntities, 2)
	})

	t.Run("Role JSON marshaling", func(t *testing.T) {
		role := adminusers.Role{
			ID:           1,
			Name:         "Custom Role",
			IsNameL10Tag: false,
			Extensions: map[string]interface{}{
				"customField": "value",
			},
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1`)
		assert.Contains(t, string(data), `"name":"Custom Role"`)
	})

	t.Run("ExecMobileAppTokens JSON marshaling", func(t *testing.T) {
		token := adminusers.ExecMobileAppTokens{
			Cloud:       "zscaler",
			OrgId:       12345,
			Name:        "Mobile Token",
			TokenId:     "token-abc-123",
			Token:       "secret-token-value",
			TokenExpiry: 1700000000,
			CreateTime:  1699000000,
			DeviceId:    "device-xyz",
			DeviceName:  "iPhone 15",
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"cloud":"zscaler"`)
		assert.Contains(t, string(data), `"tokenId":"token-abc-123"`)
		assert.Contains(t, string(data), `"deviceName":"iPhone 15"`)
	})
}

func TestAdminUsers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin users list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"loginName": "super@company.com",
				"userName": "Super Admin",
				"disabled": false,
				"isAuditor": false,
				"role": {"id": 1, "name": "Super Admin"}
			},
			{
				"id": 2,
				"loginName": "auditor@company.com",
				"userName": "Auditor User",
				"disabled": false,
				"isAuditor": true,
				"role": {"id": 3, "name": "Auditor"}
			}
		]`

		var users []adminusers.AdminUsers
		err := json.Unmarshal([]byte(jsonResponse), &users)
		require.NoError(t, err)

		assert.Len(t, users, 2)
		assert.Equal(t, "super@company.com", users[0].LoginName)
		assert.False(t, users[0].IsAuditor)
		assert.True(t, users[1].IsAuditor)
	})
}

