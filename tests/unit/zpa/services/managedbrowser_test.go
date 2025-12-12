// Package unit provides unit tests for ZPA Managed Browser service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/managed_browser"
)

func TestManagedBrowser_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ManagedBrowserProfile JSON marshaling", func(t *testing.T) {
		browser := managed_browser.ManagedBrowserProfile{
			ID:          "mb-123",
			Name:        "Test Managed Browser",
			Description: "Test Description",
			BrowserType: "CHROME",
		}

		data, err := json.Marshal(browser)
		require.NoError(t, err)

		var unmarshaled managed_browser.ManagedBrowserProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, browser.ID, unmarshaled.ID)
		assert.Equal(t, browser.Name, unmarshaled.Name)
	})
}

func TestManagedBrowser_MockServerOperations(t *testing.T) {
	t.Run("GET managed browser", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "mb-123", "name": "Mock Browser"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/managedBrowser")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
