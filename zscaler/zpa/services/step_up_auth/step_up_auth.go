package step_up_auth

import (
	"context"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig = "/zpa/mgmtconfig/v1/admin/customers/"
)

type StepAuthLevel struct {
	CreationTime         string `json:"creationTime,omitempty"`
	Delta                string `json:"delta,omitempty"`
	Description          string `json:"description,omitempty"`
	IamAuthLevelID       string `json:"iamAuthLevelId,omitempty"`
	ID                   string `json:"id,omitempty"`
	ModifiedBy           string `json:"modifiedBy,omitempty"`
	ModifiedTime         string `json:"modifiedTime,omitempty"`
	Name                 string `json:"name,omitempty"`
	ParentIamAuthLevelID string `json:"parentIamAuthLevelId,omitempty"`
	MicrotenantID        string `json:"microtenantId,omitempty"`
	MicrotenantName      string `json:"microtenantName,omitempty"`
	UserMessage          string `json:"userMessage,omitempty"`
}

func GetStepupAuthLevel(ctx context.Context, service *zscaler.Service) ([]string, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + "/stepupauthlevel"

	var result []string
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &result)
	if err != nil {
		return nil, nil, err
	}
	return result, resp, nil
}
