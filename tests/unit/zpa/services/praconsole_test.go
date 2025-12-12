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

func TestPRAConsole_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	consoleName := "Production Console"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []praconsole.PRAConsole{
			{ID: "console-001", Name: "Other Console"},
			{ID: "console-002", Name: consoleName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := praconsole.GetByName(context.Background(), service, consoleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "console-002", result.ID)
	assert.Equal(t, consoleName, result.Name)
}

func TestPRAConsole_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole"

	server.On("POST", path, common.SuccessResponse(praconsole.PRAConsole{
		ID:   "new-console-123",
		Name: "New Console",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newConsole := &praconsole.PRAConsole{
		Name: "New Console",
	}

	result, _, err := praconsole.Create(context.Background(), service, newConsole)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-console-123", result.ID)
}

func TestPRAConsole_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	consoleID := "console-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole/" + consoleID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateConsole := &praconsole.PRAConsole{
		ID:   consoleID,
		Name: "Updated Console",
	}

	resp, err := praconsole.Update(context.Background(), service, consoleID, updateConsole)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPRAConsole_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	consoleID := "console-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/praConsole/" + consoleID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := praconsole.Delete(context.Background(), service, consoleID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
