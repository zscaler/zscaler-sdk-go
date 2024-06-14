package activation

import (
	"errors"

	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services"
)

const (
	ecAdminActivateStatusEndpoint = "/ecAdminActivateStatus"
	ecAdminActivateEndpoint       = "/ecAdminActivateStatus/activate"
	ecAdminForceActivateEndpoint  = "/ecAdminActivateStatus/forcedActivate"
)

type ECAdminActivation struct {
	OrgEditStatus         string                 `json:"orgEditStatus"`
	OrgLastActivateStatus string                 `json:"orgLastActivateStatus"`
	AdminStatusMap        map[string]interface{} `json:"adminStatusMap"`
	AdminActivateStatus   string                 `json:"adminActivateStatus"`
}

func GetActivationStatus(service *services.Service) (*ECAdminActivation, error) {
	var ecAdminActivation ECAdminActivation
	err := service.Client.Read(ecAdminActivateStatusEndpoint, &ecAdminActivation)
	if err != nil {
		return nil, err
	}

	return &ecAdminActivation, nil
}

func UpdateActivationStatus(service *services.Service, activation ECAdminActivation) (*ECAdminActivation, error) {
	resp, err := service.Client.UpdateWithPut(ecAdminActivateEndpoint, activation)
	if err != nil {
		return nil, err
	}

	updateActivationStatus, ok := resp.(*ECAdminActivation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return updateActivationStatus, nil
}

func ForceActivationStatus(service *services.Service, forceActivation ECAdminActivation) (*ECAdminActivation, error) {
	resp, err := service.Client.UpdateWithPut(ecAdminForceActivateEndpoint, forceActivation)
	if err != nil {
		return nil, err
	}

	forceActivationStatus, ok := resp.(*ECAdminActivation)
	if !ok {
		return nil, errors.New("object returned from api was not an activation pointer")
	}

	return forceActivationStatus, nil
}
