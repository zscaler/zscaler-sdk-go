// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/browser_protection"
)

func TestBrowserProtection_GetActive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/browserProtection/active"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []browser_protection.BrowserProtection{{ID: "bp-001"}, {ID: "bp-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := browser_protection.GetActiveBrowserProtectionProfile(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestBrowserProtection_GetProfile_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/browserProtection"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []browser_protection.BrowserProtection{{ID: "bp-001"}, {ID: "bp-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := browser_protection.GetBrowserProtectionProfile(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
