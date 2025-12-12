// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func TestDeviceApps_Structure(t *testing.T) {
	t.Parallel()

	t.Run("App JSON marshaling", func(t *testing.T) {
		app := devices.App{
			ID:    12345,
			Name:  "Microsoft 365",
			Score: 92.5,
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Microsoft 365"`)
		assert.Contains(t, string(data), `"score":92.5`)
	})

	t.Run("App JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "Salesforce",
			"score": 88.3
		}`

		var app devices.App
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, 67890, app.ID)
		assert.Equal(t, "Salesforce", app.Name)
		assert.Equal(t, float32(88.3), app.Score)
	})
}

func TestDeviceApps_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse device apps list response", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Microsoft 365", "score": 95.0},
			{"id": 2, "name": "Salesforce", "score": 88.5},
			{"id": 3, "name": "Slack", "score": 92.0},
			{"id": 4, "name": "Zoom", "score": 78.5}
		]`

		var apps []devices.App
		err := json.Unmarshal([]byte(jsonResponse), &apps)
		require.NoError(t, err)

		assert.Len(t, apps, 4)
		assert.Equal(t, "Microsoft 365", apps[0].Name)
		assert.Equal(t, float32(95.0), apps[0].Score)
		assert.Equal(t, "Zoom", apps[3].Name)
		assert.Equal(t, float32(78.5), apps[3].Score)
	})

	t.Run("Parse single device app response", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "SAP",
			"score": 72.5
		}`

		var app devices.App
		err := json.Unmarshal([]byte(jsonResponse), &app)
		require.NoError(t, err)

		assert.Equal(t, 100, app.ID)
		assert.Equal(t, "SAP", app.Name)
		assert.Equal(t, float32(72.5), app.Score)
	})
}

