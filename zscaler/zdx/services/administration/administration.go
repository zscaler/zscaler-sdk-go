package administration

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	departmentsEndpoint = "/zdx/v1/administration/departments"
	locationsEndpoint   = "/zdx/v1/administration/locations"
)

type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Location struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type GetDepartmentsFilters struct {
	From   int    `json:"from,omitempty" url:"from,omitempty"`
	To     int    `json:"to,omitempty" url:"to,omitempty"`
	Search string `json:"search,omitempty" url:"search,omitempty"`
}

type GetLocationsFilters struct {
	From   int    `json:"from,omitempty" url:"from,omitempty"`
	To     int    `json:"to,omitempty" url:"to,omitempty"`
	Search string `json:"search,omitempty" url:"search,omitempty"`
	Q      string `json:"q,omitempty" url:"q,omitempty"`
}

// Gets the list of configured departments.
func GetDepartments(ctx context.Context, service *zscaler.Service, filters GetDepartmentsFilters) ([]Department, *http.Response, error) {
	var departments []Department
	path := departmentsEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &departments)
	if err != nil {
		return nil, nil, err
	}
	return departments, resp, nil
}

// Gets the list of configured locations.
func GetLocations(ctx context.Context, service *zscaler.Service, filters GetLocationsFilters) ([]Location, *http.Response, error) {
	var locations []Location
	path := locationsEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, filters, nil, &locations)
	if err != nil {
		return nil, nil, err
	}
	return locations, resp, nil
}
