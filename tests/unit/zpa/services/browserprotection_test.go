// Package unit provides unit tests for ZPA Browser Protection service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

func TestBrowserProtection_Structure(t *testing.T) {
	t.Parallel()

	t.Run("BrowserProtection JSON marshaling", func(t *testing.T) {
		profile := browser_protection.BrowserProtection{
			ID:          "bp-123",
			Name:        "Test Browser Protection",
			Description: "Test Description",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled browser_protection.BrowserProtection
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
	})
}

func TestBrowserProtection_MockServerOperations(t *testing.T) {
	t.Run("GET browser protection profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "bp-123", "name": "Mock Profile"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/browserProtection")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
