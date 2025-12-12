// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestNetworkServices_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/zia/api/v1/networkServices/12345"

	server.On("GET", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:          serviceID,
		Name:        "HTTPS",
		Description: "HTTPS Traffic",
		Type:        "STANDARD",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.Get(context.Background(), service, serviceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceID, result.ID)
	assert.Equal(t, "HTTPS", result.Name)
}

func TestNetworkServices_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/networkServices"

	server.On("GET", path, common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "HTTP", Type: "STANDARD"},
		{ID: 2, Name: "HTTPS", Type: "STANDARD"},
		{ID: 3, Name: "Custom Service", Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestNetworkServices_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/networkServices"

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
	path := "/zia/api/v1/networkServices/12345"

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
	path := "/zia/api/v1/networkServices/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkservices.Delete(context.Background(), service, serviceID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestNetworkServices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkServices JSON marshaling", func(t *testing.T) {
		svc := networkservices.NetworkServices{
			ID:           12345,
			Name:         "Custom HTTPS",
			Description:  "Custom HTTPS service",
			Type:         "CUSTOM",
			SrcTCPPorts:  []networkservices.NetworkPorts{{Start: 443, End: 443}},
			DestTCPPorts: []networkservices.NetworkPorts{{Start: 443, End: 443}},
		}

		data, err := json.Marshal(svc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom HTTPS"`)
	})

	t.Run("NetworkServices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "DNS",
			"type": "STANDARD",
			"description": "DNS service",
			"destTcpPorts": [{"start": 53, "end": 53}],
			"destUdpPorts": [{"start": 53, "end": 53}]
		}`

		var svc networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonData), &svc)
		require.NoError(t, err)

		assert.Equal(t, 54321, svc.ID)
		assert.Equal(t, "DNS", svc.Name)
	})
}
