package praportal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig        = "/mgmtconfig/v1/admin/customers/"
	sraPortalEndpoint = "/praPortal"
)

type SRAPortal struct {
	ID                      string `json:"id,omitempty"`
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	Enabled                 bool   `json:"enabled"`
	CName                   string `json:"cName,omitempty"`
	Domain                  string `json:"domain,omitempty"`
	CertificateID           string `json:"certificateId,omitempty"`
	CertificateName         string `json:"certificateName,omitempty"`
	CreationTime            string `json:"creationTime,omitempty"`
	ModifiedBy              string `json:"modifiedBy,omitempty"`
	ModifiedTime            string `json:"modifiedTime,omitempty"`
	UserNotification        string `json:"userNotification"`
	UserNotificationEnabled bool   `json:"userNotificationEnabled"`
}

func (service *Service) Get(portalID string) (*SRAPortal, *http.Response, error) {
	v := new(SRAPortal)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+sraPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(portalName string) (*SRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + sraPortalEndpoint
	list, resp, err := common.GetAllPagesGeneric[SRAPortal](service.Client, relativeURL, portalName)
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

func (service *Service) Create(sraPortal *SRAPortal) (*SRAPortal, *http.Response, error) {
	v := new(SRAPortal)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+sraPortalEndpoint, nil, sraPortal, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(portalID string, sraPortal *SRAPortal) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+sraPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, sraPortal, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(portalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+sraPortalEndpoint, portalID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]SRAPortal, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + sraPortalEndpoint
	list, resp, err := common.GetAllPagesGeneric[SRAPortal](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
