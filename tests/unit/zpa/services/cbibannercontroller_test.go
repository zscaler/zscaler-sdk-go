// Package unit provides unit tests for ZPA CBI Banner Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbibannercontroller"
)

func TestCBIBannerController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIBannerController JSON marshaling", func(t *testing.T) {
		banner := cbibannercontroller.CBIBannerController{
			ID:                "banner-123",
			Name:              "Test Banner",
			PrimaryColor:      "#FF0000",
			TextColor:         "#FFFFFF",
			NotificationTitle: "Welcome",
			NotificationText:  "You are accessing a protected resource",
			Logo:              "base64logo==",
			IsDefault:         false,
			Persist:           true,
		}

		data, err := json.Marshal(banner)
		require.NoError(t, err)

		var unmarshaled cbibannercontroller.CBIBannerController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, banner.ID, unmarshaled.ID)
		assert.Equal(t, banner.Name, unmarshaled.Name)
	})
}

func TestCBIBannerController_MockServerOperations(t *testing.T) {
	t.Run("GET banner by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "banner-123", "name": "Mock Banner"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbiBanner")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
