package praportal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig        = "/mgmtconfig/v1/admin/customers/"
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

	MicroTenantID string `json:"microtenantId,omitempty"`

	// The name of the Microtenant.
	MicroTenantName string `json:"microtenantName,omitempty"`
}

func (service *Service) Get(portalID string) (*PRAPortal, *http.Response, error) {
	v := new(PRAPortal)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(portalName string) (*PRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praPortalEndpoint
	list, resp, err := common.GetAllPagesGeneric[PRAPortal](service.Client, relativeURL, portalName)
	if err != nil {
		return nil, nil, err
	}
	for _, sra := range list {
		if strings.EqualFold(sra.Name, portalName) {
			return &sra, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no sra portal '%s' was found", portalName)
}

func (service *Service) Create(sraPortal *PRAPortal) (*PRAPortal, *http.Response, error) {
	v := new(PRAPortal)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+praPortalEndpoint, nil, sraPortal, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(portalID string, sraPortal *PRAPortal) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, sraPortal, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(portalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]PRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praPortalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAPortal](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
