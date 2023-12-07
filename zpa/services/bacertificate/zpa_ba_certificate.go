package bacertificate

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV1                = "/mgmtconfig/v1/admin/customers/"
	baCertificateEndpoint       = "/certificate"
	mgmtConfigV2                = "/mgmtconfig/v2/admin/customers/"
	baCertificateIssuedEndpoint = "/clientlessCertificate/issued"
)

type BaCertificate struct {
	ID                  string   `json:"id,omitempty"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	CName               string   `json:"cName,omitempty"`
	CertChain           string   `json:"certChain,omitempty"`
	CertBlob            string   `json:"certBlob,omitempty"`
	Certificate         string   `json:"certificate,omitempty"`
	CreationTime        string   `json:"creationTime,omitempty"`
	IssuedBy            string   `json:"issuedBy,omitempty"`
	IssuedTo            string   `json:"issuedTo,omitempty"`
	ModifiedBy          string   `json:"modifiedBy,omitempty"`
	ModifiedTime        string   `json:"modifiedTime,omitempty"`
	San                 []string `json:"san,omitempty"`
	SerialNo            string   `json:"serialNo,omitempty"`
	Status              string   `json:"status,omitempty"`
	ValidFromInEpochSec string   `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec   string   `json:"validToInEpochSec,omitempty"`
	MicrotenantID       string   `json:"microtenantId,omitempty"`
	MicrotenantName     string   `json:"microtenantName,omitempty"`
}

func (service *Service) Get(baCertificateID string) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetIssuedByName(CertName string) (*BaCertificate, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + baCertificateIssuedEndpoint)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BaCertificate](service.Client, relativeURL, common.Filter{Search: CertName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, baCertificate := range list {
		if baCertificate.Name == CertName {
			return &baCertificate, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no issued certificate named '%s' was found", CertName)
}

func (service *Service) GetAll() ([]BaCertificate, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + baCertificateIssuedEndpoint)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[BaCertificate](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func (service *Service) Create(baCertificate BaCertificate) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, common.Filter{MicroTenantID: service.microTenantID}, baCertificate, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Delete(baCertificateID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
