package remove_devices

import "github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"

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

func (service *Service) SoftRemoveDevices(request RemoveDevicesRequest) (*RemoveDevicesResponse, error) {
	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo("POST", softRemoveDevicesEndpoint, common.Pagination{PageSize: common.DefaultPageSize}, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (service *Service) ForceRemoveDevices(request RemoveDevicesRequest) (*RemoveDevicesResponse, error) {
	var response RemoveDevicesResponse
	_, err := service.Client.NewRequestDo("POST", forceRemoveDevicesEndpoint, common.Pagination{PageSize: common.DefaultPageSize}, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}
