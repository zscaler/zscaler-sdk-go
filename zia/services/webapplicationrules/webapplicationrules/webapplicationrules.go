package webapplicationrules

import (
	"fmt"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	webApplicationRulesEndpoint     = "/webApplicationRules"
	webApplicationRulesLiteEndpoint = "/webApplicationRules/lite"
)

type WebApplicationRules struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name,omitempty"`
	Description          string `json:"description,omitempty"`
	Actions              string `json:"actions,omitempty"`
	State                string `json:"state,omitempty"`
	Rank                 int    `json:"rank,omitempty"`
	Type                 string `json:"type,omitempty"`
	Order                int    `json:"order,omitempty"`
	TimeQuota            int    `json:"timeQuota,omitempty"`
	SizeQuota            int    `json:"sizeQuota,omitempty"`
	CascadingEnabled     bool   `json:"cascadingEnabled"`
	AccessControl        string `json:"accessControl,omitempty"`
	NumberOfApplications int    `json:"numberOfApplications"`

	EunEnabled bool `json:"eunEnabled,omitempty"`

	EunTemplateID int `json:"eunTemplateId,omitempty"`

	BrowserEunTemplateID int `json:"browserEunTemplateId,omitempty"`

	// If set to true, a predefined rule is applied
	Predefined bool `json:"predefined"`

	// If enforceTimeValidity is set to true, the URL Filtering rule is valid starting on this date and time.
	ValidityStartTime int `json:"validityStartTime,omitempty"`

	// If enforceTimeValidity is set to true, the URL Filtering rule ceases to be valid on this end date and time.
	ValidityEndTime int `json:"validityEndTime,omitempty"`

	// If enforceTimeValidity is set to true, the URL Filtering rule date and time is valid based on this time zone ID.
	ValidityTimeZoneID string `json:"validityTimeZoneId,omitempty"`

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
	Devices []common.IDNameExtensions `json:"devices"`

	// Name-ID pairs of departments for which rule must be applied
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of groups for which rule must be applied
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// The URL Filtering rule's label. Rule labels allow you to logically group your organization's policy rules. Policy rules that are not associated with a rule label are grouped under the Untagged label.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// Name-ID pairs of locations for which rule must be applied
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// Name-ID pairs of time interval during which rule must be enforced.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// Name-ID pairs of users for which rule must be applied
	SharingDomainProfiles []common.IDNameExtensions `json:"sharingDomainProfiles,omitempty"`

	FormSharingDomainProfiles []common.IDNameExtensions `json:"formSharingDomainProfiles,omitempty"`

	CloudAppInstances []common.IDNameExtensions `json:"cloudAppInstances,omitempty"`

	// The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules.
	// Note: This parameter is required for the ISOLATE action and is not applicable to other actions.
	CBIProfile   CBIProfile `json:"cbiProfile"`
	CBIProfileID int        `json:"cbiProfileId"`

	CloudAppRiskProfile CloudAppRiskProfile `json:"cloudAppRiskProfile,omitempty"`
	TenancyProfileIDs   TenancyProfileIDs   `json:"tenancyProfileIds,omitempty"`
}

type CloudAppRiskProfile struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TenancyProfileIDs struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CBIProfile struct {
	ProfileSeq int `json:"profileSeq"`

	// The universally unique identifier (UUID) for the browser isolation profile
	ID string `json:"id"`

	// Name of the browser isolation profile
	Name string `json:"name"`

	// The browser isolation profile URL
	URL string `json:"url"`

	// (Optional) Indicates whether this is a default browser isolation profile. Zscaler sets this field.
	DefaultProfile bool `json:"defaultProfile"`

	// Indicates whether sandboxMode is enabled for this profile configured in Cloud Browser Isolation
	SandboxMode bool `json:"sandboxMode"`
}

func (service *Service) Get(ruleID int) (*WebApplicationRules, error) {
	var rule WebApplicationRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", webApplicationRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}
	service.Client.Logger.Printf("[DEBUG]Returning web application rule from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetLite() ([]WebApplicationRules, error) {
	var rules []WebApplicationRules
	err := service.Client.Read(webApplicationRulesLiteEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	service.Client.Logger.Printf("[DEBUG] Returning %d web application rules lite from Get", len(rules))
	return rules, nil
}
