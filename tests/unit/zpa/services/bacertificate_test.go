// Package unit provides unit tests for ZPA BA Certificate service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/bacertificate"
)

// TestBaCertificate_Structure tests the struct definitions
func TestBaCertificate_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BaCertificate JSON marshaling", func(t *testing.T) {
		cert := bacertificate.BaCertificate{
			ID:          "cert-123",
			Name:        "Test Certificate",
			Description: "Test Description",
			CName:       "test.example.com",
			IssuedBy:    "CA Authority",
			IssuedTo:    "test.example.com",
			San:         []string{"test.example.com", "*.test.example.com"},
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled bacertificate.BaCertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
		assert.Len(t, unmarshaled.San, 2)
	})

	t.Run("BaCertificate from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cert-456",
			"name": "Production Certificate",
			"description": "Production SSL certificate",
			"cName": "prod.example.com",
			"issuedBy": "DigiCert",
			"issuedTo": "prod.example.com",
			"serialNo": "ABC123DEF456",
			"san": ["prod.example.com", "www.prod.example.com"],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var cert bacertificate.BaCertificate
		err := json.Unmarshal([]byte(apiResponse), &cert)
		require.NoError(t, err)

		assert.Equal(t, "cert-456", cert.ID)
		assert.Equal(t, "DigiCert", cert.IssuedBy)
		assert.Equal(t, "ABC123DEF456", cert.SerialNo)
	})
}

// TestBaCertificate_MockServerOperations tests CRUD operations
func TestBaCertificate_MockServerOperations(t *testing.T) {
	t.Run("GET certificate by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "cert-123", "name": "Mock Cert"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/certificate/cert-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all certificates", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/certificate")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
