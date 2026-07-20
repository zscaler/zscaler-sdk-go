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
