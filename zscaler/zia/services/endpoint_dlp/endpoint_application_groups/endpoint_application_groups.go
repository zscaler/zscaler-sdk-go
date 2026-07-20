package endpoint_application_groups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_applications"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/endpoint_dlp/endpoint_resource_group"
)

const (
	endPointApplicationGroupsEndpoint = "/zia/api/v1/endPointApplicationGroups"
)

// ApplicationGroupResource is a single application reference used in the
// create/update body of an endpoint application group. The API identifies each
// member by its application ID (the endpoint application's zappId) and its
// display name.
type ApplicationGroupResource struct {
	AppID string `json:"appId,omitempty"`
	Name  string `json:"name,omitempty"`
}

// EndpointApplicationGroups is the create/update request/response body for an
// endpoint application group.
//
// Note: the list endpoint (GetAll) returns a different shape — see
// common.EndPointApplicationGroups — where each member is expressed through the
// shared endPointApplications list rather than the writable resources list.
type EndpointApplicationGroups struct {
	ID            int                        `json:"id,omitempty"`
	Channel       string                     `json:"channel,omitempty"`
	Name          string                     `json:"name,omitempty"`
	Description   string                     `json:"description,omitempty"`
	ResourceCount int                        `json:"resourceCount,omitempty"`
	Resources     []ApplicationGroupResource `json:"resources,omitempty"`
}

func Create(ctx context.Context, service *zscaler.Service, group *EndpointApplicationGroups) (*EndpointApplicationGroups, error) {
	resp, err := service.Client.Create(ctx, endPointApplicationGroupsEndpoint, *group)
	if err != nil {
		return nil, err
	}

	createdEndpointApplicationGroup, ok := resp.(*EndpointApplicationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an endpoint application group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning endpoint application group from create: %d", createdEndpointApplicationGroup.ID)
	return createdEndpointApplicationGroup, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, group *EndpointApplicationGroups) (*EndpointApplicationGroups, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", endPointApplicationGroupsEndpoint, groupID), *group)
	if err != nil {
		return nil, err
	}
	updatedEndpointApplicationGroup, _ := resp.(*EndpointApplicationGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning endpoint application group from update: %d", updatedEndpointApplicationGroup.ID)
	return updatedEndpointApplicationGroup, nil
}

func UpdateApplicationGroupResources(ctx context.Context, service *zscaler.Service, groupID int, association *endpoint_resource_group.EndpointDlpGroupToResourceAssociation) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%d/resources", endPointApplicationGroupsEndpoint, groupID)
	_, err := service.Client.UpdateWithPut(ctx, endpoint, *association)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] updated applications association for the specified tag group: %d", groupID)
	return nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", endPointApplicationGroupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll retrieves the list of application groups. This endpoint supports a
// maximum page size of 5000, which is used to minimize the number of requests.
//
// The list endpoint returns each group in the shared common.EndPointApplicationGroups
// shape (groupId plus the endPointApplications members), which is intentionally
// different from the create/update body. There is no dedicated get-by-id
// endpoint, so Get/GetByName resolve a single group from this list.
func GetAll(ctx context.Context, service *zscaler.Service) ([]common.EndPointApplicationGroups, error) {
	var groups []common.EndPointApplicationGroups
	err := common.ReadAllPages(ctx, service.Client, endPointApplicationGroupsEndpoint, &groups, 5000)
	return groups, err
}

// Get resolves a single application group by its ID from the list endpoint.
func Get(ctx context.Context, service *zscaler.Service, groupID int) (*common.EndPointApplicationGroups, error) {
	groups, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for i := range groups {
		if groups[i].GroupID == groupID {
			return &groups[i], nil
		}
	}
	return nil, fmt.Errorf("no endpoint application group found with id: %d", groupID)
}

// GetByName resolves a single application group by its name from the list endpoint.
func GetByName(ctx context.Context, service *zscaler.Service, name string) (*common.EndPointApplicationGroups, error) {
	groups, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for i := range groups {
		if strings.EqualFold(groups[i].Name, name) {
			return &groups[i], nil
		}
	}
	return nil, fmt.Errorf("no endpoint application group found with name: %s", name)
}

func GetApplicationGroupPolicies(ctx context.Context, service *zscaler.Service, resourceIDs []int) ([]endpoint_applications.ApplicationPolicies, error) {
	var policies []endpoint_applications.ApplicationPolicies
	endpoint := endPointApplicationGroupsEndpoint + "/policies"

	queryParams := url.Values{}
	for _, id := range resourceIDs {
		queryParams.Add("resourceId", strconv.Itoa(id))
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &policies)
	return policies, err
}
