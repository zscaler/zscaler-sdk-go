package activation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	activationStatusEndpoint = "/zia/api/v1/status"
	activationEndpoint       = "/zia/api/v1/status/activate"
	eusaStatusEndpoint       = "/zia/api/v1/eusaStatus"
)

type Activation struct {
	Status string `json:"status"`
}

type ZiaEusaStatus struct {
	ID             int                      `json:"id,omitempty"`
	Version        *common.IDNameExtensions `json:"version,omitempty"`
	AcceptedStatus bool                     `json:"acceptedStatus,omitempty"`
}

func GetActivationStatus(ctx context.Context, service *zscaler.Service) (*Activation, error) {
	var activation Activation
	err := service.Client.Read(ctx, activationStatusEndpoint, &activation)
	if err != nil {
		return nil, err
	}

	return &activation, nil
}

func CreateActivation(ctx context.Context, service *zscaler.Service, activation Activation) (*Activation, error) {
	resp, err := service.Client.Create(ctx, activationEndpoint, activation)
	if err != nil {
		return nil, err
	}

	createdActivation, ok := resp.(*Activation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return createdActivation, nil
}

func GetEusaStatus(ctx context.Context, service *zscaler.Service) (*ZiaEusaStatus, error) {
	var status ZiaEusaStatus
	err := service.Client.Read(ctx, eusaStatusEndpoint+"/latest", &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func UpdateEusaStatus(ctx context.Context, service *zscaler.Service, statusID int, eusaStatus *ZiaEusaStatus) (*ZiaEusaStatus, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", eusaStatusEndpoint, statusID), *eusaStatus)
	if err != nil {
		return nil, nil, err
	}
	updatedStatus, _ := resp.(*ZiaEusaStatus)

	service.Client.GetLogger().Printf("[DEBUG]returning eusa status from update: %d", updatedStatus.ID)
	return updatedStatus, nil, nil
}
