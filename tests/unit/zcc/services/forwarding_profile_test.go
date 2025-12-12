// Package services provides unit tests for ZCC services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/forwarding_profile"
)

func TestForwardingProfile_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingProfile JSON marshaling", func(t *testing.T) {
		profile := forwarding_profile.ForwardingProfile{
			ID:                       forwarding_profile.IntOrString(123),
			Name:                     "Default Forwarding Profile",
			Active:                   "true",
			ConditionType:            1,
			EnableLWFDriver:          "true",
			EnableSplitVpnTN:         1,
			EvaluateTrustedNetwork:   1,
			PredefinedTnAll:          true,
			PredefinedTrustedNetworks: false,
			TrustedSubnets:           "10.0.0.0/8",
			TrustedGateways:          "192.168.1.1",
			DnsServers:               "8.8.8.8",
		}

		data, err := json.Marshal(profile)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"Default Forwarding Profile"`)
		assert.Contains(t, string(data), `"active":"true"`)
		assert.Contains(t, string(data), `"conditionType":1`)
	})

	t.Run("ForwardingProfile JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 456,
			"name": "Custom Forwarding Profile",
			"active": "true",
			"conditionType": 2,
			"enableLWFDriver": "false",
			"enableSplitVpnTN": 0,
			"evaluateTrustedNetwork": 1,
			"predefinedTnAll": false,
			"predefinedTrustedNetworks": true,
			"trustedSubnets": "172.16.0.0/12",
			"trustedGateways": "10.0.0.1",
			"forwardingProfileActions": [],
			"forwardingProfileZpaActions": []
		}`

		var profile forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonData), &profile)
		require.NoError(t, err)

		assert.Equal(t, forwarding_profile.IntOrString(456), profile.ID)
		assert.Equal(t, "Custom Forwarding Profile", profile.Name)
		assert.Equal(t, "true", profile.Active)
		assert.Equal(t, 2, profile.ConditionType)
		assert.True(t, profile.PredefinedTrustedNetworks)
	})

	t.Run("ForwardingProfileAction JSON marshaling", func(t *testing.T) {
		action := forwarding_profile.ForwardingProfileAction{
			ActionType:         1,
			NetworkType:        1,
			RedirectWebTraffic: 1,
			SystemProxy:        0,
			EnablePacketTunnel: 1,
			PrimaryTransport:   1,
			DTLSTimeout:        30,
			TLSTimeout:         30,
			UDPTimeout:         30,
			MtuForZadapter:     1500,
		}

		data, err := json.Marshal(action)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"actionType":1`)
		assert.Contains(t, string(data), `"networkType":1`)
		assert.Contains(t, string(data), `"redirectWebTraffic":1`)
	})

	t.Run("SystemProxyData JSON marshaling", func(t *testing.T) {
		proxyData := forwarding_profile.SystemProxyData{
			EnableAutoDetect:        1,
			EnablePAC:               0,
			EnableProxyServer:       1,
			ProxyServerAddress:      "proxy.example.com",
			ProxyServerPort:         "8080",
			BypassProxyForPrivateIP: 1,
			ProxyAction:             1,
		}

		data, err := json.Marshal(proxyData)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"proxyServerAddress":"proxy.example.com"`)
		assert.Contains(t, string(data), `"proxyServerPort":"8080"`)
	})

	t.Run("IntOrString unmarshaling from int", func(t *testing.T) {
		jsonData := `{"id": 123}`
		var result struct {
			ID forwarding_profile.IntOrString `json:"id"`
		}
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.Equal(t, forwarding_profile.IntOrString(123), result.ID)
	})

	t.Run("IntOrString unmarshaling from string", func(t *testing.T) {
		jsonData := `{"id": "456"}`
		var result struct {
			ID forwarding_profile.IntOrString `json:"id"`
		}
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)
		assert.Equal(t, forwarding_profile.IntOrString(456), result.ID)
	})
}

func TestForwardingProfile_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse forwarding profile list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "Default Profile",
				"active": "true",
				"conditionType": 1,
				"forwardingProfileActions": [
					{
						"actionType": 1,
						"networkType": 1,
						"redirectWebTraffic": 1
					}
				],
				"forwardingProfileZpaActions": []
			},
			{
				"id": 2,
				"name": "Trusted Network Profile",
				"active": "true",
				"conditionType": 2,
				"forwardingProfileActions": [],
				"forwardingProfileZpaActions": [
					{
						"actionType": 2,
						"networkType": 1
					}
				]
			}
		]`

		var profiles []forwarding_profile.ForwardingProfile
		err := json.Unmarshal([]byte(jsonResponse), &profiles)
		require.NoError(t, err)

		assert.Len(t, profiles, 2)
		assert.Equal(t, "Default Profile", profiles[0].Name)
		assert.Len(t, profiles[0].ForwardingProfileActions, 1)
		assert.Equal(t, "Trusted Network Profile", profiles[1].Name)
		assert.Len(t, profiles[1].ForwardingProfileZpaActions, 1)
	})
}

