// Package unit provides unit tests for ZPA CBI Regions service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/cloudbrowserisolation/cbiregions"
)

func TestCBIRegions_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CBIRegions JSON marshaling", func(t *testing.T) {
		region := cbiregions.CBIRegions{
			ID:   "region-123",
			Name: "US West",
		}

		data, err := json.Marshal(region)
		require.NoError(t, err)

		var unmarshaled cbiregions.CBIRegions
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, region.ID, unmarshaled.ID)
		assert.Equal(t, region.Name, unmarshaled.Name)
	})
}

func TestCBIRegions_MockServerOperations(t *testing.T) {
	t.Run("GET all regions", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "1", "name": "US West"}, {"id": "2", "name": "US East"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/cbiRegions")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
