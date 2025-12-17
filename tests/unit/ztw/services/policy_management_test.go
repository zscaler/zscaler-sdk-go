// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_dns_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_log_rules"
)

// =====================================================
// Forwarding Rules SDK Function Tests
// =====================================================

func TestForwardingRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/ztw/api/v1/ecRules/ecRdr/12345"

	server.On("GET", path, common.SuccessResponse(forwarding_rules.ForwardingRules{
		ID:            ruleID,
		Name:          "Forward-to-ZIA",
		ForwardMethod: "ZIA",
		State:         "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "ZIA", result.ForwardMethod)
}

func TestForwardingRules_GetRulesByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Forward-to-ZIA"
	path := "/ztw/api/v1/ecRules/ecRdr"

	server.On("GET", path, common.SuccessResponse([]forwarding_rules.ForwardingRules{
		{ID: 1, Name: "Other Rule", ForwardMethod: "DIRECT"},
		{ID: 2, Name: ruleName, ForwardMethod: "ZIA"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.GetRulesByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestForwardingRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecRules/ecRdr"

	server.On("GET", path, common.SuccessResponse([]forwarding_rules.ForwardingRules{
		{ID: 1, Name: "Rule 1", ForwardMethod: "ZIA"},
		{ID: 2, Name: "Rule 2", ForwardMethod: "DIRECT"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestForwardingRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecRules/ecRdr"

	server.On("POST", path, common.SuccessResponse(forwarding_rules.ForwardingRules{
		ID:            99999,
		Name:          "New Rule",
		ForwardMethod: "ZIA",
		State:         "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &forwarding_rules.ForwardingRules{
		Name:          "New Rule",
		ForwardMethod: "ZIA",
	}

	result, err := forwarding_rules.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestForwardingRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/ztw/api/v1/ecRules/ecRdr/12345"

	server.On("PUT", path, common.SuccessResponse(forwarding_rules.ForwardingRules{
		ID:            ruleID,
		Name:          "Updated Rule",
		ForwardMethod: "DIRECT",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &forwarding_rules.ForwardingRules{
		ID:            ruleID,
		Name:          "Updated Rule",
		ForwardMethod: "DIRECT",
	}

	result, err := forwarding_rules.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Rule", result.Name)
}

func TestForwardingRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/ztw/api/v1/ecRules/ecRdr/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = forwarding_rules.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

// =====================================================
// Traffic DNS Rules SDK Function Tests
// =====================================================

func TestTrafficDNSRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/ztw/api/v1/ecRules/ecDns/12345"

	server.On("GET", path, common.SuccessResponse(traffic_dns_rules.ECDNSRules{
		ID:     ruleID,
		Name:   "DNS Rule",
		Action: "ALLOW",
		State:  "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_dns_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestTrafficDNSRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecRules/ecDns"

	server.On("GET", path, common.SuccessResponse([]traffic_dns_rules.ECDNSRules{
		{ID: 1, Name: "DNS Rule 1", Action: "ALLOW"},
		{ID: 2, Name: "DNS Rule 2", Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_dns_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Traffic Log Rules SDK Function Tests
// =====================================================

func TestTrafficLogRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/ztw/api/v1/ecRules/self/12345"

	server.On("GET", path, common.SuccessResponse(traffic_log_rules.ECTrafficLogRules{
		ID:    ruleID,
		Name:  "Log Rule",
		State: "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_log_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestTrafficLogRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ecRules/self"

	server.On("GET", path, common.SuccessResponse([]traffic_log_rules.ECTrafficLogRules{
		{ID: 1, Name: "Log Rule 1", State: "ENABLED"},
		{ID: 2, Name: "Log Rule 2", State: "DISABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := traffic_log_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests
// =====================================================

func TestForwardingRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingRules JSON marshaling", func(t *testing.T) {
		rule := forwarding_rules.ForwardingRules{
			ID:            12345,
			Name:          "Forward-to-ZIA",
			Description:   "Forward traffic to ZIA",
			Type:          "EC_RDR",
			Order:         1,
			Rank:          7,
			ForwardMethod: "ZIA",
			State:         "ENABLED",
			DefaultRule:   false,
			SrcIps:        []string{"10.0.0.0/8"},
			DestAddresses: []string{"192.168.1.0/24"},
			DestCountries: []string{"US", "CA"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"forwardMethod":"ZIA"`)
		assert.Contains(t, string(data), `"state":"ENABLED"`)
	})

	t.Run("ForwardingRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Direct-Traffic",
			"description": "Direct traffic rule",
			"type": "FORWARDING",
			"order": 2,
			"rank": 5,
			"forwardMethod": "DIRECT",
			"state": "ENABLED",
			"defaultRule": false,
			"srcIps": ["172.16.0.0/12"],
			"destAddresses": ["8.8.8.8", "8.8.4.4"],
			"nwApplications": ["DNS", "HTTP"],
			"locations": [
				{"id": 100, "name": "Location-1"}
			],
			"ecGroups": [
				{"id": 200, "name": "EC-Group-1"}
			]
		}`

		var rule forwarding_rules.ForwardingRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "DIRECT", rule.ForwardMethod)
		assert.Len(t, rule.SrcIps, 1)
		assert.Len(t, rule.NwApplications, 2)
		assert.Len(t, rule.Locations, 1)
	})

	t.Run("ForwardingRulesCountQuery structure", func(t *testing.T) {
		query := forwarding_rules.ForwardingRulesCountQuery{
			PredefinedRuleCount: true,
			RuleName:            "test-rule",
			RuleOrder:           "1",
			RuleDescription:     "test description",
			RuleForwardMethod:   "ZIA",
			Location:            "US",
		}

		assert.True(t, query.PredefinedRuleCount)
		assert.Equal(t, "test-rule", query.RuleName)
	})
}

func TestTrafficDNSRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ECDNSRules JSON marshaling", func(t *testing.T) {
		rule := traffic_dns_rules.ECDNSRules{
			ID:          12345,
			Name:        "DNS-Rule-1",
			Description: "DNS forwarding rule",
			Type:        "EC_DNS",
			Action:      "ALLOW",
			Order:       1,
			Rank:        7,
			State:       "ENABLED",
			Predefined:  false,
			DefaultRule: false,
			SrcIps:      []string{"10.0.0.0/8"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"ALLOW"`)
		assert.Contains(t, string(data), `"type":"EC_DNS"`)
	})

	t.Run("ECDNSRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "DNS-Rule-2",
			"description": "Block DNS rule",
			"type": "EC_DNS",
			"action": "BLOCK",
			"order": 2,
			"state": "ENABLED",
			"destAddresses": ["1.1.1.1", "1.0.0.1"],
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"srcIpGroups": [
				{"id": 200, "name": "Internal-IPs"}
			]
		}`

		var rule traffic_dns_rules.ECDNSRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "BLOCK", rule.Action)
		assert.Len(t, rule.DestAddresses, 2)
		assert.Len(t, rule.Locations, 1)
	})
}

func TestTrafficLogRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ECTrafficLogRules JSON marshaling", func(t *testing.T) {
		rule := traffic_log_rules.ECTrafficLogRules{
			ID:            12345,
			Name:          "Log-Rule-1",
			Description:   "Traffic logging rule",
			Order:         1,
			Rank:          7,
			State:         "ENABLED",
			Type:          "EC_SELF",
			ForwardMethod: "ECSELF",
			DefaultRule:   false,
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"EC_SELF"`)
		assert.Contains(t, string(data), `"forwardMethod":"ECSELF"`)
	})

	t.Run("ECTrafficLogRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Log-Rule-2",
			"description": "Advanced logging rule",
			"order": 2,
			"state": "DISABLED",
			"type": "EC_SELF",
			"locations": [
				{"id": 100, "name": "Branch-1"},
				{"id": 101, "name": "Branch-2"}
			],
			"ecGroups": [
				{"id": 200, "name": "EC-Group-1"}
			],
			"lastModifiedTime": 1699000000
		}`

		var rule traffic_log_rules.ECTrafficLogRules
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "DISABLED", rule.State)
		assert.Len(t, rule.Locations, 2)
		assert.Len(t, rule.ECGroups, 1)
	})
}

func TestPolicyManagement_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse forwarding rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule-1", "forwardMethod": "ZIA", "state": "ENABLED"},
			{"id": 2, "name": "Rule-2", "forwardMethod": "DIRECT", "state": "ENABLED"},
			{"id": 3, "name": "Rule-3", "forwardMethod": "DROP", "state": "DISABLED"}
		]`

		var rules []forwarding_rules.ForwardingRules
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "ZIA", rules[0].ForwardMethod)
		assert.Equal(t, "DISABLED", rules[2].State)
	})

	t.Run("Parse count response", func(t *testing.T) {
		jsonResponse := `{"count": 42}`

		var response forwarding_rules.ForwardingRulesCountResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 42, response.Count)
	})
}

