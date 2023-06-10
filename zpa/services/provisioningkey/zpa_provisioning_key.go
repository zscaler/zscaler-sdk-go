package provisioningkey

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
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
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) Get(associationType, provisioningKeyID string) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("GET", url, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	v.AssociationType = associationType
	return v, resp, nil
}

// GET --> mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey
func (service *Service) GetByName(associationType, name string) (*ProvisioningKey, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey", associationType)
	list, resp, err := common.GetAllPagesGeneric[ProvisioningKey](service.Client, relativeURL, name)
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
func (service *Service) Create(associationType string, provisioningKey *ProvisioningKey) (*ProvisioningKey, *http.Response, error) {
	v := new(ProvisioningKey)
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey", associationType)
	resp, err := service.Client.NewRequestDo("POST", path, nil, provisioningKey, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func (service *Service) Update(associationType, provisioningKeyID string, provisioningKey *ProvisioningKey) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, provisioningKey, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> /mgmtconfig/v1/admin/customers/{customerId}/associationType/{associationType}/provisioningKey/{provisioningKeyId}
func (service *Service) Delete(associationType, provisioningKeyID string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey/%s", associationType, provisioningKeyID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetByNameAllAssociations(name string) (p *ProvisioningKey, assoc_type string, resp *http.Response, err error) {
	for _, assassociation_type := range ProvisioningKeyAssociationTypes {
		p, resp, err = service.GetByName(assassociation_type, name)
		if err == nil {
			assoc_type = assassociation_type
			break
		}
	}
	p.AssociationType = assoc_type
	return p, assoc_type, resp, err
}

func (service *Service) GetByIDAllAssociations(id string) (p *ProvisioningKey, assoc_type string, resp *http.Response, err error) {
	for _, assassociation_type := range ProvisioningKeyAssociationTypes {
		p, resp, err = service.Get(assassociation_type, id)
		if err == nil {
			assoc_type = assassociation_type
			break
		}
	}
	p.AssociationType = assoc_type
	return p, assoc_type, resp, err
}

func (service *Service) GetAllByAssociationType(associationType string) ([]ProvisioningKey, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/associationType/%s/provisioningKey", associationType)
	list, _, err := common.GetAllPagesGeneric[ProvisioningKey](service.Client, relativeURL, "")
	if err != nil {
		return nil, err
	}
	for i := range list {
		list[i].AssociationType = associationType
	}
	return list, nil
}

func (service *Service) GetAll() (list []ProvisioningKey, err error) {
	for _, assassociation_type := range ProvisioningKeyAssociationTypes {
		items, _ := service.GetAllByAssociationType(assassociation_type)
		if len(items) > 0 {
			list = append(list, items...)
		}
	}
	return list, nil
}
