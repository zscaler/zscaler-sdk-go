package users

import (
	"fmt"
	"net/http"
)

const (
	usersEndpoint = "users"
)

type User struct {
	ID      int       `json:"id"`
	Name    string    `json:"name,omitempty"`
	Email   string    `json:"email,omitempty"`
	Devices []Devices `json:"devices,omitempty"`
}

type Devices struct {
	ID           int            `json:"id"`
	Name         string         `json:"name,omitempty"`
	UserLocation []UserLocation `json:"geo_loc,omitempty"`
	ZSLocation   []ZSLocation   `json:"zs_loc,omitempty"`
}

type UserLocation struct {
	ID      int    `json:"id"`
	City    string `json:"city,omitempty"`
	Region  string `json:"region,omitempty"`
	Country string `json:"country,omitempty"`
	GeoType string `json:"geo_type,omitempty"`
}

type ZSLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// Gets user details including the device information, active geolocations, and Zscaler locations. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) Get(userID string) (*User, *http.Response, error) {
	v := new(User)
	path := fmt.Sprintf("%v/%v", usersEndpoint, userID)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list of all active users, their devices, active geolocations, and Zscaler locations. If the time range is not specified, the endpoint defaults to the last 2 hours.
func (service *Service) GetAll(filters GetUsersFilters) ([]User, *http.Response, error) {
	var v struct {
		NextOffSet interface{} `json:"next_offset"`
		List       []User      `json:"users"`
	}

	relativeURL := usersEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v.List, resp, nil
}
