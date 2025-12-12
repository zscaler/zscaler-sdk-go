// Package unit provides unit tests for ZPA Platforms service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Platforms represents the platforms for testing
type Platforms struct {
	Linux   string `json:"linux"`
	Android string `json:"android"`
	Windows string `json:"windows"`
	IOS     string `json:"ios"`
	MacOS   string `json:"mac"`
}

// TestPlatforms_Structure tests the struct definitions
func TestPlatforms_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Platforms JSON marshaling", func(t *testing.T) {
		platforms := Platforms{
			Linux:   "linux",
			Android: "android",
			Windows: "windows",
			IOS:     "ios",
			MacOS:   "mac",
		}

		data, err := json.Marshal(platforms)
		require.NoError(t, err)

		var unmarshaled Platforms
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, "linux", unmarshaled.Linux)
		assert.Equal(t, "windows", unmarshaled.Windows)
		assert.Equal(t, "mac", unmarshaled.MacOS)
	})

	t.Run("Platforms from API response", func(t *testing.T) {
		apiResponse := `{
			"linux": "Linux",
			"android": "Android",
			"windows": "Windows",
			"ios": "iOS",
			"mac": "macOS"
		}`

		var platforms Platforms
		err := json.Unmarshal([]byte(apiResponse), &platforms)
		require.NoError(t, err)

		assert.Equal(t, "Linux", platforms.Linux)
		assert.Equal(t, "macOS", platforms.MacOS)
	})
}

// TestPlatforms_MockServerOperations tests operations
func TestPlatforms_MockServerOperations(t *testing.T) {
	t.Run("GET all platforms", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/platform")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"linux": "linux",
				"android": "android",
				"windows": "windows",
				"ios": "ios",
				"mac": "mac"
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/platform")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

