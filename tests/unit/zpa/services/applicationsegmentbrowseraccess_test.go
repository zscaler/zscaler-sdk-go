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

	// GetAll filters results - only returns items where len(ClientlessApps) > 0
	// So we need to include ClientlessApps in the mock data
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentbrowseraccess.BrowserAccess{
			{
				ID:   "ba-001",
				Name: "Browser Access 1",
				ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{
					{ID: "clientless-1", Name: "Clientless App 1"},
				},
			},
			{
				ID:   "ba-002",
				Name: "Browser Access 2",
				ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{
					{ID: "clientless-2", Name: "Clientless App 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentbrowseraccess.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegmentBrowserAccess_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appName := "Production Browser Access"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentbrowseraccess.BrowserAccess{
			{
				ID:   "ba-001",
				Name: "Other App",
				ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{
					{ID: "clientless-1", Name: "Clientless App 1"},
				},
			},
			{
				ID:   "ba-002",
				Name: appName,
				ClientlessApps: []applicationsegmentbrowseraccess.ClientlessApps{
					{ID: "clientless-2", Name: "Clientless App 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentbrowseraccess.GetByName(context.Background(), service, appName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "ba-002", result.ID)
	assert.Equal(t, appName, result.Name)
}

func TestApplicationSegmentBrowserAccess_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("POST", path, common.SuccessResponse(applicationsegmentbrowseraccess.BrowserAccess{
		ID:   "new-ba-123",
		Name: "New Browser Access",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newApp := applicationsegmentbrowseraccess.BrowserAccess{
		Name:           "New Browser Access",
		SegmentGroupID: "sg-001",
	}

	result, _, err := applicationsegmentbrowseraccess.Create(context.Background(), service, newApp)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-ba-123", result.ID)
}

func TestApplicationSegmentBrowserAccess_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "ba-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	// Update calls Get first, so we need to mock the GET request too
	server.On("GET", path, common.SuccessResponse(applicationsegmentbrowseraccess.BrowserAccess{
		ID:   appID,
		Name: "Original Browser Access",
	}))

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateApp := &applicationsegmentbrowseraccess.BrowserAccess{
		ID:   appID,
		Name: "Updated Browser Access",
	}

	resp, err := applicationsegmentbrowseraccess.Update(context.Background(), service, appID, updateApp)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegmentBrowserAccess_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "ba-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := applicationsegmentbrowseraccess.Delete(context.Background(), service, appID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
