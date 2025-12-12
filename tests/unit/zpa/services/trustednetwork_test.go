// Package unit provides unit tests for ZPA Trusted Network service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TrustedNetworkResource represents the trusted network structure for testing
type TrustedNetworkResource struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	NetworkID         string `json:"networkId,omitempty"`
	CreationTime      string `json:"creationTime,omitempty"`
	ModifiedBy        string `json:"modifiedBy,omitempty"`
	ModifiedTime      string `json:"modifiedTime,omitempty"`
	MicroTenantID     string `json:"microtenantId,omitempty"`
	MicroTenantName   string `json:"microtenantName,omitempty"`
	ZscalerCloud      string `json:"zscalerCloud,omitempty"`
	MasterCustomerID  string `json:"masterCustomerId,omitempty"`
}

// TestTrustedNetwork_Structure tests the struct definitions
func TestTrustedNetwork_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TrustedNetworkResource JSON marshaling", func(t *testing.T) {
		network := TrustedNetworkResource{
			ID:               "tn-123",
			Name:             "Corporate Network",
			NetworkID:        "corp-network-id-12345",
			MicroTenantID:    "mt-001",
			MicroTenantName:  "Production",
			ZscalerCloud:     "zscaler.net",
			MasterCustomerID: "customer-123",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		var unmarshaled TrustedNetworkResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, network.ID, unmarshaled.ID)
		assert.Equal(t, network.Name, unmarshaled.Name)
		assert.Equal(t, network.NetworkID, unmarshaled.NetworkID)
	})

	t.Run("TrustedNetworkResource JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "tn-456",
			"name": "Branch Office Network",
			"networkId": "branch-network-id-67890",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-002",
			"microtenantName": "Branch",
			"zscalerCloud": "zscaler.net",
			"masterCustomerId": "customer-456"
		}`

		var network TrustedNetworkResource
		err := json.Unmarshal([]byte(apiResponse), &network)
		require.NoError(t, err)

		assert.Equal(t, "tn-456", network.ID)
		assert.Equal(t, "Branch Office Network", network.Name)
		assert.Equal(t, "branch-network-id-67890", network.NetworkID)
		assert.NotEmpty(t, network.CreationTime)
		assert.NotEmpty(t, network.ZscalerCloud)
	})
}

// TestTrustedNetwork_ResponseParsing tests parsing of various API responses
func TestTrustedNetwork_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse trusted network list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Network A", "networkId": "net-a-12345"},
				{"id": "2", "name": "Network B", "networkId": "net-b-67890"},
				{"id": "3", "name": "Network C", "networkId": "net-c-11111"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []TrustedNetworkResource `json:"list"`
			TotalPages int                      `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "Network A", listResp.List[0].Name)
		assert.Equal(t, "net-a-12345", listResp.List[0].NetworkID)
	})
}

// TestTrustedNetwork_MockServerOperations tests CRUD operations with mock server
func TestTrustedNetwork_MockServerOperations(t *testing.T) {
	t.Run("GET trusted network by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/trustedNetwork/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "tn-123",
				"name": "Mock Trusted Network",
				"networkId": "mock-network-id"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/trustedNetwork/tn-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all trusted networks", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Network A", "networkId": "net-a"},
					{"id": "2", "name": "Network B", "networkId": "net-b"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/trustedNetwork")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestTrustedNetwork_ErrorHandling tests error scenarios
func TestTrustedNetwork_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Trusted Network Not Found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code": "NOT_FOUND", "message": "Trusted network not found"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/trustedNetwork/nonexistent")
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// TestTrustedNetwork_SpecialCases tests edge cases
func TestTrustedNetwork_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Network ID format validation", func(t *testing.T) {
		// Network IDs typically follow a specific format
		network := TrustedNetworkResource{
			ID:        "tn-123",
			Name:      "Formatted Network",
			NetworkID: "ABC123-DEF456-GHI789",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		var unmarshaled TrustedNetworkResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "ABC123-DEF456-GHI789", unmarshaled.NetworkID)
	})

	t.Run("Microtenant configuration", func(t *testing.T) {
		network := TrustedNetworkResource{
			ID:              "tn-123",
			Name:            "Tenant Network",
			MicroTenantID:   "mt-001",
			MicroTenantName: "Tenant One",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"microtenantId":"mt-001"`)
		assert.Contains(t, string(data), `"microtenantName":"Tenant One"`)
	})

	t.Run("Zscaler cloud configuration", func(t *testing.T) {
		clouds := []string{"zscaler.net", "zscalerone.net", "zscalertwo.net", "zscloud.net"}

		for _, cloud := range clouds {
			network := TrustedNetworkResource{
				ID:           "tn-" + cloud,
				Name:         cloud + " Network",
				ZscalerCloud: cloud,
			}

			data, err := json.Marshal(network)
			require.NoError(t, err)

			assert.Contains(t, string(data), cloud)
		}
	})
}

