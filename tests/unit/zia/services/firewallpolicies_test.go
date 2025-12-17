// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestFirewallFilteringRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallFilteringRules/12345"

	server.On("GET", path, common.SuccessResponse(filteringrules.FirewallFilteringRules{
		ID:          ruleID,
		Name:        "Block Malicious IPs",
		Description: "Block known malicious IP addresses",
		Action:      "BLOCK",
		State:       "ENABLED",
		Order:       1,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filteringrules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "Block Malicious IPs", result.Name)
	assert.Equal(t, "BLOCK", result.Action)
}

func TestFirewallFilteringRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallFilteringRules"

	server.On("GET", path, common.SuccessResponse([]filteringrules.FirewallFilteringRules{
		{ID: 1, Name: "Rule 1", Action: "ALLOW", State: "ENABLED"},
		{ID: 2, Name: "Rule 2", Action: "BLOCK", State: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filteringrules.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestFirewallFilteringRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallFilteringRules"

	server.On("POST", path, common.SuccessResponse(filteringrules.FirewallFilteringRules{
		ID:     99999,
		Name:   "New Rule",
		Action: "ALLOW",
		State:  "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := filteringrules.FirewallFilteringRules{
		Name:   "New Rule",
		Action: "ALLOW",
		State:  "ENABLED",
	}

	result, err := filteringrules.Create(context.Background(), service, &newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestFirewallFilteringRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallFilteringRules/12345"

	server.On("PUT", path, common.SuccessResponse(filteringrules.FirewallFilteringRules{
		ID:     ruleID,
		Name:   "Updated Rule",
		Action: "BLOCK",
		State:  "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := filteringrules.FirewallFilteringRules{
		ID:     ruleID,
		Name:   "Updated Rule",
		Action: "BLOCK",
	}

	result, err := filteringrules.Update(context.Background(), service, ruleID, &updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Rule", result.Name)
}

func TestFirewallFilteringRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallFilteringRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = filteringrules.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

func TestFirewallFilteringRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Block Malicious IPs"
	path := "/zia/api/v1/firewallFilteringRules"

	server.On("GET", path, common.SuccessResponse([]filteringrules.FirewallFilteringRules{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: 2, Name: ruleName, Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filteringrules.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestFirewallFilteringRules_GetFilteringRuleCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallFilteringRules/count"

	server.On("GET", path, common.SuccessResponse(15))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := filteringrules.GetFirewallFilteringRuleCount(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Equal(t, 15, result)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestFirewallFilteringRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("FirewallFilteringRules JSON marshaling", func(t *testing.T) {
		rule := filteringrules.FirewallFilteringRules{
			ID:                12345,
			Name:              "Block Malicious IPs",
			Order:             1,
			Rank:              7,
			AccessControl:     "READ_WRITE",
			EnableFullLogging: true,
			Action:            "BLOCK",
			State:             "ENABLED",
			Description:       "Block known malicious IP addresses",
			SrcIps:            []string{"10.0.0.0/8"},
			DestAddresses:     []string{"192.168.1.0/24"},
			DestCountries:     []string{"US", "CA"},
			NwApplications:    []string{"HTTP", "HTTPS"},
			DefaultRule:       false,
			Predefined:        false,
			DeviceTrustLevels: []string{"TRUSTED", "UNKNOWN"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
		assert.Contains(t, string(data), `"state":"ENABLED"`)
		assert.Contains(t, string(data), `"enableFullLogging":true`)
	})

	t.Run("FirewallFilteringRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Allow Office Traffic",
			"order": 2,
			"rank": 5,
			"action": "ALLOW",
			"state": "ENABLED",
			"description": "Allow traffic from office locations",
			"srcIps": ["172.16.0.0/12"],
			"destAddresses": ["8.8.8.8", "8.8.4.4"],
			"destIpCategories": ["DNS"],
			"destCountries": ["US"],
			"sourceCountries": ["US"],
			"excludeSrcCountries": false,
			"nwApplications": ["DNS"],
			"defaultRule": false,
			"predefined": false,
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "Engineering"}
			],
			"groups": [
				{"id": 300, "name": "Developers"}
			],
			"nwServices": [
				{"id": 400, "name": "DNS"}
			],
			"lastModifiedTime": 1699000000
		}`

		var rule filteringrules.FirewallFilteringRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "ALLOW", rule.Action)
		assert.Len(t, rule.SrcIps, 1)
		assert.Len(t, rule.Locations, 1)
		assert.Len(t, rule.Departments, 1)
		assert.Len(t, rule.NwServices, 1)
	})

	t.Run("GetAllFilterOptions structure", func(t *testing.T) {
		opts := filteringrules.GetAllFilterOptions{
			PredefinedRuleCount: true,
			RuleName:            "test-rule",
			RuleLabel:           "test-label",
			RuleLabelId:         123,
			RuleOrder:           "1",
			RuleDescription:     "test description",
			RuleAction:          "ALLOW",
			Location:            "HQ",
			Department:          "Engineering",
			Group:               "Developers",
			User:                "user@company.com",
			Device:              "device-123",
			DeviceGroup:         "Mobile Devices",
			DeviceTrustLevel:    "TRUSTED",
			SrcIps:              "10.0.0.0/8",
			DestAddresses:       "192.168.1.0/24",
			SrcIpGroups:         "Internal IPs",
			DestIpGroups:        "External IPs",
			NwApplication:       "HTTP",
			NwServices:          "Web",
			DestIpCategories:    "CDN",
		}

		assert.True(t, opts.PredefinedRuleCount)
		assert.Equal(t, "test-rule", opts.RuleName)
		assert.Equal(t, 123, opts.RuleLabelId)
	})
}

func TestFirewallFilteringRules_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse firewall rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "action": "ALLOW", "state": "ENABLED", "order": 1},
			{"id": 2, "name": "Rule 2", "action": "BLOCK", "state": "ENABLED", "order": 2},
			{"id": 3, "name": "Rule 3", "action": "ALLOW", "state": "DISABLED", "order": 3}
		]`

		var rules []filteringrules.FirewallFilteringRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "BLOCK", rules[1].Action)
		assert.Equal(t, "DISABLED", rules[2].State)
	})

	t.Run("Parse rule with workload groups", func(t *testing.T) {
		jsonResponse := `{
			"id": 100,
			"name": "Workload Rule",
			"action": "ALLOW",
			"workloadGroups": [
				{"id": 1, "name": "Production"},
				{"id": 2, "name": "Development"}
			],
			"zpaAppSegments": [
				{"id": 10, "name": "ZPA Segment 1", "externalId": "zpa-1"}
			]
		}`

		var rule filteringrules.FirewallFilteringRules
		err := json.Unmarshal([]byte(jsonResponse), &rule)
		require.NoError(t, err)

		assert.Len(t, rule.WorkloadGroups, 2)
		assert.Len(t, rule.ZPAAppSegments, 1)
	})
}

