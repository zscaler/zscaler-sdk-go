package publicapi

import "github.com/zscaler/zscaler-sdk-go/zpa/services/common"

const (
	getDevicesEndpoint        = "/public/v1/getDevices"
	softRemoveDevicesEndpoint = "/public/v1/removeDevices"
	forceRemoveDevicesEndpoin = "/public/v1/forceRemoveDevices"
)

type GetDevices struct {
	AgentVersion            string `json:"agentVersion"`
	CompanyName             string `json:"companyName"`
	ConfigDownloadTime      string `json:"config_download_time"`
	DeregistrationTimestamp string `json:"deregistrationTimestamp"`
	Detail                  string `json:"detail"`
	DownloadCount           int    `json:"download_count"`
	HardwareFingerprint     string `json:"hardwareFingerprint"`
	KeepAliveTime           string `json:"keepAliveTime"`
	LastSeenTime            string `json:"last_seen_time"`
	MacAddress              string `json:"macAddress"`
	MachineHostname         string `json:"machineHostname"`
	Manufacturer            string `json:"manufacturer"`
	OsVersion               string `json:"osVersion"`
	Owner                   string `json:"owner"`
	PolicyName              string `json:"policyName"`
	RegistrationState       string `json:"registrationState"`
	RegistrationTime        string `json:"registration_time"`
	State                   int    `json:"state"`
	TunnelVersion           string `json:"tunnelVersion"`
	Type                    int    `json:"type"`
	Udid                    string `json:"udid"`
	UpmVersion              string `json:"upmVersion"`
	User                    string `json:"user"`
	VpnState                int    `json:"vpnState"`
	ZappArch                string `json:"zappArch"`
}

type RemoveDevicesRequest struct {
	ClientConnectorVersion []string `json:"clientConnectorVersion,omitempty"`
	OsType                 int      `json:"osType,omitempty"`
	Udids                  []string `json:"udids,omitempty"`
	UserName               string   `json:"userName,omitempty"`
}

type RemoveDevicesResponse struct {
	DevicesRemoved int    `json:"devicesRemoved,omitempty"`
	ErrorMsg       string `json:"errorMsg,omitempty"`
}

func (service *Service) GetAll() ([]GetDevices, error) {
	var devices []GetDevices
	_, err := service.Client.NewRequestDo("GET", getDevicesEndpoint, common.Pagination{PageSize: common.DefaultPageSize}, nil, &devices)
	if err != nil {
		return nil, err
	}
	return devices, err
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
	_, err := service.Client.NewRequestDo("POST", forceRemoveDevicesEndpoin, common.Pagination{PageSize: common.DefaultPageSize}, &request, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}
