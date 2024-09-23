package samlattribute

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v2/admin/customers/"
	mgmtConfigV1          = "/zpa/mgmtconfig/v1/admin/customers/"
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

func Get(service *zscaler.Service, samlAttributeID string) (*SamlAttribute, *http.Response, error) {
	v := new(SamlAttribute)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+samlAttributeEndpoint, samlAttributeID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(service *zscaler.Service, samlAttrName string) (*SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.GetCustomerID() + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](service.Client, relativeURL, "")
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

func GetAll(service *zscaler.Service) ([]SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.GetCustomerID() + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
