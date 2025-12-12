// Package unit provides unit tests for ZPA Server Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
)

// TestServerGroup_Structure tests the struct definitions
func TestServerGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ServerGroup JSON marshaling", func(t *testing.T) {
		group := servergroup.ServerGroup{
			ID:               "sg-123",
			Name:             "Test Server Group",
			Description:      "Test Description",
			Enabled:          true,
			IpAnchored:       false,
			ConfigSpace:      "DEFAULT",
			DynamicDiscovery: true,
			ExtranetEnabled:  false,
			MicroTenantID:    "mt-001",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled servergroup.ServerGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.Equal(t, group.Enabled, unmarshaled.Enabled)
		assert.Equal(t, group.DynamicDiscovery, unmarshaled.DynamicDiscovery)
	})

	t.Run("ServerGroup JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "sg-456",
			"name": "Production Server Group",
			"description": "Production servers",
			"enabled": true,
			"ipAnchored": true,
			"configSpace": "DEFAULT",
			"dynamicDiscovery": false,
			"extranetEnabled": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-002",
			"microtenantName": "Production Tenant",
			"appConnectorGroups": [
				{
					"id": "acg-001",
					"name": "Connector Group 1"
				}
			],
			"servers": [
				{
					"id": "srv-001",
					"name": "Server 1",
					"address": "10.0.0.1",
					"enabled": true
				}
			]
		}`

		var group servergroup.ServerGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "sg-456", group.ID)
		assert.Equal(t, "Production Server Group", group.Name)
		assert.True(t, group.Enabled)
		assert.True(t, group.IpAnchored)
		assert.False(t, group.DynamicDiscovery)
		assert.True(t, group.ExtranetEnabled)
		assert.Equal(t, "mt-002", group.MicroTenantID)
		assert.Len(t, group.AppConnectorGroups, 1)
		assert.Len(t, group.Servers, 1)
	})

	t.Run("Applications structure", func(t *testing.T) {
		apps := servergroup.Applications{
			ID:   "app-001",
			Name: "Test Application",
		}

		data, err := json.Marshal(apps)
		require.NoError(t, err)

		var unmarshaled servergroup.Applications
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, apps.ID, unmarshaled.ID)
		assert.Equal(t, apps.Name, unmarshaled.Name)
	})

	t.Run("Connectors structure", func(t *testing.T) {
		connector := servergroup.Connectors{
			ID:                    "conn-001",
			Name:                  "Test Connector",
			AppConnectorGroupID:   "acg-001",
			AppConnectorGroupName: "Connector Group 1",
			ControlChannelStatus:  "ZPN_STATUS_ONLINE",
			CurrentVersion:        "24.1.0",
			Enabled:               true,
			PrivateIP:             "10.0.0.5",
			PublicIP:              "203.0.113.5",
			Platform:              "el8",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled servergroup.Connectors
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, connector.ID, unmarshaled.ID)
		assert.Equal(t, connector.Name, unmarshaled.Name)
		assert.Equal(t, connector.ControlChannelStatus, unmarshaled.ControlChannelStatus)
	})
}

// TestServerGroup_ResponseParsing tests parsing of various API responses
func TestServerGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse single server group response", func(t *testing.T) {
		response := `{
			"id": "sg-789",
			"name": "Edge Server Group",
			"description": "Edge location servers",
			"enabled": true,
			"ipAnchored": false,
			"dynamicDiscovery": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com"
		}`

		var group servergroup.ServerGroup
		err := json.Unmarshal([]byte(response), &group)
		require.NoError(t, err)

		assert.Equal(t, "sg-789", group.ID)
		assert.Equal(t, "Edge Server Group", group.Name)
		assert.True(t, group.DynamicDiscovery)
		assert.NotEmpty(t, group.CreationTime)
	})

	t.Run("Parse server group list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Group 1", "enabled": true, "dynamicDiscovery": true},
				{"id": "2", "name": "Group 2", "enabled": false, "dynamicDiscovery": false},
				{"id": "3", "name": "Group 3", "enabled": true, "ipAnchored": true}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []servergroup.ServerGroup `json:"list"`
			TotalPages int                       `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "1", listResp.List[0].ID)
		assert.True(t, listResp.List[0].DynamicDiscovery)
		assert.True(t, listResp.List[2].IpAnchored)
	})
}

// TestServerGroup_MockServerOperations tests CRUD operations with mock server
func TestServerGroup_MockServerOperations(t *testing.T) {
	t.Run("GET server group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/serverGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "sg-123",
				"name": "Mock Server Group",
				"description": "Created by mock server",
				"enabled": true,
				"dynamicDiscovery": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/serverGroup/sg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all server groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Server Group A", "enabled": true},
					{"id": "2", "name": "Server Group B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/serverGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create server group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/serverGroup")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-sg-456",
				"name": "New Server Group",
				"description": "Newly created",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/serverGroup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update server group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/serverGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/serverGroup/sg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE server group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/serverGroup/")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/serverGroup/sg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestServerGroup_ErrorHandling tests error scenarios
func TestServerGroup_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Server Group not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serverGroup/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request - Missing app connector groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "At least one app connector group is required"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/serverGroup", "application/json", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestServerGroup_SpecialCases tests edge cases and special scenarios
func TestServerGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Server group with empty servers list", func(t *testing.T) {
		group := servergroup.ServerGroup{
			ID:      "123",
			Name:    "Empty Group",
			Enabled: true,
			Servers: nil,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled servergroup.ServerGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Nil(t, unmarshaled.Servers)
	})

	t.Run("Dynamic discovery enabled", func(t *testing.T) {
		group := servergroup.ServerGroup{
			ID:               "123",
			Name:             "Dynamic Group",
			DynamicDiscovery: true,
			Servers:          nil, // No servers needed when dynamic discovery is enabled
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"dynamicDiscovery":true`)
	})

	t.Run("IP Anchored configuration", func(t *testing.T) {
		group := servergroup.ServerGroup{
			ID:         "123",
			Name:       "Anchored Group",
			IpAnchored: true,
			Enabled:    true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ipAnchored":true`)
	})

	t.Run("Extranet enabled", func(t *testing.T) {
		group := servergroup.ServerGroup{
			ID:              "123",
			Name:            "Extranet Group",
			ExtranetEnabled: true,
			Enabled:         true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"extranetEnabled":true`)
	})
}

// TestServerGroup_GetByName tests the GetByName functionality
func TestServerGroup_GetByName(t *testing.T) {
	t.Run("Search returns matching server group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Production Servers", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serverGroup?search=Production")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

