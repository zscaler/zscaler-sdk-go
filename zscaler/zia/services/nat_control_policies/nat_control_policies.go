package nat_control_policies

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
	dnatRulesEndpoint = "/zia/api/v1/dnatRules"
)

type NatControlPolicies struct {
	AccessControl       string                    `json:"accessControl,omitempty"`
	ID                  int                       `json:"id,omitempty"`
	Name                string                    `json:"name,omitempty"`
	Order               int                       `json:"order,omitempty"`
	Rank                int                       `json:"rank,omitempty"`
	Description         string                    `json:"description,omitempty"`
	State               string                    `json:"state,omitempty"`
	RedirectFqdn        string                    `json:"redirectFqdn,omitempty"`
	RedirectIp          string                    `json:"redirectIp,omitempty"`
	RedirectPort        int                       `json:"redirectPort,omitempty"`
	LastModifiedTime    int                       `json:"lastModifiedTime,omitempty"`
	TrustedResolverRule bool                      `json:"trustedResolverRule,omitempty"`
	EnableFullLogging   bool                      `json:"enableFullLogging,omitempty"`
	Predefined          bool                      `json:"predefined,omitempty"`
	DefaultRule         bool                      `json:"defaultRule,omitempty"`
	DestAddresses       []string                  `json:"destAddresses,omitempty"`
	SrcIps              []string                  `json:"srcIps,omitempty"`
	DestCountries       []string                  `json:"destCountries,omitempty"`
	DestIpCategories    []string                  `json:"destIpCategories,omitempty"`
	ResCategories       []string                  `json:"resCategories,omitempty"`
	Locations           []common.IDNameExtensions `json:"locations,omitempty"`
	LocationGroups      []common.IDNameExtensions `json:"locationGroups,omitempty"`
	Groups              []common.IDNameExtensions `json:"groups,omitempty"`
	Departments         []common.IDNameExtensions `json:"departments,omitempty"`
	Users               []common.IDNameExtensions `json:"users,omitempty"`
	TimeWindows         []common.IDNameExtensions `json:"timeWindows,omitempty"`
	SrcIpGroups         []common.IDNameExtensions `json:"srcIpGroups,omitempty"`
	SrcIpv6Groups       []common.IDNameExtensions `json:"srcIpv6Groups,omitempty"`
	DestIpGroups        []common.IDNameExtensions `json:"destIpGroups,omitempty"`
	DestIpv6Groups      []common.IDNameExtensions `json:"destIpv6Groups,omitempty"`
	NwServices          []common.IDNameExtensions `json:"nwServices,omitempty"`
	NwServiceGroups     []common.IDNameExtensions `json:"nwServiceGroups,omitempty"`
	LastModifiedBy      *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	Devices             []common.IDNameExtensions `json:"devices,omitempty"`
	DeviceGroups        []common.IDNameExtensions `json:"deviceGroups,omitempty"`
	Labels              []common.IDNameExtensions `json:"labels,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*NatControlPolicies, error) {
	var rule NatControlPolicies
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dnatRulesEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning NAT Control rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*NatControlPolicies, error) {
	var rules []NatControlPolicies
	err := common.ReadAllPages(ctx, service.Client, dnatRulesEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no NAT Control rule rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *NatControlPolicies) (*NatControlPolicies, error) {
	resp, err := service.Client.Create(ctx, dnatRulesEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*NatControlPolicies)
	if !ok {
		return nil, errors.New("object returned from api was not a NAT Control rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning NAT Control rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *NatControlPolicies) (*NatControlPolicies, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", dnatRulesEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*NatControlPolicies)
	service.Client.GetLogger().Printf("[DEBUG]returning NAT Control rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", dnatRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]NatControlPolicies, error) {
	var rules []NatControlPolicies
	err := common.ReadAllPages(ctx, service.Client, dnatRulesEndpoint, &rules)
	return rules, err
}
