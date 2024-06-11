package applications

import (
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

const (
	appsEndpoint = "v1/apps"
)

type Apps struct {
	ID              int              `json:"id"`
	Name            string           `json:"name,omitempty"`
	Score           float64          `json:"score,omitempty"`
	MostImpactedGeo *MostImpactedGeo `json:"most_impacted_geo,omitempty"`
	Stats           *Stats           `json:"stats,omitempty"`
}

type MostImpactedGeo struct {
	ID      string `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	GeoType string `json:"geo_type,omitempty"`
}

type Stats struct {
	ActiveUsers   int `json:"active_users"`
	ActiveDevices int `json:"active_devices"`
	NumPoor       int `json:"num_poor"`
	NumOkay       int `json:"num_okay"`
	NumGood       int `json:"num_good"`
}

// Lists all active applications configured for a tenant.
// The endpoint gets each application's ZDX score (default for the last 2 hours), most impacted location, and the total number of users impacted.
// To learn more, see About the ZDX Dashboard at https://help.zscaler.com/zdx/about-zdx-dashboard.
func (service *Service) GetAllApps(filters common.GetFromToFilters) ([]Apps, *http.Response, error) {
	var apps []Apps
	path := appsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &apps)
	if err != nil {
		return nil, nil, err
	}
	return apps, resp, nil
}

func (service *Service) GetApp(appID string, filters common.GetFromToFilters) (*Apps, *http.Response, error) {
	var app Apps
	path := appsEndpoint + "/" + appID
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &app)
	if err != nil {
		return nil, nil, err
	}
	return &app, resp, nil
}
