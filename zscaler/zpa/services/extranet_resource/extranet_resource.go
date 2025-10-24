package extranet_resource

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig               = "/zpa/mgmtconfig/v1/admin/customers/"
	extranetResourceEndpoint = "/extranetResource/partner"
)

func GetExtranetResourcePartner(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + extranetResourceEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
