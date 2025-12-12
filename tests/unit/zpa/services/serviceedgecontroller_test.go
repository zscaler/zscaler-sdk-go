// Package unit provides unit tests for ZPA Service Edge Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgecontroller"
)

func TestServiceEdgeController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ServiceEdgeController JSON marshaling", func(t *testing.T) {
		edge := serviceedgecontroller.ServiceEdgeController{
			Description:          "Test Description",
			Enabled:              true,
			ServiceEdgeGroupID:   "seg-001",
			ServiceEdgeGroupName: "Test Group",
			ControlChannelStatus: "ZPN_STATUS_ONLINE",
			CurrentVersion:       "24.1.0",
		}

		data, err := json.Marshal(edge)
		require.NoError(t, err)

		var unmarshaled serviceedgecontroller.ServiceEdgeController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, edge.Description, unmarshaled.Description)
		assert.True(t, unmarshaled.Enabled)
	})
}

func TestServiceEdgeController_MockServerOperations(t *testing.T) {
	t.Run("GET service edge by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"enabled": true}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serviceEdge")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestServiceEdgeController_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("BulkDeleteRequest", func(t *testing.T) {
		req := serviceedgecontroller.BulkDeleteRequest{
			IDs: []string{"se-1", "se-2", "se-3"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		var unmarshaled serviceedgecontroller.BulkDeleteRequest
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.IDs, 3)
	})
}
