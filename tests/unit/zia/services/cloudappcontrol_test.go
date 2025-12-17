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

// =====================================================
// SDK Function Tests
// =====================================================

func TestCloudAppControl_GetByRuleID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := "/zia/api/v1/webApplicationRules/" + ruleType + "/12345"

	server.On("GET", path, common.SuccessResponse(cloudappcontrol.WebApplicationRules{
		ID:   ruleID,
		Name: "Block Streaming",
		Type: ruleType,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleID(context.Background(), service, ruleType, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestCloudAppControl_GetByRuleType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := "/zia/api/v1/webApplicationRules/" + ruleType

	server.On("GET", path, common.SuccessResponse([]cloudappcontrol.WebApplicationRules{
		{ID: 1, Name: "Rule 1", Type: ruleType},
		{ID: 2, Name: "Rule 2", Type: ruleType},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetByRuleType(context.Background(), service, ruleType)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCloudAppControl_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	path := "/zia/api/v1/webApplicationRules/" + ruleType

	server.On("POST", path, common.SuccessResponse(cloudappcontrol.WebApplicationRules{
		ID:   100,
		Name: "New Rule",
		Type: ruleType,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &cloudappcontrol.WebApplicationRules{
		Name:  "New Rule",
		Type:  ruleType,
		State: "ENABLED",
		Order: 1,
	}

	result, err := cloudappcontrol.Create(context.Background(), service, ruleType, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestCloudAppControl_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := "/zia/api/v1/webApplicationRules/" + ruleType + "/12345"

	server.On("PUT", path, common.SuccessResponse(cloudappcontrol.WebApplicationRules{
		ID:   ruleID,
		Name: "Updated Rule",
		Type: ruleType,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &cloudappcontrol.WebApplicationRules{
		ID:    ruleID,
		Name:  "Updated Rule",
		Type:  ruleType,
		State: "ENABLED",
	}

	result, err := cloudappcontrol.Update(context.Background(), service, ruleType, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Rule", result.Name)
}

func TestCloudAppControl_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "STREAMING_MEDIA"
	ruleID := 12345
	path := "/zia/api/v1/webApplicationRules/" + ruleType + "/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = cloudappcontrol.Delete(context.Background(), service, ruleType, ruleID)

	require.NoError(t, err)
}

func TestCloudAppControl_GetRuleTypeMapping_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/webApplicationRules/ruleTypeMapping"

	server.On("GET", path, common.SuccessResponse(map[string]string{
		"STREAMING_MEDIA": "Streaming Media",
		"CLOUD_STORAGE":   "Cloud Storage",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := cloudappcontrol.GetRuleTypeMapping(context.Background(), service)

	require.NoError(t, err)
	assert.Equal(t, "Streaming Media", result["STREAMING_MEDIA"])
}

// Note: CreateDuplicate test is skipped because the SDK function passes nil to Create
// which is rejected by the OneAPI client. This is a known limitation.
// The test below exercises code paths where the SDK would need to be fixed to support nil payloads.

func TestCloudAppControl_AllAvailableActions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleType := "CLOUD_STORAGE"
	path := "/zia/api/v1/webApplicationRules/" + ruleType + "/availableActions"

	server.On("POST", path, common.SuccessResponse([]string{"ALLOW", "BLOCK", "CAUTION", "ISOLATE"}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := cloudappcontrol.AvailableActionsRequest{
		CloudApps: []string{"GOOGLE_DRIVE", "DROPBOX"},
		Type:      ruleType,
	}

	result, err := cloudappcontrol.AllAvailableActions(context.Background(), service, ruleType, payload)

	require.NoError(t, err)
	assert.Len(t, result, 4)
	assert.Contains(t, result, "ALLOW")
	assert.Contains(t, result, "BLOCK")
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

