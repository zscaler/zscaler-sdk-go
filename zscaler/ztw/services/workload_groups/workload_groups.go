package workload_groups

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/common"
)

const (
	workloadGroupsEndpoint = "/ztw/api/v1/workloadGroups"
)

type WorkloadGroup struct {
	// A unique identifier assigned to the workload group
	ID int `json:"id"`

	// The name of the workload group
	Name string `json:"name,omitempty"`

	// The description of the workload group
	Description string `json:"description,omitempty"`

	// The workload group expression containing tag types, tags, and their relationships.
	Expression string `json:"expression,omitempty"`

	// Timestamp when the workload group was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Information about the admin user that last modified the ZPA gateway
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	WorkloadTagExpression WorkloadTagExpression `json:"expressionJson,omitempty"`
}

// The workload group expression containing tag types, tags, and their relationships represented in a JSON format.
type WorkloadTagExpression struct {
	ExpressionContainers []ExpressionContainer `json:"expressionContainers"`
}

type ExpressionContainer struct {
	// The tag type selected from a predefined list
	TagType string `json:"tagType,omitempty"`

	// The operator (either AND or OR) used to create logical relationships among tag types
	Operator string `json:"operator,omitempty"`

	// Contains one or more tags and the logical operator used to combine the tags within a tag type
	TagContainer TagContainer `json:"tagContainer"`
}

type TagContainer struct {
	//One or more tags, each consisting of a key-value pair, selected within a tag type.
	// If multiple tags are present within a tag type, they are combined using a logical operator.
	// Note: A maximum of 8 tags can be added to a workload group, irrespective of the number of tag types present.
	Tags []Tags `json:"tags"`

	// The logical operator (either AND or OR) used to combine the tags within a tag type
	Operator string `json:"operator,omitempty"`
}

// The list of tags (key-value pairs) selected within a tag type
type Tags struct {
	// The key component present in the key-value pair contained in a tag
	Key string `json:"key,omitempty"`
	// The value component present in the key-value pair contained in a tag
	Value string `json:"value,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, workloadID int) (*WorkloadGroup, error) {
	var workloadGroup WorkloadGroup
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", workloadGroupsEndpoint, workloadID), &workloadGroup)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning workload group from Get: %d", workloadGroup.ID)
	return &workloadGroup, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, workloadName string) (*WorkloadGroup, error) {
	var workloadGroups []WorkloadGroup
	err := common.ReadAllPages(ctx, service.Client, workloadGroupsEndpoint, &workloadGroups)
	if err != nil {
		return nil, err
	}
	for _, workloadGroup := range workloadGroups {
		if strings.EqualFold(workloadGroup.Name, workloadName) {
			return &workloadGroup, nil
		}
	}
	return nil, fmt.Errorf("no workload group found with name: %s", workloadName)
}

/*
func Create(ctx context.Context, service *zscaler.Service, groups *WorkloadGroup) (*WorkloadGroup, *http.Response, error) {
	resp, err := service.Client.Create(ctx, workloadGroupsEndpoint, *groups)
	if err != nil {
		return nil, nil, err
	}

	createdWorkloadGroup, ok := resp.(*WorkloadGroup)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a workload group pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new workload group from create: %d", createdWorkloadGroup.ID)
	return createdWorkloadGroup, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, workloadGroupID int, workloadGroup *WorkloadGroup) (*WorkloadGroup, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", workloadGroupsEndpoint, workloadGroupID), *workloadGroup)
	if err != nil {
		return nil, nil, err
	}
	updatedWorkloadGroup, _ := resp.(*WorkloadGroup)

	service.Client.GetLogger().Printf("[DEBUG]returning updates workload group from update: %d", updatedWorkloadGroup.ID)
	return updatedWorkloadGroup, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, workloadGroupID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", workloadGroupsEndpoint, workloadGroupID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
*/

func GetAll(ctx context.Context, service *zscaler.Service) ([]WorkloadGroup, error) {
	var workloadGroups []WorkloadGroup
	err := common.ReadAllPages(ctx, service.Client, workloadGroupsEndpoint, &workloadGroups)
	return workloadGroups, err
}
