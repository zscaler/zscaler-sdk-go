// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/remote_assistance"
)

const remoteAssistancePath = "/zia/api/v1/remoteAssistance"

func TestRemoteAssistance_GetRemoteAssistance_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", remoteAssistancePath, common.SuccessResponse(remote_assistance.RemoteAssistance{
		ViewOnlyUntil:       1717200000000,
		FullAccessUntil:     1725148800000,
		UsernameObfuscated:  true,
		DeviceInfoObfuscate: false,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := remote_assistance.GetRemoteAssistance(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, int64(1717200000000), result.ViewOnlyUntil)
	assert.True(t, result.UsernameObfuscated)
	assert.False(t, result.DeviceInfoObfuscate)
}

func TestRemoteAssistance_GetRemoteAssistance_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", remoteAssistancePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := remote_assistance.GetRemoteAssistance(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestRemoteAssistance_UpdateRemoteAssistance_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", remoteAssistancePath, common.SuccessResponse(remote_assistance.RemoteAssistance{
		ViewOnlyUntil:       1717200000000,
		FullAccessUntil:     1725148800000,
		UsernameObfuscated:  false,
		DeviceInfoObfuscate: true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	settings := remote_assistance.RemoteAssistance{
		ViewOnlyUntil:       1717200000000,
		FullAccessUntil:     1725148800000,
		UsernameObfuscated:  false,
		DeviceInfoObfuscate: true,
	}

	result, _, err := remote_assistance.UpdateRemoteAssistance(context.Background(), service, settings)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.DeviceInfoObfuscate)
}

func TestRemoteAssistance_Structure(t *testing.T) {
	t.Parallel()

	jsonData := `{
		"viewOnlyUntil": 1717200000000,
		"fullAccessUntil": 1725148800000,
		"usernameObfuscated": true,
		"deviceInfoObfuscate": false
	}`

	var settings remote_assistance.RemoteAssistance
	err := json.Unmarshal([]byte(jsonData), &settings)
	require.NoError(t, err)

	assert.True(t, settings.UsernameObfuscated)
}
