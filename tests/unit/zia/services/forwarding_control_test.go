// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ForwardingRule represents a forwarding control rule
type ForwardingRule struct {
	ID              int      `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	Order           int      `json:"order,omitempty"`
	Rank            int      `json:"rank,omitempty"`
	State           string   `json:"state,omitempty"`
	Type            string   `json:"type,omitempty"`
	ForwardMethod   string   `json:"forwardMethod,omitempty"`
	SrcIps          []string `json:"srcIps,omitempty"`
	DestAddresses   []string `json:"destAddresses,omitempty"`
	DestCountries   []string `json:"destCountries,omitempty"`
	DestIpCategories []string `json:"destIpCategories,omitempty"`
	NwApplications  []string `json:"nwApplications,omitempty"`
}

// ZPAGateway represents a ZPA gateway configuration
type ZPAGateway struct {
	ID              int        `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	Type            string     `json:"type,omitempty"`
	ZPAServerGroup  *ZPAServerGroup `json:"zpaServerGroup,omitempty"`
	ZPAAppSegments  []ZPAAppSegment `json:"zpaAppSegments,omitempty"`
	LastModifiedTime int       `json:"lastModifiedTime,omitempty"`
}

// ZPAServerGroup represents a ZPA server group
type ZPAServerGroup struct {
	ExternalID string `json:"externalId,omitempty"`
	Name       string `json:"name,omitempty"`
}

// ZPAAppSegment represents a ZPA application segment
type ZPAAppSegment struct {
	ExternalID string `json:"externalId,omitempty"`
	Name       string `json:"name,omitempty"`
}

// ProxyGateway represents a proxy gateway configuration
type ProxyGateway struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Type            string `json:"type,omitempty"`
	PrimaryProxy    string `json:"primaryProxy,omitempty"`
	SecondaryProxy  string `json:"secondaryProxy,omitempty"`
	PrimaryPort     int    `json:"primaryPort,omitempty"`
	SecondaryPort   int    `json:"secondaryPort,omitempty"`
	FailClosed      bool   `json:"failClosed,omitempty"`
}

func TestForwardingRule_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingRule JSON marshaling", func(t *testing.T) {
		rule := ForwardingRule{
			ID:            12345,
			Name:          "Forward to ZPA",
			Description:   "Forward internal traffic to ZPA",
			Order:         1,
			Rank:          7,
			State:         "ENABLED",
			Type:          "FORWARDING",
			ForwardMethod: "ZPA",
			SrcIps:        []string{"10.0.0.0/8"},
			DestAddresses: []string{"192.168.1.0/24"},
			DestCountries: []string{"US"},
		}

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"forwardMethod":"ZPA"`)
	})

	t.Run("ForwardingRule JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Direct Traffic",
			"description": "Allow direct traffic",
			"order": 2,
			"rank": 5,
			"state": "ENABLED",
			"type": "FORWARDING",
			"forwardMethod": "DIRECT",
			"srcIps": ["172.16.0.0/12"],
			"destAddresses": ["8.8.8.8"],
			"destCountries": ["US", "CA"],
			"destIpCategories": ["DNS"],
			"nwApplications": ["DNS"]
		}`

		var rule ForwardingRule
		err := json.Unmarshal([]byte(jsonData), &rule)
		require.NoError(t, err)

		assert.Equal(t, 54321, rule.ID)
		assert.Equal(t, "DIRECT", rule.ForwardMethod)
		assert.Len(t, rule.DestCountries, 2)
	})
}

func TestZPAGateway_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ZPAGateway JSON marshaling", func(t *testing.T) {
		gateway := ZPAGateway{
			ID:          12345,
			Name:        "ZPA Gateway 1",
			Description: "Primary ZPA gateway",
			Type:        "ZPA",
			ZPAServerGroup: &ZPAServerGroup{
				ExternalID: "zpa-sg-123",
				Name:       "Server Group 1",
			},
			ZPAAppSegments: []ZPAAppSegment{
				{ExternalID: "zpa-app-1", Name: "App 1"},
				{ExternalID: "zpa-app-2", Name: "App 2"},
			},
		}

		data, err := json.Marshal(gateway)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"zpaServerGroup"`)
		assert.Contains(t, string(data), `"zpaAppSegments"`)
	})

	t.Run("ZPAGateway JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "ZPA Gateway 2",
			"type": "ZPA",
			"zpaServerGroup": {
				"externalId": "zpa-sg-456",
				"name": "Server Group 2"
			},
			"zpaAppSegments": [
				{"externalId": "zpa-app-3", "name": "App 3"}
			],
			"lastModifiedTime": 1699000000
		}`

		var gateway ZPAGateway
		err := json.Unmarshal([]byte(jsonData), &gateway)
		require.NoError(t, err)

		assert.Equal(t, 54321, gateway.ID)
		assert.NotNil(t, gateway.ZPAServerGroup)
		assert.Len(t, gateway.ZPAAppSegments, 1)
	})
}

func TestProxyGateway_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ProxyGateway JSON marshaling", func(t *testing.T) {
		proxy := ProxyGateway{
			ID:             12345,
			Name:           "Proxy Gateway 1",
			Description:    "Primary proxy gateway",
			Type:           "PROXYCHAIN",
			PrimaryProxy:   "proxy1.company.com",
			SecondaryProxy: "proxy2.company.com",
			PrimaryPort:    8080,
			SecondaryPort:  8081,
			FailClosed:     true,
		}

		data, err := json.Marshal(proxy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"primaryProxy":"proxy1.company.com"`)
		assert.Contains(t, string(data), `"failClosed":true`)
	})

	t.Run("ProxyGateway JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Proxy Gateway 2",
			"type": "PROXYCHAIN",
			"primaryProxy": "proxy-a.company.com",
			"secondaryProxy": "proxy-b.company.com",
			"primaryPort": 3128,
			"secondaryPort": 3129,
			"failClosed": false
		}`

		var proxy ProxyGateway
		err := json.Unmarshal([]byte(jsonData), &proxy)
		require.NoError(t, err)

		assert.Equal(t, 54321, proxy.ID)
		assert.Equal(t, 3128, proxy.PrimaryPort)
		assert.False(t, proxy.FailClosed)
	})
}

func TestForwardingControl_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse forwarding rules list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Rule 1", "forwardMethod": "ZPA", "state": "ENABLED"},
			{"id": 2, "name": "Rule 2", "forwardMethod": "DIRECT", "state": "ENABLED"},
			{"id": 3, "name": "Rule 3", "forwardMethod": "PROXYCHAIN", "state": "DISABLED"}
		]`

		var rules []ForwardingRule
		err := json.Unmarshal([]byte(jsonResponse), &rules)
		require.NoError(t, err)

		assert.Len(t, rules, 3)
		assert.Equal(t, "ZPA", rules[0].ForwardMethod)
		assert.Equal(t, "DISABLED", rules[2].State)
	})

	t.Run("Parse ZPA gateways list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Gateway 1", "type": "ZPA"},
			{"id": 2, "name": "Gateway 2", "type": "ZPA"}
		]`

		var gateways []ZPAGateway
		err := json.Unmarshal([]byte(jsonResponse), &gateways)
		require.NoError(t, err)

		assert.Len(t, gateways, 2)
	})
}

