package networkservices

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
	networkServicesEndpoint = "/ztw/api/v1/networkServices"
)

type NetworkServices struct {
	// ID of network service.
	ID int `json:"id"`

	// Name of network service.
	Name string `json:"name,omitempty"`

	// Description of network service.
	Description string `json:"description,omitempty"`

	Tag string `json:"tag,omitempty"`
	// Source TCP ports.

	SrcTCPPorts []NetworkPorts `json:"srcTcpPorts,omitempty"`

	// Destination TCP ports.
	DestTCPPorts []NetworkPorts `json:"destTcpPorts,omitempty"`

	// Source UDP ports.
	SrcUDPPorts []NetworkPorts `json:"srcUdpPorts,omitempty"`

	// Destination UDP ports.
	DestUDPPorts []NetworkPorts `json:"destUdpPorts,omitempty"`

	// Type of network service: standard, predefined, or custom.
	// Supported values: STANDARD, PREDEFINED, CUSTOM
	Type string `json:"type,omitempty"`

	// Indicates whether name is a tag that can be used to look up the display string, typically from a localization resource bundle.
	IsNameL10nTag bool `json:"isNameL10nTag,omitempty"`

	// Indicates whether the IP group or IP pool is created in Cloud & Branch Connector (EC) or ZIA.
	CreatorContext string `json:"creatorContext,omitempty"`
}

type NetworkPorts struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, serviceID int) (*NetworkServices, error) {
	var networkServices NetworkServices
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), &networkServices)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning network services from Get: %d", networkServices.ID)
	return &networkServices, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, networkServiceName string) (*NetworkServices, error) {
	var networkServices []NetworkServices
	err := common.ReadAllPages(ctx, service.Client, networkServicesEndpoint, &networkServices)
	if err != nil {
		return nil, err
	}
	for _, networkService := range networkServices {
		if strings.EqualFold(networkService.Name, networkServiceName) {
			return &networkService, nil
		}
	}
	return nil, fmt.Errorf("no network services found with name: %s", networkServiceName)
}

func Create(ctx context.Context, service *zscaler.Service, networkService *NetworkServices) (*NetworkServices, error) {
	resp, err := service.Client.CreateResource(ctx, networkServicesEndpoint, *networkService)
	if err != nil {
		return nil, err
	}

	createdNetworkServices, ok := resp.(*NetworkServices)
	if !ok {
		return nil, errors.New("object returned from api was not a network service pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning network service from create: %d", createdNetworkServices.ID)
	return createdNetworkServices, nil
}

func Update(ctx context.Context, service *zscaler.Service, serviceID int, networkService *NetworkServices) (*NetworkServices, *http.Response, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID), *networkService)
	if err != nil {
		return nil, nil, err
	}
	updatedNetworkServices, _ := resp.(*NetworkServices)

	service.Client.GetLogger().Printf("[DEBUG]returning network service from Update: %d", updatedNetworkServices.ID)
	return updatedNetworkServices, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, serviceID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", networkServicesEndpoint, serviceID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllNetworkServices(ctx context.Context, service *zscaler.Service) ([]NetworkServices, error) {
	var networkServices []NetworkServices
	err := common.ReadAllPages(ctx, service.Client, networkServicesEndpoint, &networkServices)
	return networkServices, err
}
