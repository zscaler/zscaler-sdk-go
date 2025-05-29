package samlattribute

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfigV1          = "/zpa/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2          = "/zpa/mgmtconfig/v2/admin/customers/"
	samlAttributeEndpoint = "/samlAttribute"
)

type SamlAttribute struct {
	ID            string `json:"id,omitempty"`
	CreationTime  string `json:"creationTime,omitempty"`
	IdpID         string `json:"idpId,omitempty"`
	IdpName       string `json:"idpName,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
	Name          string `json:"name,omitempty"`
	SamlName      string `json:"samlName,omitempty"`
	Delta         string `json:"delta,omitempty"`
	UserAttribute bool   `json:"userAttribute,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, samlAttributeID string) (*SamlAttribute, *http.Response, error) {
	v := new(SamlAttribute)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+samlAttributeEndpoint, samlAttributeID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, samlAttrName string) (*SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.GetCustomerID() + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](ctx, service.Client, relativeURL, "")
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

func Create(ctx context.Context, service *zscaler.Service, samlAttribute *SamlAttribute) (*SamlAttribute, *http.Response, error) {
	v := new(SamlAttribute)
	url := fmt.Sprintf("%s%s%s", mgmtConfigV1, service.Client.GetCustomerID(), samlAttributeEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "POST", url, nil, samlAttribute, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, samlAttributeID string, attribute *SamlAttribute) (*http.Response, error) {
	url := fmt.Sprintf("%s%s%s%s", mgmtConfigV1, service.Client.GetCustomerID(), samlAttributeEndpoint+"/", samlAttributeID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", url, nil, attribute, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, samlAttributeID string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s%s%s", mgmtConfigV1, service.Client.GetCustomerID(), samlAttributeEndpoint+"/", samlAttributeID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", url, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]SamlAttribute, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.GetCustomerID() + samlAttributeEndpoint)
	list, resp, err := common.GetAllPagesGeneric[SamlAttribute](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
