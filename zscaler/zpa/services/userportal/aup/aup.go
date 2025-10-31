package aup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v1/admin/customers/"
	userPortalAUPEndpoint = "/userportal/aup"
)

type UserPortalAup struct {
	Aup             string `json:"aup,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	Description     string `json:"description,omitempty"`
	Email           string `json:"email,omitempty"`
	Enabled         bool   `json:"enabled,omitempty"`
	ID              string `json:"id,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	Name            string `json:"name,omitempty"`
	PhoneNum        string `json:"phoneNum,omitempty"`
	MicrotenantID   string `json:"microtenantId,omitempty"`
	MicrotenantName string `json:"microtenantName,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, aupID string) (*UserPortalAup, *http.Response, error) {
	v := new(UserPortalAup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+userPortalAUPEndpoint, aupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, userPortalName string) (*UserPortalAup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalAUPEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalAup](ctx, service.Client, relativeURL, common.Filter{Search: userPortalName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, userPortalName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", userPortalName)
}

func Create(ctx context.Context, service *zscaler.Service, userPortalAup *UserPortalAup) (*UserPortalAup, *http.Response, error) {
	v := new(UserPortalAup)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+userPortalAUPEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, userPortalAup, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, aupID string, userPortalAup *UserPortalAup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalAUPEndpoint, aupID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, userPortalAup, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, aupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalAUPEndpoint, aupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]UserPortalAup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalAUPEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalAup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
