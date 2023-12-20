package administration

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zdx/services/common"
)

const (
	departmentsEndpoint = "v1/administration/departments"
	locationsEndpoint   = "v1/administration/locations"
)

type Departments struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type Locations struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

type GetDepartmentsFilters struct {
	common.GetFromToFilters
	// The Zscaler location (ID). You can add multiple location IDs.
	Loc []int `json:"loc,omitempty" url:"loc,omitempty"`
	//The search string used to support search by name or department ID.
	Search string `json:"search,omitempty" url:"search,omitempty"`
}

type GetLocationsFilters struct {
	common.GetFromToFilters
	// The Zscaler location (ID). You can add multiple location IDs.
	Loc []int `json:"loc,omitempty" url:"loc,omitempty"`
	//The search string used to support search by name or department ID.
	Search string `json:"q,omitempty" url:"q,omitempty"`
}

// Gets configured departments.
func (service *Service) GetDepartments(appID string, filters GetDepartmentsFilters) (*Departments, *http.Response, error) {
	v := new(Departments)
	path := departmentsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Gets configured departments.
func (service *Service) GetLocations(appID string, filters GetLocationsFilters) (*Locations, *http.Response, error) {
	v := new(Locations)
	path := locationsEndpoint
	resp, err := service.Client.NewRequestDo("GET", path, filters, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetAllDepartments() ([]Departments, error) {
	var departments []Departments
	resp, err := service.Client.NewRequestDo("GET", departmentsEndpoint, nil, nil, &departments)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: status code %d", resp.StatusCode)
	}
	return departments, nil
}

func (service *Service) GetAllLocations() ([]Locations, error) {
	var locations []Locations
	resp, err := service.Client.NewRequestDo("GET", locationsEndpoint, nil, nil, &locations)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request error: status code %d", resp.StatusCode)
	}
	return locations, nil
}
