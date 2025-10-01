package ecgroup

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	ecGroupEndpoint     = "/ztw/api/v1/ecgroup"
	ecGroupLiteEndpoint = "/ztw/api/v1/ecgroup/lite"
)

type EcGroup struct {
	ID                    int                            `json:"id,omitempty"`
	Name                  string                         `json:"name,omitempty"`
	Description           string                         `json:"desc,omitempty"`
	DeployType            string                         `json:"deployType,omitempty"`
	Status                []string                       `json:"status,omitempty"`
	Platform              string                         `json:"platform,omitempty"`
	AWSAvailabilityZone   string                         `json:"awsAvailabilityZone,omitempty"`
	AzureAvailabilityZone string                         `json:"azureAvailabilityZone,omitempty"`
	MaxEcCount            int                            `json:"maxEcCount,omitempty"`
	TunnelMode            string                         `json:"tunnelMode,omitempty"`
	Location              *common.CommonIDNameExternalID `json:"location,omitempty"`
	ProvTemplate          *common.CommonIDNameExternalID `json:"provTemplate,omitempty"`
	ECVMs                 []common.ECVMs                 `json:"ecVMs,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ecGroupID int) (*EcGroup, error) {
	var ecGroup EcGroup
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", ecGroupEndpoint, ecGroupID), &ecGroup)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Cloud & Branch Connector Group from Get: %d", ecGroup.ID)
	return &ecGroup, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ecGroupName string) (*EcGroup, error) {
	var ecGroup []EcGroup
	// We are assuming this provisioning url name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, ecGroupEndpoint, &ecGroup)
	if err != nil {
		return nil, err
	}
	for _, ec := range ecGroup {
		if strings.EqualFold(ec.Name, ecGroupName) {
			return &ec, nil
		}
	}
	return nil, fmt.Errorf("no Cloud & Branch Connector Group found with name: %s", ecGroupName)
}

func Delete(ctx context.Context, service *zscaler.Service, ecGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", ecGroupEndpoint, ecGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]EcGroup, error) {
	var ecgroups []EcGroup
	err := common.ReadAllPages(ctx, service.Client, ecGroupEndpoint, &ecgroups)
	return ecgroups, err
}

func GetEcGroupLiteID(ctx context.Context, service *zscaler.Service, ecGroupID int) (*EcGroup, error) {
	var ecgroupLite EcGroup
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", ecGroupLiteEndpoint, ecGroupID), &ecgroupLite)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning Cloud & Branch Connector Group from Get: %d", ecgroupLite.ID)
	return &ecgroupLite, nil
}

func GetEcGroupLiteByName(ctx context.Context, service *zscaler.Service, ecGroupLiteName string) (*EcGroup, error) {
	var ecgroupLite []EcGroup
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?name=%s", ecGroupLiteEndpoint, url.QueryEscape(ecGroupLiteName)), &ecgroupLite)
	if err != nil {
		return nil, err
	}
	for _, ecgroupLite := range ecgroupLite {
		if strings.EqualFold(ecgroupLite.Name, ecGroupLiteName) {
			return &ecgroupLite, nil
		}
	}
	return nil, fmt.Errorf("no Cloud & Branch Connector Group found with name: %s", ecGroupLiteName)
}
