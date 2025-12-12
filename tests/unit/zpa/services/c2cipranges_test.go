// Package unit provides unit tests for ZPA C2C IP Ranges service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/c2c_ip_ranges"
)

func TestC2CIPRanges_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPRanges JSON marshaling", func(t *testing.T) {
		ranges := c2c_ip_ranges.IPRanges{
			ID:          "c2c-123",
			Name:        "Test IP Ranges",
			Description: "Test Description",
			Enabled:     true,
		}

		data, err := json.Marshal(ranges)
		require.NoError(t, err)

		var unmarshaled c2c_ip_ranges.IPRanges
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, ranges.ID, unmarshaled.ID)
		assert.Equal(t, ranges.Name, unmarshaled.Name)
	})
}

func TestC2CIPRanges_MockServerOperations(t *testing.T) {
	t.Run("GET C2C IP ranges", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "c2c-123", "name": "Mock Ranges"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/c2cIPRanges")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
