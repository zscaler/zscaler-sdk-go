package vzen_clusters

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
	vzenClusterEndpoint = "/zia/api/v1/virtualZenClusters"
)

type VZENClusters struct {
	// System-generated Virtual Service Edge cluster ID
	ID int `json:"id,omitempty"`

	// Name of the Virtual Service Edge cluster
	Name string `json:"name,omitempty"`

	// Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default.
	// Supported Values: ENABLED, DISABLED, DISABLED_BY_SERVICE_PROVIDER, NOT_PROVISIONED_IN_SERVICE_PROVIDER
	Status string `json:"status,omitempty"`

	// The Virtual Service Edge cluster IP address.
	// In a Virtual Service Edge cluster, the cluster IP address provides fault tolerance and is used to listen for user traffic.
	// This interface doesn't explicitly get an IP address.
	// The cluster IP address must be in the same VLAN as the proxy and load balancer IP addresses.
	IpAddress string `json:"ipAddress,omitempty"`

	// The Virtual Service Edge cluster subnet mask
	SubnetMask string `json:"subnetMask,omitempty"`

	// The IP address of the default gateway to the internet
	DefaultGateway string `json:"defaultGateway,omitempty"`

	// The Virtual Service Edge cluster type
	// See https://help.zscaler.com/zia/service-edges#/virtualZenClusters-post for all types
	Type string `json:"type,omitempty"`

	// A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
	IpSecEnabled bool `json:"ipSecEnabled,omitempty"`

	// The Virtual Service Edge instances you want to include in the cluster.
	// A Virtual Service Edge cluster must contain at least two Virtual Service Edge instances.
	VirtualZenNodes []common.IDNameExternalID `json:"virtualZenNodes,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, clusterID int) (*VZENClusters, error) {
	var vzenCluster VZENClusters
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", vzenClusterEndpoint, clusterID), &vzenCluster)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning vzen cluster from Get: %d", vzenCluster.ID)
	return &vzenCluster, nil
}

func GetClusterByName(ctx context.Context, service *zscaler.Service, clusterName string) (*VZENClusters, error) {
	var vzenClusters []VZENClusters
	err := common.ReadAllPages(ctx, service.Client, vzenClusterEndpoint, &vzenClusters)
	if err != nil {
		return nil, err
	}
	for _, vzenCluster := range vzenClusters {
		if strings.EqualFold(vzenCluster.Name, clusterName) {
			return &vzenCluster, nil
		}
	}
	return nil, fmt.Errorf("no vzen cluster found with name: %s", clusterName)
}

func Create(ctx context.Context, service *zscaler.Service, vzenClusters *VZENClusters) (*VZENClusters, *http.Response, error) {
	resp, err := service.Client.Create(ctx, vzenClusterEndpoint, *vzenClusters)
	if err != nil {
		return nil, nil, err
	}

	createdCluster, ok := resp.(*VZENClusters)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a vzen cluster pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new vzen cluster from create: %d", createdCluster.ID)
	return createdCluster, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, clusterID int, vzenClusters *VZENClusters) (*VZENClusters, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", vzenClusterEndpoint, clusterID), *vzenClusters)
	if err != nil {
		return nil, nil, err
	}
	updatedCluster, _ := resp.(*VZENClusters)

	service.Client.GetLogger().Printf("[DEBUG]returning updates vzen cluster from update: %d", updatedCluster.ID)
	return updatedCluster, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, clusterID int) (*http.Response, error) {
	// Step 1: Fetch the current cluster
	cluster, err := Get(ctx, service, clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve VZEN cluster with ID %d: %w", clusterID, err)
	}

	// Step 2: Prepare minimal update payload with empty VirtualZenNodes to detach
	updatePayload := &VZENClusters{
		ID:              cluster.ID,
		Name:            cluster.Name,
		Status:          cluster.Status,
		Type:            cluster.Type,
		IpAddress:       cluster.IpAddress,
		SubnetMask:      cluster.SubnetMask,
		DefaultGateway:  cluster.DefaultGateway,
		IpSecEnabled:    cluster.IpSecEnabled,
		VirtualZenNodes: []common.IDNameExternalID{}, // Detach all service edges
	}

	// Step 3: Update the cluster to remove associations
	_, _, err = Update(ctx, service, clusterID, updatePayload)
	if err != nil {
		return nil, fmt.Errorf("failed to update VZEN cluster %d before deletion: %w", clusterID, err)
	}

	// Step 4: Delete the cluster
	err = service.Client.Delete(ctx, fmt.Sprintf("%s/%d", vzenClusterEndpoint, clusterID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete VZEN cluster %d: %w", clusterID, err)
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]VZENClusters, error) {
	var vzenClusters []VZENClusters
	err := common.ReadAllPages(ctx, service.Client, vzenClusterEndpoint, &vzenClusters)
	return vzenClusters, err
}
