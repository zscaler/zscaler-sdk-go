package admin_sso_controller

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig       = "/zpa/mgmtconfig/v1/admin/customers/"
	ssoLoginEndpoint = "/v2/ssoLoginOptions"
)

type AdminSSOLoginOptions struct {
	SSOLoginOnly bool `json:"ssologinonly,omitempty"`
}

func GetSSOLoginController(ctx context.Context, service *zscaler.Service) (*AdminSSOLoginOptions, *http.Response, error) {
	v := new(AdminSSOLoginOptions)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + ssoLoginEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func UpdateSSOLoginController(ctx context.Context, service *zscaler.Service, settings *AdminSSOLoginOptions) (*AdminSSOLoginOptions, *http.Response, error) {
	v := &AdminSSOLoginOptions{
		SSOLoginOnly: settings.SSOLoginOnly, // carry through original value in case of 204
	}
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+ssoLoginEndpoint, settings, v, nil)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
