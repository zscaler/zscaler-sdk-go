package ipdestinationgroups

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	ipDestinationGroupsEndpoint = "/ipDestinationGroups"
)

type IPDestinationGroups struct {
	ID           int      `json:"id"`
	Name         string   `json:"name,omitempty"`
	Type         string   `json:"type,omitempty"`
	Addresses    []string `json:"addresses,omitempty"`
	Description  string   `json:"description,omitempty"`
	IPCategories []string `json:"ipCategories,omitempty"`
	Countries    []string `json:"countries,omitempty"`
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
