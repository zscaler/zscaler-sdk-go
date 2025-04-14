package forwardinggateways

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
	forwardGatewayEndpoint = "/ztw/api/v1/gateways"
)

type ECGateway struct {
	// A unique identifier assigned to the gateway
	ID int `json:"id"`

	// Name of the gateway.
	Name string `json:"name,omitempty"`

	// Description of the gateway.
	Description string `json:"description,omitempty"`

	// indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable.
	FailClosed bool `json:"failClosed,omitempty"`

	// indicates that traffic must be dropped when both primary and secondary proxies defined in the gateway are unreachable.
	ManualPrimary string `json:"failClomanualPrimarysed,omitempty"`

	// Specifies the secondary proxy through which traffic must be forwarded.
	ManualSecondary string `json:"manualSecondary,omitempty"`

	// Not applicable to Cloud & Branch Connector.
	SubCloudPrimary *common.CommonIDNameExternalID `json:"subcloudPrimary,omitempty"`

	// Not applicable to Cloud & Branch Connector.
	SubCloudSecondary *common.CommonIDNameExternalID `json:"subcloudSecondary,omitempty"`

	// Type of the primary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through selected data centers
	// or override (MANUAL_OVERRIDE) that forwards traffic through a specified IP address or domain.
	// Supported Values: "NONE", "AUTO", "MANUAL_OVERRIDE", "SUBCLOUD", "VZEN", "PZEN", "DC"
	PrimaryType string `json:"primaryType,omitempty"`

	// Type of the secondary proxy, such as automatic proxy (AUTO), manual proxy (DC) that forwards traffic through selected data centers,
	// or override (MANUAL_OVERRIDE) that forwards traffic through a specified IP address or domain.
	// Supported Values: "NONE", "AUTO", "MANUAL_OVERRIDE", "SUBCLOUD", "VZEN", "PZEN", "DC"
	SecondaryType string `json:"secondaryType,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.CommonIDNameExternalID `json:"lastModifiedBy,omitempty"`

	// Timestamp when the ZPA gateway was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ecGroupID int) (*ECGateway, error) {
	var ecGW ECGateway
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", forwardGatewayEndpoint, ecGroupID), &ecGW)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning forwarding gateway from Get: %d", ecGW.ID)
	return &ecGW, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ecGWName string) (*ECGateway, error) {
	var ecGW []ECGateway
	// We are assuming this provisioning url name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, forwardGatewayEndpoint, &ecGW)
	if err != nil {
		return nil, err
	}
	for _, ec := range ecGW {
		if strings.EqualFold(ec.Name, ecGWName) {
			return &ec, nil
		}
	}
	return nil, fmt.Errorf("no forwarding gateway found with name: %s", ecGWName)
}

func Create(ctx context.Context, service *zscaler.Service, rules *ECGateway) (*ECGateway, error) {
	resp, err := service.Client.CreateResource(ctx, forwardGatewayEndpoint, *rules)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ECGateway)
	if !ok {
		return nil, errors.New("object returned from api was not a forwarding gateway pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning forwarding gateway from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *ECGateway) (*ECGateway, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", forwardGatewayEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedGateways, _ := resp.(*ECGateway)
	service.Client.GetLogger().Printf("[DEBUG]returning forwarding gateway from update: %d", updatedGateways.ID)
	return updatedGateways, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ecGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", forwardGatewayEndpoint, ecGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]ECGateway, error) {
	var ecGWs []ECGateway
	err := common.ReadAllPages(ctx, service.Client, forwardGatewayEndpoint+"/lite", &ecGWs)
	return ecGWs, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ECGateway, error) {
	var ecGWs []ECGateway
	err := common.ReadAllPages(ctx, service.Client, forwardGatewayEndpoint, &ecGWs)
	return ecGWs, err
}
