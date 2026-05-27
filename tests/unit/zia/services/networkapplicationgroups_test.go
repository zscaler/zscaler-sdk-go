package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplicationgroups"
)

const networkAppGroupsPath = "/zia/api/v1/networkApplicationGroups"

func sampleNetworkApplicationGroup() networkapplicationgroups.NetworkApplicationGroups {
	return networkapplicationgroups.NetworkApplicationGroups{
		Name:                "tests-nw-app-group",
		Description:         "tests-nw-app-group",
		NetworkApplications: []string{"APNS", "APPSTORE", "DICT"},
	}
}

func TestNetworkApplicationGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	group := sampleNetworkApplicationGroup()
	group.ID = groupID

	server.On("GET", networkAppGroupsPath+"/100", common.SuccessResponse(group))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplicationgroups.GetNetworkApplicationGroups(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, []string{"APNS", "APPSTORE", "DICT"}, result.NetworkApplications)
}

func TestNetworkApplicationGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "tests-nw-app-group"
	server.On("GET", networkAppGroupsPath, common.SuccessResponse([]networkapplicationgroups.NetworkApplicationGroups{
		{ID: 1, Name: "Other Group"},
		func() networkapplicationgroups.NetworkApplicationGroups {
			g := sampleNetworkApplicationGroup()
			g.ID = 100
			return g
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplicationgroups.GetNetworkApplicationGroupsByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestNetworkApplicationGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleNetworkApplicationGroup()
	created.ID = 99999

	server.On("POST", networkAppGroupsPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleNetworkApplicationGroup()
	result, err := networkapplicationgroups.Create(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestNetworkApplicationGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	updated := sampleNetworkApplicationGroup()
	updated.ID = groupID
	updated.Name = "updated-nw-app-group"

	server.On("PUT", networkAppGroupsPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := networkapplicationgroups.Update(context.Background(), service, groupID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-nw-app-group", result.Name)
}

func TestNetworkApplicationGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", networkAppGroupsPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkapplicationgroups.Delete(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestNetworkApplicationGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkAppGroupsPath, common.SuccessResponse([]networkapplicationgroups.NetworkApplicationGroups{
		sampleNetworkApplicationGroup(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkapplicationgroups.GetAllNetworkApplicationGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestNetworkApplicationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		group := sampleNetworkApplicationGroup()
		group.ID = 100

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"APNS"`)
	})
}
