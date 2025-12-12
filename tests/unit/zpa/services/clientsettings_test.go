// Package unit provides unit tests for ZPA Client Settings service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ClientSettings represents the client settings for testing
type ClientSettings struct {
	ID                           string `json:"id,omitempty"`
	Name                         string `json:"name,omitempty"`
	ClientCertificateType        string `json:"clientCertificateType,omitempty"`
	SingningCertExpiryInEpochSec string `json:"singningCertExpiryInEpochSec,omitempty"`
	EnrollmentCertId             string `json:"enrollmentCertId,omitempty"`
	EnrollmentCertName           string `json:"enrollmentCertName,omitempty"`
	CreationTime                 string `json:"creationTime,omitempty"`
	ModifiedBy                   string `json:"modifiedBy,omitempty"`
}

// TestClientSettings_Structure tests the struct definitions
func TestClientSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ClientSettings JSON marshaling", func(t *testing.T) {
		settings := ClientSettings{
			ID:                           "cs-123",
			Name:                         "ZAPP Client Settings",
			ClientCertificateType:        "ZAPP_CLIENT",
			SingningCertExpiryInEpochSec: "1735689600",
			EnrollmentCertId:             "ec-001",
			EnrollmentCertName:           "Zscaler Enrollment Cert",
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		var unmarshaled ClientSettings
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, settings.ID, unmarshaled.ID)
		assert.Equal(t, settings.ClientCertificateType, unmarshaled.ClientCertificateType)
		assert.Equal(t, settings.EnrollmentCertName, unmarshaled.EnrollmentCertName)
	})

	t.Run("ClientSettings from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cs-456",
			"name": "Isolation Client Settings",
			"clientCertificateType": "ISOLATION_CLIENT",
			"singningCertExpiryInEpochSec": "1767225600",
			"enrollmentCertId": "ec-002",
			"enrollmentCertName": "Isolation Enrollment Cert",
			"creationTime": "1609459200000",
			"modifiedBy": "admin@example.com"
		}`

		var settings ClientSettings
		err := json.Unmarshal([]byte(apiResponse), &settings)
		require.NoError(t, err)

		assert.Equal(t, "cs-456", settings.ID)
		assert.Equal(t, "ISOLATION_CLIENT", settings.ClientCertificateType)
		assert.Equal(t, "ec-002", settings.EnrollmentCertId)
	})
}

// TestClientSettings_ResponseParsing tests parsing of API responses
func TestClientSettings_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse client settings list response", func(t *testing.T) {
		response := `[
			{"id": "1", "name": "Settings 1", "clientCertificateType": "ZAPP_CLIENT"},
			{"id": "2", "name": "Settings 2", "clientCertificateType": "ISOLATION_CLIENT"},
			{"id": "3", "name": "Settings 3", "clientCertificateType": "APP_PROTECTION"}
		]`

		var settings []ClientSettings
		err := json.Unmarshal([]byte(response), &settings)
		require.NoError(t, err)

		assert.Len(t, settings, 3)
		assert.Equal(t, "ZAPP_CLIENT", settings[0].ClientCertificateType)
		assert.Equal(t, "ISOLATION_CLIENT", settings[1].ClientCertificateType)
		assert.Equal(t, "APP_PROTECTION", settings[2].ClientCertificateType)
	})
}

// TestClientSettings_MockServerOperations tests operations
func TestClientSettings_MockServerOperations(t *testing.T) {
	t.Run("GET client settings by type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/clientSetting")
			assert.Equal(t, "ZAPP_CLIENT", r.URL.Query().Get("type"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[
				{"id": "1", "name": "ZAPP Settings", "clientCertificateType": "ZAPP_CLIENT"}
			]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/clientSetting?type=ZAPP_CLIENT")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all client settings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/clientSetting/all")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "all-settings",
				"name": "All Client Settings"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/clientSetting/all")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create client settings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-cs-456",
				"name": "New Settings"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/clientSetting", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE client settings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/clientSetting", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestClientSettings_ClientTypes tests valid client types
func TestClientSettings_ClientTypes(t *testing.T) {
	t.Parallel()

	t.Run("Valid client setting types", func(t *testing.T) {
		validTypes := []string{
			"ZAPP_CLIENT",
			"ISOLATION_CLIENT",
			"APP_PROTECTION",
		}

		for _, clientType := range validTypes {
			settings := ClientSettings{
				ID:                    "cs-" + clientType,
				Name:                  clientType + " Settings",
				ClientCertificateType: clientType,
			}

			data, err := json.Marshal(settings)
			require.NoError(t, err)

			assert.Contains(t, string(data), clientType)
		}
	})
}

