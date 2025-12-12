// Package unit provides unit tests for ZPA Workload Tag Group service
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

func TestWorkloadTagGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CommonSummary JSON marshaling", func(t *testing.T) {
		group := common.CommonSummary{
			ID:   "wtg-123",
			Name: "Production Tags",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled common.CommonSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
	})
}

func TestWorkloadTagGroup_MockServerOperations(t *testing.T) {
	t.Run("GET workload tag group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "wtg-123", "name": "Mock Group"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/workloadTagGroup/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
