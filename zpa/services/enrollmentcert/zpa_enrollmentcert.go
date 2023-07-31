package enrollmentcert

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v2/admin/customers/"
	enrollmentCertEndpoint = "/enrollmentCert"
)

type EnrollmentCert struct {
	AllowSigning            bool   `json:"allowSigning,omitempty"`
	Cname                   string `json:"cName,omitempty"`
	Certificate             string `json:"certificate,omitempty"`
	ClientCertType          string `json:"clientCertType,omitempty"`
	CreationTime            string `json:"creationTime,omitempty"`
	CSR                     string `json:"csr,omitempty"`
	Description             string `json:"description,omitempty"`
	ID                      string `json:"id,omitempty"`
	IssuedBy                string `json:"issuedBy,omitempty"`
	IssuedTo                string `json:"issuedTo,omitempty"`
	ModifiedBy              string `json:"modifiedBy,omitempty"`
	ModifiedTime            string `json:"modifiedTime,omitempty"`
	Name                    string `json:"name,omitempty"`
	ParentCertID            string `json:"parentCertId,omitempty"`
	ParentCertName          string `json:"parentCertName,omitempty"`
	PrivateKey              string `json:"privateKey,omitempty"`
	PrivateKeyPresent       bool   `json:"privateKeyPresent,omitempty"`
	SerialNo                string `json:"serialNo,omitempty"`
	ValidFromInEpochSec     string `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec       string `json:"validToInEpochSec,omitempty"`
	ZrsaEncryptedPrivateKey string `json:"zrsaencryptedprivatekey,omitempty"`
	ZrsaEncryptedSessionKey string `json:"zrsaencryptedsessionkey,omitempty"`
}

func (service *Service) Get(id string) (*EnrollmentCert, *http.Response, error) {
	v := new(EnrollmentCert)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+enrollmentCertEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(certName string) (*EnrollmentCert, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + enrollmentCertEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[EnrollmentCert](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, cert := range list {
		if strings.EqualFold(cert.Name, certName) {
			return &cert, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no signing certificate named '%s' was found", certName)
}

func (service *Service) GetAll() ([]EnrollmentCert, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + enrollmentCertEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[EnrollmentCert](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
