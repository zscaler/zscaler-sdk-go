package endpoint_dlp_rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
)

const (
	endPointDlpRulesEndpoint = "/zia/api/v1/endPointDlpRules"
)

type EndpointDlpRules struct {
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
	SubRules                  []EndpointDlpRules                 `json:"subRules,omitempty"`
}

type FileTypeCategories struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Parent string `json:"parent,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*EndpointDlpRules, error) {
	var endpointDlpRule EndpointDlpRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", endPointDlpRulesEndpoint, ruleID), &endpointDlpRule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning web dlp rule from Get: %d", endpointDlpRule.ID)
	return &endpointDlpRule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*EndpointDlpRules, error) {
	endpointDlpRules, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	for _, endpointDlpRule := range endpointDlpRules {
		if strings.EqualFold(endpointDlpRule.Name, ruleName) {
			return &endpointDlpRule, nil
		}
	}
	return nil, fmt.Errorf("no endpoint dlp rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, dlpRule *EndpointDlpRules) (*EndpointDlpRules, *http.Response, error) {
	resp, err := service.Client.Create(ctx, endPointDlpRulesEndpoint, *dlpRule)
	if err != nil {
		return nil, nil, err
	}

	createdEndpointDlpRule, ok := resp.(*EndpointDlpRules)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a endpoint dlp rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new endpoint dlp rule from create: %d", createdEndpointDlpRule.ID)
	return createdEndpointDlpRule, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, dlpRule *EndpointDlpRules) (*EndpointDlpRules, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", endPointDlpRulesEndpoint, ruleID), *dlpRule)
	if err != nil {
		return nil, nil, err
	}
	updatedOutboundEmail, _ := resp.(*EndpointDlpRules)

	service.Client.GetLogger().Printf("[DEBUG]returning updates outbound email dlp policy from update: %d", updatedOutboundEmail.ID)
	return updatedOutboundEmail, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", endPointDlpRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]EndpointDlpRules, error) {
	var endpointDlpRules []EndpointDlpRules
	err := service.Client.Read(ctx, endPointDlpRulesEndpoint, &endpointDlpRules)
	return endpointDlpRules, err
}

// GetFileTypeCategories retrieves the list of file types that are available in
// the Endpoint DLP policy rule criteria.
//
// The external parameter is optional. The file types available to select in an
// Endpoint DLP policy rule vary depending on whether Content Matching is enabled
// and DLP engines are used: setting true retrieves the list of file types
// available when DLP engines are used with Content Matching enabled, while
// setting false retrieves a different set of file types available when DLP
// engines are not used. Pass nil to omit the parameter.
func GetFileTypeCategories(ctx context.Context, service *zscaler.Service, external *bool) ([]FileTypeCategories, error) {
	var fileTypeCategories []FileTypeCategories
	endpoint := endPointDlpRulesEndpoint + "/fileTypeCategories"

	if external != nil {
		queryParams := url.Values{}
		queryParams.Set("external", strconv.FormatBool(*external))
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &fileTypeCategories)
	return fileTypeCategories, err
}
