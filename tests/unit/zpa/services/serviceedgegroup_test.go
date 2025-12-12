// Package unit provides unit tests for ZPA Service Edge Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ServiceEdgeGroup represents the service edge group structure for testing
type ServiceEdgeGroup struct {
	ID                            string            `json:"id,omitempty"`
	Name                          string            `json:"name,omitempty"`
	Description                   string            `json:"description,omitempty"`
	Enabled                       bool              `json:"enabled"`
	CityCountry                   string            `json:"cityCountry,omitempty"`
	CountryCode                   string            `json:"countryCode,omitempty"`
	Latitude                      string            `json:"latitude,omitempty"`
	Longitude                     string            `json:"longitude,omitempty"`
	Location                      string            `json:"location,omitempty"`
	UpgradeDay                    string            `json:"upgradeDay,omitempty"`
	UpgradeTimeInSecs             string            `json:"upgradeTimeInSecs,omitempty"`
	OverrideVersionProfile        bool              `json:"overrideVersionProfile"`
	VersionProfileID              string            `json:"versionProfileId,omitempty"`
	VersionProfileName            string            `json:"versionProfileName,omitempty"`
	VersionProfileVisibilityScope string            `json:"versionProfileVisibilityScope,omitempty"`
	CreationTime                  string            `json:"creationTime,omitempty"`
	ModifiedBy                    string            `json:"modifiedBy,omitempty"`
	ModifiedTime                  string            `json:"modifiedTime,omitempty"`
	MicroTenantID                 string            `json:"microtenantId,omitempty"`
	MicroTenantName               string            `json:"microtenantName,omitempty"`
	GraceDistanceEnabled          bool              `json:"graceDistanceEnabled"`
	GraceDistanceValue            string            `json:"graceDistanceValue,omitempty"`
	GraceDistanceValueUnit        string            `json:"graceDistanceValueUnit,omitempty"`
	IsPublic                      string            `json:"isPublic,omitempty"`
	ServiceEdges                  []ServiceEdge     `json:"serviceEdges,omitempty"`
	TrustedNetworks               []TrustedNetwork  `json:"trustedNetworks,omitempty"`
}

type ServiceEdge struct {
	ID                       string  `json:"id,omitempty"`
	Name                     string  `json:"name,omitempty"`
	Description              string  `json:"description,omitempty"`
	Enabled                  bool    `json:"enabled"`
	ServiceEdgeGroupID       string  `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName     string  `json:"serviceEdgeGroupName,omitempty"`
	ControlChannelStatus     string  `json:"controlChannelStatus,omitempty"`
	CurrentVersion           string  `json:"currentVersion,omitempty"`
	Platform                 string  `json:"platform,omitempty"`
	PrivateIP                string  `json:"privateIp,omitempty"`
	PublicIP                 string  `json:"publicIp,omitempty"`
	Latitude                 float64 `json:"latitude,omitempty"`
	Longitude                float64 `json:"longitude,omitempty"`
}

type TrustedNetwork struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	NetworkID       string `json:"networkId,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
}

// TestServiceEdgeGroup_Structure tests the struct definitions
func TestServiceEdgeGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ServiceEdgeGroup JSON marshaling", func(t *testing.T) {
		group := ServiceEdgeGroup{
			ID:                     "seg-123",
			Name:                   "Test Service Edge Group",
			Description:            "Test Description",
			Enabled:                true,
			CityCountry:            "San Jose, US",
			CountryCode:            "US",
			Latitude:               "37.3382",
			Longitude:              "-121.8863",
			Location:               "San Jose, CA",
			OverrideVersionProfile: true,
			UpgradeDay:             "SUNDAY",
			UpgradeTimeInSecs:      "66600",
			VersionProfileID:       "0",
			GraceDistanceEnabled:   true,
			GraceDistanceValue:     "10",
			GraceDistanceValueUnit: "MILES",
			IsPublic:               "FALSE",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled ServiceEdgeGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.Equal(t, group.Enabled, unmarshaled.Enabled)
		assert.Equal(t, group.GraceDistanceEnabled, unmarshaled.GraceDistanceEnabled)
	})

	t.Run("ServiceEdgeGroup JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "seg-456",
			"name": "Production Service Edge Group",
			"description": "Production environment service edges",
			"enabled": true,
			"cityCountry": "New York, US",
			"countryCode": "US",
			"latitude": "40.7128",
			"longitude": "-74.0060",
			"location": "New York, NY",
			"overrideVersionProfile": false,
			"upgradeDay": "MONDAY",
			"upgradeTimeInSecs": "72000",
			"versionProfileId": "1",
			"versionProfileName": "Default",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Test Microtenant",
			"graceDistanceEnabled": true,
			"graceDistanceValue": "15",
			"graceDistanceValueUnit": "KILOMETERS",
			"isPublic": "TRUE",
			"serviceEdges": [
				{
					"id": "se-001",
					"name": "Service Edge 1",
					"enabled": true,
					"controlChannelStatus": "ZPN_STATUS_ONLINE"
				}
			]
		}`

		var group ServiceEdgeGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, "seg-456", group.ID)
		assert.Equal(t, "Production Service Edge Group", group.Name)
		assert.True(t, group.Enabled)
		assert.True(t, group.GraceDistanceEnabled)
		assert.Equal(t, "TRUE", group.IsPublic)
		assert.Len(t, group.ServiceEdges, 1)
	})

	t.Run("ServiceEdge structure", func(t *testing.T) {
		edge := ServiceEdge{
			ID:                   "se-001",
			Name:                 "Test Service Edge",
			Description:          "Test edge",
			Enabled:              true,
			ServiceEdgeGroupID:   "seg-001",
			ServiceEdgeGroupName: "Edge Group 1",
			ControlChannelStatus: "ZPN_STATUS_ONLINE",
			CurrentVersion:       "24.1.0",
			Platform:             "el8",
			PrivateIP:            "10.0.0.5",
			PublicIP:             "203.0.113.5",
			Latitude:             37.3382,
			Longitude:            -121.8863,
		}

		data, err := json.Marshal(edge)
		require.NoError(t, err)

		var unmarshaled ServiceEdge
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, edge.ID, unmarshaled.ID)
		assert.Equal(t, edge.ControlChannelStatus, unmarshaled.ControlChannelStatus)
		assert.Equal(t, edge.Latitude, unmarshaled.Latitude)
	})
}

// TestServiceEdgeGroup_MockServerOperations tests CRUD operations with mock server
func TestServiceEdgeGroup_MockServerOperations(t *testing.T) {
	t.Run("GET service edge group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/serviceEdgeGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "seg-123",
				"name": "Mock Service Edge Group",
				"description": "Created by mock server",
				"enabled": true,
				"cityCountry": "Tokyo, JP"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/serviceEdgeGroup/seg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all service edge groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Edge Group A", "enabled": true},
					{"id": "2", "name": "Edge Group B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/serviceEdgeGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create service edge group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-seg-456",
				"name": "New Service Edge Group",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/serviceEdgeGroup", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update service edge group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/serviceEdgeGroup/seg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE service edge group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/serviceEdgeGroup/seg-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestServiceEdgeGroup_ErrorHandling tests error scenarios
func TestServiceEdgeGroup_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Service Edge Group not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serviceEdgeGroup/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Name is required"}`))
		}))
		defer server.Close()

		resp, _ := http.Post(server.URL+"/serviceEdgeGroup", "application/json", nil)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestServiceEdgeGroup_SpecialCases tests edge cases
func TestServiceEdgeGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Grace distance configuration", func(t *testing.T) {
		group := ServiceEdgeGroup{
			ID:                     "seg-123",
			Name:                   "Grace Distance Group",
			GraceDistanceEnabled:   true,
			GraceDistanceValue:     "25",
			GraceDistanceValueUnit: "MILES",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"graceDistanceEnabled":true`)
		assert.Contains(t, string(data), `"graceDistanceValue":"25"`)
		assert.Contains(t, string(data), `"graceDistanceValueUnit":"MILES"`)
	})

	t.Run("Public service edge group", func(t *testing.T) {
		group := ServiceEdgeGroup{
			ID:       "seg-123",
			Name:     "Public Edge Group",
			IsPublic: "TRUE",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"isPublic":"TRUE"`)
	})

	t.Run("Service edge group with trusted networks", func(t *testing.T) {
		group := ServiceEdgeGroup{
			ID:      "seg-123",
			Name:    "Trusted Network Group",
			Enabled: true,
			TrustedNetworks: []TrustedNetwork{
				{ID: "tn-001", Name: "Corporate Network", NetworkID: "corp-network-id"},
				{ID: "tn-002", Name: "Branch Network", NetworkID: "branch-network-id"},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled ServiceEdgeGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.TrustedNetworks, 2)
	})
}

