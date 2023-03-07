package applications

import (
	"net/http"
)

const (
	appsEndpoint = "v1/apps"
)

type Apps struct {
	ID              int              `json:"id"`
	Name            string           `json:"name,omitempty"`
	Score           int              `json:"score,omitempty"`
	MostImpactedGeo *MostImpactedGeo `json:"most_impacted_geo,omitempty"`
	Stats           *Stats           `json:"stats,omitempty"`
}

type MostImpactedGeo struct {
	ID      int    `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	State   string `json:"state,omitempty"`
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
func (service *Service) GetAllApps(filters GetAppsFilters) (*Apps, *http.Response, error) {
	v := new(Apps)
	path := appsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the application’s ZDX score (for the last 2 hours), most impacted locations, and the total number of users impacted.
// To learn more, see About the ZDX Score at https://help.zscaler.com/zdx/about-zdx-score#user.
func (service *Service) GetApp(appID string, filters GetAppsFilters) (*Apps, *http.Response, error) {
	v := new(Apps)
	path := appsEndpoint + "/" + appID
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
