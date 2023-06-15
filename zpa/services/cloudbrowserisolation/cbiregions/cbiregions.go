package cbiregions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig                 = "/cbiconfig/cbi/api/customers/"
	cbiRegionsEndpoint string = "/regions"
)

type CBIRegions struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (service *Service) Get(RegionID string) (*CBIRegions, *http.Response, error) {
	v := new(CBIRegions)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiRegionsEndpoint, RegionID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(cbiRegionName string) (*CBIRegions, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiRegionsEndpoint
	list, resp, err := common.GetAllPagesGeneric[CBIRegions](service.Client, relativeURL, "")
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

func (service *Service) GetAll() ([]CBIRegions, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiRegionsEndpoint
	list, resp, err := common.GetAllPagesGeneric[CBIRegions](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
