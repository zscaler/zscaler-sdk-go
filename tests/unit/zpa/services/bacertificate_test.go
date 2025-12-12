// Package unit provides unit tests for ZPA Browser Access Certificate service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BACertificate represents the browser access certificate structure for testing
type BACertificate struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	Description         string `json:"description,omitempty"`
	Certificate         string `json:"certificate,omitempty"`
	CertBlob            string `json:"certBlob,omitempty"`
	IssuedBy            string `json:"issuedBy,omitempty"`
	IssuedTo            string `json:"issuedTo,omitempty"`
	SerialNo            string `json:"serialNo,omitempty"`
	ValidFromInEpochSec string `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec   string `json:"validToInEpochSec,omitempty"`
	Status              string `json:"status,omitempty"`
	CName               string `json:"cName,omitempty"`
	CertChain           string `json:"certChain,omitempty"`
	PublicKey           string `json:"publicKey,omitempty"`
	San                 string `json:"san,omitempty"`
	CreationTime        string `json:"creationTime,omitempty"`
	ModifiedBy          string `json:"modifiedBy,omitempty"`
	ModifiedTime        string `json:"modifiedTime,omitempty"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
	MicroTenantName     string `json:"microtenantName,omitempty"`
}

// TestBACertificate_Structure tests the struct definitions
func TestBACertificate_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BACertificate JSON marshaling", func(t *testing.T) {
		cert := BACertificate{
			ID:          "ba-123",
			Name:        "webapp.example.com",
			Description: "Certificate for web application",
			IssuedBy:    "DigiCert",
			IssuedTo:    "Example Corp",
			SerialNo:    "1234567890ABCDEF",
			Status:      "ACTIVE",
			CName:       "webapp.example.com",
			San:         "webapp.example.com,www.webapp.example.com",
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled BACertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cert.ID, unmarshaled.ID)
		assert.Equal(t, cert.Name, unmarshaled.Name)
		assert.Equal(t, cert.Status, unmarshaled.Status)
		assert.Equal(t, cert.CName, unmarshaled.CName)
	})

	t.Run("BACertificate JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ba-456",
			"name": "api.example.com",
			"description": "API endpoint certificate",
			"issuedBy": "Let's Encrypt",
			"issuedTo": "api.example.com",
			"serialNo": "FEDCBA0987654321",
			"validFromInEpochSec": "1609459200",
			"validToInEpochSec": "1704153600",
			"status": "ACTIVE",
			"cName": "api.example.com",
			"san": "api.example.com,api-backup.example.com",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var cert BACertificate
		err := json.Unmarshal([]byte(apiResponse), &cert)
		require.NoError(t, err)

		assert.Equal(t, "ba-456", cert.ID)
		assert.Equal(t, "api.example.com", cert.Name)
		assert.Equal(t, "ACTIVE", cert.Status)
		assert.Contains(t, cert.San, "api-backup.example.com")
	})
}

// TestBACertificate_ResponseParsing tests parsing of various API responses
func TestBACertificate_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse BA certificate list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "app1.example.com", "status": "ACTIVE"},
				{"id": "2", "name": "app2.example.com", "status": "ACTIVE"},
				{"id": "3", "name": "expired.example.com", "status": "EXPIRED"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []BACertificate `json:"list"`
			TotalPages int             `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "ACTIVE", listResp.List[0].Status)
		assert.Equal(t, "EXPIRED", listResp.List[2].Status)
	})
}

// TestBACertificate_MockServerOperations tests CRUD operations with mock server
func TestBACertificate_MockServerOperations(t *testing.T) {
	t.Run("GET BA certificate by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/clientlessCertificate/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ba-123",
				"name": "Mock Certificate",
				"status": "ACTIVE"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/clientlessCertificate/ba-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all BA certificates", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Cert A", "status": "ACTIVE"},
					{"id": "2", "name": "Cert B", "status": "ACTIVE"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/clientlessCertificate")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST upload BA certificate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-ba-456",
				"name": "new-cert.example.com",
				"status": "ACTIVE"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/clientlessCertificate", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE BA certificate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/clientlessCertificate/ba-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestBACertificate_SpecialCases tests edge cases
func TestBACertificate_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Certificate status types", func(t *testing.T) {
		statuses := []string{
			"ACTIVE",
			"EXPIRED",
			"EXPIRING_SOON",
			"REVOKED",
			"PENDING",
		}

		for _, status := range statuses {
			cert := BACertificate{
				ID:     "ba-" + status,
				Name:   status + ".example.com",
				Status: status,
			}

			data, err := json.Marshal(cert)
			require.NoError(t, err)

			assert.Contains(t, string(data), status)
		}
	})

	t.Run("Multiple SANs", func(t *testing.T) {
		cert := BACertificate{
			ID:   "ba-123",
			Name: "multi-san.example.com",
			San:  "multi-san.example.com,www.multi-san.example.com,api.multi-san.example.com,*.multi-san.example.com",
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled BACertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Contains(t, unmarshaled.San, "*.multi-san.example.com")
	})

	t.Run("Wildcard certificate", func(t *testing.T) {
		cert := BACertificate{
			ID:    "ba-123",
			Name:  "*.example.com",
			CName: "*.example.com",
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		assert.Contains(t, string(data), "*.example.com")
	})

	t.Run("Certificate validity period", func(t *testing.T) {
		cert := BACertificate{
			ID:                  "ba-123",
			Name:                "validity-test.example.com",
			ValidFromInEpochSec: "1609459200",  // 2021-01-01
			ValidToInEpochSec:   "1735689600",  // 2025-01-01
		}

		data, err := json.Marshal(cert)
		require.NoError(t, err)

		var unmarshaled BACertificate
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "1609459200", unmarshaled.ValidFromInEpochSec)
		assert.Equal(t, "1735689600", unmarshaled.ValidToInEpochSec)
	})
}

