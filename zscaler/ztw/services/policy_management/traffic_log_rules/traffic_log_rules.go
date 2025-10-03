package traffic_log_rules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	logRulesEndpoint = "/ztw/api/v1/ecRules/self"
)

type ECTrafficLogRules struct {
	// A unique identifier assigned to the forwarding rule
	ID int `json:"id,omitempty"`

	// The name of the forwarding rule
	Name string `json:"name,omitempty"`

	// Additional information about the forwarding rule
	Description string `json:"description,omitempty"`

	// The order of execution for the forwarding rule order
	Order int `json:"order,omitempty"`

	// Admin rank assigned to the forwarding rule
	Rank int `json:"rank,omitempty"`

	// Indicates whether the forwarding rule is enabled or disabled
	// Supported Values: DISABLED and ENABLED
	State string `json:"state,omitempty"`

	// The rule type selected from the available options
	Type string `json:"type,omitempty"`

	// The rule action selected from the available options
	ForwardMethod string `json:"forwardMethod,omitempty"`

	DefaultRule bool `json:"defaultRule,omitempty"`

	// Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.
	ProxyGateway *common.CommonIDName `json:"proxyGateway,omitempty"`

	// Timestamp when the rule was last modified. This field is not applicable for POST or PUT request.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Admin user that last modified the rule. This field is not applicable for POST or PUT request.
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies
	ECGroups []common.IDNameExtensions `json:"ecGroups,omitempty"`
}

type TrafficLogRulesCountQuery struct {
	PredefinedRuleCount bool
	RuleName            string
	RuleOrder           string
	RuleDescription     string
	Location            string
}

// TrafficLogRulesCountResponse defines the expected response structure
type TrafficLogRulesCountResponse struct {
	Count int `json:"count"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*ECTrafficLogRules, error) {
	var rule ECTrafficLogRules
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", logRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning forwarding dns rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetRulesByName(ctx context.Context, service *zscaler.Service, ruleName string) (*ECTrafficLogRules, error) {
	var rules []ECTrafficLogRules
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, logRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no dns rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rules *ECTrafficLogRules) (*ECTrafficLogRules, error) {
	resp, err := service.Client.CreateResource(ctx, logRulesEndpoint, *rules)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ECTrafficLogRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning dns rules from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *ECTrafficLogRules) (*ECTrafficLogRules, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", logRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*ECTrafficLogRules)
	service.Client.GetLogger().Printf("[DEBUG]returning forwarding dns rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", logRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ECTrafficLogRules, error) {
	var rules []ECTrafficLogRules
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, logRulesEndpoint, &rules)
	return rules, err
}

/*
// GetEcRDRCount retrieves the count of forwarding dns rules using optional filters
func GetEcRDRCount(ctx context.Context, service *zscaler.Service, params *TrafficLogRulesCountQuery) (*TrafficLogRulesCountResponse, error) {
	// Build query string
	query := url.Values{}
	if params != nil {
		query.Set("predefinedRuleCount", strconv.FormatBool(params.PredefinedRuleCount))
		if params.RuleName != "" {
			query.Set("ruleName", params.RuleName)
		}
		if params.RuleOrder != "" {
			query.Set("ruleOrder", params.RuleOrder)
		}
		if params.RuleDescription != "" {
			query.Set("ruleDescription", params.RuleDescription)
		}
		if params.Location != "" {
			query.Set("location", params.Location)
		}
	}

	endpoint := logRulesEndpoint + "/count"
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var result TrafficLogRulesCountResponse
	err := service.Client.ReadResource(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve forwarding dns rule count: %w", err)
	}
	return &result, nil
}
*/
