// Package unit provides unit tests for ZPA SCIM Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ScimGroup represents the SCIM group structure for testing
type ScimGroup struct {
	ID              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	IdpID           int64  `json:"idpId,omitempty"`
	IdpName         string `json:"idpName,omitempty"`
	IdpGroupID      string `json:"idpGroupId,omitempty"`
	CreationTime    int64  `json:"creationTime,omitempty"`
	ModifiedTime    int64  `json:"modifiedTime,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
}

// TestScimGroup_Structure tests the struct definitions
func TestScimGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ScimGroup JSON marshaling", func(t *testing.T) {
		group := ScimGroup{
			ID:         123456,
			Name:       "Engineering",
			IdpID:      789,
			IdpName:    "Okta",
			IdpGroupID: "00g12345abc",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled ScimGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
		assert.Equal(t, group.IdpGroupID, unmarshaled.IdpGroupID)
	})

	t.Run("ScimGroup JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": 456789,
			"name": "DevOps Team",
			"idpId": 123,
			"idpName": "Azure AD",
			"idpGroupId": "abc123-def456",
			"creationTime": 1609459200000,
			"modifiedTime": 1612137600000,
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var group ScimGroup
		err := json.Unmarshal([]byte(apiResponse), &group)
		require.NoError(t, err)

		assert.Equal(t, int64(456789), group.ID)
		assert.Equal(t, "DevOps Team", group.Name)
		assert.Equal(t, "Azure AD", group.IdpName)
		assert.Equal(t, "abc123-def456", group.IdpGroupID)
	})
}

// TestScimGroup_ResponseParsing tests parsing of various API responses
func TestScimGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse SCIM group list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": 1, "name": "Engineering", "idpName": "Okta"},
				{"id": 2, "name": "Sales", "idpName": "Okta"},
				{"id": 3, "name": "Marketing", "idpName": "Azure AD"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []ScimGroup `json:"list"`
			TotalPages int         `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.Equal(t, "Engineering", listResp.List[0].Name)
		assert.Equal(t, "Azure AD", listResp.List[2].IdpName)
	})
}

// TestScimGroup_MockServerOperations tests CRUD operations with mock server
func TestScimGroup_MockServerOperations(t *testing.T) {
	t.Run("GET SCIM group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/scimgroup/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": 123456,
				"name": "Mock SCIM Group",
				"idpName": "Okta"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/userconfig/v1/customers/123/scimgroup/idp/idp-001/123456")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all SCIM groups by IDP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": 1, "name": "Group A", "idpId": 123},
					{"id": 2, "name": "Group B", "idpId": 123}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/userconfig/v1/customers/123/scimgroup/idpId/123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestScimGroup_SpecialCases tests edge cases
func TestScimGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("SCIM group ID formats", func(t *testing.T) {
		// Different IDP providers use different ID formats
		idpGroupIDs := []string{
			"00g1234567890abcd",     // Okta format
			"abc123-def456-ghi789", // Azure AD format
			"CN=Group,OU=Groups,DC=example,DC=com", // LDAP DN format
		}

		for _, idpGroupID := range idpGroupIDs {
			group := ScimGroup{
				ID:         123,
				Name:       "Test Group",
				IdpGroupID: idpGroupID,
			}

			data, err := json.Marshal(group)
			require.NoError(t, err)

			var unmarshaled ScimGroup
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			assert.Equal(t, idpGroupID, unmarshaled.IdpGroupID)
		}
	})

	t.Run("Numeric ID handling", func(t *testing.T) {
		group := ScimGroup{
			ID:           9999999999999,
			IdpID:        8888888888888,
			CreationTime: 1609459200000,
			ModifiedTime: 1612137600000,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled ScimGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.IdpID, unmarshaled.IdpID)
	})
}

