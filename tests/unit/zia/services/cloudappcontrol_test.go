// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudappcontrol"
)

const webApplicationRulesBase = "/zia/api/v1/webApplicationRules"

// sampleStreamingMediaRule mirrors the integration test payload in cloudappcontrol_test.go.
func sampleStreamingMediaRule(name string) cloudappcontrol.WebApplicationRules {
	return cloudappcontrol.WebApplicationRules{
		Name:         name,
		Description:  name,
		Order:        1,
		Rank:         7,
		State:        "ENABLED",
		Type:         "STREAMING_MEDIA",
		Applications: []string{"YOUTUBE", "GOOGLE_STREAMING"},
		Actions:      []string{"ALLOW_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_UPLOAD"},
	}
}

// dropboxAvailableActions mirrors the integration TestAllAvailableActions expected response.
var dropboxAvailableActions = []string{
	"ALLOW_FILE_SHARE_CREATE",
	"ALLOW_FILE_SHARE_DELETE",
	"ALLOW_FILE_SHARE_DOWNLOAD",
	"ALLOW_FILE_SHARE_EDIT",
	"ALLOW_FILE_SHARE_INVITE",
	"ALLOW_FILE_SHARE_RENAME",
	"ALLOW_FILE_SHARE_SHARE",
	"DENY_FILE_SHARE_CREATE",
	"DENY_FILE_SHARE_DELETE",
	"DENY_FILE_SHARE_DOWNLOAD",
	"DENY_FILE_SHARE_EDIT",
	"DENY_FILE_SHARE_INVITE",
	"DENY_FILE_SHARE_RENAME",
	"DENY_FILE_SHARE_SHARE",
	"FILE_SHARE_CONDITIONAL_ACCESS",
}

// =====================================================
// SDK Function Tests
// =====================================================

func TestCloudAppControl_GetByRuleID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := webApplicationRulesBase + "/" + ruleType + "/12345"

	rule := sampleStreamingMediaRule("tests-streaming-rule")
	rule.ID = ruleID

	server.On("GET", path, common.SuccessResponse(rule))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleID(context.Background(), service, ruleType, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "tests-streaming-rule", result.Name)
	assert.Equal(t, "ENABLED", result.State)
	assert.Equal(t, []string{"YOUTUBE", "GOOGLE_STREAMING"}, result.Applications)
	assert.Equal(t, []string{"ALLOW_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_UPLOAD"}, result.Actions)
}

func TestCloudAppControl_GetByRuleID_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA/9999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleID(context.Background(), service, "STREAMING_MEDIA", 9999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_GetByRuleType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := webApplicationRulesBase + "/" + ruleType

	rule := sampleStreamingMediaRule("tests-streaming-rule")
	rule.ID = 1

	server.On("GET", path, common.SuccessResponse([]cloudappcontrol.WebApplicationRules{
		rule,
		func() cloudappcontrol.WebApplicationRules {
			r := sampleStreamingMediaRule("tests-streaming-rule-2")
			r.ID = 2
			return r
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleType(context.Background(), service, ruleType)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "STREAMING_MEDIA", result[0].Type)
}

func TestCloudAppControl_GetByRuleType_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleType(context.Background(), service, "STREAMING_MEDIA")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := webApplicationRulesBase + "/" + ruleType

	created := sampleStreamingMediaRule("tests-new-rule")
	created.ID = 100

	server.On("POST", path, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := sampleStreamingMediaRule("tests-new-rule")

	result, err := cloudappcontrol.Create(context.Background(), service, ruleType, &newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
	assert.Equal(t, "tests-new-rule", result.Name)
	assert.Equal(t, []string{"YOUTUBE", "GOOGLE_STREAMING"}, result.Applications)
}

func TestCloudAppControl_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA"
	server.On("POST", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := sampleStreamingMediaRule("tests-fail-rule")
	result, err := cloudappcontrol.Create(context.Background(), service, "STREAMING_MEDIA", &newRule)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_Create_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA"
	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := sampleStreamingMediaRule("tests-no-body")
	result, err := cloudappcontrol.Create(context.Background(), service, "STREAMING_MEDIA", &newRule)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "object returned from api was not a rule Pointer")
}

func TestCloudAppControl_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := webApplicationRulesBase + "/" + ruleType + "/12345"

	updated := sampleStreamingMediaRule("tests-updated-rule")
	updated.ID = ruleID

	server.On("PUT", path, common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := sampleStreamingMediaRule("tests-updated-rule")
	updateRule.ID = ruleID

	result, err := cloudappcontrol.Update(context.Background(), service, ruleType, ruleID, &updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tests-updated-rule", result.Name)
	assert.Equal(t, 7, result.Rank)
}

func TestCloudAppControl_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA/12345"
	server.On("PUT", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := sampleStreamingMediaRule("tests-updated-rule")
	result, err := cloudappcontrol.Update(context.Background(), service, "STREAMING_MEDIA", 12345, &updateRule)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := webApplicationRulesBase + "/" + ruleType + "/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = cloudappcontrol.Delete(context.Background(), service, ruleType, ruleID)

	require.NoError(t, err)
}

func TestCloudAppControl_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/STREAMING_MEDIA/12345"
	server.On("DELETE", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = cloudappcontrol.Delete(context.Background(), service, "STREAMING_MEDIA", 12345)

	require.Error(t, err)
}

func TestCloudAppControl_CreateDuplicate_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// CreateDuplicate passes nil to Client.Create, which the OneAPI client
	// rejects before any HTTP call is made.
	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.CreateDuplicate(context.Background(), service, "STREAMING_MEDIA", 12345, "tests-duplicate")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_GetRuleTypeMapping_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/ruleTypeMapping"

	server.On("GET", path, common.SuccessResponse(map[string]string{
		"Webmail":               "WEBMAIL",
		"Social Networking":     "SOCIAL_NETWORKING",
		"Finance":               "FINANCE",
		"Legal":                 "LEGAL",
		"AI & ML Applications":  "AI_ML",
		"Streaming Media":       "STREAMING_MEDIA",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetRuleTypeMapping(context.Background(), service)

	require.NoError(t, err)
	assert.Equal(t, "WEBMAIL", result["Webmail"])
	assert.Equal(t, "SOCIAL_NETWORKING", result["Social Networking"])
	assert.Equal(t, "STREAMING_MEDIA", result["Streaming Media"])
}

func TestCloudAppControl_GetRuleTypeMapping_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := webApplicationRulesBase + "/ruleTypeMapping"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetRuleTypeMapping(context.Background(), service)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_AllAvailableActions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := webApplicationRulesBase + "/" + ruleType + "/availableActions"

	server.On("POST", path, common.SuccessResponse(dropboxAvailableActions))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := cloudappcontrol.AvailableActionsRequest{
		CloudApps: []string{"DROPBOX"},
		Type:      "ANY",
	}

	result, err := cloudappcontrol.AllAvailableActions(context.Background(), service, ruleType, payload)

	require.NoError(t, err)
	assert.Equal(t, dropboxAvailableActions, result)
}

func TestCloudAppControl_AllAvailableActions_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := webApplicationRulesBase + "/" + ruleType + "/availableActions"
	server.On("POST", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := cloudappcontrol.AvailableActionsRequest{
		CloudApps: []string{"DROPBOX"},
		Type:      "ANY",
	}

	result, err := cloudappcontrol.AllAvailableActions(context.Background(), service, ruleType, payload)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestCloudAppControl_AllAvailableActions_InvalidResponse_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := webApplicationRulesBase + "/" + ruleType + "/availableActions"
	server.On("POST", path, common.SuccessResponse(`not-a-json-array`))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := cloudappcontrol.AvailableActionsRequest{
		CloudApps: []string{"DROPBOX"},
		Type:      "ANY",
	}

	result, err := cloudappcontrol.AllAvailableActions(context.Background(), service, ruleType, payload)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestCloudAppControl_Structure(t *testing.T) {
	t.Parallel()

	t.Run("WebApplicationRules JSON marshaling", func(t *testing.T) {
		rule := cloudappcontrol.WebApplicationRules{
			ID:               12345,
			Name:             "Block Cloud Storage",
			Description:      "Block cloud storage applications",
			Actions:          []string{"BLOCK"},
			State:            "ENABLED",
			Rank:             7,
			Type:             "CLOUD_STORAGE",
			Order:            1,
			TimeQuota:        60,
			SizeQuota:        1024,
			CascadingEnabled: true,
			Predefined:       false,
			Applications:     []string{"GOOGLE_DRIVE", "DROPBOX", "BOX"},
			DeviceTrustLevels: []string{"HIGH_TRUST", "MEDIUM_TRUST"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"CLOUD_STORAGE"`)
		assert.Contains(t, string(data), `"cascadingEnabled":true`)
	})

	t.Run("WebApplicationRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Allow Collaboration Apps",
			"description": "Allow collaboration applications",
			"actions": ["ALLOW"],
			"state": "ENABLED",
			"type": "COLLABORATION",
			"order": 2,
			"applications": ["SLACK", "TEAMS", "ZOOM"],
			"numberOfApplications": 3,
			"eunEnabled": true,
			"eunTemplateId": 100,
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "Engineering"}
			],
			"cbiProfile": {
				"id": "cbi-uuid",
				"name": "Default Profile",
				"url": "https://cbi.zscaler.com",
				"profileSeq": 1,
				"defaultProfile": true,
				"sandboxMode": false
			}
		}`

		var rule cloudappcontrol.WebApplicationRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, 3, rule.NumberOfApplications)
		assert.True(t, rule.EunEnabled)
		assert.Equal(t, "cbi-uuid", rule.CBIProfile.ID)
	})

	t.Run("CBIProfile JSON marshaling", func(t *testing.T) {
		profile := cloudappcontrol.CBIProfile{
			ID:             "cbi-profile-uuid",
			Name:           "Isolation Profile",
			URL:            "https://isolation.zscaler.com",
			ProfileSeq:     1,
			DefaultProfile: true,
			SandboxMode:    false,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"cbi-profile-uuid"`)
		assert.Contains(t, string(data), `"defaultProfile":true`)
	})

	t.Run("CloudAppInstances JSON marshaling", func(t *testing.T) {
		instance := cloudappcontrol.CloudAppInstances{
			ID:   12345,
			Name: "Corporate Google Workspace",
			Type: "GOOGLE_DRIVE",
		}

		data, err := json.Marshal(instance)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"GOOGLE_DRIVE"`)
	})

	t.Run("CloudApp JSON marshaling", func(t *testing.T) {
		app := cloudappcontrol.CloudApp{
			Val:                 100,
			WebApplicationClass: "FILE_SHARING",
			BackendName:         "google_drive",
			OriginalName:        "Google Drive",
			Name:                "Google Drive",
			Deprecated:          false,
			Misc:                false,
			AppNotReady:         false,
			UnderMigration:      false,
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"val":100`)
		assert.Contains(t, string(data), `"webApplicationClass":"FILE_SHARING"`)
	})

	t.Run("AvailableActionsRequest JSON marshaling", func(t *testing.T) {
		req := cloudappcontrol.AvailableActionsRequest{
			CloudApps: []string{"GOOGLE_DRIVE", "DROPBOX"},
			Type:      "CLOUD_STORAGE",
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"cloudApps":["GOOGLE_DRIVE","DROPBOX"]`)
		assert.Contains(t, string(data), `"type":"CLOUD_STORAGE"`)
	})
}

func TestCloudAppControl_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse web application rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "type": "CLOUD_STORAGE", "state": "ENABLED"},
			{"id": 2, "name": "Rule 2", "type": "COLLABORATION", "state": "ENABLED"},
			{"id": 3, "name": "Rule 3", "type": "STREAMING_MEDIA", "state": "DISABLED"}
		]`

		var rules []cloudappcontrol.WebApplicationRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "STREAMING_MEDIA", rules[2].Type)
	})

	t.Run("Parse rule type mapping", func(t *testing.T) {
		jsonResponse := `{
			"CLOUD_STORAGE": "Cloud Storage",
			"COLLABORATION": "Collaboration",
			"STREAMING_MEDIA": "Streaming Media"
		}`

		var mapping map[string]string
		err := json.Unmarshal([]byte(jsonResponse), &mapping)
		require.NoError(t, err)

		assert.Equal(t, "Cloud Storage", mapping["CLOUD_STORAGE"])
	})
}

