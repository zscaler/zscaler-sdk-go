package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_icap_servers"
)

const icapServersPath = "/zia/api/v1/icapServers"

func TestDLPICAPServers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverID := 100
	server.On("GET", icapServersPath+"/100", common.SuccessResponse(dlp_icap_servers.DLPICAPServers{
		ID:     serverID,
		Name:   "ZS_BD_ICAP_01",
		URL:    "icaps://icap.example.com:1344",
		Status: "ENABLED",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_icap_servers.Get(context.Background(), service, serverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serverID, result.ID)
	assert.Equal(t, "ZS_BD_ICAP_01", result.Name)
	assert.Equal(t, "icaps://icap.example.com:1344", result.URL)
	assert.Equal(t, "ENABLED", result.Status)
}

func TestDLPICAPServers_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serverName := "ZS_BD_ICAP_01"
	server.On("GET", icapServersPath, common.SuccessResponse([]dlp_icap_servers.DLPICAPServers{
		{ID: 1, Name: "Other ICAP", URL: "icaps://other.example.com:1344", Status: "DISABLED"},
		{ID: 100, Name: serverName, URL: "icaps://icap.example.com:1344", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_icap_servers.GetByName(context.Background(), service, serverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serverName, result.Name)
}

func TestDLPICAPServers_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", icapServersPath, common.SuccessResponse([]dlp_icap_servers.DLPICAPServers{
		{ID: 100, Name: "ZS_BD_ICAP_01", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_icap_servers.GetByName(context.Background(), service, "ThisIcapServerDoesNotExist")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPICAPServers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", icapServersPath, common.SuccessResponse([]dlp_icap_servers.DLPICAPServers{
		{ID: 100, Name: "ZS_BD_ICAP_01", URL: "icaps://icap.example.com:1344", Status: "ENABLED"},
		{ID: 101, Name: "ZS_BD_ICAP_02", URL: "icaps://icap2.example.com:1344", Status: "DISABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_icap_servers.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestDLPICAPServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		srv := dlp_icap_servers.DLPICAPServers{
			ID:     100,
			Name:   "ZS_BD_ICAP_01",
			URL:    "icaps://icap.example.com:1344",
			Status: "ENABLED",
		}

		data, err := json.Marshal(srv)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"ZS_BD_ICAP_01"`)
		assert.Contains(t, string(data), `"status":"ENABLED"`)
	})
}
