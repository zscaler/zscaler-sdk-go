package urlfilteringpolicies

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	urlFilteringPoliciesEndpoint = "/zia/api/v1/urlFilteringRules"
)

type URLFilteringRule struct {
	// URL Filtering Rule ID
	ID int `json:"id,omitempty"`

	// Rule Name
	Name string `json:"name,omitempty"`

	// Order of execution of rule with respect to other URL Filtering rules
	Order int `json:"order,omitempty"`

	// Protocol criteria
	Protocols []string `json:"protocols,omitempty"`

	// List of URL categories for which rule must be applied
	URLCategories []string `json:"urlCategories"`

	UserRiskScoreLevels []string `json:"userRiskScoreLevels,omitempty"`

	// Rule State
	State string `json:"state,omitempty"`

	UserAgentTypes []string `json:"userAgentTypes,omitempty"`

	// Admin rank of the admin who creates this rule
	Rank int `json:"rank,omitempty"`

	// Request method for which the rule must be applied. If not set, rule is applied to all methods
	RequestMethods []string `json:"requestMethods,omitempty"`

	// URL of end user notification page to be displayed when the rule is matched. Not applicable if either 'overrideUsers' or 'overrideGroups' is specified.
	EndUserNotificationURL string `json:"endUserNotificationUrl,omitempty"`

	// When set to true, a 'BLOCK' action triggered by the rule could be overridden. If true and both overrideGroup and overrideUsers are not set, the BLOCK triggered by this rule could be overridden for any users. If blockOverride is not set, 'BLOCK' action cannot be overridden.
	BlockOverride bool `json:"blockOverride,omitempty"`

	// Time quota in minutes, after which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.
	TimeQuota int `json:"timeQuota,omitempty"`

	// Size quota in KB beyond which the URL Filtering rule is applied. If not set, no quota is enforced. If a policy rule action is set to 'BLOCK', this field is not applicable.
	SizeQuota int `json:"sizeQuota,omitempty"`

	// Additional information about the URL Filtering rule
	Description string `json:"description,omitempty"`

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

	// Action taken when traffic matches rule criteria
	Action string `json:"action,omitempty"`

	// If set to true, the CIPA Compliance rule is enabled
	Ciparule bool `json:"ciparule,omitempty"`

	// List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.
	DeviceTrustLevels []string `json:"deviceTrustLevels,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`

	// Who modified the rule last
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Name-ID pairs of users for which this rule can be overridden. Applicable only if blockOverride is set to 'true', action is 'BLOCK' and overrideGroups is not set.If this overrideUsers is not set, 'BLOCK' action can be overridden for any user.
	OverrideUsers []common.IDNameExtensions `json:"overrideUsers,omitempty"`

	// Name-ID pairs of groups for which this rule can be overridden. Applicable only if blockOverride is set to 'true' and action is 'BLOCK'. If this overrideGroups is not set, 'BLOCK' action can be overridden for any group.
	OverrideGroups []common.IDNameExtensions `json:"overrideGroups,omitempty"`

	// Name-ID pairs of the location groups to which the rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// The URL Filtering rule's label. Rule labels allow you to logically group your organization's policy rules. Policy rules that are not associated with a rule label are grouped under the Untagged label.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// Name-ID pairs of locations for which rule must be applied
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of groups for which rule must be applied
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// Name-ID pairs of departments for which rule must be applied
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of users for which rule must be applied
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// Source IP address groups for which the rule is applicable.
	// If not set, the rule is not restricted to a specific source IP address group.
	SourceIPGroups []common.IDNameExtensions `json:"sourceIpGroups,omitempty"`

	// Name-ID pairs of time interval during which rule must be enforced.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// The list of preconfigured workload groups to which the policy must be applied.
	WorkloadGroups []common.IDName `json:"workloadGroups,omitempty"`

	// The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules.
	// Note: This parameter is required for the ISOLATE action and is not applicable to other actions.
	CBIProfile   CBIProfile `json:"cbiProfile"`
	CBIProfileID int        `json:"cbiProfileId"`
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
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*URLFilteringRule, error) {
	var urlFilteringPolicies URLFilteringRule
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID), &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning url filtering rules from Get: %d", urlFilteringPolicies.ID)
	return &urlFilteringPolicies, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, urlFilteringPolicyName string) (*URLFilteringRule, error) {
	var urlFilteringPolicies []URLFilteringRule
	err := common.ReadAllPages(ctx, service.Client, urlFilteringPoliciesEndpoint, &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}
	for _, urlFilteringPolicy := range urlFilteringPolicies {
		if strings.EqualFold(urlFilteringPolicy.Name, urlFilteringPolicyName) {
			return &urlFilteringPolicy, nil
		}
	}
	return nil, fmt.Errorf("no url filtering rule found with name: %s", urlFilteringPolicyName)
}

func Create(ctx context.Context, service *zscaler.Service, ruleID *URLFilteringRule) (*URLFilteringRule, error) {
	resp, err := service.Client.Create(ctx, urlFilteringPoliciesEndpoint, *ruleID)
	if err != nil {
		return nil, err
	}

	createdURLFilteringRule, ok := resp.(*URLFilteringRule)
	if !ok {
		return nil, errors.New("object returned from api was not a url filtering rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning url filtering rule from create: %d", createdURLFilteringRule.ID)
	return createdURLFilteringRule, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rule *URLFilteringRule) (*URLFilteringRule, *http.Response, error) {
	// Add debug log to print the rule object
	service.Client.GetLogger().Printf("[DEBUG] Updating URL Filtering Rule with ID %d: %+v", ruleID, rule)
	if rule.CBIProfile.ID == "" || rule.CBIProfileID == 0 {
		// If CBIProfile object is empty, fetch it using GetByName as Get by ID is not currently returnign the full CBIProfile object with the uuid ID
		var urlFilteringPolicies []URLFilteringRule
		err := common.ReadAllPages(ctx, service.Client, urlFilteringPoliciesEndpoint, &urlFilteringPolicies)
		if err != nil {
			return nil, nil, err
		}
		for _, urlFilteringPolicy := range urlFilteringPolicies {
			if urlFilteringPolicy.ID == ruleID {
				rule.CBIProfile = urlFilteringPolicy.CBIProfile
				rule.CBIProfileID = urlFilteringPolicy.CBIProfileID
			}
		}
	}
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID), *rule)
	if err != nil {
		return nil, nil, err
	}
	updatedURLFilteringRule, _ := resp.(*URLFilteringRule)

	service.Client.GetLogger().Printf("[DEBUG] returning URL filtering rule from update: %d", updatedURLFilteringRule.ID)
	return updatedURLFilteringRule, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", urlFilteringPoliciesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll returns the all rules.
func GetAll(ctx context.Context, service *zscaler.Service) ([]URLFilteringRule, error) {
	var urlFilteringPolicies []URLFilteringRule
	err := common.ReadAllPages(ctx, service.Client, urlFilteringPoliciesEndpoint, &urlFilteringPolicies)
	if err != nil {
		return nil, err
	}
	return urlFilteringPolicies, nil
}
