// Package unit provides unit tests for ZPA Location Controller service
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

func TestLocationController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CommonSummary JSON marshaling", func(t *testing.T) {
		summary := common.CommonSummary{
			ID:   "loc-123",
			Name: "Test Location",
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

func TestLocationController_MockServerOperations(t *testing.T) {
	t.Run("GET location summary", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "loc-123", "name": "Mock Location"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/location/summary")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
