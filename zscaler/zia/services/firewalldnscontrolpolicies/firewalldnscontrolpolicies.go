package firewalldnscontrolpolicies

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
	firewallDnsRulesEndpoint = "/zia/api/v1/firewallDnsRules"
)

type FirewallDNSRules struct {
	// Unique identifier for the Firewall Filtering policy rule
	ID int `json:"id,omitempty"`

	// Name of the Firewall Filtering policy rule
	Name string `json:"name,omitempty"`

	// Rule order number of the Firewall Filtering policy rule
	Order int `json:"order,omitempty"`

	// Admin rank of the Firewall Filtering policy rule
	Rank int `json:"rank,omitempty"`

	// The admin’s access privilege to this rule based on the assigned role
	AccessControl string `json:"accessControl,omitempty"`

	// The action the Firewall Filtering policy rule takes when packets match the rule
	// Supported Values: "ALLOW", "BLOCK", "REDIR_REQ", "REDIR_RES", "REDIR_ZPA", "REDIR_REQ_DOH", "REDIR_REQ_KEEP_SENDER", "REDIR_REQ_TCP", "REDIR_REQ_UDP","BLOCK_WITH_RESPONSE"
	Action string `json:"action,omitempty"`

	// Determines whether the Firewall Filtering policy rule is enabled or disabled
	State string `json:"state,omitempty"`

	// Additional information about the rule
	Description string `json:"description,omitempty"`

	// The IP address to which the traffic will be redirected to when the DNAT rule is triggered. If not set, no redirection is done to specific IP addresses.
	RedirectIP string `json:"redirectIp,omitempty"`

	// Specifies the DNS response code to be sent to the client when the action is configured to block and send response code
	BlockResponseCode string `json:"blockResponseCode,omitempty"`

	// Timestamp when the rule was last modified. Ignored if the request is POST or PUT. For GET, ignored if or the rule is current version.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// The admin who last modified the rule
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
	SrcIps []string `json:"srcIps,omitempty"`

	// Destination IP addresses or FQDNs to which the rule applies. If not set, the rule is not restricted to a
	// specific destination IP address. Each IP entry can be a single IP address, CIDR (e.g., 10.10.33.0/24), or an IP range (e.g., 10.10.33.1-10.10.33.10).
	DestAddresses []string `json:"destAddresses,omitempty"`

	// IP address categories of destination for which the DNAT rule is applicable. If not set, the rule is not restricted to specific destination IP categories.
	DestIpCategories []string `json:"destIpCategories,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	DestCountries []string `json:"destCountries,omitempty"`

	// The countries of origin of traffic for which the rule is applicable. If not set, the rule is not restricted to specific source countries.
	SourceCountries []string `json:"sourceCountries,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	ResCategories []string `json:"resCategories,omitempty"`

	// DNS tunnels and network applications to which the rule applies
	Applications []string `json:"applications,omitempty"`

	// DNS request types to which the rule applies
	DNSRuleRequestTypes []string `json:"dnsRuleRequestTypes,omitempty"`

	// The protocols to which the rules applies
	// Supported Values: "ANY_RULE", "SMRULEF_CASCADING_ALLOWED", "TCP_RULE", "UDP_RULE", "DOHTTPS_RULE"
	Protocols []string `json:"protocols,omitempty"`

	// If set to true, the default rule is applied
	DefaultRule bool `json:"defaultRule,omitempty"`

	// A Boolean value that indicates whether packet capture (PCAP) is enabled or not
	CapturePCAP bool `json:"capturePCAP"`

	// A Boolean field that indicates that the rule is predefined by using a true value
	Predefined bool `json:"predefined,omitempty"`

	// DNS application groups to which the rule applies
	ApplicationGroups []common.IDNameExtensions `json:"applicationGroups,omitempty"`

	// The DNS gateway used to redirect traffic, specified when the rule action is to redirect DNS request to an external DNS service.
	DNSGateway *common.IDName `json:"dnsGateway,omitempty"`

	// The ZPA IP pool specified when the rule action is to resolve domain names of ZPA applications to an ephemeral IP address from a preconfigured IP pool.
	ZPAIPGroup *common.IDName `json:"zpaIpGroup"`

	// EDNS ECS object which resolves DNS request
	EDNSEcsObject *common.IDName `json:"ednsEcsObject,omitempty"`

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

	// Labels that are applicable to the rule.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// User-defined destination IP address groups on which the rule is applied. If not set, the rule is not restricted to a specific destination IP address group.
	// Note: For organizations that have enabled IPv6, the destIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	DestIpGroups []common.IDNameExtensions `json:"destIpGroups,omitempty"`

	// Destination IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
	DestIpv6Groups []common.IDNameExtensions `json:"destIpv6Groups,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	SrcIpGroups []common.IDNameExtensions `json:"srcIpGroups,omitempty"`

	// Source IPv6 address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IPv6 address group.
	SrcIpv6Groups []common.IDNameExtensions `json:"srcIpv6Groups,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*FirewallDNSRules, error) {
	var rule FirewallDNSRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", firewallDnsRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning firewall dns rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*FirewallDNSRules, error) {
	var rules []FirewallDNSRules
	err := common.ReadAllPages(ctx, service.Client, firewallDnsRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no firewall dns rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *FirewallDNSRules) (*FirewallDNSRules, error) {
	resp, err := service.Client.Create(ctx, firewallDnsRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*FirewallDNSRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning firewall dns rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *FirewallDNSRules) (*FirewallDNSRules, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", firewallDnsRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*FirewallDNSRules)
	service.Client.GetLogger().Printf("[DEBUG]returning firewall dns rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

/*
	func Create(ctx context.Context, service *zscaler.Service, rule *FirewallDNSRules) (*FirewallDNSRules, error) {
		//Validate the rule before creating
		// if err := validateFirewallDNSRules(rule); err != nil {
		// 	return nil, fmt.Errorf("validation failed for FirewallDNSRules: %w", err)
		// }

		// Proceed with creating the rule
		resp, err := service.Client.Create(ctx, firewallDnsRulesEndpoint, *rule)
		if err != nil {
			return nil, err
		}

		createdRules, ok := resp.(*FirewallDNSRules)
		if !ok {
			return nil, errors.New("object returned from api was not a rule Pointer")
		}

		service.Client.GetLogger().Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
		return createdRules, nil
	}
*/

/*
func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *FirewallDNSRules) (*FirewallDNSRules, error) {
	// Validate the rule before updating
	// if err := validateFirewallDNSRules(rules); err != nil {
	// 	return nil, fmt.Errorf("validation failed for FirewallDNSRules: %w", err)
	// }

	// Proceed with updating the rule
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", firewallDnsRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}

	updatedRules, ok := resp.(*FirewallDNSRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning firewall ips rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}
*/

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", firewallDnsRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]FirewallDNSRules, error) {
	var rules []FirewallDNSRules
	err := common.ReadAllPages(ctx, service.Client, firewallDnsRulesEndpoint, &rules)
	return rules, err
}

/*
func validateFirewallDNSRules(rule *FirewallDNSRules) error {
	switch rule.Action {
	case "REDIR_REQ_KEEP_SENDER":
		// Validate DNSGateway is not nil and contains a valid ID or Name
		if rule.DNSGateway == nil || rule.DNSGateway.ID == 0 || rule.DNSGateway.Name == "" {
			return errors.New("dnsGateway must be provided with a valid ID and Name when action is REDIR_REQ_KEEP_SENDER")
		}
		// Validate Protocols is not empty
		if len(rule.Protocols) == 0 {
			return errors.New("protocols must be provided when action is REDIR_REQ_KEEP_SENDER")
		}
	case "REDIR_REQ_DOH", "REDIR_REQ_TCP", "REDIR_REQ_UDP":
		// Validate DNSGateway is not nil and contains a valid ID or Name
		if rule.DNSGateway == nil || rule.DNSGateway.ID == 0 || rule.DNSGateway.Name == "" {
			return fmt.Errorf("dnsGateway must be provided with a valid ID and Name when action is %s", rule.Action)
		}
	case "REDIR_ZPA":
		// Validate ZPAIPGroup is not nil and contains a valid ID
		if rule.ZPAIPGroup == nil || rule.ZPAIPGroup.ID == 0 {
			return errors.New("zpaIpGroup must be provided with a valid ID when action is REDIR_ZPA")
		}
	case "REDIR_RES":
		// Validate RedirectIP is not empty
		if rule.RedirectIP == "" {
			return errors.New("redirectIp must be provided when action is REDIR_RES")
		}
	}
	return nil
}
*/
