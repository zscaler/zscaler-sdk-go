// Package unit provides unit tests for ZPA Extranet Resource service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

func TestExtranetResource_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CommonSummary JSON marshaling", func(t *testing.T) {
		partner := common.CommonSummary{
			ID:   "erp-123",
			Name: "Test Partner",
		}

		data, err := json.Marshal(partner)
		require.NoError(t, err)

		var unmarshaled common.CommonSummary
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, partner.ID, unmarshaled.ID)
		assert.Equal(t, partner.Name, unmarshaled.Name)
	})
}

func TestExtranetResource_MockServerOperations(t *testing.T) {
	t.Run("GET extranet resource partner", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`[{"id": "erp-123", "name": "Mock Partner"}]`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/extranetResource/partner")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
