// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/devicegroups"
)

func TestDeviceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceGroups JSON marshaling", func(t *testing.T) {
		group := devicegroups.DeviceGroups{
			ID:          12345,
			Name:        "Windows Devices",
			GroupType:   "BYOD",
			Description: "All Windows BYOD devices",
			OSType:      "WINDOWS_OS",
			Predefined:  false,
			DeviceCount: 150,
			DeviceNames: "device1,device2,device3",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"osType":"WINDOWS_OS"`)
		assert.Contains(t, string(data), `"deviceCount":150`)
	})

	t.Run("DeviceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "iOS Devices",
			"groupType": "MANAGED",
			"description": "All managed iOS devices",
			"osType": "IOS",
			"predefined": true,
			"deviceCount": 75
		}`

		var group devicegroups.DeviceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.True(t, group.Predefined)
		assert.Equal(t, "IOS", group.OSType)
	})

	t.Run("Devices JSON marshaling", func(t *testing.T) {
		device := devicegroups.Devices{
			ID:              12345,
			Name:            "johns-macbook",
			DeviceGroupType: "BYOD",
			DeviceModel:     "MacBook Pro",
			OSType:          "MAC_OS",
			OSVersion:       "14.2.1",
			Description:     "John's work laptop",
			OwnerUserId:     100,
			OwnerName:       "john.doe@company.com",
			HostName:        "johns-macbook.local",
		}

		data, err := json.Marshal(device)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"deviceModel":"MacBook Pro"`)
		assert.Contains(t, string(data), `"ownerUserId":100`)
	})

	t.Run("Devices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "android-tablet",
			"deviceGroupType": "MANAGED",
			"deviceModel": "Galaxy Tab S9",
			"osType": "ANDROID",
			"osVersion": "14.0",
			"description": "Company tablet",
			"ownerUserId": 200,
			"ownerName": "jane.smith@company.com",
			"hostName": "android-tablet.local"
		}`

		var device devicegroups.Devices
		err := json.Unmarshal([]byte(jsonData), &device)
		require.NoError(t, err)

		assert.Equal(t, 67890, device.ID)
		assert.Equal(t, "ANDROID", device.OSType)
		assert.Equal(t, "Galaxy Tab S9", device.DeviceModel)
	})
}

func TestDeviceGroups_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse device groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Windows", "osType": "WINDOWS_OS", "deviceCount": 100},
			{"id": 2, "name": "macOS", "osType": "MAC_OS", "deviceCount": 50},
			{"id": 3, "name": "iOS", "osType": "IOS", "deviceCount": 75}
		]`

		var groups []devicegroups.DeviceGroups
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, 75, groups[2].DeviceCount)
	})

	t.Run("Parse devices list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "device1", "osType": "WINDOWS_OS"},
			{"id": 2, "name": "device2", "osType": "MAC_OS"},
			{"id": 3, "name": "device3", "osType": "ANDROID"}
		]`

		var devices []devicegroups.Devices
		err := json.Unmarshal([]byte(jsonResponse), &devices)
		require.NoError(t, err)

		assert.Len(t, devices, 3)
	})
}

