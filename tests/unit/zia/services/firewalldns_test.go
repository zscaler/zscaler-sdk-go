// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewalldnscontrolpolicies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallipscontrolpolicies"
)

// =====================================================
// Firewall DNS Control Policies - SDK Function Tests
// =====================================================

func TestFirewallDNS_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallDnsRules/12345"

	server.On("GET", path, common.SuccessResponse(firewalldnscontrolpolicies.FirewallDNSRules{
		ID:     ruleID,
		Name:   "Block Malicious DNS",
		Action: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewalldnscontrolpolicies.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestFirewallDNS_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "Block Malicious DNS"
	path := "/zia/api/v1/firewallDnsRules"

	server.On("GET", path, common.SuccessResponse([]firewalldnscontrolpolicies.FirewallDNSRules{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: 2, Name: ruleName, Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewalldnscontrolpolicies.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestFirewallDNS_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallDnsRules"

	server.On("POST", path, common.SuccessResponse(firewalldnscontrolpolicies.FirewallDNSRules{
		ID:     100,
		Name:   "New DNS Rule",
		Action: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &firewalldnscontrolpolicies.FirewallDNSRules{
		Name:   "New DNS Rule",
		Action: "BLOCK",
		State:  "ENABLED",
		Order:  1,
	}

	result, err := firewalldnscontrolpolicies.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestFirewallDNS_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallDnsRules/12345"

	server.On("PUT", path, common.SuccessResponse(firewalldnscontrolpolicies.FirewallDNSRules{
		ID:     ruleID,
		Name:   "Updated DNS Rule",
		Action: "ALLOW",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &firewalldnscontrolpolicies.FirewallDNSRules{
		ID:     ruleID,
		Name:   "Updated DNS Rule",
		Action: "ALLOW",
		State:  "ENABLED",
		Order:  1,
	}

	result, err := firewalldnscontrolpolicies.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated DNS Rule", result.Name)
}

func TestFirewallDNS_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallDnsRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = firewalldnscontrolpolicies.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

func TestFirewallDNS_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallDnsRules"

	server.On("GET", path, common.SuccessResponse([]firewalldnscontrolpolicies.FirewallDNSRules{
		{ID: 1, Name: "Rule 1", Action: "BLOCK"},
		{ID: 2, Name: "Rule 2", Action: "ALLOW"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewalldnscontrolpolicies.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Firewall IPS Control Policies - SDK Function Tests
// =====================================================

func TestFirewallIPS_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallIpsRules/12345"

	server.On("GET", path, common.SuccessResponse(firewallipscontrolpolicies.FirewallIPSRules{
		ID:     ruleID,
		Name:   "IPS Rule",
		Action: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewallipscontrolpolicies.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
}

func TestFirewallIPS_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "IPS Rule"
	path := "/zia/api/v1/firewallIpsRules"

	server.On("GET", path, common.SuccessResponse([]firewallipscontrolpolicies.FirewallIPSRules{
		{ID: 1, Name: "Other Rule", Action: "ALLOW"},
		{ID: 2, Name: ruleName, Action: "BLOCK"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewallipscontrolpolicies.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
}

func TestFirewallIPS_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallIpsRules"

	server.On("POST", path, common.SuccessResponse(firewallipscontrolpolicies.FirewallIPSRules{
		ID:     100,
		Name:   "New IPS Rule",
		Action: "BLOCK",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newRule := &firewallipscontrolpolicies.FirewallIPSRules{
		Name:   "New IPS Rule",
		Action: "BLOCK",
		State:  "ENABLED",
		Order:  1,
	}

	result, err := firewallipscontrolpolicies.Create(context.Background(), service, newRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestFirewallIPS_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallIpsRules/12345"

	server.On("PUT", path, common.SuccessResponse(firewallipscontrolpolicies.FirewallIPSRules{
		ID:     ruleID,
		Name:   "Updated IPS Rule",
		Action: "MONITOR",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateRule := &firewallipscontrolpolicies.FirewallIPSRules{
		ID:     ruleID,
		Name:   "Updated IPS Rule",
		Action: "MONITOR",
		State:  "ENABLED",
		Order:  1,
	}

	result, err := firewallipscontrolpolicies.Update(context.Background(), service, ruleID, updateRule)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated IPS Rule", result.Name)
}

func TestFirewallIPS_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	path := "/zia/api/v1/firewallIpsRules/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = firewallipscontrolpolicies.Delete(context.Background(), service, ruleID)

	require.NoError(t, err)
}

func TestFirewallIPS_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/firewallIpsRules"

	server.On("GET", path, common.SuccessResponse([]firewallipscontrolpolicies.FirewallIPSRules{
		{ID: 1, Name: "Rule 1", Action: "BLOCK"},
		{ID: 2, Name: "Rule 2", Action: "MONITOR"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := firewallipscontrolpolicies.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests
// =====================================================

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
