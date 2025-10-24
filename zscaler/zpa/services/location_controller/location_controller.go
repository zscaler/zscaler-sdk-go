package location_controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v1/admin/customers/"
	locationEndpoint      = "/location"
	locationGroupEndpoint = "/locationGroup"
)

func GetLocationExtranetResource(ctx context.Context, service *zscaler.Service, zpnErID string) (*common.CommonSummary, *http.Response, error) {
	v := new(common.CommonSummary)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+locationEndpoint+"/extranetResource", zpnErID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetLocationSummary(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + locationEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetLocationGroupExtranetResource(ctx context.Context, service *zscaler.Service, zpnErID string) (*common.LocationGroupDTO, *http.Response, error) {
	v := new(common.LocationGroupDTO)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+locationGroupEndpoint+"/extranetResource", zpnErID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
