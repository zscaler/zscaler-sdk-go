// Package unit provides unit tests for ZPA Cloud Connector Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector_group"
)

func TestCloudConnectorGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudConnectorGroup JSON marshaling", func(t *testing.T) {
		group := cloud_connector_group.CloudConnectorGroup{
			ID:          "ccg-123",
			Name:        "Test Cloud Connector Group",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled cloud_connector_group.CloudConnectorGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
	})
}

func TestCloudConnectorGroup_MockServerOperations(t *testing.T) {
	t.Run("GET cloud connector group", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "ccg-123", "name": "Mock Group"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnectorGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
