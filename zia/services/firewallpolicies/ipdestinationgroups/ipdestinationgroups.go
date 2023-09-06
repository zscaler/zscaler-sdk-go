package ipdestinationgroups

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	ipDestinationGroupsEndpoint = "/ipDestinationGroups"
)

type IPDestinationGroups struct {
	// Unique identifer for the destination IP group
	ID int `json:"id"`

	// Destination IP group name
	Name string `json:"name,omitempty"`

	// Additional information about the destination IP group
	Description string `json:"description,omitempty"`

	// Destination IP group type (i.e., the group can contain destination IP addresses or FQDNs)
	Type string `json:"type,omitempty"`

	// Destination IP addresses, FQDNs, or wildcard FQDNs added to the group.
	Addresses []string `json:"addresses,omitempty"`

	// Destination IP address URL categories. You can identify destinations based on the URL category of the domain.
	IPCategories []string `json:"ipCategories,omitempty"`

	// Destination IP address countries. You can identify destinations based on the location of a server.
	Countries []string `json:"countries,omitempty"`

	// If set to true, the destination IP address group is non-editable. This field is applicable only to predefined IP address groups, which cannot be modified.
	IsNonEditable bool `json:"isNonEditable,omitempty"`
}

func (service *Service) Get(ipGroupID int) (*IPDestinationGroups, error) {
	var ipDestinationGroups IPDestinationGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID), &ipDestinationGroups)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning ip destination group from Get: %d", ipDestinationGroups.ID)
	return &ipDestinationGroups, nil
}

func (service *Service) GetByName(ipDestinationGroupsName string) (*IPDestinationGroups, error) {
	var ipDestinationGroups []IPDestinationGroups
	err := common.ReadAllPages(service.Client, ipDestinationGroupsEndpoint, &ipDestinationGroups)
	if err != nil {
		return nil, err
	}
	for _, ipDestinationGroup := range ipDestinationGroups {
		if strings.EqualFold(ipDestinationGroup.Name, ipDestinationGroupsName) {
			return &ipDestinationGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip destination group found with name: %s", ipDestinationGroupsName)
}

func (service *Service) Create(ipGroupID *IPDestinationGroups) (*IPDestinationGroups, error) {
	resp, err := service.Client.Create(ipDestinationGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPDestinationGroups, ok := resp.(*IPDestinationGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip destination group pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning ip destination group from create: %d", createdIPDestinationGroups.ID)
	return createdIPDestinationGroups, nil
}

func (service *Service) Update(ipGroupID int, ipGroup *IPDestinationGroups) (*IPDestinationGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID), *ipGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedIPDestinationGroups, _ := resp.(*IPDestinationGroups)

	service.Client.Logger.Printf("[DEBUG]returning ip destination group from update: %d", updatedIPDestinationGroups.ID)
	return updatedIPDestinationGroups, nil, nil
}

func (service *Service) Delete(ipGroupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", ipDestinationGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (service *Service) GetAll() ([]IPDestinationGroups, error) {
	var ipDestinationGroups []IPDestinationGroups
	err := common.ReadAllPages(service.Client, ipDestinationGroupsEndpoint, &ipDestinationGroups)
	return ipDestinationGroups, err
}
