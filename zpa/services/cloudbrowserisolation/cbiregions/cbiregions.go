package cbiregions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

const (
	cbiConfig          = "/cbiconfig/cbi/api/customers/"
	cbiRegionsEndpoint = "/regions"
)

type CBIRegions struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// The current API does not seem to support search by Name
func GetByName(service *services.Service, cbiRegionName string) (*CBIRegions, *http.Response, error) {
	list, resp, err := GetAll(service)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, cbiRegionName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no region named '%s' was found", cbiRegionName)
}

func GetAll(service *services.Service) ([]CBIRegions, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiRegionsEndpoint
	var list []CBIRegions
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
