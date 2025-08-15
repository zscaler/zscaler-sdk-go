package portal_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	userPortalEndpoint = "/userPortal"
)

type UserPortalController struct {
	CertificateId           string `json:"certificateId,omitempty"`
	CertificateName         string `json:"certificateName,omitempty"`
	CreationTime            string `json:"creationTime,omitempty"`
	Description             string `json:"description,omitempty"`
	Domain                  string `json:"domain,omitempty"`
	Enabled                 bool   `json:"enabled,omitempty"`
	ExtDomain               string `json:"extDomain,omitempty"`
	ExtDomainName           string `json:"extDomainName,omitempty"`
	ExtDomainTranslation    string `json:"extDomainTranslation,omitempty"`
	ExtLabel                string `json:"extLabel,omitempty"`
	GetcName                string `json:"getcName,omitempty"`
	ID                      string `json:"id,omitempty"`
	ImageData               string `json:"imageData,omitempty"`
	ModifiedBy              string `json:"modifiedBy,omitempty"`
	ModifiedTime            string `json:"modifiedTime,omitempty"`
	Name                    string `json:"name,omitempty"`
	MicrotenantId           string `json:"microtenantId,omitempty"`
	MicrotenantName         string `json:"microtenantName,omitempty"`
	UserNotification        string `json:"userNotification,omitempty"`
	UserNotificationEnabled bool   `json:"userNotificationEnabled,omitempty"`
	ManagedByZS             bool   `json:"managedByZs,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, portalID string) (*UserPortalController, *http.Response, error) {
	v := new(UserPortalController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+userPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, portalName string) (*UserPortalController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalController](ctx, service.Client, relativeURL, common.Filter{Search: portalName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, portalName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no user portal named '%s' was found", portalName)
}

func Create(ctx context.Context, service *zscaler.Service, controllerGroup UserPortalController) (*UserPortalController, *http.Response, error) {
	v := new(UserPortalController)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+userPortalEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, portalID string, controllerGroup *UserPortalController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, controllerGroup, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, portalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+userPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]UserPortalController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + userPortalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[UserPortalController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
