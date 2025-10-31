package extranet_resource

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig               = "/zpa/mgmtconfig/v1/admin/customers/"
	extranetResourceEndpoint = "/extranetResource/partner"
)

func GetExtranetResourcePartner(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + extranetResourceEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetExtranetResourcePartnerByName(ctx context.Context, service *zscaler.Service, extranetName string) (*common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + extranetResourceEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, extranetName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no extranet resource named '%s' was found", extranetName)
}
