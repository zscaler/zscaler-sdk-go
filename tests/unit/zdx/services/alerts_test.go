// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/alerts"
)

func TestAlerts_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Alert JSON marshaling", func(t *testing.T) {
		alert := alerts.Alert{
			ID:              12345,
			RuleName:        "High CPU Usage Alert",
			Severity:        "critical",
			AlertType:       "device_health",
			AlertStatus:     "active",
			NumGeolocations: 3,
			NumDevices:      150,
			StartedOn:       1699900000,
			EndedOn:         0,
			Application: alerts.Application{
				ID:   1,
				Name: "Microsoft 365",
			},
			Departments: []alerts.Department{
				{ID: 1, Name: "Engineering", NumDevices: 75},
				{ID: 2, Name: "Sales", NumDevices: 50},
			},
			Locations: []alerts.Location{
				{
					ID:         1,
					Name:       "San Jose HQ",
					NumDevices: 100,
					Groups: []alerts.Group{
						{ID: 1, Name: "Building A"},
					},
				},
			},
			Geolocations: []alerts.Geolocation{
				{ID: "US-CA", Name: "California, US", NumDevices: 100},
			},
		}

		data, err := json.Marshal(alert)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"rule_name":"High CPU Usage Alert"`)
		assert.Contains(t, string(data), `"severity":"critical"`)
		assert.Contains(t, string(data), `"alert_status":"active"`)
	})

	t.Run("Alert JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"rule_name": "Network Latency Alert",
			"severity": "warning",
			"alert_type": "network",
			"alert_status": "resolved",
			"num_geolocations": 2,
			"num_devices": 50,
			"started_on": 1699800000,
			"ended_on": 1699850000,
			"application": {
				"id": 2,
				"name": "Salesforce"
			},
			"departments": [
				{"id": 1, "name": "Support", "num_devices": 25}
			],
			"locations": [],
			"geolocations": []
		}`

		var alert alerts.Alert
		err := json.Unmarshal([]byte(jsonData), &alert)
		require.NoError(t, err)

		assert.Equal(t, 67890, alert.ID)
		assert.Equal(t, "Network Latency Alert", alert.RuleName)
		assert.Equal(t, "warning", alert.Severity)
		assert.Equal(t, "resolved", alert.AlertStatus)
		assert.Equal(t, 50, alert.NumDevices)
		assert.Equal(t, "Salesforce", alert.Application.Name)
		assert.Len(t, alert.Departments, 1)
	})

	t.Run("AlertsResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"alerts": [
				{"id": 1, "rule_name": "Alert 1", "severity": "critical"},
				{"id": 2, "rule_name": "Alert 2", "severity": "warning"}
			],
			"next_offset": "abc123"
		}`

		var response alerts.AlertsResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Alerts, 2)
		assert.Equal(t, "abc123", response.NextOffset)
		assert.Equal(t, "Alert 1", response.Alerts[0].RuleName)
	})

	t.Run("Device JSON marshaling", func(t *testing.T) {
		device := alerts.Device{
			ID:        12345,
			Name:      "LAPTOP-001",
			UserID:    1001,
			UserName:  "john.doe",
			UserEmail: "john.doe@example.com",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"LAPTOP-001"`)
		assert.Contains(t, string(data), `"userEmail":"john.doe@example.com"`)
	})

	t.Run("AffectedDevicesResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"devices": [
				{"id": 1, "name": "Device 1", "userid": 100, "userName": "user1", "userEmail": "user1@example.com"},
				{"id": 2, "name": "Device 2", "userid": 101, "userName": "user2", "userEmail": "user2@example.com"}
			],
			"next_offset": "offset123"
		}`

		var response alerts.AffectedDevicesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Devices, 2)
		assert.Equal(t, "offset123", response.NextOffset)
		assert.Equal(t, "Device 1", response.Devices[0].Name)
		assert.Equal(t, "user1@example.com", response.Devices[0].UserEmail)
	})
}

func TestAlerts_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse ongoing alerts response", func(t *testing.T) {
		jsonResponse := `{
			"alerts": [
				{
					"id": 1001,
					"rule_name": "CPU Usage > 90%",
					"severity": "critical",
					"alert_type": "device_health",
					"alert_status": "active",
					"num_devices": 25,
					"started_on": 1699900000,
					"application": {"id": 1, "name": "System"},
					"departments": [
						{"id": 1, "name": "Engineering", "num_devices": 15},
						{"id": 2, "name": "Design", "num_devices": 10}
					]
				},
				{
					"id": 1002,
					"rule_name": "Network Packet Loss",
					"severity": "warning",
					"alert_type": "network",
					"alert_status": "active",
					"num_devices": 10,
					"started_on": 1699890000
				}
			],
			"next_offset": ""
		}`

		var response alerts.AlertsResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Alerts, 2)
		assert.Empty(t, response.NextOffset)
		
		// Check first alert
		assert.Equal(t, 1001, response.Alerts[0].ID)
		assert.Equal(t, "critical", response.Alerts[0].Severity)
		assert.Equal(t, 25, response.Alerts[0].NumDevices)
		assert.Len(t, response.Alerts[0].Departments, 2)
		
		// Check second alert
		assert.Equal(t, 1002, response.Alerts[1].ID)
		assert.Equal(t, "warning", response.Alerts[1].Severity)
	})

	t.Run("Parse historical alerts response", func(t *testing.T) {
		jsonResponse := `{
			"alerts": [
				{
					"id": 2001,
					"rule_name": "Disk Space Low",
					"severity": "warning",
					"alert_status": "resolved",
					"started_on": 1699700000,
					"ended_on": 1699750000,
					"num_devices": 5
				}
			],
			"next_offset": "page2"
		}`

		var response alerts.AlertsResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.Alerts, 1)
		assert.Equal(t, "page2", response.NextOffset)
		assert.Equal(t, "resolved", response.Alerts[0].AlertStatus)
		assert.NotZero(t, response.Alerts[0].EndedOn)
	})
}

