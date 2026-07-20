// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_applications"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_custom_apps"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
)

const endPointApplicationsBasePath = "/zia/api/v1/endPointApplications"

// =====================================================
// GetCustomApps
// =====================================================

func TestEndpointApplications_GetCustomApps_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApps"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Empty(t, q.Get("search"))
		assert.Empty(t, q.Get("osType"))
		return common.SuccessResponse([]endpoint_custom_apps.EndpointApplications{
			{ResourceID: 1, ApplicationName: "App A", OsType: "WINDOWS_OS"},
			{ResourceID: 2, ApplicationName: "App B", OsType: "MAC_OS"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_custom_apps.GetCustomApps(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "App A", result[0].ApplicationName)
}

func TestEndpointApplications_GetCustomApps_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApps"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "chrome", q.Get("search"))
		assert.Equal(t, "WINDOWS_OS", q.Get("osType"))
		return common.SuccessResponse([]endpoint_custom_apps.EndpointApplications{
			{ResourceID: 10, ApplicationName: "Chrome", OsType: "WINDOWS_OS"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_custom_apps.GetCustomAppsFilterOptions{
		Search: "chrome",
		OsType: "WINDOWS_OS",
	}
	result, err := endpoint_custom_apps.GetCustomApps(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 10, result[0].ResourceID)
}

func TestEndpointApplications_GetCustomApps_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApps"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_custom_apps.GetCustomApps(context.Background(), service, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetCustomApp
// =====================================================

func TestEndpointApplications_GetCustomApp_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApp/55"

	server.On("GET", path, common.SuccessResponse(endpoint_custom_apps.EndpointApplications{
		ResourceID:      55,
		ApplicationName: "Custom App",
		OsType:          "WINDOWS_OS",
		ApplicationType: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_custom_apps.GetCustomApp(context.Background(), service, 55)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 55, result.ResourceID)
	assert.Equal(t, "Custom App", result.ApplicationName)
	assert.Equal(t, "CUSTOM", result.ApplicationType)
}

func TestEndpointApplications_GetCustomApp_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApp/999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_custom_apps.GetCustomApp(context.Background(), service, 999)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// GetApplicationCount
// =====================================================

func TestEndpointApplications_GetApplicationCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/count"
	server.On("GET", path, common.SuccessResponse(42))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := endpoint_applications.GetApplicationCount(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Equal(t, 42, count)
}

func TestEndpointApplications_GetApplicationCount_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/count"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "app", q.Get("search"))
		assert.Equal(t, "MAC_OS", q.Get("osType"))
		assert.Equal(t, "CUSTOM", q.Get("applicationType"))
		return common.SuccessResponse(7)
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_applications.GetApplicationCountFilterOptions{
		Search:          "app",
		OsType:          "MAC_OS",
		ApplicationType: "CUSTOM",
	}
	count, err := endpoint_applications.GetApplicationCount(context.Background(), service, opts)

	require.NoError(t, err)
	assert.Equal(t, 7, count)
}

func TestEndpointApplications_GetApplicationCount_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/count"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := endpoint_applications.GetApplicationCount(context.Background(), service, nil)

	require.Error(t, err)
	assert.Equal(t, 0, count)
}

// =====================================================
// GetCloudAppsCount
// =====================================================

func TestEndpointApplications_GetCloudAppsCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/cloudApps/count"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "WELL_KNOWN", q.Get("applicationType"))
		return common.SuccessResponse(123)
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_applications.GetApplicationCountFilterOptions{
		ApplicationType: "WELL_KNOWN",
	}
	count, err := endpoint_applications.GetCloudAppsCount(context.Background(), service, opts)

	require.NoError(t, err)
	assert.Equal(t, 123, count)
}

func TestEndpointApplications_GetCloudAppsCount_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/cloudApps/count"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	count, err := endpoint_applications.GetCloudAppsCount(context.Background(), service, nil)

	require.Error(t, err)
	assert.Equal(t, 0, count)
}

// =====================================================
// GetApplicationPolicies
// =====================================================

func TestEndpointApplications_GetApplicationPolicies_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/policies"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		ids := r.URL.Query()["resourceId"]
		assert.ElementsMatch(t, []string{"1", "2"}, ids)
		return common.SuccessResponse([]endpoint_applications.ApplicationPolicies{
			{RuleName: "Rule A", RuleType: "DLP", ID: "100"},
			{RuleName: "Rule B", RuleType: "FIREWALL"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetApplicationPolicies(context.Background(), service, []int{1, 2})

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "Rule A", result[0].RuleName)
	assert.Equal(t, "DLP", result[0].RuleType)
}

func TestEndpointApplications_GetApplicationPolicies_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/policies"
	server.On("GET", path, common.SuccessResponse([]endpoint_applications.ApplicationPolicies{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetApplicationPolicies(context.Background(), service, []int{99})

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestEndpointApplications_GetApplicationPolicies_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/policies"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetApplicationPolicies(context.Background(), service, []int{1})

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetCategoriesWithNonEmptyApps
// =====================================================

func TestEndpointApplications_GetCategoriesWithNonEmptyApps_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/getCategoriesWithNonEmptyApps"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Empty(t, q.Get("search"))
		assert.Empty(t, q.Get("osType"))
		return common.SuccessResponse([]string{"WELLKNOWN", "CUSTOM"})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetCategoriesWithNonEmptyApps(context.Background(), service, nil)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "WELLKNOWN", result[0])
}

func TestEndpointApplications_GetCategoriesWithNonEmptyApps_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/getCategoriesWithNonEmptyApps"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "well", q.Get("search"))
		assert.Equal(t, "WINDOWS_OS", q.Get("osType"))
		return common.SuccessResponse([]string{"WELLKNOWN"})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_applications.GetCategoriesWithNonEmptyAppsFilterOptions{
		Search: "well",
		OsType: "WINDOWS_OS",
	}
	result, err := endpoint_applications.GetCategoriesWithNonEmptyApps(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "WELLKNOWN", result[0])
}

func TestEndpointApplications_GetCategoriesWithNonEmptyApps_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/getCategoriesWithNonEmptyApps"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetCategoriesWithNonEmptyApps(context.Background(), service, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetAllEndpointApplications / GetAllEndpointApplicationsLite
// =====================================================

func TestEndpointApplications_GetAllEndpointApplications_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Empty(t, q.Get("search"))
		return common.SuccessResponse([]map[string]interface{}{
			{"resourceId": 1, "applicationName": "App A"},
			{"resourceId": 2, "applicationName": "App B"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetAllEndpointApplications(context.Background(), service, nil)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ResourceID)
	assert.Equal(t, "App A", result[0].ApplicationName)
}

func TestEndpointApplications_GetAllEndpointApplications_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "chrome", q.Get("search"))
		assert.Equal(t, "WINDOWS_OS", q.Get("osType"))
		assert.Equal(t, "CUSTOM", q.Get("applicationType"))
		return common.SuccessResponse([]map[string]interface{}{
			{"resourceId": 10, "applicationName": "Chrome"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_applications.GetApplicationCountFilterOptions{
		Search:          "chrome",
		OsType:          "WINDOWS_OS",
		ApplicationType: "CUSTOM",
	}
	result, err := endpoint_applications.GetAllEndpointApplications(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 10, result[0].ResourceID)
}

func TestEndpointApplications_GetAllEndpointApplications_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", endPointApplicationsBasePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetAllEndpointApplications(context.Background(), service, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

func TestEndpointApplications_GetAllEndpointApplicationsLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/lite"

	server.On("GET", path, common.SuccessResponse([]map[string]interface{}{
		{"resourceId": 5, "applicationName": "Lite App"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetAllEndpointApplicationsLite(context.Background(), service, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 5, result[0].ResourceID)
}

func TestEndpointApplications_GetAllEndpointApplicationsLite_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/lite"

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "edge", q.Get("search"))
		assert.Equal(t, "MAC_OS", q.Get("osType"))
		return common.SuccessResponse([]map[string]interface{}{
			{"resourceId": 6, "applicationName": "Edge"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	opts := &endpoint_applications.GetApplicationCountFilterOptions{
		Search: "edge",
		OsType: "MAC_OS",
	}
	result, err := endpoint_applications.GetAllEndpointApplicationsLite(context.Background(), service, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 6, result[0].ResourceID)
}

func TestEndpointApplications_GetAllEndpointApplicationsLite_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", endPointApplicationsBasePath+"/lite", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_applications.GetAllEndpointApplicationsLite(context.Background(), service, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// Custom Apps CRUD (endpoint_custom_apps)
// =====================================================

func TestEndpointCustomApps_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApp"

	server.OnFunc("POST", path, func(_ *http.Request, body []byte) common.MockResponse {
		assert.Contains(t, string(body), "My Custom App")
		return common.SuccessResponse(endpoint_resource.EndpointResource{
			ID:   777,
			Name: "My Custom App",
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_custom_apps.Create(context.Background(), service, &endpoint_resource.EndpointResource{
		Name: "My Custom App",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 777, result.ID)
}

func TestEndpointCustomApps_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", endPointApplicationsBasePath+"/customApp", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_custom_apps.Create(context.Background(), service, &endpoint_resource.EndpointResource{
		Name: "My Custom App",
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointCustomApps_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", endPointApplicationsBasePath+"/customApp/777", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_custom_apps.Update(context.Background(), service, 777, &endpoint_resource.EndpointResource{
		ID: 777,
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointCustomApps_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", endPointApplicationsBasePath+"/customApp/999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_custom_apps.Delete(context.Background(), service, 999)

	require.Error(t, err)
}

func TestEndpointCustomApps_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApp/777"

	server.On("PUT", path, common.SuccessResponse(endpoint_resource.EndpointResource{
		ID:   777,
		Name: "Updated Custom App",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_custom_apps.Update(context.Background(), service, 777, &endpoint_resource.EndpointResource{
		ID:   777,
		Name: "Updated Custom App",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Custom App", result.Name)
}

func TestEndpointCustomApps_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := endPointApplicationsBasePath + "/customApp/777"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_custom_apps.Delete(context.Background(), service, 777)

	require.NoError(t, err)
}
