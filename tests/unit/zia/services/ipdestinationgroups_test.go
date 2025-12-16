// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipdestinationgroups"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestIPDestinationGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipDestinationGroups/12345"

	server.On("GET", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:          groupID,
		Name:        "Corporate Servers",
		Description: "Internal server IPs",
		Type:        "DSTN_IP",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Corporate Servers", result.Name)
}

func TestIPDestinationGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Corporate Servers"
	path := "/zia/api/v1/ipDestinationGroups"

	server.On("GET", path, common.SuccessResponse([]ipdestinationgroups.IPDestinationGroups{
		{ID: 1, Name: "Other Group"},
		{ID: 2, Name: groupName},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestIPDestinationGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipDestinationGroups"

	server.On("GET", path, common.SuccessResponse([]ipdestinationgroups.IPDestinationGroups{
		{ID: 1, Name: "Group 1", Type: "DSTN_IP"},
		{ID: 2, Name: "Group 2", Type: "DSTN_FQDN"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipdestinationgroups.GetAll(context.Background(), service, "")

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIPDestinationGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipDestinationGroups"

	server.On("POST", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:   99999,
		Name: "New IP Group",
		Type: "DSTN_IP",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newGroup := &ipdestinationgroups.IPDestinationGroups{
		Name: "New IP Group",
		Type: "DSTN_IP",
	}

	result, err := ipdestinationgroups.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestIPDestinationGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipDestinationGroups/12345"

	server.On("PUT", path, common.SuccessResponse(ipdestinationgroups.IPDestinationGroups{
		ID:   groupID,
		Name: "Updated IP Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGroup := &ipdestinationgroups.IPDestinationGroups{
		ID:   groupID,
		Name: "Updated IP Group",
	}

	result, _, err := ipdestinationgroups.Update(context.Background(), service, groupID, updateGroup, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated IP Group", result.Name)
}

func TestIPDestinationGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipDestinationGroups/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = ipdestinationgroups.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

func TestIPDestinationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPDestinationGroups JSON marshaling", func(t *testing.T) {
		group := ipdestinationgroups.IPDestinationGroups{
			ID:          12345,
			Name:        "Corporate Servers",
			Description: "Internal server IPs",
			Type:        "DSTN_IP",
			Addresses:   []string{"10.0.0.0/8", "192.168.0.0/16"},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Corporate Servers"`)
		assert.Contains(t, string(data), `"type":"DSTN_IP"`)
	})

	t.Run("IPDestinationGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "External Servers",
			"description": "External server addresses",
			"type": "DSTN_FQDN",
			"addresses": ["*.example.com", "api.company.com"],
			"ipCategories": ["CATEGORY_1"],
			"countries": ["US", "CA"]
		}`

		var group ipdestinationgroups.IPDestinationGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "DSTN_FQDN", group.Type)
		assert.Len(t, group.Addresses, 2)
	})
}

