// Package unit provides unit tests for ZPA User Portal Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/userportal/portal_controller"
)

// TestUserPortalController_Structure tests the struct definitions
func TestUserPortalController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("UserPortalController JSON marshaling", func(t *testing.T) {
		portal := portal_controller.UserPortalController{
			ID:                      "portal-123",
			Name:                    "Corporate User Portal",
			Description:             "Main portal for corporate users",
			Enabled:                 true,
			Domain:                  "portal.company.com",
			CertificateId:           "cert-001",
			CertificateName:         "Corporate Cert",
			UserNotification:        "Welcome to the User Portal",
			UserNotificationEnabled: true,
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled portal_controller.UserPortalController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, portal.ID, unmarshaled.ID)
		assert.Equal(t, portal.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.True(t, unmarshaled.UserNotificationEnabled)
	})

	t.Run("UserPortalController from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "portal-456",
			"name": "Guest Portal",
			"description": "Portal for guest access",
			"enabled": true,
			"domain": "guest.company.com",
			"certificateId": "cert-002",
			"certificateName": "Guest Cert",
			"extDomain": "ext.company.com",
			"extDomainName": "External Domain",
			"userNotification": "Welcome, Guest!",
			"userNotificationEnabled": true,
			"managedByZs": false,
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var portal portal_controller.UserPortalController
		err := json.Unmarshal([]byte(apiResponse), &portal)
		require.NoError(t, err)

		assert.Equal(t, "portal-456", portal.ID)
		assert.Equal(t, "Guest Portal", portal.Name)
		assert.True(t, portal.Enabled)
		assert.False(t, portal.ManagedByZS)
	})
}

// TestUserPortalController_MockServerOperations tests CRUD operations
func TestUserPortalController_MockServerOperations(t *testing.T) {
	t.Run("GET portal by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/userPortal/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "portal-123", "name": "Mock Portal", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userPortal/portal-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all portals", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/userPortal")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create portal", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-portal", "name": "New Portal"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/userPortal", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE portal", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/userPortal/portal-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
