package alerts

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services"
)

const (
	baseEndpoint            = "v1/alerts"
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

// GetOngoingAlerts retrieves ongoing alerts
func GetOngoingAlerts(service *services.Service) (*AlertsResponse, *http.Response, error) {
	var response AlertsResponse
	resp, err := service.Client.NewRequestDo("GET", ongoingEndpoint, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetHistoricalAlerts retrieves historical alerts
func GetHistoricalAlerts(service *services.Service) (*AlertsResponse, *http.Response, error) {
	var response AlertsResponse
	resp, err := service.Client.NewRequestDo("GET", historicalEndpoint, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetAlert retrieves a specific alert by ID
func GetAlert(service *services.Service, alertID string) (*Alert, *http.Response, error) {
	var response Alert
	path := fmt.Sprintf(alertEndpoint, alertID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}

// GetAffectedDevices retrieves the affected devices for a specific alert by ID
func GetAffectedDevices(service *services.Service, alertID string) (*AlertsResponse, *http.Response, error) {
	var response AlertsResponse
	path := fmt.Sprintf(affectedDevicesEndpoint, alertID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, err
	}
	return &response, resp, nil
}
