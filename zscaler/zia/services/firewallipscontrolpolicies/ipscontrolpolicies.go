package firewallipscontrolpolicies

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	firewallIpsRulesEndpoint = "/zia/api/v1/firewallIpsRules"
)

type FirewallIPSRules struct {
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

	// The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
	SourceCountries []string `json:"sourceCountries,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	ResCategories []string `json:"resCategories,omitempty"`

	// If set to true, the default rule is applied
	DefaultRule bool `json:"defaultRule"`

	// A Boolean value that indicates whether packet capture (PCAP) is enabled or not
	CapturePCAP bool `json:"capturePCAP"`

	// A Boolean field that indicates that the rule is predefined by using a true value
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

	// Destination IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
	DestIpv6Groups []common.IDNameExtensions `json:"destIpv6Groups,omitempty"`

	// User-defined network services on which the rule is applied. If not set, the rule is not restricted to a specific network service.
	NwServices []common.IDNameExtensions `json:"nwServices,omitempty"`

	// User-defined network service applications on which the rule is applied. If not set, the rule is not restricted to a specific network service application.
	NwServiceGroups []common.IDNameExtensions `json:"nwServiceGroups,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	SrcIpGroups []common.IDNameExtensions `json:"srcIpGroups,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`

	// Advanced threat categories to which the rule applies
	ThreatCategories []common.IDNameExtensions `json:"threatCategories,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.ZPAAppSegments `json:"zpaAppSegments"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*FirewallIPSRules, error) {
	var rule FirewallIPSRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", firewallIpsRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning firewall ips rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*FirewallIPSRules, error) {
	var rules []FirewallIPSRules
	err := common.ReadAllPages(ctx, service.Client, firewallIpsRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no firewall ips rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *FirewallIPSRules) (*FirewallIPSRules, error) {
	resp, err := service.Client.Create(ctx, firewallIpsRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*FirewallIPSRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *FirewallIPSRules) (*FirewallIPSRules, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", firewallIpsRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*FirewallIPSRules)
	service.Client.GetLogger().Printf("[DEBUG]returning firewall ips rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", firewallIpsRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]FirewallIPSRules, error) {
	var rules []FirewallIPSRules
	err := common.ReadAllPages(ctx, service.Client, firewallIpsRulesEndpoint, &rules)
	return rules, err
}
