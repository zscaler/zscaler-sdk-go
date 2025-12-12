// Package unit provides unit tests for ZPA Cloud Connector service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloud_connector"
)

func TestCloudConnector_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudConnector JSON marshaling", func(t *testing.T) {
		connector := cloud_connector.CloudConnector{
			ID:          "cc-123",
			Name:        "Test Cloud Connector",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(connector)
		require.NoError(t, err)

		var unmarshaled cloud_connector.CloudConnector
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, connector.ID, unmarshaled.ID)
		assert.Equal(t, connector.Name, unmarshaled.Name)
	})
}

func TestCloudConnector_MockServerOperations(t *testing.T) {
	t.Run("GET cloud connector", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "cc-123", "name": "Mock Connector"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cloudConnector")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
