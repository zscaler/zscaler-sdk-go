package remove_devices

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
)

const (
	softRemoveDevicesEndpoint  = "/zcc/papi/public/v1/removeDevices"
	forceRemoveDevicesEndpoint = "/zcc/papi/public/v1/forceRemoveDevices"
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
func SoftRemoveDevices(service *zscaler.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", softRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewZccRequestDo("POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// ForceRemoveDevices force removes the enrolled devices from the portal
func ForceRemoveDevices(service *zscaler.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", forceRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo("POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}
