// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

func TestIPDestinationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPDestinationGroups JSON marshaling", func(t *testing.T) {
		group := ipdestinationgroups.IPDestinationGroups{
			ID:          12345,
			Name:        "External-Servers",
			Description: "External server IP addresses",
			Type:        "DSTN_IP",
			Addresses:   []string{"192.168.1.0/24", "10.0.0.1", "example.com"},
			Countries:   []string{"US", "CA", "GB"},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"External-Servers"`)
		assert.Contains(t, string(data), `"type":"DSTN_IP"`)
	})

	t.Run("IPDestinationGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Cloud-Services",
			"description": "Cloud provider IP ranges",
			"type": "DSTN_FQDN",
			"addresses": ["*.amazonaws.com", "*.azure.com"],
			"ipCategories": ["CLOUD_SERVICES", "SAAS"],
			"countries": ["US"],
			"isNonEditable": true
		}`

		var group ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "DSTN_FQDN", group.Type)
		assert.Len(t, group.Addresses, 2)
		assert.Len(t, group.IPCategories, 2)
		assert.True(t, group.IsNonEditable)
	})
}

func TestNetworkServices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkServices JSON marshaling", func(t *testing.T) {
		service := networkservices.NetworkServices{
			ID:          12345,
			Name:        "Custom-HTTP-Service",
			Description: "Custom HTTP traffic on non-standard ports",
			Type:        "CUSTOM",
			SrcTCPPorts: []networkservices.NetworkPorts{
				{Start: 1024, End: 65535},
			},
			DestTCPPorts: []networkservices.NetworkPorts{
				{Start: 8080, End: 8080},
				{Start: 8443, End: 8443},
			},
		}

		data, err := json.Marshal(service)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"CUSTOM"`)
		assert.Contains(t, string(data), `"srcTcpPorts"`)
		assert.Contains(t, string(data), `"destTcpPorts"`)
	})

	t.Run("NetworkServices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "DNS-Service",
			"description": "DNS traffic",
			"type": "PREDEFINED",
			"srcUdpPorts": [
				{"start": 1024, "end": 65535}
			],
			"destUdpPorts": [
				{"start": 53, "end": 53}
			],
			"destTcpPorts": [
				{"start": 53, "end": 53}
			],
			"isNameL10nTag": false,
			"creatorContext": "ZIA"
		}`

		var service networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonData), &service)
		require.NoError(t, err)

		assert.Equal(t, 54321, service.ID)
		assert.Equal(t, "PREDEFINED", service.Type)
		assert.Len(t, service.SrcUDPPorts, 1)
		assert.Len(t, service.DestUDPPorts, 1)
		assert.Equal(t, 53, service.DestUDPPorts[0].Start)
	})

	t.Run("NetworkPorts JSON marshaling", func(t *testing.T) {
		ports := networkservices.NetworkPorts{
			Start: 443,
			End:   443,
		}

		data, err := json.Marshal(ports)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"start":443`)
		assert.Contains(t, string(data), `"end":443`)
	})
}

func TestPolicyResources_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse IP destination groups list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Group-1", "type": "DSTN_IP"},
			{"id": 2, "name": "Group-2", "type": "DSTN_FQDN"},
			{"id": 3, "name": "Group-3", "type": "DSTN_IP"}
		]`

		var groups []ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonResponse), &groups)
		require.NoError(t, err)

		assert.Len(t, groups, 3)
		assert.Equal(t, "DSTN_FQDN", groups[1].Type)
	})

	t.Run("Parse network services list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "HTTP", "type": "PREDEFINED"},
			{"id": 2, "name": "HTTPS", "type": "PREDEFINED"},
			{"id": 3, "name": "Custom-App", "type": "CUSTOM"}
		]`

		var services []networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonResponse), &services)
		require.NoError(t, err)

		assert.Len(t, services, 3)
		assert.Equal(t, "CUSTOM", services[2].Type)
	})
}

