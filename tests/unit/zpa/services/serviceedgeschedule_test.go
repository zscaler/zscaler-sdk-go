// Package unit provides unit tests for ZPA Service Edge Schedule service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SEAssistantSchedule represents the assistant schedule for testing
type SEAssistantSchedule struct {
	ID                string `json:"id,omitempty"`
	CustomerID        string `json:"customerId"`
	DeleteDisabled    bool   `json:"deleteDisabled"`
	Enabled           bool   `json:"enabled"`
	Frequency         string `json:"frequency"`
	FrequencyInterval string `json:"frequencyInterval"`
}

// TestServiceEdgeSchedule_Structure tests the struct definitions
func TestServiceEdgeSchedule_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AssistantSchedule JSON marshaling", func(t *testing.T) {
		schedule := SEAssistantSchedule{
			ID:                "sched-123",
			CustomerID:        "cust-001",
			DeleteDisabled:    false,
			Enabled:           true,
			Frequency:         "DAYS",
			FrequencyInterval: "7",
		}

		data, err := json.Marshal(schedule)
		require.NoError(t, err)

		var unmarshaled SEAssistantSchedule
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, schedule.ID, unmarshaled.ID)
		assert.True(t, unmarshaled.Enabled)
		assert.Equal(t, "7", unmarshaled.FrequencyInterval)
	})

	t.Run("AssistantSchedule from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "sched-456",
			"customerId": "cust-002",
			"deleteDisabled": true,
			"enabled": true,
			"frequency": "DAYS",
			"frequencyInterval": "14"
		}`

		var schedule SEAssistantSchedule
		err := json.Unmarshal([]byte(apiResponse), &schedule)
		require.NoError(t, err)

		assert.Equal(t, "sched-456", schedule.ID)
		assert.True(t, schedule.Enabled)
		assert.True(t, schedule.DeleteDisabled)
		assert.Equal(t, "14", schedule.FrequencyInterval)
	})
}

// TestServiceEdgeSchedule_MockServerOperations tests operations
func TestServiceEdgeSchedule_MockServerOperations(t *testing.T) {
	t.Run("GET schedule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/serviceEdgeSchedule")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "sched-123", "enabled": true, "frequencyInterval": "7"}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/serviceEdgeSchedule")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create schedule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{"id": "new-sched", "enabled": true}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/serviceEdgeSchedule", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("PUT update schedule", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "PUT", r.Method)
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("PUT", server.URL+"/serviceEdgeSchedule/sched-123", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestServiceEdgeSchedule_SpecialCases tests edge cases
func TestServiceEdgeSchedule_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("Valid frequency intervals", func(t *testing.T) {
		validIntervals := []string{"5", "7", "14", "30", "60", "90"}

		for _, interval := range validIntervals {
			schedule := SEAssistantSchedule{
				ID:                "sched-test",
				Enabled:           true,
				FrequencyInterval: interval,
			}

			data, err := json.Marshal(schedule)
			require.NoError(t, err)

			assert.Contains(t, string(data), interval)
		}
	})

	t.Run("Disabled schedule", func(t *testing.T) {
		schedule := SEAssistantSchedule{
			ID:      "sched-disabled",
			Enabled: false,
		}

		data, err := json.Marshal(schedule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enabled":false`)
	})
}

