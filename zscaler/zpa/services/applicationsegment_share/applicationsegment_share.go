package applicationsegment_share

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig         = "/zpa/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type AppSegmentSharedToMicrotenant struct {
	ApplicationID       string   `json:"applicationId,omitempty"`
	ShareToMicrotenants []string `json:"shareToMicrotenants,omitempty"`
	MicroTenantID       string   `json:"microtenantId,omitempty"`
}

func AppSegmentMicrotenantShare(service *zscaler.Service, applicationID string, appSegmentRequest AppSegmentSharedToMicrotenant) (*http.Response, error) {
	microTenantID := appSegmentRequest.MicroTenantID
	if microTenantID == "" && service.MicroTenantID() != nil {
		microTenantID = *service.MicroTenantID()
	}

	relativeURL := fmt.Sprintf("%s%s%s/%s/share", mgmtConfig, service.Client.GetCustomerID(), appSegmentEndpoint, applicationID)

	// Add microTenantID to the filter if it's provided
	filter := common.Filter{}
	if microTenantID != "" {
		filter.MicroTenantID = &microTenantID
	}

	resp, err := service.Client.NewRequestDo("PUT", relativeURL, filter, appSegmentRequest, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
