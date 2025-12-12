// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/nat_control_policies"
)

func TestNatControlPolicies_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NatControlPolicies JSON marshaling", func(t *testing.T) {
		policy := nat_control_policies.NatControlPolicies{
			ID:                  12345,
			Name:                "NAT Policy 1",
			Order:               1,
			Rank:                7,
			Description:         "Redirect traffic to specific server",
			State:               "ENABLED",
			RedirectFqdn:        "proxy.company.com",
			RedirectIp:          "192.168.1.100",
			RedirectPort:        8080,
			TrustedResolverRule: true,
			EnableFullLogging:   true,
			Predefined:          false,
			DefaultRule:         false,
			DestAddresses:       []string{"10.0.0.0/8"},
			SrcIps:              []string{"192.168.0.0/16"},
			DestCountries:       []string{"US", "CA"},
		}

		data, err := json.Marshal(policy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"redirectFqdn":"proxy.company.com"`)
		assert.Contains(t, string(data), `"redirectPort":8080`)
	})

	t.Run("NatControlPolicies JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "NAT Policy 2",
			"order": 2,
			"state": "ENABLED",
			"description": "DNAT for specific destinations",
			"redirectIp": "10.0.0.100",
			"redirectPort": 443,
			"trustedResolverRule": false,
			"enableFullLogging": true,
			"predefined": false,
			"defaultRule": false,
			"destAddresses": ["172.16.0.0/12"],
			"destCountries": ["US"],
			"destIpCategories": ["DATACENTER"],
			"resCategories": ["DNS"],
			"locations": [
				{"id": 100, "name": "HQ"}
			],
			"departments": [
				{"id": 200, "name": "IT"}
			],
			"srcIpGroups": [
				{"id": 300, "name": "Internal Networks"}
			],
			"destIpGroups": [
				{"id": 400, "name": "External Servers"}
			],
			"nwServices": [
				{"id": 500, "name": "HTTPS"}
			]
		}`

		var policy nat_control_policies.NatControlPolicies
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.Equal(t, 54321, policy.ID)
		assert.Equal(t, "10.0.0.100", policy.RedirectIp)
		assert.Len(t, policy.Locations, 1)
		assert.Len(t, policy.SrcIpGroups, 1)
	})

	t.Run("NatControlPolicies default rule", func(t *testing.T) {
		jsonData := `{
			"id": 1,
			"name": "Default NAT Rule",
			"order": 100,
			"state": "ENABLED",
			"defaultRule": true,
			"predefined": true
		}`

		var policy nat_control_policies.NatControlPolicies
		err := json.Unmarshal([]byte(jsonData), &policy)
		require.NoError(t, err)

		assert.True(t, policy.DefaultRule)
		assert.True(t, policy.Predefined)
	})
}

func TestNatControlPolicies_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse NAT control policies list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Policy 1", "state": "ENABLED", "order": 1},
			{"id": 2, "name": "Policy 2", "state": "ENABLED", "order": 2},
			{"id": 3, "name": "Default", "state": "ENABLED", "order": 100, "defaultRule": true}
		]`

		var policies []nat_control_policies.NatControlPolicies
		err := json.Unmarshal([]byte(jsonResponse), &policies)
		require.NoError(t, err)

		assert.Len(t, policies, 3)
		assert.True(t, policies[2].DefaultRule)
	})
}

