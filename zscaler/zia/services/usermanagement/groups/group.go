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
	var groups []Groups
	page := 1

	// Construct the endpoint with the search parameter
	endpointWithSearch := fmt.Sprintf("%s?search=%s&%s", groupsEndpoint, url.QueryEscape(targetGroup), common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)))

	for {
		err := common.ReadPage(ctx, service.Client, endpointWithSearch, page, &groups)
		if err != nil {
			return nil, err
		}

		// Iterate over the groups and check if the name matches the targetGroup
		for _, group := range groups {
			if strings.EqualFold(group.Name, targetGroup) {
				return &group, nil
			}
		}

		// Break the loop if there are no more pages
		if len(groups) < common.GetPageSize() {
			break
		}
		page++
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

func GetAllGroups(ctx context.Context, service *zscaler.Service) ([]Groups, error) {
	var groups []Groups
	err := common.ReadAllPages(ctx, service.Client, groupsEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &groups)
	return groups, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]Groups, error) {
	var groups []Groups
	err := common.ReadAllPages(ctx, service.Client, groupsEndpoint+"/lite", &groups)
	return groups, err
}
