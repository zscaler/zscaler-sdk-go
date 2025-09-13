package vzen_nodes

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
	vzenNodeEndpoint = "/zia/api/v1/virtualZenNodes"
)

type VZENNodes struct {
	// System-generated identifier for the Virtual Service Edge. Ignored for POST and PUT requests.
	ID int `json:"id,omitempty"`

	// The Zscaler service gateway ID
	ZGatewayID int `json:"zgatewayId,omitempty"`

	// The Virtual Service Edge name
	Name string `json:"name,omitempty"`

	// Indicates the Virtual Service Edge status. Supported values: ENABLED, DISABLED, DISABLED_BY_SERVICE_PROVIDER, NOT_PROVISIONED_IN_SERVICE_PROVIDER, IN_TRIAL
	Status string `json:"status,omitempty"`

	// Represents the Virtual Service Edge instances deployed for production purposes
	InProduction bool `json:"inProduction,omitempty"`

	// The IP address to which you forward the traffic.
	// All user and server workload traffic is forwarded to the proxy IP address of the Virtual Service Edge.
	// If the Virtual Service Edge has to receive and service traffic from users or workloads over the internet, ensure that this IP address has access to both the internet and users.
	IPAddress string `json:"ipAddress,omitempty"`

	// The corresponding subnet mask
	SubnetMask string `json:"subnetMask,omitempty"`

	// The IP address of the default gateway to the internet
	DefaultGateway string `json:"defaultGateway,omitempty"`

	// The Virtual Service Edge subscription type
	Type string `json:"type,omitempty"`

	// A Boolean value that indicates whether IPSec is enabled. Enable this option to terminate IPSec traffic from the client at the Virtual Service Edge node.
	IPSecEnabled bool `json:"ipSecEnabled,omitempty"`

	// A Boolean value that indicates whether or not the On-Demand Support Tunnel is enabled
	OnDemandSupportTunnelEnabled bool `json:"onDemandSupportTunnelEnabled,omitempty"`

	// A Boolean value that indicates whether or not a support tunnel for Zscaler Support is enabled.
	// Enable this option to allow the service to establish a support tunnel for Zscaler Support to access the Virtual Service Edge.
	EstablishSupportTunnelEnabled bool `json:"establishSupportTunnelEnabled,omitempty"`

	// The IP address of the load balancer. This field is applicable only when the 'deploymentMode' field is set to CLUSTER.
	LoadBalancerIPAddress string `json:"loadBalancerIpAddress,omitempty"`

	// Specifies the deployment mode. Select either STANDALONE or CLUSTER if you have the VMware ESXi platform. Otherwise, select only STANDALONE.
	// Supporteed Values: STANDALONE, CLUSTER
	DeploymentMode string `json:"deploymentMode,omitempty"`

	// Virtual Service Edge cluster name
	ClusterName string `json:"clusterName,omitempty"`

	// The Virtual Service Edge SKU type
	// Supported Values: SMALL, MEDIUM, LARGE
	VzenSkuType string `json:"vzenSkuType,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, nodeID int) (*VZENNodes, error) {
	var vzenNode VZENNodes
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", vzenNodeEndpoint, nodeID), &vzenNode)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning vzen node from Get: %d", vzenNode.ID)
	return &vzenNode, nil
}

func GetNodeByName(ctx context.Context, service *zscaler.Service, nodeName string) (*VZENNodes, error) {
	var vzenNodes []VZENNodes
	err := common.ReadAllPages(ctx, service.Client, vzenNodeEndpoint, &vzenNodes)
	if err != nil {
		return nil, err
	}
	for _, vzenNode := range vzenNodes {
		if strings.EqualFold(vzenNode.Name, nodeName) {
			return &vzenNode, nil
		}
	}
	return nil, fmt.Errorf("no vzen node found with name: %s", nodeName)
}

func Create(ctx context.Context, service *zscaler.Service, vzenNodes *VZENNodes) (*VZENNodes, *http.Response, error) {
	resp, err := service.Client.Create(ctx, vzenNodeEndpoint, *vzenNodes)
	if err != nil {
		return nil, nil, err
	}

	createdNode, ok := resp.(*VZENNodes)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a vzen node pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new vzen node from create: %d", createdNode.ID)
	return createdNode, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, nodeID int, vzenNodes *VZENNodes) (*VZENNodes, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", vzenNodeEndpoint, nodeID), *vzenNodes)
	if err != nil {
		return nil, nil, err
	}
	updatedNode, _ := resp.(*VZENNodes)

	service.Client.GetLogger().Printf("[DEBUG]returning updates vzen node from update: %d", updatedNode.ID)
	return updatedNode, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleLabelID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", vzenNodeEndpoint, ruleLabelID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]VZENNodes, error) {
	var vzenNodes []VZENNodes
	err := common.ReadAllPages(ctx, service.Client, vzenNodeEndpoint, &vzenNodes)
	return vzenNodes, err
}
