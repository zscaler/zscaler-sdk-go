package dns_gateways

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dnsGatewaysEndpoint = "/zia/api/v1/dnsGateways"
)

type DNSGateways struct {
	ID                  int                      `json:"id,omitempty"`
	Name                string                   `json:"name,omitempty"`
	DnsGatewayType      string                   `json:"dnsGatewayType,omitempty"`
	PrimaryIpOrFqdn     string                   `json:"primaryIpOrFqdn,omitempty"`
	PrimaryPorts        []int                    `json:"primaryPorts,omitempty"`
	SecondaryIpOrFqdn   string                   `json:"secondaryIpOrFqdn,omitempty"`
	SecondaryPorts      []int                    `json:"secondaryPorts,omitempty"`
	Protocols           []string                 `json:"protocols,omitempty"`
	FailureBehavior     string                   `json:"failureBehavior,omitempty"`
	LastModifiedTime    int                      `json:"lastModifiedTime,omitempty"`
	LastModifiedBy      *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`
	AutoCreated         bool                     `json:"autoCreated,omitempty"`
	NatZtrGateway       bool                     `json:"natZtrGateway,omitempty"`
	DnsGatewayProtocols []string                 `json:"dnsGatewayProtocols,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, gwID int) (*DNSGateways, error) {
	var dnsGateway DNSGateways
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dnsGatewaysEndpoint, gwID), &dnsGateway)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning dns gatewayfrom Get: %d", dnsGateway.ID)
	return &dnsGateway, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, gwName string) (*DNSGateways, error) {
	var dnsGateways []DNSGateways
	err := common.ReadAllPages(ctx, service.Client, dnsGatewaysEndpoint, &dnsGateways)
	if err != nil {
		return nil, err
	}
	for _, dnsGateway := range dnsGateways {
		if strings.EqualFold(dnsGateway.Name, gwName) {
			return &dnsGateway, nil
		}
	}
	return nil, fmt.Errorf("no dns gateway found with name: %s", gwName)
}

func Create(ctx context.Context, service *zscaler.Service, gwID *DNSGateways) (*DNSGateways, *http.Response, error) {
	resp, err := service.Client.Create(ctx, dnsGatewaysEndpoint, *gwID)
	if err != nil {
		return nil, nil, err
	}

	createdDnsGW, ok := resp.(*DNSGateways)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a dns gateway pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new dns gateway from create: %d", createdDnsGW.ID)
	return createdDnsGW, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, gwID int, dnsGWs *DNSGateways) (*DNSGateways, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dnsGatewaysEndpoint, gwID), *dnsGWs)
	if err != nil {
		return nil, nil, err
	}
	updatedDnsGW, _ := resp.(*DNSGateways)

	service.Client.GetLogger().Printf("[DEBUG]returning updates dns gateway  from update: %d", updatedDnsGW.ID)
	return updatedDnsGW, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, gwID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dnsGatewaysEndpoint, gwID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]DNSGateways, error) {
	var dnsGateways []DNSGateways
	err := common.ReadAllPages(ctx, service.Client, dnsGatewaysEndpoint+"/lite", &dnsGateways)
	return dnsGateways, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DNSGateways, error) {
	var dnsGateways []DNSGateways
	err := common.ReadAllPages(ctx, service.Client, dnsGatewaysEndpoint, &dnsGateways)
	return dnsGateways, err
}
