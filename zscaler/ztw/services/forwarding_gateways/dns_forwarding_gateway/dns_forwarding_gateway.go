package dns_forwarding_gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	dnsGatewayEndpoint = "/ztw/api/v1/dnsGateways"
)

type DNSGateway struct {
	// A unique identifier assigned to the gateway
	ID int `json:"id"`

	// Name of the gateway.
	Name string `json:"name,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	Type string `json:"type,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	FailureBehavior string `json:"failureBehavior,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	DNSGatewayType string `json:"dnsGatewayType,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	PrimaryIP string `json:"primaryIp,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	SecondaryIP string `json:"secondaryIp,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	ECDNSGatewayOptionsPrimary string `json:"ecDnsGatewayOptionsPrimary,omitempty"`

	// Type of the gateway. Supported types are ZIA and ECSELF (Log and Control gateway).
	ECDNSGatewayOptionsSecondary string `json:"ecDnsGatewayOptionsSecondary,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Timestamp when the ZPA gateway was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ecGroupID int) (*DNSGateway, *http.Response, error) {
	var dnsGW DNSGateway
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, ecGroupID), &dnsGW)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Returning forwarding dns gateway from Get: %d", dnsGW.ID)
	return &dnsGW, nil, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, dnsGWName string) (*DNSGateway, error) {
	var dnsGW []DNSGateway
	// We are assuming this provisioning url name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayEndpoint, &dnsGW)
	if err != nil {
		return nil, err
	}
	for _, ec := range dnsGW {
		if strings.EqualFold(ec.Name, dnsGWName) {
			return &ec, nil
		}
	}
	return nil, fmt.Errorf("no forwarding dns gateway found with name: %s", dnsGWName)
}

func Create(ctx context.Context, service *zscaler.Service, rules *DNSGateway) (*DNSGateway, *http.Response, error) {
	resp, err := service.Client.CreateResource(ctx, dnsGatewayEndpoint, *rules)
	if err != nil {
		return nil, nil, err
	}

	createdRules, ok := resp.(*DNSGateway)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a forwarding dns gateway pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning forwarding dns gateway from create: %d", createdRules.ID)
	return createdRules, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *DNSGateway) (*DNSGateway, *http.Response, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, ruleID), *rules)
	if err != nil {
		return nil, nil, err
	}
	updatedGateways, _ := resp.(*DNSGateway)
	service.Client.GetLogger().Printf("[DEBUG]returning forwarding dns gateway from update: %d", updatedGateways.ID)
	return updatedGateways, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ecGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", dnsGatewayEndpoint, ecGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]DNSGateway, error) {
	var dnsGWs []DNSGateway
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayEndpoint+"/lite", &dnsGWs)
	return dnsGWs, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DNSGateway, error) {
	var dnsGWs []DNSGateway
	err := common.ReadAllPages(ctx, service.Client, dnsGatewayEndpoint, &dnsGWs)
	return dnsGWs, err
}
