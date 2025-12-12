// Package unit provides unit tests for ZPA Workload Tag Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// WorkloadTagGroupSummary represents the workload tag group summary for testing
type WorkloadTagGroupSummary struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// TestWorkloadTagGroup_Structure tests the struct definitions
func TestWorkloadTagGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WorkloadTagGroupSummary JSON marshaling", func(t *testing.T) {
		tag := WorkloadTagGroupSummary{
			ID:   "wtg-123",
			Name: "Production Workloads",
		}

		data, err := json.Marshal(tag)
		require.NoError(t, err)

		var unmarshaled WorkloadTagGroupSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, tag.ID, unmarshaled.ID)
		assert.Equal(t, tag.Name, unmarshaled.Name)
	})

	t.Run("WorkloadTagGroupSummary from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "wtg-456",
			"name": "Development Workloads"
		}`

		var tag WorkloadTagGroupSummary
		err := json.Unmarshal([]byte(apiResponse), &tag)
		require.NoError(t, err)

		assert.Equal(t, "wtg-456", tag.ID)
		assert.Equal(t, "Development Workloads", tag.Name)
	})
}

// TestWorkloadTagGroup_ResponseParsing tests parsing of API responses
func TestWorkloadTagGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse workload tag group list response", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Production"},
				{"id": "2", "name": "Staging"},
				{"id": "3", "name": "Development"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []WorkloadTagGroupSummary `json:"list"`
			TotalPages int                       `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
	})
}

// TestWorkloadTagGroup_MockServerOperations tests operations
func TestWorkloadTagGroup_MockServerOperations(t *testing.T) {
	t.Run("GET all workload tag groups", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/workloadTagGroup/summary")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Tag Group 1"},
					{"id": "2", "name": "Tag Group 2"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/workloadTagGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// TestWorkloadTagGroup_SpecialCases tests edge cases
func TestWorkloadTagGroup_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Common workload tag group names", func(t *testing.T) {
		names := []string{
			"Production",
			"Staging",
			"Development",
			"QA",
			"Testing",
			"Sandbox",
		}

		for _, name := range names {
			tag := WorkloadTagGroupSummary{
				ID:   "wtg-" + name,
				Name: name,
			}

			data, err := json.Marshal(tag)
			require.NoError(t, err)

			assert.Contains(t, string(data), name)
		}
	})
}

