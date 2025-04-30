package bandwidth_control_rules

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
	bandwidthControlEndpoint = "/zia/api/v1/bandwidthControlRules"
)

type BandwidthControlRules struct {
	ID                int                       `json:"id,omitempty"`
	Name              string                    `json:"name,omitempty"`
	Order             int                       `json:"order,omitempty"`
	State             string                    `json:"state,omitempty"`
	Description       string                    `json:"description,omitempty"`
	MaxBandwidth      int                       `json:"maxBandwidth,omitempty"`
	MinBandwidth      int                       `json:"minBandwidth,omitempty"`
	Rank              int                       `json:"rank,omitempty"`
	LastModifiedTime  int                       `json:"lastModifiedTime,omitempty"`
	AccessControl     string                    `json:"accessControl,omitempty"`
	DefaultRule       bool                      `json:"defaultRule,omitempty"`
	Protocols         []string                  `json:"protocols,omitempty"`
	DeviceTrustLevels []string                  `json:"deviceTrustLevels,omitempty"`
	LastModifiedBy    *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	BandwidthClasses  []common.IDNameExtensions `json:"bandwidthClasses,omitempty"`
	LocationGroups    []common.IDNameExtensions `json:"locationGroups,omitempty"`
	Labels            []common.IDNameExtensions `json:"labels,omitempty"`
	Devices           []common.IDNameExtensions `json:"devices,omitempty"`
	DeviceGroups      []common.IDNameExtensions `json:"deviceGroups,omitempty"`
	Locations         []common.IDNameExtensions `json:"locations,omitempty"`
	TimeWindows       []common.IDNameExtensions `json:"timeWindows,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*BandwidthControlRules, error) {
	var rule BandwidthControlRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", bandwidthControlEndpoint, ruleID), &rule)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning Bandwidth Control rule from Get: %d", rule.ID)
	return &rule, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*BandwidthControlRules, error) {
	var rules []BandwidthControlRules
	err := common.ReadAllPages(ctx, service.Client, bandwidthControlEndpoint, &rules)
	if err != nil {
		return nil, err
	}
	for _, rule := range rules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no Bandwidth Control rule rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, rule *BandwidthControlRules) (*BandwidthControlRules, error) {
	resp, err := service.Client.Create(ctx, bandwidthControlEndpoint, *rule)
	if err != nil {
		return nil, err
	}

	createdRules, ok := resp.(*BandwidthControlRules)
	if !ok {
		return nil, errors.New("object returned from api was not a Bandwidth Control rule Pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning Bandwidth Control rule from create: %d", createdRules.ID)
	return createdRules, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, rules *BandwidthControlRules) (*BandwidthControlRules, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", bandwidthControlEndpoint, ruleID), *rules)
	if err != nil {
		return nil, err
	}
	updatedRules, _ := resp.(*BandwidthControlRules)
	service.Client.GetLogger().Printf("[DEBUG]returning Bandwidth Control rule from update: %d", updatedRules.ID)
	return updatedRules, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", bandwidthControlEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAllLite(ctx context.Context, service *zscaler.Service) ([]BandwidthControlRules, error) {
	var profiles []BandwidthControlRules
	err := common.ReadAllPages(ctx, service.Client, bandwidthControlEndpoint+"/lite", &profiles)
	return profiles, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]BandwidthControlRules, error) {
	var rules []BandwidthControlRules
	err := common.ReadAllPages(ctx, service.Client, bandwidthControlEndpoint, &rules)
	return rules, err
}
