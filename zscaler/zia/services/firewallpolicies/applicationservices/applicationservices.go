package applicationservices

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	appServicesLiteEndpoint = "/zia/api/v1/appServices/lite"
)

type ApplicationServicesLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	NameL10nTag bool   `json:"nameL10nTag"`
}

func GetByName(ctx context.Context, service *zscaler.Service, serviceName string) (*ApplicationServicesLite, error) {
	var appServicesLite []ApplicationServicesLite
	err := common.ReadAllPages(ctx, service.Client, appServicesLiteEndpoint, &appServicesLite)
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]ApplicationServicesLite, error) {
	var appServices []ApplicationServicesLite
	err := common.ReadAllPages(ctx, service.Client, appServicesLiteEndpoint, &appServices)
	return appServices, err
}
