// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/devicegroups"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDeviceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/deviceGroups"

	server.On("GET", path, common.SuccessResponse([]devicegroups.DeviceGroups{
		{ID: 1, Name: "Mobile Devices", Description: "All mobile devices"},
		{ID: 2, Name: "Windows Devices", Description: "All Windows devices"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevicesGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDeviceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/deviceGroups"

	server.On("GET", path, common.SuccessResponse([]devicegroups.DeviceGroups{
		{ID: 1, Name: "Mobile Devices", Description: "All mobile devices"},
		{ID: 2, Name: "Windows Devices", Description: "All Windows devices"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetDeviceGroupByName(context.Background(), service, "Mobile Devices")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Mobile Devices", result.Name)
}

func TestDevices_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/deviceGroups/devices"

	server.On("GET", path, common.SuccessResponse([]devicegroups.Devices{
		{ID: 1, Name: "Device 1", DeviceModel: "iPhone 15"},
		{ID: 2, Name: "Device 2", DeviceModel: "Galaxy S24"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := devicegroups.GetAllDevices(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDeviceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceGroups JSON marshaling", func(t *testing.T) {
		group := devicegroups.DeviceGroups{
			ID:          12345,
			Name:        "Corporate Devices",
			Description: "All corporate managed devices",
			GroupType:   "ZCC_OS",
			OSType:      "WINDOWS_OS",
			Predefined:  false,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Corporate Devices"`)
	})

	t.Run("DeviceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Mobile Devices",
			"groupType": "ZCC_OS",
			"description": "All mobile devices",
			"osType": "IOS",
			"predefined": true,
			"deviceCount": 150
		}`

		var group devicegroups.DeviceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "Mobile Devices", group.Name)
		assert.True(t, group.Predefined)
	})

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := devicegroups.Devices{
			ID:          12345,
			Name:        "LAPTOP-001",
			DeviceModel: "MacBook Pro",
			OSType:      "MAC_OS",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"LAPTOP-001"`)
	})
}
