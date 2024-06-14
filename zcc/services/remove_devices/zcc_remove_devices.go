package remove_devices

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zcc/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zcc/services/common"
)

const (
	softRemoveDevicesEndpoint  = "/public/v1/removeDevices"
	forceRemoveDevicesEndpoint = "/public/v1/forceRemoveDevices"
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
func SoftRemoveDevices(service *services.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", softRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo("POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

// ForceRemoveDevices force removes the enrolled devices from the portal
func ForceRemoveDevices(service *services.Service, request RemoveDevicesRequest, pageSize int) (*RemoveDevicesResponse, error) {
	pagination := common.NewPagination(pageSize)
	fullURL := fmt.Sprintf("%s?pageSize=%d", forceRemoveDevicesEndpoint, pagination.PageSize)

	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo("POST", fullURL, nil, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}
