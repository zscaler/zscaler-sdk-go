// Package unit provides unit tests for ZPA Enrollment Certificate service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// EnrollmentCert represents the enrollment certificate structure for testing
type EnrollmentCert struct {
	ID                        string `json:"id,omitempty"`
	Name                      string `json:"name,omitempty"`
	Description               string `json:"description,omitempty"`
	Certificate               string `json:"certificate,omitempty"`
	ClientCertType            string `json:"clientCertType,omitempty"`
	IssuedBy                  string `json:"issuedBy,omitempty"`
	IssuedTo                  string `json:"issuedTo,omitempty"`
	SerialNo                  string `json:"serialNo,omitempty"`
	ValidFromInEpochSec       string `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec         string `json:"validToInEpochSec,omitempty"`
	CertChain                 string `json:"certChain,omitempty"`
	AllowSigning              bool   `json:"allowSigning"`
	ZrsaEncryptedPrivateKey   string `json:"zrsaEncryptedPrivateKey,omitempty"`
	ZrsaEncryptedSessionKey   string `json:"zrsaEncryptedSessionKey,omitempty"`
	CreationTime              string `json:"creationTime,omitempty"`
	ModifiedBy                string `json:"modifiedBy,omitempty"`
	ModifiedTime              string `json:"modifiedTime,omitempty"`
	MicroTenantID             string `json:"microtenantId,omitempty"`
	MicroTenantName           string `json:"microtenantName,omitempty"`
}

// TestEnrollmentCert_Structure tests the struct definitions
func TestEnrollmentCert_Structure(t *testing.T) {
	t.Parallel()

	t.Run("EnrollmentCert JSON marshaling", func(t *testing.T) {
		cert := EnrollmentCert{
			ID:             "ec-123",
			Name:           "Connector Certificate",
			Description:    "Certificate for App Connectors",
			ClientCertType: "CONNECTOR",
			IssuedBy:       "Zscaler Root CA",
			IssuedTo:       "App Connector",
			SerialNo:       "1234567890",
			AllowSigning:   true,
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled EnrollmentCert
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
		assert.Equal(t, cert.ClientCertType, unmarshaled.ClientCertType)
		assert.True(t, unmarshaled.AllowSigning)
	})

	t.Run("EnrollmentCert JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ec-456",
			"name": "Service Edge Certificate",
			"description": "Certificate for Service Edges",
			"clientCertType": "SERVICE_EDGE",
			"issuedBy": "Zscaler Root CA",
			"issuedTo": "Service Edge",
			"serialNo": "9876543210",
			"validFromInEpochSec": "1609459200",
			"validToInEpochSec": "1704153600",
			"allowSigning": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var cert EnrollmentCert
		err := json.Unmarshal([]byte(apiResponse), &cert)
		require.NoError(t, err)

		assert.Equal(t, "ec-456", cert.ID)
		assert.Equal(t, "Service Edge Certificate", cert.Name)
		assert.Equal(t, "SERVICE_EDGE", cert.ClientCertType)
		assert.True(t, cert.AllowSigning)
		assert.NotEmpty(t, cert.ValidFromInEpochSec)
		assert.NotEmpty(t, cert.ValidToInEpochSec)
	})
}

// TestEnrollmentCert_ResponseParsing tests parsing of various API responses
func TestEnrollmentCert_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse enrollment cert list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Connector Cert", "clientCertType": "CONNECTOR", "allowSigning": true},
				{"id": "2", "name": "Service Edge Cert", "clientCertType": "SERVICE_EDGE", "allowSigning": true},
				{"id": "3", "name": "Browser Access Cert", "clientCertType": "BROWSER_ACCESS", "allowSigning": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []EnrollmentCert `json:"list"`
			TotalPages int              `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "CONNECTOR", listResp.List[0].ClientCertType)
		assert.False(t, listResp.List[2].AllowSigning)
	})
}

// TestEnrollmentCert_MockServerOperations tests CRUD operations with mock server
func TestEnrollmentCert_MockServerOperations(t *testing.T) {
	t.Run("GET enrollment cert by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/enrollmentCert/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ec-123",
				"name": "Mock Certificate",
				"clientCertType": "CONNECTOR",
				"allowSigning": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/enrollmentCert/ec-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all enrollment certs", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Cert A", "clientCertType": "CONNECTOR"},
					{"id": "2", "name": "Cert B", "clientCertType": "SERVICE_EDGE"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/enrollmentCert")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestEnrollmentCert_SpecialCases tests edge cases
func TestEnrollmentCert_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Client cert types", func(t *testing.T) {
		certTypes := []string{
			"CONNECTOR",
			"SERVICE_EDGE",
			"BROWSER_ACCESS",
			"BRANCH_CONNECTOR",
			"CLIENT_CONNECTOR",
			"CLOUD_CONNECTOR",
			"ISOLATION_CLIENT",
		}

		for _, certType := range certTypes {
			cert := EnrollmentCert{
				ID:             "ec-" + certType,
				Name:           certType + " Certificate",
				ClientCertType: certType,
			}

			data, err := json.Marshal(cert)
			require.NoError(t, err)

			assert.Contains(t, string(data), certType)
		}
	})

	t.Run("Certificate with signing disabled", func(t *testing.T) {
		cert := EnrollmentCert{
			ID:           "ec-123",
			Name:         "Read-only Certificate",
			AllowSigning: false,
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"allowSigning":false`)
	})

	t.Run("Certificate validity period", func(t *testing.T) {
		cert := EnrollmentCert{
			ID:                  "ec-123",
			Name:                "Validity Test Cert",
			ValidFromInEpochSec: "1609459200",  // 2021-01-01
			ValidToInEpochSec:   "1735689600",  // 2025-01-01
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled EnrollmentCert
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "1609459200", unmarshaled.ValidFromInEpochSec)
		assert.Equal(t, "1735689600", unmarshaled.ValidToInEpochSec)
	})
}

