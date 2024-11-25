package applications

import (
	"fmt"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	scoreEndpoint   = "/score"
	metricsEndpoint = "/metrics"
)

// Gets the application's ZDX score trend. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetAppScores(service *services.Service, appID int, filters common.GetFromToFilters) ([]common.Metric, *http.Response, error) {
	var v []common.Metric
	var single common.Metric
	path := fmt.Sprintf("%s/%d%s", appsEndpoint, appID, scoreEndpoint)
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

/*
Gets the application's metric trend. For Web Probes, you can access Page Fetch Time, Server Response Time, DNS Time or Availability.
If not specified, it defaults to Page Fetch Time (PFT).
For CloudPath Probes, you can access latency metrics for End to End, Client - Egress, Egress - Application, ZIA Service Edge - Egress, and ZIA Service Edge - Application.
If not specified, it defaults to End to End latency.
If the time range is not specified, the endpoint defaults to the last 2 hours.
*/
// Gets the application's metric trend.
func GetAppMetrics(service *services.Service, appID int, filters common.GetFromToFilters) ([]common.Metric, *http.Response, error) {
	var v []common.Metric
	var single common.Metric
	path := fmt.Sprintf("%s/%d%s", appsEndpoint, appID, metricsEndpoint)
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
