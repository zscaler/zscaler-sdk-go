package custom_config_controller

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	mgmtConfig = "/zpa/mgmtconfig/v1/admin/customers/"
)

type ZIACloudConfig struct {
	ZIACloudDomain        string `json:"ziaCloudDomain,omitempty"`
	ZIACloudServiceApiKey string `json:"ziaCloudServiceApiKey,omitempty"`
	ZIAPassword           string `json:"ziaPassword,omitempty"`
	ZIASandboxApiToken    string `json:"ziaSandboxApiToken,omitempty"`
	ZIAUsername           string `json:"ziaUsername,omitempty"`
}

type SessionTerminationOnReath struct {
	AllowDisableSessionTerminationOnReauth bool `json:"allowDisableSessionTerminationOnReauth,omitempty"`
	SessionTerminationOnReauth             bool `json:"sessionTerminationOnReauth,omitempty"`
}

func CheckZiaCloudConfig(ctx context.Context, service *zscaler.Service) (bool, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/config/isZiaCloudConfigAvailable"

	var result bool
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &result)
	if err != nil {
		return false, nil, err
	}

	return result, resp, nil
}

func GetZIACloudConfig(ctx context.Context, service *zscaler.Service) (*ZIACloudConfig, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/config/ziaCloudConfig"

	var result ZIACloudConfig
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &result)
	if err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

func AddZIACloudConfig(ctx context.Context, service *zscaler.Service, cloudConfig *ZIACloudConfig) (*ZIACloudConfig, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/config/ziaCloudConfig"

	var result ZIACloudConfig
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, nil, cloudConfig, &result)
	if err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

func GetSessionTerminationOnReath(ctx context.Context, service *zscaler.Service) (*SessionTerminationOnReath, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/config/sessionTerminationOnReauth"

	var result SessionTerminationOnReath
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &result)
	if err != nil {
		return nil, nil, err
	}

	return &result, resp, nil
}

func UpdateSessionTerminationOnReath(ctx context.Context, service *zscaler.Service, sessionTermination *SessionTerminationOnReath) (*http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/config/sessionTerminationOnReauth"

	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, nil, sessionTermination, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
