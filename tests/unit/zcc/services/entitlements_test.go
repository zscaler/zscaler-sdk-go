// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/entitlements"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestEntitlements_GetZdxGroup_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getZdxGroupEntitlements"

	server.On("GET", path, common.SuccessResponse([]entitlements.ZdxGroupEntitlements{
		{
			TotalCount:                2,
			UpmEnableForAll:           1,
			CollectZdxLocation:        1,
			ComputeDeviceGroupsForZDX: 1,
			UpmGroupList: []entitlements.DeviceGroup{
				{GroupID: 1, GroupName: "Group 1", Active: 1},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := entitlements.GetZdxGroupEntitlements(context.Background(), service, "", 100)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 2, result[0].TotalCount)
}

func TestEntitlements_GetZpaGroup_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/getZpaGroupEntitlements"

	server.On("GET", path, common.SuccessResponse([]entitlements.ZpaGroupEntitlements{
		{
			TotalCount:                3,
			ZpaEnableForAll:           1,
			ComputeDeviceGroupsForZPA: 1,
			GroupList: []entitlements.GroupListItem{
				{GroupID: 1, GroupName: "ZPA Group 1", Active: 1, ZpaEnabled: 1},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := entitlements.GetZpaGroupEntitlements(context.Background(), service, "", 100)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, 3, result[0].TotalCount)
}

func TestEntitlements_UpdateZdxGroup_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/updateZdxGroupEntitlement"

	server.On("PUT", path, common.SuccessResponse(entitlements.ZdxGroupEntitlements{
		TotalCount:      5,
		UpmEnableForAll: 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGroup := &entitlements.ZdxGroupEntitlements{
		UpmEnableForAll: 1,
	}

	result, err := entitlements.UpdateZdxGroupEntitlements(context.Background(), service, updateGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.UpmEnableForAll)
}

func TestEntitlements_UpdateZpaGroup_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/updateZpaGroupEntitlement"

	server.On("PUT", path, common.SuccessResponse(entitlements.ZpaGroupEntitlements{
		TotalCount:      10,
		ZpaEnableForAll: 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGroup := &entitlements.ZpaGroupEntitlements{
		ZpaEnableForAll: 1,
	}

	result, err := entitlements.UpdateZpaGroupEntitlements(context.Background(), service, updateGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.ZpaEnableForAll)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestEntitlements_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZdxGroupEntitlements JSON marshaling", func(t *testing.T) {
		group := entitlements.ZdxGroupEntitlements{
			TotalCount:                5,
			UpmEnableForAll:           1,
			CollectZdxLocation:        1,
			ComputeDeviceGroupsForZDX: 0,
			LogoutZCCForZDXService:    0,
			UpmGroupList: []entitlements.DeviceGroup{
				{GroupID: 100, GroupName: "Test Group", Active: 1, UpmEnabled: 1},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"totalCount":5`)
		assert.Contains(t, string(data), `"upmEnableForAll":1`)
		assert.Contains(t, string(data), `"collectZdxLocation":1`)
	})

	t.Run("ZdxGroupEntitlements JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 10,
			"upmEnableForAll": 0,
			"collectZdxLocation": 1,
			"computeDeviceGroupsForZDX": 1,
			"logoutZCCForZDXService": 0,
			"upmGroupList": [
				{"groupId": 200, "groupName": "Test", "active": 1, "upmEnabled": 1}
			]
		}`

		var group entitlements.ZdxGroupEntitlements
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 10, group.TotalCount)
		assert.Equal(t, 0, group.UpmEnableForAll)
		assert.Len(t, group.UpmGroupList, 1)
	})

	t.Run("ZpaGroupEntitlements JSON marshaling", func(t *testing.T) {
		group := entitlements.ZpaGroupEntitlements{
			TotalCount:                3,
			ZpaEnableForAll:           1,
			ComputeDeviceGroupsForZPA: 1,
			MachineTunEnabledForAll:   0,
			GroupList: []entitlements.GroupListItem{
				{GroupID: 300, GroupName: "ZPA Test", Active: 1, ZpaEnabled: 1},
			},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"totalCount":3`)
		assert.Contains(t, string(data), `"zpaEnableForAll":1`)
	})

	t.Run("ZpaGroupEntitlements JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 15,
			"zpaEnableForAll": 1,
			"computeDeviceGroupsForZPA": 0,
			"machineTunEnabledForAll": 1,
			"groupList": [
				{"groupId": 400, "groupName": "Access Group", "active": 1, "zpaEnabled": 1}
			]
		}`

		var group entitlements.ZpaGroupEntitlements
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 15, group.TotalCount)
		assert.Equal(t, 1, group.ZpaEnableForAll)
		assert.Len(t, group.GroupList, 1)
	})

	t.Run("DeviceGroup JSON marshaling", func(t *testing.T) {
		dg := entitlements.DeviceGroup{
			GroupID:    123,
			GroupName:  "Device Group",
			Active:     1,
			AuthType:   "SAML",
			UpmEnabled: 1,
		}

		data, err := json.Marshal(dg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"groupId":123`)
		assert.Contains(t, string(data), `"groupName":"Device Group"`)
	})
}
