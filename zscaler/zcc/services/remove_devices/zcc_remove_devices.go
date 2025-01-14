package remove_devices

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	softRemoveDevicesEndpoint   = "/zcc/papi/public/v1/removeDevices"
	forceRemoveDevicesEndpoint  = "/zcc/papi/public/v1/forceRemoveDevices"
	removeMachineTunnelEndpoint = "/zcc/papi/public/v1/removeMachineTunnel"
)

type RemoveDevicesResponse struct {
	DevicesRemoved int    `json:"devicesRemoved,omitempty"`
	ErrorMsg       string `json:"errorMsg,omitempty"`
}

type RemoveDevicesRequest struct {
	ClientConnectorVersion []string `json:"clientConnectorVersion,omitempty"`
	OsType                 int      `json:"osType,omitempty"`
	Udids                  []string `json:"udids,omitempty"`
	UserName               string   `json:"userName,omitempty"`
}

// SoftRemoveDevices soft removes the enrolled devices from the portal
func SoftRemoveDevices(ctx context.Context, service *zscaler.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", softRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewZccRequestDo(ctx, "POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// ForceRemoveDevices force removes the enrolled devices from the portal
func ForceRemoveDevices(ctx context.Context, service *zscaler.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", forceRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo(ctx, "POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// RemoveMachineTunnel sends a request to remove machine tunnels for the specified hostnames or machine tokens.
func RemoveMachineTunnel(ctx context.Context, service *zscaler.Service, hostNames []string, machineToken string) (*RemoveDevicesResponse, error) {
	// Validate input
	if len(hostNames) == 0 && machineToken == "" {
		return nil, fmt.Errorf("either hostNames or machineToken must be provided")
	}

	// Construct request payload
	payload := map[string]interface{}{}
	if len(hostNames) > 0 {
		payload["hostName"] = strings.Join(hostNames, ",") // Ensure hostnames are joined as a comma-separated string
	}
	if machineToken != "" {
		payload["machineToken"] = machineToken
	}

	// Make the request
	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo(ctx, "POST", removeMachineTunnelEndpoint, nil, payload, &response)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	return &response, nil
}
