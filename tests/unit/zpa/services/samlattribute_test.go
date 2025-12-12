// Package unit provides unit tests for ZPA SAML Attribute service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SamlAttribute represents the SAML attribute structure for testing
type SamlAttribute struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	IdpID           string `json:"idpId,omitempty"`
	IdpName         string `json:"idpName,omitempty"`
	SamlName        string `json:"samlName,omitempty"`
	UserAttribute   bool   `json:"userAttribute"`
	CreationTime    string `json:"creationTime,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
}

// TestSamlAttribute_Structure tests the struct definitions
func TestSamlAttribute_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SamlAttribute JSON marshaling", func(t *testing.T) {
		attr := SamlAttribute{
			ID:            "sa-123",
			Name:          "department",
			IdpID:         "idp-001",
			IdpName:       "Okta",
			SamlName:      "urn:oid:2.16.840.1.113730.3.1.3",
			UserAttribute: true,
		}

		data, err := json.Marshal(attr)
		require.NoError(t, err)

		var unmarshaled SamlAttribute
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, attr.ID, unmarshaled.ID)
		assert.Equal(t, attr.Name, unmarshaled.Name)
		assert.Equal(t, attr.SamlName, unmarshaled.SamlName)
		assert.True(t, unmarshaled.UserAttribute)
	})

	t.Run("SamlAttribute JSON unmarshaling from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "sa-456",
			"name": "email",
			"idpId": "idp-002",
			"idpName": "Azure AD",
			"samlName": "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress",
			"userAttribute": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000",
			"modifiedBy": "admin@example.com",
			"microtenantId": "mt-001",
			"microtenantName": "Production"
		}`

		var attr SamlAttribute
		err := json.Unmarshal([]byte(apiResponse), &attr)
		require.NoError(t, err)

		assert.Equal(t, "sa-456", attr.ID)
		assert.Equal(t, "email", attr.Name)
		assert.True(t, attr.UserAttribute)
		assert.Contains(t, attr.SamlName, "emailaddress")
	})
}

// TestSamlAttribute_ResponseParsing tests parsing of various API responses
func TestSamlAttribute_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse SAML attribute list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "department", "samlName": "dept", "userAttribute": true},
				{"id": "2", "name": "email", "samlName": "email", "userAttribute": true},
				{"id": "3", "name": "group", "samlName": "memberOf", "userAttribute": false}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []SamlAttribute `json:"list"`
			TotalPages int             `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
		assert.True(t, listResp.List[0].UserAttribute)
		assert.False(t, listResp.List[2].UserAttribute)
	})
}

// TestSamlAttribute_MockServerOperations tests CRUD operations with mock server
func TestSamlAttribute_MockServerOperations(t *testing.T) {
	t.Run("GET SAML attribute by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/samlAttribute/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "sa-123",
				"name": "Mock SAML Attribute",
				"userAttribute": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v2/admin/customers/123/samlAttribute/sa-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET SAML attributes by IDP ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/idp/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Attr A", "idpId": "idp-001"},
					{"id": "2", "name": "Attr B", "idpId": "idp-001"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/samlAttribute/idp/idp-001")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestSamlAttribute_SpecialCases tests edge cases
func TestSamlAttribute_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Common SAML attribute names", func(t *testing.T) {
		commonAttrs := []struct {
			name     string
			samlName string
		}{
			{"email", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"},
			{"department", "urn:oid:2.16.840.1.113730.3.1.3"},
			{"groups", "http://schemas.microsoft.com/ws/2008/06/identity/claims/groups"},
			{"displayName", "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/displayname"},
		}

		for _, attr := range commonAttrs {
			samlAttr := SamlAttribute{
				ID:       "sa-" + attr.name,
				Name:     attr.name,
				SamlName: attr.samlName,
			}

			data, err := json.Marshal(samlAttr)
			require.NoError(t, err)

			assert.Contains(t, string(data), attr.name)
		}
	})

	t.Run("Non-user attribute", func(t *testing.T) {
		attr := SamlAttribute{
			ID:            "sa-123",
			Name:          "Group Membership",
			UserAttribute: false,
		}

		data, err := json.Marshal(attr)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"userAttribute":false`)
	})
}

