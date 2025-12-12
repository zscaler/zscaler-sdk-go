// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminuserrolemgmt/roles"
)

func TestAdminUsers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminUsers JSON marshaling", func(t *testing.T) {
		admin := admins.AdminUsers{
			ID:                     12345,
			LoginName:              "admin@company.com",
			UserName:               "Admin User",
			Email:                  "admin@company.com",
			Comments:               "Super admin",
			Disabled:               false,
			IsNonEditable:          false,
			IsPasswordLoginAllowed: true,
			IsAuditor:              false,
			AdminScopeType:         "ORGANIZATION",
			Role: &admins.Role{
				ID:   100,
				Name: "Super Admin",
			},
		}

		data, err := json.Marshal(admin)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"loginName":"admin@company.com"`)
		assert.Contains(t, string(data), `"adminScopeType":"ORGANIZATION"`)
	})

	t.Run("AdminUsers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"loginName": "auditor@company.com",
			"userName": "Auditor User",
			"email": "auditor@company.com",
			"isAuditor": true,
			"isPasswordLoginAllowed": true,
			"isPasswordExpired": false,
			"isSecurityReportCommEnabled": true,
			"isServiceUpdateCommEnabled": false,
			"isExecMobileAppEnabled": true,
			"adminScopeType": "DEPARTMENT",
			"adminScopeScopeEntities": [
				{"id": 1, "name": "Engineering"}
			],
			"role": {
				"id": 200,
				"name": "Auditor"
			},
			"execMobileAppTokens": [
				{
					"cloud": "zscloud",
					"orgId": 1000,
					"name": "Token1",
					"tokenId": "token-123"
				}
			]
		}`

		var admin admins.AdminUsers
		err := json.Unmarshal([]byte(jsonData), &admin)
		require.NoError(t, err)

		assert.Equal(t, 54321, admin.ID)
		assert.True(t, admin.IsAuditor)
		assert.Equal(t, "DEPARTMENT", admin.AdminScopeType)
		assert.Len(t, admin.AdminScopeEntities, 1)
		assert.NotNil(t, admin.Role)
		assert.Len(t, admin.ExecMobileAppTokens, 1)
	})

	t.Run("Role JSON marshaling", func(t *testing.T) {
		role := admins.Role{
			ID:           100,
			Name:         "Super Admin",
			IsNameL10Tag: false,
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":100`)
		assert.Contains(t, string(data), `"name":"Super Admin"`)
	})

	t.Run("ExecMobileAppTokens JSON marshaling", func(t *testing.T) {
		token := admins.ExecMobileAppTokens{
			Cloud:       "zscloud",
			OrgId:       12345,
			Name:        "Mobile Token",
			TokenId:     "token-abc-123",
			Token:       "secret-token",
			TokenExpiry: 1699000000,
			CreateTime:  1698000000,
			DeviceId:    "device-123",
			DeviceName:  "iPhone 15",
		}

		data, err := json.Marshal(token)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"cloud":"zscloud"`)
		assert.Contains(t, string(data), `"tokenId":"token-abc-123"`)
	})

	t.Run("PasswordExpiry JSON marshaling", func(t *testing.T) {
		expiry := admins.PasswordExpiry{
			PasswordExpirationEnabled: true,
			PasswordExpiryDays:        90,
		}

		data, err := json.Marshal(expiry)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"passwordExpirationEnabled":true`)
		assert.Contains(t, string(data), `"passwordExpiryDays":90`)
	})
}

func TestAdminRoles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminRoles JSON marshaling", func(t *testing.T) {
		role := roles.AdminRoles{
			ID:               12345,
			Rank:             7,
			Name:             "Super Admin",
			PolicyAccess:     "READ_WRITE",
			AlertingAccess:   "READ_WRITE",
			UsernameAccess:   "READ_WRITE",
			DeviceInfoAccess: "READ_WRITE",
			DashboardAccess:  "READ_WRITE",
			ReportAccess:     "READ_WRITE",
			AnalysisAccess:   "READ_WRITE",
			AdminAcctAccess:  "READ_WRITE",
			IsAuditor:        false,
			Permissions:      []string{"ADMIN_ACCOUNT_READ_WRITE", "POLICY_READ_WRITE"},
			IsNonEditable:    false,
			LogsLimit:        "UNRESTRICTED",
			RoleType:         "EXECUTIVE_INSIGHTS_BUSINESS_ADMIN",
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"policyAccess":"READ_WRITE"`)
		assert.Contains(t, string(data), `"permissions"`)
	})

	t.Run("AdminRoles JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"rank": 5,
			"name": "Security Admin",
			"policyAccess": "READ_WRITE",
			"alertingAccess": "READ_ONLY",
			"usernameAccess": "NONE",
			"deviceInfoAccess": "READ_ONLY",
			"dashboardAccess": "READ_WRITE",
			"reportAccess": "READ_WRITE",
			"analysisAccess": "READ_ONLY",
			"adminAcctAccess": "NONE",
			"isAuditor": false,
			"permissions": ["POLICY_READ_WRITE", "REPORT_READ"],
			"featurePermissions": {
				"FIREWALL": "READ_WRITE",
				"SSL_INSPECTION": "READ_ONLY"
			},
			"isNonEditable": false,
			"logsLimit": "LAST_30_DAYS",
			"roleType": "STANDARD"
		}`

		var role roles.AdminRoles
		err := json.Unmarshal([]byte(jsonData), &role)
		require.NoError(t, err)

		assert.Equal(t, 54321, role.ID)
		assert.Equal(t, 5, role.Rank)
		assert.Equal(t, "READ_WRITE", role.PolicyAccess)
		assert.Equal(t, "NONE", role.UsernameAccess)
		assert.Len(t, role.Permissions, 2)
		assert.NotNil(t, role.FeaturePermissions)
	})
}

func TestAdminUserRoleMgmt_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin users list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "loginName": "admin1@company.com", "userName": "Admin 1", "isAuditor": false},
			{"id": 2, "loginName": "admin2@company.com", "userName": "Admin 2", "isAuditor": false},
			{"id": 3, "loginName": "auditor@company.com", "userName": "Auditor", "isAuditor": true}
		]`

		var admins []admins.AdminUsers
		err := json.Unmarshal([]byte(jsonResponse), &admins)
		require.NoError(t, err)

		assert.Len(t, admins, 3)
		assert.True(t, admins[2].IsAuditor)
	})

	t.Run("Parse admin roles list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Super Admin", "policyAccess": "READ_WRITE", "permissions": ["ALL"]},
			{"id": 2, "name": "Read Only", "policyAccess": "READ_ONLY", "permissions": ["READ"]}
		]`

		var rolesList []roles.AdminRoles
		err := json.Unmarshal([]byte(jsonResponse), &rolesList)
		require.NoError(t, err)

		assert.Len(t, rolesList, 2)
		assert.Equal(t, "READ_ONLY", rolesList[1].PolicyAccess)
	})
}

