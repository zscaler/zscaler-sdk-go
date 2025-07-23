package filetypecontrol

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
	fileTypeControlEndpoint = "/zia/api/v1/fileTypeRules"
)

type FileTypeRules struct {
	// System generated identifier for a file-type policy
	ID int `json:"id,omitempty"`

	// The name of the File Type rule
	Name string `json:"name,omitempty"`

	// Additional information about the File Type rule
	Description string `json:"description,omitempty"`

	// Rule State
	State string `json:"state,omitempty"`

	// Order of policy execution with respect to other file-type policies. Order N indicates N-th file-type policy is evaluated. This field is not applicable to the Lite API.
	Order int `json:"order,omitempty"`

	// Action taken when traffic matches policy. This field is not applicable to the Lite API.
	// Supported Values: "BLOCK", "CAUTION", "ALLOW"
	FilteringAction string `json:"filteringAction,omitempty"`

	// Time quota in minutes, after which the policy must be applied. If not set, no time quota is enforced. Ignored if action is BLOCK.
	TimeQuota int `json:"timeQuota,omitempty"`

	// Size quota in KB, beyond which the policy must be applied. If not set, size quota is not enforced. Ignored if action is BLOCK.
	SizeQuota int `json:"sizeQuota,omitempty"`

	// The access privilege for this DLP policy rule based on the admin's state.
	AccessControl string `json:"accessControl,omitempty"`

	// Admin rank of the admin who creates this rule
	Rank int `json:"rank,omitempty"`

	// A Boolean value that indicates whether packet capture (PCAP) is enabled or not
	CapturePCAP bool `json:"capturePCAP"`

	// File operation performed. This field is not applicable to the Lite API.
	Operation string `json:"operation"`

	// Flag to check whether a file has active content or not
	ActiveContent bool `json:"activeContent"`

	// Flag to check whether a file has active content or not
	Unscannable bool `json:"unscannable"`

	BrowserEunTemplateID int `json:"browserEunTemplateId,omitempty"`

	// The list of cloud applications to which the File Type Control policy rule must be applied
	// New to Review that.
	CloudApplications []string `json:"cloudApplications,omitempty"`

	// The list of file types to which the Sandbox Rule must be applied.
	FileTypes []string `json:"fileTypes,omitempty"`

	// Minimum file size (in KB) used for evaluation of the FTP rule
	MinSize int `json:"minSize,omitempty"`

	// Maximum file size (in KB) used for evaluation of the FTP rule
	MaxSize int `json:"maxSize,omitempty"`

	// Protocol for the given rule. This field is not applicable to the Lite API.
	Protocols []string `json:"protocols,omitempty"`

	// The list of URL categories to which the DLP policy rule must be applied.
	URLCategories []string `json:"urlCategories,omitempty"`

	// When the rule was last modified
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Who modified the rule last
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// Name-ID pairs of locations for which rule must be applied
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// Name-ID pairs of the location groups to which the rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// Name-ID pairs of groups for which rule must be applied
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// Name-ID pairs of departments for which rule must be applied
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// Name-ID pairs of users for which rule must be applied
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// Name-ID pairs of time interval during which rule must be enforced.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// The URL Filtering rule's label. Rule labels allow you to logically group your organization's policy rules. Policy rules that are not associated with a rule label are grouped under the Untagged label.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// This field is applicable for devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	DeviceGroups []common.IDNameExtensions `json:"deviceGroups"`

	// Name-ID pairs of devices for which rule must be applied. Specifies devices that are managed using Zscaler Client Connector. If no value is set, this field is ignored during the policy evaluation.
	Devices []common.IDNameExtensions `json:"devices"`

	// List of device trust levels for which the rule must be applied. This field is applicable for devices that are managed using Zscaler Client Connector. The trust levels are assigned to the devices based on your posture configurations in the Zscaler Client Connector Portal. If no value is set, this field is ignored during the policy evaluation.
	DeviceTrustLevels []string `json:"deviceTrustLevels,omitempty"`

	// The list of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ZPA Gateway forwarding method.
	ZPAAppSegments []common.ZPAAppSegments `json:"zpaAppSegments"`
}

func Get(ctx context.Context, service *zscaler.Service, ruleID int) (*FileTypeRules, error) {
	var fileTypes FileTypeRules
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", fileTypeControlEndpoint, ruleID), &fileTypes)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]Returning file type control rules from Get: %d", fileTypes.ID)
	return &fileTypes, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, ruleName string) (*FileTypeRules, error) {
	var fileTypeControlRules []FileTypeRules
	err := common.ReadAllPages(ctx, service.Client, fileTypeControlEndpoint, &fileTypeControlRules)
	if err != nil {
		return nil, err
	}
	for _, fileTypeControlRule := range fileTypeControlRules {
		if strings.EqualFold(fileTypeControlRule.Name, ruleName) {
			return &fileTypeControlRule, nil
		}
	}
	return nil, fmt.Errorf("no file type control rule found with name: %s", ruleName)
}

func Create(ctx context.Context, service *zscaler.Service, ruleID *FileTypeRules) (*FileTypeRules, error) {
	resp, err := service.Client.Create(ctx, fileTypeControlEndpoint, *ruleID)
	if err != nil {
		return nil, err
	}

	createdFileTypeControlRule, ok := resp.(*FileTypeRules)
	if !ok {
		return nil, errors.New("object returned from api was not a file type control rule pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning file type control rule from create: %d", createdFileTypeControlRule.ID)
	return createdFileTypeControlRule, nil
}

func Update(ctx context.Context, service *zscaler.Service, ruleID int, fileTypeRules *FileTypeRules) (*FileTypeRules, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", fileTypeControlEndpoint, ruleID), *fileTypeRules)
	if err != nil {
		return nil, err
	}
	updatedFileTypeControlRule, _ := resp.(*FileTypeRules)

	service.Client.GetLogger().Printf("[DEBUG]returning updates from file type control rule from update: %d", updatedFileTypeControlRule.ID)
	return updatedFileTypeControlRule, nil
}

func Delete(ctx context.Context, service *zscaler.Service, ruleID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", fileTypeControlEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]FileTypeRules, error) {
	var fileTypeRules []FileTypeRules
	err := common.ReadAllPages(ctx, service.Client, fileTypeControlEndpoint, &fileTypeRules)
	return fileTypeRules, err
}
