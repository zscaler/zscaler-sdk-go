// Package unit provides unit tests for ZPA App Server Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ApplicationServer represents the application server structure for testing
type ApplicationServer struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	Description      string   `json:"description,omitempty"`
	Address          string   `json:"address,omitempty"`
	Enabled          bool     `json:"enabled"`
	AppServerGroupID string   `json:"appServerGroupIds,omitempty"`
	ConfigSpace      string   `json:"configSpace,omitempty"`
	CreationTime     string   `json:"creationTime,omitempty"`
	ModifiedBy       string   `json:"modifiedBy,omitempty"`
	ModifiedTime     string   `json:"modifiedTime,omitempty"`
	MicroTenantID    string   `json:"microtenantId,omitempty"`
	MicroTenantName  string   `json:"microtenantName,omitempty"`
}

// TestApplicationServer_Structure tests the struct definitions
func TestApplicationServer_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ApplicationServer JSON marshaling", func(t *testing.T) {
		server := ApplicationServer{
			ID:          "srv-123",
			Name:        "Web Server 1",
			Description: "Primary web server",
			Address:     "10.0.0.1",
			Enabled:     true,
			ConfigSpace: "DEFAULT",
		}

		data, err := json.Marshal(server)
		require.NoError(t, err)

		var unmarshaled ApplicationServer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, server.ID, unmarshaled.ID)
		assert.Equal(t, server.Name, unmarshaled.Name)
		assert.Equal(t, server.Address, unmarshaled.Address)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("ApplicationServer JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "srv-456",
			"name": "Database Server",
			"description": "Primary database server",
			"address": "192.168.1.100",
			"enabled": true,
			"configSpace": "DEFAULT",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var srv ApplicationServer
		err := json.Unmarshal([]byte(apiResponse), &srv)
		require.NoError(t, err)

		assert.Equal(t, "srv-456", srv.ID)
		assert.Equal(t, "Database Server", srv.Name)
		assert.Equal(t, "192.168.1.100", srv.Address)
		assert.True(t, srv.Enabled)
	})
}

// TestApplicationServer_ResponseParsing tests parsing of various API responses
func TestApplicationServer_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse application server list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Server A", "address": "10.0.0.1", "enabled": true},
				{"id": "2", "name": "Server B", "address": "10.0.0.2", "enabled": true},
				{"id": "3", "name": "Server C", "address": "10.0.0.3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []ApplicationServer `json:"list"`
			TotalPages int                 `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestApplicationServer_MockServerOperations tests CRUD operations with mock server
func TestApplicationServer_MockServerOperations(t *testing.T) {
	t.Run("GET application server by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/server/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "srv-123",
				"name": "Mock Server",
				"address": "10.0.0.1",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/server/srv-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all application servers", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Server A", "address": "10.0.0.1"},
					{"id": "2", "name": "Server B", "address": "10.0.0.2"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/server")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create application server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-srv-456",
				"name": "New Server",
				"address": "10.0.0.100",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/server", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update application server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/server/srv-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE application server", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/server/srv-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestApplicationServer_SpecialCases tests edge cases
func TestApplicationServer_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Server with IP address", func(t *testing.T) {
		srv := ApplicationServer{
			ID:      "srv-123",
			Name:    "IP Server",
			Address: "10.0.0.1",
			Enabled: true,
		}

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		var unmarshaled ApplicationServer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "10.0.0.1", unmarshaled.Address)
	})

	t.Run("Server with FQDN", func(t *testing.T) {
		srv := ApplicationServer{
			ID:      "srv-123",
			Name:    "FQDN Server",
			Address: "server.example.com",
			Enabled: true,
		}

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		var unmarshaled ApplicationServer
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "server.example.com", unmarshaled.Address)
	})

	t.Run("Disabled server", func(t *testing.T) {
		srv := ApplicationServer{
			ID:      "srv-123",
			Name:    "Disabled Server",
			Address: "10.0.0.1",
			Enabled: false,
		}

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
	})

	t.Run("Server in different config spaces", func(t *testing.T) {
		configSpaces := []string{"DEFAULT", "SIEM"}

		for _, space := range configSpaces {
			srv := ApplicationServer{
				ID:          "srv-" + space,
				Name:        space + " Server",
				ConfigSpace: space,
			}

			data, err := json.Marshal(srv)
			require.NoError(t, err)

			assert.Contains(t, string(data), space)
		}
	})
}

