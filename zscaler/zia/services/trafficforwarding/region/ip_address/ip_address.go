package ip_address

import (
	"context"
	"fmt"
	"net/url"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	byIPAdddressEndpoint = "/zia/api/v1/region/byIPAddress"
)

type ByIPAddress struct {
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

func GetByIPAddress(ctx context.Context, service *zscaler.Service, ipAddress string) (*ByIPAddress, error) {
	var ip ByIPAddress
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", byIPAdddressEndpoint, url.QueryEscape(ipAddress)), &ip)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}
