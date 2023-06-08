package applicationservices

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	appServicesLiteEndpoint = "/appServices/lite"
)

type ApplicationServicesLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	NameL10nTag bool   `json:"nameL10nTag"`
}

func (service *Service) Get(serviceID int) (*ApplicationServicesLite, error) {
	var appServicesLite ApplicationServicesLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", appServicesLiteEndpoint, serviceID), &appServicesLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning application services from Get: %d", appServicesLite.ID)
	return &appServicesLite, nil
}

func (service *Service) GetByName(serviceName string) (*ApplicationServicesLite, error) {
	var appServicesLite []ApplicationServicesLite
	err := common.ReadAllPagesWithFilters(service.Client, appServicesLiteEndpoint, map[string]string{"search": serviceName}, &appServicesLite)
	if err != nil {
		return nil, err
	}
	for _, appServicesLite := range appServicesLite {
		if strings.EqualFold(appServicesLite.Name, serviceName) {
			return &appServicesLite, nil
		}
	}
	return nil, fmt.Errorf("no application services found with name: %s", serviceName)
}
