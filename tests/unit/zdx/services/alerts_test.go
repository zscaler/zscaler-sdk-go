// Package services provides unit tests for ZDX alerts service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/alerts"
	zdxcommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestAlerts_GetOngoingAlerts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/alerts/ongoing"

	server.On("GET", path, common.SuccessResponse(alerts.AlertsResponse{
		Alerts: []alerts.Alert{
			{ID: 1, RuleName: "High Latency Alert", Severity: "critical", AlertStatus: "active"},
			{ID: 2, RuleName: "Packet Loss Alert", Severity: "warning", AlertStatus: "active"},
		},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := alerts.GetOngoingAlerts(context.Background(), service, zdxcommon.GetFromToFilters{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Alerts, 2)
	assert.Equal(t, "High Latency Alert", result.Alerts[0].RuleName)
	assert.Equal(t, "critical", result.Alerts[0].Severity)
}

func TestAlerts_GetHistoricalAlerts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/alerts/historical"

	server.On("GET", path, common.SuccessResponse(alerts.AlertsResponse{
		Alerts: []alerts.Alert{
			{ID: 10, RuleName: "Past Latency Alert", Severity: "warning", AlertStatus: "resolved", StartedOn: 1700000000, EndedOn: 1700100000},
			{ID: 11, RuleName: "Past DNS Alert", Severity: "critical", AlertStatus: "resolved", StartedOn: 1700200000, EndedOn: 1700300000},
		},
		NextOffset: "page2",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	filters := zdxcommon.GetFromToFilters{
		From: 1699900000,
		To:   1700400000,
	}

	result, _, err := alerts.GetHistoricalAlerts(context.Background(), service, filters)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Alerts, 2)
	assert.Equal(t, "Past Latency Alert", result.Alerts[0].RuleName)
	assert.Equal(t, "page2", result.NextOffset)
}

func TestAlerts_GetAlert_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/alerts/12345"

	server.On("GET", path, common.SuccessResponse(alerts.Alert{
		ID:          12345,
		RuleName:    "Critical Application Alert",
		Severity:    "critical",
		AlertType:   "application",
		AlertStatus: "active",
		NumDevices:  150,
		StartedOn:   1700000000,
		Application: alerts.Application{
			ID:   100,
			Name: "Office 365",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := alerts.GetAlert(context.Background(), service, "12345")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 12345, result.ID)
	assert.Equal(t, "Critical Application Alert", result.RuleName)
	assert.Equal(t, "Office 365", result.Application.Name)
}

func TestAlerts_GetAffectedDevices_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/alerts/12345/affected_devices"

	server.On("GET", path, common.SuccessResponse(alerts.AffectedDevicesResponse{
		Devices: []alerts.Device{
			{ID: 1, Name: "LAPTOP-001", UserName: "john.doe", UserEmail: "john.doe@example.com"},
			{ID: 2, Name: "LAPTOP-002", UserName: "jane.smith", UserEmail: "jane.smith@example.com"},
			{ID: 3, Name: "DESKTOP-001", UserName: "bob.wilson", UserEmail: "bob.wilson@example.com"},
		},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := alerts.GetAffectedDevices(context.Background(), service, "12345", zdxcommon.GetFromToFilters{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Devices, 3)
	assert.Equal(t, "LAPTOP-001", result.Devices[0].Name)
	assert.Equal(t, "john.doe@example.com", result.Devices[0].UserEmail)
}

func TestAlerts_GetOngoingAlerts_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/alerts/ongoing"

	server.On("GET", path, common.SuccessResponse(alerts.AlertsResponse{
		Alerts:     []alerts.Alert{},
		NextOffset: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := alerts.GetOngoingAlerts(context.Background(), service, zdxcommon.GetFromToFilters{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result.Alerts, 0)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestAlerts_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Alert JSON marshaling", func(t *testing.T) {
		alert := alerts.Alert{
			ID:              123,
			RuleName:        "Test Alert",
			Severity:        "critical",
			AlertType:       "network",
			AlertStatus:     "active",
			NumGeolocations: 5,
			NumDevices:      100,
			StartedOn:       1700000000,
			Application: alerts.Application{
				ID:   1,
				Name: "Test App",
			},
		}

		data, err := json.Marshal(alert)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":123`)
		assert.Contains(t, string(data), `"rule_name":"Test Alert"`)
		assert.Contains(t, string(data), `"severity":"critical"`)
		assert.Contains(t, string(data), `"num_devices":100`)
	})

	t.Run("Alert JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"rule_name": "Latency Alert",
			"severity": "warning",
			"alert_type": "application",
			"alert_status": "resolved",
			"num_geolocations": 3,
			"num_devices": 50,
			"started_on": 1700000000,
			"ended_on": 1700100000,
			"application": {
				"id": 10,
				"name": "Salesforce"
			}
		}`

		var alert alerts.Alert
		err := json.Unmarshal([]byte(jsonData), &alert)
		require.NoError(t, err)

		assert.Equal(t, 456, alert.ID)
		assert.Equal(t, "Latency Alert", alert.RuleName)
		assert.Equal(t, "Salesforce", alert.Application.Name)
		assert.Equal(t, 50, alert.NumDevices)
	})

	t.Run("AlertsResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"alerts": [
				{"id": 1, "rule_name": "Alert 1", "severity": "critical"},
				{"id": 2, "rule_name": "Alert 2", "severity": "warning"}
			],
			"next_offset": "page2"
		}`

		var response alerts.AlertsResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Alerts, 2)
		assert.Equal(t, "page2", response.NextOffset)
	})

	t.Run("AffectedDevicesResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"devices": [
				{"id": 1, "name": "Device1", "userName": "user1", "userEmail": "user1@test.com"},
				{"id": 2, "name": "Device2", "userName": "user2", "userEmail": "user2@test.com"}
			],
			"next_offset": ""
		}`

		var response alerts.AffectedDevicesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Len(t, response.Devices, 2)
		assert.Equal(t, "Device1", response.Devices[0].Name)
		assert.Equal(t, "user1@test.com", response.Devices[0].UserEmail)
	})

	t.Run("Device JSON marshaling", func(t *testing.T) {
		device := alerts.Device{
			ID:        100,
			Name:      "LAPTOP-XYZ",
			UserID:    500,
			UserName:  "john.doe",
			UserEmail: "john.doe@company.com",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":100`)
		assert.Contains(t, string(data), `"name":"LAPTOP-XYZ"`)
		assert.Contains(t, string(data), `"userEmail":"john.doe@company.com"`)
	})

	t.Run("Department JSON marshaling", func(t *testing.T) {
		dept := alerts.Department{
			ID:         1,
			Name:       "Engineering",
			NumDevices: 50,
		}

		data, err := json.Marshal(dept)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":1`)
		assert.Contains(t, string(data), `"name":"Engineering"`)
		assert.Contains(t, string(data), `"num_devices":50`)
	})

	t.Run("Location with Groups JSON marshaling", func(t *testing.T) {
		location := alerts.Location{
			ID:         100,
			Name:       "San Jose",
			NumDevices: 200,
			Groups: []alerts.Group{
				{ID: 1, Name: "Group A"},
				{ID: 2, Name: "Group B"},
			},
		}

		data, err := json.Marshal(location)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":100`)
		assert.Contains(t, string(data), `"name":"San Jose"`)
		assert.Contains(t, string(data), `"groups"`)
	})

	t.Run("Geolocation JSON marshaling", func(t *testing.T) {
		geo := alerts.Geolocation{
			ID:         "US-CA",
			Name:       "California, US",
			NumDevices: 150,
		}

		data, err := json.Marshal(geo)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"US-CA"`)
		assert.Contains(t, string(data), `"name":"California, US"`)
		assert.Contains(t, string(data), `"num_devices":150`)
	})
}

func TestAlerts_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse complex alert with nested structures", func(t *testing.T) {
		jsonResponse := `{
			"id": 999,
			"rule_name": "Complex Alert",
			"severity": "critical",
			"alert_type": "application",
			"alert_status": "active",
			"num_geolocations": 2,
			"num_devices": 75,
			"started_on": 1700000000,
			"application": {
				"id": 50,
				"name": "Microsoft Teams"
			},
			"departments": [
				{"id": 1, "name": "Engineering", "num_devices": 40},
				{"id": 2, "name": "Sales", "num_devices": 35}
			],
			"locations": [
				{"id": 100, "name": "San Jose", "num_devices": 50, "groups": [{"id": 1, "name": "Building A"}]},
				{"id": 101, "name": "New York", "num_devices": 25, "groups": []}
			],
			"geolocations": [
				{"id": "US-CA", "name": "California, US", "num_devices": 50},
				{"id": "US-NY", "name": "New York, US", "num_devices": 25}
			]
		}`

		var alert alerts.Alert
		err := json.Unmarshal([]byte(jsonResponse), &alert)
		require.NoError(t, err)

		assert.Equal(t, 999, alert.ID)
		assert.Equal(t, "Complex Alert", alert.RuleName)
		assert.Equal(t, "Microsoft Teams", alert.Application.Name)
		assert.Len(t, alert.Departments, 2)
		assert.Len(t, alert.Locations, 2)
		assert.Len(t, alert.Geolocations, 2)
		assert.Equal(t, "Building A", alert.Locations[0].Groups[0].Name)
	})
}
