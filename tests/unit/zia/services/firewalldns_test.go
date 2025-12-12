// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallipscontrolpolicies"
)

func TestFirewallDNSRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("FirewallDNSRules JSON marshaling", func(t *testing.T) {
		policy := firewalldnscontrolpolicies.FirewallDNSRules{
			ID:          12345,
			Name:        "Block Malicious DNS",
			Order:       1,
			State:       "ENABLED",
			Action:      "BLOCK",
			Description: "Block DNS queries to malicious domains",
			Rank:        7,
			Predefined:  false,
			DefaultRule: false,
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
	})

	t.Run("FirewallDNSRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Allow Internal DNS",
			"order": 2,
			"state": "ENABLED",
			"action": "ALLOW",
			"description": "Allow DNS to internal servers",
			"rank": 5,
			"predefined": false,
			"defaultRule": false,
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "IT"}
			],
			"destAddresses": ["10.0.0.53", "10.0.0.54"],
			"protocols": ["ANY_RULE"]
		}`

		var policy firewalldnscontrolpolicies.FirewallDNSRules
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, 54321, policy.ID)
		assert.Equal(t, "ALLOW", policy.Action)
		assert.Len(t, policy.Locations, 1)
	})
}

func TestFirewallIPSRules_Structure(t *testing.T) {
	t.Parallel()

	t.Run("FirewallIPSRules JSON marshaling", func(t *testing.T) {
		policy := firewallipscontrolpolicies.FirewallIPSRules{
			ID:          12345,
			Name:        "IPS Block High Risk",
			Order:       1,
			State:       "ENABLED",
			Action:      "BLOCK",
			Description: "Block high risk IPS signatures",
			Rank:        7,
			Predefined:  false,
			DefaultRule: false,
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"action":"BLOCK"`)
	})

	t.Run("FirewallIPSRules JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "IPS Monitor Medium Risk",
			"order": 2,
			"state": "ENABLED",
			"action": "MONITOR",
			"description": "Monitor medium risk signatures",
			"rank": 5,
			"predefined": false,
			"defaultRule": false,
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"srcIpGroups": [
				{"id": 200, "name": "Internal Networks"}
			],
			"destIpGroups": [
				{"id": 300, "name": "External Networks"}
			],
			"protocols": ["TCP_RULE", "UDP_RULE"]
		}`

		var policy firewallipscontrolpolicies.FirewallIPSRules
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, 54321, policy.ID)
		assert.Equal(t, "MONITOR", policy.Action)
	})
}

func TestFirewallDNSIPS_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse DNS control rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Policy 1", "action": "BLOCK", "state": "ENABLED"},
			{"id": 2, "name": "Policy 2", "action": "ALLOW", "state": "ENABLED"},
			{"id": 3, "name": "Default", "action": "ALLOW", "defaultRule": true}
		]`

		var policies []firewalldnscontrolpolicies.FirewallDNSRules
		err := json.Unmarshal([]byte(jsonResponse), &policies)
		require.NoError(t, err)

		assert.Len(t, policies, 3)
		assert.True(t, policies[2].DefaultRule)
	})

	t.Run("Parse IPS control rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "IPS Rule 1", "action": "BLOCK"},
			{"id": 2, "name": "IPS Rule 2", "action": "MONITOR"},
			{"id": 3, "name": "IPS Default", "action": "ALLOW", "defaultRule": true}
		]`

		var policies []firewallipscontrolpolicies.FirewallIPSRules
		err := json.Unmarshal([]byte(jsonResponse), &policies)
		require.NoError(t, err)

		assert.Len(t, policies, 3)
	})
}
