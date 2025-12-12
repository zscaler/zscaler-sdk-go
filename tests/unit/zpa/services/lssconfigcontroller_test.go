// Package unit provides unit tests for ZPA LSS Config Controller service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func TestLSSConfigController_Structure(t *testing.T) {
	t.Parallel()

	t.Run("LSSResource JSON marshaling", func(t *testing.T) {
		lss := lssconfigcontroller.LSSResource{
			ID: "lss-123",
			LSSConfig: &lssconfigcontroller.LSSConfig{
				Name:       "Test LSS Config",
				Enabled:    true,
				Format:     "JSON",
				LSSHost:    "lss.example.com",
				LSSPort:    "11000",
			},
		}

		data, err := json.Marshal(lss)
		require.NoError(t, err)

		var unmarshaled lssconfigcontroller.LSSResource
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, lss.ID, unmarshaled.ID)
	})
}

func TestLSSConfigController_MockServerOperations(t *testing.T) {
	t.Run("GET LSS config by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "lss-123"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/lssConfig")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
