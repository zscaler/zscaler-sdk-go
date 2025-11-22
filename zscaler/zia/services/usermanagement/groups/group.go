package groups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	groupsEndpoint = "/zia/api/v1/groups"
)

type Groups struct {
	// Unique identfier for the group
	ID int `json:"id"`

	// Group name
	Name string `json:"name,omitempty"`

	// Unique identfier for the identity provider (IdP)
	IdpID int `json:"idpId"`

	// Additional information about the group
	Comments string `json:"comments,omitempty"`

	// Additional information about the group
	IsSystemDefined bool `json:"isSystemDefined,omitempty"`
}

func GetGroups(ctx context.Context, service *zscaler.Service, groupID int) (*Groups, error) {
	var groups Groups
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", groupsEndpoint, groupID), &groups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning Groups from Get: %d", groups.ID)
	return &groups, nil
}

func GetGroupByName(ctx context.Context, service *zscaler.Service, targetGroup string) (*Groups, error) {
	// Use GetAllGroups with search parameter to leverage built-in pagination
	opts := &GetAllGroupsFilterOptions{
		Search: targetGroup,
	}

	groups, err := GetAllGroups(ctx, service, opts)
	if err != nil {
		return nil, err
	}

	// Iterate over the groups and check if the name matches the targetGroup exactly
	for _, group := range groups {
		if strings.EqualFold(group.Name, targetGroup) {
			return &group, nil
		}
	}

	return nil, fmt.Errorf("no group found with name: %s", targetGroup)
}

func Create(ctx context.Context, service *zscaler.Service, groupID *Groups) (*Groups, *http.Response, error) {
	resp, err := service.Client.Create(ctx, groupsEndpoint, *groupID)
	if err != nil {
		return nil, nil, err
	}

	createdGroup, ok := resp.(*Groups)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new group from create: %d", createdGroup.ID)
	return createdGroup, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, groupID int, groups *Groups) (*Groups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", groupsEndpoint, groupID), *groups)
	if err != nil {
		return nil, nil, err
	}
	updatedGroup, _ := resp.(*Groups)

	service.Client.GetLogger().Printf("[DEBUG]returning updates group from update: %d", updatedGroup.ID)
	return updatedGroup, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", groupsEndpoint, groupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAllGroupsFilterOptions represents optional query parameters for GetAllGroups
type GetAllGroupsFilterOptions struct {
	// Search string used to match against a group's name or comments attributes
	Search string
	// The string value defined by the group name or other applicable attributes
	DefinedBy string
}

func GetAllGroups(ctx context.Context, service *zscaler.Service, opts *GetAllGroupsFilterOptions) ([]Groups, error) {
	var groups []Groups
	endpoint := groupsEndpoint

	// Build query parameters from filter options
	queryParams := url.Values{}
	if opts != nil {
		if opts.Search != "" {
			queryParams.Add("search", opts.Search)
		}
		if opts.DefinedBy != "" {
			queryParams.Add("definedBy", opts.DefinedBy)
		}
	}

	// Add sort parameters using service defaults (always use service.SortBy and service.SortOrder)
	if string(service.SortBy) != "" {
		queryParams.Add("sortBy", string(service.SortBy))
	}
	if string(service.SortOrder) != "" {
		queryParams.Add("sortOrder", string(service.SortOrder))
	}

	// Build endpoint with query parameters
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &groups)
	return groups, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]Groups, error) {
	var groups []Groups
	err := common.ReadAllPages(ctx, service.Client, groupsEndpoint+"/lite", &groups)
	return groups, err
}
