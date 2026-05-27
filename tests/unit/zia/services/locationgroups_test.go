package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/location/locationgroups"
)

const locationGroupsPath = "/zia/api/v1/locations/groups"
const locationGroupsCountPath = "/zia/api/v1/locations/groups/count"

func TestLocationGroups_GetLocationGroup_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	server.On("GET", locationGroupsPath+"/100", common.SuccessResponse(locationgroups.LocationGroup{
		ID:        groupID,
		Name:      "US Locations",
		GroupType: "STATIC_GROUP",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroup(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "STATIC_GROUP", result.GroupType)
}

func TestLocationGroups_GetLocationGroupByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "US Locations"
	server.On("GET", locationGroupsPath, common.SuccessResponse([]locationgroups.LocationGroup{
		{ID: 1, Name: "Other Group", GroupType: "DYNAMIC_GROUP"},
		{ID: 100, Name: groupName, GroupType: "STATIC_GROUP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroupByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestLocationGroups_GetGroupType_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", locationGroupsPath, common.SuccessResponse([]locationgroups.LocationGroup{
		{ID: 100, Name: "Dynamic Group", GroupType: "DYNAMIC_GROUP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetGroupType(context.Background(), service, "DYNAMIC_GROUP")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "DYNAMIC_GROUP", result.GroupType)
}

func TestLocationGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", locationGroupsPath, common.SuccessResponse([]locationgroups.LocationGroup{
		{ID: 1, Name: "Static Group", GroupType: "STATIC_GROUP"},
		{ID: 2, Name: "Dynamic Group", GroupType: "DYNAMIC_GROUP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetAll(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestLocationGroups_GetLocationGroupCount_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", locationGroupsCountPath, common.SuccessResponse(5))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroupCount(context.Background(), service, nil)

	require.NoError(t, err)
	assert.Equal(t, 5, result)
}

func TestLocationGroups_GetLocationGroupCount_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	name := "US Locations"
	groupType := "STATIC_GROUP"
	fetchLocations := true
	server.On("GET", locationGroupsCountPath, common.SuccessResponse(2))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroupCount(context.Background(), service, &locationgroups.GetAllFilterOptions{
		Name:           &name,
		GroupType:      &groupType,
		FetchLocations: &fetchLocations,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, result)
}

func TestLocationGroups_GetLocationGroup_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", locationGroupsPath+"/99999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := locationgroups.GetLocationGroup(context.Background(), service, 99999)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestLocationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		group := locationgroups.LocationGroup{
			ID:        100,
			Name:      "US Locations",
			GroupType: "STATIC_GROUP",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"groupType":"STATIC_GROUP"`)
	})
}
