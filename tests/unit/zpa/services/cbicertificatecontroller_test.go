// Package unit provides unit tests for ZPA CBI Certificate Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CBICertificate represents the CBI certificate for testing
type CBICertificate struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	PEM       string `json:"pem,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

// TestCBICertificate_Structure tests the struct definitions
func TestCBICertificate_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBICertificate JSON marshaling", func(t *testing.T) {
		cert := CBICertificate{
			ID:        "cert-123",
			Name:      "Corporate Root CA",
			PEM:       "-----BEGIN CERTIFICATE-----\nMIIDXTCC...\n-----END CERTIFICATE-----",
			IsDefault: false,
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled CBICertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
		assert.Contains(t, unmarshaled.PEM, "BEGIN CERTIFICATE")
	})

	t.Run("CBICertificate from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cert-456",
			"name": "Intermediate CA",
			"pem": "-----BEGIN CERTIFICATE-----\nMIIEpDCCA...\n-----END CERTIFICATE-----",
			"isDefault": true
		}`

		var cert CBICertificate
		err := json.Unmarshal([]byte(apiResponse), &cert)
		require.NoError(t, err)

		assert.Equal(t, "cert-456", cert.ID)
		assert.Equal(t, "Intermediate CA", cert.Name)
		assert.True(t, cert.IsDefault)
	})
}

// TestCBICertificate_ResponseParsing tests parsing of API responses
func TestCBICertificate_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse certificate list response", func(t *testing.T) {
		response := `[
			{"id": "1", "name": "Root CA", "isDefault": true},
			{"id": "2", "name": "Intermediate CA", "isDefault": false},
			{"id": "3", "name": "Custom CA", "isDefault": false}
		]`

		var certs []CBICertificate
		err := json.Unmarshal([]byte(response), &certs)
		require.NoError(t, err)

		assert.Len(t, certs, 3)
		assert.True(t, certs[0].IsDefault)
	})
}

// TestCBICertificate_MockServerOperations tests CRUD operations
func TestCBICertificate_MockServerOperations(t *testing.T) {
	t.Run("GET certificate by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/certificates/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "cert-123", "name": "Mock Cert", "isDefault": false}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/certificates/cert-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all certificates", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `[{"id": "1", "name": "Cert A"}, {"id": "2", "name": "Cert B"}]`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbi/api/customers/123/certificates")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create certificate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-cert", "name": "New Certificate"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/cbi/api/customers/123/certificate", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update certificate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/certificates/cert-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE certificate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/certificates/cert-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

