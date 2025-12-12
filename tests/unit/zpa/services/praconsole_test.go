// Package unit provides unit tests for ZPA PRA Console service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

func TestPRAConsole_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PRAConsole JSON marshaling", func(t *testing.T) {
		console := praconsole.PRAConsole{
			ID:          "pc-123",
			Name:        "Test Console",
			Description: "Test Description",
			Enabled:     true,
			IconText:    "base64icon==",
		}

		data, err := json.Marshal(console)
		require.NoError(t, err)

		var unmarshaled praconsole.PRAConsole
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, console.ID, unmarshaled.ID)
		assert.Equal(t, console.Name, unmarshaled.Name)
	})
}

func TestPRAConsole_MockServerOperations(t *testing.T) {
	t.Run("GET console by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "pc-123", "name": "Mock Console"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praConsole")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
