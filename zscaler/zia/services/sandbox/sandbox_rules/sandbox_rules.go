package sandbox_rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

const (
	sandboxEndpoint = "/zia/api/v1/sandboxRules"
)

type SandboxRules struct {
	// Unique identifier generated for the rule
	ID int `json:"id,omitempty"`

	// The name of the Sandbox rule
	Name string `json:"name,omitempty"`

	// Additional information about the Sandbox rule
	Description string `json:"description,omitempty"`

	// Rule State
	State string `json:"state,omitempty"`

	// Policy rules are evaluated in ascending numerical order (Rule 1 before Rule 2, and so on), and this field specifies the order of execution for the rule.
	Order int `json:"order,omitempty"`

	// The action configured for the rule that must take place if the traffic matches the rule criteria. Supported values: "ALLOW" and "BLOCK"
	BaRuleAction string `json:"baRuleAction,omitempty"`

	// A Boolean value indicating whether a First-Time Action is specifically configured for the rule.
	FirstTimeEnable bool `json:"firstTimeEnable"`

	// The action that must take place when users download unknown files for the first time
	// Supported values: "ALLOW_SCAN", "QUARANTINE", "ALLOW_NOSCAN", "QUARANTINE_ISOLATE"
	FirstTimeOperation string `json:"firstTimeOperation"`

	// When set to true, this indicates that 'Machine Learning Intelligence Action' checkbox has been checked on
	MLActionEnabled bool `json:"mlActionEnabled"`

	ByThreatScore int `json:"byThreatScore,omitempty"`

	// The access privilege for this DLP policy rule based on the admin's state.
	AccessControl string `json:"accessControl,omitempty"`

	// The protocols to which the rule applies
	Protocols []string `json:"protocols,omitempty"`

	// Admin rank of the admin who creates this rule
	Rank int `json:"rank,omitempty"`

	// The threat categories to which the rule applies
	BaPolicyCategories []string `json:"baPolicyCategories,omitempty"`

	// The list of file types to which the Sandbox Rule must be applied.
	FileTypes []string `json:"fileTypes,omitempty"`

	// When the rule was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Who modified the rule last
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Name-ID pairs of locations for which rule must be applied
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// Name-ID pairs of groups for which rule must be applied
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// Name-ID pairs of departments for which rule must be applied
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of users for which rule must be applied
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// Name-ID pairs of time interval during which rule must be enforced.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// The URL Filtering rule's label. Rule labels allow you to logically group your organization's policy rules. Policy rules that are not associated with a rule label are grouped under the Untagged label.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`

	// The list of URL categories to which the DLP policy rule must be applied.
	URLCategories []string `json:"urlCategories,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.ZPAAppSegments `json:"zpaAppSegments"`

	// The cloud browser isolation profile to which the ISOLATE action is applied in the URL Filtering Policy rules.
	// Note: This parameter is required for the ISOLATE action and is not applicable to other actions.
	CBIProfile   urlfilteringpolicies.CBIProfile `json:"cbiProfile"`
	CBIProfileID int                             `json:"cbiProfileId"`

	// If set to true, the default rule is applied
	DefaultRule bool `json:"defaultRule"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*SandboxRules, error) {
	var rule SandboxRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", sandboxEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning Sandbox rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*SandboxRules, error) {
	var rules []SandboxRules
	err := common.ReadAllPages(ctx, service.Client, sandboxEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no firewall rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *SandboxRules) (*SandboxRules, error) {
	resp, err := service.Client.Create(ctx, sandboxEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*SandboxRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning sandbox rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rule *SandboxRules) (*SandboxRules, *http.Response, error) {
	// Add debug log to print the rule object
	service.Client.GetLogger().Printf("[DEBUG] Updating Sandbox rule with ID %d: %+v", ruleID, rule)
	if rule.CBIProfile.ID == "" || rule.CBIProfileID == 0 {
		// If CBIProfile object is empty, fetch it using GetByName as Get by ID is not currently returnign the full CBIProfile object with the uuid ID
		var sandboxRules []SandboxRules
		err := common.ReadAllPages(ctx, service.Client, sandboxEndpoint, &sandboxRules)
		if err != nil {
			return nil, nil, err
		}
		for _, sandboxPolicy := range sandboxRules {
			if sandboxPolicy.ID == ruleID {
				rule.CBIProfile = sandboxPolicy.CBIProfile
				rule.CBIProfileID = sandboxPolicy.CBIProfileID
			}
		}
	}
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", sandboxEndpoint, ruleID), *rule)
	if err != nil {
		return nil, nil, err
	}
	updatedSandboxRule, _ := resp.(*SandboxRules)

	service.Client.GetLogger().Printf("[DEBUG] returning Sandbox rule from update: %d", updatedSandboxRule.ID)
	return updatedSandboxRule, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", sandboxEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll returns the all rules.
func GetAll(ctx context.Context, service *zscaler.Service) ([]SandboxRules, error) {
	var sandboxPolicies []SandboxRules
	err := common.ReadAllPages(ctx, service.Client, sandboxEndpoint, &sandboxPolicies)
	if err != nil {
		return nil, err
	}
	return sandboxPolicies, nil
}
