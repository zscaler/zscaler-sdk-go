// Package unit provides unit tests for ZPA Cloud Connector service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CloudConnector represents the cloud connector for testing
type CloudConnector struct {
	ID                     string                 `json:"id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	Description            string                 `json:"description,omitempty"`
	EdgeConnectorGroupID   string                 `json:"edgeConnectorGroupId,omitempty"`
	EdgeConnectorGroupName string                 `json:"edgeConnectorGroupName,omitempty"`
	Enabled                bool                   `json:"enabled,omitempty"`
	Fingerprint            string                 `json:"fingerprint,omitempty"`
	IpAcl                  []string               `json:"ipAcl,omitempty"`
	IssuedCertID           string                 `json:"issuedCertId,omitempty"`
	CreationTime           string                 `json:"creationTime,omitempty"`
	ModifiedBy             string                 `json:"modifiedBy,omitempty"`
	ModifiedTime           int                    `json:"modifiedTime,omitempty"`
	EnrollmentCert         map[string]interface{} `json:"enrollmentCert,omitempty"`
}

// TestCloudConnector_Structure tests the struct definitions
func TestCloudConnector_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudConnector JSON marshaling", func(t *testing.T) {
		cc := CloudConnector{
			ID:                     "cc-123",
			Name:                   "AWS Cloud Connector",
			Description:            "Cloud connector for AWS VPC",
			EdgeConnectorGroupID:   "ecg-001",
			EdgeConnectorGroupName: "AWS Edge Group",
			Enabled:                true,
			Fingerprint:            "ABC123DEF456",
			IpAcl:                  []string{"10.0.0.0/8"},
			IssuedCertID:           "cert-001",
		}

		data, err := json.Marshal(cc)
		require.NoError(t, err)

		var unmarshaled CloudConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cc.ID, unmarshaled.ID)
		assert.Equal(t, cc.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("CloudConnector from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "cc-456",
			"name": "Azure Cloud Connector",
			"description": "Cloud connector for Azure VNet",
			"edgeConnectorGroupId": "ecg-002",
			"edgeConnectorGroupName": "Azure Edge Group",
			"enabled": true,
			"fingerprint": "XYZ789GHI012",
			"ipAcl": ["172.16.0.0/12"],
			"issuedCertId": "cert-002",
			"creationTime": "1609459200000",
			"modifiedTime": 1612137600000,
			"modifiedBy": "admin@example.com",
			"enrollmentCert": {
				"id": "ec-001",
				"name": "Cloud Cert"
			}
		}`

		var cc CloudConnector
		err := json.Unmarshal([]byte(apiResponse), &cc)
		require.NoError(t, err)

		assert.Equal(t, "cc-456", cc.ID)
		assert.Equal(t, "Azure Cloud Connector", cc.Name)
		assert.True(t, cc.Enabled)
		assert.NotNil(t, cc.EnrollmentCert)
	})
}

// TestCloudConnector_ResponseParsing tests parsing of API responses
func TestCloudConnector_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse cloud connector list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "AWS CC", "enabled": true},
				{"id": "2", "name": "Azure CC", "enabled": true},
				{"id": "3", "name": "GCP CC", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []CloudConnector `json:"list"`
			TotalPages int              `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestCloudConnector_MockServerOperations tests operations
func TestCloudConnector_MockServerOperations(t *testing.T) {
	t.Run("GET all cloud connectors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/cloudConnector")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Cloud Connector A", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnector")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET cloud connector by name (search)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "AWS CC", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnector?search=AWS")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestCloudConnector_SpecialCases tests edge cases
func TestCloudConnector_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Cloud connector with multiple IP ACLs", func(t *testing.T) {
		cc := CloudConnector{
			ID:      "cc-123",
			Name:    "Multi-ACL CC",
			Enabled: true,
			IpAcl:   []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		}

		data, err := json.Marshal(cc)
		require.NoError(t, err)

		var unmarshaled CloudConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.IpAcl, 3)
	})

	t.Run("Cloud connector with enrollment cert", func(t *testing.T) {
		cc := CloudConnector{
			ID:           "cc-123",
			Name:         "Enrolled CC",
			IssuedCertID: "cert-001",
			EnrollmentCert: map[string]interface{}{
				"id":   "ec-001",
				"name": "Cloud Enrollment Cert",
			},
		}

		data, err := json.Marshal(cc)
		require.NoError(t, err)

		var unmarshaled CloudConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotNil(t, unmarshaled.EnrollmentCert)
	})
}

