// Package unit provides unit tests for ZPA Cloud Connector Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CloudConnectorGroup represents the cloud connector group for testing
type CloudConnectorGroup struct {
	ID              string                    `json:"id,omitempty"`
	Name            string                    `json:"name,omitempty"`
	Description     string                    `json:"description,omitempty"`
	Enabled         bool                      `json:"enabled,omitempty"`
	GeolocationID   string                    `json:"geoLocationId,omitempty"`
	ZiaCloud        string                    `json:"ziaCloud,omitempty"`
	ZiaOrgid        string                    `json:"ziaOrgId,omitempty"`
	ZnfGroupType    string                    `json:"znfGroupType,omitempty"`
	CloudConnectors []CloudConnectorInGroup   `json:"cloudConnectors,omitempty"`
	CreationTime    string                    `json:"creationTime,omitempty"`
	ModifiedBy      string                    `json:"modifiedBy,omitempty"`
	ModifiedTime    string                    `json:"modifiedTime,omitempty"`
}

// CloudConnectorInGroup represents a cloud connector within a group
type CloudConnectorInGroup struct {
	ID              string                 `json:"id,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Enabled         bool                   `json:"enabled,omitempty"`
	Fingerprint     string                 `json:"fingerprint,omitempty"`
	IPACL           []string               `json:"ipAcl,omitempty"`
	IssuedCertID    string                 `json:"issuedCertId,omitempty"`
	SigningCert     map[string]interface{} `json:"signingCert,omitempty"`
	MicroTenantID   string                 `json:"microtenantId,omitempty"`
	MicroTenantName string                 `json:"microtenantName,omitempty"`
	CreationTime    string                 `json:"creationTime,omitempty"`
	ModifiedBy      string                 `json:"modifiedBy,omitempty"`
	ModifiedTime    string                 `json:"modifiedTime,omitempty"`
}

// TestCloudConnectorGroup_Structure tests the struct definitions
func TestCloudConnectorGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudConnectorGroup JSON marshaling", func(t *testing.T) {
		ccg := CloudConnectorGroup{
			ID:            "ccg-123",
			Name:          "AWS Cloud Connector Group",
			Description:   "Cloud connector group for AWS",
			Enabled:       true,
			GeolocationID: "geo-001",
			ZiaCloud:      "zscaler.net",
			ZiaOrgid:      "org-001",
			ZnfGroupType:  "CLOUD_CONNECTOR",
			CloudConnectors: []CloudConnectorInGroup{
				{
					ID:      "cc-001",
					Name:    "AWS Connector 1",
					Enabled: true,
				},
			},
		}

		data, err := json.Marshal(ccg)
		require.NoError(t, err)

		var unmarshaled CloudConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ccg.ID, unmarshaled.ID)
		assert.Equal(t, ccg.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.CloudConnectors, 1)
	})

	t.Run("CloudConnectorGroup from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "ccg-456",
			"name": "Azure Cloud Connector Group",
			"description": "Cloud connector group for Azure",
			"enabled": true,
			"geoLocationId": "geo-002",
			"ziaCloud": "zscalerone.net",
			"ziaOrgId": "org-002",
			"znfGroupType": "CLOUD_CONNECTOR",
			"cloudConnectors": [
				{
					"id": "cc-002",
					"name": "Azure Connector 1",
					"enabled": true,
					"fingerprint": "ABC123",
					"ipAcl": ["10.0.0.0/8"],
					"microtenantId": "mt-001",
					"microtenantName": "Production"
				},
				{
					"id": "cc-003",
					"name": "Azure Connector 2",
					"enabled": true,
					"fingerprint": "DEF456"
				}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var ccg CloudConnectorGroup
		err := json.Unmarshal([]byte(apiResponse), &ccg)
		require.NoError(t, err)

		assert.Equal(t, "ccg-456", ccg.ID)
		assert.Equal(t, "Azure Cloud Connector Group", ccg.Name)
		assert.True(t, ccg.Enabled)
		assert.Len(t, ccg.CloudConnectors, 2)
		assert.Equal(t, "mt-001", ccg.CloudConnectors[0].MicroTenantID)
	})
}

// TestCloudConnectorGroup_ResponseParsing tests parsing of API responses
func TestCloudConnectorGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse cloud connector group list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Group 1", "enabled": true, "cloudConnectors": [{"id": "cc-1"}]},
				{"id": "2", "name": "Group 2", "enabled": true, "cloudConnectors": []},
				{"id": "3", "name": "Group 3", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []CloudConnectorGroup `json:"list"`
			TotalPages int                   `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Len(t, listResp.List[0].CloudConnectors, 1)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestCloudConnectorGroup_MockServerOperations tests operations
func TestCloudConnectorGroup_MockServerOperations(t *testing.T) {
	t.Run("GET cloud connector group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/cloudConnectorGroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "ccg-123",
				"name": "Mock Group",
				"enabled": true,
				"cloudConnectors": []
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnectorGroup/ccg-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all cloud connector groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Group A"},
					{"id": "2", "name": "Group B"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnectorGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET cloud connector group summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/summary")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Summary Group A"},
					{"id": "2", "name": "Summary Group B"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnectorGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestCloudConnectorGroup_SpecialCases tests edge cases
func TestCloudConnectorGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Cloud connector group with ZIA integration", func(t *testing.T) {
		ccg := CloudConnectorGroup{
			ID:       "ccg-123",
			Name:     "ZIA Integrated Group",
			ZiaCloud: "zscaler.net",
			ZiaOrgid: "12345",
		}

		data, err := json.Marshal(ccg)
		require.NoError(t, err)

		var unmarshaled CloudConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "zscaler.net", unmarshaled.ZiaCloud)
		assert.Equal(t, "12345", unmarshaled.ZiaOrgid)
	})

	t.Run("Cloud connector with signing cert", func(t *testing.T) {
		cc := CloudConnectorInGroup{
			ID:           "cc-123",
			Name:         "Signed Connector",
			IssuedCertID: "cert-001",
			SigningCert: map[string]interface{}{
				"id":   "sc-001",
				"name": "Signing Certificate",
			},
		}

		data, err := json.Marshal(cc)
		require.NoError(t, err)

		var unmarshaled CloudConnectorInGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.NotNil(t, unmarshaled.SigningCert)
	})

	t.Run("ZNF group types", func(t *testing.T) {
		groupTypes := []string{"CLOUD_CONNECTOR", "EDGE_CONNECTOR"}

		for _, groupType := range groupTypes {
			ccg := CloudConnectorGroup{
				ID:           "ccg-" + groupType,
				Name:         groupType + " Group",
				ZnfGroupType: groupType,
			}

			data, err := json.Marshal(ccg)
			require.NoError(t, err)

			assert.Contains(t, string(data), groupType)
		}
	})
}

