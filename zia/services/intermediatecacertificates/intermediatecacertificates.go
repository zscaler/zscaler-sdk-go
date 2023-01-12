package intermediatecacertificates

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	intermediateCaCertificatesEndpoint = "/intermediateCaCertificate"
	intCADownloadAttestationEndpoint   = "/intermediateCaCertificate/downloadAttestation"
	intCADownloadCSREndpoint           = "/intermediateCaCertificate/downloadCsr"
	intCADownloadPublicKeyEndpoint     = "/intermediateCaCertificate/downloadPublicKey"
	intCAGenerateCSREndpoint           = "/intermediateCaCertificate/generateCsr"
	intCAFinalizeCSREndpoint           = "/intermediateCaCertificate/finalizeCert"
	intCAKeyPairEndpoint               = "/intermediateCaCertificate/keyPair"
	intCACertMakeDefaultEndpoint       = "/intermediateCaCertificate/makeDefault"
	intCAReadyToUseEndpoint            = "/intermediateCaCertificate/readyToUse"
	intCAShowCertEndpoint              = "/intermediateCaCertificate/showCert"
	intCAShowCSREndpoint               = "/intermediateCaCertificate/showCsr"
	intCAUploadCert                    = "/intermediateCaCertificate/uploadCert"
	intCAUploadCertChain               = "/intermediateCaCertificate/uploadCertChain"
	intCAVerifyKeyAttestation          = "/intermediateCaCertificate/verifyKeyAttestation"
)

type IntermediateCACertificate struct {
	ID                         int    `json:"id"`
	Name                       string `json:"name,omitempty"`
	Description                string `json:"description,omitempty"`
	Type                       string `json:"type,omitempty"`
	Region                     string `json:"region,omitempty"`
	Status                     string `json:"status,omitempty"`
	DefaultCertificate         bool   `json:"defaultCertificate,omitempty"`
	CertStartDate              int    `json:"certStartDate,omitempty"`
	CertExpDate                int    `json:"certExpDate,omitempty"`
	CurrentState               string `json:"currentState,omitempty"`
	PublicKey                  string `json:"publicKey,omitempty"`
	KeyGenerationTime          int    `json:"keyGenerationTime,omitempty"`
	HSMAttestationVerifiedTime int    `json:"hsmAttestationVerifiedTime,omitempty"`
	CSRFileName                string `json:"csrFileName,omitempty"`
	CSRGenerationTime          int    `json:"csrGenerationTime,omitempty"`
}

type CertSigningRequest struct {
	CertID               int    `json:"certId"`
	CSRFileName          string `json:"csrFileName,omitempty"`
	CommName             string `json:"commName,omitempty"`
	ORGName              string `json:"orgName,omitempty"`
	DeptName             string `json:"deptName,omitempty"`
	City                 string `json:"city,omitempty"`
	State                string `json:"state,omitempty"`
	Country              string `json:"country,omitempty"`
	KeySize              int    `json:"keySize,omitempty"`
	SignatureAlgorithm   string `json:"signatureAlgorithm,omitempty"`
	PathLengthConstraint int    `json:"pathLengthConstraint,omitempty"`
}

func (service *Service) Get(certID int) (*IntermediateCACertificate, error) {
	var intermediateCACertificate IntermediateCACertificate
	err := service.Client.Read(fmt.Sprintf("%s/%d", intermediateCaCertificatesEndpoint, certID), &intermediateCACertificate)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning intermediate ca certificate from Get: %d", intermediateCACertificate.ID)
	return &intermediateCACertificate, nil
}

func (service *Service) GetByName(certName string) (*IntermediateCACertificate, error) {
	var intermediateCACertificate []IntermediateCACertificate
	err := common.ReadAllPages(service.Client, intermediateCaCertificatesEndpoint, &intermediateCACertificate)
	if err != nil {
		return nil, err
	}
	for _, certificate := range intermediateCACertificate {
		if strings.EqualFold(certificate.Name, certName) {
			return &certificate, nil
		}
	}
	return nil, fmt.Errorf("no intermediate ca certificate found with name: %s", certName)
}

func (service *Service) GetDownloadAttestation(certID int) (*IntermediateCACertificate, error) {
	var downloadAttestation IntermediateCACertificate
	err := service.Client.Read(fmt.Sprintf("%s/%d", intCADownloadAttestationEndpoint, certID), &downloadAttestation)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning downloaded attestation from Get: %d", downloadAttestation.ID)
	return &downloadAttestation, nil
}

func (service *Service) GetDownloadCSR(certID int) (*IntermediateCACertificate, error) {
	var downloadCSR IntermediateCACertificate
	err := service.Client.Read(fmt.Sprintf("%s/%d", intCADownloadCSREndpoint, certID), &downloadCSR)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning downloaded csr from Get: %d", downloadCSR.ID)
	return &downloadCSR, nil
}

func (service *Service) GetDownloadPublicKey(certID int) (*IntermediateCACertificate, error) {
	var downloadPublicKey IntermediateCACertificate
	err := service.Client.Read(fmt.Sprintf("%s/%d", intCADownloadPublicKeyEndpoint, certID), &downloadPublicKey)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning downloaded public key from Get: %d", downloadPublicKey.ID)
	return &downloadPublicKey, nil
}

func (service *Service) GetIntCAReadyToUse() (*IntermediateCACertificate, error) {
	var readyToUse IntermediateCACertificate
	err := service.Client.Read((intCAReadyToUseEndpoint), &readyToUse)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning downloaded public key from Get: %d", readyToUse.ID)
	return &readyToUse, nil
}

func (service *Service) GetShowCert(certID int) (*CertSigningRequest, error) {
	var showCert CertSigningRequest
	err := service.Client.Read(fmt.Sprintf("%s/%d", intCAShowCertEndpoint, certID), &showCert)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning info about signed intrermediate CA certificates from Get: %d", showCert.CertID)
	return &showCert, nil
}

func (service *Service) GetShowCSR(certID int) (*CertSigningRequest, error) {
	var showCSR CertSigningRequest
	err := service.Client.Read(fmt.Sprintf("%s/%d", intCAShowCSREndpoint, certID), &showCSR)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning info about signed intermediate CA certificates from Get: %d", showCSR.CertID)
	return &showCSR, nil
}

func (service *Service) GetAll() ([]IntermediateCACertificate, error) {
	var intermediateCACertificate []IntermediateCACertificate
	err := common.ReadAllPages(service.Client, intermediateCaCertificatesEndpoint, &intermediateCACertificate)
	return intermediateCACertificate, err
}

func (service *Service) CreateIntCACertificate(cert *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intermediateCaCertificatesEndpoint, *cert)
	if err != nil {
		return nil, err
	}

	createdIntermediateCACert, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from create: %d", createdIntermediateCACert.ID)
	return createdIntermediateCACert, nil
}

func (service *Service) CreateIntCAGenerateCSR(cert *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAGenerateCSREndpoint, *cert)
	if err != nil {
		return nil, err
	}

	createdIntCAGenerateCSR, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from create: %d", createdIntCAGenerateCSR.ID)
	return createdIntCAGenerateCSR, nil
}

func (service *Service) CreateIntCAFinalizeCert(cert *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAFinalizeCSREndpoint, *cert)
	if err != nil {
		return nil, err
	}

	createdIntCAFinalizeCSR, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from create: %d", createdIntCAFinalizeCSR.ID)
	return createdIntCAFinalizeCSR, nil
}

func (service *Service) CreateIntCAKeyPair(keyPair *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAKeyPairEndpoint, *keyPair)
	if err != nil {
		return nil, err
	}

	createdIntCAKeyPair, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from create: %d", createdIntCAKeyPair.ID)
	return createdIntCAKeyPair, nil
}

func (service *Service) CreateUploadCert(certID *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAUploadCert, *certID)
	if err != nil {
		return nil, err
	}

	createdIntCAUploadCert, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning uploaded customer intermediate ca certificate from create: %d", createdIntCAUploadCert.ID)
	return createdIntCAUploadCert, nil
}

func (service *Service) CreateUploadCertChain(certID *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAUploadCertChain, *certID)
	if err != nil {
		return nil, err
	}

	createdIntCAUploadChain, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning uploaded certificate chain from create: %d", createdIntCAUploadChain.ID)
	return createdIntCAUploadChain, nil
}

func (service *Service) CreateVerifyKeyAttestation(certID *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intCAVerifyKeyAttestation, *certID)
	if err != nil {
		return nil, err
	}

	createdVerifyKeyAttestation, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning key attestation from create: %d", createdVerifyKeyAttestation.ID)
	return createdVerifyKeyAttestation, nil
}

func (service *Service) Update(certID int, certificates *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", intermediateCaCertificatesEndpoint, certID), *certificates)
	if err != nil {
		return nil, err
	}
	updatedIntermediateCaCert, _ := resp.(*IntermediateCACertificate)
	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from update: %d", updatedIntermediateCaCert.ID)
	return updatedIntermediateCaCert, nil
}

func (service *Service) UpdateMakeDefault(certID int, certificates *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", intCACertMakeDefaultEndpoint, certID), *certificates)
	if err != nil {
		return nil, err
	}
	updatedIntermediateCaCert, _ := resp.(*IntermediateCACertificate)
	service.Client.Logger.Printf("[DEBUG]returning default certificate from update: %d", updatedIntermediateCaCert.ID)
	return updatedIntermediateCaCert, nil
}
func (service *Service) Delete(certID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", intermediateCaCertificatesEndpoint, certID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
