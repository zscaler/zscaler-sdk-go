package devices

import "github.com/zscaler/zscaler-sdk-go/zdx/services/common"

const (
	deviceHealthMetricsEndpoint = "/health-metrics"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/health-metrics-get
Gets the health metrics trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
The health metrics include CPU, Memory, Disk I/O, Network I/O, Wi-Fi, Network Bandwidth, etc.
*/

type HealthMetrics struct {
	Category  string      `json:"category,omitempty"`
	Instances []Instances `json:"instances,omitempty"`
}

type Instances struct {
	Name    string    `json:"metric,omitempty"`
	Metrics []Metrics `json:"metrics,omitempty"`
}

type MetricSeries struct {
	Metric     string              `json:"metric,omitempty"`
	Unit       string              `json:"unit,omitempty"`
	DataPoints []common.DataPoints `json:"datapoints"`
}
