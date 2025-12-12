// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/activation"
)

func TestActivation_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Activation JSON marshaling", func(t *testing.T) {
		act := activation.Activation{
			Status: "ACTIVE",
		}

		data, err := json.Marshal(act)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"status":"ACTIVE"`)
	})

	t.Run("Activation JSON unmarshaling", func(t *testing.T) {
		jsonData := `{"status":"PENDING"}`

		var act activation.Activation
		err := json.Unmarshal([]byte(jsonData), &act)
		require.NoError(t, err)

		assert.Equal(t, "PENDING", act.Status)
	})

	t.Run("ZiaEusaStatus JSON marshaling", func(t *testing.T) {
		eusa := activation.ZiaEusaStatus{
			ID:             12345,
			AcceptedStatus: true,
		}

		data, err := json.Marshal(eusa)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"acceptedStatus":true`)
	})

	t.Run("ZiaEusaStatus JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"version": {
				"id": 1,
				"name": "v1.0"
			},
			"acceptedStatus": false
		}`

		var eusa activation.ZiaEusaStatus
		err := json.Unmarshal([]byte(jsonData), &eusa)
		require.NoError(t, err)

		assert.Equal(t, 54321, eusa.ID)
		assert.False(t, eusa.AcceptedStatus)
		assert.NotNil(t, eusa.Version)
	})
}

func TestActivation_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse activation status response", func(t *testing.T) {
		jsonResponse := `{"status":"ACTIVE"}`

		var act activation.Activation
		err := json.Unmarshal([]byte(jsonResponse), &act)
		require.NoError(t, err)

		assert.Equal(t, "ACTIVE", act.Status)
	})

	t.Run("Parse various activation statuses", func(t *testing.T) {
		statuses := []string{"ACTIVE", "PENDING", "INPROGRESS", "NONE"}

		for _, status := range statuses {
			jsonData := `{"status":"` + status + `"}`
			var act activation.Activation
			err := json.Unmarshal([]byte(jsonData), &act)
			require.NoError(t, err)
			assert.Equal(t, status, act.Status)
		}
	})
}

