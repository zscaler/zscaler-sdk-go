// Package unit provides unit tests for ZPA Provisioning Key service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ProvisioningKey represents the provisioning key structure for testing
type ProvisioningKey struct {
	ID                     string `json:"id,omitempty"`
	Name                   string `json:"name,omitempty"`
	AssociationType        string `json:"associationType,omitempty"`
	ProvisioningKey        string `json:"provisioningKey,omitempty"`
	AppConnectorGroupID    string `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName  string `json:"appConnectorGroupName,omitempty"`
	ServiceEdgeGroupID     string `json:"serviceEdgeGroupId,omitempty"`
	ServiceEdgeGroupName   string `json:"serviceEdgeGroupName,omitempty"`
	ZcomponentID           string `json:"zcomponentId,omitempty"`
	ZcomponentName         string `json:"zcomponentName,omitempty"`
	EnrollmentCertID       string `json:"enrollmentCertId,omitempty"`
	EnrollmentCertName     string `json:"enrollmentCertName,omitempty"`
	MaxUsage               string `json:"maxUsage,omitempty"`
	UsageCount             string `json:"usageCount,omitempty"`
	Enabled                bool   `json:"enabled"`
	CreationTime           string `json:"creationTime,omitempty"`
	ModifiedBy             string `json:"modifiedBy,omitempty"`
	ModifiedTime           string `json:"modifiedTime,omitempty"`
	MicroTenantID          string `json:"microtenantId,omitempty"`
	MicroTenantName        string `json:"microtenantName,omitempty"`
	IPACLExcluded          bool   `json:"ipAclExcluded"`
	UIConfig               string `json:"uiConfig,omitempty"`
}

// TestProvisioningKey_Structure tests the struct definitions
func TestProvisioningKey_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ProvisioningKey JSON marshaling for connector", func(t *testing.T) {
		key := ProvisioningKey{
			ID:                    "pk-123",
			Name:                  "Test Provisioning Key",
			AssociationType:       "CONNECTOR_GRP",
			ProvisioningKey:       "ABC123DEF456",
			AppConnectorGroupID:   "acg-001",
			AppConnectorGroupName: "Connector Group 1",
			EnrollmentCertID:      "cert-001",
			EnrollmentCertName:    "Connector Certificate",
			MaxUsage:              "10",
			UsageCount:            "3",
			Enabled:               true,
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		var unmarshaled ProvisioningKey
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, key.ID, unmarshaled.ID)
		assert.Equal(t, key.Name, unmarshaled.Name)
		assert.Equal(t, "CONNECTOR_GRP", unmarshaled.AssociationType)
		assert.Equal(t, key.ProvisioningKey, unmarshaled.ProvisioningKey)
	})

	t.Run("ProvisioningKey JSON marshaling for service edge", func(t *testing.T) {
		key := ProvisioningKey{
			ID:                   "pk-456",
			Name:                 "Service Edge Provisioning Key",
			AssociationType:      "SERVICE_EDGE_GRP",
			ProvisioningKey:      "XYZ789ABC123",
			ServiceEdgeGroupID:   "seg-001",
			ServiceEdgeGroupName: "Service Edge Group 1",
			EnrollmentCertID:     "cert-002",
			EnrollmentCertName:   "Service Edge Certificate",
			MaxUsage:             "20",
			UsageCount:           "5",
			Enabled:              true,
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		var unmarshaled ProvisioningKey
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "SERVICE_EDGE_GRP", unmarshaled.AssociationType)
		assert.Equal(t, key.ServiceEdgeGroupID, unmarshaled.ServiceEdgeGroupID)
	})

	t.Run("ProvisioningKey JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pk-789",
			"name": "Production Key",
			"associationType": "CONNECTOR_GRP",
			"provisioningKey": "PROD-KEY-12345",
			"appConnectorGroupId": "acg-002",
			"appConnectorGroupName": "Production Connectors",
			"zcomponentId": "zcomp-001",
			"zcomponentName": "Component 1",
			"enrollmentCertId": "cert-003",
			"enrollmentCertName": "Production Certificate",
			"maxUsage": "100",
			"usageCount": "25",
			"enabled": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production Tenant",
			"ipAclExcluded": false
		}`

		var key ProvisioningKey
		err := json.Unmarshal([]byte(apiResponse), &key)
		require.NoError(t, err)

		assert.Equal(t, "pk-789", key.ID)
		assert.Equal(t, "Production Key", key.Name)
		assert.Equal(t, "CONNECTOR_GRP", key.AssociationType)
		assert.Equal(t, "100", key.MaxUsage)
		assert.Equal(t, "25", key.UsageCount)
		assert.True(t, key.Enabled)
		assert.False(t, key.IPACLExcluded)
	})
}

// TestProvisioningKey_ResponseParsing tests parsing of various API responses
func TestProvisioningKey_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse provisioning key list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Key 1", "associationType": "CONNECTOR_GRP", "enabled": true, "maxUsage": "10", "usageCount": "2"},
				{"id": "2", "name": "Key 2", "associationType": "SERVICE_EDGE_GRP", "enabled": true, "maxUsage": "20", "usageCount": "5"},
				{"id": "3", "name": "Key 3", "associationType": "CONNECTOR_GRP", "enabled": false, "maxUsage": "5", "usageCount": "5"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []ProvisioningKey `json:"list"`
			TotalPages int               `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "CONNECTOR_GRP", listResp.List[0].AssociationType)
		assert.Equal(t, "SERVICE_EDGE_GRP", listResp.List[1].AssociationType)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestProvisioningKey_MockServerOperations tests CRUD operations with mock server
func TestProvisioningKey_MockServerOperations(t *testing.T) {
	t.Run("GET provisioning key by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/provisioningKey/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "pk-123",
				"name": "Mock Provisioning Key",
				"associationType": "CONNECTOR_GRP",
				"provisioningKey": "MOCK-KEY-12345",
				"enabled": true,
				"maxUsage": "50",
				"usageCount": "10"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/associationType/CONNECTOR_GRP/provisioningKey/pk-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all provisioning keys by association type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/CONNECTOR_GRP/provisioningKey")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Key A", "enabled": true},
					{"id": "2", "name": "Key B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/associationType/CONNECTOR_GRP/provisioningKey")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create provisioning key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-pk-456",
				"name": "New Provisioning Key",
				"associationType": "CONNECTOR_GRP",
				"provisioningKey": "NEW-KEY-67890",
				"enabled": true,
				"maxUsage": "25",
				"usageCount": "0"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/provisioningKey", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update provisioning key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/provisioningKey/pk-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE provisioning key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/provisioningKey/pk-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestProvisioningKey_ErrorHandling tests error scenarios
func TestProvisioningKey_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Provisioning Key not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/provisioningKey/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("400 Bad Request - Max usage exceeded", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": "INVALID_REQUEST", "message": "Max usage has been reached"}`))
		}))
		defer server.Close()

		resp, _ := http.Get(server.URL + "/provisioningKey/pk-123")
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// TestProvisioningKey_SpecialCases tests edge cases
func TestProvisioningKey_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Key at max usage", func(t *testing.T) {
		key := ProvisioningKey{
			ID:         "pk-123",
			Name:       "Maxed Out Key",
			MaxUsage:   "10",
			UsageCount: "10",
			Enabled:    true,
		}

		assert.Equal(t, key.MaxUsage, key.UsageCount)
	})

	t.Run("Disabled provisioning key", func(t *testing.T) {
		key := ProvisioningKey{
			ID:      "pk-123",
			Name:    "Disabled Key",
			Enabled: false,
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
	})

	t.Run("IP ACL excluded key", func(t *testing.T) {
		key := ProvisioningKey{
			ID:            "pk-123",
			Name:          "ACL Excluded Key",
			IPACLExcluded: true,
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ipAclExcluded":true`)
	})

	t.Run("UI config", func(t *testing.T) {
		key := ProvisioningKey{
			ID:       "pk-123",
			Name:     "UI Config Key",
			UIConfig: "NEW_APP_CONNECTOR_UI",
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"uiConfig":"NEW_APP_CONNECTOR_UI"`)
	})
}

