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
}

func (service *Service) Get(credentialID string) (*Credential, *http.Response, error) {
	v := new(Credential)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(credentialName string) (*Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + credentialEndpoint
	list, resp, err := common.GetAllPagesGeneric[Credential](service.Client, relativeURL, credentialName)
	if err != nil {
		return nil, nil, err
	}
	for _, cred := range list {
		if strings.EqualFold(cred.Name, credentialName) {
			return &cred, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", credentialName)
}

func (service *Service) Create(credential *Credential) (*Credential, *http.Response, error) {
	v := new(Credential)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, nil, credential, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(credentialID string, credentialRequest *Credential) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, credentialRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(credentialID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + credentialEndpoint
	list, resp, err := common.GetAllPagesGeneric[Credential](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
