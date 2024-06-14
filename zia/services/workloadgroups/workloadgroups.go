package workloadgroups

import (
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	workloadGroupsEndpoint = "/workloadGroups"
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

	// Information about the admin that last modified the workload group
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

func Get(service *services.Service, workloadID int) (*WorkloadGroup, error) {
	var workloadGroup WorkloadGroup
	err := service.Client.Read(fmt.Sprintf("%s/%d", workloadGroupsEndpoint, workloadID), &workloadGroup)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning workload group from Get: %d", workloadGroup.ID)
	return &workloadGroup, nil
}

func GetByName(service *services.Service, workloadName string) (*WorkloadGroup, error) {
	var workloadGroups []WorkloadGroup
	err := common.ReadAllPages(service.Client, workloadGroupsEndpoint, &workloadGroups)
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

func GetAll(service *services.Service) ([]WorkloadGroup, error) {
	var workloadGroups []WorkloadGroup
	err := common.ReadAllPages(service.Client, workloadGroupsEndpoint, &workloadGroups)
	return workloadGroups, err
}
