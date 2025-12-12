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
		ID:             "schedule-123",
		Enabled:        true,
		Frequency:      "7",
		DeleteDisabled: false,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := serviceedgeschedule.GetSchedule(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.Enabled)
}
