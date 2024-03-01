package credentialcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	credentialEndpoint = "/credential"
)

type Credential struct {
	ID                      string `json:"id,omitempty"`
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	LastCredentialResetTime string `json:"lastCredentialResetTime,omitempty"`
	CredentialType          string `json:"credentialType,omitempty"`
	Passphrase              string `json:"passphrase,omitempty"`
	Password                string `json:"password,omitempty"`
	PrivateKey              string `json:"privateKey,omitempty"`
	UserDomain              string `json:"userDomain,omitempty"`
	UserName                string `json:"userName,omitempty"`
	CreationTime            string `json:"creationTime,omitempty"`
	ModifiedBy              string `json:"modifiedBy,omitempty"`
	ModifiedTime            string `json:"modifiedTime,omitempty"`
	MicroTenantID           string `json:"microtenantId,omitempty"`
	MicroTenantName         string `json:"microtenantName,omitempty"`
	TargetMicrotenantId     string `json:"targetMicrotenantId,omitempty"`
}

func (service *Service) Get(credentialID string) (*Credential, *http.Response, error) {
	v := new(Credential)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(credentialName string) (*Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[Credential](service.Client, relativeURL, common.Filter{Search: credentialName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, cred := range list {
		if strings.EqualFold(cred.Name, credentialName) {
			return &cred, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no credential controller named '%s' was found", credentialName)
}

func (service *Service) Create(credential *Credential) (*Credential, *http.Response, error) {
	v := new(Credential)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, common.Filter{MicroTenantID: service.microTenantID}, credential, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(credentialID string, credentialRequest *Credential) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, credentialRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(credentialID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

/*
	func (service *Service) CredentialMove(credentialID string, targetMicrotenantId string) (*http.Response, error) {
		// Construct the URL using the credentialEndpoint const and append "/move"
		relativeURL := fmt.Sprintf("%s%s%s/%s/move", mgmtConfig, service.Client.Config.CustomerID, credentialEndpoint, credentialID)

		// Append the targetMicrotenantId as a query parameter
		if targetMicrotenantId != "" {
			relativeURL += "?targetMicrotenantId=" + targetMicrotenantId
		}

		// Make the POST request with an empty body since the API expects an empty body for this operation
		resp, err := service.Client.NewRequestDo("POST", relativeURL, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
*/
func (service *Service) GetAll() ([]Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[Credential](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
