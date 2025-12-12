// Package unit provides unit tests for ZPA Private Cloud Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PrivateCloudController represents the private cloud controller for testing
type PrivateCloudController struct {
	ID                              string   `json:"id,omitempty"`
	Name                            string   `json:"name,omitempty"`
	Description                     string   `json:"description,omitempty"`
	Enabled                         bool     `json:"enabled,omitempty"`
	ControlChannelStatus            string   `json:"controlChannelStatus,omitempty"`
	CurrentVersion                  string   `json:"currentVersion,omitempty"`
	ExpectedVersion                 string   `json:"expectedVersion,omitempty"`
	PrivateIP                       string   `json:"privateIp,omitempty"`
	PublicIP                        string   `json:"publicIp,omitempty"`
	Latitude                        string   `json:"latitude,omitempty"`
	Longitude                       string   `json:"longitude,omitempty"`
	Location                        string   `json:"location,omitempty"`
	Platform                        string   `json:"platform,omitempty"`
	UpgradeStatus                   string   `json:"upgradeStatus,omitempty"`
	PrivateCloudControllerGroupID   string   `json:"privateCloudControllerGroupId,omitempty"`
	PrivateCloudControllerGroupName string   `json:"privateCloudControllerGroupName,omitempty"`
	MicrotenantID                   string   `json:"microtenantId,omitempty"`
	MicrotenantName                 string   `json:"microtenantName,omitempty"`
	CreationTime                    string   `json:"creationTime,omitempty"`
	ModifiedTime                    string   `json:"modifiedTime,omitempty"`
	ListenIps                       []string `json:"listenIps,omitempty"`
	PublishIps                      []string `json:"publishIps,omitempty"`
}

// TestPrivateCloudController_Structure tests the struct definitions
func TestPrivateCloudController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivateCloudController JSON marshaling", func(t *testing.T) {
		controller := PrivateCloudController{
			ID:                   "pcc-123",
			Name:                 "Private Cloud Controller 1",
			Description:          "Primary controller",
			Enabled:              true,
			ControlChannelStatus: "CONNECTED",
			CurrentVersion:       "22.1.0",
			PrivateIP:            "10.0.0.10",
			PublicIP:             "203.0.113.10",
			Location:             "San Jose, CA",
			Platform:             "RHEL",
		}

		data, err := json.Marshal(controller)
		require.NoError(t, err)

		var unmarshaled PrivateCloudController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, controller.ID, unmarshaled.ID)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "CONNECTED", unmarshaled.ControlChannelStatus)
	})

	t.Run("PrivateCloudController from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pcc-456",
			"name": "Production Controller",
			"description": "Production private cloud controller",
			"enabled": true,
			"controlChannelStatus": "CONNECTED",
			"currentVersion": "22.2.0",
			"expectedVersion": "22.3.0",
			"privateIp": "10.0.1.10",
			"publicIp": "203.0.113.20",
			"latitude": "37.3861",
			"longitude": "-122.0839",
			"location": "Mountain View, CA",
			"platform": "Ubuntu",
			"upgradeStatus": "CURRENT",
			"privateCloudControllerGroupId": "pccg-001",
			"privateCloudControllerGroupName": "US West Group",
			"listenIps": ["10.0.1.10", "10.0.1.11"],
			"publishIps": ["203.0.113.20"],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var controller PrivateCloudController
		err := json.Unmarshal([]byte(apiResponse), &controller)
		require.NoError(t, err)

		assert.Equal(t, "pcc-456", controller.ID)
		assert.True(t, controller.Enabled)
		assert.Len(t, controller.ListenIps, 2)
	})
}

// TestPrivateCloudController_MockServerOperations tests CRUD operations
func TestPrivateCloudController_MockServerOperations(t *testing.T) {
	t.Run("GET controller by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/privateCloudController/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pcc-123", "name": "Mock Controller"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/privateCloudController/pcc-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT restart controller", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			assert.Contains(t, r.URL.Path, "/restart/")
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/privateCloudController/restart/pcc-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE controller", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/privateCloudController/pcc-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

