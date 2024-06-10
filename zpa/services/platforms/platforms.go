package platforms

import (
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
)

const (
	mgmtConfig       = "/mgmtconfig/v1/admin/customers/"
	platformEndpoint = "/platform"
)

type Platforms struct {
	Linux   string `json:"linux"`
	Android string `json:"android"`
	Windows string `json:"windows"`
	IOS     string `json:"ios"`
	MacOS   string `json:"mac"`
}

func GetAllPlatforms(service *services.Service) (*Platforms, *http.Response, error) {
	v := new(Platforms)
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + platformEndpoint
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
