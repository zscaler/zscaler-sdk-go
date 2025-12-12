// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/trustednetwork"
)

func TestTrustedNetwork_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	networkID := "network-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/network/" + networkID

	server.On("GET", path, common.SuccessResponse(trustednetwork.TrustedNetwork{
		ID:        networkID,
		Name:      "Test Network",
		NetworkID: "net-001",
		Domain:    "example.com",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := trustednetwork.Get(context.Background(), service, networkID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, networkID, result.ID)
	assert.Equal(t, "Test Network", result.Name)
}

func TestTrustedNetwork_GetByNetID_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	netID := "net-12345"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/network"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []trustednetwork.TrustedNetwork{
			{ID: "tn-001", Name: "Other Network", NetworkID: "net-001"},
			{ID: "tn-002", Name: "Target Network", NetworkID: netID},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := trustednetwork.GetByNetID(context.Background(), service, netID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tn-002", result.ID)
	assert.Equal(t, netID, result.NetworkID)
}

func TestTrustedNetwork_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	networkName := "Production Network"
	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/network"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []trustednetwork.TrustedNetwork{
			{ID: "tn-001", Name: "Other Network"},
			{ID: "tn-002", Name: networkName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := trustednetwork.GetByName(context.Background(), service, networkName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tn-002", result.ID)
	assert.Equal(t, networkName, result.Name)
}

func TestTrustedNetwork_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/network"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []trustednetwork.TrustedNetwork{
			{ID: "tn-001", Name: "Network 1"},
			{ID: "tn-002", Name: "Network 2"},
			{ID: "tn-003", Name: "Network 3"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := trustednetwork.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestTrustedNetwork_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v2/admin/customers/" + testCustomerID + "/network"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []trustednetwork.TrustedNetwork{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := trustednetwork.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no trusted network named")
}

func TestTrustedNetwork_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	networkID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/network/" + networkID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := trustednetwork.Get(context.Background(), service, networkID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
