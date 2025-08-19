package config_override

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig              = "/zpa/mgmtconfig/v1/admin/customers/"
	configOverridesEndpoint = "/configOverrides"
)

type ConfigOverrides struct {
	BrokerName     string `json:"brokerName,omitempty"`
	ConfigKey      string `json:"configKey,omitempty"`
	ConfigValue    string `json:"configValue,omitempty"`
	ConfigValueInt string `json:"configValueInt,omitempty"`
	CustomerId     string `json:"customerId,omitempty"`
	CustomerName   string `json:"customerName,omitempty"`
	Description    string `json:"description,omitempty"`
	TargetGid      string `json:"targetGid,omitempty"`
	TargetName     string `json:"targetName,omitempty"`
	TargetType     string `json:"targetType,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, configID string) (*ConfigOverrides, *http.Response, error) {
	v := new(ConfigOverrides)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+configOverridesEndpoint, configID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Create(ctx context.Context, service *zscaler.Service, configOverride *ConfigOverrides) (*ConfigOverrides, *http.Response, error) {
	v := new(ConfigOverrides)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+configOverridesEndpoint, common.Filter{}, configOverride, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, configID string, configOverride *ConfigOverrides) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+configOverridesEndpoint, configID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{}, configOverride, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ConfigOverrides, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + configOverridesEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ConfigOverrides](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
