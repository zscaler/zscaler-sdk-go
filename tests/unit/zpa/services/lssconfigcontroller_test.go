// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/lssconfigcontroller"
)

func TestLSSConfigController_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	configID := "lss-12345"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/lssConfig/" + configID

	server.On("GET", path, common.SuccessResponse(lssconfigcontroller.LSSResource{
		ID: configID,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := lssconfigcontroller.Get(context.Background(), service, configID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, configID, result.ID)
}

func TestLSSConfigController_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/lssConfig"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []lssconfigcontroller.LSSResource{{ID: "lss-001"}, {ID: "lss-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := lssconfigcontroller.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
