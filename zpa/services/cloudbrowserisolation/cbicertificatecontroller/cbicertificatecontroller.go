package cbicertificatecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig               = "/cbiconfig/cbi/api/customers/"
	cbiCertificateEndpoint  = "/certificate"
	cbiCertificatesEndpoint = "/certificates"
)

type CBICertificate struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	PEM       string `json:"pem,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

func (service *Service) Get(certificateID string) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(certificateName string) (*CBICertificate, *http.Response, error) {
	list, resp, err := service.GetAll()
	if err != nil {
		return nil, nil, err
	}
	for _, profile := range list {
		if strings.EqualFold(profile.Name, certificateName) {
			return &profile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no certificate named '%s' was found", certificateName)
}

func (service *Service) Create(cbiProfile *CBICertificate) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiCertificateEndpoint, common.Filter{MicroTenantID: service.microTenantID}, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(certificateID string, certificateRequest *CBICertificate) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, certificateRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(certificateID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]CBICertificate, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiCertificatesEndpoint
	var list []CBICertificate
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &list)
	if err != nil {
		return nil, resp, err
	}
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
