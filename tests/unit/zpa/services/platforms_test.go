// Package unit provides unit tests for ZPA Platforms service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/platforms"
)

func TestPlatforms_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Platforms JSON marshaling", func(t *testing.T) {
		platform := platforms.Platforms{
			Linux:   "linux",
			Android: "android",
			Windows: "windows",
			IOS:     "ios",
			MacOS:   "mac",
		}

		data, err := json.Marshal(platform)
		require.NoError(t, err)

		var unmarshaled platforms.Platforms
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "linux", unmarshaled.Linux)
		assert.Equal(t, "windows", unmarshaled.Windows)
	})
}

func TestPlatforms_MockServerOperations(t *testing.T) {
	t.Run("GET all platforms", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"linux": "linux", "windows": "windows"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/platforms")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
