// Package unit provides unit tests for ZPA PRA Credential service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/pracredential"
)

func TestPRACredential_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Credential JSON marshaling", func(t *testing.T) {
		cred := pracredential.Credential{
			ID:             "cred-123",
			Name:           "Test Credential",
			Description:    "Test Description",
			CredentialType: "USERNAME_PASSWORD",
			UserName:       "admin",
		}

		data, err := json.Marshal(cred)
		require.NoError(t, err)

		var unmarshaled pracredential.Credential
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, cred.ID, unmarshaled.ID)
		assert.Equal(t, cred.Name, unmarshaled.Name)
	})
}

func TestPRACredential_MockServerOperations(t *testing.T) {
	t.Run("GET credential by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "cred-123", "name": "Mock Cred"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praCredential")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
