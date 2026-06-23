package ueba_rules

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
	uebaRulesEndpoint = "/zia/api/v1/alertRuleConfiguration/uebaRules"
)

type UebaRules struct {
	ID                 int                       `json:"id,omitempty"`
	AlertName          string                    `json:"alertName,omitempty"`
	AlertClass         string                    `json:"alertClass,omitempty"`
	Status             string                    `json:"status,omitempty"`
	EventTypes         []string                  `json:"eventTypes,omitempty"`
	WithinTime         int                       `json:"withinTime,omitempty"`
	MinTimes           int                       `json:"minTimes,omitempty"`
	WindowSize         int                       `json:"windowSize,omitempty"`
	EnableUpdate       bool                      `json:"enableUpdate,omitempty"`
	UpdateWindowSize   int                       `json:"updateWindowSize,omitempty"`
	NumSystemImpacted  int                       `json:"numSystemImpacted,omitempty"`
	Deleted            bool                      `json:"deleted,omitempty"`
	LastModifiedTime   int                       `json:"lastModifiedTime,omitempty"`
	Webhooks           []int                     `json:"webhooks,omitempty"`
	EmailIds           []string                  `json:"emailIds,omitempty"`
	Channel            string                    `json:"channel,omitempty"`
	AlertType          string                    `json:"alertType,omitempty"`
	Action             string                    `json:"action,omitempty"`
	ActionThreshold    int                       `json:"actionThreshold,omitempty"`
	Countries          []string                  `json:"countries,omitempty"`
	ObjectTypes        []string                  `json:"objectTypes,omitempty"`
	DocTypes           []string                  `json:"docTypes,omitempty"`
	ActionTimeInterval int                       `json:"actionTimeInterval,omitempty"`
	Activities         []string                  `json:"activities,omitempty"`
	LastModifiedBy     *common.IDNameExtensions  `json:"lastModifiedBy,omitempty"`
	UserGroupId        *common.IDNameExtensions  `json:"userGroupId,omitempty"`
	CasbApplications   []common.IDNameExtensions `json:"casbApplications,omitempty"`
	DlpEngines         []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	Locations          []common.IDNameExtensions `json:"locations,omitempty"`
	Users              []common.IDNameExtensions `json:"users,omitempty"`
	Departments        []common.IDNameExtensions `json:"departments,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, alertID int) (*UebaRules, error) {
	var alertDefinition UebaRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", uebaRulesEndpoint, alertID), &alertDefinition)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning alert definition from Get: %d", alertDefinition.ID)
	return &alertDefinition, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*UebaRules, error) {
	// Use GetAll to leverage the single API call and built-in pagination
	uebaRules, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, uebaRules := range uebaRules {
		if strings.EqualFold(uebaRules.AlertName, ruleName) {
			return &uebaRules, nil
		}
	}
	return nil, fmt.Errorf("no ueba rules found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, uebaRules *UebaRules) (*UebaRules, *http.Response, error) {
	resp, err := service.Client.Create(ctx, uebaRulesEndpoint, *uebaRules)
	if err != nil {
		return nil, nil, err
	}

	createdUebaRules, ok := resp.(*UebaRules)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a ueba rules pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new ueba rules from create: %d", createdUebaRules.ID)
	return createdUebaRules, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, uebaRules *UebaRules) (*UebaRules, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", uebaRulesEndpoint, ruleID), *uebaRules)
	if err != nil {
		return nil, nil, err
	}
	updatedUebaRules, _ := resp.(*UebaRules)

	service.Client.GetLogger().Printf("[DEBUG]returning updates ueba rules from update: %d", updatedUebaRules.ID)
	return updatedUebaRules, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", uebaRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]UebaRules, error) {
	var uebaRules []UebaRules
	err := common.ReadAllPages(ctx, service.Client, uebaRulesEndpoint, &uebaRules)
	return uebaRules, err
}
