package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservicegroups"
)

const networkServiceGroupsPath = "/zia/api/v1/networkServiceGroups"

func sampleNetworkServiceGroup() networkservicegroups.NetworkServiceGroups {
	return networkservicegroups.NetworkServiceGroups{
		Name:        "tests-nw-svc-group",
		Description: "tests-nw-svc-group",
		Services: []networkservicegroups.Services{
			{ID: 100},
			{ID: 101},
			{ID: 102},
		},
	}
}

func TestNetworkServiceGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	group := sampleNetworkServiceGroup()
	group.ID = groupID

	server.On("GET", networkServiceGroupsPath+"/100", common.SuccessResponse(group))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetNetworkServiceGroups(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Len(t, result.Services, 3)
}

func TestNetworkServiceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "tests-nw-svc-group"
	server.On("GET", networkServiceGroupsPath, common.SuccessResponse([]networkservicegroups.NetworkServiceGroups{
		{ID: 1, Name: "Other Group"},
		func() networkservicegroups.NetworkServiceGroups {
			g := sampleNetworkServiceGroup()
			g.ID = 100
			return g
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetNetworkServiceGroupsByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestNetworkServiceGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleNetworkServiceGroup()
	created.ID = 99999

	server.On("POST", networkServiceGroupsPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleNetworkServiceGroup()
	result, err := networkservicegroups.CreateNetworkServiceGroups(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestNetworkServiceGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	updated := sampleNetworkServiceGroup()
	updated.ID = groupID
	updated.Name = "updated-nw-svc-group"

	server.On("PUT", networkServiceGroupsPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := networkservicegroups.UpdateNetworkServiceGroups(context.Background(), service, groupID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-nw-svc-group", result.Name)
}

func TestNetworkServiceGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", networkServiceGroupsPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkservicegroups.DeleteNetworkServiceGroups(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestNetworkServiceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkServiceGroupsPath, common.SuccessResponse([]networkservicegroups.NetworkServiceGroups{
		sampleNetworkServiceGroup(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestNetworkServiceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		group := sampleNetworkServiceGroup()
		group.ID = 100

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"tests-nw-svc-group"`)
	})
}
