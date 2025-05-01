package client_settings

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v1/admin/customers/"
	clientSettingEndpoint = "/clientSetting"
)

var supportedClientSettingTypes = map[string]bool{
	"ZAPP_CLIENT":      true,
	"ISOLATION_CLIENT": true,
	"APP_PROTECTION":   true,
}

type ClientSettings struct {
	ID                           string `json:"id,omitempty"`
	CreationTime                 string `json:"creationTime,omitempty"`
	ModifiedBy                   string `json:"modifiedBy,omitempty"`
	ClientCertificateType        string `json:"clientCertificateType,omitempty"`
	SingningCertExpiryInEpochSec string `json:"singningCertExpiryInEpochSec,omitempty"`
	Name                         string `json:"name,omitempty"`
	EnrollmentCertId             string `json:"enrollmentCertId,omitempty"`
	EnrollmentCertName           string `json:"enrollmentCertName,omitempty"`
}

func GetClientSettings(ctx context.Context, service *zscaler.Service, clientType *string) ([]ClientSettings, *http.Response, error) {
	var settings []ClientSettings
	baseURL := mgmtConfig + service.Client.GetCustomerID() + clientSettingEndpoint

	if clientType != nil {
		t := strings.ToUpper(strings.TrimSpace(*clientType))
		if !supportedClientSettingTypes[t] {
			return nil, nil, fmt.Errorf("invalid client setting type: %s", t)
		}
		baseURL += "?type=" + url.QueryEscape(t)
	}

	resp, err := service.Client.NewRequestDo(ctx, "GET", baseURL, nil, &settings, nil)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

func GetAllClientSettings(ctx context.Context, service *zscaler.Service) (*ClientSettings, *http.Response, error) {
	v := new(ClientSettings)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + clientSettingEndpoint + "/all"
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Create(ctx context.Context, service *zscaler.Service, settings *ClientSettings) (*ClientSettings, *http.Response, error) {
	v := new(ClientSettings)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+clientSettingEndpoint, settings, v, nil)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service) (*http.Response, error) {
	path := mgmtConfig + service.Client.GetCustomerID() + clientSettingEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
