package dnsgateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	dnsGatewayEndpoint     = "/ztw/api/v1/dnsGateways"
	dnsGatewayLiteEndpoint = "/ztw/api/v1/dnsGateways/lite"
)

type DNSGateway struct {
	ID                           int                            `json:"id,omitempty"`
	Name                         string                         `json:"name,omitempty"`
	DNSGatewayType               string                         `json:"dnsGatewayType,omitempty"`
	ECDnsGatewayOptionsPrimary   string                         `json:"ecDnsGatewayOptionsPrimary,omitempty"`
	ECDnsGatewayOptionsSecondary string                         `json:"ecDnsGatewayOptionsSecondary,omitempty"`
	FailureBehavior              string                         `json:"failureBehavior,omitempty"`
	PrimaryIP                    string                         `json:"primaryIp,omitempty"`
	SecondaryIP                  string                         `json:"secondaryIp,omitempty"`
	LastModifiedTime             int                            `json:"lastModifiedTime,omitempty"`
	LastModifiedBy               *common.CommonIDNameExternalID `json:"lastModifiedBy,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, gatewayID int) (*DNSGateway, error) {
	var gateway DNSGateway
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, gatewayID), &gateway)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning DNS gateway from Get: %d", gateway.ID)
	return &gateway, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, gatewayName string) (*DNSGateway, error) {
	var gateways []DNSGateway
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayEndpoint, &gateways)
	if err != nil {
		return nil, err
	}
	for _, gateway := range gateways {
		if strings.EqualFold(gateway.Name, gatewayName) {
			return &gateway, nil
		}
	}
	return nil, fmt.Errorf("no DNS gateway found with name: %s", gatewayName)
}

func Create(ctx context.Context, service *zscaler.Service, gateway *DNSGateway) (*DNSGateway, error) {
	resp, err := service.Client.CreateResource(ctx, dnsGatewayEndpoint, *gateway)
	if err != nil {
		return nil, err
	}
	createdGateway, ok := resp.(*DNSGateway)
	if !ok {
		return nil, errors.New("object returned from api was not a DNS gateway pointer")
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning DNS gateway from Create: %d", createdGateway.ID)
	return createdGateway, nil
}

func Update(ctx context.Context, service *zscaler.Service, gatewayID int, gateway *DNSGateway) (*DNSGateway, *http.Response, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, gatewayID), *gateway)
	if err != nil {
		return nil, nil, err
	}
	updatedGateway, _ := resp.(*DNSGateway)
	service.Client.GetLogger().Printf("[DEBUG] Returning DNS gateway from Update: %d", updatedGateway.ID)
	return updatedGateway, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, gatewayID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, gatewayID))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DNSGateway, error) {
	var gateways []DNSGateway
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayEndpoint, &gateways)
	return gateways, err
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]DNSGateway, error) {
	var gateways []DNSGateway
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayLiteEndpoint, &gateways)
	return gateways, err
}
