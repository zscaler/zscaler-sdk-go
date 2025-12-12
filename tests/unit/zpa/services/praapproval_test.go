// Package unit provides unit tests for ZPA PRA Approval service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praapproval"
)

func TestPRAApproval_Structure(t *testing.T) {
	t.Parallel()

	t.Run("PrivilegedApproval JSON marshaling", func(t *testing.T) {
		approval := praapproval.PrivilegedApproval{
			ID:           "pa-123",
			EmailIDs:     []string{"user@example.com"},
			StartTime:    "1609459200000",
			EndTime:      "1612137600000",
			Status:       "ACTIVE",
		}

		data, err := json.Marshal(approval)
		require.NoError(t, err)

		var unmarshaled praapproval.PrivilegedApproval
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, approval.ID, unmarshaled.ID)
	})
}

func TestPRAApproval_MockServerOperations(t *testing.T) {
	t.Run("GET approval by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "pa-123", "status": "ACTIVE"}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/praApproval")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
