// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/dc_exclusions"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/extranet"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/greinternalipranges"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnelinfo"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/ipv6_config"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/region/datacenter"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/region/geo_coordinates"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/region/ip_address"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/region/search"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/sub_clouds"
	virtualipaddress "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/virtualipaddress"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/trafficforwarding/vpncredentials"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestVPNCredentials_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	vpnID := 12345
	path := "/zia/api/v1/vpnCredentials/12345"

	server.On("GET", path, common.SuccessResponse(vpncredentials.VPNCredentials{
		ID:        vpnID,
		Type:      "UFQDN",
		FQDN:      "vpn.company.com",
		IPAddress: "203.0.113.1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.Get(context.Background(), service, vpnID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, vpnID, result.ID)
	assert.Equal(t, "vpn.company.com", result.FQDN)
}

func TestVPNCredentials_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vpnCredentials"

	server.On("GET", path, common.SuccessResponse([]vpncredentials.VPNCredentials{
		{ID: 1, Type: "UFQDN", FQDN: "vpn1.company.com"},
		{ID: 2, Type: "IP", IPAddress: "203.0.113.1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestVPNCredentials_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vpnCredentials"

	server.On("POST", path, common.SuccessResponse(vpncredentials.VPNCredentials{
		ID:   99999,
		Type: "UFQDN",
		FQDN: "new.vpn.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newVPN := &vpncredentials.VPNCredentials{
		Type: "UFQDN",
		FQDN: "new.vpn.com",
	}

	result, _, err := vpncredentials.Create(context.Background(), service, newVPN)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestVPNCredentials_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	vpnID := 12345
	path := "/zia/api/v1/vpnCredentials/12345"

	server.On("PUT", path, common.SuccessResponse(vpncredentials.VPNCredentials{
		ID:   vpnID,
		FQDN: "updated.vpn.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateVPN := &vpncredentials.VPNCredentials{
		ID:   vpnID,
		FQDN: "updated.vpn.com",
	}

	result, _, err := vpncredentials.Update(context.Background(), service, vpnID, updateVPN)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated.vpn.com", result.FQDN)
}

func TestVPNCredentials_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	vpnID := 12345
	path := "/zia/api/v1/vpnCredentials/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	err = vpncredentials.Delete(context.Background(), service, vpnID)

	require.NoError(t, err)
}

func TestVPNCredentials_GetByType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	vpnType := "UFQDN"
	path := "/zia/api/v1/vpnCredentials"

	server.On("GET", path, common.SuccessResponse([]vpncredentials.VPNCredentials{
		{ID: 1, Type: vpnType, FQDN: "vpn1.company.com"},
		{ID: 2, Type: "IP", IPAddress: "203.0.113.1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.GetVPNByType(context.Background(), service, vpnType, nil, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 2)
}

func TestVPNCredentials_GetByFQDN_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	fqdn := "vpn1.company.com"
	path := "/zia/api/v1/vpnCredentials"

	server.On("GET", path, common.SuccessResponse([]vpncredentials.VPNCredentials{
		{ID: 1, Type: "UFQDN", FQDN: fqdn},
		{ID: 2, Type: "UFQDN", FQDN: "vpn2.company.com"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.GetByFQDN(context.Background(), service, fqdn)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, fqdn, result.FQDN)
}

func TestVPNCredentials_GetByIP_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ipAddress := "203.0.113.1"
	path := "/zia/api/v1/vpnCredentials"

	server.On("GET", path, common.SuccessResponse([]vpncredentials.VPNCredentials{
		{ID: 1, Type: "IP", IPAddress: ipAddress},
		{ID: 2, Type: "IP", IPAddress: "198.51.100.1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.GetByIP(context.Background(), service, ipAddress)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ipAddress, result.IPAddress)
}

func TestVPNCredentials_BulkDelete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vpnCredentials/bulkDelete"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ids := []int{1, 2, 3}
	_, err = vpncredentials.BulkDelete(context.Background(), service, ids)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestVPNCredentials_Structure(t *testing.T) {
	t.Parallel()

	t.Run("VPNCredentials JSON marshaling", func(t *testing.T) {
		vpn := vpncredentials.VPNCredentials{
			ID:           12345,
			Type:         "UFQDN",
			FQDN:         "vpn.company.com",
			IPAddress:    "203.0.113.1",
			PreSharedKey: "secret-key",
			Comments:     "Primary VPN",
		}

		data, err := json.Marshal(vpn)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"UFQDN"`)
		assert.Contains(t, string(data), `"fqdn":"vpn.company.com"`)
	})

	t.Run("VPNCredentials JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"type": "IP",
			"ipAddress": "198.51.100.1",
			"preSharedKey": "another-key",
			"comments": "Backup VPN"
		}`

		var vpn vpncredentials.VPNCredentials
		err := json.Unmarshal([]byte(jsonData), &vpn)
		require.NoError(t, err)

		assert.Equal(t, 54321, vpn.ID)
		assert.Equal(t, "IP", vpn.Type)
		assert.Equal(t, "198.51.100.1", vpn.IPAddress)
	})
}

// =====================================================
// staticips
// =====================================================

func TestStaticIPs_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	staticID := 12345
	path := "/zia/api/v1/staticIP/12345"
	server.On("GET", path, common.SuccessResponse(staticips.StaticIP{
		ID: staticID, IpAddress: "104.239.238.10", Comment: "tests-static-ip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := staticips.Get(context.Background(), service, staticID)
	require.NoError(t, err)
	assert.Equal(t, "104.239.238.10", result.IpAddress)
}

func TestStaticIPs_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/staticIP"
	server.On("GET", path, common.SuccessResponse([]staticips.StaticIP{
		{ID: 1, IpAddress: "104.239.238.10"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := staticips.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestStaticIPs_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/staticIP"
	server.On("POST", path, common.SuccessResponse(staticips.StaticIP{
		ID: 99999, IpAddress: "104.239.238.20", Comment: "tests-static-ip",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newIP := &staticips.StaticIP{IpAddress: "104.239.238.20", Comment: "tests-static-ip"}
	result, _, err := staticips.Create(context.Background(), service, newIP)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestStaticIPs_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	staticID := 12345
	path := "/zia/api/v1/staticIP/12345"
	server.On("PUT", path, common.SuccessResponse(staticips.StaticIP{
		ID: staticID, Comment: "updated-comment",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &staticips.StaticIP{ID: staticID, Comment: "updated-comment"}
	result, _, err := staticips.Update(context.Background(), service, staticID, update)
	require.NoError(t, err)
	assert.Equal(t, "updated-comment", result.Comment)
}

func TestStaticIPs_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/staticIP/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = staticips.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestStaticIPs_GetByIPAddress_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ip := "104.239.238.10"
	path := "/zia/api/v1/staticIP"
	server.On("GET", path, common.SuccessResponse([]staticips.StaticIP{
		{ID: 1, IpAddress: ip},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := staticips.GetByIPAddress(context.Background(), service, ip)
	require.NoError(t, err)
	assert.Equal(t, ip, result.IpAddress)
}

// =====================================================
// gretunnels
// =====================================================

func TestGRETunnels_GetGreTunnels_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tunnelID := 12345
	path := "/zia/api/v1/greTunnels/12345"
	server.On("GET", path, common.SuccessResponse(gretunnels.GreTunnels{
		ID: tunnelID, SourceIP: "104.239.238.10", Comment: "tests-gre-tunnel",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := gretunnels.GetGreTunnels(context.Background(), service, tunnelID)
	require.NoError(t, err)
	assert.Equal(t, "104.239.238.10", result.SourceIP)
}

func TestGRETunnels_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/greTunnels"
	server.On("GET", path, common.SuccessResponse([]gretunnels.GreTunnels{
		{ID: 1, SourceIP: "104.239.238.10"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := gretunnels.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestGRETunnels_CreateGreTunnels_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/greTunnels"
	server.On("POST", path, common.SuccessResponse(gretunnels.GreTunnels{
		ID: 99999, SourceIP: "104.239.238.10",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newTunnel := &gretunnels.GreTunnels{SourceIP: "104.239.238.10", Comment: "tests-gre-tunnel"}
	result, _, err := gretunnels.CreateGreTunnels(context.Background(), service, newTunnel)
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestGRETunnels_UpdateGreTunnels_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tunnelID := 12345
	path := "/zia/api/v1/greTunnels/12345"
	server.On("PUT", path, common.SuccessResponse(gretunnels.GreTunnels{
		ID: tunnelID, Comment: "updated-gre",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	update := &gretunnels.GreTunnels{ID: tunnelID, Comment: "updated-gre"}
	result, _, err := gretunnels.UpdateGreTunnels(context.Background(), service, tunnelID, update)
	require.NoError(t, err)
	assert.Equal(t, "updated-gre", result.Comment)
}

func TestGRETunnels_DeleteGreTunnels_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/greTunnels/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = gretunnels.DeleteGreTunnels(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestGRETunnels_GetByIPAddress_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ip := "104.239.238.10"
	path := "/zia/api/v1/greTunnels"
	server.On("GET", path, common.SuccessResponse([]gretunnels.GreTunnels{
		{ID: 1, SourceIP: ip},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := gretunnels.GetByIPAddress(context.Background(), service, ip)
	require.NoError(t, err)
	assert.Equal(t, ip, result.SourceIP)
}

// =====================================================
// gretunnelinfo, greinternalipranges, ipv6_config, virtualipaddress, region
// =====================================================

func TestGRETunnelInfo_GetGRETunnelInfo_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/orgProvisioning/ipGreTunnelInfo"
	server.On("GET", path, common.SuccessResponse([]gretunnelinfo.GRETunnelInfo{
		{TunID: 1, IPaddress: "104.239.238.10", GREEnabled: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := gretunnelinfo.GetGRETunnelInfo(context.Background(), service, "104.239.238.10")
	require.NoError(t, err)
	assert.True(t, result.GREEnabled)
}

func TestGREInternalIPRanges_GetGREInternalIPRange_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/greTunnels/availableInternalIpRanges"
	server.On("GET", path, common.SuccessResponse([]greinternalipranges.GREInternalIPRange{
		{StartIPAddress: "10.0.0.0", EndIPAddress: "10.0.0.7"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := greinternalipranges.GetGREInternalIPRange(context.Background(), service, 1)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, *result, 1)
}

func TestIPv6Config_GetIPv6Config_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipv6config"
	server.On("GET", path, common.SuccessResponse(ipv6_config.IPv6Config{
		IpV6Enabled: true, DnsPrefix: "64:ff9b::/96",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipv6_config.GetIPv6Config(context.Background(), service)
	require.NoError(t, err)
	assert.True(t, result.IpV6Enabled)
}

func TestIPv6Config_GetDns64Prefix_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipv6config/dns64prefix"
	server.On("GET", path, common.SuccessResponse([]ipv6_config.IPv6ConfigPrefix{
		{ID: 1, Name: "DNS64", PrefixMask: "64:ff9b::/96"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipv6_config.GetDns64Prefix(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestVirtualIPAddress_GetZscalerVIPs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.ZscalerVIPs{
		{DataCenter: "SJC4", City: "San Jose", GREIPs: []string{"192.0.2.1"}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetZscalerVIPs(context.Background(), service, "SJC4")
	require.NoError(t, err)
	assert.Equal(t, "SJC4", result.DataCenter)
}

func TestVirtualIPAddress_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips/recommendedList"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1", DataCenter: "SJC4", CountryCode: "US"},
		{ID: 2, VirtualIp: "192.0.2.2", DataCenter: "SJC4", CountryCode: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetAll(context.Background(), service, "104.239.238.10")
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestRegion_GetDatacenterRegion_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/region/search"
	server.On("GET", path, common.SuccessResponse([]region.Regions{
		{CityName: "San Jose", StateName: "California", CountryName: "USA"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := region.GetDatacenterRegion(context.Background(), service, "San")
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIPAddress_GetByIPAddress_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/region/byIPAddress/8.8.8.8"
	server.On("GET", path, common.SuccessResponse(ip_address.ByIPAddress{
		CityName: "San Jose", StateName: "California", CountryName: "USA",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ip_address.GetByIPAddress(context.Background(), service, "8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, "San Jose", result.CityName)
}

func TestGeoCoordinates_GetByGeoCoordinates_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/region/byGeoCoordinates"
	server.On("GET", path, common.SuccessResponse(geo_coordinates.GeoCoordinates{
		CityName: "San Jose", Latitude: 37.3382, Longitude: -121.8863,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := geo_coordinates.GetByGeoCoordinates(context.Background(), service, 37.3382, -121.8863)
	require.NoError(t, err)
	assert.Equal(t, "San Jose", result.CityName)
}

func TestDatacenter_SearchByDatacenters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips/groupByDatacenter"
	server.On("GET", path, common.SuccessResponse([]datacenter.DatacenterVIPS{
		{GreVIP: []datacenter.GreVIP{{VirtualIp: "192.0.2.1", Datacenter: "SJC4"}}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := datacenter.SearchByDatacenters(context.Background(), service, ziacommon.DatacenterSearchParameters{
		RoutableIP: true,
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestExtranet_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet"
	server.On("GET", path, common.SuccessResponse([]extranet.Extranet{
		{ID: 1, Name: "Branch Extranet"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestExtranet_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet/12345"
	server.On("GET", path, common.SuccessResponse(extranet.Extranet{
		ID: 12345, Name: "Branch Extranet",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.Get(context.Background(), service, 12345)
	require.NoError(t, err)
	assert.Equal(t, "Branch Extranet", result.Name)
}

func TestDCExclusions_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dcExclusions"
	server.On("GET", path, common.SuccessResponse([]dc_exclusions.DCExclusions{
		{DcID: 1, Description: "SJC4 Exclusion"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dc_exclusions.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSubClouds_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/subclouds"
	server.On("GET", path, common.SuccessResponse([]sub_clouds.SubClouds{
		{ID: 1, Name: "US West"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sub_clouds.GetAll(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestIPv6Config_GetNat64Prefix_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipv6config/nat64prefix"
	server.On("GET", path, common.SuccessResponse([]ipv6_config.IPv6ConfigPrefix{
		{ID: 1, Name: "NAT64", PrefixMask: "64:ff9b::/96"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipv6_config.GetNat64Prefix(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestVPNCredentials_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vpnCredentials/99999"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vpncredentials.Get(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestVPNCredentials_BulkDelete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vpnCredentials/bulkDelete"
	server.On("POST", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = vpncredentials.BulkDelete(context.Background(), service, []int{1, 2})
	require.Error(t, err)
}

func TestExtranet_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet"
	server.On("POST", path, common.SuccessResponse(extranet.Extranet{
		ID: 99999, Name: "Branch Extranet",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := extranet.Create(context.Background(), service, &extranet.Extranet{Name: "Branch Extranet"})
	require.NoError(t, err)
	assert.Equal(t, 99999, result.ID)
}

func TestExtranet_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet"
	server.On("PUT", path, common.SuccessResponse(extranet.Extranet{
		ID: 12345, Name: "Updated Extranet",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.Update(context.Background(), service, 12345, &extranet.Extranet{Name: "Updated Extranet"})
	require.NoError(t, err)
	assert.Equal(t, "Updated Extranet", result.Name)
}

func TestExtranet_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet/12345"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = extranet.Delete(context.Background(), service, 12345)
	require.NoError(t, err)
}

func TestExtranet_GetExtranetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "Branch Extranet"
	path := "/zia/api/v1/extranet"
	server.On("GET", path, common.SuccessResponse([]extranet.Extranet{
		{ID: 1, Name: name},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.GetExtranetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestExtranet_GetLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet/lite"
	server.On("GET", path, common.SuccessResponse([]ziacommon.IDNameExternalID{
		{ID: 1, Name: "Branch Extranet"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.GetLite(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestExtranet_GetAll_WithOptions_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/extranet"
	server.On("GET", path, common.SuccessResponse([]extranet.Extranet{
		{ID: 1, Name: "Branch Extranet"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := extranet.GetAll(context.Background(), service, &extranet.GetAllOptions{
		Search: "Branch", OrderBy: "name", Order: "asc",
	})
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDCExclusions_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dcExclusions"
	server.On("GET", path, common.SuccessResponse([]dc_exclusions.DCExclusions{
		{DcID: 1, DcName: &ziacommon.IDNameExtensions{Name: "SJC4"}},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dc_exclusions.GetByName(context.Background(), service, "SJC4")
	require.NoError(t, err)
	assert.Equal(t, 1, result.DcID)
}

func TestDCExclusions_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dcExclusions"
	server.On("POST", path, common.SuccessResponse([]dc_exclusions.DCExclusions{
		{DcID: 1, Description: "Maintenance window"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := dc_exclusions.Create(context.Background(), service, &dc_exclusions.DCExclusions{
		DcID: 1, Description: "Maintenance window",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, result.DcID)
}

func TestDCExclusions_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dcExclusions"
	server.On("PUT", path, common.SuccessResponse(dc_exclusions.DCExclusions{
		DcID: 1, Description: "Updated exclusion",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := dc_exclusions.Update(context.Background(), service, &dc_exclusions.DCExclusions{
		DcID: 1, Description: "Updated exclusion",
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated exclusion", result.Description)
}

func TestDCExclusions_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/dcExclusions/1"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dc_exclusions.Delete(context.Background(), service, 1)
	require.NoError(t, err)
}

func TestDCExclusions_GetDatacenters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/datacenters"
	server.On("GET", path, common.SuccessResponse([]dc_exclusions.Datacenter{
		{ID: 1, Name: "SJC4", City: "San Jose"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dc_exclusions.GetDatacenters(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestSubClouds_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/subclouds/isLastDcInCountry//1"
	server.On("GET", path, common.SuccessResponse(sub_clouds.SubCloudCountryDCExclusionInfo{
		ID: 1, Country: "US",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sub_clouds.Get(context.Background(), service, 1)
	require.NoError(t, err)
	assert.Equal(t, "US", result.Country)
}

func TestSubClouds_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "US West"
	path := "/zia/api/v1/subclouds"
	server.On("GET", path, common.SuccessResponse([]sub_clouds.SubClouds{
		{ID: 1, Name: name},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := sub_clouds.GetByName(context.Background(), service, name)
	require.NoError(t, err)
	assert.Equal(t, name, result.Name)
}

func TestSubClouds_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/subclouds/1"
	server.On("PUT", path, common.SuccessResponse(sub_clouds.SubClouds{
		ID: 1, Name: "US West Updated",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := sub_clouds.Update(context.Background(), service, 1, &sub_clouds.SubClouds{Name: "US West Updated"})
	require.NoError(t, err)
	assert.Equal(t, "US West Updated", result.Name)
}

func TestSubClouds_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/subclouds/1"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = sub_clouds.Delete(context.Background(), service, 1)
	require.NoError(t, err)
}

func TestVirtualIPAddress_GetZSGREVirtualIPList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips/recommendedList"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1", CountryCode: "US"},
		{ID: 2, VirtualIp: "192.0.2.2", CountryCode: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetZSGREVirtualIPList(context.Background(), service, "203.0.113.1", 2)
	require.NoError(t, err)
	assert.Len(t, *result, 2)
}

func TestVirtualIPAddress_GetPairZSGREVirtualIPsWithinCountry_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips/recommendedList"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1", CountryCode: "US"},
		{ID: 2, VirtualIp: "192.0.2.2", CountryCode: "US"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetPairZSGREVirtualIPsWithinCountry(context.Background(), service, "203.0.113.1", "US")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(*result), 2)
}

func TestVirtualIPAddress_GetVIPRecommendedList_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips/recommendedList"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1"},
		{ID: 2, VirtualIp: "192.0.2.2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetVIPRecommendedList(context.Background(), service,
		virtualipaddress.WithRoutableIP(true),
		virtualipaddress.WithSourceIP("203.0.113.1"),
	)
	require.NoError(t, err)
	assert.Len(t, *result, 2)
}

func TestVirtualIPAddress_GetAllSourceIPs_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	staticPath := "/zia/api/v1/staticIP"
	vipPath := "/zia/api/v1/vips/recommendedList"

	server.On("GET", staticPath, common.SuccessResponse([]staticips.StaticIP{
		{ID: 1, IpAddress: "203.0.113.1"},
	}))
	server.On("GET", vipPath, common.SuccessResponse([]virtualipaddress.GREVirtualIPList{
		{ID: 1, VirtualIp: "192.0.2.1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetAllSourceIPs(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestVirtualIPAddress_GetZscalerVIPs_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/vips"
	server.On("GET", path, common.SuccessResponse([]virtualipaddress.ZscalerVIPs{
		{DataCenter: "SJC4"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := virtualipaddress.GetZscalerVIPs(context.Background(), service, "UNKNOWN")
	require.Error(t, err)
	assert.Nil(t, result)
}
