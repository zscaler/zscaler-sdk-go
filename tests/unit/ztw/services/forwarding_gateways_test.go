// Package services provides unit tests for ZTW services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

func TestDNSForwardingGateway_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DNSGateway JSON marshaling", func(t *testing.T) {
		gw := dns_forwarding_gateway.DNSGateway{
			ID:                           12345,
			Name:                         "DNS-Gateway-1",
			Type:                         "ZIA",
			FailureBehavior:              "FAIL_OPEN",
			DNSGatewayType:               "PRIMARY",
			PrimaryIP:                    "10.0.0.1",
			SecondaryIP:                  "10.0.0.2",
			ECDNSGatewayOptionsPrimary:   "AUTO",
			ECDNSGatewayOptionsSecondary: "MANUAL",
			LastModifiedTime:             1699000000,
		}

		data, err := json.Marshal(gw)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"DNS-Gateway-1"`)
		assert.Contains(t, string(data), `"primaryIp":"10.0.0.1"`)
	})

	t.Run("DNSGateway JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "DNS-Gateway-2",
			"type": "ECSELF",
			"failureBehavior": "FAIL_CLOSED",
			"dnsGatewayType": "SECONDARY",
			"primaryIp": "192.168.1.1",
			"secondaryIp": "192.168.1.2",
			"lastModifiedTime": 1699500000
		}`

		var gw dns_forwarding_gateway.DNSGateway
		err := json.Unmarshal([]byte(jsonData), &gw)
		require.NoError(t, err)

		assert.Equal(t, 54321, gw.ID)
		assert.Equal(t, "DNS-Gateway-2", gw.Name)
		assert.Equal(t, "ECSELF", gw.Type)
		assert.Equal(t, "192.168.1.1", gw.PrimaryIP)
	})
}

func TestZIAForwardingGateway_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ECGateway JSON marshaling", func(t *testing.T) {
		gw := zia_forwarding_gateway.ECGateway{
			ID:              12345,
			Name:            "ZIA-Gateway-1",
			Description:     "Primary ZIA forwarding gateway",
			FailClosed:      true,
			ManualPrimary:   "10.0.0.1",
			ManualSecondary: "10.0.0.2",
			PrimaryType:     "AUTO",
			SecondaryType:   "MANUAL_OVERRIDE",
			Type:            "ZIA",
			FailureBehavior: "FAIL_CLOSED",
		}

		data, err := json.Marshal(gw)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"ZIA-Gateway-1"`)
		assert.Contains(t, string(data), `"failClosed":true`)
		assert.Contains(t, string(data), `"primaryType":"AUTO"`)
	})

	t.Run("ECGateway JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "ZIA-Gateway-2",
			"description": "Secondary ZIA gateway",
			"failClosed": false,
			"manualPrimary": "172.16.0.1",
			"manualSecondary": "172.16.0.2",
			"primaryType": "DC",
			"secondaryType": "SUBCLOUD",
			"type": "ECSELF",
			"lastModifiedTime": 1699500000
		}`

		var gw zia_forwarding_gateway.ECGateway
		err := json.Unmarshal([]byte(jsonData), &gw)
		require.NoError(t, err)

		assert.Equal(t, 54321, gw.ID)
		assert.Equal(t, "ZIA-Gateway-2", gw.Name)
		assert.Equal(t, "DC", gw.PrimaryType)
		assert.Equal(t, "SUBCLOUD", gw.SecondaryType)
	})
}

func TestForwardingGateways_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse DNS gateways list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "DNS-GW-1", "type": "ZIA"},
			{"id": 2, "name": "DNS-GW-2", "type": "ECSELF"},
			{"id": 3, "name": "DNS-GW-3", "type": "ZIA"}
		]`

		var gateways []dns_forwarding_gateway.DNSGateway
		err := json.Unmarshal([]byte(jsonResponse), &gateways)
		require.NoError(t, err)

		assert.Len(t, gateways, 3)
		assert.Equal(t, "DNS-GW-1", gateways[0].Name)
	})

	t.Run("Parse ZIA gateways list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "ZIA-GW-1", "primaryType": "AUTO"},
			{"id": 2, "name": "ZIA-GW-2", "primaryType": "DC"}
		]`

		var gateways []zia_forwarding_gateway.ECGateway
		err := json.Unmarshal([]byte(jsonResponse), &gateways)
		require.NoError(t, err)

		assert.Len(t, gateways, 2)
	})
}

