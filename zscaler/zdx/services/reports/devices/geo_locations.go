package devices

import (
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services"
)

const (
	geoLocationsEndpoint = "v1/active_geo"
)

type GeoLocation struct {
	ID          string     `json:"id"`
	Name        string     `json:"name,omitempty"`
	GeoType     string     `json:"geo_type,omitempty"`
	Description string     `json:"description,omitempty"`
	Children    []Children `json:"children,omitempty"`
}

type Children struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	GeoType     string `json:"geo_type,omitempty"`
}

// Gets the list of all active geolocations for the time range specified.
// If not specified, the endpoint defaults to the last 2 hours. The state and city data is retrieved only for the US.
func GetGeoLocations(service *services.Service, filters GeoLocationFilter) ([]GeoLocation, *http.Response, error) {
	var v []GeoLocation
	path := geoLocationsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
