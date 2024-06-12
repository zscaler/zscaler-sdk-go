package cbicertificatecontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
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

func Get(service *services.Service, certificateID string) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *services.Service, certificateName string) (*CBICertificate, *http.Response, error) {
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

func Create(service *services.Service, cbiProfile *CBICertificate) (*CBICertificate, *http.Response, error) {
	v := new(CBICertificate)
	resp, err := service.Client.NewRequestDo("POST", cbiConfig+service.Client.Config.CustomerID+cbiCertificateEndpoint, nil, cbiProfile, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *services.Service, certificateID string, certificateRequest *CBICertificate) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, certificateRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *services.Service, certificateID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", cbiConfig+service.Client.Config.CustomerID+cbiCertificatesEndpoint, certificateID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *services.Service) ([]CBICertificate, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + cbiCertificatesEndpoint
	var list []CBICertificate
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &list)
	if err != nil {
		return nil, resp, err
	}
	return list, resp, nil
}
