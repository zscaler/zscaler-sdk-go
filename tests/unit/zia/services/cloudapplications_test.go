// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudapplications/cloudapplications"
)

const (
	cloudAppPolicyPath    = "/zia/api/v1/cloudApplications/policy"
	cloudAppSSLPolicyPath = "/zia/api/v1/cloudApplications/sslPolicy"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestCloudApplications_GetCloudApplicationPolicy_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppPolicyPath, common.SuccessResponse([]cloudapplications.CloudApplications{
		{App: "GMAIL", AppName: "Gmail", Parent: "WEB_MAIL", ParentName: "Webmail"},
		{App: "OUTLOOK", AppName: "Outlook", Parent: "WEB_MAIL", ParentName: "Webmail"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass": []interface{}{"WEB_MAIL"},
	}

	result, err := cloudapplications.GetCloudApplicationPolicy(context.Background(), service, params)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Gmail", result[0].AppName)
	assert.Equal(t, "WEB_MAIL", result[0].Parent)
	assert.Equal(t, "Webmail", result[0].ParentName)
}

func TestCloudApplications_GetCloudApplicationPolicy_WithGroupResults_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", cloudAppPolicyPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Contains(t, r.URL.RawQuery, "appClass=WEB_MAIL")
		assert.Contains(t, r.URL.RawQuery, "groupResults=true")
		return common.SuccessResponse([]cloudapplications.CloudApplications{
			{App: "GMAIL", AppName: "Gmail", Parent: "WEB_MAIL", ParentName: "Webmail"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass":     []interface{}{"WEB_MAIL"},
		"groupResults": true,
	}

	result, err := cloudapplications.GetCloudApplicationPolicy(context.Background(), service, params)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCloudApplications_GetCloudApplicationPolicy_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppPolicyPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass": []interface{}{"WEB_MAIL"},
	}

	result, err := cloudapplications.GetCloudApplicationPolicy(context.Background(), service, params)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching cloud application policies")
}

func TestCloudApplications_GetCloudApplicationPolicy_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppPolicyPath, common.SuccessResponse([]cloudapplications.CloudApplications{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudapplications.GetCloudApplicationPolicy(context.Background(), service, map[string]interface{}{})

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestCloudApplications_GetCloudApplicationSSLPolicy_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppSSLPolicyPath, common.SuccessResponse([]cloudapplications.CloudApplications{
		{App: "FACEBOOK", AppName: "Facebook", Parent: "SOCIAL_NETWORKING", ParentName: "Social Networking"},
		{App: "LINKEDIN", AppName: "LinkedIn", Parent: "SOCIAL_NETWORKING", ParentName: "Social Networking"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass": []interface{}{"SOCIAL_NETWORKING"},
	}

	result, err := cloudapplications.GetCloudApplicationSSLPolicy(context.Background(), service, params)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Facebook", result[0].AppName)
	assert.Equal(t, "SOCIAL_NETWORKING", result[0].Parent)
}

func TestCloudApplications_GetCloudApplicationSSLPolicy_WithGroupResults_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", cloudAppSSLPolicyPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Contains(t, r.URL.RawQuery, "appClass=SOCIAL_NETWORKING")
		assert.Contains(t, r.URL.RawQuery, "groupResults=false")
		return common.SuccessResponse([]cloudapplications.CloudApplications{
			{App: "TWITTER", AppName: "Twitter", Parent: "SOCIAL_NETWORKING", ParentName: "Social Networking"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass":     []interface{}{"SOCIAL_NETWORKING"},
		"groupResults": false,
	}

	result, err := cloudapplications.GetCloudApplicationSSLPolicy(context.Background(), service, params)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestCloudApplications_GetCloudApplicationSSLPolicy_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", cloudAppSSLPolicyPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	params := map[string]interface{}{
		"appClass": []interface{}{"SOCIAL_NETWORKING"},
	}

	result, err := cloudapplications.GetCloudApplicationSSLPolicy(context.Background(), service, params)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "error fetching cloud application SSL policies")
}

// =====================================================
// Structure Tests
// =====================================================

func TestCloudApplications_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CloudApplications JSON marshaling", func(t *testing.T) {
		app := cloudapplications.CloudApplications{
			App:        "GMAIL",
			AppName:    "Gmail",
			Parent:     "WEB_MAIL",
			ParentName: "Webmail",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"app":"GMAIL"`)
		assert.Contains(t, string(data), `"appName":"Gmail"`)
		assert.Contains(t, string(data), `"parent":"WEB_MAIL"`)
		assert.Contains(t, string(data), `"parentName":"Webmail"`)
	})

	t.Run("CloudApplications JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"app": "FACEBOOK",
			"appName": "Facebook",
			"parent": "SOCIAL_NETWORKING",
			"parentName": "Social Networking"
		}`

		var app cloudapplications.CloudApplications
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, "FACEBOOK", app.App)
		assert.Equal(t, "Facebook", app.AppName)
		assert.Equal(t, "SOCIAL_NETWORKING", app.Parent)
		assert.Equal(t, "Social Networking", app.ParentName)
	})
}

func TestCloudApplications_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse cloud application policy list", func(t *testing.T) {
		jsonResponse := `[
			{"app": "GMAIL", "appName": "Gmail", "parent": "WEB_MAIL", "parentName": "Webmail"},
			{"app": "OUTLOOK", "appName": "Outlook", "parent": "WEB_MAIL", "parentName": "Webmail"}
		]`

		var apps []cloudapplications.CloudApplications
		err := json.Unmarshal([]byte(jsonResponse), &apps)
		require.NoError(t, err)

		assert.Len(t, apps, 2)
		assert.Equal(t, "Gmail", apps[0].AppName)
	})
}
