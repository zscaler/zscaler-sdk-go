package filteringrules

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
)

const (
	firewallRulesEndpoint = "/zia/api/v1/firewallFilteringRules"
)

type FirewallFilteringRules struct {
	// Unique identifier for the Firewall Filtering policy rule
	ID int `json:"id,omitempty"`

	// Name of the Firewall Filtering policy rule
	Name string `json:"name,omitempty"`

	// Rule order number of the Firewall Filtering policy rule
	Order int `json:"order"`

	// Admin rank of the Firewall Filtering policy rule
	Rank int `json:"rank"`

	// The admin’s access privilege to this rule based on the assigned role
	AccessControl string `json:"accessControl,omitempty"`

	// A Boolean value that indicates whether full logging is enabled. A true value indicates that full logging is enabled, whereas a false value indicates that aggregate logging is enabled.
	EnableFullLogging bool `json:"enableFullLogging"`

	// The action the Firewall Filtering policy rule takes when packets match the rule
	Action string `json:"action,omitempty"`

	// Determines whether the Firewall Filtering policy rule is enabled or disabled
	State string `json:"state,omitempty"`

	// Additional information about the rule
	Description string `json:"description,omitempty"`

	// Timestamp when the rule was last modified. Ignored if the request is POST or PUT. For GET, ignored if or the rule is current version.
	LastModifiedTime int                      `json:"lastModifiedTime,omitempty"`
	LastModifiedBy   *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
	SrcIps []string `json:"srcIps,omitempty"`

	// List of destination IP addresses for which the rule is applicable. CIDR notation can be used for destination IP addresses. If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups or destIpCategories.
	DestAddresses []string `json:"destAddresses,omitempty"`

	// IP address categories of destination for which the DNAT rule is applicable. If not set, the rule is not restricted to specific destination IP categories.
	DestIpCategories []string `json:"destIpCategories,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	DestCountries []string `json:"destCountries,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	SourceCountries []string `json:"sourceCountries,omitempty"`

	// Indicates whether the countries specified in the sourceCountries field are included or excluded from the rule.
	// A true value denotes that the specified source countries are excluded from the rule.
	// A false value denotes that the rule is applied to the source countries if there is a match.
	ExcludeSrcCountries bool `json:"excludeSrcCountries,omitempty"`

	// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
	NwApplications []string `json:"nwApplications,omitempty"`

	// If set to true, the default rule is applied
	DefaultRule bool `json:"defaultRule"`

	// If set to true, a predefined rule is applied
	Predefined bool `json:"predefined"`

	// The locations to which the Firewall Filtering policy rule applies
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// The location groups to which the Firewall Filtering policy rule applies
	LocationsGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// The departments to which the Firewall Filtering policy rule applies
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// The groups to which the Firewall Filtering policy rule applies
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// The users to which the Firewall Filtering policy rule applies
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// The time interval in which the Firewall Filtering policy rule applies
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// User-defined network service application group on which the rule is applied. If not set, the rule is not restricted to a specific network service application group.
	NwApplicationGroups []common.IDNameExtensions `json:"nwApplicationGroups,omitempty"`

	// Application services on which this rule is applied
	AppServices []common.IDNameExtensions `json:"appServices,omitempty"`

	// Application service groups on which this rule is applied
	AppServiceGroups []common.IDNameExtensions `json:"appServiceGroups,omitempty"`

	// Labels that are applicable to the rule.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// User-defined destination IP address groups on which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.
	// Note: For organizations that have enabled IPv6, the destIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	DestIpGroups []common.IDNameExtensions `json:"destIpGroups,omitempty"`

	// User-defined network services on which the rule is applied. If not set, the rule is not restricted to a specific network service.
	NwServices []common.IDNameExtensions `json:"nwServices,omitempty"`

	// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
	NwServiceGroups []common.IDNameExtensions `json:"nwServiceGroups,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	SrcIpGroups []common.IDNameExtensions `json:"srcIpGroups,omitempty"`

	// List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.
	DeviceTrustLevels []string `json:"deviceTrustLevels,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`

	// The list of preconfigured workload groups to which the policy must be applied.
	WorkloadGroups []common.IDName `json:"workloadGroups,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.ZPAAppSegments `json:"zpaAppSegments"`
}

// GetAllFilterOptions represents optional filter parameters for GetAll
type GetAllFilterOptions struct {
	PredefinedRuleCount bool
	RuleName            string
	RuleLabel           string
	RuleLabelId         int
	RuleOrder           string
	RuleDescription     string
	RuleAction          string
	Location            string
	Department          string
	Group               string
	User                string
	Device              string
	DeviceGroup         string
	DeviceTrustLevel    string
	SrcIps              string
	DestAddresses       string
	SrcIpGroups         string
	DestIpGroups        string
	NwApplication       string
	NwServices          string
	DestIpCategories    string
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*FirewallFilteringRules, error) {
	var rule FirewallFilteringRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning firewall rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*FirewallFilteringRules, error) {
	// Use GetAll with RuleName filter to leverage API filtering and pagination
	opts := &GetAllFilterOptions{
		RuleName: ruleName,
	}
	rules, err := GetAll(ctx, service, opts)
	if err != nil {
		return nil, err
	}
	// API may do partial matching, so verify exact match (case-insensitive)
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no firewall rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.Create(ctx, firewallRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*FirewallFilteringRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*FirewallFilteringRules)
	service.Client.GetLogger().Printf("[DEBUG]returning firewall rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// GetAll retrieves all firewall filtering rules with optional filters.
// This endpoint supports a maximum page size of 5000.
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]FirewallFilteringRules, error) {
	var rules []FirewallFilteringRules
	endpoint := firewallRulesEndpoint

	// Build query parameters from filter options
	queryParams := url.Values{}
	if opts != nil {
		if opts.RuleName != "" {
			queryParams.Add("ruleName", opts.RuleName)
		}
		if opts.RuleLabel != "" {
			queryParams.Add("ruleLabel", opts.RuleLabel)
		}
		if opts.RuleLabelId != 0 {
			queryParams.Add("ruleLabelId", strconv.Itoa(opts.RuleLabelId))
		}
		if opts.RuleOrder != "" {
			queryParams.Add("ruleOrder", opts.RuleOrder)
		}
		if opts.RuleDescription != "" {
			queryParams.Add("ruleDescription", opts.RuleDescription)
		}
		if opts.RuleAction != "" {
			queryParams.Add("ruleAction", opts.RuleAction)
		}
		if opts.Location != "" {
			queryParams.Add("location", opts.Location)
		}
		if opts.Department != "" {
			queryParams.Add("department", opts.Department)
		}
		if opts.Group != "" {
			queryParams.Add("group", opts.Group)
		}
		if opts.User != "" {
			queryParams.Add("user", opts.User)
		}
		if opts.Device != "" {
			queryParams.Add("device", opts.Device)
		}
		if opts.DeviceGroup != "" {
			queryParams.Add("deviceGroup", opts.DeviceGroup)
		}
		if opts.DeviceTrustLevel != "" {
			queryParams.Add("deviceTrustLevel", opts.DeviceTrustLevel)
		}
		if opts.SrcIps != "" {
			queryParams.Add("srcIps", opts.SrcIps)
		}
		if opts.DestAddresses != "" {
			queryParams.Add("destAddresses", opts.DestAddresses)
		}
		if opts.SrcIpGroups != "" {
			queryParams.Add("srcIpGroups", opts.SrcIpGroups)
		}
		if opts.DestIpGroups != "" {
			queryParams.Add("destIpGroups", opts.DestIpGroups)
		}
		if opts.NwApplication != "" {
			queryParams.Add("nwApplication", opts.NwApplication)
		}
		if opts.NwServices != "" {
			queryParams.Add("nwServices", opts.NwServices)
		}
		if opts.DestIpCategories != "" {
			queryParams.Add("destIpCategories", opts.DestIpCategories)
		}
	}

	// Build endpoint with query parameters
	baseQuery := queryParams.Encode()
	if baseQuery != "" {
		endpoint += "?" + baseQuery
	}

	// Use pageSize=5000 as this endpoint supports it
	err := common.ReadAllPages(ctx, service.Client, endpoint, &rules, 5000)
	return rules, err
}

// GetFirewallFilteringRuleCount retrieves the count of firewall filtering rules using optional filters.
// The API returns a simple integer count.
func GetFirewallFilteringRuleCount(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) (int, error) {
	// Build query string
	query := url.Values{}
	if opts != nil {
		query.Set("predefinedRuleCount", strconv.FormatBool(opts.PredefinedRuleCount))
		if opts.RuleName != "" {
			query.Set("ruleName", opts.RuleName)
		}
		if opts.RuleLabel != "" {
			query.Set("ruleLabel", opts.RuleLabel)
		}
		if opts.RuleLabelId != 0 {
			query.Set("ruleLabelId", strconv.Itoa(opts.RuleLabelId))
		}
		if opts.RuleOrder != "" {
			query.Set("ruleOrder", opts.RuleOrder)
		}
		if opts.RuleDescription != "" {
			query.Set("ruleDescription", opts.RuleDescription)
		}
		if opts.RuleAction != "" {
			query.Set("ruleAction", opts.RuleAction)
		}
		if opts.Location != "" {
			query.Set("location", opts.Location)
		}
		if opts.Department != "" {
			query.Set("department", opts.Department)
		}
		if opts.Group != "" {
			query.Set("group", opts.Group)
		}
		if opts.User != "" {
			query.Set("user", opts.User)
		}
		if opts.Device != "" {
			query.Set("device", opts.Device)
		}
		if opts.DeviceGroup != "" {
			query.Set("deviceGroup", opts.DeviceGroup)
		}
		if opts.DeviceTrustLevel != "" {
			query.Set("deviceTrustLevel", opts.DeviceTrustLevel)
		}
		if opts.SrcIps != "" {
			query.Set("srcIps", opts.SrcIps)
		}
		if opts.DestAddresses != "" {
			query.Set("destAddresses", opts.DestAddresses)
		}
		if opts.SrcIpGroups != "" {
			query.Set("srcIpGroups", opts.SrcIpGroups)
		}
		if opts.DestIpGroups != "" {
			query.Set("destIpGroups", opts.DestIpGroups)
		}
		if opts.NwApplication != "" {
			query.Set("nwApplication", opts.NwApplication)
		}
		if opts.NwServices != "" {
			query.Set("nwServices", opts.NwServices)
		}
		if opts.DestIpCategories != "" {
			query.Set("destIpCategories", opts.DestIpCategories)
		}
	}

	endpoint := firewallRulesEndpoint + "/count"
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var count int
	err := service.Client.Read(ctx, endpoint, &count)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve firewall filtering rule count: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning firewall filtering rule count: %d", count)
	return count, nil
}
