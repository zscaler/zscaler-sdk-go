package policysetcontroller

import (
	"net/http"
)

const (
	platformEndpoint = "/platform"
)

type Platform struct {
	Linux   string `json:"linux"`
	Android string `json:"android"`
	Windows string `json:"windows"`
	IOS     string `json:"ios"`
	MacOS   string `json:"mac"`
}

func (service *Service) GetAllPlatforms() (*ClientTypes, *http.Response, error) {
	v := new(ClientTypes)
	relativeURL := mgmtConfig + platformEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
