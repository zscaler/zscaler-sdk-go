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
		"list":       []managed_browser.ManagedBrowserProfile{{ID: "mb-001"}, {ID: "mb-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := managed_browser.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
