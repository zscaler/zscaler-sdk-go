package applicationservices

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	appServicesLiteEndpoint = "/appServices/lite"
)

type ApplicationServicesLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	NameL10nTag bool   `json:"nameL10nTag"`
}

func GetByName(service *services.Service, serviceName string) (*ApplicationServicesLite, error) {
	var appServicesLite []ApplicationServicesLite
	err := common.ReadAllPages(service.Client, appServicesLiteEndpoint, &appServicesLite)
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

func GetAll(service *services.Service) ([]ApplicationServicesLite, error) {
	var appServices []ApplicationServicesLite
	err := common.ReadAllPages(service.Client, appServicesLiteEndpoint, &appServices)
	return appServices, err
}
