package provisioningkey

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig = "/zpa/mgmtconfig/v1/admin/customers/"
)

// TODO: because there isn't an endpoint to get all provisionning keys, we need to have all association type here
var ProvisioningKeyAssociationTypes []string = []string{
	"CONNECTOR_GRP",
	"SERVICE_EDGE_GRP",
}

type ProvisioningKey struct {
	AppConnectorGroupID   string   `json:"appConnectorGroupId,omitempty"`
	AppConnectorGroupName string   `json:"appConnectorGroupName,omitempty"`
	CreationTime          string   `json:"creationTime,omitempty"`
	Enabled               bool     `json:"enabled,omitempty"`
	ExpirationInEpochSec  string   `json:"expirationInEpochSec,omitempty"`
	ID                    string   `json:"id,omitempty"`
	IPACL                 []string `json:"ipAcl,omitempty"`
	MaxUsage              string   `json:"maxUsage,omitempty"`
	ModifiedBy            string   `json:"modifiedBy,omitempty"`
	ModifiedTime          string   `json:"modifiedTime,omitempty"`
	Name                  string   `json:"name,omitempty"`
	ProvisioningKey       string   `json:"provisioningKey,omitempty"`
	EnrollmentCertID      string   `json:"enrollmentCertId,omitempty"`
	EnrollmentCertName    string   `json:"enrollmentCertName,omitempty"`
	UIConfig              string   `json:"uiConfig,omitempty"`
	UsageCount            string   `json:"usageCount,omitempty"`
	ZcomponentID          string   `json:"zcomponentId,omitempty"`
	ZcomponentName        string   `json:"zcomponentName,omitempty"`
	AssociationType       string   `json:"associationType"`
	MicroTenantID         string   `json:"microtenantId,omitempty"`
	MicroTenantName       string   `json:"microtenantName,omitempty"`
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func Get(ctx context.Context, service *zscaler.Service, associationType, provisioningKeyID string) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	url := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", url, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	v.AssociationType = associationType
	return v, resp, nil
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func GetByName(ctx context.Context, service *zscaler.Service, associationType, name string) (*ProvisioningKey, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey", associationType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ProvisioningKey](ctx, service.Client, relativeURL, common.Filter{Search: name, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, provisioningKey := range list {
		if strings.EqualFold(provisioningKey.Name, name) {
			provisioningKey.AssociationType = associationType
			return &provisioningKey, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no Provisioning Key named '%s' was found", name)
}

// POST --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func Create(ctx context.Context, service *zscaler.Service, associationType string, provisioningKey *ProvisioningKey) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	path := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey", associationType)
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, common.Filter{MicroTenantID: service.MicroTenantID()}, provisioningKey, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func Update(ctx context.Context, service *zscaler.Service, associationType, provisioningKeyID string, provisioningKey *ProvisioningKey) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, provisioningKey, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func Delete(ctx context.Context, service *zscaler.Service, associationType, provisioningKeyID string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetByNameAllAssociations(ctx context.Context, service *zscaler.Service, name string) (p *ProvisioningKey, assoc_type string, resp *http.Response, err error) {
	for _, associationType := range ProvisioningKeyAssociationTypes {
		p, resp, err = GetByName(ctx, service, associationType, name)
		if err == nil {
			assoc_type = associationType
			break
		}
	}
	if p != nil {
		p.AssociationType = assoc_type
	}
	return p, assoc_type, resp, err
}

func GetByIDAllAssociations(ctx context.Context, service *zscaler.Service, id string) (p *ProvisioningKey, assoc_type string, resp *http.Response, err error) {
	for _, associationType := range ProvisioningKeyAssociationTypes {
		p, resp, err = Get(ctx, service, associationType, id)
		if err == nil {
			assoc_type = associationType
			break
		}
	}
	if p != nil {
		p.AssociationType = assoc_type
	}
	return p, assoc_type, resp, err
}

func GetAllByAssociationType(ctx context.Context, service *zscaler.Service, associationType string) ([]ProvisioningKey, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.GetCustomerID()+"/associationType/%s/provisioningKey", associationType)
	list, _, err := common.GetAllPagesGenericWithCustomFilters[ProvisioningKey](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].AssociationType = associationType
	}
	return list, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) (list []ProvisioningKey, err error) {
	for _, associationType := range ProvisioningKeyAssociationTypes {
		items, _ := GetAllByAssociationType(ctx, service, associationType)
		if len(items) > 0 {
			list = append(list, items...)
		}
	}
	return list, nil
}
