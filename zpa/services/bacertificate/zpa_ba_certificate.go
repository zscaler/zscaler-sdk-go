package bacertificate

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v1/zpa/services/common"
)

const (
	mgmtConfigV1                = "/mgmtconfig/v1/admin/customers/"
	baCertificateEndpoint       = "/certificate"
	mgmtConfigV2                = "/mgmtconfig/v2/admin/customers/"
	baCertificateIssuedEndpoint = "/clientlessCertificate/issued"
)

type BaCertificate struct {
	CName               string   `json:"cName,omitempty"`
	CertChain           string   `json:"certChain,omitempty"`
	CertBlob            string   `json:"certBlob,omitempty"`
	CreationTime        string   `json:"creationTime,omitempty"`
	Description         string   `json:"description,omitempty"`
	ID                  string   `json:"id,omitempty"`
	IssuedBy            string   `json:"issuedBy,omitempty"`
	IssuedTo            string   `json:"issuedTo,omitempty"`
	ModifiedBy          string   `json:"modifiedBy,omitempty"`
	ModifiedTime        string   `json:"modifiedTime,omitempty"`
	Name                string   `json:"name,omitempty"`
	San                 []string `json:"san,omitempty"`
	SerialNo            string   `json:"serialNo,omitempty"`
	Status              string   `json:"status,omitempty"`
	ValidFromInEpochSec string   `json:"validFromInEpochSec,omitempty"`
	ValidToInEpochSec   string   `json:"validToInEpochSec,omitempty"`
}

func (service *Service) Get(baCertificateID string) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetIssuedByName(certName string) (*BaCertificate, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + baCertificateIssuedEndpoint)
	list, resp, err := common.GetAllPagesGeneric[BaCertificate](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, baCertificate := range list {
		if baCertificate.Name == certName {
			return &baCertificate, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no issued certificate named '%s' was found", certName)
}

func (service *Service) GetAll() ([]BaCertificate, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + baCertificateIssuedEndpoint)
	list, resp, err := common.GetAllPagesGeneric[BaCertificate](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func (service *Service) Create(baCertificate BaCertificate) (*BaCertificate, *http.Response, error) {
	v := new(BaCertificate)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, nil, baCertificate, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(baCertificateID string, baCertificate *BaCertificate) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, baCertificate, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(baCertificateID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+baCertificateEndpoint, baCertificateID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
