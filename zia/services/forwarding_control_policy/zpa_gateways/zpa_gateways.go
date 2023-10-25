package zpa_gateways

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	zpaGatewaysEndpoint = "/zpaGateways"
)

type ZPAGateways struct {
	// A unique identifier assigned to the ZPA gateway
	ID int `json:"id,omitempty"`

	// The name of the ZPA gateway
	Name string `json:"name,omitempty"`

	// Additional information about the rule
	Description string `json:"description,omitempty"`

	// The ZPA Server Group that is configured for Source IP Anchoring
	ZPAServerGroup []common.IDNameExtensions `json:"zpaServerGroup,omitempty"`

	// All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled
	ZPAAppSegments []common.IDNameExtensions `json:"zpaAppSegments,omitempty"`

	// The ID of the ZPA tenant where Source IP Anchoring is configured
	ZPATenantId int `json:"zpaTenantId,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Timestamp when the ZPA gateway was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)
	// Supported Values: "ZPA", "ECZPA"
	Type string `json:"type,omitempty"`
}

func (service *Service) Get(ruleID int) (*ZPAGateways, error) {
	var rule ZPAGateways
	err := service.Client.Read(fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning zpa gateway from Get: %d", rule.ID)
	return &rule, nil
}

func (service *Service) GetByName(ruleName string) (*ZPAGateways, error) {
	var rules []ZPAGateways
	err := common.ReadAllPages(service.Client, zpaGatewaysEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no zpa gateway found with name: %s", ruleName)
}

func (service *Service) Create(rule *ZPAGateways) (*ZPAGateways, error) {
	resp, err := service.Client.Create(zpaGatewaysEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ZPAGateways)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning zpa gateway from create: %d", createdRules.ID)
	return createdRules, nil
}

func (service *Service) Update(ruleID int, rules *ZPAGateways) (*ZPAGateways, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*ZPAGateways)
	service.Client.Logger.Printf("[DEBUG]returning zpa gateway from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
func (service *Service) GetAll() ([]ZPAGateways, error) {
	var rules []ZPAGateways
	err := common.ReadAllPages(service.Client, zpaGatewaysEndpoint, &rules)
	return rules, err
}
