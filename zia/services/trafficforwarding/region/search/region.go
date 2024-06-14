package region

import (
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
)

const (
	regionSearchEndpoint = "/region/search"
)

type Regions struct {
	// The geographical ID of the city
	Datacenter int `json:"cityGeoId"`

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

func GetDatacenterRegion(service *services.Service, regionPrefix string) ([]Regions, error) {
	var regions []Regions
	err := service.Client.Read(fmt.Sprintf("%s?prefix=%s", regionSearchEndpoint, url.QueryEscape(regionPrefix)), &regions)
	if err != nil {
		return nil, err
	}
	return regions, nil
}
