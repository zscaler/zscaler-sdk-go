// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipsourcegroups"
)

func TestIPSourceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPSourceGroups JSON marshaling", func(t *testing.T) {
		group := ipsourcegroups.IPSourceGroups{
			ID:            12345,
			Name:          "Corporate Networks",
			Description:   "Internal corporate network ranges",
			IPAddresses:   []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"},
			IsNonEditable: false,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"ipAddresses":["10.0.0.0/8","172.16.0.0/12","192.168.0.0/16"]`)
	})

	t.Run("IPSourceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "VPN Clients",
			"description": "VPN client IP ranges",
			"ipAddresses": ["10.10.0.0/16", "10.20.0.0/16"],
			"isNonEditable": false
		}`

		var group ipsourcegroups.IPSourceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Len(t, group.IPAddresses, 2)
	})

	t.Run("IPSourceGroups predefined", func(t *testing.T) {
		jsonData := `{
			"id": 1,
			"name": "RFC 1918",
			"description": "Private IP address ranges",
			"ipAddresses": ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"],
			"isNonEditable": true
		}`

		var group ipsourcegroups.IPSourceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.True(t, group.IsNonEditable)
	})
}

func TestIPDestinationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPDestinationGroups JSON marshaling", func(t *testing.T) {
		group := ipdestinationgroups.IPDestinationGroups{
			ID:            12345,
			Name:          "Cloud Services",
			Description:   "Public cloud service IPs",
			Addresses:     []string{"52.0.0.0/8", "34.0.0.0/8"},
			Type:          "DSTN_IP",
			IsNonEditable: false,
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"DSTN_IP"`)
	})

	t.Run("IPDestinationGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Partner Networks",
			"description": "Partner organization IPs",
			"addresses": ["203.0.113.0/24", "198.51.100.0/24"],
			"type": "DSTN_IP",
			"countries": ["US", "CA"],
			"ipCategories": ["DATACENTER"],
			"isNonEditable": false
		}`

		var group ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Len(t, group.Addresses, 2)
		assert.Len(t, group.Countries, 2)
	})

	t.Run("IPDestinationGroups with FQDNs", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "API Endpoints",
			"description": "API server endpoints",
			"type": "DSTN_FQDN",
			"addresses": ["api.example.com", "gateway.example.com"]
		}`

		var group ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, "DSTN_FQDN", group.Type)
	})
}

func TestIPGroups_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse IP source groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Group 1", "ipAddresses": ["10.0.0.0/8"]},
			{"id": 2, "name": "Group 2", "ipAddresses": ["172.16.0.0/12"]},
			{"id": 3, "name": "Group 3", "ipAddresses": ["192.168.0.0/16"]}
		]`

		var groups []ipsourcegroups.IPSourceGroups
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
	})

	t.Run("Parse IP destination groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Dest 1", "type": "DSTN_IP"},
			{"id": 2, "name": "Dest 2", "type": "DSTN_FQDN"},
			{"id": 3, "name": "Dest 3", "type": "DSTN_IP"}
		]`

		var groups []ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, "DSTN_FQDN", groups[1].Type)
	})
}

