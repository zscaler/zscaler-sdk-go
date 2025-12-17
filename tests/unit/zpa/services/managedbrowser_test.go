// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/managed_browser"
)

func TestManagedBrowser_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/managedBrowser"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []managed_browser.ManagedBrowserProfile{
			{ID: "mb-001", Name: "Browser Profile 1"},
			{ID: "mb-002", Name: "Browser Profile 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := managed_browser.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestManagedBrowser_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	browserName := "Production Browser"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/managedBrowser"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []managed_browser.ManagedBrowserProfile{
			{ID: "mb-001", Name: "Other Browser"},
			{ID: "mb-002", Name: browserName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := managed_browser.GetByName(context.Background(), service, browserName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "mb-002", result.ID)
	assert.Equal(t, browserName, result.Name)
}

func TestManagedBrowser_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/managedBrowser"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []managed_browser.ManagedBrowserProfile{},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := managed_browser.GetByName(context.Background(), service, "NonExistent")

	require.Error(t, err)
	assert.Nil(t, result)
}
