// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_group"
)

const (
	dlpEndpointResourceGroupsBasePath = "/zia/api/v1/endPointDlpResourceGroups"
	dlpEndpointResourceGroupsResource = "/zia/api/v1/dlpEndpointResource"
)

// =====================================================
// GetResourceTags
// =====================================================

func TestEndpointResourceGroup_GetResourceTags_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsResource + "/42/groups"

	server.On("GET", path, common.SuccessResponse([]ziacommon.IDNameExternalID{
		{ID: 1, Name: "Tag A", ExternalID: "ext-1"},
		{ID: 2, Name: "Tag B"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetResourceGroupTag(context.Background(), service, 42)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 1, result[0].ID)
	assert.Equal(t, "Tag A", result[0].Name)
	assert.Equal(t, "ext-1", result[0].ExternalID)
}

func TestEndpointResourceGroup_GetResourceTags_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsResource + "/42/groups"
	server.On("GET", path, common.SuccessResponse([]ziacommon.IDNameExternalID{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetResourceGroupTag(context.Background(), service, 42)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestEndpointResourceGroup_GetResourceTags_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsResource + "/999/groups"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetResourceGroupTag(context.Background(), service, 999)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetResourceTagsList
// =====================================================

func TestEndpointResourceGroup_GetResourceTagsList_NoOpts_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/" + string(endpoint_resource_channel.ChannelPrinting)

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Empty(t, q.Get("name"))
		assert.Empty(t, q.Get("sortOrder"))
		assert.Empty(t, q.Get("searchResources"))
		return common.SuccessResponse([]endpoint_resource_group.DlpEndpointResourceGroups{
			{ID: 1, Name: "Tag Group A"},
			{ID: 2, Name: "Tag Group B"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetResourceGroupTagsList(context.Background(), service, endpoint_resource_channel.ChannelPrinting, nil)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestEndpointResourceGroup_GetResourceTagsList_WithFilters_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/" + string(endpoint_resource_channel.ChannelNetworkDriveTransfer)

	server.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
		q := r.URL.Query()
		assert.Equal(t, "finance", q.Get("name"))
		assert.Equal(t, "desc", q.Get("sortOrder"))
		assert.Equal(t, "true", q.Get("searchResources"))
		return common.SuccessResponse([]endpoint_resource_group.DlpEndpointResourceGroups{
			{ID: 5, Name: "Finance Tag"},
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	search := true
	opts := &endpoint_resource_group.GetResourceTagsListFilterOptions{
		Name:            "finance",
		SortOrder:       "desc",
		SearchResources: &search,
	}
	result, err := endpoint_resource_group.GetResourceGroupTagsList(context.Background(), service, endpoint_resource_channel.ChannelNetworkDriveTransfer, opts)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 5, result[0].ID)
}

func TestEndpointResourceGroup_GetResourceTagsList_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/" + string(endpoint_resource_channel.ChannelPersonalCloudStorage)
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetResourceGroupTagsList(context.Background(), service, endpoint_resource_channel.ChannelPersonalCloudStorage, nil)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// GetDlpResourcesByTag
// =====================================================

func TestEndpointResourceGroup_GetDlpResourcesByTag_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/10/resources"

	server.On("GET", path, common.SuccessResponse([]endpoint_resource.EndpointResource{
		{ID: 100, Name: "Resource One"},
		{ID: 200, Name: "Resource Two"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetDlpResourceByTag(context.Background(), service, 10)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 100, result[0].ID)
}

func TestEndpointResourceGroup_GetDlpResourcesByTag_Empty_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/10/resources"
	server.On("GET", path, common.SuccessResponse([]endpoint_resource.EndpointResource{}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetDlpResourceByTag(context.Background(), service, 10)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestEndpointResourceGroup_GetDlpResourcesByTag_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/999/resources"
	server.On("GET", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := endpoint_resource_group.GetDlpResourceByTag(context.Background(), service, 999)

	require.Error(t, err)
	assert.Empty(t, result)
}

// =====================================================
// UpdateDlpResourcesByTag
// =====================================================

func TestEndpointResourceGroup_UpdateDlpResourcesByTag_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/10/resources"

	server.OnFunc("PUT", path, func(_ *http.Request, body []byte) common.MockResponse {
		assert.Contains(t, string(body), "resourcesToBeAdded")
		assert.Contains(t, string(body), "resourcesToBeDeleted")
		return common.NoContentResponse()
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	association := &endpoint_resource_group.EndpointDlpGroupToResourceAssociation{
		ResourcesToBeAdded:   []int{1, 2, 3},
		ResourcesToBeDeleted: []int{4, 5},
	}

	_, err = endpoint_resource_group.UpdateDlpResourcesByTag(context.Background(), service, 10, association)

	require.NoError(t, err)
}

func TestEndpointResourceGroup_UpdateDlpResourcesByTag_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/999/resources"
	server.On("PUT", path, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	association := &endpoint_resource_group.EndpointDlpGroupToResourceAssociation{
		ResourcesToBeAdded: []int{1},
	}

	_, err = endpoint_resource_group.UpdateDlpResourcesByTag(context.Background(), service, 999, association)

	require.Error(t, err)
}

// =====================================================
// Create / Update / Delete
// =====================================================

func TestEndpointResourceGroup_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("POST", dlpEndpointResourceGroupsBasePath, func(_ *http.Request, body []byte) common.MockResponse {
		assert.Contains(t, string(body), "Finance Group")
		return common.SuccessResponse(endpoint_resource_group.DlpEndpointResourceGroups{
			ID:      301,
			Name:    "Finance Group",
			Channel: "PRINTING",
		})
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource_group.Create(context.Background(), service, &endpoint_resource_group.DlpEndpointResourceGroups{
		Name:    "Finance Group",
		Channel: "PRINTING",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 301, result.ID)
}

func TestEndpointResourceGroup_Create_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("POST", dlpEndpointResourceGroupsBasePath, common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource_group.Create(context.Background(), service, &endpoint_resource_group.DlpEndpointResourceGroups{
		Name: "Finance Group",
	})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointResourceGroup_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/301"

	server.On("PUT", path, common.SuccessResponse(endpoint_resource_group.DlpEndpointResourceGroups{
		ID:   301,
		Name: "Updated Group",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource_group.Update(context.Background(), service, 301, &endpoint_resource_group.DlpEndpointResourceGroups{
		ID:   301,
		Name: "Updated Group",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Group", result.Name)
}

func TestEndpointResourceGroup_Update_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("PUT", dlpEndpointResourceGroupsBasePath+"/301", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := endpoint_resource_group.Update(context.Background(), service, 301, &endpoint_resource_group.DlpEndpointResourceGroups{ID: 301})

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestEndpointResourceGroup_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := dlpEndpointResourceGroupsBasePath + "/301"
	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_resource_group.Delete(context.Background(), service, 301)

	require.NoError(t, err)
}

func TestEndpointResourceGroup_Delete_Error_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", dlpEndpointResourceGroupsBasePath+"/999", common.NotFoundResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = endpoint_resource_group.Delete(context.Background(), service, 999)

	require.Error(t, err)
}
