// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/custom_config_controller"
)

func TestCustomConfigController_CheckZiaCloudConfig_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/config/isZiaCloudConfigAvailable"

	server.On("GET", path, common.SuccessResponse(true))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := custom_config_controller.CheckZiaCloudConfig(context.Background(), service)

	require.NoError(t, err)
	assert.True(t, result)
}

func TestCustomConfigController_GetZIACloudConfig_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/config/ziaCloud"

	server.On("GET", path, common.SuccessResponse(custom_config_controller.ZIACloudConfig{
		ZIACloudDomain: "zscaler.net",
		ZIAUsername:    "admin@example.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := custom_config_controller.GetZIACloudConfig(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "zscaler.net", result.ZIACloudDomain)
}

func TestCustomConfigController_AddZIACloudConfig_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/config/ziaCloud"

	server.On("POST", path, common.SuccessResponse(custom_config_controller.ZIACloudConfig{
		ZIACloudDomain: "zscaler.net",
		ZIAUsername:    "newadmin@example.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newConfig := &custom_config_controller.ZIACloudConfig{
		ZIACloudDomain: "zscaler.net",
		ZIAUsername:    "newadmin@example.com",
	}

	result, _, err := custom_config_controller.AddZIACloudConfig(context.Background(), service, newConfig)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "newadmin@example.com", result.ZIAUsername)
}

func TestCustomConfigController_GetSessionTerminationOnReath_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/config/sessionTerminationOnReauth"

	server.On("GET", path, common.SuccessResponse(custom_config_controller.SessionTerminationOnReath{
		SessionTerminationOnReauth:             true,
		AllowDisableSessionTerminationOnReauth: false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := custom_config_controller.GetSessionTerminationOnReath(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.SessionTerminationOnReauth)
}

func TestCustomConfigController_UpdateSessionTerminationOnReath_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/config/sessionTerminationOnReauth"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateSettings := &custom_config_controller.SessionTerminationOnReath{
		SessionTerminationOnReauth: false,
	}

	resp, err := custom_config_controller.UpdateSessionTerminationOnReath(context.Background(), service, updateSettings)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
