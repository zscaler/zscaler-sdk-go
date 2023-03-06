package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zdx/services/common"
)

const (
	deviceHealthMetricsEndpoint = "health-metrics"
)

type HealthMetrics struct {
	Category  string      `json:"category,omitempty"`
	Instances []Instances `json:"instances,omitempty"`
}

type Instances struct {
	Name    string          `json:"metric,omitempty"`
	Metrics []common.Metric `json:"metrics,omitempty"`
}

// Gets the health metrics trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
// The health metrics include CPU, Memory, Disk I/O, Network I/O, Wi-Fi, Network Bandwidth, etc.
func (service *Service) GetHealthMetrics(deviceID string, filters common.GetFromToFilters) (*HealthMetrics, *http.Response, error) {
	v := new(HealthMetrics)
	path := fmt.Sprintf("%v/%v/%v", devicesEndpoint, deviceID, deviceHealthMetricsEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
