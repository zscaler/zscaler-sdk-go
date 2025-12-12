// Package unit provides unit tests for ZPA PRA Portal service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praportal"
)

func TestPRAPortal_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PRAPortal JSON marshaling", func(t *testing.T) {
		portal := praportal.PRAPortal{
			ID:            "portal-123",
			Name:          "Test Portal",
			Description:   "Test Description",
			Enabled:       true,
			Domain:        "pra.example.com",
			CertificateID: "cert-001",
		}

		data, err := json.Marshal(portal)
		require.NoError(t, err)

		var unmarshaled praportal.PRAPortal
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, portal.ID, unmarshaled.ID)
		assert.Equal(t, portal.Name, unmarshaled.Name)
	})
}

func TestPRAPortal_MockServerOperations(t *testing.T) {
	t.Run("GET portal by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "portal-123", "name": "Mock Portal"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praPortal")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
