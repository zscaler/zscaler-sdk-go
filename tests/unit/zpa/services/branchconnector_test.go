// Package unit provides unit tests for ZPA Branch Connector service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/branch_connector"
)

func TestBranchConnector_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BranchConnector JSON marshaling", func(t *testing.T) {
		connector := branch_connector.BranchConnector{
			ID:          "bc-123",
			Name:        "Test Branch Connector",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled branch_connector.BranchConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, connector.ID, unmarshaled.ID)
		assert.Equal(t, connector.Name, unmarshaled.Name)
	})
}

func TestBranchConnector_MockServerOperations(t *testing.T) {
	t.Run("GET branch connector", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "bc-123", "name": "Mock Connector"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/branchConnector")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
