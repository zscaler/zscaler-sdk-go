package groups

import (
	"context"
	"fmt"
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
}

func GetGroups(ctx context.Context, service *zscaler.Service, groupID int) (*Groups, error) {
	var groups Groups
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", groupsEndpoint, groupID), &groups)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning Groups from Get: %d", groups.ID)
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

func GetAllGroups(ctx context.Context, service *zscaler.Service) ([]Groups, error) {
	var groups []Groups
	err := common.ReadAllPages(ctx, service.Client, groupsEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &groups)
	return groups, err
}
