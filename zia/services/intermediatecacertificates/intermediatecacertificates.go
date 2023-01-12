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

func (service *Service) GetAll() ([]IntermediateCACertificate, error) {
	var intermediateCACertificate []IntermediateCACertificate
	err := common.ReadAllPages(service.Client, intermediateCaCertificatesEndpoint, &intermediateCACertificate)
	return intermediateCACertificate, err
}

func (service *Service) Create(rule *IntermediateCACertificate) (*IntermediateCACertificate, error) {
	resp, err := service.Client.Create(intermediateCaCertificatesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdIntermediateCaCert, ok := resp.(*IntermediateCACertificate)
	if !ok {
		return nil, errors.New("object returned from api was not an intermediate ca certificate Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning intermediate ca certificate from create: %d", createdIntermediateCaCert.ID)
	return createdIntermediateCaCert, nil
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

func (service *Service) Delete(certID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", intermediateCaCertificatesEndpoint, certID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
