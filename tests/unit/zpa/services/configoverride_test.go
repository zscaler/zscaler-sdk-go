// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/config_override"
)

func TestConfigOverride_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	configID := "config-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/configOverrides/" + configID

	server.On("GET", path, common.SuccessResponse(config_override.ConfigOverrides{
		ConfigKey:   "TEST_KEY",
		ConfigValue: "test-value",
		TargetType:  "CUSTOMER",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := config_override.Get(context.Background(), service, configID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "TEST_KEY", result.ConfigKey)
	assert.Equal(t, "test-value", result.ConfigValue)
}

func TestConfigOverride_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/configOverrides"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []config_override.ConfigOverrides{
			{ConfigKey: "KEY_A", ConfigValue: "value_a"},
			{ConfigKey: "KEY_B", ConfigValue: "value_b"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := config_override.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestConfigOverride_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/configOverrides"

	server.On("POST", path, common.SuccessResponse(config_override.ConfigOverrides{
		ConfigKey:   "NEW_KEY",
		ConfigValue: "new-value",
		TargetType:  "CUSTOMER",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newConfig := &config_override.ConfigOverrides{
		ConfigKey:   "NEW_KEY",
		ConfigValue: "new-value",
		TargetType:  "CUSTOMER",
	}

	result, _, err := config_override.Create(context.Background(), service, newConfig)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "NEW_KEY", result.ConfigKey)
}

func TestConfigOverride_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	configID := "config-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/configOverrides/" + configID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateConfig := &config_override.ConfigOverrides{
		ConfigKey:   "UPDATED_KEY",
		ConfigValue: "updated-value",
	}

	resp, err := config_override.Update(context.Background(), service, configID, updateConfig)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
