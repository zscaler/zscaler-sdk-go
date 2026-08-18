package customer_domain_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                       = "/zpa/mgmtconfig/v1/admin/customers/"
	customerDomainControllerEndpoint = "/v2/associationtype"
)

type CustomerDomainController struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	AssociationType string `json:"associationType"`
	Capture         bool   `json:"capture"`
	CreationTime    string `json:"creationTime"`
	Domain          string `json:"domain"`
	ModifiedBy      string `json:"modifiedBy"`
	ModifiedTime    string `json:"modifiedTime"`
	MicrotenantID   string `json:"microtenantId"`
}

// GetCustomerDomainController returns the domains associated with a customer for
// the given association type. The endpoint returns a bare JSON array (no
// pagination envelope) and supports the optional microtenantId query parameter.
func GetCustomerDomainController(ctx context.Context, service *zscaler.Service, associationType string) ([]CustomerDomainController, *http.Response, error) {
	var v []CustomerDomainController
	relativeURL := fmt.Sprintf("%s/%s/domains", mgmtConfig+service.Client.GetCustomerID()+customerDomainControllerEndpoint, associationType)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
