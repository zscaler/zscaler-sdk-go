package applicationsegment_share

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	appSegmentEndpoint = "/application"
)

type AppSegmentSharedToMicrotenant struct {
	ApplicationID       string   `json:"applicationId,omitempty"`
	ShareToMicrotenants []string `json:"shareToMicrotenants,omitempty"`
	MicroTenantID       string   `json:"microtenantId,omitempty"`
}

func (service *Service) AppSegmentMicrotenantShare(applicationID string, appSegmentRequest AppSegmentSharedToMicrotenant) (*http.Response, error) {
	microTenantID := appSegmentRequest.MicroTenantID
	if microTenantID == "" && service.microTenantID != nil {
		microTenantID = *service.microTenantID
	}

	relativeURL := fmt.Sprintf("%s%s%s/%s/share", mgmtConfig, service.Client.Config.CustomerID, appSegmentEndpoint, applicationID)

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
