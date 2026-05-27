// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/shadowitreport"
)

func TestShadowITReport_GetAllCloudAppsLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplications/lite"
	server.On("GET", path, common.SuccessResponse([]shadowitreport.CloudApplicationsAndCustomTags{
		{ID: 1, Name: "Office365"},
		{ID: 2, Name: "Dropbox"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := shadowitreport.GetAllCloudAppsLite(context.Background(), service, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestShadowITReport_GetAllCloudAppsLite_Paginated_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplications/lite"
	server.On("GET", path, common.SuccessResponse([]shadowitreport.CloudApplicationsAndCustomTags{
		{ID: 1, Name: "Office365"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	pageNumber := 1
	limit := 10
	result, err := shadowitreport.GetAllCloudAppsLite(context.Background(), service, &pageNumber, &limit)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestShadowITReport_GetAllCustomTags_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/customTags"
	server.On("GET", path, common.SuccessResponse([]shadowitreport.CloudApplicationsAndCustomTags{
		{ID: 100, Name: "Sanctioned"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := shadowitreport.GetAllCustomTags(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestShadowITReport_CreateCloudApplicationsExport_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/shadowIT/applications/export"
	server.On("POST", path, common.RawResponse([]byte("app,name\n1,Office365"), 200, map[string]string{
		"Content-Type": "text/csv",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	exportReq := shadowitreport.CloudApplicationsExport{
		Duration: "LAST_7_DAYS",
	}

	resp, err := shadowitreport.CreateCloudApplicationsExport(context.Background(), service, exportReq)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestShadowITReport_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplications/bulkUpdate"
	server.On("PUT", path, common.SuccessResponse(shadowitreport.ApplicationBulkUpdate{
		SanctionedState: "SANCTIONED",
		ApplicationIDs:  []int{1, 2},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &shadowitreport.ApplicationBulkUpdate{
		SanctionedState: "SANCTIONED",
		ApplicationIDs:  []int{1, 2},
	}

	result, err := shadowitreport.Update(context.Background(), service, update)
	require.NoError(t, err)
	assert.Equal(t, "SANCTIONED", result.SanctionedState)
}

func TestShadowITReport_Update_NoContent_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplications/bulkUpdate"
	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := shadowitreport.Update(context.Background(), service, &shadowitreport.ApplicationBulkUpdate{
		SanctionedState: "SANCTIONED",
	})
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestShadowITReport_CreateCloudApplicationsExportCSV_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/shadowIT/applications/USER/exportCsv"
	server.On("POST", path, common.SuccessResponse(shadowitreport.CloudApplicationsExportCSV{
		Duration: "LAST_7_DAYS",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := shadowitreport.CreateCloudApplicationsExportCSV(context.Background(), service, "USER", &shadowitreport.CloudApplicationsExportCSV{
		Duration: "LAST_7_DAYS",
	})
	require.NoError(t, err)
	assert.Equal(t, "LAST_7_DAYS", result.Duration)
}

func TestShadowITReport_CreateCloudApplicationsExportCSV_InvalidEntity_SDK(t *testing.T) {
	service, err := common.CreateTestService(context.Background(), common.NewTestServer(), "123456")
	require.NoError(t, err)

	_, _, err = shadowitreport.CreateCloudApplicationsExportCSV(context.Background(), service, "INVALID", &shadowitreport.CloudApplicationsExportCSV{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid entity value")
}

func TestShadowITReport_GetAllCloudAppsLite_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/cloudApplications/lite"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := shadowitreport.GetAllCloudAppsLite(context.Background(), service, nil, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}
