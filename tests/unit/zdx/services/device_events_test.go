// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func TestDeviceEvents_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceEvents JSON marshaling", func(t *testing.T) {
		deviceEvents := devices.DeviceEvents{
			TimeStamp: 1699900000,
			Events: []devices.Events{
				{
					Category:    "Network",
					Name:        "network_change",
					DisplayName: "Network Changed",
					Prev:        "WiFi - CorpWiFi",
					Curr:        "Ethernet - LAN",
				},
				{
					Category:    "System",
					Name:        "user_login",
					DisplayName: "User Login",
					Prev:        "",
					Curr:        "john.doe",
				},
			},
		}

		data, err := json.Marshal(deviceEvents)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"timestamp":1699900000`)
		assert.Contains(t, string(data), `"category":"Network"`)
		assert.Contains(t, string(data), `"display_name":"Network Changed"`)
	})

	t.Run("DeviceEvents JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"timestamp": 1699950000,
			"instances": [
				{
					"category": "Hardware",
					"name": "disk_change",
					"display_name": "Disk Changed",
					"prev": "256 GB",
					"curr": "512 GB"
				}
			]
		}`

		var deviceEvents devices.DeviceEvents
		err := json.Unmarshal([]byte(jsonData), &deviceEvents)
		require.NoError(t, err)

		assert.Equal(t, 1699950000, deviceEvents.TimeStamp)
		assert.Len(t, deviceEvents.Events, 1)
		assert.Equal(t, "Hardware", deviceEvents.Events[0].Category)
		assert.Equal(t, "Disk Changed", deviceEvents.Events[0].DisplayName)
	})

	t.Run("Events JSON marshaling", func(t *testing.T) {
		event := devices.Events{
			Category:    "Software",
			Name:        "app_update",
			DisplayName: "Application Updated",
			Prev:        "Chrome 119.0",
			Curr:        "Chrome 120.0",
		}

		data, err := json.Marshal(event)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"category":"Software"`)
		assert.Contains(t, string(data), `"name":"app_update"`)
		assert.Contains(t, string(data), `"prev":"Chrome 119.0"`)
		assert.Contains(t, string(data), `"curr":"Chrome 120.0"`)
	})
}

func TestDeviceEvents_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse device events list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"timestamp": 1699900000,
				"instances": [
					{"category": "Network", "name": "wifi_connect", "display_name": "WiFi Connected", "curr": "CorpWiFi"}
				]
			},
			{
				"timestamp": 1699903600,
				"instances": [
					{"category": "System", "name": "reboot", "display_name": "System Rebooted"},
					{"category": "Software", "name": "update", "display_name": "Software Updated"}
				]
			}
		]`

		var events []devices.DeviceEvents
		err := json.Unmarshal([]byte(jsonResponse), &events)
		require.NoError(t, err)

		assert.Len(t, events, 2)
		assert.Equal(t, 1699900000, events[0].TimeStamp)
		assert.Len(t, events[0].Events, 1)
		assert.Equal(t, 1699903600, events[1].TimeStamp)
		assert.Len(t, events[1].Events, 2)
	})
}

