package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxies"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/proxy_gateways"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/forwarding_control_policy/zpa_gateways"
)

const forwardingRulesPath = "/zia/api/v1/forwardingRules"
const proxiesPath = "/zia/api/v1/proxies"
const proxiesLitePath = "/zia/api/v1/proxies/lite"
const dedicatedIPGWLitePath = "/zia/api/v1/dedicatedIPGateways/lite"
const proxyGatewaysPath = "/zia/api/v1/proxyGateways"
const proxyGatewaysLitePath = "/zia/api/v1/proxyGateways/lite"
const zpaGatewaysPath = "/zia/api/v1/zpaGateways"

func sampleForwardingRule() forwarding_rules.ForwardingRules {
	return forwarding_rules.ForwardingRules{
		Name:          "tests-forwarding-rule",
		Description:   "tests-forwarding-rule",
		Order:         1,
		Rank:          7,
		State:         "ENABLED",
		Type:          "FORWARDING",
		ForwardMethod: "DIRECT",
		DestCountries: []string{"COUNTRY_CA", "COUNTRY_US", "COUNTRY_MX", "COUNTRY_AU", "COUNTRY_GB"},
		SrcIps:        []string{"192.168.100.10"},
	}
}

func TestForwardingRules_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	rule := sampleForwardingRule()
	rule.ID = ruleID

	server.On("GET", forwardingRulesPath+"/12345", common.SuccessResponse(rule))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.Get(context.Background(), service, ruleID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleID, result.ID)
	assert.Equal(t, "DIRECT", result.ForwardMethod)
}

func TestForwardingRules_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleName := "tests-forwarding-rule"
	server.On("GET", forwardingRulesPath, common.SuccessResponse([]forwarding_rules.ForwardingRules{
		{ID: 1, Name: "Other Rule"},
		func() forwarding_rules.ForwardingRules {
			r := sampleForwardingRule()
			r.ID = 2
			return r
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.GetByName(context.Background(), service, ruleName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, ruleName, result.Name)
}

func TestForwardingRules_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleForwardingRule()
	created.ID = 99999

	server.On("POST", forwardingRulesPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleForwardingRule()
	result, err := forwarding_rules.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestForwardingRules_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	ruleID := 12345
	updated := sampleForwardingRule()
	updated.ID = ruleID
	updated.Name = "updated-forwarding-rule"

	server.On("PUT", forwardingRulesPath+"/12345", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.Update(context.Background(), service, ruleID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-forwarding-rule", result.Name)
}

func TestForwardingRules_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", forwardingRulesPath+"/12345", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = forwarding_rules.Delete(context.Background(), service, 12345)

	require.NoError(t, err)
}

func TestForwardingRules_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", forwardingRulesPath, common.SuccessResponse([]forwarding_rules.ForwardingRules{
		sampleForwardingRule(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := forwarding_rules.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func sampleProxy() proxies.Proxies {
	return proxies.Proxies{
		Name:                  "tests-proxy",
		Description:           "tests-proxy",
		Type:                  "PROXYCHAIN",
		Address:               "192.168.1.1",
		Port:                  5000,
		InsertXauHeader:       true,
		Base64EncodeXauHeader: true,
	}
}

func TestProxies_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	proxyID := 100
	proxy := sampleProxy()
	proxy.ID = proxyID

	server.On("GET", proxiesPath+"/100", common.SuccessResponse(proxy))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxies.Get(context.Background(), service, proxyID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, proxyID, result.ID)
	assert.Equal(t, "PROXYCHAIN", result.Type)
}

func TestProxies_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	proxyName := "tests-proxy"
	server.On("GET", proxiesPath, common.SuccessResponse([]proxies.Proxies{
		func() proxies.Proxies {
			p := sampleProxy()
			p.ID = 100
			return p
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxies.GetByName(context.Background(), service, proxyName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, proxyName, result.Name)
}

func TestProxies_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleProxy()
	created.ID = 99999

	server.On("POST", proxiesPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleProxy()
	result, _, err := proxies.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestProxies_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	proxyID := 100
	updated := sampleProxy()
	updated.ID = proxyID
	updated.Name = "updated-proxy"

	server.On("PUT", proxiesPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := proxies.Update(context.Background(), service, proxyID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-proxy", result.Name)
}

func TestProxies_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", proxiesPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = proxies.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestProxies_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", proxiesPath, common.SuccessResponse([]proxies.Proxies{
		sampleProxy(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxies.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestProxies_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", proxiesLitePath, common.SuccessResponse([]proxies.Proxies{
		{ID: 100, Name: "tests-proxy", Type: "PROXYCHAIN"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxies.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestProxies_GetDedicatedIPGWLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", dedicatedIPGWLitePath, common.SuccessResponse([]proxies.DedicatedIPGateways{
		{Id: 1, Name: "Dedicated GW 1"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxies.GetDedicatedIPGWLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestProxyGateways_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwName := "Proxy Gateway 1"
	server.On("GET", proxyGatewaysPath, common.SuccessResponse([]proxy_gateways.ProxyGateways{
		{ID: 1, Name: gwName, Type: "PROXYCHAIN", FailClosed: true},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxy_gateways.GetByName(context.Background(), service, gwName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gwName, result.Name)
}

func TestProxyGateways_GetLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", proxyGatewaysLitePath, common.SuccessResponse([]proxy_gateways.ProxyGateways{
		{ID: 1, Name: "Proxy Gateway 1", Type: "PROXYCHAIN"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxy_gateways.GetLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestProxyGateways_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", proxyGatewaysPath, common.SuccessResponse([]proxy_gateways.ProxyGateways{
		{ID: 1, Name: "Proxy Gateway 1", Type: "PROXYCHAIN", FailClosed: false},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := proxy_gateways.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func sampleZPAGateway() zpa_gateways.ZPAGateways {
	return zpa_gateways.ZPAGateways{
		Name:        "tests-zpa-gateway",
		Description: "tests-zpa-gateway",
		Type:        "ZPA",
		ZPAServerGroup: zpa_gateways.ZPAServerGroup{
			ExternalID: "zpa-sg-123",
			Name:       "Server Group 1",
		},
		ZPAAppSegments: []zpa_gateways.ZPAAppSegments{
			{ExternalID: "zpa-app-1", Name: "App Segment 1"},
		},
	}
}

func TestZPAGateways_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwID := 100
	gw := sampleZPAGateway()
	gw.ID = gwID

	server.On("GET", zpaGatewaysPath+"/100", common.SuccessResponse(gw))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zpa_gateways.Get(context.Background(), service, gwID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gwID, result.ID)
	assert.Equal(t, "ZPA", result.Type)
}

func TestZPAGateways_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwName := "tests-zpa-gateway"
	server.On("GET", zpaGatewaysPath, common.SuccessResponse([]zpa_gateways.ZPAGateways{
		func() zpa_gateways.ZPAGateways {
			g := sampleZPAGateway()
			g.ID = 100
			return g
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zpa_gateways.GetByName(context.Background(), service, gwName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, gwName, result.Name)
}

func TestZPAGateways_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleZPAGateway()
	created.ID = 99999

	server.On("POST", zpaGatewaysPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleZPAGateway()
	result, err := zpa_gateways.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestZPAGateways_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	gwID := 100
	updated := sampleZPAGateway()
	updated.ID = gwID
	updated.Name = "updated-zpa-gateway"

	server.On("PUT", zpaGatewaysPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zpa_gateways.Update(context.Background(), service, gwID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-zpa-gateway", result.Name)
}

func TestZPAGateways_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", zpaGatewaysPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = zpa_gateways.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestZPAGateways_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", zpaGatewaysPath, common.SuccessResponse([]zpa_gateways.ZPAGateways{
		sampleZPAGateway(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := zpa_gateways.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestForwardingControl_Structure(t *testing.T) {
	t.Parallel()

	t.Run("ForwardingRules JSON marshaling", func(t *testing.T) {
		rule := sampleForwardingRule()
		rule.ID = 12345

		data, err := json.Marshal(rule)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"forwardMethod":"DIRECT"`)
		assert.Contains(t, string(data), `"COUNTRY_US"`)
	})

	t.Run("Proxies JSON marshaling", func(t *testing.T) {
		proxy := sampleProxy()
		proxy.ID = 100

		data, err := json.Marshal(proxy)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"insertXauHeader":true`)
	})
}
