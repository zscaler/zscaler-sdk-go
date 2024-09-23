package cbicertificatecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	cbiConfig               = "/zpa/cbiconfig/cbi/api/customers/"
	cbiCertificateEndpoint  = "/certificate"
	cbiCertificatesEndpoint = "/certificates"
)

type CBICertificate struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	PEM       string `json:"pem,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

func Get(service *zscaler.Service, certificateID string) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.GetCustomerID()+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *zscaler.Service, certificateName string) (*CBICertificate, *http.Response, error) {
	list, resp, err := GetAll(service)
	if err != nil {
		return nil, nil, err
	}
	for _, cert := range list {
		if strings.EqualFold(cert.Name, certificateName) {
			return &cert, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no certificate named '%s' was found", certificateName)
}

func GetByNameOrID(service *zscaler.Service, identifier string) (*CBICertificate, *http.Response, error) {
	// Retrieve all banners
	list, resp, err := GetAll(service)
	if err != nil {
		return nil, nil, err
	}
	// Try to find by ID
	for _, certificate := range list {
		if certificate.ID == identifier {
			return Get(service, certificate.ID)
		}
	}
	// Try to find by name
	for _, certificate := range list {
		if strings.EqualFold(certificate.Name, identifier) {
			return Get(service, certificate.ID)
		}
	}
	return nil, resp, fmt.Errorf("no isolation certificate named or with ID '%s' was found", identifier)
}

func Create(service *zscaler.Service, cbiProfile *CBICertificate) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.GetCustomerID()+cbiCertificateEndpoint, nil, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *zscaler.Service, certificateID string, certificateRequest *CBICertificate) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, certificateRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *zscaler.Service, certificateID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.GetCustomerID()+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *zscaler.Service) ([]CBICertificate, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.GetCustomerID() + cbiCertificatesEndpoint
	var list []CBICertificate
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, resp, err
	}
	return list, resp, nil
}
