// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/remove_devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestRemoveDevices_SoftRemove_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/removeDevices"

	server.On("POST", path, common.SuccessResponse(remove_devices.RemoveDevicesResponse{
		DevicesRemoved: 2,
		ErrorMsg:       "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	request := remove_devices.RemoveDevicesRequest{
		Udids:    []string{"device-001", "device-002"},
		UserName: "user1@example.com",
		OsType:   1,
	}

	result, err := remove_devices.SoftRemoveDevices(context.Background(), service, request, 100)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.DevicesRemoved)
	assert.Empty(t, result.ErrorMsg)
}

func TestRemoveDevices_ForceRemove_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/forceRemoveDevices"

	server.On("POST", path, common.SuccessResponse(remove_devices.RemoveDevicesResponse{
		DevicesRemoved: 1,
		ErrorMsg:       "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	request := remove_devices.RemoveDevicesRequest{
		Udids: []string{"device-003"},
	}

	result, err := remove_devices.ForceRemoveDevices(context.Background(), service, request, 100)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.DevicesRemoved)
}

func TestRemoveDevices_RemoveMachineTunnel_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/removeMachineTunnel"

	server.On("POST", path, common.SuccessResponse(remove_devices.RemoveDevicesResponse{
		DevicesRemoved: 1,
		ErrorMsg:       "",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	hostNames := []string{"host1.example.com"}
	machineToken := "token-123"

	result, err := remove_devices.RemoveMachineTunnel(context.Background(), service, hostNames, machineToken)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, result.DevicesRemoved)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestRemoveDevices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("RemoveDevicesRequest JSON marshaling", func(t *testing.T) {
		request := remove_devices.RemoveDevicesRequest{
			ClientConnectorVersion: []string{"4.2.0.100"},
			OsType:                 1,
			Udids:                  []string{"udid-001", "udid-002"},
			UserName:               "user@example.com",
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"clientConnectorVersion":["4.2.0.100"]`)
		assert.Contains(t, string(data), `"osType":1`)
		assert.Contains(t, string(data), `"udids":["udid-001","udid-002"]`)
		assert.Contains(t, string(data), `"userName":"user@example.com"`)
	})

	t.Run("RemoveDevicesRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"clientConnectorVersion": ["4.3.0.200"],
			"osType": 2,
			"udids": ["udid-003"],
			"userName": "user3@example.com"
		}`

		var request remove_devices.RemoveDevicesRequest
		err := json.Unmarshal([]byte(jsonData), &request)
		require.NoError(t, err)

		assert.Len(t, request.ClientConnectorVersion, 1)
		assert.Equal(t, "4.3.0.200", request.ClientConnectorVersion[0])
		assert.Len(t, request.Udids, 1)
		assert.Equal(t, "user3@example.com", request.UserName)
	})

	t.Run("RemoveDevicesResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"devicesRemoved": 5,
			"errorMsg": ""
		}`

		var response remove_devices.RemoveDevicesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 5, response.DevicesRemoved)
		assert.Empty(t, response.ErrorMsg)
	})

	t.Run("RemoveDevicesResponse with error", func(t *testing.T) {
		jsonData := `{
			"devicesRemoved": 0,
			"errorMsg": "Device not found"
		}`

		var response remove_devices.RemoveDevicesResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 0, response.DevicesRemoved)
		assert.Equal(t, "Device not found", response.ErrorMsg)
	})
}
