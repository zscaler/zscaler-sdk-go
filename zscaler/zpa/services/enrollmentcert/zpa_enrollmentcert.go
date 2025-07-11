package enrollmentcert

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfigV1           = "/zpa/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2           = "/zpa/mgmtconfig/v2/admin/customers/"
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
	MicrotenantID           string `json:"microtenantId,omitempty"`
}

type GenerateEnrollmentCSR struct {
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	ZRSAEncryptedPrivateKey string `json:"zrsaencryptedprivatekey,omitempty"`
	CSR                     string `json:"csr,omitempty"`
}

type GenerateSelfSignedCert struct {
	Name                    string `json:"name,omitempty"`
	Description             string `json:"description,omitempty"`
	ValidFromInEpochSec     string `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec       string `json:"validToInEpochSec,omitempty"`
	RootCertificateID       string `json:"rootCertificateId,omitempty"`
	MicrotenantID           string `json:"microtenantId,omitempty"`
	ZRSAEncryptedPrivateKey string `json:"zrsaencryptedprivatekey,omitempty"`
	Certificate             string `json:"certificate,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, id string) (*EnrollmentCert, *http.Response, error) {
	v := new(EnrollmentCert)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint, id)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, certName string) (*EnrollmentCert, *http.Response, error) {
	relativeURL := mgmtConfigV2 + service.Client.GetCustomerID() + enrollmentCertEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[EnrollmentCert](ctx, service.Client, relativeURL, common.Filter{Search: certName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, cert *EnrollmentCert) (*EnrollmentCert, *http.Response, error) {
	v := new(EnrollmentCert)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint, nil, cert, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, certID string, cert *EnrollmentCert) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint, certID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, nil, cert, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, certID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint, certID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]EnrollmentCert, *http.Response, error) {
	relativeURL := mgmtConfigV2 + service.Client.GetCustomerID() + enrollmentCertEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[EnrollmentCert](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GenerateCSR(ctx context.Context, service *zscaler.Service, cert *GenerateEnrollmentCSR) (*GenerateEnrollmentCSR, *http.Response, error) {
	v := new(GenerateEnrollmentCSR)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint+"/csr/generate", nil, cert, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GenerateSelfSigned(ctx context.Context, service *zscaler.Service, cert *GenerateSelfSignedCert) (*GenerateSelfSignedCert, *http.Response, error) {
	v := new(GenerateSelfSignedCert)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfigV1+service.Client.GetCustomerID()+enrollmentCertEndpoint+"/selfsigned/generate", nil, cert, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
