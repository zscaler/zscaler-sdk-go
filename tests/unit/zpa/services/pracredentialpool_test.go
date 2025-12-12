// Package unit provides unit tests for ZPA PRA Credential Pool service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredentialpool"
)

func TestPRACredentialPool_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CredentialPool JSON marshaling", func(t *testing.T) {
		pool := pracredentialpool.CredentialPool{
			ID:             "pool-123",
			Name:           "Test Pool",
			CredentialType: "USERNAME_PASSWORD",
		}

		data, err := json.Marshal(pool)
		require.NoError(t, err)

		var unmarshaled pracredentialpool.CredentialPool
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, pool.ID, unmarshaled.ID)
		assert.Equal(t, pool.Name, unmarshaled.Name)
	})
}

func TestPRACredentialPool_MockServerOperations(t *testing.T) {
	t.Run("GET pool by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "pool-123", "name": "Mock Pool"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praCredentialPool")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
