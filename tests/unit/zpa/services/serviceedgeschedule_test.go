// Package unit provides unit tests for ZPA Service Edge Schedule service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgeschedule"
)

func TestServiceEdgeSchedule_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AssistantSchedule JSON marshaling", func(t *testing.T) {
		schedule := serviceedgeschedule.AssistantSchedule{
			ID:                "sched-123",
			CustomerID:        "cust-001",
			Enabled:           true,
			Frequency:         "days",
			FrequencyInterval: "7",
			DeleteDisabled:    false,
		}

		data, err := json.Marshal(schedule)
		require.NoError(t, err)

		var unmarshaled serviceedgeschedule.AssistantSchedule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, schedule.ID, unmarshaled.ID)
		assert.True(t, unmarshaled.Enabled)
	})
}

func TestServiceEdgeSchedule_MockServerOperations(t *testing.T) {
	t.Run("GET schedule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"id": "sched-123", "enabled": true}`))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/schedule")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
