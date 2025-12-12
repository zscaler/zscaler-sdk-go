// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_dns_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policy_management/traffic_log_rules"
)

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

