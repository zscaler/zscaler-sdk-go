// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/privilegedremoteaccess/praconsole"
)

func TestPRAConsole_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	consoleID := "console-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole/" + consoleID

	server.On("GET", path, common.SuccessResponse(praconsole.PRAConsole{
		ID:   consoleID,
		Name: "Test Console",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praconsole.Get(context.Background(), service, consoleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, consoleID, result.ID)
}

func TestPRAConsole_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []praconsole.PRAConsole{{ID: "console-001"}, {ID: "console-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praconsole.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
