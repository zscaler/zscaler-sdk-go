package devices

import (
	"net/http"
)

const (
	geoLocationsEndpoint = "v1/active_geo"
)

type GeoLocation struct {
	ID          int           `json:"id"`
	Name        string        `json:"name,omitempty"`
	GeoType     string        `json:"geo_type,omitempty"`
	Description string        `json:"description,omitempty"`
	Children    []GeoLocation `json:"children,omitempty"`
}

// Gets the list of all active geolocations for the time range specified.
// If not specified, the endpoint defaults to the last 2 hours. The state and city data is retrieved only for the US.
func (service *Service) GetGeoLocations(filters GeoLocationFilter) (*GeoLocation, *http.Response, error) {
	v := new(GeoLocation)
	path := geoLocationsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
