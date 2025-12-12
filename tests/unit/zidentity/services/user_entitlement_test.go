// Package services provides unit tests for ZIdentity services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/user_entitlement"
)

func TestUserEntitlement_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Entitlements JSON marshaling", func(t *testing.T) {
		entitlement := user_entitlement.Entitlements{
			Roles: []common.IDNameDisplayName{
				{ID: "role-001", Name: "Admin", DisplayName: "Administrator"},
				{ID: "role-002", Name: "ReadOnly", DisplayName: "Read Only User"},
			},
			Scope: common.IDNameDisplayName{
				ID:          "scope-001",
				Name:        "Global",
				DisplayName: "Global Scope",
			},
			Service: user_entitlement.Service{
				ID:              "svc-zpa",
				ServiceName:     "ZPA",
				CloudName:       "zscaler",
				CloudDomainName: "zscaler.com",
				OrgName:         "acme-corp",
				OrgID:           "org-123",
			},
		}

		data, err := json.Marshal(entitlement)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"roles"`)
		assert.Contains(t, string(data), `"scope"`)
		assert.Contains(t, string(data), `"service"`)
		assert.Contains(t, string(data), `"serviceName":"ZPA"`)
	})

	t.Run("Entitlements JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"roles": [
				{"id": "role-admin", "name": "SuperAdmin", "displayName": "Super Administrator"},
				{"id": "role-user", "name": "User", "displayName": "Standard User"}
			],
			"scope": {
				"id": "scope-global",
				"name": "GlobalScope",
				"displayName": "Global Administration Scope"
			},
			"service": {
				"id": "svc-zia",
				"serviceName": "ZIA",
				"cloudName": "zscalerbeta",
				"cloudDomainName": "zscalerbeta.net",
				"orgName": "enterprise",
				"orgId": "org-456"
			}
		}`

		var entitlement user_entitlement.Entitlements
		err := json.Unmarshal([]byte(jsonData), &entitlement)
		require.NoError(t, err)

		assert.Len(t, entitlement.Roles, 2)
		assert.Equal(t, "SuperAdmin", entitlement.Roles[0].Name)
		assert.Equal(t, "GlobalScope", entitlement.Scope.Name)
		assert.Equal(t, "ZIA", entitlement.Service.ServiceName)
		assert.Equal(t, "org-456", entitlement.Service.OrgID)
	})

	t.Run("Service JSON marshaling", func(t *testing.T) {
		service := user_entitlement.Service{
			ID:              "svc-zdx",
			ServiceName:     "ZDX",
			CloudName:       "zscaler",
			CloudDomainName: "zscaler.com",
			OrgName:         "example-org",
			OrgID:           "org-789",
		}

		data, err := json.Marshal(service)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"svc-zdx"`)
		assert.Contains(t, string(data), `"serviceName":"ZDX"`)
		assert.Contains(t, string(data), `"cloudName":"zscaler"`)
		assert.Contains(t, string(data), `"cloudDomainName":"zscaler.com"`)
		assert.Contains(t, string(data), `"orgName":"example-org"`)
		assert.Contains(t, string(data), `"orgId":"org-789"`)
	})

	t.Run("Scope JSON marshaling", func(t *testing.T) {
		scope := user_entitlement.Scope{
			Scope: []common.IDNameDisplayName{
				{ID: "scope-1", Name: "Scope1", DisplayName: "First Scope"},
				{ID: "scope-2", Name: "Scope2", DisplayName: "Second Scope"},
			},
		}

		data, err := json.Marshal(scope)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"scope"`)
		assert.Contains(t, string(data), `"Scope1"`)
		assert.Contains(t, string(data), `"Scope2"`)
	})
}

func TestUserEntitlement_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse admin entitlements response", func(t *testing.T) {
		jsonResponse := `[
			{
				"roles": [
					{"id": "role-1", "name": "Admin", "displayName": "Administrator"}
				],
				"scope": {
					"id": "scope-global",
					"name": "Global",
					"displayName": "Global Scope"
				},
				"service": {
					"id": "svc-zpa",
					"serviceName": "ZPA",
					"cloudName": "zscaler",
					"orgName": "acme"
				}
			},
			{
				"roles": [
					{"id": "role-2", "name": "ReadOnly", "displayName": "Read Only"}
				],
				"scope": {
					"id": "scope-limited",
					"name": "Limited",
					"displayName": "Limited Scope"
				},
				"service": {
					"id": "svc-zia",
					"serviceName": "ZIA",
					"cloudName": "zscaler",
					"orgName": "acme"
				}
			}
		]`

		var entitlements []user_entitlement.Entitlements
		err := json.Unmarshal([]byte(jsonResponse), &entitlements)
		require.NoError(t, err)

		assert.Len(t, entitlements, 2)
		assert.Equal(t, "ZPA", entitlements[0].Service.ServiceName)
		assert.Equal(t, "Admin", entitlements[0].Roles[0].Name)
		assert.Equal(t, "ZIA", entitlements[1].Service.ServiceName)
		assert.Equal(t, "ReadOnly", entitlements[1].Roles[0].Name)
	})

	t.Run("Parse service entitlements response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": "svc-zpa",
				"serviceName": "ZPA",
				"cloudName": "zscaler",
				"cloudDomainName": "zscaler.com",
				"orgName": "enterprise-corp",
				"orgId": "org-001"
			},
			{
				"id": "svc-zia",
				"serviceName": "ZIA",
				"cloudName": "zscaler",
				"cloudDomainName": "zscaler.com",
				"orgName": "enterprise-corp",
				"orgId": "org-001"
			},
			{
				"id": "svc-zdx",
				"serviceName": "ZDX",
				"cloudName": "zscaler",
				"cloudDomainName": "zscaler.com",
				"orgName": "enterprise-corp",
				"orgId": "org-001"
			}
		]`

		var services []user_entitlement.Service
		err := json.Unmarshal([]byte(jsonResponse), &services)
		require.NoError(t, err)

		assert.Len(t, services, 3)
		assert.Equal(t, "ZPA", services[0].ServiceName)
		assert.Equal(t, "ZIA", services[1].ServiceName)
		assert.Equal(t, "ZDX", services[2].ServiceName)
		
		// All services should belong to the same org
		for _, svc := range services {
			assert.Equal(t, "enterprise-corp", svc.OrgName)
			assert.Equal(t, "org-001", svc.OrgID)
		}
	})

	t.Run("Parse entitlement with multiple roles", func(t *testing.T) {
		jsonResponse := `{
			"roles": [
				{"id": "role-super", "name": "SuperAdmin", "displayName": "Super Administrator"},
				{"id": "role-admin", "name": "Admin", "displayName": "Administrator"},
				{"id": "role-policy", "name": "PolicyAdmin", "displayName": "Policy Administrator"},
				{"id": "role-audit", "name": "Auditor", "displayName": "Auditor"}
			],
			"scope": {
				"id": "scope-all",
				"name": "AllResources",
				"displayName": "All Resources"
			},
			"service": {
				"id": "svc-zpa",
				"serviceName": "ZPA",
				"cloudName": "zscaler"
			}
		}`

		var entitlement user_entitlement.Entitlements
		err := json.Unmarshal([]byte(jsonResponse), &entitlement)
		require.NoError(t, err)

		assert.Len(t, entitlement.Roles, 4)
		roleNames := make([]string, len(entitlement.Roles))
		for i, role := range entitlement.Roles {
			roleNames[i] = role.Name
		}
		assert.Contains(t, roleNames, "SuperAdmin")
		assert.Contains(t, roleNames, "PolicyAdmin")
		assert.Contains(t, roleNames, "Auditor")
	})
}

