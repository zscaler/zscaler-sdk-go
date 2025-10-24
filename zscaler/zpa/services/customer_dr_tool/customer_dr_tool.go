package customer_dr_tool

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                    = "/zpa/mgmtconfig/v1/admin/customers/"
	customerDRToolEndpoint string = "/customerDRToolVersion"
)

type CustomerDrTool struct {
	CreationTime string `json:"creationTime,omitempty"`
	CustomerId   string `json:"customerId,omitempty"`
	ID           string `json:"id,omitempty"`
	Latest       bool   `json:"latest,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	Platform     string `json:"platform,omitempty"`
	Version      string `json:"version,omitempty"`
}

func GetCustomerDRTool(ctx context.Context, service *zscaler.Service) ([]CustomerDrTool, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + customerDRToolEndpoint
	list, resp, err := common.GetAllPagesGeneric[CustomerDrTool](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
