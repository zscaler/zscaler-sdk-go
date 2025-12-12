// Package unit provides unit tests for ZPA Enrollment Certificate service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/enrollmentcert"
)

// TestEnrollmentCert_Structure tests the struct definitions
func TestEnrollmentCert_Structure(t *testing.T) {
	t.Parallel()

	t.Run("EnrollmentCert JSON marshaling", func(t *testing.T) {
		cert := enrollmentcert.EnrollmentCert{
			ID:             "cert-123",
			Name:           "Test Enrollment Cert",
			Description:    "Test Description",
			Cname:          "enroll.example.com",
			ClientCertType: "CONNECTOR",
			AllowSigning:   true,
			IssuedBy:       "CA Authority",
			IssuedTo:       "enroll.example.com",
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled enrollmentcert.EnrollmentCert
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.AllowSigning)
	})

	t.Run("EnrollmentCert from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cert-456",
			"name": "Production Enrollment Cert",
			"description": "Production enrollment certificate",
			"cName": "enroll.prod.example.com",
			"clientCertType": "SERVICE_EDGE",
			"allowSigning": true,
			"issuedBy": "Zscaler CA",
			"issuedTo": "enroll.prod.example.com",
			"parentCertId": "parent-001",
			"parentCertName": "Root CA",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var cert enrollmentcert.EnrollmentCert
		err := json.Unmarshal([]byte(apiResponse), &cert)
		require.NoError(t, err)

		assert.Equal(t, "cert-456", cert.ID)
		assert.Equal(t, "SERVICE_EDGE", cert.ClientCertType)
		assert.True(t, cert.AllowSigning)
	})
}

// TestEnrollmentCert_MockServerOperations tests CRUD operations
func TestEnrollmentCert_MockServerOperations(t *testing.T) {
	t.Run("GET enrollment cert by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "cert-123", "name": "Mock Cert"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/enrollmentCert/cert-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all enrollment certs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/enrollmentCert")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestEnrollmentCert_SpecialCases tests edge cases
func TestEnrollmentCert_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Client cert types", func(t *testing.T) {
		types := []string{"CONNECTOR", "SERVICE_EDGE", "NP_ASSISTANT"}

		for _, certType := range types {
			cert := enrollmentcert.EnrollmentCert{
				ID:             "cert-" + certType,
				Name:           certType + " Cert",
				ClientCertType: certType,
			}

			data, err := json.Marshal(cert)
			require.NoError(t, err)
			assert.Contains(t, string(data), certType)
		}
	})
}
