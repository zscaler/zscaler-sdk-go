// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDevices_GetDeviceCleanupInfo_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getDeviceCleanupInfo"

	server.On("GET", path, common.SuccessResponse(devices.DeviceCleanupInfo{
		ID:                    "cleanup-001",
		Active:                "true",
		AutoPurgeDays:         "30",
		AutoRemovalDays:       "60",
		CompanyID:             "company-123",
		DeviceExceedLimit:     "100",
		ForceRemoveType:       "1",
		ForceRemoveTypeString: "IMMEDIATE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetDeviceCleanupInfo(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "cleanup-001", result.ID)
	assert.Equal(t, "true", result.Active)
	assert.Equal(t, "30", result.AutoPurgeDays)
}

func TestDevices_SetDeviceCleanupInfo_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/setDeviceCleanupInfo"

	server.On("PUT", path, common.SuccessResponse(devices.DeviceCleanupInfo{
		ID:              "cleanup-001",
		Active:          "true",
		AutoPurgeDays:   "45",
		AutoRemovalDays: "90",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	cleanupInfo := &devices.DeviceCleanupInfo{
		ID:              "cleanup-001",
		Active:          "true",
		AutoPurgeDays:   "45",
		AutoRemovalDays: "90",
	}

	result, err := devices.SetDeviceCleanupInfo(context.Background(), service, cleanupInfo)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "45", result.AutoPurgeDays)
}

func TestDevices_GetDeviceDetails_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getDeviceDetails"

	server.On("GET", path, common.SuccessResponse([]devices.DeviceDetails{
		{
			AgentVersion:    "4.2.0.100",
			MachineHostname: "laptop-001",
			MacAddress:      "00:11:22:33:44:55",
			OSVersion:       "Windows 10",
			State:           "ACTIVE",
			Type:            "WINDOWS",
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devices.GetDeviceDetails(context.Background(), service, "user@example.com", "device-001")

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "laptop-001", result[0].MachineHostname)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDevices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("GetDevices JSON marshaling", func(t *testing.T) {
		device := devices.GetDevices{
			AgentVersion:      "4.2.0.100",
			CompanyName:       "Test Company",
			MachineHostname:   "test-machine",
			MacAddress:        "00:11:22:33:44:55",
			OsVersion:         "Windows 10",
			Owner:             "test.user@example.com",
			User:              "test.user@example.com",
			RegistrationState: "REGISTERED",
			State:             1,
			Type:              1,
			Udid:              "device-udid-123",
			VpnState:          1,
			ZappArch:          "x64",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"agentVersion":"4.2.0.100"`)
		assert.Contains(t, string(data), `"machineHostname":"test-machine"`)
		assert.Contains(t, string(data), `"macAddress":"00:11:22:33:44:55"`)
	})

	t.Run("GetDevices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"agentVersion": "4.2.0.100",
			"companyName": "Test Company",
			"machineHostname": "test-machine",
			"macAddress": "00:11:22:33:44:55",
			"osVersion": "Windows 10",
			"owner": "test.user@example.com",
			"registrationState": "REGISTERED",
			"state": 1,
			"type": 1,
			"udid": "device-udid-123",
			"user": "test.user@example.com",
			"vpnState": 1,
			"zappArch": "x64"
		}`

		var device devices.GetDevices
		err := json.Unmarshal([]byte(jsonData), &device)
		require.NoError(t, err)

		assert.Equal(t, "4.2.0.100", device.AgentVersion)
		assert.Equal(t, "Test Company", device.CompanyName)
		assert.Equal(t, "test-machine", device.MachineHostname)
		assert.Equal(t, "00:11:22:33:44:55", device.MacAddress)
		assert.Equal(t, 1, device.State)
		assert.Equal(t, 1, device.VpnState)
	})

	t.Run("DeviceDetails JSON marshaling", func(t *testing.T) {
		details := devices.DeviceDetails{
			AgentVersion:        "4.2.0.100",
			MachineHostname:     "test-machine",
			MacAddress:          "00:11:22:33:44:55",
			OSVersion:           "Windows 10",
			Owner:               "test.user@example.com",
			UserName:            "test.user@example.com",
			State:               "ACTIVE",
			Type:                "WINDOWS",
			UniqueID:            "device-unique-123",
			HardwareFingerprint: "hw-fingerprint-abc",
		}

		data, err := json.Marshal(details)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"agent_version":"4.2.0.100"`)
		assert.Contains(t, string(data), `"machineHostname":"test-machine"`)
		assert.Contains(t, string(data), `"mac_address":"00:11:22:33:44:55"`)
	})

	t.Run("DeviceCleanupInfo JSON marshaling", func(t *testing.T) {
		cleanup := devices.DeviceCleanupInfo{
			ID:                    "cleanup-123",
			Active:                "true",
			AutoPurgeDays:         "30",
			AutoRemovalDays:       "60",
			CompanyID:             "company-456",
			DeviceExceedLimit:     "100",
			ForceRemoveType:       "1",
			ForceRemoveTypeString: "IMMEDIATE",
		}

		data, err := json.Marshal(cleanup)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"cleanup-123"`)
		assert.Contains(t, string(data), `"active":"true"`)
		assert.Contains(t, string(data), `"autoPurgeDays":"30"`)
	})
}

func TestDevices_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse devices list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"agentVersion": "4.2.0.100",
				"machineHostname": "device-1",
				"state": 1,
				"udid": "udid-001"
			},
			{
				"agentVersion": "4.2.0.101",
				"machineHostname": "device-2",
				"state": 0,
				"udid": "udid-002"
			}
		]`

		var devicesList []devices.GetDevices
		err := json.Unmarshal([]byte(jsonResponse), &devicesList)
		require.NoError(t, err)

		assert.Len(t, devicesList, 2)
		assert.Equal(t, "device-1", devicesList[0].MachineHostname)
		assert.Equal(t, "device-2", devicesList[1].MachineHostname)
		assert.Equal(t, 1, devicesList[0].State)
		assert.Equal(t, 0, devicesList[1].State)
	})

	t.Run("Parse device details response", func(t *testing.T) {
		jsonResponse := `[
			{
				"agent_version": "4.2.0.100",
				"machineHostname": "test-machine",
				"mac_address": "00:11:22:33:44:55",
				"os_version": "Windows 10",
				"state": "ACTIVE",
				"type": "WINDOWS",
				"unique_id": "device-unique-123"
			}
		]`

		var detailsList []devices.DeviceDetails
		err := json.Unmarshal([]byte(jsonResponse), &detailsList)
		require.NoError(t, err)

		assert.Len(t, detailsList, 1)
		assert.Equal(t, "test-machine", detailsList[0].MachineHostname)
		assert.Equal(t, "ACTIVE", detailsList[0].State)
	})
}
