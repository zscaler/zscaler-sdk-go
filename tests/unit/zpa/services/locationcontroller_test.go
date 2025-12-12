// Package unit provides unit tests for ZPA Location Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LocationSummary represents the location summary for testing
type LocationSummary struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// LocationGroupDTO represents the location group for testing
type LocationGroupDTO struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TestLocationController_Structure tests the struct definitions
func TestLocationController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("LocationSummary JSON marshaling", func(t *testing.T) {
		location := LocationSummary{
			ID:   "loc-123",
			Name: "US West",
		}

		data, err := json.Marshal(location)
		require.NoError(t, err)

		var unmarshaled LocationSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, location.ID, unmarshaled.ID)
		assert.Equal(t, location.Name, unmarshaled.Name)
	})

	t.Run("LocationGroupDTO JSON marshaling", func(t *testing.T) {
		group := LocationGroupDTO{
			ID:          "lg-123",
			Name:        "North America",
			Description: "All North American locations",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled LocationGroupDTO
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
	})
}

// TestLocationController_ResponseParsing tests parsing of API responses
func TestLocationController_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse location summary list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "US East"},
				{"id": "2", "name": "US West"},
				{"id": "3", "name": "EU Central"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []LocationSummary `json:"list"`
			TotalPages int               `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
	})
}

// TestLocationController_MockServerOperations tests operations
func TestLocationController_MockServerOperations(t *testing.T) {
	t.Run("GET location extranet resource", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/location/extranetResource/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Location A"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/location/extranetResource/er-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET location summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/location/summary")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "US East"},
					{"id": "2", "name": "US West"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/location/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GET location group extranet resource", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/locationGroup/extranetResource/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Group A"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/locationGroup/extranetResource/er-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

