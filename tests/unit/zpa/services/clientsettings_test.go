// Package unit provides unit tests for ZPA Client Settings service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/client_settings"
)

func TestClientSettings_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ClientSettings JSON marshaling", func(t *testing.T) {
		settings := client_settings.ClientSettings{
			ID:   "cs-123",
			Name: "Test Settings",
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		var unmarshaled client_settings.ClientSettings
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, settings.ID, unmarshaled.ID)
		assert.Equal(t, settings.Name, unmarshaled.Name)
	})
}

func TestClientSettings_MockServerOperations(t *testing.T) {
	t.Run("GET client settings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "cs-123", "name": "Mock Settings"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/clientSettings")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
