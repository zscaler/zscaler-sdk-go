package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

const (
	deviceCloudPathProbesEndpoint = "cloudpath-probes"
)

// devices/{deviceid}/apps/{appid}/cloudpath-probes
type DeviceCloudPathProbe struct {
	ID             int              `json:"id,omitempty"`
	Name           string           `json:"name,omitempty"`
	NumProbes      int              `json:"num_probes,omitempty"`
	AverageLatency []AverageLatency `json:"avg_latencies,omitempty"`
}

type AverageLatency struct {
	LegSRC  string  `json:"leg_src,omitempty"`
	LegDst  string  `json:"leg_dst,omitempty"`
	Latency float32 `json:"latency,omitempty"`
}

// /devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}
type NetworkStats struct {
	LegSRC string          `json:"leg_src,omitempty"`
	LegDst string          `json:"leg_dst,omitempty"`
	Stats  []common.Metric `json:"stats,omitempty"`
}

// /devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}/cloudpath
type CloudPathProbe struct {
	TimeStamp int         `json:"timestamp,omitempty"`
	CloudPath []CloudPath `json:"cloudpath,omitempty"` // Changed to a slice
}

type CloudPath struct {
	SRC           string  `json:"src,omitempty"`
	DST           string  `json:"dst,omitempty"`
	NumHops       int     `json:"num_hops,omitempty"`
	Latency       float32 `json:"latency,omitempty"`
	Loss          float32 `json:"loss,omitempty"`
	NumUnrespHops int     `json:"num_unresp_hops,omitempty"`
	TunnelType    int     `json:"tunnel_type,omitempty"`
	Hops          []Hops  `json:"hops,omitempty"` // Changed to a slice
}

type Hops struct {
	IP          string `json:"ip,omitempty"`
	GWMac       string `json:"gw_mac,omitempty"`
	GWMacVendor string `json:"gw_mac_vendor,omitempty"`
	PktSent     int    `json:"pkt_sent,omitempty"`
	PktRcvd     int    `json:"pkt_rcvd,omitempty"`
	LatencyMin  int    `json:"latency_min,omitempty"`
	LatencyMax  int    `json:"latency_max,omitempty"`
	LatencyAvg  int    `json:"latency_avg,omitempty"`
	LatencyDiff int    `json:"latency_diff,omitempty"`
}

// devices/{deviceid}/apps/{appid}/cloudpath-probes
// Gets the list of all active Cloud Path probes on a device.
// If the time range is not specified, the default is the previous 2 hours
func GetAllCloudPathProbes(service *services.Service, deviceID, appID int, filters common.GetFromToFilters) ([]DeviceCloudPathProbe, *http.Response, error) {
	var v []DeviceCloudPathProbe
	path := fmt.Sprintf("%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathProbesEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v) // Pass the address of v
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// /devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}/cloudpath
// Gets the Web probe's Page Fetch Time (PFT) on a device for an application. If the time range is not specified, the endpoint defaults to the previous 2 hours.
func GetDeviceAppCloudPathProbe(service *services.Service, deviceID, appID, probeID int, filters common.GetFromToFilters) ([]NetworkStats, *http.Response, error) {
	var v []NetworkStats
	path := fmt.Sprintf("%v/%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathProbesEndpoint, probeID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v) // Pass the address of v
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// /devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}/cloudpath
// Gets the Cloud Path hop data for an application on a specific device.
// Includes the summary data for the entire path like the total number of hops, packet loss, latency, and tunnel type (if available).
// It also includes a similar summary of data for each individual hop. If the time range is not specified, the endpoint defaults to the previous 2 hours.
func GetCloudPathAppDevice(service *services.Service, deviceID, appID, probeID int, filters common.GetFromToFilters) ([]CloudPathProbe, *http.Response, error) {
	var v []CloudPathProbe
	path := fmt.Sprintf("%v/%v/%v/%v/%v/%v/%s", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathProbesEndpoint, probeID, "/cloudpath")
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v) // Pass the address of v
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
