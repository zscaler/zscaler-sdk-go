package filteringrules

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	firewallRulesEndpoint = "/firewallFilteringRules"
)

type FirewallFilteringRules struct {
	// Unique identifier for the Firewall Filtering policy rule
	ID int `json:"id,omitempty"`

	// Name of the Firewall Filtering policy rule
	Name string `json:"name,omitempty"`

	// Rule order number of the Firewall Filtering policy rule
	Order int `json:"order"`

	// Admin rank of the Firewall Filtering policy rule
	Rank              int    `json:"rank"`
	AccessControl     string `json:"accessControl,omitempty"`
	EnableFullLogging bool   `json:"enableFullLogging"`

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
	WorkloadGroups []common.IDNameExtensions `json:"workloadGroups,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.ZPAAppSegments `json:"zpaAppSegments"`
}

func (service *Service) Get(ruleID int) (*FirewallFilteringRules, error) {
	var rule FirewallFilteringRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning firewall rule from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetByName(ruleName string) (*FirewallFilteringRules, error) {
	var rules []FirewallFilteringRules
	err := common.ReadAllPages(service.Client, firewallRulesEndpoint, &rules)
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

func (service *Service) Create(rule *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.Create(firewallRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*FirewallFilteringRules)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func (service *Service) Update(ruleID int, rules *FirewallFilteringRules) (*FirewallFilteringRules, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*FirewallFilteringRules)
	service.Client.Logger.Printf("[DEBUG]returning firewall rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", firewallRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]FirewallFilteringRules, error) {
	var rules []FirewallFilteringRules
	err := common.ReadAllPages(service.Client, firewallRulesEndpoint, &rules)
	return rules, err
}
