package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/dns_gateways"
)

const dnsGatewaysPath = "/zia/api/v1/dnsGateways"
const dnsGatewaysLitePath = "/zia/api/v1/dnsGateways/lite"

func sampleDNSGateway() dns_gateways.DNSGateways {
	return dns_gateways.DNSGateways{
		Name:              "tests-dns-gateway",
		PrimaryIpOrFqdn:   "8.8.8.8",
		SecondaryIpOrFqdn: "4.4.4.4",
		FailureBehavior:   "FAIL_RET_ERR",
		Protocols:         []string{"TCP", "UDP", "DOH"},
		PrimaryPorts:      []int{53, 53, 443},
		SecondaryPorts:    []int{53, 53, 443},
	}
}

func TestDNSGateways_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwID := 100
	gw := sampleDNSGateway()
	gw.ID = gwID

	server.On("GET", dnsGatewaysPath+"/100", common.SuccessResponse(gw))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_gateways.Get(context.Background(), service, gwID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gwID, result.ID)
	assert.Equal(t, "8.8.8.8", result.PrimaryIpOrFqdn)
	assert.Equal(t, "FAIL_RET_ERR", result.FailureBehavior)
}

func TestDNSGateways_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwName := "tests-dns-gateway"
	server.On("GET", dnsGatewaysPath, common.SuccessResponse([]dns_gateways.DNSGateways{
		{ID: 1, Name: "Other Gateway"},
		func() dns_gateways.DNSGateways {
			g := sampleDNSGateway()
			g.ID = 100
			return g
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_gateways.GetByName(context.Background(), service, gwName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gwName, result.Name)
}

func TestDNSGateways_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleDNSGateway()
	created.ID = 99999

	server.On("POST", dnsGatewaysPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleDNSGateway()
	result, _, err := dns_gateways.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestDNSGateways_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwID := 100
	updated := sampleDNSGateway()
	updated.ID = gwID
	updated.Name = "updated-dns-gateway"

	server.On("PUT", dnsGatewaysPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := dns_gateways.Update(context.Background(), service, gwID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-dns-gateway", result.Name)
}

func TestDNSGateways_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", dnsGatewaysPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = dns_gateways.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestDNSGateways_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dnsGatewaysPath, common.SuccessResponse([]dns_gateways.DNSGateways{
		sampleDNSGateway(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_gateways.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDNSGateways_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dnsGatewaysLitePath, common.SuccessResponse([]dns_gateways.DNSGateways{
		{ID: 100, Name: "tests-dns-gateway"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dns_gateways.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDNSGateways_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		gw := sampleDNSGateway()
		gw.ID = 100

		data, err := json.Marshal(gw)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"primaryIpOrFqdn":"8.8.8.8"`)
		assert.Contains(t, string(data), `"failureBehavior":"FAIL_RET_ERR"`)
	})
}
