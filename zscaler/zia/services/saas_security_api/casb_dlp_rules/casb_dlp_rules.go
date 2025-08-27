package casb_dlp_rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	casbDlpRulesEndpoint = "/zia/api/v1/casbDlpRules"
)

type CasbDLPRules struct {
	Type                          string                    `json:"type,omitempty"`
	ID                            int                       `json:"id,omitempty"`
	Order                         int                       `json:"order,omitempty"`
	Rank                          int                       `json:"rank,omitempty"`
	LastModifiedTime              int                       `json:"lastModifiedTime,omitempty"`
	Name                          string                    `json:"name,omitempty"`
	State                         string                    `json:"state,omitempty"`
	Action                        string                    `json:"action,omitempty"`
	Severity                      string                    `json:"severity,omitempty"`
	Description                   string                    `json:"description,omitempty"`
	BucketOwner                   string                    `json:"bucketOwner,omitempty"`
	ExternalAuditorEmail          string                    `json:"externalAuditorEmail,omitempty"`
	ContentLocation               string                    `json:"contentLocation,omitempty"`
	NumberOfInternalCollaborators string                    `json:"numberOfInternalCollaborators,omitempty"`
	NumberOfExternalCollaborators string                    `json:"numberOfExternalCollaborators,omitempty"`
	Recipient                     string                    `json:"recipient,omitempty"`
	QuarantineLocation            string                    `json:"quarantineLocation,omitempty"`
	AccessControl                 string                    `json:"accessControl,omitempty"`
	WatermarkDeleteOldVersion     bool                      `json:"watermarkDeleteOldVersion,omitempty"`
	IncludeCriteriaDomainProfile  bool                      `json:"includeCriteriaDomainProfile,omitempty"`
	IncludeEmailRecipientProfile  bool                      `json:"includeEmailRecipientProfile,omitempty"`
	WithoutContentInspection      bool                      `json:"withoutContentInspection,omitempty"`
	IncludeEntityGroups           bool                      `json:"includeEntityGroups,omitempty"`
	FileTypes                     []string                  `json:"fileTypes,omitempty"`
	CollaborationScope            []string                  `json:"collaborationScope,omitempty"`
	Domains                       []string                  `json:"domains,omitempty"`
	Components                    []string                  `json:"components,omitempty"`
	DeviceTrustLevels             []string                  `json:"deviceTrustLevels,omitempty"`
	ObjectTypes                   []common.IDNameExtensions `json:"objectTypes,omitempty"`
	Buckets                       []common.IDNameExtensions `json:"buckets,omitempty"`
	Labels                        []common.IDNameExtensions `json:"labels,omitempty"`
	IncludedDomainProfiles        []common.IDNameExtensions `json:"includedDomainProfiles,omitempty"`
	ExcludedDomainProfiles        []common.IDNameExtensions `json:"excludedDomainProfiles,omitempty"`
	CriteriaDomainProfiles        []common.IDNameExtensions `json:"criteriaDomainProfiles,omitempty"`
	EmailRecipientProfiles        []common.IDNameExtensions `json:"emailRecipientProfiles,omitempty"`
	Devices                       []common.IDNameExtensions `json:"devices,omitempty"`
	DeviceGroups                  []common.IDNameExtensions `json:"deviceGroups,omitempty"`
	EntityGroups                  []common.IDNameExtensions `json:"entityGroups,omitempty"`
	CloudAppTenants               []common.IDNameExtensions `json:"cloudAppTenants,omitempty"`
	Users                         []common.IDNameExtensions `json:"users,omitempty"`
	Groups                        []common.IDNameExtensions `json:"groups,omitempty"`
	Departments                   []common.IDNameExtensions `json:"departments,omitempty"`
	DLPEngines                    []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	LastModifiedBy                *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	Auditor                       *common.IDCustom          `json:"auditor,omitempty"`
	ZscalerIncidentReceiver       *common.IDCustom          `json:"zscalerIncidentReceiver,omitempty"`
	AuditorNotification           *common.IDCustom          `json:"auditorNotification,omitempty"`
	Tag                           *common.IDCustom          `json:"tag,omitempty"`
	WatermarkProfile              *common.IDCustom          `json:"watermarkProfile,omitempty"`
	RedactionProfile              *common.IDCustom          `json:"redactionProfile,omitempty"`
	CasbEmailLabel                *common.IDCustom          `json:"casbEmailLabel,omitempty"`
	CasbTombstoneTemplate         *common.IDCustom          `json:"casbTombstoneTemplate,omitempty"`
	Receiver                      *Receiver                 `json:"receiver,omitempty"`
}

type Receiver struct {
	ID     int                      `json:"id,omitempty"`
	Name   string                   `json:"name,omitempty"`
	Type   string                   `json:"type,omitempty"`
	Tenant *common.IDNameExtensions `json:"tenant,omitempty"`
}

func GetByRuleID(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int) (*CasbDLPRules, error) {
	var rule CasbDLPRules
	endpoint := fmt.Sprintf("%s/%d?ruleType=%s", casbDlpRulesEndpoint, ruleID, url.QueryEscape(ruleType))
	err := service.Client.Read(ctx, endpoint, &rule)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning casb dlp rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByRuleType(ctx context.Context, service *zscaler.Service, ruleType string) ([]CasbDLPRules, error) {
	var rules []CasbDLPRules
	endpoint := fmt.Sprintf("%s?ruleType=%s", casbDlpRulesEndpoint, url.QueryEscape(ruleType))
	err := service.Client.Read(ctx, endpoint, &rules)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning casb dlp rule from GetByRuleType: %+v", rules)
	return rules, nil
}

func Create(ctx context.Context, service *zscaler.Service, rule *CasbDLPRules) (*CasbDLPRules, error) {
	url := casbDlpRulesEndpoint
	resp, err := service.Client.Create(ctx, url, *rule)
	if err != nil {
		return nil, err
	}
	createdRules, ok := resp.(*CasbDLPRules)
	if !ok {
		return nil, errors.New("object returned from API was not a CasbDLPRules pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning casb dlp rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *CasbDLPRules) (*CasbDLPRules, error) {
	url := fmt.Sprintf("%s/%d", casbDlpRulesEndpoint, ruleID)
	resp, err := service.Client.UpdateWithPut(ctx, url, *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*CasbDLPRules)
	service.Client.GetLogger().Printf("[DEBUG] Returning casb dlp rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int) (*http.Response, error) {
	endpoint := fmt.Sprintf("%s/%d?ruleType=%s", casbDlpRulesEndpoint, ruleID, url.QueryEscape(ruleType))
	err := service.Client.Delete(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CasbDLPRules, error) {
	var rules []CasbDLPRules
	err := common.ReadAllPages(ctx, service.Client, casbDlpRulesEndpoint+"/all", &rules)
	return rules, err
}
