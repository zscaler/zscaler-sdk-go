package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zdx/services/common"
)

const (
	deviceWebProbesEndpoint = "web-probes"
)

type DeviceWebProbe struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	NumProbes int    `json:"num_probes,omitempty"`
	AvgScore  int    `json:"avg_score,omitempty"`
	AvgPFT    int    `json:"avg_pft,omitempty"`
}

func generateWebProbesPath(deviceID, appID string) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceWebProbesEndpoint)
}

func generateWebProbePath(deviceID, appID, probeID string) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceWebProbesEndpoint, probeID)
}

// Gets the Web Probe metrics trend on a device for an application.
// For Web Probes, you can access Page Fetch Time, Server Response Time, DNS Time, or Availability.
// If not specified, it defaults to Page Fetch Time (PFT).
// If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetWebProbes(deviceID, appID, probeID string, filters common.GetFromToFilters) (*common.Metric, *http.Response, error) {
	v := new(common.Metric)
	path := generateWebProbePath(deviceID, appID, probeID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list of all active web probes on a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetAllWebProbes(deviceID, appID string, filters common.GetFromToFilters) ([]DeviceWebProbe, *http.Response, error) {
	var v []DeviceWebProbe
	path := generateWebProbesPath(deviceID, appID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
