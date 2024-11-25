package activation

import (
	"context"
	"errors"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	activationStatusEndpoint = "/zia/api/v1/status"
	activationEndpoint       = "/zia/api/v1/status/activate"
)

type Activation struct {
	Status string `json:"status"`
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
