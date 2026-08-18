package one_identity_controller

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                    = "/zpa/mgmtconfig/v1/admin/customers/"
	oneIdentityControllerEndpoint = "/iamidpmapping"
)

type OneIdentityController struct {
	ID                     string                 `json:"id"`
	Sequence               string                 `json:"sequence"`
	IamIdpIdToIamIdMapping IamIdpIdToIamIdMapping `json:"iamIdpIdToIamIdMapping"`
}

type IamIdpIdToIamIdMapping struct {
	DeliveryTag int       `json:"deliveryTag"`
	OrgId       string    `json:"orgId"`
	Mappings    []Mapping `json:"mappings"`
}

type Mapping struct {
	IAMIDPId string `json:"iamIdpId"`
	IDPId    string `json:"idpId"`
	IDPName  string `json:"idpName"`
}

func Get(ctx context.Context, service *zscaler.Service) (*OneIdentityController, *http.Response, error) {
	v := new(OneIdentityController)
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + oneIdentityControllerEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
