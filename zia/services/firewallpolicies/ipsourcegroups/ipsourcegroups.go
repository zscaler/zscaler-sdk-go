package ipsourcegroups

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	ipSourceGroupsEndpoint = "/ipSourceGroups"
)

type IPSourceGroups struct {
	// A unique identifier of the source IP address group.
	ID int `json:"id"`

	// The name of the source IP address group.
	Name string `json:"name,omitempty"`

	// The description of the source IP address group.
	Description string `json:"description,omitempty"`

	// Source IP addresses added to the group.
	IPAddresses []string `json:"ipAddresses,omitempty"`

	// If set to true, the destination IP address group is non-editable. This field is applicable only to predefined IP address groups, which cannot be modified.
	IsNonEditable bool `json:"isNonEditable,omitempty"`
}

func (service *Service) Get(ipGroupID int) (*IPSourceGroups, error) {
	var ipSourceGroups IPSourceGroups
	err := service.Client.Read(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), &ipSourceGroups)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning ip source groupfrom Get: %d", ipSourceGroups.ID)
	return &ipSourceGroups, nil
}

func (service *Service) GetByName(ipSourceGroupsName string) (*IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := common.ReadAllPages(service.Client, ipSourceGroupsEndpoint, &ipSourceGroups)
	if err != nil {
		return nil, err
	}
	for _, ipSourceGroup := range ipSourceGroups {
		if strings.EqualFold(ipSourceGroup.Name, ipSourceGroupsName) {
			return &ipSourceGroup, nil
		}
	}
	return nil, fmt.Errorf("no ip source group found with name: %s", ipSourceGroupsName)
}

func (service *Service) Create(ipGroupID *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.Create(ipSourceGroupsEndpoint, *ipGroupID)
	if err != nil {
		return nil, err
	}

	createdIPSourceGroups, ok := resp.(*IPSourceGroups)
	if !ok {
		return nil, errors.New("object returned from api was not an ip source group pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning ip source group from create: %d", createdIPSourceGroups.ID)
	return createdIPSourceGroups, nil
}

func (service *Service) Update(ipGroupID int, ipGroup *IPSourceGroups) (*IPSourceGroups, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID), *ipGroup)
	if err != nil {
		return nil, err
	}
	updatedIPSourceGroups, _ := resp.(*IPSourceGroups)

	service.Client.Logger.Printf("[DEBUG]returning ip source group from update: %d", updatedIPSourceGroups.ID)
	return updatedIPSourceGroups, nil
}

func (service *Service) Delete(ipGroupID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", ipSourceGroupsEndpoint, ipGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (service *Service) GetAll() ([]IPSourceGroups, error) {
	var ipSourceGroups []IPSourceGroups
	err := common.ReadAllPages(service.Client, ipSourceGroupsEndpoint, &ipSourceGroups)
	return ipSourceGroups, err
}
