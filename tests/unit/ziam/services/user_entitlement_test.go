// Package services provides unit tests for ZIdentity services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcommon "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ziam/services/user_entitlement"
)

func TestUserEntitlement_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Entitlements JSON marshaling", func(t *testing.T) {
		entitlement := user_entitlement.Entitlements{
			Roles: []common.IDNameDisplayName{
				{ID: "role-001", Name: "Admin", DisplayName: "Administrator"},
				{ID: "role-002", Name: "ReadOnly", DisplayName: "Read Only User"},
			},
			Scope: &common.IDNameDisplayName{
				ID:          "scope-001",
				Name:        "Global",
				DisplayName: "Global Scope",
			},
			Service: &user_entitlement.Service{
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
		require.NotNil(t, entitlement.Scope)
		assert.Equal(t, "GlobalScope", entitlement.Scope.Name)
		require.NotNil(t, entitlement.Service)
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

	// Scope is a pointer so that "no scope" is distinguishable from a scope
	// whose id happens to be empty. A ZIAM-service entitlement carries no
	// `scope` key at all, which as a value type was indistinguishable from
	// `{"id": ""}`.
	t.Run("Absent scope unmarshals to nil", func(t *testing.T) {
		jsonData := `{
			"roles": [
				{"id": "phi1fv30307oc", "name": "Users Admin"}
			],
			"service": {
				"id": "hhlm44rd307qf",
				"serviceName": "ZIAM",
				"orgId": "0"
			}
		}`

		var entitlement user_entitlement.Entitlements
		err := json.Unmarshal([]byte(jsonData), &entitlement)
		require.NoError(t, err)

		assert.Nil(t, entitlement.Scope)
		require.NotNil(t, entitlement.Service)
		assert.Equal(t, "ZIAM", entitlement.Service.ServiceName)
	})

	t.Run("ServiceEntitlement wraps a service object", func(t *testing.T) {
		entitlement := user_entitlement.ServiceEntitlement{
			Service: &user_entitlement.Service{ID: "svc-zpa", ServiceName: "ZPA"},
		}

		data, err := json.Marshal(entitlement)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"service":{`)
		assert.Contains(t, string(data), `"serviceName":"ZPA"`)
	})
}

func TestUserEntitlement_ResponseParsing(t *testing.T) {
	t.Parallel()

	// The payload below is a verbatim admin-entitlements response from a
	// tenant. It is worth pinning exactly because it exercises the shape that
	// broke the old value-typed Scope: the first entitlement is scoped, the
	// second — a ZIAM-service entitlement — has no `scope` key at all.
	t.Run("Parse admin entitlements response", func(t *testing.T) {
		jsonResponse := `[
			{
				"roles": [
					{
						"id": "hplm44rqv07h6",
						"name": "ZPA Administrator"
					}
				],
				"scope": {
					"id": "mplm44rqi07jb"
				},
				"service": {
					"id": "hhlm44rae07ib",
					"serviceName": "ZPA",
					"cloudName": "ZPABETA",
					"cloudDomainName": "zpabeta.net",
					"orgName": "William Guilherme - Internal",
					"orgId": "72058304855015424"
				}
			},
			{
				"roles": [
					{
						"id": "phi1fv30307oc",
						"name": "Users Admin"
					}
				],
				"service": {
					"id": "hhlm44rd307qf",
					"serviceName": "ZIAM",
					"orgId": "0"
				}
			}
		]`

		var entitlements []user_entitlement.Entitlements
		err := json.Unmarshal([]byte(jsonResponse), &entitlements)
		require.NoError(t, err)

		require.Len(t, entitlements, 2)

		require.NotNil(t, entitlements[0].Service)
		assert.Equal(t, "ZPA", entitlements[0].Service.ServiceName)
		assert.Equal(t, "zpabeta.net", entitlements[0].Service.CloudDomainName)
		assert.Equal(t, "72058304855015424", entitlements[0].Service.OrgID)
		assert.Equal(t, "ZPA Administrator", entitlements[0].Roles[0].Name)
		require.NotNil(t, entitlements[0].Scope)
		assert.Equal(t, "mplm44rqi07jb", entitlements[0].Scope.ID)

		require.NotNil(t, entitlements[1].Service)
		assert.Equal(t, "ZIAM", entitlements[1].Service.ServiceName)
		assert.Equal(t, "Users Admin", entitlements[1].Roles[0].Name)
		// No scope key, and therefore no scope — not an empty one.
		assert.Nil(t, entitlements[1].Scope)
		// Fields the element omits stay empty rather than being invented.
		assert.Empty(t, entitlements[1].Service.CloudName)
	})

	// Each element of the service-entitlements array wraps a `service` object,
	// per the published schema. Reading it as a bare array of services matched
	// no keys and silently produced one zero-valued Service per entitlement.
	t.Run("Parse service entitlements response", func(t *testing.T) {
		jsonResponse := `[
			{
				"service": {
					"id": "svc-zpa",
					"serviceName": "ZPA",
					"cloudName": "zscaler",
					"cloudDomainName": "zscaler.com",
					"orgName": "enterprise-corp",
					"orgId": "org-001"
				}
			},
			{
				"service": {
					"id": "svc-zia",
					"serviceName": "ZIA",
					"cloudName": "zscaler",
					"cloudDomainName": "zscaler.com",
					"orgName": "enterprise-corp",
					"orgId": "org-001"
				}
			},
			{
				"service": {
					"id": "svc-zdx",
					"serviceName": "ZDX",
					"cloudName": "zscaler",
					"cloudDomainName": "zscaler.com",
					"orgName": "enterprise-corp",
					"orgId": "org-001"
				}
			}
		]`

		var entitlements []user_entitlement.ServiceEntitlement
		err := json.Unmarshal([]byte(jsonResponse), &entitlements)
		require.NoError(t, err)

		require.Len(t, entitlements, 3)
		require.NotNil(t, entitlements[0].Service)
		assert.Equal(t, "ZPA", entitlements[0].Service.ServiceName)
		assert.Equal(t, "ZIA", entitlements[1].Service.ServiceName)
		assert.Equal(t, "ZDX", entitlements[2].Service.ServiceName)

		// All services should belong to the same org
		for _, entitlement := range entitlements {
			require.NotNil(t, entitlement.Service)
			assert.Equal(t, "enterprise-corp", entitlement.Service.OrgName)
			assert.Equal(t, "org-001", entitlement.Service.OrgID)
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

		require.Len(t, entitlement.Roles, 4)
		roleNames := make([]string, len(entitlement.Roles))
		for i, role := range entitlement.Roles {
			roleNames[i] = role.Name
		}
		assert.Contains(t, roleNames, "SuperAdmin")
		assert.Contains(t, roleNames, "PolicyAdmin")
		assert.Contains(t, roleNames, "Auditor")
	})
}

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestUserEntitlement_GetAdminEntitlement_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/ziam/admin/api/v1/users/" + userID + "/admin-entitlements"

	server.On("GET", path, testcommon.SuccessResponse([]user_entitlement.Entitlements{
		{
			Roles: []common.IDNameDisplayName{
				{ID: "role-1", Name: "Admin", DisplayName: "Administrator"},
			},
			Scope: &common.IDNameDisplayName{
				ID:   "scope-global",
				Name: "Global",
			},
			Service: &user_entitlement.Service{
				ID:          "svc-zpa",
				ServiceName: "ZPA",
				CloudName:   "zscaler",
			},
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	results, err := user_entitlement.GetAdminEntitlement(context.Background(), service, userID)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results, 1)
	require.NotNil(t, results[0].Service)
	assert.Equal(t, "ZPA", results[0].Service.ServiceName)
	assert.Equal(t, "Admin", results[0].Roles[0].Name)
}

func TestUserEntitlement_GetServiceEntitlement_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	userID := "user-12345"
	path := "/ziam/admin/api/v1/users/" + userID + "/service-entitlements"

	server.On("GET", path, testcommon.SuccessResponse([]user_entitlement.ServiceEntitlement{
		{Service: &user_entitlement.Service{ID: "svc-zpa", ServiceName: "ZPA", CloudName: "zscaler"}},
		{Service: &user_entitlement.Service{ID: "svc-zia", ServiceName: "ZIA", CloudName: "zscaler"}},
		{Service: &user_entitlement.Service{ID: "svc-zdx", ServiceName: "ZDX", CloudName: "zscaler"}},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	results, err := user_entitlement.GetServiceEntitlement(context.Background(), service, userID)
	require.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results, 3)
	for _, result := range results {
		require.NotNil(t, result.Service)
	}
	assert.Equal(t, "ZPA", results[0].Service.ServiceName)
	assert.Equal(t, "ZIA", results[1].Service.ServiceName)
	assert.Equal(t, "ZDX", results[2].Service.ServiceName)
}
