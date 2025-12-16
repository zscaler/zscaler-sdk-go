// Package services provides unit tests for ZTW services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

// =====================================================
// IP Destination Groups SDK Function Tests
// =====================================================

func TestIPDestinationGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/ztw/api/v1/ipDestinationGroups/12345"

	server.On("GET", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:   groupID,
		Name: "External-Servers",
		Type: "DSTN_IP",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
}

func TestIPDestinationGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "External-Servers"
	path := "/ztw/api/v1/ipDestinationGroups"

	server.On("GET", path, common.SuccessResponse([]ipdestinationgroups.IPDestinationGroups{
		{ID: 1, Name: "Other Group", Type: "DSTN_FQDN"},
		{ID: 2, Name: groupName, Type: "DSTN_IP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestIPDestinationGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ipDestinationGroups"

	server.On("GET", path, common.SuccessResponse([]ipdestinationgroups.IPDestinationGroups{
		{ID: 1, Name: "Group 1", Type: "DSTN_IP"},
		{ID: 2, Name: "Group 2", Type: "DSTN_FQDN"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIPDestinationGroups_GetAllLite_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ipDestinationGroups/lite"

	server.On("GET", path, common.SuccessResponse([]ipdestinationgroups.IPDestinationGroups{
		{ID: 1, Name: "Group 1"},
		{ID: 2, Name: "Group 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.GetAllLite(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIPDestinationGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/ipDestinationGroups"

	server.On("POST", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:   99999,
		Name: "New Group",
		Type: "DSTN_IP",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newGroup := &ipdestinationgroups.IPDestinationGroups{
		Name: "New Group",
		Type: "DSTN_IP",
	}

	result, err := ipdestinationgroups.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestIPDestinationGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/ztw/api/v1/ipDestinationGroups/12345"

	server.On("PUT", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:   groupID,
		Name: "Updated Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGroup := &ipdestinationgroups.IPDestinationGroups{
		ID:   groupID,
		Name: "Updated Group",
	}

	result, _, err := ipdestinationgroups.Update(context.Background(), service, groupID, updateGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Group", result.Name)
}

func TestIPDestinationGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/ztw/api/v1/ipDestinationGroups/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = ipdestinationgroups.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
}

// =====================================================
// Network Services SDK Function Tests
// =====================================================

func TestNetworkServices_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/ztw/api/v1/networkServices/12345"

	server.On("GET", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:   serviceID,
		Name: "Custom-HTTP-Service",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.Get(context.Background(), service, serviceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceID, result.ID)
}

func TestNetworkServices_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceName := "Custom-HTTP-Service"
	path := "/ztw/api/v1/networkServices"

	server.On("GET", path, common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "Other Service", Type: "PREDEFINED"},
		{ID: 2, Name: serviceName, Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetByName(context.Background(), service, serviceName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceName, result.Name)
}

func TestNetworkServices_GetAllNetworkServices_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/networkServices"

	server.On("GET", path, common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "Service 1", Type: "CUSTOM"},
		{ID: 2, Name: "Service 2", Type: "PREDEFINED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestNetworkServices_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/ztw/api/v1/networkServices"

	server.On("POST", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:   99999,
		Name: "New Service",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newSvc := &networkservices.NetworkServices{
		Name: "New Service",
		Type: "CUSTOM",
	}

	result, err := networkservices.Create(context.Background(), service, newSvc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestNetworkServices_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/ztw/api/v1/networkServices/12345"

	server.On("PUT", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:   serviceID,
		Name: "Updated Service",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateSvc := &networkservices.NetworkServices{
		ID:   serviceID,
		Name: "Updated Service",
	}

	result, _, err := networkservices.Update(context.Background(), service, serviceID, updateSvc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Service", result.Name)
}

func TestNetworkServices_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/ztw/api/v1/networkServices/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkservices.Delete(context.Background(), service, serviceID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

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

