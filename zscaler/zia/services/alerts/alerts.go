package alerts

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	alertsEndpoint = "/zia/api/v1/alertSubscriptions"
)

type AlertSubscriptions struct {
	ID               int      `json:"id,omitempty"`
	Description      string   `json:"description,omitempty"`
	Email            string   `json:"email,omitempty"`
	Deleted          bool     `json:"deleted,omitempty"`
	Pt0Severities    []string `json:"pt0Severities,omitempty"`
	SecureSeverities []string `json:"secureSeverities,omitempty"`
	ManageSeverities []string `json:"manageSeverities,omitempty"`
	ComplySeverities []string `json:"complySeverities,omitempty"`
	SystemSeverities []string `json:"systemSeverities,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, subscriptionID int) (*AlertSubscriptions, error) {
	var subscription AlertSubscriptions
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", alertsEndpoint, subscriptionID), &subscription)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning alert subscription from Get: %d", subscription.ID)
	return &subscription, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AlertSubscriptions, error) {
	var alerts []AlertSubscriptions
	err := common.ReadAllPages(ctx, service.Client, alertsEndpoint, &alerts)
	return alerts, err
}

func Create(ctx context.Context, service *zscaler.Service, alerts *AlertSubscriptions) (*AlertSubscriptions, *http.Response, error) {
	resp, err := service.Client.Create(ctx, alertsEndpoint, *alerts)
	if err != nil {
		return nil, nil, err
	}

	createdAlert, ok := resp.(*AlertSubscriptions)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a alert subscription pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new alert subscription from create: %d", createdAlert.ID)
	return createdAlert, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, subscriptionID int, alerts *AlertSubscriptions) (*AlertSubscriptions, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", alertsEndpoint, subscriptionID), *alerts)
	if err != nil {
		return nil, nil, err
	}
	updatedAlert, _ := resp.(*AlertSubscriptions)

	service.Client.GetLogger().Printf("[DEBUG]returning updates alert subscription from update: %d", updatedAlert.ID)
	return updatedAlert, nil, nil
}
