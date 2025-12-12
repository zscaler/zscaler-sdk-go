// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/step_up_auth"
)

func TestStepUpAuth_GetStepupAuthLevel_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/stepupauthlevel"

	// The actual function returns []string (list of auth level names/IDs)
	server.On("GET", path, common.SuccessResponse([]string{"Level1", "Level2", "Level3"}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := step_up_auth.GetStepupAuthLevel(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Contains(t, result, "Level1")
}

func TestStepUpAuth_GetStepupAuthLevel_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/stepupauthlevel"

	server.On("GET", path, common.SuccessResponse([]string{}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := step_up_auth.GetStepupAuthLevel(context.Background(), service)

	require.NoError(t, err)
	assert.Empty(t, result)
}
