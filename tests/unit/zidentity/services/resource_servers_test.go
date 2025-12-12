// Package services provides unit tests for ZIdentity services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcommon "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	resourceservers "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/resource_servers"
)

func TestResourceServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ResourceServers JSON marshaling", func(t *testing.T) {
		server := resourceservers.ResourceServers{
			ID:          "rs-123",
			Name:        "ZPA API",
			DisplayName: "Zscaler Private Access API",
			Description: "API for managing ZPA resources",
			PrimaryAud:  "https://api.zscaler.com",
			DefaultApi:  true,
			ServiceScopes: []resourceservers.ServiceScopes{
				{
					Service: resourceservers.Service{
						ID:          "svc-001",
						Name:        "zpa",
						DisplayName: "Zscaler Private Access",
						CloudName:   "zscaler",
						OrgName:     "acme-corp",
					},
					Scopes: []resourceservers.Scopes{
						{ID: "scope-001", Name: "read:policies"},
						{ID: "scope-002", Name: "write:policies"},
					},
				},
			},
		}

		data, err := json.Marshal(server)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"rs-123"`)
		assert.Contains(t, string(data), `"name":"ZPA API"`)
		assert.Contains(t, string(data), `"primaryAud":"https://api.zscaler.com"`)
		assert.Contains(t, string(data), `"defaultApi":true`)
		assert.Contains(t, string(data), `"serviceScopes"`)
	})

	t.Run("ResourceServers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "rs-456",
			"name": "ZIA API",
			"displayName": "Zscaler Internet Access API",
			"description": "API for managing ZIA resources",
			"primaryAud": "https://api.zscalerbeta.net",
			"defaultApi": false,
			"serviceScopes": [
				{
					"service": {
						"id": "svc-zia",
						"name": "zia",
						"displayName": "Zscaler Internet Access",
						"cloudName": "zscalerbeta",
						"orgName": "example-org"
					},
					"scopes": [
						{"id": "scope-r1", "name": "read:rules"},
						{"id": "scope-w1", "name": "write:rules"},
						{"id": "scope-a1", "name": "admin:rules"}
					]
				}
			]
		}`

		var server resourceservers.ResourceServers
		err := json.Unmarshal([]byte(jsonData), &server)
		require.NoError(t, err)

		assert.Equal(t, "rs-456", server.ID)
		assert.Equal(t, "ZIA API", server.Name)
		assert.Equal(t, "Zscaler Internet Access API", server.DisplayName)
		assert.False(t, server.DefaultApi)
		assert.Len(t, server.ServiceScopes, 1)
		assert.Equal(t, "zia", server.ServiceScopes[0].Service.Name)
		assert.Len(t, server.ServiceScopes[0].Scopes, 3)
	})

	t.Run("ServiceScopes JSON marshaling", func(t *testing.T) {
		serviceScope := resourceservers.ServiceScopes{
			Service: resourceservers.Service{
				ID:          "svc-test",
				Name:        "test-service",
				DisplayName: "Test Service",
				CloudName:   "test-cloud",
				OrgName:     "test-org",
			},
			Scopes: []resourceservers.Scopes{
				{ID: "s1", Name: "read"},
				{ID: "s2", Name: "write"},
			},
		}

		data, err := json.Marshal(serviceScope)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"service"`)
		assert.Contains(t, string(data), `"scopes"`)
		assert.Contains(t, string(data), `"name":"test-service"`)
	})

	t.Run("Service JSON marshaling", func(t *testing.T) {
		service := resourceservers.Service{
			ID:          "svc-zdx",
			Name:        "zdx",
			DisplayName: "Zscaler Digital Experience",
			CloudName:   "zscaler",
			OrgName:     "enterprise-org",
		}

		data, err := json.Marshal(service)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"svc-zdx"`)
		assert.Contains(t, string(data), `"name":"zdx"`)
		assert.Contains(t, string(data), `"displayName":"Zscaler Digital Experience"`)
		assert.Contains(t, string(data), `"cloudName":"zscaler"`)
	})

	t.Run("Scopes JSON marshaling", func(t *testing.T) {
		scope := resourceservers.Scopes{
			ID:   "scope-admin",
			Name: "admin:all",
		}

		data, err := json.Marshal(scope)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"scope-admin"`)
		assert.Contains(t, string(data), `"name":"admin:all"`)
	})
}

func TestResourceServers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse resource servers list response", func(t *testing.T) {
		jsonResponse := `{
			"results_total": 3,
			"pageOffset": 0,
			"pageSize": 100,
			"next_link": "",
			"records": [
				{
					"id": "rs-001",
					"name": "ZPA API",
					"defaultApi": true
				},
				{
					"id": "rs-002",
					"name": "ZIA API",
					"defaultApi": false
				},
				{
					"id": "rs-003",
					"name": "ZDX API",
					"defaultApi": false
				}
			]
		}`

		var response common.PaginationResponse[resourceservers.ResourceServers]
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.ResultsTotal)
		assert.Len(t, response.Records, 3)
		assert.Equal(t, "ZPA API", response.Records[0].Name)
		assert.True(t, response.Records[0].DefaultApi)
		assert.False(t, response.Records[1].DefaultApi)
	})

	t.Run("Parse resource server with multiple service scopes", func(t *testing.T) {
		jsonResponse := `{
			"id": "rs-multi",
			"name": "Multi-Service API",
			"displayName": "Multi-Service Resource Server",
			"primaryAud": "https://api.example.com",
			"defaultApi": true,
			"serviceScopes": [
				{
					"service": {
						"id": "svc-1",
						"name": "service-a",
						"displayName": "Service A"
					},
					"scopes": [
						{"id": "s1", "name": "scope-a-read"},
						{"id": "s2", "name": "scope-a-write"}
					]
				},
				{
					"service": {
						"id": "svc-2",
						"name": "service-b",
						"displayName": "Service B"
					},
					"scopes": [
						{"id": "s3", "name": "scope-b-read"}
					]
				}
			]
		}`

		var server resourceservers.ResourceServers
		err := json.Unmarshal([]byte(jsonResponse), &server)
		require.NoError(t, err)

		assert.Equal(t, "rs-multi", server.ID)
		assert.Len(t, server.ServiceScopes, 2)
		assert.Equal(t, "service-a", server.ServiceScopes[0].Service.Name)
		assert.Len(t, server.ServiceScopes[0].Scopes, 2)
		assert.Equal(t, "service-b", server.ServiceScopes[1].Service.Name)
		assert.Len(t, server.ServiceScopes[1].Scopes, 1)
	})
}

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestResourceServers_Get_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	resourceID := "rs-12345"
	path := "/admin/api/v1/resource-servers/" + resourceID

	server.On("GET", path, testcommon.SuccessResponse(resourceservers.ResourceServers{
		ID:          resourceID,
		Name:        "ZPA API",
		DisplayName: "Zscaler Private Access API",
		Description: "API for managing ZPA resources",
		PrimaryAud:  "https://api.zscaler.com",
		DefaultApi:  true,
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := resourceservers.Get(context.Background(), service, resourceID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, resourceID, result.ID)
	assert.Equal(t, "ZPA API", result.Name)
	assert.True(t, result.DefaultApi)
}

func TestResourceServers_GetAll_SDK(t *testing.T) {
	server := testcommon.NewTestServer()
	defer server.Close()

	path := "/admin/api/v1/resource-servers"

	server.On("GET", path, testcommon.SuccessResponse(common.PaginationResponse[resourceservers.ResourceServers]{
		ResultsTotal: 2,
		PageOffset:   0,
		PageSize:     100,
		Records: []resourceservers.ResourceServers{
			{ID: "rs-1", Name: "ZPA API", DefaultApi: true},
			{ID: "rs-2", Name: "ZIA API", DefaultApi: false},
		},
	}))

	service, err := testcommon.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	results, err := resourceservers.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, 2)
	assert.Equal(t, "ZPA API", results[0].Name)
}

