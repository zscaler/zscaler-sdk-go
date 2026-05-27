package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplications"
)

const networkApplicationsPath = "/zia/api/v1/networkApplications"

func TestNetworkApplications_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkApplicationsPath+"/APNS", common.SuccessResponse(networkapplications.NetworkApplications{
		ID:             "APNS",
		ParentCategory: "GENERAL",
		Description:    "Apple Push Notification Service",
		Deprecated:     false,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplications.GetNetworkApplication(context.Background(), service, "APNS", "")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "APNS", result.ID)
	assert.Equal(t, "GENERAL", result.ParentCategory)
}

func TestNetworkApplications_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkApplicationsPath, common.SuccessResponse([]networkapplications.NetworkApplications{
		{ID: "APNS", ParentCategory: "GENERAL", Description: "Apple Push Notification Service"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplications.GetByName(context.Background(), service, "APNS", "en-US")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "APNS", result.ID)
}

func TestNetworkApplications_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkApplicationsPath, common.SuccessResponse([]networkapplications.NetworkApplications{
		{ID: "APNS", ParentCategory: "GENERAL"},
		{ID: "GARP", ParentCategory: "GENERAL"},
		{ID: "DIAMETER", ParentCategory: "GENERAL"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplications.GetAll(context.Background(), service, "en-US")

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestNetworkApplications_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		app := networkapplications.NetworkApplications{
			ID:             "APNS",
			ParentCategory: "GENERAL",
			Description:    "Apple Push Notification Service",
			Deprecated:     false,
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"APNS"`)
		assert.Contains(t, string(data), `"deprecated":false`)
	})
}
