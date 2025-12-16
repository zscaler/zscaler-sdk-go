// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
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
