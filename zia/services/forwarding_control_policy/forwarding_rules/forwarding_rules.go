package forwarding_rules

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	forwardingRulesEndpoint = "/forwardingRules"
)

type ForwardingRules struct {
	// A unique identifier assigned to the forwarding rule
	ID int `json:"id,omitempty"`

	// The name of the forwarding rule
	Name string `json:"name,omitempty"`

	// Additional information about the forwarding rule
	Description string `json:"description,omitempty"`

	// The rule type selected from the available options
	// Supported Values: "FIREWALL", "DNS", "DNAT", "SNAT", "FORWARDING", "INTRUSION_PREVENTION", "EC_DNS", "EC_RDR", "EC_SELF", "DNS_RESPONSE"
	Type string `json:"type,omitempty"`

	// The order of execution for the forwarding rule order
	Order int `json:"order"`

	// Admin rank assigned to the forwarding rule
	Rank int `json:"rank"`

	// Name-ID pairs of the locations to which the forwarding rule applies. If not set, the rule is applied to all locations.
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the forwarding rule applies
	LocationsGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// Name-ID pairs of the Zscaler Cloud Connector groups to which the forwarding rule applies
	ECGroups []common.IDNameExtensions `json:"ecGroups,omitempty"`

	// Name-ID pairs of the departments to which the forwarding rule applies. If not set, the rule applies to all departments.
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of the user groups to which the forwarding rule applies. If not set, the rule applies to all groups.
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// Name-ID pairs of the users to which the forwarding rule applies. If not set, user criteria is ignored during policy enforcement.
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// The type of traffic forwarding method selected from the available options
	// Supported Values: "INVALID", "DIRECT", "PROXYCHAIN", "ZIA", "ZPA", "ECZPA", "ECSELF", "DROP"
	ForwardMethod string `json:"forwardMethod,omitempty"`

	// Indicates whether the forwarding rule is enabled or disabled
	// Supported Values: DISABLED and ENABLED
	State string `json:"state,omitempty"`

	// Timestamp when the rule was last modified. This field is not applicable for POST or PUT request.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Admin user that last modified the rule. This field is not applicable for POST or PUT request.
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// User-defined source IP addresses for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address.
	SrcIps []string `json:"srcIps,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	// Note: For organizations that have enabled IPv6, the srcIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	SrcIpGroups []common.IDNameExtensions `json:"srcIpGroups,omitempty"`

	// Source IP address groups for which the rule is applicable. If not set, the rule is not restricted to a specific source IP address group.
	// Note: For organizations that have enabled IPv6, the srcIpv6Groups field lists the IPv6 source address groups for which the rule is applicable.
	SrcIpv6Groups []common.IDNameExtensions `json:"srcIpv6Groups,omitempty"`

	// List of destination IP addresses or FQDNs for which the rule is applicable. CIDR notation can be used for destination IP addresses.
	//  If not set, the rule is not restricted to a specific destination addresses unless specified by destCountries, destIpGroups, or destIpCategories.
	DestAddresses []string `json:"destAddresses,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	DestIpCategories []string `json:"destIpCategories,omitempty"`

	// List of destination IP categories to which the rule applies. If not set, the rule is not restricted to specific destination IP categories.
	ResCategories []string `json:"resCategories,omitempty"`

	// Destination countries for which the rule is applicable. If not set, the rule is not restricted to specific destination countries.
	DestCountries []string `json:"destCountries,omitempty"`

	// User-defined destination IP address groups to which the rule is applied.
	// If not set, the rule is not restricted to a specific destination IP address group.
	DestIpGroups []common.IDNameExtensions `json:"destIpGroups,omitempty"`

	// Destination IPv6 address groups for which the rule is applicable.
	// If not set, the rule is not restricted to a specific source IPv6 address group.
	DestIpv6Groups []common.IDNameExtensions `json:"destIpv6Groups,omitempty"`

	// User-defined network services to which the rule applies. If not set, the rule is not restricted to a specific network service.
	// Note: When the forwarding method is Proxy Chaining, only TCP-based network services are considered for policy match .
	NwServices []common.IDNameExtensions `json:"nwServices,omitempty"`

	// User-defined network service group to which the rule applies.
	// If not set, the rule is not restricted to a specific network service group.
	NwServiceGroups []common.IDNameExtensions `json:"nwServiceGroups,omitempty"`

	// Labels that are applicable to the rule.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// User-defined network service application groups to which the rule applied.
	// If not set, the rule is not restricted to a specific network service application group.
	NwApplicationGroups []common.IDNameExtensions `json:"nwApplicationGroups,omitempty"`

	AppServiceGroups []common.IDNameExtensions `json:"appServiceGroups,omitempty"`

	// The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.
	ProxyGateway *common.IDName `json:"proxyGateway,omitempty"`

	// The ZPA Server Group for which this rule is applicable.
	// Only the Server Groups that are associated with the selected Application Segments are allowed.
	// This field is applicable only for the ZPA forwarding method.
	ZPAGateway *common.IDName `json:"zpaGateway,omitempty"`

	// The proxy gateway for which the rule is applicable. This field is applicable only for the Proxy Chaining forwarding method.
	ZPAAppSegments []ZPAAppSegments `json:"zpaAppSegments"`

	// List of ZPA Application Segments for which this rule is applicable.
	// This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).
	ZPAApplicationSegments []ZPAApplicationSegments `json:"zpaApplicationSegments,omitempty"`

	// List of ZPA Application Segment Groups for which this rule is applicable.
	// This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector).
	ZPAApplicationSegmentGroups []ZPAApplicationSegmentGroups `json:"zpaApplicationSegmentGroups,omitempty"`

	// The predefined ZPA Broker Rule generated by Zscaler (readonly: true)
	ZPABrokerRule bool `json:"zpaBrokerRule,omitempty"`
}

type ZPAAppSegments struct {
	// A unique identifier assigned to the Application Segment
	ID int `json:"id"`

	// The name of the Application Segment
	Name string `json:"name,omitempty"`

	// Indicates the external ID. Applicable only when this reference is of an external entity.
	ExternalID string `json:"externalId"`

	// ID of the ZPA tenant where the Application Segment is configured.
	ZPATenantId string `json:"zpaTenantId,omitempty"`
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

func (service *Service) Get(ruleID int) (*ForwardingRules, error) {
	var rule ForwardingRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning forwarding rule from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetByName(ruleName string) (*ForwardingRules, error) {
	var rules []ForwardingRules
	err := common.ReadAllPages(service.Client, forwardingRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no forwarding rule found with name: %s", ruleName)
}

func (service *Service) Create(rule *ForwardingRules) (*ForwardingRules, error) {
	resp, err := service.Client.Create(forwardingRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ForwardingRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func (service *Service) Update(ruleID int, rules *ForwardingRules) (*ForwardingRules, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*ForwardingRules)
	service.Client.Logger.Printf("[DEBUG]returning forwarding rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", forwardingRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]ForwardingRules, error) {
	var rules []ForwardingRules
	err := common.ReadAllPages(service.Client, forwardingRulesEndpoint, &rules)
	return rules, err
}
