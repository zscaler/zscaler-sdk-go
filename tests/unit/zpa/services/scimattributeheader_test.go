// Package unit provides unit tests for ZPA SCIM Attribute Header service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ScimAttributeHeader represents the SCIM attribute header for testing
type ScimAttributeHeader struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	IdpID           string   `json:"idpId,omitempty"`
	DataType        string   `json:"dataType,omitempty"`
	SchemaURI       string   `json:"schemaURI,omitempty"`
	CanonicalValues []string `json:"canonicalValues,omitempty"`
	CaseSensitive   bool     `json:"caseSensitive,omitempty"`
	MultiValued     bool     `json:"multivalued,omitempty"`
	Mutability      string   `json:"mutability,omitempty"`
	Required        bool     `json:"required,omitempty"`
	Returned        string   `json:"returned,omitempty"`
	Uniqueness      bool     `json:"uniqueness,omitempty"`
	CreationTime    string   `json:"creationTime,omitempty"`
	ModifiedTime    string   `json:"modifiedTime,omitempty"`
}

// TestScimAttributeHeader_Structure tests the struct definitions
func TestScimAttributeHeader_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ScimAttributeHeader JSON marshaling", func(t *testing.T) {
		attr := ScimAttributeHeader{
			ID:              "scim-123",
			Name:            "department",
			Description:     "User department",
			IdpID:           "idp-001",
			DataType:        "String",
			SchemaURI:       "urn:ietf:params:scim:schemas:core:2.0:User",
			CanonicalValues: []string{"Engineering", "Sales", "Marketing"},
			CaseSensitive:   false,
			MultiValued:     false,
			Mutability:      "readWrite",
			Required:        false,
		}

		data, err := json.Marshal(attr)
		require.NoError(t, err)

		var unmarshaled ScimAttributeHeader
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, attr.ID, unmarshaled.ID)
		assert.Equal(t, attr.Name, unmarshaled.Name)
		assert.Len(t, unmarshaled.CanonicalValues, 3)
	})

	t.Run("ScimAttributeHeader from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "scim-456",
			"name": "userName",
			"description": "Unique identifier for the user",
			"idpId": "idp-002",
			"dataType": "String",
			"schemaURI": "urn:ietf:params:scim:schemas:core:2.0:User",
			"caseSensitive": true,
			"multivalued": false,
			"mutability": "readOnly",
			"required": true,
			"returned": "always",
			"uniqueness": true,
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var attr ScimAttributeHeader
		err := json.Unmarshal([]byte(apiResponse), &attr)
		require.NoError(t, err)

		assert.Equal(t, "scim-456", attr.ID)
		assert.Equal(t, "userName", attr.Name)
		assert.True(t, attr.Required)
		assert.True(t, attr.Uniqueness)
	})
}

// TestScimAttributeHeader_MockServerOperations tests operations
func TestScimAttributeHeader_MockServerOperations(t *testing.T) {
	t.Run("GET SCIM attribute by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/scimattribute/")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "scim-123", "name": "department"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/idp/idp-001/scimattribute/scim-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET all SCIM attributes by IdP ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": [{"id": "1", "name": "attr1"}, {"id": "2", "name": "attr2"}], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/idp/idp-001/scimattribute")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET SCIM attribute values", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"list": ["value1", "value2", "value3"], "totalPages": 1}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/scimattribute/idpId/idp-001/attributeId/scim-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestScimAttributeHeader_SpecialCases tests edge cases
func TestScimAttributeHeader_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Common SCIM attributes", func(t *testing.T) {
		attrs := []string{
			"userName",
			"emails",
			"name.familyName",
			"name.givenName",
			"displayName",
			"groups",
			"department",
			"title",
		}

		for _, name := range attrs {
			attr := ScimAttributeHeader{
				ID:   "scim-" + name,
				Name: name,
			}

			data, err := json.Marshal(attr)
			require.NoError(t, err)

			assert.Contains(t, string(data), name)
		}
	})

	t.Run("Data types", func(t *testing.T) {
		dataTypes := []string{"String", "Boolean", "Integer", "DateTime", "Complex"}

		for _, dt := range dataTypes {
			attr := ScimAttributeHeader{
				ID:       "scim-dt",
				Name:     "test",
				DataType: dt,
			}

			data, err := json.Marshal(attr)
			require.NoError(t, err)

			assert.Contains(t, string(data), dt)
		}
	})
}

