package appservicegroups

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	appServicesGroupLiteEndpoint = "/appServiceGroups/lite"
)

type ApplicationServicesGroupLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	NameL10nTag bool   `json:"nameL10nTag"`
}

func (service *Service) Get(serviceGroupID int) (*ApplicationServicesGroupLite, error) {
	var appServicesGroupLite ApplicationServicesGroupLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", appServicesGroupLiteEndpoint, serviceGroupID), &appServicesGroupLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning app services group from Get: %d", appServicesGroupLite.ID)
	return &appServicesGroupLite, nil
}

func (service *Service) GetByName(serviceGroupName string) (*ApplicationServicesGroupLite, error) {
	var appServicesGroupLite []ApplicationServicesGroupLite
	err := common.ReadAllPages(service.Client, appServicesGroupLiteEndpoint, &appServicesGroupLite)
	if err != nil {
		return nil, err
	}
	for _, appServicesGroupLite := range appServicesGroupLite {
		if strings.EqualFold(appServicesGroupLite.Name, serviceGroupName) {
			return &appServicesGroupLite, nil
		}
	}
	return nil, fmt.Errorf("no app services group found with name: %s", serviceGroupName)
}

func (service *Service) GetAll() ([]ApplicationServicesGroupLite, error) {
	var appServiceGroups []ApplicationServicesGroupLite
	err := common.ReadAllPages(service.Client, appServicesGroupLiteEndpoint, &appServiceGroups)
	return appServiceGroups, err
}
