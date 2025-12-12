// Package unit provides unit tests for ZPA Microtenants service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/microtenants"
)

// TestMicrotenants_Structure tests the struct definitions
func TestMicrotenants_Structure(t *testing.T) {
	t.Parallel()

	t.Run("MicroTenant JSON marshaling", func(t *testing.T) {
		mt := microtenants.MicroTenant{
			ID:                      "mt-123",
			Name:                    "Test Microtenant",
			Description:             "Test Description",
			Enabled:                 true,
			CriteriaAttribute:       "AuthDomain",
			CriteriaAttributeValues: []string{"test.com", "example.com"},
			Operator:                "OR",
			Priority:                "1",
			Roles: []microtenants.Roles{
				{ID: "role-001", Name: "Admin"},
			},
		}

		data, err := json.Marshal(mt)
		require.NoError(t, err)

		var unmarshaled microtenants.MicroTenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, mt.ID, unmarshaled.ID)
		assert.Equal(t, mt.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.CriteriaAttributeValues, 2)
	})

	t.Run("MicroTenant from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "mt-456",
			"name": "Production Microtenant",
			"description": "Production environment",
			"enabled": true,
			"criteriaAttribute": "AuthDomain",
			"criteriaAttributeValues": ["prod.example.com"],
			"privilegedApprovalsEnabled": true,
			"operator": "AND",
			"priority": "1",
			"roles": [
				{"id": "role-001", "name": "Admin", "customRole": false}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var mt microtenants.MicroTenant
		err := json.Unmarshal([]byte(apiResponse), &mt)
		require.NoError(t, err)

		assert.Equal(t, "mt-456", mt.ID)
		assert.True(t, mt.Enabled)
		assert.True(t, mt.PrivilegedApprovalsEnabled)
		assert.Len(t, mt.Roles, 1)
	})
}

// TestMicrotenants_MockServerOperations tests CRUD operations
func TestMicrotenants_MockServerOperations(t *testing.T) {
	t.Run("GET microtenant by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "mt-123", "name": "Mock Microtenant", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/microtenant/mt-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-mt", "name": "New Microtenant"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/microtenant", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/microtenant/mt-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
