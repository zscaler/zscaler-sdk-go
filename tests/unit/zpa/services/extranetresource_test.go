// Package unit provides unit tests for ZPA Extranet Resource service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ExtranetResourceSummary represents the extranet resource summary for testing
type ExtranetResourceSummary struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestExtranetResource_Structure tests the struct definitions
func TestExtranetResource_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ExtranetResourceSummary JSON marshaling", func(t *testing.T) {
		resource := ExtranetResourceSummary{
			ID:   "er-123",
			Name: "Partner Extranet",
		}

		data, err := json.Marshal(resource)
		require.NoError(t, err)

		var unmarshaled ExtranetResourceSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, resource.ID, unmarshaled.ID)
		assert.Equal(t, resource.Name, unmarshaled.Name)
	})

	t.Run("ExtranetResourceSummary from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "er-456",
			"name": "Vendor Extranet Resource"
		}`

		var resource ExtranetResourceSummary
		err := json.Unmarshal([]byte(apiResponse), &resource)
		require.NoError(t, err)

		assert.Equal(t, "er-456", resource.ID)
		assert.Equal(t, "Vendor Extranet Resource", resource.Name)
	})
}

// TestExtranetResource_ResponseParsing tests parsing of API responses
func TestExtranetResource_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse extranet resource list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Partner A"},
				{"id": "2", "name": "Partner B"},
				{"id": "3", "name": "Partner C"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []ExtranetResourceSummary `json:"list"`
			TotalPages int                       `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
	})
}

// TestExtranetResource_MockServerOperations tests operations
func TestExtranetResource_MockServerOperations(t *testing.T) {
	t.Run("GET all extranet resource partners", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/extranetResource/partner")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Partner A"},
					{"id": "2", "name": "Partner B"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/extranetResource/partner")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

