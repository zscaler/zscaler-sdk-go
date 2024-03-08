package zpath

import (
	"net/http"
)

const (
	mgmtConfig    = "/mgmtconfig/v1/admin"
	zpathEndpoint = "/zpathCloud/getAltClouds"
)

func (service *Service) GetAltCloud() ([]string, *http.Response, error) {
	relativeURL := mgmtConfig + zpathEndpoint
	var cloudEndpoints []string

	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &cloudEndpoints)
	if err != nil {
		return nil, nil, err
	}
	service.Client.Config.Logger.Printf("[INFO] got alternate clouds: %#v", cloudEndpoints)
	return cloudEndpoints, resp, nil
}
