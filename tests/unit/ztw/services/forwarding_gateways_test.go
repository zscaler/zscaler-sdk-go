// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/dns_forwarding_gateway"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/forwarding_gateways/zia_forwarding_gateway"
)

// =====================================================
// DNS Forwarding Gateway SDK Function Tests
// =====================================================

func TestDNSForwardingGateway_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/dnsGateways/12345"

	server.On("GET", path, common.SuccessResponse(dns_forwarding_gateway.DNSGateway{
		ID:        gatewayID,
		Name:      "DNS-Gateway-1",
		Type:      "ZIA",
		PrimaryIP: "10.0.0.1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := dns_forwarding_gateway.Get(context.Background(), service, gatewayID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gatewayID, result.ID)
	assert.Equal(t, "DNS-Gateway-1", result.Name)
}

func TestDNSForwardingGateway_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayName := "DNS-Gateway-1"
	path := "/ztw/api/v1/dnsGateways"

	server.On("GET", path, common.SuccessResponse([]dns_forwarding_gateway.DNSGateway{
		{ID: 1, Name: "Other Gateway", Type: "ECSELF"},
		{ID: 2, Name: gatewayName, Type: "ZIA"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_forwarding_gateway.GetByName(context.Background(), service, gatewayName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gatewayName, result.Name)
}

func TestDNSForwardingGateway_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/dnsGateways"

	server.On("GET", path, common.SuccessResponse([]dns_forwarding_gateway.DNSGateway{
		{ID: 1, Name: "Gateway 1", Type: "ZIA"},
		{ID: 2, Name: "Gateway 2", Type: "ECSELF"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_forwarding_gateway.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDNSForwardingGateway_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/dnsGateways/lite"

	server.On("GET", path, common.SuccessResponse([]dns_forwarding_gateway.DNSGateway{
		{ID: 1, Name: "Gateway 1"},
		{ID: 2, Name: "Gateway 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_forwarding_gateway.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDNSForwardingGateway_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/dnsGateways"

	server.On("POST", path, common.SuccessResponse(dns_forwarding_gateway.DNSGateway{
		ID:        99999,
		Name:      "New Gateway",
		Type:      "ZIA",
		PrimaryIP: "10.0.0.1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newGW := &dns_forwarding_gateway.DNSGateway{
		Name:      "New Gateway",
		Type:      "ZIA",
		PrimaryIP: "10.0.0.1",
	}

	result, _, err := dns_forwarding_gateway.Create(context.Background(), service, newGW)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDNSForwardingGateway_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/dnsGateways/12345"

	server.On("PUT", path, common.SuccessResponse(dns_forwarding_gateway.DNSGateway{
		ID:   gatewayID,
		Name: "Updated Gateway",
		Type: "ZIA",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGW := &dns_forwarding_gateway.DNSGateway{
		ID:   gatewayID,
		Name: "Updated Gateway",
		Type: "ZIA",
	}

	result, _, err := dns_forwarding_gateway.Update(context.Background(), service, gatewayID, updateGW)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Gateway", result.Name)
}

func TestDNSForwardingGateway_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/dnsGateways/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dns_forwarding_gateway.Delete(context.Background(), service, gatewayID)

	require.NoError(t, err)
}

// =====================================================
// ZIA Forwarding Gateway SDK Function Tests
// =====================================================

func TestZIAForwardingGateway_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/gateways/12345"

	server.On("GET", path, common.SuccessResponse(zia_forwarding_gateway.ECGateway{
		ID:          gatewayID,
		Name:        "ZIA-Gateway-1",
		PrimaryType: "AUTO",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := zia_forwarding_gateway.Get(context.Background(), service, gatewayID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gatewayID, result.ID)
}

func TestZIAForwardingGateway_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayName := "ZIA-Gateway-1"
	path := "/ztw/api/v1/gateways"

	server.On("GET", path, common.SuccessResponse([]zia_forwarding_gateway.ECGateway{
		{ID: 1, Name: "Other Gateway"},
		{ID: 2, Name: gatewayName},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zia_forwarding_gateway.GetByName(context.Background(), service, gatewayName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gatewayName, result.Name)
}

func TestZIAForwardingGateway_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/gateways"

	server.On("GET", path, common.SuccessResponse([]zia_forwarding_gateway.ECGateway{
		{ID: 1, Name: "Gateway 1"},
		{ID: 2, Name: "Gateway 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zia_forwarding_gateway.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestZIAForwardingGateway_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/gateways"

	server.On("POST", path, common.SuccessResponse(zia_forwarding_gateway.ECGateway{
		ID:          99999,
		Name:        "New ZIA Gateway",
		PrimaryType: "AUTO",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newGW := &zia_forwarding_gateway.ECGateway{
		Name:        "New ZIA Gateway",
		PrimaryType: "AUTO",
	}

	result, _, err := zia_forwarding_gateway.Create(context.Background(), service, newGW)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestZIAForwardingGateway_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/gateways/12345"

	server.On("PUT", path, common.SuccessResponse(zia_forwarding_gateway.ECGateway{
		ID:   gatewayID,
		Name: "Updated ZIA Gateway",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGW := &zia_forwarding_gateway.ECGateway{
		ID:   gatewayID,
		Name: "Updated ZIA Gateway",
	}

	result, _, err := zia_forwarding_gateway.Update(context.Background(), service, gatewayID, updateGW)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated ZIA Gateway", result.Name)
}

func TestZIAForwardingGateway_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gatewayID := 12345
	path := "/ztw/api/v1/gateways/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zia_forwarding_gateway.Delete(context.Background(), service, gatewayID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

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

