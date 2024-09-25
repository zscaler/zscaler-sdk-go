package appservicegroups

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	appServicesGroupLiteEndpoint = "/zia/api/v1/appServiceGroups/lite"
)

type ApplicationServicesGroupLite struct {
	ID          int    `json:"id"`
	Name        string `json:"name,omitempty"`
	NameL10nTag bool   `json:"nameL10nTag"`
}

func GetByName(ctx context.Context, service *zscaler.Service, serviceGroupName string) (*ApplicationServicesGroupLite, error) {
	var appServicesGroupLite []ApplicationServicesGroupLite
	err := common.ReadAllPages(ctx, service.Client, appServicesGroupLiteEndpoint, &appServicesGroupLite)
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]ApplicationServicesGroupLite, error) {
	var appServiceGroups []ApplicationServicesGroupLite
	err := common.ReadAllPages(ctx, service.Client, appServicesGroupLiteEndpoint, &appServiceGroups)
	return appServiceGroups, err
}
