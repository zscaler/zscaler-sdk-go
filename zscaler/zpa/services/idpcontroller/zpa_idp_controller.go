package idpcontroller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfigV1          = "/zpa/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2          = "/zpa/mgmtconfig/v2/admin/customers/"
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

func Get(ctx context.Context, service *zscaler.Service, IdpID string) (*IdpController, *http.Response, error) {
	v := new(IdpController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+idpControllerEndpoint, IdpID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, idpName string) (*IdpController, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.GetCustomerID() + idpControllerEndpoint)
	list, resp, err := common.GetAllPagesGeneric[IdpController](ctx, service.Client, relativeURL, "")
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

func GetAll(ctx context.Context, service *zscaler.Service) ([]IdpController, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.GetCustomerID() + idpControllerEndpoint)
	list, resp, err := common.GetAllPagesGeneric[IdpController](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
