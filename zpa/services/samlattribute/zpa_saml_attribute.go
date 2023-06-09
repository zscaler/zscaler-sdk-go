package samlattribute

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig            = "/mgmtconfig/v2/admin/customers/"
	mgmtConfigV1          = "/mgmtconfig/v1/admin/customers/"
	samlAttributeEndpoint = "/samlAttribute"
)

type SamlAttribute struct {
	CreationTime  string `json:"creationTime,omitempty"`
	ID            string `json:"id,omitempty"`
	IdpID         string `json:"idpId,omitempty"`
	IdpName       string `json:"idpName,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
	Name          string `json:"name,omitempty"`
	SamlName      string `json:"samlName,omitempty"`
	UserAttribute bool   `json:"userAttribute,omitempty"`
}

func (service *Service) Get(samlAttributeID string) (*SamlAttribute, *http.Response, error) {
	v := new(SamlAttribute)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.Config.CustomerID+samlAttributeEndpoint, samlAttributeID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(samlAttrName string) (*SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](service.Client, relativeURL, samlAttrName)
	if err != nil {
		return nil, nil, err
	}
	for _, samlAttribute := range list {
		if samlAttribute.Name == samlAttrName {
			return &samlAttribute, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no saml attribute named '%s' was found", samlAttrName)
}

func (service *Service) GetAll() ([]SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
