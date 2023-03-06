package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zdx/services/common"
)

const (
	deviceCloudPathEndpoint = "v1/cloudpath"
)

type DeviceCloudPath struct {
	TimeStamp int         `json:"timestamp,omitempty"`
	CloudPath []CloudPath `json:"cloudpath,omitempty"`
}

type CloudPath struct {
	SRC           string `json:"src,omitempty"`
	DST           string `json:"dst,omitempty"`
	NumHops       int    `json:"num_hops,omitempty"`
	Latency       int    `json:"latency,omitempty"`
	Loss          int    `json:"loss,omitempty"`
	NumUnrespHops int    `json:"num_unresp_hops,omitempty"`
	TunnelType    int    `json:"tunnel_type,omitempty"`
	Hops          []Hops `json:"hops,omitempty"`
}

type Hops struct {
	IP          string `json:"ip,omitempty"`
	GWMac       string `json:"gw_mac,omitempty"`
	GWMacVendor string `json:"gw_mac_vendor,omitempty"`
	PkgSent     int    `json:"pkt_sent,omitempty"`
	PkgRcvd     int    `json:"pkt_rcvd,omitempty"`
	LatencyMin  int    `json:"latency_min,omitempty"`
	LatencyMax  int    `json:"latency_max,omitempty"`
	LatencyAvg  int    `json:"latency_avg,omitempty"`
	LatencyDiff int    `json:"latency_diff,omitempty"`
}

// Gets the Cloud Path hop data for an application on a specific device.
// Includes the summary data for the entire path like the total number of hops, packet loss, latency, and tunnel type (if available).
// It also includes a similar summary of data for each individual hop.
// If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetDeviceCloudPath(deviceID, appID, probeID string, filters common.GetFromToFilters) (*DeviceCloudPath, *http.Response, error) {
	v := new(DeviceCloudPath)
	path := fmt.Sprintf("%v/%v/%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathEndpoint, probeID, deviceCloudPathEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
