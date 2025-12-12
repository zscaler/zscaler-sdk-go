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

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/activeBrowserProtectionProfile
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/activeBrowserProtectionProfile"

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

	// Correct path: /zpa/mgmtconfig/v1/admin/customers/{customerId}/browserProtectionProfile
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/browserProtectionProfile"

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

func TestBrowserProtection_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	profileName := "Production Profile"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/browserProtectionProfile"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []browser_protection.BrowserProtection{
			{ID: "bp-001", Name: "Other Profile"},
			{ID: "bp-002", Name: profileName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := browser_protection.GetBrowserProtectionProfileByName(context.Background(), service, profileName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "bp-002", result.ID)
	assert.Equal(t, profileName, result.Name)
}
