package forwardingrules

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	forwardingRulesEndpoint = "/ztw/api/v1/ecRules/ecRdr"
)

type ForwardingRules struct {
	// A unique identifier assigned to the forwarding rule
	ID int `json:"id,omitempty"`

	// The name of the forwarding rule
	Name string `json:"name,omitempty"`

	// Access permission available for the current user to the rule.
	AccessControl string `json:"accessControl,omitempty"`

	// Additional information about the forwarding rule
	Description string `json:"description,omitempty"`

	// The rule type selected from the available options
	// Supported Values: "FIREWALL", "DNS", "DNAT", "SNAT", "FORWARDING", "INTRUSION_PREVENTION", "EC_DNS", "EC_RDR", "EC_SELF", "DNS_RESPONSE"
	Type string `json:"type,omitempty"`

	// The order of execution for the forwarding rule order
	Order int `json:"order"`

	// Admin rank assigned to the forwarding rule
	Rank int `json:"rank"`

	// TThe forwarding method used in the rule, indicating whether the traffic is sent to ZIA, ZPA, directly to the destination (DIRECT), or dropped (DROP).
	// Supported Values: "INVALID", "DIRECT", "PROXYCHAIN", "ZIA", "ZPA", "ECZPA", "ECSELF", "DROP", "ENATDEDIP", "GEOIP"
	ForwardMethod string `json:"forwardMethod,omitempty"`

	// This parameter was deprecated and is no longer configurable.
	// Supported Values: "SMRULEF_ZPA_BROKERS_RULE", "SMRULEF_APPC_DYNAMIC_SRC_IPGROUP", "SMRULEF_EXCL_SRC_IP", "BALANCED_RULE", "BESTLINK_RULE"
	WanSelection string `json:"wanSelection,omitempty"`

	// Indicates whether the forwarding rule is enabled or disabled
	// Supported Values: DISABLED and ENABLED
	State string `json:"state,omitempty"`

	// Not applicable to Cloud & Branch Connector.
	BlockResponseCode string `json:"blockResponseCode,omitempty"`

	// Timestamp when the rule was last modified. This field is not applicable for POST or PUT request.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
	SrcIps []string `json:"srcIps,omitempty"`

	// List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses.
	//  If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.
	DestAddresses []string `json:"destAddresses,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	DestIpCategories []string `json:"destIpCategories,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	ResCategories []string `json:"resCategories,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	DestCountries []string `json:"destCountries,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	SourceCountries []string `json:"sourceCountries,omitempty"`

	// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
	NwApplications []string `json:"nwApplications,omitempty"`

	// Not applicable to Cloud & Branch Connector.
	SourceIpGroupExclusion bool `json:"sourceIpGroupExclusion,omitempty"`

	// The predefined ZPA Broker Rule generated by Zscaler (readonly: true)
	ZPABrokerRule bool `json:"zpaBrokerRule,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.CommonIDNameExternalID `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.CommonIDNameExternalID `json:"devices"`

	// Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
	Locations []common.CommonIDNameExternalID `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the forwarding rule applies
	LocationsGroups []common.CommonIDNameExternalID `json:"locationGroups,omitempty"`

	// Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies
	ECGroups []common.CommonIDNameExternalID `json:"ecGroups,omitempty"`

	// Name-ID pairs of the departments to which the forwarding rule applies. If not set, the rule applies to all departments.
	Departments []common.CommonIDNameExternalID `json:"departments,omitempty"`

	// Name-ID pairs of the user groups to which the forwarding rule applies. If not set, the rule applies to all groups.
	Groups []common.CommonIDNameExternalID `json:"groups,omitempty"`

	// Name-ID pairs of the users to which the forwarding rule applies. If not set, user criteria is ignored during policy enforcement.
	Users []common.CommonIDNameExternalID `json:"users,omitempty"`

	// Admin user that last modified the rule. This field is not applicable for POST or PUT request.
	LastModifiedBy *common.CommonIDNameExternalID `json:"lastModifiedBy,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	// Note: For organizations that have enabled IPv6, the srcIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	SrcIpGroups []common.CommonIDNameExternalID `json:"srcIpGroups,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	// Note: For organizations that have enabled IPv6, the srcIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	SrcIpv6Groups []common.CommonIDNameExternalID `json:"srcIpv6Groups,omitempty"`

	// User-defined destination IP address groups to which the rule is applied.
	// If not set, the rule is not restricted to a specific destination IP address group.
	DestIpGroups []common.CommonIDNameExternalID `json:"destIpGroups,omitempty"`

	// Destination IPv6 address groups for which the rule is applicable.
	// If not set, the rule is not restricted to a specific source IPv6 address group.
	DestIpv6Groups []common.IDNameExtensions `json:"destIpv6Groups,omitempty"`

	// User-defined network services to which the rule applies. If not set, the rule is not restricted to a specific network service.
	// Note: When the forwarding method is Proxy Chaining, only TCP-based network services are considered for policy match .
	NwServices []common.CommonIDNameExternalID `json:"nwServices,omitempty"`

	// User-defined network service group to which the rule applies.
	// If not set, the rule is not restricted to a specific network service group.
	NwServiceGroups []common.CommonIDNameExternalID `json:"nwServiceGroups,omitempty"`

	// Labels that are applicable to the rule.
	Labels []common.CommonIDNameExternalID `json:"labels,omitempty"`

	// User-defined network service application groups to which the rule applied.
	// If not set, the rule is not restricted to a specific network service application group.
	NwApplicationGroups []common.CommonIDNameExternalID `json:"nwApplicationGroups,omitempty"`

	AppServiceGroups []common.CommonIDNameExternalID `json:"appServiceGroups,omitempty"`

	// The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.
	ProxyGateway *common.CommonIDName `json:"proxyGateway,omitempty"`

	// The ZPA Server Group for which this rule is applicable.
	// Only the Server Groups that are associated with the selected Application Segments are allowed.
	// This field is applicable only for the ZPA forwarding method.
	ZPAGateway *common.CommonIDName `json:"zpaGateway,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.CommonZPAIDNameID `json:"zpaAppSegments"`

	// List of ZPA Application Segments for which this rule is applicable.
	// This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).
	ZPAApplicationSegments []ZPAApplicationSegments `json:"zpaApplicationSegments,omitempty"`

	// List of ZPA Application Segment Groups for which this rule is applicable.
	// This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).
	ZPAApplicationSegmentGroups []ZPAApplicationSegmentGroups `json:"zpaApplicationSegmentGroups,omitempty"`
}

type ZPAApplicationSegments struct {
	// A unique identifier assigned to the Application Segment
	ID int `json:"id,omitempty"`

	// The name of the Application Segment
	Name string `json:"name,omitempty"`

	// Additional information about the Application Segment
	Description string `json:"description,omitempty"`

	// ID of the ZPA tenant where the Application Segment is configured
	ZPAID int `json:"zpaId,omitempty"`

	// Indicates whether the ZPA Application Segment has been deleted
	Deleted bool `json:"deleted,omitempty"`
}

type ZPAApplicationSegmentGroups struct {
	// A unique identifier assigned to the Application Segment Group
	ID int `json:"id,omitempty"`

	// The name of the Application Segment Group
	Name string `json:"name,omitempty"`

	// ID of the ZPA tenant where the Application Segment is configured
	ZPAID int `json:"zpaId,omitempty"`

	// Indicates whether the ZPA Application Segment has been deleted
	Deleted bool `json:"deleted,omitempty"`

	// The number of ZPA Application Segments in the group
	ZPAAppSegmentsCount int `json:"zpaAppSegmentsCount,omitempty"`
}

type ForwardingRulesCountQuery struct {
	PredefinedRuleCount bool
	RuleName            string
	RuleOrder           string
	RuleDescription     string
	RuleForwardMethod   string
	Location            string
}

// ForwardingRulesCountResponse defines the expected response structure
type ForwardingRulesCountResponse struct {
	Count int `json:"count"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*ForwardingRules, error) {
	var rule ForwardingRules
	err := service.Client.ReadResource(ctx, fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning forwarding rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetRulesByName(ctx context.Context, service *zscaler.Service, ruleName string) (*ForwardingRules, error) {
	var rules []ForwardingRules
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, forwardingRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rules *ForwardingRules) (*ForwardingRules, error) {
	resp, err := service.Client.CreateResource(ctx, forwardingRulesEndpoint, *rules)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ForwardingRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rules from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *ForwardingRules) (*ForwardingRules, error) {
	resp, err := service.Client.UpdateWithPutResource(ctx, fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*ForwardingRules)
	service.Client.GetLogger().Printf("[DEBUG]returning forwarding rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.DeleteResource(ctx, fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ForwardingRules, error) {
	var rules []ForwardingRules
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(ctx, service.Client, forwardingRulesEndpoint, &rules)
	return rules, err
}

// GetEcRDRCount retrieves the count of forwarding rules using optional filters
func GetEcRDRCount(ctx context.Context, service *zscaler.Service, params *ForwardingRulesCountQuery) (*ForwardingRulesCountResponse, error) {
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
		if params.RuleForwardMethod != "" {
			query.Set("ruleForwardMethod", params.RuleForwardMethod)
		}
		if params.Location != "" {
			query.Set("location", params.Location)
		}
	}

	endpoint := forwardingRulesEndpoint + "/count"
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	var result ForwardingRulesCountResponse
	err := service.Client.ReadResource(ctx, endpoint, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve forwarding rule count: %w", err)
	}
	return &result, nil
}
