// Package unit provides unit tests for ZPA App Connector Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
)

// TestAppConnectorGroup_Structure tests the struct definitions
func TestAppConnectorGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppConnectorGroup JSON marshaling", func(t *testing.T) {
		group := appconnectorgroup.AppConnectorGroup{
			ID:                     "123",
			Name:                   "Test Connector Group",
			Description:            "Test Description",
			Enabled:                true,
			CityCountry:            "San Jose, US",
			CountryCode:            "US",
			Latitude:               "37.3382",
			Longitude:              "-121.8863",
			Location:               "San Jose, CA",
			OverrideVersionProfile: true,
			PRAEnabled:             false,
			WAFDisabled:            true,
			UpgradeDay:             "SUNDAY",
			UpgradeTimeInSecs:      "66600",
			VersionProfileID:       "0",
			TCPQuickAckApp:         true,
			TCPQuickAckAssistant:   true,
			UseInDrMode:            false,
			MicroTenantID:          "tenant-123",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled appconnectorgroup.AppConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.Equal(t, group.Enabled, unmarshaled.Enabled)
		assert.Equal(t, group.CityCountry, unmarshaled.CityCountry)
		assert.Equal(t, group.Latitude, unmarshaled.Latitude)
		assert.Equal(t, group.Longitude, unmarshaled.Longitude)
	})

	t.Run("AppConnectorGroup JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "456",
			"name": "Production Connector Group",
			"description": "Production environment connectors",
			"enabled": true,
			"cityCountry": "New York, US",
			"countryCode": "US",
			"latitude": "40.7128",
			"longitude": "-74.0060",
			"location": "New York, NY",
			"overrideVersionProfile": false,
			"praEnabled": true,
			"wafDisabled": false,
			"upgradeDay": "MONDAY",
			"upgradeTimeInSecs": "72000",
			"versionProfileId": "1",
			"versionProfileName": "Default",
			"tcpQuickAckApp": true,
			"tcpQuickAckAssistant": false,
			"tcpQuickAckReadAssistant": true,
			"lssAppConnectorGroup": false,
			"useInDrMode": true,
			"microtenantId": "mt-001",
			"microtenantName": "Test Microtenant",
			"connectorGroupType": "DEFAULT",
			"dnsQueryType": "IPV4_IPV6",
			"serverGroups": [
				{
					"id": "sg-001",
					"name": "Server Group 1",
					"enabled": true
				}
			]
		}`

		var group appconnectorgroup.AppConnectorGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "456", group.ID)
		assert.Equal(t, "Production Connector Group", group.Name)
		assert.True(t, group.Enabled)
		assert.True(t, group.PRAEnabled)
		assert.False(t, group.WAFDisabled)
		assert.True(t, group.UseInDrMode)
		assert.Equal(t, "mt-001", group.MicroTenantID)
		assert.Equal(t, "Test Microtenant", group.MicroTenantName)
		assert.Equal(t, "DEFAULT", group.ConnectorGroupType)
		assert.Equal(t, "IPV4_IPV6", group.DNSQueryType)
		assert.Len(t, group.AppServerGroup, 1)
	})

	t.Run("NPAssistantGroup structure", func(t *testing.T) {
		npGroup := appconnectorgroup.NPAssistantGroup{
			ID:                  "np-001",
			AppConnectorGroupID: "acg-001",
			MTU:                 "1500",
			LanSubnets: []appconnectorgroup.LanSubnet{
				{
					ID:     "subnet-001",
					Name:   "LAN Subnet 1",
					Subnet: "10.0.0.0/24",
				},
			},
		}

		data, err := json.Marshal(npGroup)
		require.NoError(t, err)

		var unmarshaled appconnectorgroup.NPAssistantGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, npGroup.ID, unmarshaled.ID)
		assert.Equal(t, npGroup.MTU, unmarshaled.MTU)
		assert.Len(t, unmarshaled.LanSubnets, 1)
		assert.Equal(t, "10.0.0.0/24", unmarshaled.LanSubnets[0].Subnet)
	})
}

// TestAppConnectorGroup_ResponseParsing tests parsing of various API responses
func TestAppConnectorGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse single app connector group response", func(t *testing.T) {
		response := `{
			"id": "789",
			"name": "Edge Connector Group",
			"description": "Edge location connectors",
			"enabled": true,
			"cityCountry": "London, GB",
			"countryCode": "GB",
			"latitude": "51.5074",
			"longitude": "-0.1278",
			"location": "London",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com"
		}`

		var group appconnectorgroup.AppConnectorGroup
		err := json.Unmarshal([]byte(response), &group)
		require.NoError(t, err)

		assert.Equal(t, "789", group.ID)
		assert.Equal(t, "Edge Connector Group", group.Name)
		assert.Equal(t, "London, GB", group.CityCountry)
		assert.Equal(t, "GB", group.CountryCode)
		assert.NotEmpty(t, group.CreationTime)
	})

	t.Run("Parse app connector group list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Group 1", "enabled": true},
				{"id": "2", "name": "Group 2", "enabled": false},
				{"id": "3", "name": "Group 3", "enabled": true}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []appconnectorgroup.AppConnectorGroup `json:"list"`
			TotalPages int                                   `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "1", listResp.List[0].ID)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[1].Enabled)
	})
}

// TestAppConnectorGroup_MockServerOperations tests CRUD operations with mock server
func TestAppConnectorGroup_MockServerOperations(t *testing.T) {
	t.Run("GET app connector group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/appConnectorGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "acg-123",
				"name": "Mock Connector Group",
				"description": "Created by mock server",
				"enabled": true,
				"cityCountry": "Tokyo, JP",
				"countryCode": "JP"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		// Verify server responds correctly
		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/appConnectorGroup/acg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all app connector groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Group A", "enabled": true, "cityCountry": "San Jose, US"},
					{"id": "2", "name": "Group B", "enabled": true, "cityCountry": "London, GB"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/appConnectorGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create app connector group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/appConnectorGroup")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-acg-456",
				"name": "New Connector Group",
				"description": "Newly created",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/appConnectorGroup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update app connector group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/appConnectorGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/appConnectorGroup/acg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE app connector group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Contains(t, r.URL.Path, "/appConnectorGroup/")

			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/appConnectorGroup/acg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppConnectorGroup_ErrorHandling tests error scenarios
func TestAppConnectorGroup_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "App Connector Group not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/appConnectorGroup/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request - Invalid payload", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Name is required"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/appConnectorGroup", "application/json", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("409 Conflict - Duplicate name", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"code": "DUPLICATE_NAME", "message": "An app connector group with this name already exists"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/appConnectorGroup", "application/json", nil)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

// TestAppConnectorGroup_SpecialCases tests edge cases and special scenarios
func TestAppConnectorGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Empty connectors list", func(t *testing.T) {
		group := appconnectorgroup.AppConnectorGroup{
			ID:         "123",
			Name:       "Empty Group",
			Enabled:    true,
			Connectors: nil,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled appconnectorgroup.AppConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Nil(t, unmarshaled.Connectors)
	})

	t.Run("Boolean fields default values", func(t *testing.T) {
		// When unmarshaling, boolean fields should default to false if not present
		response := `{"id": "123", "name": "Test"}`

		var group appconnectorgroup.AppConnectorGroup
		err := json.Unmarshal([]byte(response), &group)
		require.NoError(t, err)

		assert.False(t, group.Enabled)
		assert.False(t, group.OverrideVersionProfile)
		assert.False(t, group.PRAEnabled)
		assert.False(t, group.WAFDisabled)
		assert.False(t, group.TCPQuickAckApp)
		assert.False(t, group.UseInDrMode)
	})

	t.Run("Microtenant configuration", func(t *testing.T) {
		group := appconnectorgroup.AppConnectorGroup{
			ID:              "123",
			Name:            "Tenant Group",
			MicroTenantID:   "mt-001",
			MicroTenantName: "Tenant One",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"microtenantId":"mt-001"`)
		assert.Contains(t, string(data), `"microtenantName":"Tenant One"`)
	})

	t.Run("Version profile configuration", func(t *testing.T) {
		group := appconnectorgroup.AppConnectorGroup{
			ID:                            "123",
			Name:                          "Versioned Group",
			OverrideVersionProfile:        true,
			VersionProfileID:              "2",
			VersionProfileName:            "Custom Profile",
			VersionProfileVisibilityScope: "CUSTOMER",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled appconnectorgroup.AppConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.OverrideVersionProfile)
		assert.Equal(t, "2", unmarshaled.VersionProfileID)
		assert.Equal(t, "Custom Profile", unmarshaled.VersionProfileName)
	})
}

// TestAppConnectorGroup_GetByName tests the GetByName functionality
func TestAppConnectorGroup_GetByName(t *testing.T) {
	t.Run("Search returns matching group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Production Connectors", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/appConnectorGroup?search=Production")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Search returns empty when no match", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [],
				"totalPages": 0
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/appConnectorGroup?search=NonExistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

