package geo_coordinates

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	geoCoordinatesEndpoint = "/zia/api/v1/region/byGeoCoordinates"
)

type GeoCoordinates struct {
	// The geographical ID of the city
	CityGeoId int `json:"cityGeoId"`

	// The geographical ID of the state
	StateGeoId int `json:"stateGeoId"`

	// The latitude coordinate of the city
	Latitude float64 `json:"latitude"`

	// The longitude coordinate of the city
	Longitude float64 `json:"longitude"`

	// The name of the city
	CityName string `json:"cityName"`

	// The name of the state, province, or territory of a country
	StateName string `json:"stateName"`

	// The name of the country
	CountryName string `json:"countryName"`

	// The ISO standard two-letter country code
	CountryCode string `json:"countryCode"`

	// The postal code
	PostalCode string `json:"postalCode"`

	// The ISO standard two-letter continent code
	ContinentCode string `json:"continentCode"`
}

func GetByGeoCoordinates(service *zscaler.Service, latitude, longitude float64) (*GeoCoordinates, error) {
	var region GeoCoordinates
	queryParams := fmt.Sprintf("latitude=%f&longitude=%f", latitude, longitude)
	fullEndpoint := fmt.Sprintf("%s?%s", geoCoordinatesEndpoint, queryParams)

	err := service.Client.Read(fullEndpoint, &region)
	if err != nil {
		return nil, err
	}
	return &region, nil
}
