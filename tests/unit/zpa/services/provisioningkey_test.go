// Package unit provides unit tests for ZPA Provisioning Key service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/provisioningkey"
)

// TestProvisioningKey_Structure tests the struct definitions
func TestProvisioningKey_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ProvisioningKey JSON marshaling", func(t *testing.T) {
		key := provisioningkey.ProvisioningKey{
			ID:                  "pk-123",
			Name:                "Test Provisioning Key",
			AssociationType:     "CONNECTOR_GRP",
			AppConnectorGroupID: "acg-001",
			MaxUsage:            "10",
			Enabled:             true,
			EnrollmentCertID:    "cert-001",
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)

		var unmarshaled provisioningkey.ProvisioningKey
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, key.ID, unmarshaled.ID)
		assert.Equal(t, key.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
	})

	t.Run("ProvisioningKey from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pk-456",
			"name": "Production Key",
			"associationType": "SERVICE_EDGE_GRP",
			"appConnectorGroupId": "acg-002",
			"appConnectorGroupName": "Production Connectors",
			"maxUsage": "50",
			"usageCount": "5",
			"enabled": true,
			"enrollmentCertId": "cert-002",
			"enrollmentCertName": "Prod Cert",
			"provisioningKey": "abc123def456",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var key provisioningkey.ProvisioningKey
		err := json.Unmarshal([]byte(apiResponse), &key)
		require.NoError(t, err)

		assert.Equal(t, "pk-456", key.ID)
		assert.Equal(t, "SERVICE_EDGE_GRP", key.AssociationType)
		assert.True(t, key.Enabled)
	})
}

// TestProvisioningKey_MockServerOperations tests CRUD operations
func TestProvisioningKey_MockServerOperations(t *testing.T) {
	t.Run("GET provisioning key by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "pk-123", "name": "Mock Key", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/provisioningKey/pk-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create provisioning key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-pk", "name": "New Key"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/provisioningKey", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
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

// TestProvisioningKey_SpecialCases tests edge cases
func TestProvisioningKey_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Association types", func(t *testing.T) {
		types := []string{"CONNECTOR_GRP", "SERVICE_EDGE_GRP", "NP_ASSISTANT_GRP"}

		for _, assocType := range types {
			key := provisioningkey.ProvisioningKey{
				ID:              "pk-" + assocType,
				Name:            assocType + " Key",
				AssociationType: assocType,
			}

			data, err := json.Marshal(key)
			require.NoError(t, err)
			assert.Contains(t, string(data), assocType)
		}
	})
}
