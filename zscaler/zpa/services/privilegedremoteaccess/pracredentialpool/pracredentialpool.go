package pracredentialpool

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig         = "/zpa/waap-pra-config/v1/admin/customers/"
	credentialEndpoint = "/credential-pool"
)

type CredentialPool struct {

	// The unique identifier of the privileged credential
	ID string `json:"id,omitempty"`

	//The name of the privileged credential.
	Name string `json:"name,omitempty"`

	// The protocol type that was designated for that particular privileged credential.
	// The protocol type options are SSH, RDP, and VNC. Each protocol type has its own credential requirements.
	CredentialType string `json:"credentialType,omitempty"`

	PRACredentials []common.CommonIDName `json:"credentials"`

	CredentialMappingCount string `json:"credentialMappingCount,omitempty"`

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
}

func Get(ctx context.Context, service *zscaler.Service, credentialID string) (*CredentialPool, *http.Response, error) {
	v := new(CredentialPool)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, credentialName string) (*CredentialPool, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[CredentialPool](ctx, service.Client, relativeURL, common.Filter{Search: credentialName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, credential *CredentialPool) (*CredentialPool, *http.Response, error) {
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, credential, nil)
	if err != nil {
		return nil, resp, err
	}

	// Do a follow-up GetAll to find the created resource by name
	all, _, err := GetAll(ctx, service)
	if err != nil {
		return nil, resp, fmt.Errorf("credential pool created, but failed to fetch for ID lookup: %w", err)
	}

	for _, c := range all {
		if strings.EqualFold(c.Name, credential.Name) {
			return &c, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("credential pool created, but could not locate ID for name: %s", credential.Name)
}

func Update(ctx context.Context, service *zscaler.Service, credentialID string, credentialRequest *CredentialPool) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, credentialRequest, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, credentialID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+credentialEndpoint, credentialID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CredentialPool, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + credentialEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[CredentialPool](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
