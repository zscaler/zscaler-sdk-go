package endpoint_dlp_sub_rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

const (
	endPointDlpSubRulesEndpoint = "/zia/api/v1/endPointDlpRules"
)

type EndpointDlpSubRules struct {
	ID                        int                                `json:"id,omitempty"`
	Name                      string                             `json:"name,omitempty"`
	State                     string                             `json:"state,omitempty"`
	Order                     int                                `json:"order,omitempty"`
	Rank                      int                                `json:"rank,omitempty"`
	FileTypes                 []string                           `json:"fileTypes,omitempty"`
	DataTransferMethod        string                             `json:"dataTransferMethod,omitempty"`
	Description               string                             `json:"description,omitempty"`
	MinSize                   int                                `json:"minSize,omitempty"`
	DeviceTrustLevels         []string                           `json:"deviceTrustLevels,omitempty"`
	Action                    string                             `json:"action,omitempty"`
	ExternalAuditorEmail      string                             `json:"externalAuditorEmail,omitempty"`
	LastModifiedTime          int                                `json:"lastModifiedTime,omitempty"`
	ParentRule                int                                `json:"parentRule,omitempty"`
	Severity                  string                             `json:"severity,omitempty"`
	UserRiskScoreLevels       []string                           `json:"userRiskScoreLevels,omitempty"`
	EunEnabled                bool                               `json:"eunEnabled,omitempty"`
	EunTemplateId             int                                `json:"eunTemplateId,omitempty"`
	UcTemplateId              int                                `json:"ucTemplateId,omitempty"`
	NetworkType               string                             `json:"networkType,omitempty"`
	WithoutContentInspection  bool                               `json:"withoutContentInspection,omitempty"`
	NotificationTemplate      *common.IDNameExtensions           `json:"notificationTemplate,omitempty"`
	Auditor                   *common.IDNameExtensions           `json:"auditor,omitempty"`
	LastModifiedBy            *common.IDNameExtensions           `json:"lastModifiedBy,omitempty"`
	Resources                 []common.IDNameExtensions          `json:"resources,omitempty"`
	ResourceGroups            []common.IDNameExtensions          `json:"resourceGroups,omitempty"`
	Labels                    []common.IDNameExtensions          `json:"labels,omitempty"`
	DlpEngines                []common.IDNameExtensions          `json:"dlpEngines,omitempty"`
	Users                     []common.IDNameExtensions          `json:"users,omitempty"`
	Groups                    []common.IDNameExtensions          `json:"groups,omitempty"`
	Departments               []common.IDNameExtensions          `json:"departments,omitempty"`
	Devices                   []common.IDNameExtensions          `json:"devices,omitempty"`
	DeviceGroups              []common.IDNameExtensions          `json:"deviceGroups,omitempty"`
	Receiver                  *dlp_web_rules.Receiver            `json:"receiver,omitempty"`
	EndPointApplications      []common.EndPointApplications      `json:"endPointApplications,omitempty"`
	EndPointApplicationGroups []common.EndPointApplicationGroups `json:"endPointApplicationGroups,omitempty"`
}

func Create(ctx context.Context, service *zscaler.Service, ruleID int, dlpRule *EndpointDlpSubRules) (*EndpointDlpSubRules, *http.Response, error) {
	resp, err := service.Client.Create(ctx, fmt.Sprintf("%s/%d/subRule", endPointDlpSubRulesEndpoint, ruleID), *dlpRule)
	if err != nil {
		return nil, nil, err
	}

	createdSubRule, ok := resp.(*EndpointDlpSubRules)
	if !ok {
		return nil, nil, errors.New("object returned from api was not an endpoint dlp sub-rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new endpoint dlp sub-rule from create: %d", createdSubRule.ID)
	return createdSubRule, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID, subRuleID int, dlpRule *EndpointDlpSubRules) (*EndpointDlpSubRules, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d/subRule/%d", endPointDlpSubRulesEndpoint, ruleID, subRuleID), *dlpRule)
	if err != nil {
		return nil, nil, err
	}
	updatedSubRule, _ := resp.(*EndpointDlpSubRules)

	service.Client.GetLogger().Printf("[DEBUG]returning updated endpoint dlp sub-rule from update: %d", updatedSubRule.ID)
	return updatedSubRule, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID, subRuleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d/subRule/%d", endPointDlpSubRulesEndpoint, ruleID, subRuleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
