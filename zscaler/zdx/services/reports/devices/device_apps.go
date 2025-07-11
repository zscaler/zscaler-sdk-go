package devices

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	deviceAppsEndpoint = "apps"
)

type App struct {
	ID    int     `json:"id,omitempty"`
	Name  string  `json:"name,omitempty"`
	Score float32 `json:"score,omitempty"`
}

// Gets the application's ZDX score trend for a device. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetDeviceApp(ctx context.Context, service *zscaler.Service, deviceID, appID string, filters common.GetFromToFilters) (*App, *http.Response, error) {
	v := new(App)
	path := fmt.Sprintf("%v/%v/%v/%v", devicesEndpoint, deviceID, deviceAppsEndpoint, appID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list all active applications for a device. The endpoint gets the ZDX score each application. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetDeviceAllApps(ctx context.Context, service *zscaler.Service, deviceID string, filters common.GetFromToFilters) ([]App, *http.Response, error) {
	var v []App
	relativeURL := devicesEndpoint + "/" + deviceID + "/" + deviceAppsEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
