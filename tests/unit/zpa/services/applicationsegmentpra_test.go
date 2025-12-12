// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentpra"
)

func TestApplicationSegmentPRA_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "pra-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("GET", path, common.SuccessResponse(applicationsegmentpra.AppSegmentPRA{
		ID:   appID,
		Name: "Test PRA App",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentpra.Get(context.Background(), service, appID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, appID, result.ID)
}

func TestApplicationSegmentPRA_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	// GetAll filters results - only returns items where len(PRAApps) > 0
	// So we need to include PRAApps in the mock data
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentpra.AppSegmentPRA{
			{
				ID:   "pra-001",
				Name: "PRA App 1",
				PRAApps: []applicationsegmentpra.PRAApps{
					{ID: "praapp-1", Name: "PRA Sub App 1"},
				},
			},
			{
				ID:   "pra-002",
				Name: "PRA App 2",
				PRAApps: []applicationsegmentpra.PRAApps{
					{ID: "praapp-2", Name: "PRA Sub App 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentpra.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegmentPRA_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appName := "Production PRA App"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentpra.AppSegmentPRA{
			{
				ID:   "pra-001",
				Name: "Other App",
				PRAApps: []applicationsegmentpra.PRAApps{
					{ID: "praapp-1", Name: "PRA Sub App 1"},
				},
			},
			{
				ID:   "pra-002",
				Name: appName,
				PRAApps: []applicationsegmentpra.PRAApps{
					{ID: "praapp-2", Name: "PRA Sub App 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentpra.GetByName(context.Background(), service, appName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "pra-002", result.ID)
	assert.Equal(t, appName, result.Name)
}

func TestApplicationSegmentPRA_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("POST", path, common.SuccessResponse(applicationsegmentpra.AppSegmentPRA{
		ID:   "new-pra-123",
		Name: "New PRA App",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newApp := applicationsegmentpra.AppSegmentPRA{
		Name:           "New PRA App",
		SegmentGroupID: "sg-001",
	}

	result, _, err := applicationsegmentpra.Create(context.Background(), service, newApp)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-pra-123", result.ID)
}

func TestApplicationSegmentPRA_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "pra-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	// Update calls Get first, so we need to mock the GET request too
	server.On("GET", path, common.SuccessResponse(applicationsegmentpra.AppSegmentPRA{
		ID:   appID,
		Name: "Original PRA App",
	}))

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateApp := &applicationsegmentpra.AppSegmentPRA{
		ID:   appID,
		Name: "Updated PRA App",
	}

	resp, err := applicationsegmentpra.Update(context.Background(), service, appID, updateApp)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegmentPRA_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "pra-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := applicationsegmentpra.Delete(context.Background(), service, appID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}
