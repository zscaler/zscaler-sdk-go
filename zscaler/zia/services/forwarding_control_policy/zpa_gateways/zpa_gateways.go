package zpa_gateways

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
	zpaGatewaysEndpoint = "/zia/api/v1/zpaGateways"
)

type ZPAGateways struct {
	// A unique identifier assigned to the ZPA gateway
	ID int `json:"id"`

	// The name of the ZPA gateway
	Name string `json:"name,omitempty"`

	// Additional details about the ZPA gateway
	Description string `json:"description,omitempty"`

	// The ZPA Server Group that is configured for Source IP Anchoring
	ZPAServerGroup ZPAServerGroup `json:"zpaServerGroup,omitempty"`

	// All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled
	ZPAAppSegments []ZPAAppSegments `json:"zpaAppSegments,omitempty"`

	// The ID of the ZPA tenant where Source IP Anchoring is configured
	ZPATenantId int `json:"zpaTenantId,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Timestamp when the ZPA gateway was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Indicates whether the ZPA gateway is configured for Zscaler Internet Access (using option ZPA) or Zscaler Cloud Connector (using option ECZPA)
	// Supported Values: "ZPA", "ECZPA"
	Type string `json:"type"`
}

// The ZPA Server Group that is configured for Source IP Anchoring
type ZPAServerGroup struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The name of the Application Segment
	Name string `json:"name,omitempty"`

	// An external identifier used for an entity that is managed outside of ZIA.
	// Examples include zpaServerGroup and zpaAppSegments.
	// This field is not applicable to ZIA-managed entities.
	ExternalID string `json:"externalId,omitempty"`

	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// All the Application Segments that are associated with the selected ZPA Server Group for which Source IP Anchoring is enabled
type ZPAAppSegments struct {
	// Identifier that uniquely identifies an entity
	ID int `json:"id,omitempty"`

	// The name of the Application Segment
	Name string `json:"name,omitempty"`

	// An external identifier used for an entity that is managed outside of ZIA.
	// Examples include zpaServerGroup and zpaAppSegments.
	// This field is not applicable to ZIA-managed entities.
	ExternalID string `json:"externalId,omitempty"`

	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, gatewayID int) (*ZPAGateways, error) {
	var rule ZPAGateways
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, gatewayID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning zpa gateway from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*ZPAGateways, error) {
	var rules []ZPAGateways
	err := common.ReadAllPages(ctx, service.Client, zpaGatewaysEndpoint, &rules)
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

func Create(ctx context.Context, service *zscaler.Service, rule *ZPAGateways) (*ZPAGateways, error) {
	resp, err := service.Client.Create(ctx, zpaGatewaysEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*ZPAGateways)
	if !ok {
		return nil, errors.New("object returned from api was not a rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning zpa gateway from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, gatewayID int, rules *ZPAGateways) (*ZPAGateways, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, gatewayID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*ZPAGateways)
	service.Client.GetLogger().Printf("[DEBUG]returning zpa gateway from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, gatewayID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", zpaGatewaysEndpoint, gatewayID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ZPAGateways, error) {
	var rules []ZPAGateways
	err := common.ReadAllPages(ctx, service.Client, zpaGatewaysEndpoint, &rules)
	return rules, err
}
