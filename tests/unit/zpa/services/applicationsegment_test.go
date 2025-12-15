// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/applicationsegment"
)

func TestApplicationSegment_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("GET", path, common.SuccessResponse(applicationsegment.ApplicationSegmentResource{
		ID:             appID,
		Name:           "Test Application",
		Description:    "Test description",
		Enabled:        true,
		DomainNames:    []string{"app.example.com"},
		SegmentGroupID: "sg-001",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := applicationsegment.Get(context.Background(), service, appID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, appID, result.ID)
	assert.Equal(t, "Test Application", result.Name)
}

func TestApplicationSegment_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appName := "Production App"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegment.ApplicationSegmentResource{
			{ID: "app-001", Name: "Other App", Enabled: true},
			{ID: "app-002", Name: appName, Enabled: true},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetByName(context.Background(), service, appName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "app-002", result.ID)
}

func TestApplicationSegment_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("POST", path, common.SuccessResponse(applicationsegment.ApplicationSegmentResource{
		ID:   "new-app-123",
		Name: "New Application",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	newApp := applicationsegment.ApplicationSegmentResource{
		Name:           "New Application",
		SegmentGroupID: "sg-001",
	}

	result, _, err := applicationsegment.Create(context.Background(), service, newApp)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-app-123", result.ID)
}

func TestApplicationSegment_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	updateApp := applicationsegment.ApplicationSegmentResource{
		ID:   appID,
		Name: "Updated Application",
	}

	resp, err := applicationsegment.Update(context.Background(), service, appID, updateApp)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegment_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := applicationsegment.Delete(context.Background(), service, appID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegment_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []applicationsegment.ApplicationSegmentResource{
			{ID: "app-001", Name: "App 1"},
			{ID: "app-002", Name: "App 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegment_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.Get(context.Background(), service, appID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestApplicationSegment_GetApplicationSummary_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/summary"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []map[string]interface{}{
			{"id": "app-001", "name": "App 1"},
			{"id": "app-002", "name": "App 2"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetApplicationSummary(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegment_GetApplicationCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/configured/count"

	server.On("GET", path, common.SuccessResponse([]applicationsegment.ApplicationCountResponse{
		{AppsConfigured: "10", ConfiguredDateInEpochSeconds: "1704067200"},
		{AppsConfigured: "5", ConfiguredDateInEpochSeconds: "1704153600"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetApplicationCount(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "10", result[0].AppsConfigured)
}

func TestApplicationSegment_GetCurrentAndMaxLimit_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/count/currentAndMaxLimit"

	server.On("GET", path, common.SuccessResponse(applicationsegment.ApplicationCurrentMaxLimitResponse{
		MaxAppsLimit:     "1000",
		CurrentAppsCount: "250",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetCurrentAndMaxLimit(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "1000", result.MaxAppsLimit)
	assert.Equal(t, "250", result.CurrentAppsCount)
}

func TestApplicationSegment_GetApplicationMappings_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/mappings"

	server.On("GET", path, common.SuccessResponse([]applicationsegment.ApplicationMappings{
		{Name: "server-group-1", Type: "SERVER_GROUP"},
		{Name: "server-group-2", Type: "SERVER_GROUP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetApplicationMappings(context.Background(), service, appID)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "server-group-1", result[0].Name)
}

func TestApplicationSegment_GetWeightedLoadBalancerConfig_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/weightedLbConfig"

	server.On("GET", path, common.SuccessResponse(applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         appID,
		WeightedLoadBalancing: true,
		ApplicationToServerGroupMaps: []applicationsegment.ApplicationToServerGroupMapping{
			{ID: "sg-001", Name: "Server Group 1", Weight: "50", Passive: false},
			{ID: "sg-002", Name: "Server Group 2", Weight: "50", Passive: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := applicationsegment.GetWeightedLoadBalancerConfig(context.Background(), service, appID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.True(t, result.WeightedLoadBalancing)
	assert.Len(t, result.ApplicationToServerGroupMaps, 2)
}

func TestApplicationSegment_UpdateWeightedLoadBalancerConfig_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	appID := "app-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/" + appID + "/weightedLbConfig"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	config := applicationsegment.WeightedLoadBalancerConfig{
		ApplicationID:         appID,
		WeightedLoadBalancing: true,
	}

	_, resp, err := applicationsegment.UpdateWeightedLoadBalancerConfig(context.Background(), service, appID, config)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegment_GetMultiMatchUnsupportedReferences_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/multimatchUnsupportedReferences"

	server.On("POST", path, common.SuccessResponse([]applicationsegment.MultiMatchUnsupportedReferencesResponse{
		{ID: "ref-001", AppSegmentName: "App 1"},
		{ID: "ref-002", AppSegmentName: "App 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	domainNames := applicationsegment.MultiMatchUnsupportedReferencesPayload{"example.com", "test.com"}
	result, _, err := applicationsegment.GetMultiMatchUnsupportedReferences(context.Background(), service, domainNames)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestApplicationSegment_UpdatebulkUpdateMultiMatch_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/bulkUpdateMultiMatch"

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	bulkUpdate := applicationsegment.BulkUpdateMultiMatchPayload{
		ApplicationIDs: []int{1, 2, 3},
		MatchStyle:     "MULTIPLE",
	}

	resp, err := applicationsegment.UpdatebulkUpdateMultiMatch(context.Background(), service, bulkUpdate)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestApplicationSegment_ApplicationValidation_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/application/validate"

	server.On("POST", path, common.SuccessResponse(applicationsegment.ApplicationValidationError{
		ID:     "",
		Reason: "",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	validateReq := applicationsegment.ApplicationSegmentResource{
		Name:        "Test App",
		DomainNames: []string{"app.example.com"},
	}

	result, _, err := applicationsegment.ApplicationValidation(context.Background(), service, validateReq)

	require.NoError(t, err)
	require.NotNil(t, result)
}
