// Package unit provides unit tests for ZPA services
package unit

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func TestTagGroup_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagGroupID := "tg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/" + tagGroupID

	server.On("GET", path, common.SuccessResponse(tag_group.TagGroup{
		ID:          tagGroupID,
		Name:        "Test Tag Group",
		Description: "Test description",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, resp, err := tag_group.Get(context.Background(), service, tagGroupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotNil(t, resp)
	assert.Equal(t, tagGroupID, result.ID)
	assert.Equal(t, "Test Tag Group", result.Name)
}

func TestTagGroup_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagGroupName := "Production Tag Group"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_group.TagGroup{
			{ID: "tg-001", Name: "Other Group"},
			{ID: "tg-002", Name: tagGroupName},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_group.GetByName(context.Background(), service, tagGroupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "tg-002", result.ID)
	assert.Equal(t, tagGroupName, result.Name)
}

func TestTagGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup"

	server.On("POST", path, common.SuccessResponse(tag_group.TagGroup{
		ID:          "new-tg-123",
		Name:        "New Tag Group",
		Description: "Created via unit test",
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_group.Create(context.Background(), service, tag_group.TagGroup{
		Name:        "New Tag Group",
		Description: "Created via unit test",
		Tags:        []tag_group.Tag{},
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "new-tg-123", result.ID)
	assert.Equal(t, "New Tag Group", result.Name)
}

func TestTagGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagGroupID := "tg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/" + tagGroupID

	server.On("PUT", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_group.Update(context.Background(), service, tagGroupID, &tag_group.TagGroup{
		ID:          tagGroupID,
		Name:        "Updated Tag Group",
		Description: "Updated description",
		Tags:        []tag_group.Tag{},
	})

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagGroupID := "tg-12345"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/" + tagGroupID

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	resp, err := tag_group.Delete(context.Background(), service, tagGroupID)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestTagGroup_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list": []tag_group.TagGroup{
			{ID: "tg-001", Name: "Group 1"},
			{ID: "tg-002", Name: "Group 2"},
			{ID: "tg-003", Name: "Group 3"},
		},
		"totalPages": 1,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_group.GetAll(context.Background(), service)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, 3)
}

func TestTagGroup_GetByName_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/search"

	server.On("POST", path, common.SuccessResponse(map[string]interface{}{
		"list":       []tag_group.TagGroup{},
		"totalPages": 0,
	}))

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_group.GetByName(context.Background(), service, "NonExistent")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "no tag group named")
}

func TestTagGroup_Get_NotFound_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	tagGroupID := "nonexistent-id"
	path := "/zpa/mgmtconfig/v1/admin/customers/" + testCustomerID + "/tagGroup/" + tagGroupID

	server.On("GET", path, common.MockResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"id": "resource.not.found", "message": "Resource not found"}`,
	})

	service, err := common.CreateTestService(context.Background(), server, testCustomerID)
	require.NoError(t, err)

	result, _, err := tag_group.Get(context.Background(), service, tagGroupID)

	assert.Error(t, err)
	assert.Nil(t, result)
}
