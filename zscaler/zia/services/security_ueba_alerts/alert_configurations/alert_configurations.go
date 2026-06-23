package alert_configurations

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
	alertConfigurationEndpoint = "/zia/api/v1/alertRuleConfiguration/rules"
)

type AlertConfigurationRule struct {
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
	UserGroupID        *common.IDNameExtensions  `json:"userGroupId,omitempty"`
	CasbApplications   []common.IDNameExtensions `json:"casbApplications,omitempty"`
	DlpEngines         []common.IDNameExtensions `json:"dlpEngines,omitempty"`
	Locations          []common.IDNameExtensions `json:"locations,omitempty"`
	Users              []common.IDNameExtensions `json:"users,omitempty"`
	Departments        []common.IDNameExtensions `json:"departments,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, alertID int) (*AlertConfigurationRule, error) {
	var alertDefinition AlertConfigurationRule
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", alertConfigurationEndpoint, alertID), &alertDefinition)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning alert configuration rule from Get: %d", alertDefinition.ID)
	return &alertDefinition, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*AlertConfigurationRule, error) {
	// Use GetAll to leverage the single API call and built-in pagination
	alertConfiguration, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, alertDefinition := range alertConfiguration {
		if strings.EqualFold(alertDefinition.AlertName, ruleName) {
			return &alertDefinition, nil
		}
	}
	return nil, fmt.Errorf("no alert configuration rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, alertDefinitions *AlertConfigurationRule) (*AlertConfigurationRule, *http.Response, error) {
	resp, err := service.Client.Create(ctx, alertConfigurationEndpoint, *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}

	createdAlertDefinition, ok := resp.(*AlertConfigurationRule)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a alert configuration rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new alert configuration rule from create: %d", createdAlertDefinition.ID)
	return createdAlertDefinition, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, alertDefinitionID int, alertDefinitions *AlertConfigurationRule) (*AlertConfigurationRule, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", alertConfigurationEndpoint, alertDefinitionID), *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}
	updatedAlertDefinition, _ := resp.(*AlertConfigurationRule)

	service.Client.GetLogger().Printf("[DEBUG]returning updates alert configuration rule from update: %d", updatedAlertDefinition.ID)
	return updatedAlertDefinition, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, alertDefinitionID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", alertConfigurationEndpoint, alertDefinitionID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AlertConfigurationRule, error) {
	var alertConfiguration []AlertConfigurationRule
	err := common.ReadAllPages(ctx, service.Client, alertConfigurationEndpoint, &alertConfiguration)
	return alertConfiguration, err
}
