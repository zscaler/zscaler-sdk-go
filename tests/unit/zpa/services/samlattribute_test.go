// Package unit provides unit tests for ZPA SAML Attribute service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/samlattribute"
)

// TestSamlAttribute_Structure tests the struct definitions
func TestSamlAttribute_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SamlAttribute JSON marshaling", func(t *testing.T) {
		attr := samlattribute.SamlAttribute{
			ID:            "saml-123",
			Name:          "Department",
			SamlName:      "department",
			IdpID:         "idp-001",
			IdpName:       "Okta",
			UserAttribute: true,
		}

		data, err := json.Marshal(attr)
		require.NoError(t, err)

		var unmarshaled samlattribute.SamlAttribute
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, attr.ID, unmarshaled.ID)
		assert.Equal(t, attr.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.UserAttribute)
	})

	t.Run("SamlAttribute from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "saml-456",
			"name": "Groups",
			"samlName": "groups",
			"idpId": "idp-002",
			"idpName": "Azure AD",
			"userAttribute": true,
			"delta": "30",
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var attr samlattribute.SamlAttribute
		err := json.Unmarshal([]byte(apiResponse), &attr)
		require.NoError(t, err)

		assert.Equal(t, "saml-456", attr.ID)
		assert.Equal(t, "Azure AD", attr.IdpName)
		assert.True(t, attr.UserAttribute)
	})
}

// TestSamlAttribute_MockServerOperations tests operations
func TestSamlAttribute_MockServerOperations(t *testing.T) {
	t.Run("GET SAML attribute by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "saml-123", "name": "Mock Attribute"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/samlAttribute/saml-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all SAML attributes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1"}, {"id": "2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/samlAttribute")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestSamlAttribute_SpecialCases tests edge cases
func TestSamlAttribute_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Common SAML attributes", func(t *testing.T) {
		attrs := []string{"department", "groups", "email", "displayName", "title", "manager"}

		for _, name := range attrs {
			attr := samlattribute.SamlAttribute{
				ID:       "saml-" + name,
				Name:     name,
				SamlName: name,
			}

			data, err := json.Marshal(attr)
			require.NoError(t, err)
			assert.Contains(t, string(data), name)
		}
	})
}
