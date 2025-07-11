package praportal

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig        = "/zpa/mgmtconfig/v1/admin/customers/"
	praPortalEndpoint = "/praPortal"
)

type PRAPortal struct {
	// The unique identifier of the privileged portal.
	ID string `json:"id,omitempty"`

	// The name of the privileged portal.
	Name string `json:"name,omitempty"`

	// The description of the privileged portal.
	Description string `json:"description,omitempty"`

	// Whether or not the privileged portal is enabled.
	Enabled bool `json:"enabled"`

	// The canonical name (CNAME DNS records) associated with the privileged portal.
	CName string `json:"cName,omitempty"`

	// The domain of the privileged portal.
	Domain string `json:"domain,omitempty"`

	// The unique identifier of the certificate.
	CertificateID string `json:"certificateId,omitempty"`

	// The name of the certificate.
	CertificateName string `json:"certificateName,omitempty"`

	// The time the privileged portal is created.
	CreationTime string `json:"creationTime,omitempty"`

	// The unique identifier of the tenant who modified the privileged portal.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// The time the privileged portal is modified.
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// The notification message displayed in the banner of the privileged portallink, if enabled.
	UserNotification string `json:"userNotification"`

	// Indicates if the Notification Banner is enabled (true) or disabled (false).
	UserNotificationEnabled bool `json:"userNotificationEnabled"`

	ExtDomain string `json:"extDomain"`

	ExtDomainName string `json:"extDomainName"`

	ExtDomainTranslation string `json:"extDomainTranslation"`

	ExtLabel string `json:"extLabel"`

	UserPortalGid string `json:"userPortalGid,omitempty"`

	UserPortalName string `json:"userPortalName,omitempty"`

	GetcName string `json:"getcName,omitempty"`

	// The name of the Microtenant.
	MicroTenantName string `json:"microtenantName,omitempty"`

	// The name of the Microtenant.
	MicroTenantID string `json:"microtenantId,omitempty"`

	ObjectType string `json:"objectType,omitempty"`

	Action string `json:"action,omitempty"`

	CertManagedByZsRadio string `json:"certManagedByZsRadio,omitempty"`

	IsSRAPortal bool `json:"isSRAPortal,omitempty"`

	ManagedByZs bool `json:"managedByZs,omitempty"`

	ScopeName string `json:"scopeName,omitempty"`

	HideInfoTooltip bool `json:"hideInfoTooltip,omitempty"`

	RestrictedEntity bool `json:"restrictedEntity,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, portalID string) (*PRAPortal, *http.Response, error) {
	v := new(PRAPortal)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, portalName string) (*PRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + praPortalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAPortal](ctx, service.Client, relativeURL, common.Filter{Search: portalName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, sra := range list {
		if strings.EqualFold(sra.Name, portalName) {
			return &sra, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no pra portal '%s' was found", portalName)
}

func Create(ctx context.Context, service *zscaler.Service, sraPortal *PRAPortal) (*PRAPortal, *http.Response, error) {
	v := new(PRAPortal)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+praPortalEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, sraPortal, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, portalID string, sraPortal *PRAPortal) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, sraPortal, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, portalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + praPortalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAPortal](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
