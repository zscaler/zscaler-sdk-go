package activation

import (
	"context"
	"errors"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	ecAdminActivateStatusEndpoint = "/ztw/api/v1/ecAdminActivateStatus"
	ecAdminActivateEndpoint       = "/ztw/api/v1/ecAdminActivateStatus/activate"
	ecAdminForceActivateEndpoint  = "/ztw/api/v1/ecAdminActivateStatus/forcedActivate"
)

type ECAdminActivation struct {
	OrgEditStatus         string                 `json:"orgEditStatus"`
	OrgLastActivateStatus string                 `json:"orgLastActivateStatus"`
	AdminStatusMap        map[string]interface{} `json:"adminStatusMap"`
	AdminActivateStatus   string                 `json:"adminActivateStatus"`
}

func GetActivationStatus(ctx context.Context, service *zscaler.Service) (*ECAdminActivation, error) {
	var ecAdminActivation ECAdminActivation
	err := service.Client.Read(ctx, ecAdminActivateStatusEndpoint, &ecAdminActivation)
	if err != nil {
		return nil, err
	}

	return &ecAdminActivation, nil
}

func UpdateActivationStatus(ctx context.Context, service *zscaler.Service, activation ECAdminActivation) (*ECAdminActivation, error) {
	resp, err := service.Client.UpdateWithPut(ctx, ecAdminActivateEndpoint, activation)
	if err != nil {
		return nil, err
	}

	updateActivationStatus, ok := resp.(*ECAdminActivation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return updateActivationStatus, nil
}

func ForceActivationStatus(ctx context.Context, service *zscaler.Service, forceActivation ECAdminActivation) (*ECAdminActivation, error) {
	resp, err := service.Client.UpdateWithPut(ctx, ecAdminForceActivateEndpoint, forceActivation)
	if err != nil {
		return nil, err
	}

	forceActivationStatus, ok := resp.(*ECAdminActivation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return forceActivationStatus, nil
}
