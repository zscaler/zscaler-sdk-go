package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_incident_receiver_servers"
)

const incidentReceiverPath = "/zia/api/v1/incidentReceiverServers"

func TestDLPIncidentReceiverServers_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	receiverID := 100
	server.On("GET", incidentReceiverPath+"/100", common.SuccessResponse(dlp_incident_receiver_servers.IncidentReceiverServers{
		ID:     receiverID,
		Name:   "ZS_BD_INC_RECEIVER_01",
		URL:    "icaps://incident.example.com:1344",
		Status: "ENABLED",
		Flags:  0,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_incident_receiver_servers.Get(context.Background(), service, receiverID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, receiverID, result.ID)
	assert.Equal(t, "ZS_BD_INC_RECEIVER_01", result.Name)
	assert.Equal(t, "ENABLED", result.Status)
}

func TestDLPIncidentReceiverServers_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	receiverName := "ZS_BD_INC_RECEIVER_01"
	server.On("GET", incidentReceiverPath, common.SuccessResponse([]dlp_incident_receiver_servers.IncidentReceiverServers{
		{ID: 1, Name: "Other Receiver", Status: "DISABLED"},
		{ID: 100, Name: receiverName, URL: "icaps://incident.example.com:1344", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_incident_receiver_servers.GetByName(context.Background(), service, receiverName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, receiverName, result.Name)
}

func TestDLPIncidentReceiverServers_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", incidentReceiverPath, common.SuccessResponse([]dlp_incident_receiver_servers.IncidentReceiverServers{
		{ID: 100, Name: "ZS_BD_INC_RECEIVER_01", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_incident_receiver_servers.GetByName(context.Background(), service, "ThisReceiverDoesNotExist")

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDLPIncidentReceiverServers_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", incidentReceiverPath, common.SuccessResponse([]dlp_incident_receiver_servers.IncidentReceiverServers{
		{ID: 100, Name: "ZS_BD_INC_RECEIVER_01", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := dlp_incident_receiver_servers.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestDLPIncidentReceiverServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		receiver := dlp_incident_receiver_servers.IncidentReceiverServers{
			ID:     100,
			Name:   "ZS_BD_INC_RECEIVER_01",
			URL:    "icaps://incident.example.com:1344",
			Status: "ENABLED",
		}

		data, err := json.Marshal(receiver)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"ZS_BD_INC_RECEIVER_01"`)
		assert.Contains(t, string(data), `"status":"ENABLED"`)
	})
}
