// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentbrowseraccess"
)

func TestApplicationSegmentBrowserAccess_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "ba-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("GET", path, common.SuccessResponse(applicationsegmentbrowseraccess.BrowserAccess{
		ID:   appID,
		Name: "Test Browser Access",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentbrowseraccess.Get(context.Background(), service, appID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, appID, result.ID)
}

func TestApplicationSegmentBrowserAccess_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []applicationsegmentbrowseraccess.BrowserAccess{{ID: "ba-001"}, {ID: "ba-002"}},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentbrowseraccess.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}
