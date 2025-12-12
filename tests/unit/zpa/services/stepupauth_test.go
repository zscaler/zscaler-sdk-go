// Package unit provides unit tests for ZPA Step Up Auth service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/step_up_auth"
)

func TestStepUpAuth_Structure(t *testing.T) {
	t.Parallel()

	t.Run("StepAuthLevel JSON marshaling", func(t *testing.T) {
		auth := step_up_auth.StepAuthLevel{
			ID:          "sua-123",
			Description: "High Security Level",
		}

		data, err := json.Marshal(auth)
		require.NoError(t, err)

		var unmarshaled step_up_auth.StepAuthLevel
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, auth.ID, unmarshaled.ID)
		assert.Equal(t, auth.Description, unmarshaled.Description)
	})
}

func TestStepUpAuth_MockServerOperations(t *testing.T) {
	t.Run("GET step up auth level", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "sua-123", "description": "Level 1"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/stepUpAuth")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
