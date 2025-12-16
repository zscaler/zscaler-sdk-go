// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/ipsourcegroups"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestIPSourceGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipSourceGroups/12345"

	server.On("GET", path, common.SuccessResponse(ipsourcegroups.IPSourceGroups{
		ID:          groupID,
		Name:        "Corporate IPs",
		Description: "Corporate IP addresses",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipsourcegroups.Get(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Equal(t, "Corporate IPs", result.Name)
}

func TestIPSourceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "Corporate IPs"
	path := "/zia/api/v1/ipSourceGroups"

	server.On("GET", path, common.SuccessResponse([]ipsourcegroups.IPSourceGroups{
		{ID: 1, Name: "Other Group"},
		{ID: 2, Name: groupName},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipsourcegroups.GetByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestIPSourceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipSourceGroups"

	server.On("GET", path, common.SuccessResponse([]ipsourcegroups.IPSourceGroups{
		{ID: 1, Name: "Group 1"},
		{ID: 2, Name: "Group 2"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := ipsourcegroups.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestIPSourceGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/ipSourceGroups"

	server.On("POST", path, common.SuccessResponse(ipsourcegroups.IPSourceGroups{
		ID:   99999,
		Name: "New IP Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newGroup := &ipsourcegroups.IPSourceGroups{
		Name: "New IP Group",
	}

	result, err := ipsourcegroups.Create(context.Background(), service, newGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestIPSourceGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipSourceGroups/12345"

	server.On("PUT", path, common.SuccessResponse(ipsourcegroups.IPSourceGroups{
		ID:   groupID,
		Name: "Updated IP Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateGroup := &ipsourcegroups.IPSourceGroups{
		ID:   groupID,
		Name: "Updated IP Group",
	}

	result, err := ipsourcegroups.Update(context.Background(), service, groupID, updateGroup)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated IP Group", result.Name)
}

func TestIPSourceGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 12345
	path := "/zia/api/v1/ipSourceGroups/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = ipsourcegroups.Delete(context.Background(), service, groupID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

func TestIPSourceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IPSourceGroups JSON marshaling", func(t *testing.T) {
		group := ipsourcegroups.IPSourceGroups{
			ID:          12345,
			Name:        "Corporate IPs",
			Description: "Corporate IP addresses",
			IPAddresses: []string{"10.0.0.0/8", "192.168.0.0/16"},
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Corporate IPs"`)
	})

	t.Run("IPSourceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Branch Office IPs",
			"description": "Branch office IP addresses",
			"ipAddresses": ["172.16.0.0/12"]
		}`

		var group ipsourcegroups.IPSourceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Equal(t, "Branch Office IPs", group.Name)
	})
}
