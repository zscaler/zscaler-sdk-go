package alerts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

const (
	baseEndpoint            = "/zdx/v1/alerts"
	ongoingEndpoint         = baseEndpoint + "/ongoing"
	historicalEndpoint      = baseEndpoint + "/historical"
	alertEndpoint           = baseEndpoint + "/%s"
	affectedDevicesEndpoint = baseEndpoint + "/%s/affected_devices"
)

type AlertsResponse struct {
	Alerts     []Alert `json:"alerts"`
	NextOffset string  `json:"next_offset"`
}

type Alert struct {
	ID              int           `json:"id"`
	RuleName        string        `json:"rule_name,omitempty"`
	Severity        string        `json:"severity,omitempty"`
	AlertType       string        `json:"alert_type,omitempty"`
	AlertStatus     string        `json:"alert_status,omitempty"`
	NumGeolocations int           `json:"num_geolocations,omitempty"`
	NumDevices      int           `json:"num_devices,omitempty"`
	StartedOn       int           `json:"started_on,omitempty"`
	EndedOn         int           `json:"ended_on,omitempty"`
	Application     Application   `json:"application,omitempty"`
	Departments     []Department  `json:"departments,omitempty"`
	Locations       []Location    `json:"locations,omitempty"`
	Geolocations    []Geolocation `json:"geolocations,omitempty"`
}

type Application struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Department struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	NumDevices int    `json:"num_devices"`
}

type Geolocation struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	NumDevices int    `json:"num_devices"`
}

type Location struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	NumDevices int     `json:"num_devices"`
	Groups     []Group `json:"groups"`
}

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Device struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	UserID    int    `json:"userid"`
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
}

type AffectedDevicesResponse struct {
	Devices    []Device `json:"devices"`
	NextOffset string   `json:"next_offset"`
}

// GetOngoingAlerts retrieves ongoing alerts with optional filters
func GetOngoingAlerts(ctx context.Context, service *zscaler.Service, filters common.GetFromToFilters) (*AlertsResponse, *http.Response, error) {
	var response AlertsResponse
	resp, err := service.Client.NewRequestDo(ctx, "GET", ongoingEndpoint, filters, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetHistoricalAlerts retrieves historical alerts
// Gets the list of alert history rules defined across an organization.
// All alert history rules are returned if the search filter is not specified.
// The default is set to the previous 2 hours. Alert history rules have an Ended On date.
// The Ended On date is used to pull alerts in conjunction with the provided filters.
// Cannot exceed the 14-day time range limit for alert rules.
func GetHistoricalAlerts(ctx context.Context, service *zscaler.Service, filters common.GetFromToFilters) (*AlertsResponse, *http.Response, error) {
	var response AlertsResponse
	resp, err := service.Client.NewRequestDo(ctx, "GET", historicalEndpoint, filters, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetAlert retrieves a specific alert by ID
func GetAlert(ctx context.Context, service *zscaler.Service, alertID string) (*Alert, *http.Response, error) {
	var response Alert
	path := fmt.Sprintf(alertEndpoint, alertID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetAffectedDevices retrieves the affected devices for a specific alert by ID
func GetAffectedDevices(ctx context.Context, service *zscaler.Service, alertID string, filters common.GetFromToFilters) (*AffectedDevicesResponse, *http.Response, error) {
	var response AffectedDevicesResponse
	path := fmt.Sprintf(affectedDevicesEndpoint, alertID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}
