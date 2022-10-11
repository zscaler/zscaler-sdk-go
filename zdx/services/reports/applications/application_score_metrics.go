package applications

import "github.com/zscaler/zscaler-sdk-go/zdx/services/common"

const (
	scoreEndpoint   = "/score"
	metricsEndpoint = "/metrics"
)

/*
https://help.zscaler.com/zdx/reports#/apps/{appid}/score-get
Gets the application’s ZDX score trend. If the time range is not specified, the endpoint defaults to the last 2 hours.

https://help.zscaler.com/zdx/reports#/apps/{appid}/metrics-get
Gets the application’s metric trend. For Web Probes, you can access Page Fetch Time, Server Response Time, DNS Time or Availability.
If not specified, it defaults to Page Fetch Time (PFT).
For CloudPath Probes, you can access latency metrics for End to End, Client - Egress, Egress - Application, ZIA Service Edge - Egress, and ZIA Service Edge - Application.
If not specified, it defaults to End to End latency.
If the time range is not specified, the endpoint defaults to the last 2 hours.
*/

type ModelSeries struct {
	Metric     string              `json:"metric,omitempty"`
	Unit       string              `json:"unit,omitempty"`
	DataPoints []common.DataPoints `json:"datapoints"`
}
