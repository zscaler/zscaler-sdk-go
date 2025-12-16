// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestURLFilteringRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/urlFilteringRules"

	// Get uses GetAll internally and finds by ID
	server.On("GET", path, common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: ruleID, Name: "Block Social Media", Action: "BLOCK", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := urlfilteringpolicies.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "Block Social Media", result.Name)
}

func TestURLFilteringRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/urlFilteringRules"

	server.On("GET", path, common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{
		{ID: 1, Name: "Rule 1", Action: "BLOCK"},
		{ID: 2, Name: "Rule 2", Action: "ALLOW"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := urlfilteringpolicies.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestURLFilteringRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Block Social Media"
	path := "/zia/api/v1/urlFilteringRules"

	server.On("GET", path, common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: 2, Name: ruleName, Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := urlfilteringpolicies.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestURLFilteringRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/urlFilteringRules"

	server.On("POST", path, common.SuccessResponse(urlfilteringpolicies.URLFilteringRule{
		ID:     99999,
		Name:   "New URL Rule",
		Action: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &urlfilteringpolicies.URLFilteringRule{
		Name:   "New URL Rule",
		Action: "BLOCK",
	}

	result, err := urlfilteringpolicies.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestURLFilteringRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/urlFilteringRules/12345"
	listPath := "/zia/api/v1/urlFilteringRules"

	// Mock GetAll call that Update performs to fetch CBIProfile
	server.On("GET", listPath, common.SuccessResponse([]urlfilteringpolicies.URLFilteringRule{
		{ID: ruleID, Name: "Updated URL Rule", CBIProfile: &urlfilteringpolicies.CBIProfile{ID: "cbi-id"}},
	}))

	server.On("PUT", path, common.SuccessResponse(urlfilteringpolicies.URLFilteringRule{
		ID:   ruleID,
		Name: "Updated URL Rule",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &urlfilteringpolicies.URLFilteringRule{
		ID:   ruleID,
		Name: "Updated URL Rule",
	}

	result, _, err := urlfilteringpolicies.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated URL Rule", result.Name)
}

func TestURLFilteringRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/urlFilteringRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = urlfilteringpolicies.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

func TestURLFilteringPolicies_Structure(t *testing.T) {
	t.Parallel()

	t.Run("URLFilteringRule JSON marshaling", func(t *testing.T) {
		rule := urlfilteringpolicies.URLFilteringRule{
			ID:                  12345,
			Name:                "Block Social Media",
			Order:               1,
			Protocols:           []string{"HTTPS_RULE", "HTTP_RULE"},
			URLCategories:       []string{"SOCIAL_NETWORKING", "STREAMING_MEDIA"},
			State:               "ENABLED",
			Action:              "BLOCK",
			Rank:                7,
			BlockOverride:       true,
			TimeQuota:           60,
			SizeQuota:           1024,
			EnforceTimeValidity: true,
			Ciparule:            false,
			DeviceTrustLevels:   []string{"HIGH_TRUST", "MEDIUM_TRUST"},
			UserRiskScoreLevels: []string{"HIGH", "CRITICAL"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
		assert.Contains(t, string(data), `"blockOverride":true`)
	})

	t.Run("URLFilteringRule JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Allow Business Apps",
			"order": 2,
			"protocols": ["HTTPS_RULE"],
			"urlCategories": ["BUSINESS", "FINANCE"],
			"state": "ENABLED",
			"action": "ALLOW",
			"rank": 5,
			"requestMethods": ["GET", "POST"],
			"sourceCountries": ["US", "CA"],
			"endUserNotificationUrl": "https://notify.company.com",
			"validityStartTime": 1699000000,
			"validityEndTime": 1699999999,
			"validityTimeZoneId": "America/Los_Angeles",
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "IT"}
			],
			"groups": [
				{"id": 300, "name": "Admins"}
			],
			"users": [
				{"id": 400, "name": "john.doe@company.com"}
			],
			"cbiProfile": {
				"id": "cbi-profile-uuid",
				"name": "Default CBI Profile",
				"url": "https://cbi.zscaler.com",
				"profileSeq": 1
			}
		}`

		var rule urlfilteringpolicies.URLFilteringRule
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "ALLOW", rule.Action)
		assert.NotNil(t, rule.CBIProfile)
		assert.Equal(t, "cbi-profile-uuid", rule.CBIProfile.ID)
		assert.Len(t, rule.Locations, 1)
	})

	t.Run("CBIProfile JSON marshaling", func(t *testing.T) {
		profile := urlfilteringpolicies.CBIProfile{
			ID:         "cbi-uuid-12345",
			Name:       "Isolation Profile",
			URL:        "https://isolation.zscaler.com",
			ProfileSeq: 1,
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"cbi-uuid-12345"`)
		assert.Contains(t, string(data), `"profileSeq":1`)
	})

	t.Run("URLAdvancedPolicySettings JSON marshaling", func(t *testing.T) {
		settings := urlfilteringpolicies.URLAdvancedPolicySettings{
			EnableDynamicContentCat:           true,
			ConsiderEmbeddedSites:             true,
			EnforceSafeSearch:                 true,
			EnableOffice365:                   true,
			EnableUcaasZoom:                   true,
			EnableChatGptPrompt:               true,
			EnableMicrosoftCoPilotPrompt:      true,
			EnableNewlyRegisteredDomains:      true,
			BlockSkype:                        false,
			EnableBlockOverrideForNonAuthUser: true,
			SafeSearchApps:                    []string{"GOOGLE", "BING", "YOUTUBE"},
		}

		data, err := json.Marshal(settings)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"enableDynamicContentCat":true`)
		assert.Contains(t, string(data), `"enforceSafeSearch":true`)
	})
}

func TestURLFilteringPolicies_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse URL filtering rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "action": "BLOCK", "state": "ENABLED"},
			{"id": 2, "name": "Rule 2", "action": "ALLOW", "state": "ENABLED"},
			{"id": 3, "name": "Rule 3", "action": "CAUTION", "state": "DISABLED"}
		]`

		var rules []urlfilteringpolicies.URLFilteringRule
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "CAUTION", rules[2].Action)
	})
}

