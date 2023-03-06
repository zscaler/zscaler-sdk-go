package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zdx/services/common"
)

const (
	deviceCloudPathProbesEndpoint = "cloudpath-probes"
)

type DeviceCloudPathProbe struct {
	ID             int              `json:"id,omitempty"`
	Name           string           `json:"name,omitempty"`
	NumProbes      int              `json:"num_probes,omitempty"`
	AverageLatency []AverageLatency `json:"avg_latencies,omitempty"`
}

type AverageLatency struct {
	LegSRC  string `json:"leg_src,omitempty"`
	LegDst  string `json:"leg_dst,omitempty"`
	Latency int    `json:"latency,omitempty"`
}

type DeviceCloudPathProbeMetric struct {
	LegSRC string          `json:"leg_src,omitempty"`
	LegDst string          `json:"leg_dst,omitempty"`
	Stats  []common.Metric `json:"stats,omitempty"`
}

// Gets the CloudPath Probe's metric trend on a device for an application.
// For Cloud Path Probes, you can access latency metrics for End to End, Client - Egress, Egress - Application, ZIA Service Edge- Egress, and ZIA Service Edge - Application.
// If not specified, it defaults to End to End latency.
// If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetCloudPathProbesMetric(deviceID, appID, probeID string, filters common.GetFromToFilters) (*DeviceCloudPathProbeMetric, *http.Response, error) {
	v := new(DeviceCloudPathProbeMetric)
	path := fmt.Sprintf("%v/%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathProbesEndpoint, probeID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list of all active Cloud Path probes on a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetAllCloudPathProbes(deviceID, appID string, filters common.GetFromToFilters) ([]DeviceCloudPathProbe, *http.Response, error) {
	var v []DeviceCloudPathProbe
	path := fmt.Sprintf("%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceCloudPathProbesEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
