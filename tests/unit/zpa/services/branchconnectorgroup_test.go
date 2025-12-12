// Package unit provides unit tests for ZPA Branch Connector Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestBranchConnectorGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CommonSummary JSON marshaling", func(t *testing.T) {
		summary := common.CommonSummary{
			ID:   "bcg-123",
			Name: "Test Group",
		}

		data, err := json.Marshal(summary)
		require.NoError(t, err)

		var unmarshaled common.CommonSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, summary.ID, unmarshaled.ID)
		assert.Equal(t, summary.Name, unmarshaled.Name)
	})
}

func TestBranchConnectorGroup_MockServerOperations(t *testing.T) {
	t.Run("GET branch connector group summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "bcg-123", "name": "Mock Group"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/branchConnectorGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
