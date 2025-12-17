// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestNSSServers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := "/zia/api/v1/nssServers/12345"

	server.On("GET", path, common.SuccessResponse(nss_servers.NSSServers{
		ID:     nssID,
		Name:   "NSS Server 1",
		Status: "UP",
		Type:   "NSS_FOR_WEB",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.Get(context.Background(), service, nssID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, nssID, result.ID)
	assert.Equal(t, "NSS Server 1", result.Name)
}

func TestNSSServers_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverName := "NSS Server 1"
	path := "/zia/api/v1/nssServers"

	server.On("GET", path, common.SuccessResponse([]nss_servers.NSSServers{
		{ID: 1, Name: "Other Server", Status: "UP"},
		{ID: 2, Name: serverName, Status: "UP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetByName(context.Background(), service, serverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, serverName, result.Name)
}

func TestNSSServers_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/nssServers"

	server.On("POST", path, common.SuccessResponse(nss_servers.NSSServers{
		ID:   100,
		Name: "New NSS Server",
		Type: "NSS_FOR_WEB",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newServer := &nss_servers.NSSServers{
		Name: "New NSS Server",
		Type: "NSS_FOR_WEB",
	}

	result, err := nss_servers.Create(context.Background(), service, newServer)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
}

func TestNSSServers_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := "/zia/api/v1/nssServers/12345"

	server.On("PUT", path, common.SuccessResponse(nss_servers.NSSServers{
		ID:   nssID,
		Name: "Updated NSS Server",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateServer := &nss_servers.NSSServers{
		ID:   nssID,
		Name: "Updated NSS Server",
	}

	result, err := nss_servers.Update(context.Background(), service, nssID, updateServer)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated NSS Server", result.Name)
}

func TestNSSServers_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := "/zia/api/v1/nssServers/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nss_servers.Delete(context.Background(), service, nssID)

	require.NoError(t, err)
}

func TestNSSServers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/nssServers"

	server.On("GET", path, common.SuccessResponse([]nss_servers.NSSServers{
		{ID: 1, Name: "Server 1", Status: "UP"},
		{ID: 2, Name: "Server 2", Status: "UP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

// =====================================================
// Structure Tests
// =====================================================

func TestNSSServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NSSServers JSON marshaling", func(t *testing.T) {
		srv := nss_servers.NSSServers{
			ID:        12345,
			Name:      "NSS Server 1",
			Status:    "UP",
			State:     "ENABLED",
			Type:      "NSS_FOR_WEB",
			IcapSvrId: 100,
		}

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"NSS_FOR_WEB"`)
		assert.Contains(t, string(data), `"status":"UP"`)
	})

	t.Run("NSSServers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "NSS Server 2",
			"status": "DOWN",
			"state": "DISABLED",
			"type": "NSS_FOR_FIREWALL",
			"icapSvrId": 200
		}`

		var srv nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonData), &srv)
		require.NoError(t, err)

		assert.Equal(t, 54321, srv.ID)
		assert.Equal(t, "DOWN", srv.Status)
		assert.Equal(t, "NSS_FOR_FIREWALL", srv.Type)
	})

	t.Run("NSSServers with VZEN type", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "VZEN Server",
			"status": "UP",
			"state": "ENABLED",
			"type": "VZEN"
		}`

		var srv nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonData), &srv)
		require.NoError(t, err)

		assert.Equal(t, "VZEN", srv.Type)
	})
}

func TestNSSServers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse NSS servers list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Server 1", "status": "UP", "type": "NSS_FOR_WEB"},
			{"id": 2, "name": "Server 2", "status": "UP", "type": "NSS_FOR_FIREWALL"},
			{"id": 3, "name": "Server 3", "status": "DOWN", "type": "VZEN"}
		]`

		var servers []nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonResponse), &servers)
		require.NoError(t, err)

		assert.Len(t, servers, 3)
		assert.Equal(t, "DOWN", servers[2].Status)
	})
}

