// Package unit provides unit tests for ZPA PRA Portal service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PRAPortal represents the PRA portal for testing
type PRAPortal struct {
	ID                      string   `json:"id,omitempty"`
	Name                    string   `json:"name,omitempty"`
	Description             string   `json:"description,omitempty"`
	Enabled                 bool     `json:"enabled"`
	CName                   string   `json:"cName,omitempty"`
	Domain                  string   `json:"domain,omitempty"`
	CertificateID           string   `json:"certificateId,omitempty"`
	CertificateName         string   `json:"certificateName,omitempty"`
	UserNotification        string   `json:"userNotification"`
	UserNotificationEnabled bool     `json:"userNotificationEnabled"`
	ExtDomain               string   `json:"extDomain"`
	ExtDomainName           string   `json:"extDomainName"`
	ExtDomainTranslation    string   `json:"extDomainTranslation"`
	ExtLabel                string   `json:"extLabel"`
	UserPortalGid           string   `json:"userPortalGid,omitempty"`
	UserPortalName          string   `json:"userPortalName,omitempty"`
	MicroTenantID           string   `json:"microtenantId,omitempty"`
	MicroTenantName         string   `json:"microtenantName,omitempty"`
	IsSRAPortal             bool     `json:"isSRAPortal,omitempty"`
	ManagedByZs             bool     `json:"managedByZs,omitempty"`
	ApprovalReviewers       []string `json:"approvalReviewers,omitempty"`
	CreationTime            string   `json:"creationTime,omitempty"`
	ModifiedBy              string   `json:"modifiedBy,omitempty"`
	ModifiedTime            string   `json:"modifiedTime,omitempty"`
}

// TestPRAPortal_Structure tests the struct definitions
func TestPRAPortal_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PRAPortal JSON marshaling", func(t *testing.T) {
		portal := PRAPortal{
			ID:                      "portal-123",
			Name:                    "Corporate PRA Portal",
			Description:             "Main privileged access portal",
			Enabled:                 true,
			CName:                   "pra.corp.example.com",
			Domain:                  "pra.zscaler.com",
			CertificateID:           "cert-001",
			CertificateName:         "Corporate Cert",
			UserNotification:        "Welcome to the PRA Portal",
			UserNotificationEnabled: true,
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled PRAPortal
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, portal.ID, unmarshaled.ID)
		assert.Equal(t, portal.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.True(t, unmarshaled.UserNotificationEnabled)
	})

	t.Run("PRAPortal from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "portal-456",
			"name": "Guest PRA Portal",
			"description": "Portal for guest access",
			"enabled": true,
			"cName": "guest-pra.example.com",
			"domain": "guest.zscaler.com",
			"certificateId": "cert-002",
			"certificateName": "Guest Cert",
			"userNotification": "Guest access portal",
			"userNotificationEnabled": false,
			"extDomain": "ext.zscaler.com",
			"extDomainName": "External Domain",
			"extDomainTranslation": "translation",
			"extLabel": "External",
			"userPortalGid": "upg-001",
			"userPortalName": "User Portal",
			"microtenantId": "mt-001",
			"microtenantName": "Production",
			"isSRAPortal": true,
			"managedByZs": false,
			"approvalReviewers": ["admin@example.com", "security@example.com"],
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "1612137600000"
		}`

		var portal PRAPortal
		err := json.Unmarshal([]byte(apiResponse), &portal)
		require.NoError(t, err)

		assert.Equal(t, "portal-456", portal.ID)
		assert.Equal(t, "Guest PRA Portal", portal.Name)
		assert.True(t, portal.Enabled)
		assert.True(t, portal.IsSRAPortal)
		assert.Len(t, portal.ApprovalReviewers, 2)
	})
}

// TestPRAPortal_ResponseParsing tests parsing of API responses
func TestPRAPortal_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse portal list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Portal 1", "enabled": true},
				{"id": "2", "name": "Portal 2", "enabled": true},
				{"id": "3", "name": "Portal 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []PRAPortal `json:"list"`
			TotalPages int         `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestPRAPortal_MockServerOperations tests CRUD operations
func TestPRAPortal_MockServerOperations(t *testing.T) {
	t.Run("GET portal by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/praPortal/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "portal-123", "name": "Mock Portal", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praPortal/portal-123")
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

		resp, err := http.Get(server.URL + "/praPortal")
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

		resp, err := http.Post(server.URL+"/praPortal", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update portal", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/praPortal/portal-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE portal", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/praPortal/portal-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestPRAPortal_SpecialCases tests edge cases
func TestPRAPortal_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Portal with notification enabled", func(t *testing.T) {
		portal := PRAPortal{
			ID:                      "portal-notif",
			Name:                    "Notification Portal",
			UserNotification:        "Important: This is a privileged access session",
			UserNotificationEnabled: true,
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled PRAPortal
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.True(t, unmarshaled.UserNotificationEnabled)
		assert.Contains(t, unmarshaled.UserNotification, "privileged access")
	})

	t.Run("Portal with approval reviewers", func(t *testing.T) {
		portal := PRAPortal{
			ID:      "portal-approval",
			Name:    "Approval Portal",
			Enabled: true,
			ApprovalReviewers: []string{
				"admin1@example.com",
				"admin2@example.com",
				"security@example.com",
			},
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled PRAPortal
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.ApprovalReviewers, 3)
	})

	t.Run("Managed by Zscaler portal", func(t *testing.T) {
		portal := PRAPortal{
			ID:          "portal-managed",
			Name:        "Managed Portal",
			ManagedByZs: true,
			IsSRAPortal: true,
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"managedByZs":true`)
		assert.Contains(t, string(data), `"isSRAPortal":true`)
	})

	t.Run("Disabled portal", func(t *testing.T) {
		portal := PRAPortal{
			ID:      "portal-disabled",
			Name:    "Disabled Portal",
			Enabled: false,
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled PRAPortal
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})
}

