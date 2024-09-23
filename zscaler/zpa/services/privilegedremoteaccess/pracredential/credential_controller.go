package pracredential

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	credentialEndpoint = "/credential"
)

type Credential struct {
	// The unique identifier of the privileged credential
	ID string `json:"id,omitempty"`

	//The name of the privileged credential.
	Name string `json:"name,omitempty"`

	// The description of the privileged credential.
	Description string `json:"description,omitempty"`

	// The time the privileged credential was last reset.
	LastCredentialResetTime string `json:"lastCredentialResetTime,omitempty"`

	// The protocol type that was designated for that particular privileged credential.
	// The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements.
	CredentialType string `json:"credentialType,omitempty"`

	// The password that is used to protect the SSH private key. This field is optional.
	Passphrase string `json:"passphrase,omitempty"`

	// The password associated with the username for the login you want to use for the privileged credential.
	Password string `json:"password,omitempty"`

	// The SSH private key associated with the username for the login you want to use for the privileged credential.
	PrivateKey string `json:"privateKey,omitempty"`

	// The domain name associated with the username.
	// You can also include the domain name as part of the username.
	// The domain name only needs to be specified with logging in to an RDP console that is connected to an Active Directory Domain.
	UserDomain string `json:"userDomain,omitempty"`

	// The username for the login you want to use for the privileged credential.
	UserName string `json:"userName,omitempty"`

	// The time the privileged credential is created.
	CreationTime string `json:"creationTime,omitempty"`

	// The unique identifier of the tenant who modified the privileged credential.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// The time the privileged credential is modified.
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// The unique identifier of the Microtenant for the ZPA tenant.
	// If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.
	// Pass microtenantId as null to retrieve data from all customers associated with the tenant.
	MicroTenantID string `json:"microtenantId,omitempty"`

	// The name of the Microtenant.
	MicroTenantName string `json:"microtenantName,omitempty"`

	// The unique identifier of the target Microtenant that the privileged credential is being moved to.
	TargetMicrotenantId string `json:"targetMicrotenantId,omitempty"`
}

func Get(service *zscaler.Service, credentialID string) (*Credential, *http.Response, error) {
	v := new(Credential)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *zscaler.Service, credentialName string) (*Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[Credential](service.Client, relativeURL, common.Filter{Search: credentialName, MicroTenantID: service.MicroTenantID()})
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

func Create(service *zscaler.Service, credential *Credential) (*Credential, *http.Response, error) {
	v := new(Credential)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, credential, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(service *zscaler.Service, credentialID string, credentialRequest *Credential) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, credentialRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(service *zscaler.Service, credentialID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func CredentialMove(service *zscaler.Service, credentialID string, targetMicrotenantId string) (*http.Response, error) {
	// Construct the URL using the credentialEndpoint const and append "/move"
	relativeURL := fmt.Sprintf("%s%s%s/%s/move", mgmtConfig, service.Client.GetCustomerID(), credentialEndpoint, credentialID)

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

func GetAll(service *zscaler.Service) ([]Credential, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[Credential](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
