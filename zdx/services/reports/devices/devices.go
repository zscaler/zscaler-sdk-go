package devices

import (
	"fmt"
	"net/http"
)

const (
	devicesEndpoint = "devices"
)

type DeviceDetail struct {
	ID       int       `json:"id"`
	Name     string    `json:"name,omitempty"`
	Hardware *Hardware `json:"hardware,omitempty"`
	Network  []Network `json:"network,omitempty"`
	Software *Software `json:"software,omitempty"`
}

type Hardware struct {
	HWModel     string `json:"hw_model,omitempty"`
	HWMFG       string `json:"hw_mfg,omitempty"`
	HWType      string `json:"hw_type,omitempty"`
	HWSerial    string `json:"hw_serial,omitempty"`
	TotMem      string `json:"tot_mem,omitempty"`
	GPU         string `json:"gpu,omitempty"`
	DiskSize    string `json:"disk_size,omitempty"`
	DiskModel   string `json:"disk_model,omitempty"`
	DiskType    string `json:"disk_type,omitempty"`
	CPUMFG      string `json:"cpu_mfg,omitempty"`
	CPUModel    string `json:"cpu_model,omitempty"`
	SpeedGHZ    int    `json:"speed_ghz,omitempty"`
	LogicalProc int    `json:"logical_proc,omitempty"`
	NumCores    int    `json:"num_cores,omitempty"`
}

type Network struct {
	NetType     string `json:"net_type,omitempty"`
	Status      string `json:"status,omitempty"`
	IPv4        string `json:"ipv4,omitempty"`
	IPv6        string `json:"ipv6,omitempty"`
	DNSSRVS     string `json:"dns_srvs,omitempty"`
	DNSSuffix   string `json:"dns_suffix,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	MAC         string `json:"mac,omitempty"`
	GUID        string `json:"guid,omitempty"`
	WiFiAdapter string `json:"wifi_adapter,omitempty"`
	WiFiType    string `json:"wifi_type,omitempty"`
	SSID        string `json:"ssid,omitempty"`
	Channel     string `json:"channel,omitempty"`
	BSSID       string `json:"bssid,omitempty"`
}

type Software struct {
	OSName        string `json:"os_name,omitempty"`
	OSVer         string `json:"os_ver,omitempty"`
	Hostname      string `json:"hostname,omitempty"`
	NetBios       string `json:"netbios,omitempty"`
	User          string `json:"user,omitempty"`
	ClientConnVer string `json:"client_conn_ver,omitempty"`
	ZDXVer        string `json:"zdx_ver,omitempty"`
}

// Gets the device details including the device model information, tunnel type, network, and software details. The JSON must contain the user ID and email address to associate the device to a user. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) Get(deviceID string) (*DeviceDetail, *http.Response, error) {
	v := new(DeviceDetail)
	path := fmt.Sprintf("%v/%v", devicesEndpoint, deviceID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list of all active devices and its basic details. The JSON must contain the user's ID and email address to associate the device to the user. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetAll(filters GetDevicesFilters) ([]DeviceDetail, *http.Response, error) {
	var v struct {
		NextOffSet interface{}    `json:"next_offset"`
		List       []DeviceDetail `json:"devices"`
	}

	relativeURL := devicesEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v.List, resp, nil
}
