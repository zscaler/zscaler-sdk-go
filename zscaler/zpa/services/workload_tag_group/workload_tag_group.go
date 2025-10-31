package workload_tag_group

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
	workloadTagGroupEndpoint = "/workloadTagGroup/summary"
)

func GetWorkloadTagGroup(ctx context.Context, service *zscaler.Service) ([]common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + workloadTagGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[common.CommonSummary](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, workloadTagGroupName string) (*common.CommonSummary, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + workloadTagGroupEndpoint
	list, resp, err := common.GetAllPagesGeneric[common.CommonSummary](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, workloadTagGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no workload tag group named '%s' was found", workloadTagGroupName)
}
