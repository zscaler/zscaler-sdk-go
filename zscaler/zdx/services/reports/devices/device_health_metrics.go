package devices

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
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

// /devices/42827781/apps/1/health-metrics?from=1718247199&to=1718254399
// Gets the health metrics trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
// The health metrics include CPU, Memory, Disk I/O, Network I/O, Wi-Fi, Network Bandwidth, etc.
func GetHealthMetrics(service *services.Service, deviceID int, filters common.GetFromToFilters) ([]HealthMetrics, *http.Response, error) {
	var v []HealthMetrics
	path := fmt.Sprintf("%v/%v/%v", devicesEndpoint, deviceID, deviceHealthMetricsEndpoint)
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
