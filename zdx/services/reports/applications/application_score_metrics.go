package applications

import (
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

const (
	scoreEndpoint   = "v1/score"
	metricsEndpoint = "v1/metrics"
)

// Gets the application's ZDX score trend. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetAppScores(appID string, filters GetAppsFilters) (*common.Metric, *http.Response, error) {
	v := new(common.Metric)
	path := appsEndpoint + "/" + appID + scoreEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

/*
Gets the application's metric trend. For Web Probes, you can access Page Fetch Time, Server Response Time, DNS Time or Availability.
If not specified, it defaults to Page Fetch Time (PFT).
For CloudPath Probes, you can access latency metrics for End to End, Client - Egress, Egress - Application, ZIA Service Edge - Egress, and ZIA Service Edge - Application.
If not specified, it defaults to End to End latency.
If the time range is not specified, the endpoint defaults to the last 2 hours.
*/
func (service *Service) GetAppMetrics(appID string, filters GetAppsFilters) (*common.Metric, *http.Response, error) {
	v := new(common.Metric)
	path := appsEndpoint + "/" + appID + metricsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
