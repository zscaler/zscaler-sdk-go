// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/np_client"
)

func TestNPClient_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/vpnConnectedUsers"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []np_client.NPClient{
			{Id: 1, UserName: "user1@example.com"},
			{Id: 2, UserName: "user2@example.com"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := np_client.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "user1@example.com", result[0].UserName)
}

func TestNPClient_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	userName := "testuser@example.com"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/vpnConnectedUsers"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list": []np_client.NPClient{
			{Id: 1, UserName: "other@example.com"},
			{Id: 2, UserName: userName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := np_client.GetByName(context.Background(), service, userName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.Id)
	assert.Equal(t, userName, result.UserName)
}

func TestNPClient_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/vpnConnectedUsers"

	server.On("GET", path, common.SuccessResponse(map[string]interface{}{
		"list":       []np_client.NPClient{},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := np_client.GetByName(context.Background(), service, "nonexistent@example.com")

	require.Error(t, err)
	assert.Nil(t, result)
}

