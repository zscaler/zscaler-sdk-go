package alert_definitions

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
	alertDefinitionsEndpoint = "/zia/api/v1/alertDefinitions"
)

type AlertDefinitions struct {
	ID                   int                      `json:"id,omitempty"`
	Status               string                   `json:"status,omitempty"`
	AlertName            string                   `json:"alertName,omitempty"`
	Occurrence           string                   `json:"occurrence,omitempty"`
	TrafficChangePercent int                      `json:"trafficChangePercent,omitempty"`
	Interval             string                   `json:"interval,omitempty"`
	Scope                string                   `json:"scope,omitempty"`
	Entity               *common.IDNameExtensions `json:"entity,omitempty"`
	Severity             string                   `json:"severity,omitempty"`
	Comments             string                   `json:"comments,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, alertID int) (*AlertDefinitions, error) {
	var alertDefinition AlertDefinitions
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", alertDefinitionsEndpoint, alertID), &alertDefinition)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning alert definition from Get: %d", alertDefinition.ID)
	return &alertDefinition, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, alertName string) (*AlertDefinitions, error) {
	// Use GetAll to leverage the single API call and built-in pagination
	alertDefinitions, err := GetAll(ctx, service)
	if err != nil {
		return nil, err
	}
	// Search for exact match (case-insensitive)
	for _, alertDefinition := range alertDefinitions {
		if strings.EqualFold(alertDefinition.AlertName, alertName) {
			return &alertDefinition, nil
		}
	}
	return nil, fmt.Errorf("no alert definition found with name: %s", alertName)
}

func Create(ctx context.Context, service *zscaler.Service, alertDefinitions *AlertDefinitions) (*AlertDefinitions, *http.Response, error) {
	resp, err := service.Client.Create(ctx, alertDefinitionsEndpoint, *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}

	createdAlertDefinition, ok := resp.(*AlertDefinitions)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a alert definition pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new alert definition from create: %d", createdAlertDefinition.ID)
	return createdAlertDefinition, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, alertDefinitionID int, alertDefinitions *AlertDefinitions) (*AlertDefinitions, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", alertDefinitionsEndpoint, alertDefinitionID), *alertDefinitions)
	if err != nil {
		return nil, nil, err
	}
	updatedAlertDefinition, _ := resp.(*AlertDefinitions)

	service.Client.GetLogger().Printf("[DEBUG]returning updates alert definition from update: %d", updatedAlertDefinition.ID)
	return updatedAlertDefinition, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, alertDefinitionID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", alertDefinitionsEndpoint, alertDefinitionID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AlertDefinitions, error) {
	var alertDefinitions []AlertDefinitions
	err := common.ReadAllPages(ctx, service.Client, alertDefinitionsEndpoint, &alertDefinitions)
	return alertDefinitions, err
}
