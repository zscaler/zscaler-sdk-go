// Package unit provides unit tests for ZPA Isolation Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/isolationprofile"
)

func TestIsolationProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IsolationProfile JSON marshaling", func(t *testing.T) {
		profile := isolationprofile.IsolationProfile{
			ID:          "iso-123",
			Name:        "Test Isolation Profile",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled isolationprofile.IsolationProfile
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
	})
}

func TestIsolationProfile_MockServerOperations(t *testing.T) {
	t.Run("GET isolation profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "iso-123", "name": "Mock Profile"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/isolationProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
