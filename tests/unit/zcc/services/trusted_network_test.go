// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/trusted_network"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestTrustedNetwork_GetMultiple_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webTrustedNetwork/listByCompany"

	server.On("GET", path, common.SuccessResponse(trusted_network.TrustedNetworksResponse{
		TotalCount: 2,
		TrustedNetworkContracts: []trusted_network.TrustedNetwork{
			{ID: "tn-001", NetworkName: "Corporate HQ", Active: true, CompanyID: "company-123"},
			{ID: "tn-002", NetworkName: "Branch Office", Active: true, CompanyID: "company-123"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := trusted_network.GetMultipleTrustedNetworks(context.Background(), service, "", "", nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.TotalCount)
	assert.Len(t, result.TrustedNetworkContracts, 2)
}

func TestTrustedNetwork_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webTrustedNetwork/listByCompany"

	server.On("GET", path, common.SuccessResponse(trusted_network.TrustedNetworksResponse{
		TotalCount: 2,
		TrustedNetworkContracts: []trusted_network.TrustedNetwork{
			{ID: "tn-001", NetworkName: "Corporate HQ", Active: true},
			{ID: "tn-002", NetworkName: "Branch Office", Active: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := trusted_network.GetTrustedNetworkByName(context.Background(), service, "Corporate HQ")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tn-001", result.ID)
	assert.Equal(t, "Corporate HQ", result.NetworkName)
}

func TestTrustedNetwork_GetByID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webTrustedNetwork/listByCompany"

	server.On("GET", path, common.SuccessResponse(trusted_network.TrustedNetworksResponse{
		TotalCount: 2,
		TrustedNetworkContracts: []trusted_network.TrustedNetwork{
			{ID: "tn-001", NetworkName: "Corporate HQ", Active: true},
			{ID: "tn-002", NetworkName: "Branch Office", Active: false},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := trusted_network.GetTrustedNetworkByID(context.Background(), service, "tn-002")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tn-002", result.ID)
}

func TestTrustedNetwork_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Mock the create endpoint
	createPath := "/zcc/papi/public/v1/webTrustedNetwork/create"
	server.On("POST", createPath, common.SuccessResponse(trusted_network.TrustedNetwork{
		ID:          "tn-new",
		NetworkName: "New Network",
		Active:      true,
		CompanyID:   "company-123",
	}))

	// Mock the listByCompany endpoint (called after create to verify)
	listPath := "/zcc/papi/public/v1/webTrustedNetwork/listByCompany"
	server.On("GET", listPath, common.SuccessResponse(trusted_network.TrustedNetworksResponse{
		TotalCount: 1,
		TrustedNetworkContracts: []trusted_network.TrustedNetwork{
			{ID: "tn-new", NetworkName: "New Network", Active: true, CompanyID: "company-123"},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newNetwork := &trusted_network.TrustedNetwork{
		NetworkName:    "New Network",
		Active:         true,
		TrustedSubnets: "192.168.1.0/24",
	}

	result, _, err := trusted_network.CreateTrustedNetwork(context.Background(), service, newNetwork)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tn-new", result.ID)
}

func TestTrustedNetwork_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// Mock the edit endpoint
	editPath := "/zcc/papi/public/v1/webTrustedNetwork/edit"
	server.On("PUT", editPath, common.SuccessResponse(trusted_network.TrustedNetwork{
		ID:          "tn-001",
		NetworkName: "Updated Network",
		Active:      true,
	}))

	// Mock the listByCompany endpoint (called after update to verify)
	listPath := "/zcc/papi/public/v1/webTrustedNetwork/listByCompany"
	server.On("GET", listPath, common.SuccessResponse(trusted_network.TrustedNetworksResponse{
		TotalCount: 1,
		TrustedNetworkContracts: []trusted_network.TrustedNetwork{
			{ID: "tn-001", NetworkName: "Updated Network", Active: true},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateNetwork := &trusted_network.TrustedNetwork{
		ID:          "tn-001",
		NetworkName: "Updated Network",
		Active:      true,
	}

	result, _, err := trusted_network.UpdateTrustedNetwork(context.Background(), service, updateNetwork)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Network", result.NetworkName)
}

func TestTrustedNetwork_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zcc/papi/public/v1/webTrustedNetwork/tn-001/delete"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = trusted_network.DeleteTrustedNetwork(context.Background(), service, "tn-001")

	require.NoError(t, err)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestTrustedNetwork_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TrustedNetwork JSON marshaling", func(t *testing.T) {
		tn := trusted_network.TrustedNetwork{
			ID:                 "tn-123",
			NetworkName:        "Corporate Network",
			Active:             true,
			CompanyID:          "company-456",
			ConditionType:      1,
			DnsServers:         "8.8.8.8,8.8.4.4",
			DnsSearchDomains:   "corp.example.com",
			TrustedSubnets:     "192.168.0.0/16,10.0.0.0/8",
			TrustedGateways:    "192.168.1.1",
			TrustedDhcpServers: "192.168.1.2",
			TrustedEgressIps:   "203.0.113.1",
			Ssids:              "CorpWiFi,GuestWiFi",
			Hostnames:          "server1.corp.com",
			CreatedBy:          "admin@example.com",
			EditedBy:           "admin@example.com",
		}

		data, err := json.Marshal(tn)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"tn-123"`)
		assert.Contains(t, string(data), `"networkName":"Corporate Network"`)
		assert.Contains(t, string(data), `"active":true`)
		assert.Contains(t, string(data), `"trustedSubnets":"192.168.0.0/16,10.0.0.0/8"`)
	})

	t.Run("TrustedNetwork JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "tn-456",
			"networkName": "Branch Office",
			"active": false,
			"companyId": "company-789",
			"conditionType": 2,
			"dnsServers": "1.1.1.1",
			"trustedSubnets": "172.16.0.0/12",
			"ssids": "BranchWiFi",
			"createdBy": "admin@example.com"
		}`

		var tn trusted_network.TrustedNetwork
		err := json.Unmarshal([]byte(jsonData), &tn)
		require.NoError(t, err)

		assert.Equal(t, "tn-456", tn.ID)
		assert.Equal(t, "Branch Office", tn.NetworkName)
		assert.False(t, tn.Active)
		assert.Equal(t, "1.1.1.1", tn.DnsServers)
	})

	t.Run("TrustedNetworksResponse JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"totalCount": 3,
			"trustedNetworkContracts": [
				{"id": "tn-1", "networkName": "Network 1", "active": true},
				{"id": "tn-2", "networkName": "Network 2", "active": true},
				{"id": "tn-3", "networkName": "Network 3", "active": false}
			]
		}`

		var response trusted_network.TrustedNetworksResponse
		err := json.Unmarshal([]byte(jsonData), &response)
		require.NoError(t, err)

		assert.Equal(t, 3, response.TotalCount)
		assert.Len(t, response.TrustedNetworkContracts, 3)
	})
}
