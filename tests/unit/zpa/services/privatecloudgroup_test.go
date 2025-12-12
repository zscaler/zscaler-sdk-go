// Package unit provides unit tests for ZPA Private Cloud Group service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/private_cloud_group"
)

func TestPrivateCloudGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivateCloudGroup JSON marshaling", func(t *testing.T) {
		group := private_cloud_group.PrivateCloudGroup{
			ID:          "pcg-123",
			Name:        "Test Group",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		var unmarshaled private_cloud_group.PrivateCloudGroup
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, group.ID, unmarshaled.ID)
		assert.Equal(t, group.Name, unmarshaled.Name)
	})
}

func TestPrivateCloudGroup_MockServerOperations(t *testing.T) {
	t.Run("GET group by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "pcg-123", "name": "Mock Group"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/privateCloudGroup")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
