// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func TestDeviceWebProbes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceWebProbe JSON marshaling", func(t *testing.T) {
		probe := devices.DeviceWebProbe{
			ID:        12345,
			Name:      "office365.com",
			NumProbes: 100,
			AvgScore:  92.5,
			AvgPFT:    250.3,
		}

		data, err := json.Marshal(probe)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"office365.com"`)
		assert.Contains(t, string(data), `"num_probes":100`)
		assert.Contains(t, string(data), `"avg_score":92.5`)
		assert.Contains(t, string(data), `"avg_pft":250.3`)
	})

	t.Run("DeviceWebProbe JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "salesforce.com",
			"num_probes": 50,
			"avg_score": 88.0,
			"avg_pft": 320.5
		}`

		var probe devices.DeviceWebProbe
		err := json.Unmarshal([]byte(jsonData), &probe)
		require.NoError(t, err)

		assert.Equal(t, 67890, probe.ID)
		assert.Equal(t, "salesforce.com", probe.Name)
		assert.Equal(t, 50, probe.NumProbes)
		assert.Equal(t, float32(88.0), probe.AvgScore)
		assert.Equal(t, float32(320.5), probe.AvgPFT)
	})
}

func TestDeviceWebProbes_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse web probes list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "microsoft.com",
				"num_probes": 100,
				"avg_score": 95.0,
				"avg_pft": 200.0
			},
			{
				"id": 2,
				"name": "google.com",
				"num_probes": 80,
				"avg_score": 98.0,
				"avg_pft": 150.0
			},
			{
				"id": 3,
				"name": "slack.com",
				"num_probes": 60,
				"avg_score": 90.0,
				"avg_pft": 280.0
			}
		]`

		var probes []devices.DeviceWebProbe
		err := json.Unmarshal([]byte(jsonResponse), &probes)
		require.NoError(t, err)

		assert.Len(t, probes, 3)
		assert.Equal(t, "microsoft.com", probes[0].Name)
		assert.Equal(t, float32(95.0), probes[0].AvgScore)
		assert.Equal(t, "google.com", probes[1].Name)
		assert.Equal(t, float32(150.0), probes[1].AvgPFT)
	})
}

