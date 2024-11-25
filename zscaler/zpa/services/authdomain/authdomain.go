package authdomain

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	authDomainEndpoint = "/authDomains"
)

type AuthDomain struct {
	AuthDomains []string `json:"authDomains"`
}

func GetAllAuthDomains(ctx context.Context, service *zscaler.Service) (*AuthDomain, *http.Response, error) {
	v := new(AuthDomain)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + authDomainEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
