// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/admin_roles"
)

func TestAdminRoles_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AdminRole JSON marshaling", func(t *testing.T) {
		role := admin_roles.AdminRole{
			ID:                           "role-123",
			RoleName:                     "Super Administrator",
			CompanyID:                    "company-456",
			AdminManagement:              "FULL",
			Dashboard:                    "FULL",
			DeviceOverview:               "FULL",
			DeviceGroups:                 "FULL",
			AuditLogs:                    "FULL",
			AuthSetting:                  "FULL",
			TrustedNetwork:               "FULL",
			ForwardingProfile:            "FULL",
			ClientConnectorAppStore:      "FULL",
			ClientConnectorIDP:           "FULL",
			ClientConnectorSupport:       "FULL",
			ClientConnectorNotifications: "FULL",
			IsEditable:                   false,
			CreatedBy:                    "system",
			UpdatedBy:                    "admin@example.com",
		}

		data, err := json.Marshal(role)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"role-123"`)
		assert.Contains(t, string(data), `"roleName":"Super Administrator"`)
		assert.Contains(t, string(data), `"adminManagement":"FULL"`)
		assert.Contains(t, string(data), `"isEditable":false`)
	})

	t.Run("AdminRole JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "role-789",
			"roleName": "Custom Admin",
			"companyId": "company-456",
			"adminManagement": "READ",
			"dashboard": "FULL",
			"deviceOverview": "READ",
			"deviceGroups": "NONE",
			"auditLogs": "READ",
			"authSetting": "NONE",
			"trustedNetwork": "FULL",
			"forwardingProfile": "READ",
			"clientConnectorAppStore": "NONE",
			"clientConnectorIdp": "READ",
			"clientConnectorSupport": "FULL",
			"clientConnectorNotifications": "FULL",
			"isEditable": true,
			"createdBy": "admin@example.com",
			"updatedBy": "admin@example.com"
		}`

		var role admin_roles.AdminRole
		err := json.Unmarshal([]byte(jsonData), &role)
		require.NoError(t, err)

		assert.Equal(t, "role-789", role.ID)
		assert.Equal(t, "Custom Admin", role.RoleName)
		assert.Equal(t, "READ", role.AdminManagement)
		assert.Equal(t, "FULL", role.Dashboard)
		assert.Equal(t, "NONE", role.DeviceGroups)
		assert.True(t, role.IsEditable)
	})
}

func TestAdminRoles_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin roles list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "role-001",
				"roleName": "Super Administrator",
				"isEditable": false,
				"adminManagement": "FULL",
				"dashboard": "FULL"
			},
			{
				"id": "role-002",
				"roleName": "Read Only Admin",
				"isEditable": true,
				"adminManagement": "NONE",
				"dashboard": "READ"
			},
			{
				"id": "role-003",
				"roleName": "Device Manager",
				"isEditable": true,
				"adminManagement": "NONE",
				"dashboard": "READ",
				"deviceOverview": "FULL",
				"deviceGroups": "FULL"
			}
		]`

		var roles []admin_roles.AdminRole
		err := json.Unmarshal([]byte(jsonResponse), &roles)
		require.NoError(t, err)

		assert.Len(t, roles, 3)
		assert.Equal(t, "Super Administrator", roles[0].RoleName)
		assert.False(t, roles[0].IsEditable)
		assert.Equal(t, "Read Only Admin", roles[1].RoleName)
		assert.True(t, roles[1].IsEditable)
		assert.Equal(t, "Device Manager", roles[2].RoleName)
		assert.Equal(t, "FULL", roles[2].DeviceGroups)
	})
}

