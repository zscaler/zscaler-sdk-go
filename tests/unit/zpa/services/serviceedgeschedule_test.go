// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgeschedule"
)

func TestServiceEdgeSchedule_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeSchedule"

	server.On("GET", path, common.SuccessResponse(serviceedgeschedule.AssistantSchedule{
		ID:                testCustomerID,
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
		DeleteDisabled:    false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgeschedule.GetSchedule(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
	assert.Equal(t, "7", result.FrequencyInterval)
}

func TestServiceEdgeSchedule_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeSchedule"

	server.On("POST", path, common.SuccessResponse(serviceedgeschedule.AssistantSchedule{
		ID:                testCustomerID,
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
		DeleteDisabled:    false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newSchedule := serviceedgeschedule.AssistantSchedule{
		Enabled:           true,
		FrequencyInterval: "7",
		Frequency:         "7",
	}

	result, _, err := serviceedgeschedule.CreateSchedule(context.Background(), service, newSchedule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
}

func TestServiceEdgeSchedule_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	scheduleID := testCustomerID
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/serviceEdgeSchedule/" + scheduleID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateSchedule := &serviceedgeschedule.AssistantSchedule{
		ID:                scheduleID,
		Enabled:           true,
		FrequencyInterval: "14",
		Frequency:         "14",
	}

	resp, err := serviceedgeschedule.UpdateSchedule(context.Background(), service, scheduleID, updateSchedule)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
