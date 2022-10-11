package devices

import "github.com/zscaler/zscaler-sdk-go/zdx/services/common"

const (
	deviceQualityMetricsEndpoint = "/call-quality-metrics"
)

/*
https://help.zscaler.com/zdx/reports#/devices/{deviceid}/apps/{appid}/cloudpath-probes/{probeid}/cloudpath-get
Gets the Cloud Path hop data for an application on a specific device.
Includes the summary data for the entire path like the total number of hops, packet loss, latency, and tunnel type (if available).
It also includes a similar summary of data for each individual hop. If the time range is not specified, the endpoint defaults to the last 2 hours.
*/

type CallQualityMetrics struct {
	MeetID        string    `json:"meet_id,omitempty"`
	MeetSessionID string    `json:"meet_session_id,omitempty"`
	MeetSubject   string    `json:"meet_subject,omitempty"`
	Metrics       []Metrics `json:"metrics,omitempty"`
}

type Metrics struct {
	Metric     string              `json:"metric,omitempty"`
	Unit       string              `json:"unit,omitempty"`
	DataPoints []common.DataPoints `json:"datapoints"`
}
