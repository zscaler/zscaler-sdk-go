package devices

import (
	"fmt"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	deviceWebProbesEndpoint = "web-probes"
)

type DeviceWebProbe struct {
	ID        int     `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	NumProbes int     `json:"num_probes,omitempty"`
	AvgScore  float32 `json:"avg_score,omitempty"`
	AvgPFT    float32 `json:"avg_pft,omitempty"`
}

func generateWebProbesPath(deviceID, appID int) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceWebProbesEndpoint)
}

func generateWebProbePath(deviceID, appID, probeID int) string {
	return fmt.Sprintf("%v/%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceWebProbesEndpoint, probeID)
}

// Gets the Web Probe metrics trend on a device for an application.
// For Web Probes, you can access Page Fetch Time, Server Response Time, DNS Time, or Availability.
// If not specified, it defaults to Page Fetch Time (PFT).
// If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetWebProbes(service *services.Service, deviceID, appID, probeID int, filters common.GetFromToFilters) ([]common.Metric, *http.Response, error) {
	var v []common.Metric
	var single common.Metric
	path := generateWebProbePath(deviceID, appID, probeID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v)
	if err == nil {
		return v, resp, nil
	}

	// If unmarshalling to an array fails, try unmarshalling to a single object
	resp, err = service.Client.NewRequestDo("GET", path, filters, nil, &single)
	if err == nil {
		v = append(v, single)
		return v, resp, nil
	}

	return nil, nil, err
}

// Gets the list of all active web probes on a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetAllWebProbes(service *services.Service, deviceID, appID int, filters common.GetFromToFilters) ([]DeviceWebProbe, *http.Response, error) {
	var v []DeviceWebProbe
	path := generateWebProbesPath(deviceID, appID)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v) // Pass the address of v
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
