package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	usersEndpoint = "/zdx/v1/users"
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
	ID           string  `json:"id"`
	City         string  `json:"city,omitempty"`
	State        string  `json:"state,omitempty"`
	Country      string  `json:"country,omitempty"`
	GeoLat       float32 `json:"geo_lat,omitempty"`
	GeoLong      float32 `json:"geo_long,omitempty"`
	GeoDetection string  `json:"geo_detection,omitempty"`
}

type ZSLocation struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// Gets user details including the device information, active geolocations, and Zscaler locations. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetUser(ctx context.Context, service *zscaler.Service, userID string) (*User, *http.Response, error) {
	v := new(User)
	path := fmt.Sprintf("%v/%v", usersEndpoint, userID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets the list of all active users, their devices, active geolocations, and Zscaler locations. If the time range is not specified, the endpoint defaults to the last 2 hours.
func GetAllUsers(ctx context.Context, service *zscaler.Service, filters GetUsersFilters) ([]User, *http.Response, error) {
	var v struct {
		NextOffSet interface{} `json:"next_offset"`
		List       []User      `json:"users"`
	}

	relativeURL := usersEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, filters, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v.List, resp, nil
}
