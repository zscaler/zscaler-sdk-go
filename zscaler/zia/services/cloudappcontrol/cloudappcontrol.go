package cloudappcontrol

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	webApplicationRulesEndpoint = "/zia/api/v1/webApplicationRules"
)

type WebApplicationRules struct {
	ID                   int      `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Actions              []string `json:"actions,omitempty"`
	State                string   `json:"state,omitempty"`
	Rank                 int      `json:"rank,omitempty"`
	Type                 string   `json:"type,omitempty"`
	Order                int      `json:"order,omitempty"`
	TimeQuota            int      `json:"timeQuota,omitempty"`
	SizeQuota            int      `json:"sizeQuota,omitempty"`
	CascadingEnabled     bool     `json:"cascadingEnabled,omitempty"`
	AccessControl        string   `json:"accessControl,omitempty"`
	Applications         []string `json:"applications,omitempty"`
	NumberOfApplications int      `json:"numberOfApplications,omitempty"`

	EunEnabled bool `json:"eunEnabled,omitempty"`

	EunTemplateID int `json:"eunTemplateId,omitempty"`

	BrowserEunTemplateID int `json:"browserEunTemplateId,omitempty"`

	// If set to true, a predefined rule is applied
	Predefined bool `json:"predefined,omitempty"`

	// If enforceTimeValidity is set to true, the URL Filtering rule is valid starting on this date and time.
	ValidityStartTime int `json:"validityStartTime,omitempty"`

	// If enforceTimeValidity is set to true, the URL Filtering rule ceases to be valid on this end date and time.
	ValidityEndTime int `json:"validityEndTime,omitempty"`

	// If enforceTimeValidity is set to true, the URL Filtering rule date and time is valid based on this time zone ID.
	ValidityTimeZoneID string `json:"validityTimeZoneId,omitempty"`

	UserAgentTypes []string `json:"userAgentTypes,omitempty"`

	// When the rule was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Enforce a set a validity time period for the URL Filtering rule. To learn more, see Configuring the URL Filtering Policy.
	EnforceTimeValidity bool `json:"enforceTimeValidity,omitempty"`

	// List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.
	DeviceTrustLevels []string `json:"deviceTrustLevels,omitempty"`

	UserRiskScoreLevels []string `json:"userRiskScoreLevels,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices,omitempty"`

	// Name-ID pairs of departments for which rule must be applied
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of groups for which rule must be applied
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// The URL Filtering rule's label. Rule labels allow you to logically group your organization's policy rules. Policy rules that are not associated with a rule label are grouped under the Untagged label.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// Name-ID pairs of users for which rule must be applied
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// Name-ID pairs of locations for which rule must be applied
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// Name-ID pairs of time interval during which rule must be enforced.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	CloudAppInstances []CloudAppInstances `json:"cloudAppInstances,omitempty"`

	//FormSharingDomainProfiles []common.CommonEntityReferences `json:"formSharingDomainProfiles,omitempty"`
	//SharingDomainProfiles     []common.CommonEntityReferences `json:"sharingDomainProfiles,omitempty"`
	TenancyProfileIDs   []common.IDNameExtensions `json:"tenancyProfileIds,omitempty"`
	CloudAppRiskProfile *common.IDCustom          `json:"cloudAppRiskProfile,omitempty"`

	// The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules.
	// Note: This parameter is required for the ISOLATE action and is not applicable to other actions.
	CBIProfile CBIProfile `json:"cbiProfile,omitempty"`
}

type CBIProfile struct {
	ProfileSeq int `json:"profileSeq,omitempty"`

	// The universally unique identifier (UUID) for the browser isolation profile
	ID string `json:"id,omitempty"`

	// Name of the browser isolation profile
	Name string `json:"name,omitempty"`

	// The browser isolation profile URL
	URL string `json:"url,omitempty"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile,omitempty"`

	// Indicates whether sandboxMode is enabled for this profile configured in Cloud Browser Isolation
	SandboxMode bool `json:"sandboxMode,omitempty"`
}

type CloudAppInstances struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type CloudApp struct {
	Val                 int    `json:"val,omitempty"`
	WebApplicationClass string `json:"webApplicationClass,omitempty"`
	BackendName         string `json:"backendName,omitempty"`
	OriginalName        string `json:"originalName,omitempty"`
	Name                string `json:"name,omitempty"`
	Deprecated          bool   `json:"deprecated,omitempty"`
	Misc                bool   `json:"misc,omitempty"`
	AppNotReady         bool   `json:"appNotReady,omitempty"`
	UnderMigration      bool   `json:"underMigration,omitempty"`
	AppCatModified      bool   `json:"appCatModified,omitempty"`
}

type AvailableActionsRequest struct {
	CloudApps []string `json:"cloudApps,omitempty"`
	Type      string   `json:"type,omitempty"`
}

func GetByRuleID(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int) (*WebApplicationRules, error) {
	var rule WebApplicationRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s/%d", webApplicationRulesEndpoint, ruleType, ruleID), &rule)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG]Returning web application rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByRuleType(ctx context.Context, service *zscaler.Service, ruleType string) ([]WebApplicationRules, error) {
	var rules []WebApplicationRules
	url := fmt.Sprintf("%s/%s", webApplicationRulesEndpoint, ruleType)
	err := service.Client.Read(ctx, url, &rules)
	if err != nil {
		return nil, err
	}
	service.Client.GetLogger().Printf("[DEBUG] Returning web application rules from GetByRuleType: %+v", rules)
	return rules, nil
}

func Create(ctx context.Context, service *zscaler.Service, ruleType string, rule *WebApplicationRules) (*WebApplicationRules, error) {
	url := fmt.Sprintf("%s/%s", webApplicationRulesEndpoint, ruleType)
	resp, err := service.Client.Create(ctx, url, *rule)
	if err != nil {
		return nil, err
	}
	createdRules, ok := resp.(*WebApplicationRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int, rules *WebApplicationRules) (*WebApplicationRules, error) {
	url := fmt.Sprintf("%s/%s/%d", webApplicationRulesEndpoint, ruleType, ruleID)
	resp, err := service.Client.UpdateWithPut(ctx, url, *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*WebApplicationRules)
	service.Client.GetLogger().Printf("[DEBUG]returning forwarding rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s/%d", webApplicationRulesEndpoint, ruleType, ruleID)
	err := service.Client.Delete(ctx, url)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CreateDuplicate(ctx context.Context, service *zscaler.Service, ruleType string, ruleID int, newName string) (*WebApplicationRules, error) {
	url := fmt.Sprintf("%s/%s/duplicate/%d?name=%s", webApplicationRulesEndpoint, ruleType, ruleID, newName)
	resp, err := service.Client.Create(ctx, url, nil) // Assuming the body is nil for duplication
	if err != nil {
		return nil, err
	}
	createdRule, ok := resp.(*WebApplicationRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG] returning rule from create duplicate: %d", createdRule.ID)
	return createdRule, nil
}

func AllAvailableActions(ctx context.Context, service *zscaler.Service, ruleType string, payload AvailableActionsRequest) ([]string, error) {
	service.Client.GetLogger().Printf("[DEBUG] AllAvailableActions called with ruleType: %s and payload: %+v", ruleType, payload)

	// Marshal the payload into a JSON string
	payloadData, err := json.Marshal(payload)
	if err != nil {
		service.Client.GetLogger().Printf("[DEBUG] error marshalling payload: %v", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/%s/availableActions", webApplicationRulesEndpoint, ruleType)

	// Use CreateWithRawPayload to send the request
	resp, err := service.Client.CreateWithRawPayload(ctx, url, string(payloadData))
	if err != nil {
		service.Client.GetLogger().Printf("[DEBUG] error creating request: %v", err)
		return nil, err
	}

	// Unmarshal the response into a slice of strings
	var availableActions []string
	if err := json.Unmarshal(resp, &availableActions); err != nil {
		service.Client.GetLogger().Printf("[DEBUG] error unmarshalling response: %v", err)
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] returning available actions: %+v", availableActions)
	return availableActions, nil
}

// GetRuleTypeMapping retrieves the rule type mapping from the API.
func GetRuleTypeMapping(ctx context.Context, service *zscaler.Service) (map[string]string, error) {
	// Initialize a map to hold the response
	var ruleTypeMapping map[string]string

	// Perform the GET request
	err := service.Client.Read(ctx, webApplicationRulesEndpoint+"/ruleTypeMapping", &ruleTypeMapping)
	if err != nil {
		return nil, err
	}

	// Log the retrieved data
	service.Client.GetLogger().Printf("[DEBUG] Returning web application rule type mapping: %+v", ruleTypeMapping)

	return ruleTypeMapping, nil
}
