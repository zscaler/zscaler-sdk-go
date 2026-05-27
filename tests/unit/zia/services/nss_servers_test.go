// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

const nssServersPath = "/zia/api/v1/nssServers"

// sampleNSSServer mirrors the integration test payload in nss_servers_test.go.
func sampleNSSServer(name string) nss_servers.NSSServers {
	return nss_servers.NSSServers{
		Name:   name,
		Status: "ENABLED",
		Type:   "NSS_FOR_FIREWALL",
	}
}

// =====================================================
// SDK Function Tests
// =====================================================

func TestNSSServers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := nssServersPath + "/12345"

	srv := sampleNSSServer("tests-nss-firewall-server")
	srv.ID = nssID
	srv.State = "ENABLED"

	server.On("GET", path, common.SuccessResponse(srv))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.Get(context.Background(), service, nssID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, nssID, result.ID)
	assert.Equal(t, "tests-nss-firewall-server", result.Name)
	assert.Equal(t, "ENABLED", result.Status)
	assert.Equal(t, "NSS_FOR_FIREWALL", result.Type)
}

func TestNSSServers_Get_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath+"/9999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.Get(context.Background(), service, 9999)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNSSServers_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverName := "tests-nss-firewall-server"

	srv := sampleNSSServer(serverName)
	srv.ID = 2

	server.On("GET", nssServersPath, common.SuccessResponse([]nss_servers.NSSServers{
		{ID: 1, Name: "other-server", Status: "ENABLED", Type: "NSS_FOR_WEB"},
		srv,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetByName(context.Background(), service, serverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)
	assert.Equal(t, serverName, result.Name)
	assert.Equal(t, "NSS_FOR_FIREWALL", result.Type)
}

func TestNSSServers_GetByName_CaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath, common.SuccessResponse([]nss_servers.NSSServers{
		{ID: 5, Name: "Tests-NSS-Firewall-Server", Status: "ENABLED", Type: "NSS_FOR_FIREWALL"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetByName(context.Background(), service, "TESTS-NSS-FIREWALL-SERVER")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Tests-NSS-Firewall-Server", result.Name)
}

func TestNSSServers_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath, common.SuccessResponse([]nss_servers.NSSServers{
		{ID: 1, Name: "existing-server"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetByName(context.Background(), service, "non_existent_name")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no nss server found with name: non_existent_name")
}

func TestNSSServers_GetByName_APIError_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetByName(context.Background(), service, "any")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNSSServers_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := nssServersPath

	created := sampleNSSServer("tests-new-nss-server")
	created.ID = 100

	server.On("POST", path, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newServer := sampleNSSServer("tests-new-nss-server")

	result, err := nss_servers.Create(context.Background(), service, &newServer)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 100, result.ID)
	assert.Equal(t, "ENABLED", result.Status)
	assert.Equal(t, "NSS_FOR_FIREWALL", result.Type)
}

func TestNSSServers_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", nssServersPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newServer := sampleNSSServer("tests-fail")
	result, err := nss_servers.Create(context.Background(), service, &newServer)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNSSServers_Create_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", nssServersPath, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newServer := sampleNSSServer("tests-no-body")
	result, err := nss_servers.Create(context.Background(), service, &newServer)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "object returned from api was not a nss server Pointer")
}

func TestNSSServers_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := nssServersPath + "/12345"

	updated := sampleNSSServer("tests-updated-nss-server")
	updated.ID = nssID

	server.On("PUT", path, common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateServer := sampleNSSServer("tests-updated-nss-server")
	updateServer.ID = nssID

	result, err := nss_servers.Update(context.Background(), service, nssID, &updateServer)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tests-updated-nss-server", result.Name)
}

func TestNSSServers_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", nssServersPath+"/12345", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateServer := sampleNSSServer("tests-updated-nss-server")
	result, err := nss_servers.Update(context.Background(), service, 12345, &updateServer)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNSSServers_Update_UnexpectedResponseType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", nssServersPath+"/12345", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateServer := sampleNSSServer("tests-updated-nss-server")
	result, err := nss_servers.Update(context.Background(), service, 12345, &updateServer)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "object returned from api was not a nss server Pointer")
}

func TestNSSServers_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nssID := 12345
	path := nssServersPath + "/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nss_servers.Delete(context.Background(), service, nssID)

	require.NoError(t, err)
}

func TestNSSServers_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", nssServersPath+"/12345", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = nss_servers.Delete(context.Background(), service, 12345)

	require.Error(t, err)
}

func TestNSSServers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath, common.SuccessResponse([]nss_servers.NSSServers{
		sampleNSSServer("server-1"),
		sampleNSSServer("server-2"),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestNSSServers_GetAll_WithTypeFilter_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", nssServersPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "type=NSS_FOR_FIREWALL", r.URL.RawQuery)
		return common.SuccessResponse([]nss_servers.NSSServers{
			sampleNSSServer("firewall-server"),
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	serverType := "NSS_FOR_FIREWALL"
	result, err := nss_servers.GetAll(context.Background(), service, &serverType)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "NSS_FOR_FIREWALL", result[0].Type)
}

func TestNSSServers_GetAll_WithTypeFilterCaseInsensitive_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", nssServersPath, func(r *http.Request, _ []byte) common.MockResponse {
		assert.Equal(t, "type=NSS_FOR_WEB", r.URL.RawQuery)
		return common.SuccessResponse([]nss_servers.NSSServers{
			{Name: "web-server", Status: "ENABLED", Type: "NSS_FOR_WEB"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	serverType := " nss_for_web "
	result, err := nss_servers.GetAll(context.Background(), service, &serverType)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestNSSServers_GetAll_InvalidType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	invalid := "INVALID_TYPE"
	result, err := nss_servers.GetAll(context.Background(), service, &invalid)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid server type: INVALID_TYPE")
}

func TestNSSServers_GetAll_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", nssServersPath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := nss_servers.GetAll(context.Background(), service, nil)

	require.Error(t, err)
	assert.Nil(t, result)
}

// =====================================================
// Structure Tests
// =====================================================

func TestNSSServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NSSServers JSON marshaling", func(t *testing.T) {
		srv := sampleNSSServer("tests-nss-firewall-server")
		srv.ID = 12345
		srv.State = "ENABLED"
		srv.IcapSvrId = 100

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"NSS_FOR_FIREWALL"`)
		assert.Contains(t, string(data), `"status":"ENABLED"`)
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

