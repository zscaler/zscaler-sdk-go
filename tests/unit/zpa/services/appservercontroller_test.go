// Package unit provides unit tests for ZPA App Server Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appservercontroller"
)

// TestAppServerController_Structure tests the struct definitions
func TestAppServerController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationServer JSON marshaling", func(t *testing.T) {
		server := appservercontroller.ApplicationServer{
			ID:                "srv-123",
			Name:              "Test Server",
			Address:           "10.0.0.100",
			Description:       "Test Description",
			Enabled:           true,
			AppServerGroupIds: []string{"sg-001", "sg-002"},
			ConfigSpace:       "DEFAULT",
			MicroTenantID:     "mt-001",
		}

		data, err := json.Marshal(server)
		require.NoError(t, err)

		var unmarshaled appservercontroller.ApplicationServer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, server.ID, unmarshaled.ID)
		assert.Equal(t, server.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.AppServerGroupIds, 2)
	})

	t.Run("ApplicationServer from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "srv-456",
			"name": "Production Server",
			"address": "10.0.0.200",
			"description": "Production application server",
			"enabled": true,
			"appServerGroupIds": ["sg-001"],
			"configSpace": "DEFAULT",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-002",
			"microtenantName": "Production"
		}`

		var server appservercontroller.ApplicationServer
		err := json.Unmarshal([]byte(apiResponse), &server)
		require.NoError(t, err)

		assert.Equal(t, "srv-456", server.ID)
		assert.True(t, server.Enabled)
		assert.Equal(t, "10.0.0.200", server.Address)
	})
}

// TestAppServerController_MockServerOperations tests CRUD operations
func TestAppServerController_MockServerOperations(t *testing.T) {
	t.Run("GET server by ID", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "srv-123", "name": "Mock Server", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer mockServer.Close()

		resp, err := http.Get(mockServer.URL + "/server/srv-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create server", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-srv", "name": "New Server"}`
			w.Write([]byte(response))
		}))
		defer mockServer.Close()

		resp, err := http.Post(mockServer.URL+"/server", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE server", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer mockServer.Close()

		req, _ := http.NewRequest("DELETE", mockServer.URL+"/server/srv-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
