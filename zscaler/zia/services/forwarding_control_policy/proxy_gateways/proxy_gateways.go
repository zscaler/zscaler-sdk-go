package proxy_gateways

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	proxyGatewaysEndpoint    = "/zia/api/v1/proxyGateways"
	proxyGatewayLiteEndpoint = "/zia/api/v1/proxyGateways/lite"
)

type ProxyGateways struct {
	// A unique identifier assigned to the Proxy gateway
	ID int `json:"id"`

	// The name of the Proxy gateway
	Name string `json:"name,omitempty"`

	// Additional details about the Proxy gateway
	Description string `json:"description,omitempty"`

	// The primary proxy for the gateway. This field is not applicable to the Lite API.
	PrimaryProxy *common.IDNameExternalID `json:"primaryProxy,omitempty"`

	// The seconday proxy for the gateway. This field is not applicable to the Lite API.
	SecondaryProxy *common.IDNameExternalID `json:"secondaryProxy,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Timestamp when the ZPA gateway was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Indicates whether fail close is enabled to drop the traffic or disabled to allow the traffic when both primary and secondary proxies defined in this gateway are unreachable.
	FailClosed bool `json:"failClosed"`

	// Proxy type: Supported Values: "PROXYCHAIN", "ZIA", "ECSELF"
	Type string `json:"type"`
}

// func Get(ctx context.Context, service *zscaler.Service, gatewayID int) (*ProxyGateways, error) {
// 	var rule ProxyGateways
// 	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", proxyGatewaysEndpoint, gatewayID), &rule)
// 	if err != nil {
// 		return nil, err
// 	}

// 	service.Client.GetLogger().Printf("[DEBUG]Returning zpa gateway from Get: %d", rule.ID)
// 	return &rule, nil
// }

func GetByName(ctx context.Context, service *zscaler.Service, gwName string) (*ProxyGateways, error) {
	var proxyGWs []ProxyGateways
	err := common.ReadAllPages(ctx, service.Client, proxyGatewaysEndpoint, &proxyGWs)
	if err != nil {
		return nil, err
	}
	for _, proxyGW := range proxyGWs {
		if strings.EqualFold(proxyGW.Name, gwName) {
			return &proxyGW, nil
		}
	}
	return nil, fmt.Errorf("no zpa gateway found with name: %s", gwName)
}

func GetLite(ctx context.Context, service *zscaler.Service) ([]ProxyGateways, error) {
	var proxyGWs []ProxyGateways
	err := common.ReadAllPages(ctx, service.Client, proxyGatewayLiteEndpoint, &proxyGWs)
	return proxyGWs, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ProxyGateways, error) {
	var proxyGWs []ProxyGateways
	err := common.ReadAllPages(ctx, service.Client, proxyGatewaysEndpoint, &proxyGWs)
	return proxyGWs, err
}
