package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	deviceQualityMetricsEndpoint = "call-quality-metrics"
)

type CallQualityMetrics struct {
	MeetID        string          `json:"meet_id,omitempty"`
	MeetSessionID string          `json:"meet_session_id,omitempty"`
	MeetSubject   string          `json:"meet_subject,omitempty"`
	Metrics       []common.Metric `json:"metrics,omitempty"`
}

// Gets the Call Quality metric trend for a device for a CQM application.
// If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetQualityMetrics(ctx context.Context, service *zscaler.Service, deviceID, appID int, filters common.GetFromToFilters) ([]CallQualityMetrics, *http.Response, error) {
	var v []CallQualityMetrics
	path := fmt.Sprintf("%v/%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID, deviceQualityMetricsEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
