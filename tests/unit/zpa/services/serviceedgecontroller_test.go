// Package unit provides unit tests for ZPA Service Edge Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ServiceEdgeController represents the service edge controller for testing
type ServiceEdgeController struct {
	ID                   string   `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Enabled              bool     `json:"enabled,omitempty"`
	ServiceEdgeGroupID   string   `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName string   `json:"serviceEdgeGroupName,omitempty"`
	ControlChannelStatus string   `json:"controlChannelStatus,omitempty"`
	CurrentVersion       string   `json:"currentVersion,omitempty"`
	ExpectedVersion      string   `json:"expectedVersion,omitempty"`
	PrivateIP            string   `json:"privateIp,omitempty"`
	PublicIP             string   `json:"publicIp,omitempty"`
	Latitude             string   `json:"latitude,omitempty"`
	Longitude            string   `json:"longitude,omitempty"`
	Location             string   `json:"location,omitempty"`
	Platform             string   `json:"platform,omitempty"`
	UpgradeStatus        string   `json:"upgradeStatus,omitempty"`
	ListenIPs            []string `json:"listenIps,omitempty"`
	PublishIPs           []string `json:"publishIps,omitempty"`
	MicroTenantID        string   `json:"microtenantId,omitempty"`
	MicroTenantName      string   `json:"microtenantName,omitempty"`
	CreationTime         string   `json:"creationTime,omitempty"`
	ModifiedTime         string   `json:"modifiedTime,omitempty"`
}

// TestServiceEdgeController_Structure tests the struct definitions
func TestServiceEdgeController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ServiceEdgeController JSON marshaling", func(t *testing.T) {
		edge := ServiceEdgeController{
			ID:                   "se-123",
			Name:                 "Service Edge 1",
			Description:          "Primary service edge",
			Enabled:              true,
			ServiceEdgeGroupID:   "seg-001",
			ServiceEdgeGroupName: "US West Group",
			ControlChannelStatus: "CONNECTED",
			CurrentVersion:       "22.1.0",
			PrivateIP:            "10.0.0.100",
			PublicIP:             "203.0.113.100",
			Location:             "San Jose, CA",
			Platform:             "Ubuntu",
		}

		data, err := json.Marshal(edge)
		require.NoError(t, err)

		var unmarshaled ServiceEdgeController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, edge.ID, unmarshaled.ID)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "CONNECTED", unmarshaled.ControlChannelStatus)
	})

	t.Run("ServiceEdgeController from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "se-456",
			"name": "Service Edge 2",
			"description": "Secondary service edge",
			"enabled": true,
			"serviceEdgeGroupId": "seg-002",
			"serviceEdgeGroupName": "EU Central Group",
			"controlChannelStatus": "CONNECTED",
			"currentVersion": "22.2.0",
			"expectedVersion": "22.3.0",
			"privateIp": "10.0.1.100",
			"publicIp": "203.0.113.200",
			"latitude": "50.1109",
			"longitude": "8.6821",
			"location": "Frankfurt, Germany",
			"platform": "RHEL",
			"upgradeStatus": "CURRENT",
			"listenIps": ["10.0.1.100", "10.0.1.101"],
			"publishIps": ["203.0.113.200"],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var edge ServiceEdgeController
		err := json.Unmarshal([]byte(apiResponse), &edge)
		require.NoError(t, err)

		assert.Equal(t, "se-456", edge.ID)
		assert.True(t, edge.Enabled)
		assert.Len(t, edge.ListenIPs, 2)
	})
}

// TestServiceEdgeController_MockServerOperations tests CRUD operations
func TestServiceEdgeController_MockServerOperations(t *testing.T) {
	t.Run("GET service edge by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/serviceEdge/")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "se-123", "name": "Mock Edge"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serviceEdge/se-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST bulk delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/bulkDelete")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/serviceEdge/bulkDelete", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE service edge", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/serviceEdge/se-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

