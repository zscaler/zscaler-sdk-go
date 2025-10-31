package location_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v1/admin/customers/"
	locationEndpoint      = "/location"
	locationGroupEndpoint = "/locationGroup"
)

func GetLocationExtranetResource(ctx context.Context, service *zscaler.Service, zpnErID string) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+locationEndpoint+"/extranetResource", zpnErID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetLocationSummary(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + locationEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetLocationSummaryByName(ctx context.Context, service *zscaler.Service, locationName string) (*common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + locationEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, locationName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no location named '%s' was found", locationName)
}

func GetLocationGroupExtranetResource(ctx context.Context, service *zscaler.Service, zpnErID string) ([]common.LocationGroupDTO, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+locationGroupEndpoint+"/extranetResource", zpnErID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.LocationGroupDTO](ctx, service.Client, relativeURL, common.Filter{})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
