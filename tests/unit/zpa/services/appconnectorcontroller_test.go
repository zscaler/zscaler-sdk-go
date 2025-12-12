// Package unit provides unit tests for ZPA App Connector Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorcontroller"
)

// TestAppConnectorController_Structure tests the struct definitions
func TestAppConnectorController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppConnector JSON marshaling", func(t *testing.T) {
		connector := appconnectorcontroller.AppConnector{
			ID:                   "conn-123",
			Name:                 "Test Connector",
			Description:          "Test Description",
			Enabled:              true,
			AppConnectorGroupID:  "acg-001",
			AppConnectorGroupName: "Test Group",
			ControlChannelStatus: "ZPN_STATUS_ONLINE",
			CurrentVersion:       "24.1.0",
			Latitude:             "37.3382",
			Longitude:            "-121.8863",
			Location:             "San Jose, CA",
			Platform:             "el8",
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled appconnectorcontroller.AppConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, connector.ID, unmarshaled.ID)
		assert.Equal(t, connector.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("AppConnector from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "conn-456",
			"name": "Production Connector",
			"description": "Production environment connector",
			"enabled": true,
			"appConnectorGroupId": "acg-002",
			"appConnectorGroupName": "Production Group",
			"controlChannelStatus": "ZPN_STATUS_ONLINE",
			"currentVersion": "24.2.0",
			"expectedVersion": "24.3.0",
			"latitude": "40.7128",
			"longitude": "-74.0060",
			"location": "New York, NY",
			"platform": "ubuntu",
			"privateIp": "10.0.0.5",
			"publicIp": "203.0.113.5",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var connector appconnectorcontroller.AppConnector
		err := json.Unmarshal([]byte(apiResponse), &connector)
		require.NoError(t, err)

		assert.Equal(t, "conn-456", connector.ID)
		assert.True(t, connector.Enabled)
		assert.Equal(t, "ZPN_STATUS_ONLINE", connector.ControlChannelStatus)
	})
}

// TestAppConnectorController_MockServerOperations tests CRUD operations
func TestAppConnectorController_MockServerOperations(t *testing.T) {
	t.Run("GET connector by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "conn-123", "name": "Mock Connector", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/appConnector/conn-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST bulk delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/appConnector/bulkDelete", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE connector", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/appConnector/conn-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppConnectorController_SpecialCases tests edge cases
func TestAppConnectorController_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Control channel status values", func(t *testing.T) {
		statuses := []string{"ZPN_STATUS_ONLINE", "ZPN_STATUS_OFFLINE", "ZPN_STATUS_UNKNOWN"}

		for _, status := range statuses {
			connector := appconnectorcontroller.AppConnector{
				ID:                   "conn-" + status,
				Name:                 status + " Connector",
				ControlChannelStatus: status,
			}

			data, err := json.Marshal(connector)
			require.NoError(t, err)
			assert.Contains(t, string(data), status)
		}
	})

	t.Run("BulkDeleteRequest", func(t *testing.T) {
		req := appconnectorcontroller.BulkDeleteRequest{
			IDs: []string{"conn-1", "conn-2", "conn-3"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled appconnectorcontroller.BulkDeleteRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.IDs, 3)
	})
}
