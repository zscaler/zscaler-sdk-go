// Package unit provides unit tests for ZPA Branch Connector service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BranchConnector represents the branch connector for testing
type BranchConnector struct {
	ID                       string                 `json:"id,omitempty"`
	Name                     string                 `json:"name,omitempty"`
	Description              string                 `json:"description,omitempty"`
	BranchConnectorGroupID   string                 `json:"branchConnectorGroupId,omitempty"`
	BranchConnectorGroupName string                 `json:"branchConnectorGroupName,omitempty"`
	EdgeConnectorGroupID     string                 `json:"edgeConnectorGroupId,omitempty"`
	EdgeConnectorGroupName   string                 `json:"edgeConnectorGroupName,omitempty"`
	Enabled                  bool                   `json:"enabled,omitempty"`
	Fingerprint              string                 `json:"fingerprint,omitempty"`
	IpAcl                    []string               `json:"ipAcl,omitempty"`
	IssuedCertID             string                 `json:"issuedCertId,omitempty"`
	CreationTime             string                 `json:"creationTime,omitempty"`
	ModifiedBy               string                 `json:"modifiedBy,omitempty"`
	ModifiedTime             string                 `json:"modifiedTime,omitempty"`
	EnrollmentCert           map[string]interface{} `json:"enrollmentCert,omitempty"`
}

// TestBranchConnector_Structure tests the struct definitions
func TestBranchConnector_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BranchConnector JSON marshaling", func(t *testing.T) {
		bc := BranchConnector{
			ID:                       "bc-123",
			Name:                     "Branch Office Connector",
			Description:              "Main branch office connector",
			BranchConnectorGroupID:   "bcg-001",
			BranchConnectorGroupName: "Branch Group 1",
			EdgeConnectorGroupID:     "ecg-001",
			EdgeConnectorGroupName:   "Edge Group 1",
			Enabled:                  true,
			Fingerprint:              "ABC123DEF456",
			IpAcl:                    []string{"10.0.0.0/8", "192.168.0.0/16"},
			IssuedCertID:             "cert-001",
		}

		data, err := json.Marshal(bc)
		require.NoError(t, err)

		var unmarshaled BranchConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, bc.ID, unmarshaled.ID)
		assert.Equal(t, bc.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.IpAcl, 2)
	})

	t.Run("BranchConnector from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "bc-456",
			"name": "Remote Office Connector",
			"description": "Remote office branch connector",
			"branchConnectorGroupId": "bcg-002",
			"branchConnectorGroupName": "Remote Branches",
			"edgeConnectorGroupId": "ecg-002",
			"edgeConnectorGroupName": "Remote Edges",
			"enabled": true,
			"fingerprint": "XYZ789GHI012",
			"ipAcl": ["172.16.0.0/12"],
			"issuedCertId": "cert-002",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"enrollmentCert": {
				"id": "ec-001",
				"name": "Branch Cert"
			}
		}`

		var bc BranchConnector
		err := json.Unmarshal([]byte(apiResponse), &bc)
		require.NoError(t, err)

		assert.Equal(t, "bc-456", bc.ID)
		assert.Equal(t, "Remote Office Connector", bc.Name)
		assert.True(t, bc.Enabled)
		assert.NotNil(t, bc.EnrollmentCert)
	})
}

// TestBranchConnector_ResponseParsing tests parsing of API responses
func TestBranchConnector_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse branch connector list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Branch 1", "enabled": true},
				{"id": "2", "name": "Branch 2", "enabled": true},
				{"id": "3", "name": "Branch 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []BranchConnector `json:"list"`
			TotalPages int               `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestBranchConnector_MockServerOperations tests operations with mock server
func TestBranchConnector_MockServerOperations(t *testing.T) {
	t.Run("GET all branch connectors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/branchConnector")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Branch A", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/branchConnector")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET branch connector by name (search)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			search := r.URL.Query().Get("search")
			assert.NotEmpty(t, search)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Main Branch", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/branchConnector?search=Main")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestBranchConnector_SpecialCases tests edge cases
func TestBranchConnector_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Branch connector with IP ACLs", func(t *testing.T) {
		bc := BranchConnector{
			ID:      "bc-123",
			Name:    "ACL Branch",
			Enabled: true,
			IpAcl:   []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
		}

		data, err := json.Marshal(bc)
		require.NoError(t, err)

		var unmarshaled BranchConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.IpAcl, 3)
	})

	t.Run("Disabled branch connector", func(t *testing.T) {
		bc := BranchConnector{
			ID:      "bc-123",
			Name:    "Disabled Branch",
			Enabled: false,
		}

		data, err := json.Marshal(bc)
		require.NoError(t, err)

		// With omitempty, false won't appear
		var unmarshaled BranchConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.False(t, unmarshaled.Enabled)
	})
}

