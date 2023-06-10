package applicationservicesgroup

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
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

	service.Client.Logger.Printf("[DEBUG]Returning application services group from Get: %d", appServicesGroupLite.ID)
	return &appServicesGroupLite, nil
}

func (service *Service) GetByName(serviceGroupName string) (*ApplicationServicesGroupLite, error) {
	var appServicesGroupLite []ApplicationServicesGroupLite
	err := common.ReadAllPagesWithFilters(service.Client, appServicesGroupLiteEndpoint, map[string]string{"search": serviceGroupName}, &appServicesGroupLite)
	if err != nil {
		return nil, err
	}
	for _, appServicesGroupLite := range appServicesGroupLite {
		if strings.EqualFold(appServicesGroupLite.Name, serviceGroupName) {
			return &appServicesGroupLite, nil
		}
	}
	return nil, fmt.Errorf("no application services group found with name: %s", serviceGroupName)
}
