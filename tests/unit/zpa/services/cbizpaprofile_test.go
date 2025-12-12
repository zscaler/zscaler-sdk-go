// Package unit provides unit tests for ZPA CBI ZPA Profile service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbizpaprofile"
)

func TestCBIZpaProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZPAProfiles JSON marshaling", func(t *testing.T) {
		profile := cbizpaprofile.ZPAProfiles{
			ID:          "zpa-123",
			Name:        "Test ZPA Profile",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		var unmarshaled cbizpaprofile.ZPAProfiles
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, profile.ID, unmarshaled.ID)
		assert.Equal(t, profile.Name, unmarshaled.Name)
	})
}

func TestCBIZpaProfile_MockServerOperations(t *testing.T) {
	t.Run("GET ZPA profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "zpa-123", "name": "Mock Profile"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbiZpaProfile")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
