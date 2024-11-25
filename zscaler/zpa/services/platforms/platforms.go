package platforms

import (
	"context"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig       = "/zpa/mgmtconfig/v1/admin/customers/"
	platformEndpoint = "/platform"
)

type Platforms struct {
	Linux   string `json:"linux"`
	Android string `json:"android"`
	Windows string `json:"windows"`
	IOS     string `json:"ios"`
	MacOS   string `json:"mac"`
}

func GetAllPlatforms(ctx context.Context, service *zscaler.Service) (*Platforms, *http.Response, error) {
	v := new(Platforms)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + platformEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
