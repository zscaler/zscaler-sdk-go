// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/entitlements"
)

func TestEntitlements_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZdxGroupEntitlements JSON marshaling", func(t *testing.T) {
		zdxGroup := entitlements.ZdxGroupEntitlements{
			CollectZdxLocation:        1,
			ComputeDeviceGroupsForZDX: 1,
			LogoutZCCForZDXService:    0,
			TotalCount:                2,
			UpmEnableForAll:           1,
			UpmDeviceGroupList: []entitlements.DeviceGroup{
				{
					GroupID:    1,
					GroupName:  "All Devices",
					Active:     1,
					AuthType:   "SAML",
					UpmEnabled: 1,
				},
			},
			UpmGroupList: []entitlements.DeviceGroup{
				{
					GroupID:    2,
					GroupName:  "Engineering",
					Active:     1,
					AuthType:   "SAML",
					UpmEnabled: 1,
				},
			},
		}

		data, err := json.Marshal(zdxGroup)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"collectZdxLocation":1`)
		assert.Contains(t, string(data), `"computeDeviceGroupsForZDX":1`)
		assert.Contains(t, string(data), `"upmEnableForAll":1`)
	})

	t.Run("ZdxGroupEntitlements JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"collectZdxLocation": 1,
			"computeDeviceGroupsForZDX": 1,
			"logoutZCCForZDXService": 0,
			"totalCount": 3,
			"upmEnableForAll": 0,
			"upmDeviceGroupList": [
				{
					"groupId": 1,
					"groupName": "Windows Devices",
					"active": 1,
					"authType": "SAML",
					"upmEnabled": 1
				}
			],
			"upmGroupList": []
		}`

		var zdxGroup entitlements.ZdxGroupEntitlements
		err := json.Unmarshal([]byte(jsonData), &zdxGroup)
		require.NoError(t, err)

		assert.Equal(t, 1, zdxGroup.CollectZdxLocation)
		assert.Equal(t, 3, zdxGroup.TotalCount)
		assert.Len(t, zdxGroup.UpmDeviceGroupList, 1)
		assert.Equal(t, "Windows Devices", zdxGroup.UpmDeviceGroupList[0].GroupName)
	})

	t.Run("ZpaGroupEntitlements JSON marshaling", func(t *testing.T) {
		zpaGroup := entitlements.ZpaGroupEntitlements{
			ComputeDeviceGroupsForZPA: 1,
			MachineTunEnabledForAll:   0,
			TotalCount:                2,
			ZpaEnableForAll:           1,
			DeviceGroupList: []entitlements.DeviceGroupItem{
				{
					GroupID:    1,
					GroupName:  "All Devices",
					Active:     1,
					AuthType:   "SAML",
					ZpaEnabled: 1,
				},
			},
			GroupList: []entitlements.GroupListItem{
				{
					GroupID:    2,
					GroupName:  "Engineering",
					Active:     1,
					AuthType:   "SAML",
					ZpaEnabled: 1,
				},
			},
		}

		data, err := json.Marshal(zpaGroup)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"computeDeviceGroupsForZPA":1`)
		assert.Contains(t, string(data), `"zpaEnableForAll":1`)
	})

	t.Run("DeviceGroup JSON marshaling", func(t *testing.T) {
		group := entitlements.DeviceGroup{
			GroupID:    123,
			GroupName:  "Test Group",
			Active:     1,
			AuthType:   "OIDC",
			UpmEnabled: 1,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"groupId":123`)
		assert.Contains(t, string(data), `"groupName":"Test Group"`)
		assert.Contains(t, string(data), `"authType":"OIDC"`)
	})
}

func TestEntitlements_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse ZDX group entitlements response", func(t *testing.T) {
		jsonResponse := `{
			"collectZdxLocation": 1,
			"computeDeviceGroupsForZDX": 1,
			"logoutZCCForZDXService": 0,
			"totalCount": 5,
			"upmEnableForAll": 1,
			"upmDeviceGroupList": [
				{"groupId": 1, "groupName": "Group A", "active": 1, "upmEnabled": 1},
				{"groupId": 2, "groupName": "Group B", "active": 1, "upmEnabled": 0}
			],
			"upmGroupList": [
				{"groupId": 3, "groupName": "Group C", "active": 1, "upmEnabled": 1}
			]
		}`

		var zdxGroup entitlements.ZdxGroupEntitlements
		err := json.Unmarshal([]byte(jsonResponse), &zdxGroup)
		require.NoError(t, err)

		assert.Equal(t, 5, zdxGroup.TotalCount)
		assert.Len(t, zdxGroup.UpmDeviceGroupList, 2)
		assert.Len(t, zdxGroup.UpmGroupList, 1)
		assert.Equal(t, "Group A", zdxGroup.UpmDeviceGroupList[0].GroupName)
		assert.Equal(t, 1, zdxGroup.UpmDeviceGroupList[0].UpmEnabled)
		assert.Equal(t, 0, zdxGroup.UpmDeviceGroupList[1].UpmEnabled)
	})

	t.Run("Parse ZPA group entitlements response", func(t *testing.T) {
		jsonResponse := `{
			"computeDeviceGroupsForZPA": 1,
			"machineTunEnabledForAll": 1,
			"totalCount": 4,
			"zpaEnableForAll": 0,
			"deviceGroupList": [
				{"groupId": 1, "groupName": "Device Group 1", "active": 1, "zpaEnabled": 1}
			],
			"groupList": [
				{"groupId": 2, "groupName": "User Group 1", "active": 1, "zpaEnabled": 0},
				{"groupId": 3, "groupName": "User Group 2", "active": 0, "zpaEnabled": 1}
			]
		}`

		var zpaGroup entitlements.ZpaGroupEntitlements
		err := json.Unmarshal([]byte(jsonResponse), &zpaGroup)
		require.NoError(t, err)

		assert.Equal(t, 4, zpaGroup.TotalCount)
		assert.Equal(t, 1, zpaGroup.MachineTunEnabledForAll)
		assert.Len(t, zpaGroup.DeviceGroupList, 1)
		assert.Len(t, zpaGroup.GroupList, 2)
	})
}

