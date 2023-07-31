package idpcontroller

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v2/admin/customers/"
	mgmtConfigV1          = "/mgmtconfig/v1/admin/customers/"
	idpControllerEndpoint = "/idp"
)

type IdpController struct {
	AdminSpSigningCertID        string         `json:"adminSpSigningCertId,omitempty"`
	AutoProvision               string         `json:"autoProvision,omitempty"`
	CreationTime                string         `json:"creationTime,omitempty"`
	Description                 string         `json:"description,omitempty"`
	DisableSamlBasedPolicy      bool           `json:"disableSamlBasedPolicy"`
	Domainlist                  []string       `json:"domainList,omitempty"`
	EnableScimBasedPolicy       bool           `json:"enableScimBasedPolicy"`
	EnableArbitraryAuthDomains  string         `json:"enableArbitraryAuthDomains"`
	Enabled                     bool           `json:"enabled"`
	ForceAuth                   bool           `json:"forceAuth"`
	ID                          string         `json:"id,omitempty"`
	IdpEntityID                 string         `json:"idpEntityId,omitempty"`
	LoginHint                   bool           `json:"loginHint,omitempty"`
	LoginNameAttribute          string         `json:"loginNameAttribute,omitempty"`
	LoginURL                    string         `json:"loginUrl,omitempty"`
	ModifiedBy                  string         `json:"modifiedBy,omitempty"`
	ModifiedTime                string         `json:"modifiedTime,omitempty"`
	Name                        string         `json:"name,omitempty"`
	ReauthOnUserUpdate          bool           `json:"reauthOnUserUpdate"`
	RedirectBinding             bool           `json:"redirectBinding"`
	ScimEnabled                 bool           `json:"scimEnabled"`
	ScimServiceProviderEndpoint string         `json:"scimServiceProviderEndpoint,omitempty"`
	ScimSharedSecretExists      bool           `json:"scimSharedSecretExists,omitempty"`
	SignSamlRequest             string         `json:"signSamlRequest,,omitempty"`
	SsoType                     []string       `json:"ssoType,omitempty"`
	UseCustomSpMetadata         bool           `json:"useCustomSPMetadata"`
	UserSpSigningCertID         string         `json:"userSpSigningCertId,omitempty"`
	AdminMetadata               *AdminMetadata `json:"adminMetadata,omitempty"`
	UserMetadata                *UserMetadata  `json:"userMetadata,omitempty"`
}

type AdminMetadata struct {
	CertificateURL string `json:"certificateUrl"`
	SpBaseURL      string `json:"spBaseUrl"`
	SpEntityID     string `json:"spEntityId"`
	SpMetadataURL  string `json:"spMetadataUrl"`
	SpPostURL      string `json:"spPostUrl"`
}

type UserMetadata struct {
	CertificateURL string `json:"certificateUrl,omitempty"`
	SpBaseURL      string `json:"spBaseUrl"`
	SpEntityID     string `json:"spEntityId,omitempty"`
	SpMetadataURL  string `json:"spMetadataUrl,omitempty"`
	SpPostURL      string `json:"spPostUrl,omitempty"`
}

func (service *Service) Get(IdpID string) (*IdpController, *http.Response, error) {
	v := new(IdpController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+idpControllerEndpoint, IdpID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(idpName string) (*IdpController, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + idpControllerEndpoint)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IdpController](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, idpController := range list {
		if idpController.Name == idpName {
			return &idpController, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no Idp-Controller named '%s' was found", idpName)
}

func (service *Service) GetAll() ([]IdpController, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + idpControllerEndpoint)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IdpController](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
