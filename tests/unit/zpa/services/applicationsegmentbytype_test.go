// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbytype"
)

func TestApplicationSegmentByType_GetByApplicationType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appName := "TestApp"
	appType := "BROWSER_ACCESS"
	// The SDK uses this path pattern with query params
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/getAppsByType"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentbytype.AppSegmentBaseAppDto{
			{ID: "app-001", Name: appName, Enabled: true},
			{ID: "app-002", Name: "OtherApp", Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentbytype.GetByApplicationType(context.Background(), service, appName, appType, false)

	require.NoError(t, err)
	require.NotNil(t, result)
	// GetByApplicationType filters results
}

func TestApplicationSegmentByType_DeleteByApplicationType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	appType := "BROWSER_ACCESS"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/deleteAppByType"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := applicationsegmentbytype.DeleteByApplicationType(context.Background(), service, appID, appType)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

