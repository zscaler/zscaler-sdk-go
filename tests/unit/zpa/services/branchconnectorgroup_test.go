// Package unit provides unit tests for ZPA Branch Connector Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// BranchConnectorGroupSummary represents the branch connector group summary for testing
type BranchConnectorGroupSummary struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// TestBranchConnectorGroup_Structure tests the struct definitions
func TestBranchConnectorGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BranchConnectorGroupSummary JSON marshaling", func(t *testing.T) {
		summary := BranchConnectorGroupSummary{
			ID:          "bcg-123",
			Name:        "Main Branch Group",
			Description: "Primary branch connector group",
		}

		data, err := json.Marshal(summary)
		require.NoError(t, err)

		var unmarshaled BranchConnectorGroupSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, summary.ID, unmarshaled.ID)
		assert.Equal(t, summary.Name, unmarshaled.Name)
	})

	t.Run("BranchConnectorGroupSummary from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "bcg-456",
			"name": "Remote Branch Group",
			"description": "Group for remote branch connectors"
		}`

		var summary BranchConnectorGroupSummary
		err := json.Unmarshal([]byte(apiResponse), &summary)
		require.NoError(t, err)

		assert.Equal(t, "bcg-456", summary.ID)
		assert.Equal(t, "Remote Branch Group", summary.Name)
	})
}

// TestBranchConnectorGroup_ResponseParsing tests parsing of API responses
func TestBranchConnectorGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse branch connector group summary list", func(t *testing.T) {
		response := `{
			"list": [
				{"id": "1", "name": "Group 1"},
				{"id": "2", "name": "Group 2"},
				{"id": "3", "name": "Group 3"}
			],
			"totalPages": 1
		}`

		type ListResponse struct {
			List       []BranchConnectorGroupSummary `json:"list"`
			TotalPages int                           `json:"totalPages"`
		}

		var listResp ListResponse
		err := json.Unmarshal([]byte(response), &listResp)
		require.NoError(t, err)

		assert.Len(t, listResp.List, 3)
	})
}

// TestBranchConnectorGroup_MockServerOperations tests operations
func TestBranchConnectorGroup_MockServerOperations(t *testing.T) {
	t.Run("GET branch connector group summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/branchConnectorGroup/summary")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"list": [
					{"id": "1", "name": "Summary Group A"},
					{"id": "2", "name": "Summary Group B"}
				],
				"totalPages": 1
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/branchConnectorGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

