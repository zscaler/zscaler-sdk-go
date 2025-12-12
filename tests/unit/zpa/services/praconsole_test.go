// Package unit provides unit tests for ZPA PRA Console service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PRAConsole represents the PRA console for testing
type PRAConsole struct {
	ID              string              `json:"id,omitempty"`
	Name            string              `json:"name,omitempty"`
	Description     string              `json:"description,omitempty"`
	Enabled         bool                `json:"enabled"`
	IconText        string              `json:"iconText,omitempty"`
	CreationTime    string              `json:"creationTime,omitempty"`
	ModifiedBy      string              `json:"modifiedBy,omitempty"`
	ModifiedTime    string              `json:"modifiedTime,omitempty"`
	MicroTenantID   string              `json:"microtenantId,omitempty"`
	MicroTenantName string              `json:"microtenantName,omitempty"`
	PRAApplication  PRAConsoleApp       `json:"praApplication,omitempty"`
	PRAPortals      []PRAConsolePortal  `json:"praPortals"`
}

// PRAConsoleApp represents the PRA application reference
type PRAConsoleApp struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// PRAConsolePortal represents the PRA portal reference
type PRAConsolePortal struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestPRAConsole_Structure tests the struct definitions
func TestPRAConsole_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PRAConsole JSON marshaling", func(t *testing.T) {
		console := PRAConsole{
			ID:          "console-123",
			Name:        "SSH Console",
			Description: "SSH access console",
			Enabled:     true,
			IconText:    "base64icondata==",
			PRAApplication: PRAConsoleApp{
				ID:   "app-001",
				Name: "SSH Application",
			},
			PRAPortals: []PRAConsolePortal{
				{ID: "portal-001", Name: "Main Portal"},
			},
		}

		data, err := json.Marshal(console)
		require.NoError(t, err)

		var unmarshaled PRAConsole
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, console.ID, unmarshaled.ID)
		assert.Equal(t, console.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "app-001", unmarshaled.PRAApplication.ID)
		assert.Len(t, unmarshaled.PRAPortals, 1)
	})

	t.Run("PRAConsole from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "console-456",
			"name": "RDP Console",
			"description": "Remote Desktop console",
			"enabled": true,
			"iconText": "cmRwaWNvbg==",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"praApplication": {
				"id": "app-002",
				"name": "RDP Application"
			},
			"praPortals": [
				{"id": "portal-001", "name": "Corporate Portal"},
				{"id": "portal-002", "name": "Guest Portal"}
			]
		}`

		var console PRAConsole
		err := json.Unmarshal([]byte(apiResponse), &console)
		require.NoError(t, err)

		assert.Equal(t, "console-456", console.ID)
		assert.Equal(t, "RDP Console", console.Name)
		assert.True(t, console.Enabled)
		assert.Equal(t, "RDP Application", console.PRAApplication.Name)
		assert.Len(t, console.PRAPortals, 2)
	})
}

// TestPRAConsole_ResponseParsing tests parsing of API responses
func TestPRAConsole_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse console list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Console 1", "enabled": true},
				{"id": "2", "name": "Console 2", "enabled": true},
				{"id": "3", "name": "Console 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []PRAConsole `json:"list"`
			TotalPages int          `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestPRAConsole_MockServerOperations tests CRUD operations
func TestPRAConsole_MockServerOperations(t *testing.T) {
	t.Run("GET console by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/praConsole/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "console-123", "name": "Mock Console", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praConsole/console-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET console by portal ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/praPortal/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "console-123", "name": "Portal Console"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praConsole/praPortal/portal-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create console", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-console", "name": "New Console"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/praConsole", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST bulk create consoles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/bulk")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[{"id": "console-1"}, {"id": "console-2"}]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/praConsole/bulk", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update console", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/praConsole/console-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE console", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/praConsole/console-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestPRAConsole_SpecialCases tests edge cases
func TestPRAConsole_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Console with multiple portals", func(t *testing.T) {
		console := PRAConsole{
			ID:      "console-multi",
			Name:    "Multi-Portal Console",
			Enabled: true,
			PRAPortals: []PRAConsolePortal{
				{ID: "portal-1", Name: "Portal 1"},
				{ID: "portal-2", Name: "Portal 2"},
				{ID: "portal-3", Name: "Portal 3"},
			},
		}

		data, err := json.Marshal(console)
		require.NoError(t, err)

		var unmarshaled PRAConsole
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.PRAPortals, 3)
	})

	t.Run("Disabled console", func(t *testing.T) {
		console := PRAConsole{
			ID:      "console-disabled",
			Name:    "Disabled Console",
			Enabled: false,
		}

		data, err := json.Marshal(console)
		require.NoError(t, err)

		var unmarshaled PRAConsole
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})
}

