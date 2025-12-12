// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StaticIP represents a static IP configuration
type StaticIP struct {
	ID               int    `json:"id,omitempty"`
	IPAddress        string `json:"ipAddress,omitempty"`
	GeoOverride      bool   `json:"geoOverride,omitempty"`
	Latitude         string `json:"latitude,omitempty"`
	Longitude        string `json:"longitude,omitempty"`
	RoutableIP       bool   `json:"routableIP,omitempty"`
	LastModifiedTime int    `json:"lastModifiedTime,omitempty"`
	Comment          string `json:"comment,omitempty"`
	Managed          bool   `json:"managed,omitempty"`
}

// GRETunnel represents a GRE tunnel configuration
type GRETunnel struct {
	ID                  int      `json:"id,omitempty"`
	SourceIP            string   `json:"sourceIp,omitempty"`
	PrimaryDestVip      *DestVip `json:"primaryDestVip,omitempty"`
	SecondaryDestVip    *DestVip `json:"secondaryDestVip,omitempty"`
	InternalIPRange     string   `json:"internalIpRange,omitempty"`
	WithinCountry       bool     `json:"withinCountry,omitempty"`
	Comment             string   `json:"comment,omitempty"`
	IPUnnumbered        bool     `json:"ipUnnumbered,omitempty"`
	LastModificationTime int     `json:"lastModificationTime,omitempty"`
}

// DestVip represents a destination VIP
type DestVip struct {
	ID         int    `json:"id,omitempty"`
	VirtualIP  string `json:"virtualIp,omitempty"`
	PrivateIP  string `json:"privateIp,omitempty"`
	Datacenter string `json:"datacenter,omitempty"`
	City       string `json:"city,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
}

// VPNCredential represents VPN credentials
type VPNCredential struct {
	ID           int    `json:"id,omitempty"`
	Type         string `json:"type,omitempty"`
	FQDN         string `json:"fqdn,omitempty"`
	IPAddress    string `json:"ipAddress,omitempty"`
	PreSharedKey string `json:"preSharedKey,omitempty"`
	Comments     string `json:"comments,omitempty"`
}

func TestTrafficForwarding_StaticIP(t *testing.T) {
	t.Parallel()

	t.Run("StaticIP JSON marshaling", func(t *testing.T) {
		staticIP := StaticIP{
			ID:          12345,
			IPAddress:   "203.0.113.10",
			GeoOverride: true,
			Latitude:    "37.7749",
			Longitude:   "-122.4194",
			RoutableIP:  true,
			Comment:     "Primary static IP",
			Managed:     false,
		}

		data, err := json.Marshal(staticIP)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"ipAddress":"203.0.113.10"`)
		assert.Contains(t, string(data), `"geoOverride":true`)
	})

	t.Run("StaticIP JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"ipAddress": "198.51.100.20",
			"geoOverride": false,
			"routableIP": true,
			"lastModifiedTime": 1699000000,
			"comment": "Secondary IP",
			"managed": true
		}`

		var staticIP StaticIP
		err := json.Unmarshal([]byte(jsonData), &staticIP)
		require.NoError(t, err)

		assert.Equal(t, 54321, staticIP.ID)
		assert.Equal(t, "198.51.100.20", staticIP.IPAddress)
		assert.True(t, staticIP.Managed)
	})
}

func TestTrafficForwarding_GRETunnel(t *testing.T) {
	t.Parallel()

	t.Run("GRETunnel JSON marshaling", func(t *testing.T) {
		tunnel := GRETunnel{
			ID:              12345,
			SourceIP:        "203.0.113.1",
			InternalIPRange: "10.0.0.0/8",
			WithinCountry:   true,
			Comment:         "Primary GRE tunnel",
			IPUnnumbered:    false,
			PrimaryDestVip: &DestVip{
				ID:          100,
				VirtualIP:   "185.46.212.88",
				PrivateIP:   "10.0.0.1",
				Datacenter:  "US-SJC",
				City:        "San Jose",
				CountryCode: "US",
			},
		}

		data, err := json.Marshal(tunnel)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"sourceIp":"203.0.113.1"`)
		assert.Contains(t, string(data), `"withinCountry":true`)
		assert.Contains(t, string(data), `"primaryDestVip"`)
	})

	t.Run("GRETunnel JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"sourceIp": "198.51.100.1",
			"primaryDestVip": {
				"id": 200,
				"virtualIp": "185.46.212.90",
				"datacenter": "US-LAX",
				"city": "Los Angeles",
				"countryCode": "US"
			},
			"secondaryDestVip": {
				"id": 201,
				"virtualIp": "185.46.212.91",
				"datacenter": "US-SFO",
				"city": "San Francisco",
				"countryCode": "US"
			},
			"internalIpRange": "172.16.0.0/12",
			"withinCountry": false,
			"comment": "Backup tunnel",
			"ipUnnumbered": true
		}`

		var tunnel GRETunnel
		err := json.Unmarshal([]byte(jsonData), &tunnel)
		require.NoError(t, err)

		assert.Equal(t, 54321, tunnel.ID)
		assert.NotNil(t, tunnel.PrimaryDestVip)
		assert.NotNil(t, tunnel.SecondaryDestVip)
		assert.True(t, tunnel.IPUnnumbered)
	})
}

func TestTrafficForwarding_VPNCredential(t *testing.T) {
	t.Parallel()

	t.Run("VPNCredential UFQDN JSON marshaling", func(t *testing.T) {
		vpn := VPNCredential{
			ID:           12345,
			Type:         "UFQDN",
			FQDN:         "vpn.company.com@zscaler.net",
			PreSharedKey: "secret-psk-key",
			Comments:     "Primary VPN",
		}

		data, err := json.Marshal(vpn)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"type":"UFQDN"`)
		assert.Contains(t, string(data), `"fqdn":"vpn.company.com@zscaler.net"`)
	})

	t.Run("VPNCredential IP JSON marshaling", func(t *testing.T) {
		vpn := VPNCredential{
			ID:           54321,
			Type:         "IP",
			IPAddress:    "203.0.113.50",
			PreSharedKey: "another-psk-key",
			Comments:     "IP-based VPN",
		}

		data, err := json.Marshal(vpn)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"type":"IP"`)
		assert.Contains(t, string(data), `"ipAddress":"203.0.113.50"`)
	})
}

func TestTrafficForwarding_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse static IPs list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "ipAddress": "203.0.113.1", "routableIP": true},
			{"id": 2, "ipAddress": "203.0.113.2", "routableIP": true},
			{"id": 3, "ipAddress": "203.0.113.3", "routableIP": false}
		]`

		var staticIPs []StaticIP
		err := json.Unmarshal([]byte(jsonResponse), &staticIPs)
		require.NoError(t, err)

		assert.Len(t, staticIPs, 3)
		assert.False(t, staticIPs[2].RoutableIP)
	})

	t.Run("Parse GRE tunnels list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "sourceIp": "10.0.0.1", "withinCountry": true},
			{"id": 2, "sourceIp": "10.0.0.2", "withinCountry": false}
		]`

		var tunnels []GRETunnel
		err := json.Unmarshal([]byte(jsonResponse), &tunnels)
		require.NoError(t, err)

		assert.Len(t, tunnels, 2)
		assert.True(t, tunnels[0].WithinCountry)
	})

	t.Run("Parse VPN credentials list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "type": "UFQDN", "fqdn": "vpn1.company.com"},
			{"id": 2, "type": "IP", "ipAddress": "203.0.113.10"}
		]`

		var vpnCreds []VPNCredential
		err := json.Unmarshal([]byte(jsonResponse), &vpnCreds)
		require.NoError(t, err)

		assert.Len(t, vpnCreds, 2)
		assert.Equal(t, "UFQDN", vpnCreds[0].Type)
		assert.Equal(t, "IP", vpnCreds[1].Type)
	})
}

