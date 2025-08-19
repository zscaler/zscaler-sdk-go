package api_keys

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig      = "/zpa/mgmtconfig/v1/admin/customers/"
	apiKeysEndpoint = "/apiKeys"
)

type APIKeys struct {
	ClientID             string `json:"clientId,omitempty"`
	ClientSecret         string `json:"clientSecret,omitempty"`
	CreationTime         string `json:"creationTime,omitempty"`
	Enabled              bool   `json:"enabled,omitempty"`
	IamClientId          string `json:"iamClientId,omitempty"`
	ID                   string `json:"id,omitempty"`
	IsLocked             bool   `json:"isLocked,omitempty"`
	ModifiedBy           string `json:"modifiedBy,omitempty"`
	ModifiedTime         string `json:"modifiedTime,omitempty"`
	Name                 string `json:"name,omitempty"`
	PinSessionEnabled    bool   `json:"pinSessionEnabled,omitempty"`
	ReadOnly             bool   `json:"readOnly,omitempty"`
	RestrictionType      string `json:"restrictionType,omitempty"`
	RoleID               string `json:"roleId,omitempty"`
	MicrotenantId        string `json:"microtenantId,omitempty"`
	MicrotenantName      string `json:"microtenantName,omitempty"`
	SyncVersion          string `json:"syncVersion,omitempty"`
	TokenExpiryTimeInSec string `json:"tokenExpiryTimeInSec,omitempty"`
	ZscalerManaged       bool   `json:"zscalerManaged,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, keyID string) (*APIKeys, *http.Response, error) {
	v := new(APIKeys)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+apiKeysEndpoint, keyID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, keyName string) (*APIKeys, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + apiKeysEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[APIKeys](ctx, service.Client, relativeURL, common.Filter{Search: keyName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, keyName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no api key named '%s' was found", keyName)
}

func Create(ctx context.Context, service *zscaler.Service, apiKey APIKeys) (*APIKeys, *http.Response, error) {
	v := new(APIKeys)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+apiKeysEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, apiKey, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, keyID string, apiKey *APIKeys) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+apiKeysEndpoint, keyID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, apiKey, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, keyID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+apiKeysEndpoint, keyID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]APIKeys, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + apiKeysEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[APIKeys](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
