package networkservicegroups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

const (
	networkServiceGroupsEndpoint = "/ztw/api/v1/networkServiceGroups"
)

type NetworkServiceGroups struct {
	// ID of the network service group.
	ID int `json:"id"`
	// Name of the network service group.
	Name string `json:"name,omitempty"`
	// The network services included in the group.
	Services []Services `json:"services,omitempty"`
	// Description of the network service group.
	Description string `json:"description,omitempty"`
	// Indicates whether the IP group or IP pool is created in Cloud & Branch Connector (EC) or ZIA.
	CreatorContext string `json:"creatorContext,omitempty"`
}

type Services struct {
	ID            int                            `json:"id"`
	Name          string                         `json:"name,omitempty"`
	Tag           string                         `json:"tag,omitempty"`
	SrcTCPPorts   []networkservices.NetworkPorts `json:"srcTcpPorts,omitempty"`
	DestTCPPorts  []networkservices.NetworkPorts `json:"destTcpPorts,omitempty"`
	SrcUDPPorts   []networkservices.NetworkPorts `json:"srcUdpPorts,omitempty"`
	DestUDPPorts  []networkservices.NetworkPorts `json:"destUdpPorts,omitempty"`
	Type          string                         `json:"type,omitempty"`
	Description   string                         `json:"description,omitempty"`
	IsNameL10nTag bool                           `json:"isNameL10nTag,omitempty"`
}

func GetNetworkServiceGroups(ctx context.Context, service *zscaler.Service, serviceGroupID int) (*NetworkServiceGroups, error) {
	var networkServiceGroups NetworkServiceGroups
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID), &networkServiceGroups)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning network service groups from Get: %d", networkServiceGroups.ID)
	return &networkServiceGroups, nil
}

func GetNetworkServiceGroupsByName(ctx context.Context, service *zscaler.Service, serviceGroupsName string) (*NetworkServiceGroups, error) {
	var networkServiceGroups []NetworkServiceGroups
	err := common.ReadAllPages(ctx, service.Client, networkServiceGroupsEndpoint, &networkServiceGroups)
	if err != nil {
		return nil, err
	}
	for _, networkServiceGroup := range networkServiceGroups {
		if strings.EqualFold(networkServiceGroup.Name, serviceGroupsName) {
			return &networkServiceGroup, nil
		}
	}
	return nil, fmt.Errorf("no network service groups found with name: %s", serviceGroupsName)
}

func CreateNetworkServiceGroups(ctx context.Context, service *zscaler.Service, networkServiceGroups *NetworkServiceGroups) (*NetworkServiceGroups, error) {
	resp, err := service.Client.CreateResource(ctx, networkServiceGroupsEndpoint, *networkServiceGroups)
	if err != nil {
		return nil, err
	}

	createdNetworkServiceGroups, ok := resp.(*NetworkServiceGroups)
	if !ok {
		return nil, errors.New("object returned from api was not a network service groups pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning network service groups from create: %d", createdNetworkServiceGroups.ID)
	return createdNetworkServiceGroups, nil
}

func UpdateNetworkServiceGroups(ctx context.Context, service *zscaler.Service, serviceGroupID int, networkServiceGroups *NetworkServiceGroups) (*NetworkServiceGroups, *http.Response, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID), *networkServiceGroups)
	if err != nil {
		return nil, nil, err
	}
	updatedNetworkServiceGroups, _ := resp.(*NetworkServiceGroups)

	service.Client.GetLogger().Printf("[DEBUG]returning network service groups from Update: %d", updatedNetworkServiceGroups.ID)
	return updatedNetworkServiceGroups, nil, nil
}

func DeleteNetworkServiceGroups(ctx context.Context, service *zscaler.Service, serviceGroupID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", networkServiceGroupsEndpoint, serviceGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllNetworkServiceGroups(ctx context.Context, service *zscaler.Service) ([]NetworkServiceGroups, error) {
	var networkServiceGroups []NetworkServiceGroups
	err := common.ReadAllPages(ctx, service.Client, networkServiceGroupsEndpoint, &networkServiceGroups)
	return networkServiceGroups, err
}
