package groups

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	groupsEndpoint = "/groups"
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

func (service *Service) GetGroups(groupID int) (*Groups, error) {
	var groups Groups
	err := service.Client.Read(fmt.Sprintf("%s/%d", groupsEndpoint, groupID), &groups)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning Groups from Get: %d", groups.ID)
	return &groups, nil
}

func (service *Service) GetGroupByName(targetGroup string) (*Groups, error) {
	var groups []Groups
	page := 1

	// Construct the endpoint with the search parameter
	endpointWithSearch := fmt.Sprintf("%s?search=%s", groupsEndpoint, url.QueryEscape(targetGroup))

	for {
		err := common.ReadPage(service.Client, endpointWithSearch, page, &groups)
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

func (service *Service) GetAllGroups() ([]Groups, error) {
	var groups []Groups
	err := common.ReadAllPages(service.Client, groupsEndpoint, &groups)
	return groups, err
}
