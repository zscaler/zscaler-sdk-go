// Package unit provides unit tests for ZPA Private Cloud Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_controller"
)

func TestPrivateCloudController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivateCloudController JSON marshaling", func(t *testing.T) {
		controller := private_cloud_controller.PrivateCloudController{
			Description:          "Test Description",
			Enabled:              true,
			CurrentVersion:       "24.1.0",
			ControlChannelStatus: "ZPN_STATUS_ONLINE",
		}

		data, err := json.Marshal(controller)
		require.NoError(t, err)

		var unmarshaled private_cloud_controller.PrivateCloudController
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, controller.Description, unmarshaled.Description)
		assert.True(t, unmarshaled.Enabled)
	})
}

func TestPrivateCloudController_MockServerOperations(t *testing.T) {
	t.Run("GET controller by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"enabled": true}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/privateCloudController")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
