package cloud_application_instances

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	cloudApplicationInstancesEndpoint = "/cloudApplicationInstances"
)

type CloudApplicationInstances struct {
	InstanceID          int                      `json:"instanceId"`
	InstanceType        string                   `json:"instanceType,omitempty"`
	InstanceName        string                   `json:"instanceName,omitempty"`
	ModifiedBy          *common.IDNameExtensions `json:"modifiedBy,omitempty"`
	ModifiedAt          int                      `json:"modifiedAt,omitempty"`
	InstanceIdentifiers []InstanceIdentifiers    `json:"instanceIdentifiers,omitempty"`
}

type InstanceIdentifiers struct {
	InstanceID             int                      `json:"instanceId"`
	InstanceIdentifier     string                   `json:"instanceIdentifier,omitempty"`
	InstanceIdentifierName string                   `json:"instanceIdentifierName,omitempty"`
	IdentifierType         string                   `json:"identifierType,omitempty"`
	ModifiedBy             *common.IDNameExtensions `json:"modifiedBy,omitempty"`
	ModifiedAt             int                      `json:"modifiedAt,omitempty"`
}

func (service *Service) GetCloudApplicationInstanceID(cloudInstanceID int) (*CloudApplicationInstances, error) {
	var cloudInstance CloudApplicationInstances
	err := service.Client.Read(fmt.Sprintf("%s/%d", cloudApplicationInstancesEndpoint, cloudInstanceID), &cloudInstance)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning cloud instance ID from Get: %d", cloudInstance.InstanceID)
	return &cloudInstance, nil
}

func (service *Service) GetCloudApplicationInstanceByName(cloudInstanceName string) (*CloudApplicationInstances, error) {
	var cloudInstances []CloudApplicationInstances
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?name=%s", cloudApplicationInstancesEndpoint, url.QueryEscape(cloudInstanceName)), &cloudInstances)
	if err != nil {
		return nil, err
	}
	for _, cloudInstance := range cloudInstances {
		if strings.EqualFold(cloudInstance.InstanceName, cloudInstanceName) {
			return &cloudInstance, nil
		}
	}
	return nil, fmt.Errorf("no cloud instance found with name: %s", cloudInstanceName)
}
