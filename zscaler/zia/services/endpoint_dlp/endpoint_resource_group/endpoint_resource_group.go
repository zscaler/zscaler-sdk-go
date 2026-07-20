package endpoint_resource_group

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_channel"
)

const (
	endPointDlpResourceGroupsEndpoint = "/zia/api/v1/endPointDlpResourceGroups"
	dlpEndpointResourceEndpoint       = "/zia/api/v1/dlpEndpointResource"
)

type DlpEndpointResourceGroups struct {
	ID            int                                  `json:"id,omitempty"`
	Channel       string                               `json:"channel,omitempty"`
	Name          string                               `json:"name,omitempty"`
	Description   string                               `json:"description,omitempty"`
	Resources     []endpoint_resource.EndpointResource `json:"resources,omitempty"`
	ResourceCount int                                  `json:"resourceCount,omitempty"`
}

// EndpointDlpGroupToResourceAssociation represents the set of DLP resources to
// be associated with or removed from an existing tag group.
type EndpointDlpGroupToResourceAssociation struct {
	// ResourcesToBeAdded holds the IDs of DLP resources to be added to the tag group.
	ResourcesToBeAdded []int `json:"resourcesToBeAdded,omitempty"`

	// ResourcesToBeDeleted holds the IDs of DLP resources to be removed from the tag group.
	ResourcesToBeDeleted []int `json:"resourcesToBeDeleted,omitempty"`
}

// GetResourceTagsListFilterOptions holds the optional query parameters supported
// by GetResourceTagsList.
type GetResourceTagsListFilterOptions struct {
	// Name is the search string used to filter the list by DLP resource name or
	// other fields. Optional.
	Name string

	// SortOrder specifies the sorting order for the list by ascending or
	// descending order of the DLP resource tag names. Optional.
	SortOrder string

	// SearchResources must be set to true to include search strings via the Name
	// parameter. Optional.
	SearchResources *bool
}

func Create(ctx context.Context, service *zscaler.Service, resourceGroup *DlpEndpointResourceGroups) (*DlpEndpointResourceGroups, *http.Response, error) {
	resp, err := service.Client.Create(ctx, endPointDlpResourceGroupsEndpoint, *resourceGroup)
	if err != nil {
		return nil, nil, err
	}

	createdEndpointDlpRule, ok := resp.(*DlpEndpointResourceGroups)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a endpoint dlp resource group  pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new endpoint dlp resource group  from create: %d", createdEndpointDlpRule.ID)
	return createdEndpointDlpRule, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, resourceGroup *DlpEndpointResourceGroups) (*DlpEndpointResourceGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", endPointDlpResourceGroupsEndpoint, groupID), *resourceGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedResourceGroup, _ := resp.(*DlpEndpointResourceGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning updates endpoint dlp resource group from update: %d", updatedResourceGroup.ID)
	return updatedResourceGroup, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", endPointDlpResourceGroupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetResourceTags retrieves the list of tags, in name-ID pairs, to which the
// specified DLP resource is associated. Tags are channel-specific. The id
// parameter is required.
func GetResourceGroupTag(ctx context.Context, service *zscaler.Service, id int) ([]common.IDNameExternalID, error) {
	var tags []common.IDNameExternalID
	endpoint := fmt.Sprintf("%s/%d/groups", dlpEndpointResourceEndpoint, id)
	err := service.Client.Read(ctx, endpoint, &tags)
	return tags, err
}

// GetResourceTagsList retrieves the list of DLP resource tags added for the
// specified channel.
//
// The channel parameter is required and is part of the endpoint path. The name,
// sortOrder, and searchResources parameters are optional. The endpoint supports
// pagination, so common.ReadAllPages is used to aggregate all pages; the page
// and pageSize query parameters are handled internally by the pagination helper.
func GetResourceGroupTagsList(ctx context.Context, service *zscaler.Service, channel endpoint_resource_channel.Channel, opts *GetResourceTagsListFilterOptions) ([]DlpEndpointResourceGroups, error) {
	var resourceGroups []DlpEndpointResourceGroups
	endpoint := fmt.Sprintf("%s/%s", endPointDlpResourceGroupsEndpoint, url.PathEscape(string(channel)))

	queryParams := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			queryParams.Set("name", opts.Name)
		}
		if opts.SortOrder != "" {
			queryParams.Set("sortOrder", opts.SortOrder)
		}
		if opts.SearchResources != nil {
			queryParams.Set("searchResources", strconv.FormatBool(*opts.SearchResources))
		}
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &resourceGroups)
	return resourceGroups, err
}

func GetDlpResourceByTag(ctx context.Context, service *zscaler.Service, groupID int) ([]endpoint_resource.EndpointResource, error) {
	var tags []endpoint_resource.EndpointResource
	endpoint := fmt.Sprintf("%s/%d/resources", endPointDlpResourceGroupsEndpoint, groupID)
	err := service.Client.Read(ctx, endpoint, &tags)
	return tags, err
}

// UpdateDlpResourcesByTag updates the DLP resources association for the tag group
// with the specified group ID. The groupID parameter is required, and the
// resources to add and/or remove are supplied in the request body.
func UpdateDlpResourcesByTag(ctx context.Context, service *zscaler.Service, groupID int, association *EndpointDlpGroupToResourceAssociation) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%d/resources", endPointDlpResourceGroupsEndpoint, groupID)
	_, err := service.Client.UpdateWithPut(ctx, endpoint, *association)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] updated dlp resources association for tag group: %d", groupID)
	return nil, nil
}
