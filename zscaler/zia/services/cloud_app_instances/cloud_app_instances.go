package cloud_app_instances

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
	cloudInstancesEndpoint = "/zia/api/v1/cloudApplicationInstances"
)

type CloudApplicationInstances struct {
	InstanceID          int                      `json:"instanceId,omitempty"`
	InstanceType        string                   `json:"instanceType,omitempty"`
	InstanceName        string                   `json:"instanceName,omitempty"`
	ModifiedBy          *common.IDNameExtensions `json:"modifiedBy,omitempty"`
	ModifiedAt          int                      `json:"modifiedAt,omitempty"`
	InstanceIdentifiers []InstanceIdentifiers    `json:"instanceIdentifiers,omitempty"`
}

type InstanceIdentifiers struct {
	InstanceID             int                      `json:"instanceId,omitempty"`
	InstanceIdentifier     string                   `json:"instanceIdentifier,omitempty"`
	InstanceIdentifierName string                   `json:"instanceIdentifierName,omitempty"`
	IdentifierType         string                   `json:"identifierType,omitempty"`
	ModifiedAt             int                      `json:"modifiedAt,omitempty"`
	ModifiedBy             *common.IDNameExtensions `json:"modifiedBy,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, instanceID int) (*CloudApplicationInstances, error) {
	var cloudInstance CloudApplicationInstances
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", cloudInstancesEndpoint, instanceID), &cloudInstance)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning cloud instances from Get: %d", cloudInstance.InstanceID)
	return &cloudInstance, nil
}

func GetInstanceByName(ctx context.Context, service *zscaler.Service, instanceName string) (*CloudApplicationInstances, error) {
	var cloudInstances []CloudApplicationInstances
	err := common.ReadAllPages(ctx, service.Client, cloudInstancesEndpoint, &cloudInstances)
	if err != nil {
		return nil, err
	}
	for _, cloudInstance := range cloudInstances {
		if strings.EqualFold(cloudInstance.InstanceName, instanceName) {
			return &cloudInstance, nil
		}
	}
	return nil, fmt.Errorf("no cloud instance found with name: %s", instanceName)
}

func Create(ctx context.Context, service *zscaler.Service, instanceID *CloudApplicationInstances) (*CloudApplicationInstances, *http.Response, error) {
	resp, err := service.Client.Create(ctx, cloudInstancesEndpoint, *instanceID)
	if err != nil {
		return nil, nil, err
	}

	createdCloudInstance, ok := resp.(*CloudApplicationInstances)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a cloud instance pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new cloud instance from create: %d", createdCloudInstance.InstanceID)
	return createdCloudInstance, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, instanceID int, cloudInstance *CloudApplicationInstances) (*CloudApplicationInstances, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", cloudInstancesEndpoint, instanceID), *cloudInstance)
	if err != nil {
		return nil, nil, err
	}
	updatedCloudInstance, _ := resp.(*CloudApplicationInstances)

	service.Client.GetLogger().Printf("[DEBUG]returning updates cloud instance from update: %d", updatedCloudInstance.InstanceID)
	return updatedCloudInstance, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, instanceID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", cloudInstancesEndpoint, instanceID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CloudApplicationInstances, error) {
	var cloudInstances []CloudApplicationInstances
	err := common.ReadAllPages(ctx, service.Client, cloudInstancesEndpoint, &cloudInstances)
	return cloudInstances, err
}
