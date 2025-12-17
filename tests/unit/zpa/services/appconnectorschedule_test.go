// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorschedule"
)

func TestAppConnectorSchedule_GetSchedule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connectorSchedule"

	server.On("GET", path, common.SuccessResponse(appconnectorschedule.AssistantSchedule{
		ID:                testCustomerID,
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
		DeleteDisabled:    false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := appconnectorschedule.GetSchedule(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
	assert.Equal(t, "7", result.FrequencyInterval)
}

func TestAppConnectorSchedule_CreateSchedule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connectorSchedule"

	server.On("POST", path, common.SuccessResponse(appconnectorschedule.AssistantSchedule{
		ID:                testCustomerID,
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
		DeleteDisabled:    false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newSchedule := appconnectorschedule.AssistantSchedule{
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
	}

	result, _, err := appconnectorschedule.CreateSchedule(context.Background(), service, newSchedule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
}

func TestAppConnectorSchedule_UpdateSchedule_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	schedulerID := "schedule-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/connectorSchedule/" + schedulerID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateSchedule := &appconnectorschedule.AssistantSchedule{
		ID:                schedulerID,
		Enabled:           true,
		FrequencyInterval: "14",
		Frequency:         "14",
	}

	resp, err := appconnectorschedule.UpdateSchedule(context.Background(), service, schedulerID, updateSchedule)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
