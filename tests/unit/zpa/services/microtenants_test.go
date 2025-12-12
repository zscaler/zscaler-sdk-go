// Package unit provides unit tests for ZPA Microtenants service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Microtenant represents the microtenant structure for testing
type Microtenant struct {
	ID                         string `json:"id,omitempty"`
	Name                       string `json:"name,omitempty"`
	Description                string `json:"description,omitempty"`
	Enabled                    bool   `json:"enabled"`
	CriteriaAttribute          string `json:"criteriaAttribute,omitempty"`
	CriteriaAttributeValues    string `json:"criteriaAttributeValues,omitempty"`
	OperatorCriteriaAttribute  string `json:"operatorCriteriaAttribute,omitempty"`
	CreationTime               string `json:"creationTime,omitempty"`
	ModifiedBy                 string `json:"modifiedBy,omitempty"`
	ModifiedTime               string `json:"modifiedTime,omitempty"`
	Privilege                  string `json:"privilege,omitempty"`
}

// MicrotenantUser represents a user in a microtenant
type MicrotenantUser struct {
	ID              string   `json:"id,omitempty"`
	Email           string   `json:"email,omitempty"`
	DisplayName     string   `json:"displayName,omitempty"`
	Username        string   `json:"username,omitempty"`
	RoleIDs         []string `json:"roleIds,omitempty"`
	MicrotenantID   string   `json:"microtenantId,omitempty"`
	MicrotenantName string   `json:"microtenantName,omitempty"`
}

// TestMicrotenant_Structure tests the struct definitions
func TestMicrotenant_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Microtenant JSON marshaling", func(t *testing.T) {
		tenant := Microtenant{
			ID:                        "mt-123",
			Name:                      "Engineering Tenant",
			Description:               "Tenant for engineering team",
			Enabled:                   true,
			CriteriaAttribute:         "AuthDomain",
			CriteriaAttributeValues:   "engineering.example.com",
			OperatorCriteriaAttribute: "EQUALS",
		}

		data, err := json.Marshal(tenant)
		require.NoError(t, err)

		var unmarshaled Microtenant
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, tenant.ID, unmarshaled.ID)
		assert.Equal(t, tenant.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "AuthDomain", unmarshaled.CriteriaAttribute)
	})

	t.Run("Microtenant JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "mt-456",
			"name": "Sales Tenant",
			"description": "Tenant for sales team",
			"enabled": true,
			"criteriaAttribute": "AuthDomain",
			"criteriaAttributeValues": "sales.example.com",
			"operatorCriteriaAttribute": "EQUALS",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"privilege": "FULL_ACCESS"
		}`

		var tenant Microtenant
		err := json.Unmarshal([]byte(apiResponse), &tenant)
		require.NoError(t, err)

		assert.Equal(t, "mt-456", tenant.ID)
		assert.Equal(t, "Sales Tenant", tenant.Name)
		assert.True(t, tenant.Enabled)
		assert.Equal(t, "FULL_ACCESS", tenant.Privilege)
	})

	t.Run("MicrotenantUser structure", func(t *testing.T) {
		user := MicrotenantUser{
			ID:              "user-001",
			Email:           "user@example.com",
			DisplayName:     "Test User",
			Username:        "testuser",
			RoleIDs:         []string{"role-1", "role-2"},
			MicrotenantID:   "mt-001",
			MicrotenantName: "Engineering",
		}

		data, err := json.Marshal(user)
		require.NoError(t, err)

		var unmarshaled MicrotenantUser
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, user.Email, unmarshaled.Email)
		assert.Len(t, unmarshaled.RoleIDs, 2)
	})
}

// TestMicrotenant_ResponseParsing tests parsing of various API responses
func TestMicrotenant_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse microtenant list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Tenant A", "enabled": true},
				{"id": "2", "name": "Tenant B", "enabled": true},
				{"id": "3", "name": "Tenant C", "enabled": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []Microtenant `json:"list"`
			TotalPages int           `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].Enabled)
		assert.False(t, listResp.List[2].Enabled)
	})
}

// TestMicrotenant_MockServerOperations tests CRUD operations with mock server
func TestMicrotenant_MockServerOperations(t *testing.T) {
	t.Run("GET microtenant by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/microtenants/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "mt-123",
				"name": "Mock Microtenant",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/microtenants/mt-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all microtenants", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Tenant A", "enabled": true},
					{"id": "2", "name": "Tenant B", "enabled": true}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/microtenants")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-mt-456",
				"name": "New Microtenant",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/microtenants", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/microtenants/mt-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE microtenant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/microtenants/mt-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestMicrotenant_SpecialCases tests edge cases
func TestMicrotenant_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Criteria attribute types", func(t *testing.T) {
		criteriaAttrs := []string{
			"AuthDomain",
			"SCIM_GROUP",
			"SAML_ATTRIBUTE",
		}

		for _, attr := range criteriaAttrs {
			tenant := Microtenant{
				ID:                "mt-" + attr,
				Name:              attr + " Tenant",
				CriteriaAttribute: attr,
			}

			data, err := json.Marshal(tenant)
			require.NoError(t, err)

			assert.Contains(t, string(data), attr)
		}
	})

	t.Run("Operator criteria attribute values", func(t *testing.T) {
		operators := []string{
			"EQUALS",
			"CONTAINS",
			"STARTS_WITH",
			"ENDS_WITH",
		}

		for _, op := range operators {
			tenant := Microtenant{
				ID:                        "mt-op-" + op,
				Name:                      op + " Tenant",
				OperatorCriteriaAttribute: op,
			}

			data, err := json.Marshal(tenant)
			require.NoError(t, err)

			assert.Contains(t, string(data), op)
		}
	})

	t.Run("Privilege levels", func(t *testing.T) {
		privileges := []string{
			"FULL_ACCESS",
			"READ_ONLY",
			"DELEGATED",
		}

		for _, priv := range privileges {
			tenant := Microtenant{
				ID:        "mt-priv-" + priv,
				Name:      priv + " Tenant",
				Privilege: priv,
			}

			data, err := json.Marshal(tenant)
			require.NoError(t, err)

			assert.Contains(t, string(data), priv)
		}
	})

	t.Run("Disabled microtenant", func(t *testing.T) {
		tenant := Microtenant{
			ID:      "mt-123",
			Name:    "Disabled Tenant",
			Enabled: false,
		}

		data, err := json.Marshal(tenant)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
	})
}

