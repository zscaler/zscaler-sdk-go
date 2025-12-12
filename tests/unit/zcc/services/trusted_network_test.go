// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/trusted_network"
)

func TestTrustedNetwork_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TrustedNetwork JSON marshaling", func(t *testing.T) {
		network := trusted_network.TrustedNetwork{
			ID:                     "tn-123",
			Active:                 true,
			CompanyID:              "company-456",
			ConditionType:          1,
			NetworkName:            "Corporate Network",
			DnsServers:             "8.8.8.8,8.8.4.4",
			DnsSearchDomains:       "corp.example.com",
			Ssids:                  "CorpWiFi,GuestWiFi",
			TrustedGateways:        "192.168.1.1",
			TrustedSubnets:         "10.0.0.0/8,172.16.0.0/12",
			TrustedDhcpServers:     "192.168.1.100",
			TrustedEgressIps:       "203.0.113.1",
			Hostnames:              "gateway.corp.example.com",
			ResolvedIpsForHostname: "10.0.0.1",
			CreatedBy:              "admin@example.com",
			EditedBy:               "admin@example.com",
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"tn-123"`)
		assert.Contains(t, string(data), `"active":true`)
		assert.Contains(t, string(data), `"networkName":"Corporate Network"`)
		assert.Contains(t, string(data), `"dnsServers":"8.8.8.8,8.8.4.4"`)
	})

	t.Run("TrustedNetwork JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "tn-789",
			"active": true,
			"companyId": "company-456",
			"conditionType": 2,
			"networkName": "Branch Office",
			"dnsServers": "1.1.1.1",
			"dnsSearchDomains": "branch.example.com",
			"ssids": "BranchWiFi",
			"trustedGateways": "192.168.10.1",
			"trustedSubnets": "192.168.10.0/24",
			"trustedDhcpServers": "192.168.10.100",
			"guid": "guid-abc-123"
		}`

		var network trusted_network.TrustedNetwork
		err := json.Unmarshal([]byte(jsonData), &network)
		require.NoError(t, err)

		assert.Equal(t, "tn-789", network.ID)
		assert.True(t, network.Active)
		assert.Equal(t, "Branch Office", network.NetworkName)
		assert.Equal(t, 2, network.ConditionType)
		assert.Equal(t, "1.1.1.1", network.DnsServers)
		assert.Equal(t, "BranchWiFi", network.Ssids)
		assert.Equal(t, "guid-abc-123", network.Guid)
	})

	t.Run("TrustedNetworksResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 2,
			"trustedNetworkContracts": [
				{
					"id": "tn-001",
					"networkName": "Network 1",
					"active": true
				},
				{
					"id": "tn-002",
					"networkName": "Network 2",
					"active": false
				}
			]
		}`

		var response trusted_network.TrustedNetworksResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 2, response.TotalCount)
		assert.Len(t, response.TrustedNetworkContracts, 2)
		assert.Equal(t, "Network 1", response.TrustedNetworkContracts[0].NetworkName)
		assert.True(t, response.TrustedNetworkContracts[0].Active)
		assert.Equal(t, "Network 2", response.TrustedNetworkContracts[1].NetworkName)
		assert.False(t, response.TrustedNetworkContracts[1].Active)
	})
}

func TestTrustedNetwork_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse trusted networks list response", func(t *testing.T) {
		jsonResponse := `{
			"totalCount": 3,
			"trustedNetworkContracts": [
				{
					"id": "tn-001",
					"networkName": "Headquarters",
					"active": true,
					"conditionType": 1,
					"ssids": "HQ-WiFi",
					"trustedSubnets": "10.0.0.0/8"
				},
				{
					"id": "tn-002",
					"networkName": "Branch Office 1",
					"active": true,
					"conditionType": 2,
					"dnsServers": "192.168.1.1"
				},
				{
					"id": "tn-003",
					"networkName": "Remote Office",
					"active": false,
					"conditionType": 1
				}
			]
		}`

		var response trusted_network.TrustedNetworksResponse
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.TotalCount)
		assert.Len(t, response.TrustedNetworkContracts, 3)
		
		// Check first network
		assert.Equal(t, "Headquarters", response.TrustedNetworkContracts[0].NetworkName)
		assert.True(t, response.TrustedNetworkContracts[0].Active)
		assert.Equal(t, "HQ-WiFi", response.TrustedNetworkContracts[0].Ssids)
		
		// Check second network
		assert.Equal(t, "Branch Office 1", response.TrustedNetworkContracts[1].NetworkName)
		assert.Equal(t, "192.168.1.1", response.TrustedNetworkContracts[1].DnsServers)
		
		// Check third network
		assert.Equal(t, "Remote Office", response.TrustedNetworkContracts[2].NetworkName)
		assert.False(t, response.TrustedNetworkContracts[2].Active)
	})
}

