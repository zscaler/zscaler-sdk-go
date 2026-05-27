// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegmentinspection"
)

func TestApplicationSegmentInspection_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "insp-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("GET", path, common.SuccessResponse(applicationsegmentinspection.AppSegmentInspection{
		ID:   appID,
		Name: "Test Inspection App",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentinspection.Get(context.Background(), service, appID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, appID, result.ID)
}

func TestApplicationSegmentInspection_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	// GetAll filters results - only returns items where len(InspectionAppDto) > 0
	// So we need to include InspectionAppDto in the mock data
	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentinspection.AppSegmentInspection{
			{
				ID:   "insp-001",
				Name: "Inspection App 1",
				InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
					{ID: "dto-1", Name: "Inspection DTO 1"},
				},
			},
			{
				ID:   "insp-002",
				Name: "Inspection App 2",
				InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
					{ID: "dto-2", Name: "Inspection DTO 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentinspection.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegmentInspection_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appName := "Production Inspection App"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegmentinspection.AppSegmentInspection{
			{
				ID:   "insp-001",
				Name: "Other App",
				InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
					{ID: "dto-1", Name: "Inspection DTO 1"},
				},
			},
			{
				ID:   "insp-002",
				Name: appName,
				InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
					{ID: "dto-2", Name: "Inspection DTO 2"},
				},
			},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegmentinspection.GetByName(context.Background(), service, appName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "insp-002", result.ID)
	assert.Equal(t, appName, result.Name)
}

func TestApplicationSegmentInspection_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("POST", path, common.SuccessResponse(applicationsegmentinspection.AppSegmentInspection{
		ID:   "new-insp-123",
		Name: "New Inspection App",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newApp := applicationsegmentinspection.AppSegmentInspection{
		Name:           "New Inspection App",
		SegmentGroupID: "sg-001",
	}

	result, _, err := applicationsegmentinspection.Create(context.Background(), service, newApp)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-insp-123", result.ID)
}

func TestApplicationSegmentInspection_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "insp-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	// Update calls Get first, so we need to mock the GET request too
	server.On("GET", path, common.SuccessResponse(applicationsegmentinspection.AppSegmentInspection{
		ID:   appID,
		Name: "Original Inspection App",
	}))

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateApp := &applicationsegmentinspection.AppSegmentInspection{
		ID:   appID,
		Name: "Updated Inspection App",
	}

	resp, err := applicationsegmentinspection.Update(context.Background(), service, appID, updateApp)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegmentInspection_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "insp-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := applicationsegmentinspection.Delete(context.Background(), service, appID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegmentInspection_GetByName_NotFound_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	path := common.ZPAPath(api.CustomerID, "application")

	api.On("GET", path, common.SuccessResponse(common.ZPAList([]applicationsegmentinspection.AppSegmentInspection{
		{ID: "x", Name: "Other", InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{{ID: "d", Name: "D"}}},
	})))

	result, _, err := applicationsegmentinspection.GetByName(context.Background(), api.Service, "Unknown")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestApplicationSegmentInspection_Update_AppsConfigInjection_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	appID := "insp-parent"
	path := common.ZPAPath(api.CustomerID, "application", appID)

	api.On("GET", path, common.SuccessResponse(applicationsegmentinspection.AppSegmentInspection{
		ID:   appID,
		Name: "Parent",
		InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
			{ID: "inspect-child-77", Name: "Child Inspect"},
		},
	}))
	api.On("PUT", path, common.NoContentResponse())

	upd := &applicationsegmentinspection.AppSegmentInspection{
		CommonAppsDto: applicationsegmentinspection.CommonAppsDto{
			AppsConfig: []applicationsegmentinspection.AppsConfig{
				{Name: "Child Inspect", Enabled: true, InspectAppID: "will-be-overwritten"},
			},
		},
	}

	resp, err := applicationsegmentinspection.Update(context.Background(), api.Service, appID, upd)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestApplicationSegmentInspection_Update_ClearsEmptyAppsConfig_SDK(t *testing.T) {
	api := common.NewZPATest(t)
	appID := "insp-clear"
	path := common.ZPAPath(api.CustomerID, "application", appID)

	api.On("GET", path, common.SuccessResponse(applicationsegmentinspection.AppSegmentInspection{
		ID:   appID,
		Name: "Parent",
		InspectionAppDto: []applicationsegmentinspection.InspectionAppDto{
			{ID: "i1", Name: "X"},
		},
	}))
	api.On("PUT", path, common.NoContentResponse())

	resp, err := applicationsegmentinspection.Update(context.Background(), api.Service, appID,
		&applicationsegmentinspection.AppSegmentInspection{
			CommonAppsDto: applicationsegmentinspection.CommonAppsDto{},
		})
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
