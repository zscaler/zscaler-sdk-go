// Package unit provides unit tests for ZPA User Portal Link service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_link"
)

// TestUserPortalLink_Structure tests the struct definitions
func TestUserPortalLink_Structure(t *testing.T) {
	t.Parallel()

	t.Run("UserPortalLink JSON marshaling", func(t *testing.T) {
		link := portal_link.UserPortalLink{
			ID:            "link-123",
			Name:          "Application Link",
			Description:   "Link to internal application",
			Enabled:       true,
			ApplicationID: "app-001",
			Link:          "https://app.internal.com",
			LinkPath:      "/app",
			Protocol:      "HTTPS",
			UserPortalID:  "portal-001",
			UserPortals: []portal_controller.UserPortalController{
				{ID: "portal-001", Name: "Main Portal"},
			},
		}

		data, err := json.Marshal(link)
		require.NoError(t, err)

		var unmarshaled portal_link.UserPortalLink
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, link.ID, unmarshaled.ID)
		assert.Equal(t, link.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.UserPortals, 1)
	})

	t.Run("UserPortalLink from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "link-456",
			"name": "HR Portal Link",
			"description": "Link to HR application",
			"enabled": true,
			"applicationId": "app-002",
			"link": "https://hr.company.com",
			"linkPath": "/hr",
			"protocol": "HTTPS",
			"iconText": "base64icon==",
			"userPortalId": "portal-002",
			"userPortals": [
				{"id": "portal-001", "name": "Corporate Portal"},
				{"id": "portal-002", "name": "HR Portal"}
			],
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var link portal_link.UserPortalLink
		err := json.Unmarshal([]byte(apiResponse), &link)
		require.NoError(t, err)

		assert.Equal(t, "link-456", link.ID)
		assert.Equal(t, "HR Portal Link", link.Name)
		assert.True(t, link.Enabled)
		assert.Len(t, link.UserPortals, 2)
	})

	t.Run("PortalLinkBulkRequest JSON marshaling", func(t *testing.T) {
		bulk := portal_link.PortalLinkBulkRequest{
			UserPortalLinks: []portal_link.UserPortalLink{
				{ID: "link-1", Name: "Link 1"},
				{ID: "link-2", Name: "Link 2"},
			},
			UserPortals: []portal_controller.UserPortalController{
				{ID: "portal-1"},
			},
		}

		data, err := json.Marshal(bulk)
		require.NoError(t, err)

		var unmarshaled portal_link.PortalLinkBulkRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UserPortalLinks, 2)
		assert.Len(t, unmarshaled.UserPortals, 1)
	})
}

// TestUserPortalLink_MockServerOperations tests CRUD operations
func TestUserPortalLink_MockServerOperations(t *testing.T) {
	t.Run("GET link by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/userPortalLink/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "link-123", "name": "Mock Link", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userPortalLink/link-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET links by user portal ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/userPortal/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "link-123", "name": "Portal Link"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userPortalLink/userPortal/portal-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create link", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-link", "name": "New Link"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/userPortalLink", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST bulk create links", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Contains(t, r.URL.Path, "/bulk")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"userPortalLinks": [{"id": "1"}, {"id": "2"}]}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/userPortalLink/bulk", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE link", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/userPortalLink/link-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestUserPortalLink_SpecialCases tests edge cases
func TestUserPortalLink_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Protocol types", func(t *testing.T) {
		protocols := []string{"HTTP", "HTTPS", "RDP", "SSH", "VNC"}

		for _, protocol := range protocols {
			link := portal_link.UserPortalLink{
				ID:       "link-" + protocol,
				Name:     protocol + " Link",
				Protocol: protocol,
			}

			data, err := json.Marshal(link)
			require.NoError(t, err)

			assert.Contains(t, string(data), protocol)
		}
	})

	t.Run("Link with multiple portals", func(t *testing.T) {
		link := portal_link.UserPortalLink{
			ID:   "link-multi",
			Name: "Multi-Portal Link",
			UserPortals: []portal_controller.UserPortalController{
				{ID: "portal-1", Name: "Portal 1"},
				{ID: "portal-2", Name: "Portal 2"},
				{ID: "portal-3", Name: "Portal 3"},
			},
		}

		data, err := json.Marshal(link)
		require.NoError(t, err)

		var unmarshaled portal_link.UserPortalLink
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.UserPortals, 3)
	})
}
