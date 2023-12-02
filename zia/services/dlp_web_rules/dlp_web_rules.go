package dlp_web_rules

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	webDlpRulesEndpoint = "/webDlpRules"
)

type WebDLPRules struct {
	// The unique identifier for the DLP policy rule.
	ID int `json:"id,omitempty"`

	// The rule order of execution for the DLP policy rule with respect to other rules.
	Order int `json:"order,omitempty"`
	// The access privilege for this DLP policy rule based on the admin's state.
	AccessControl string `json:"accessControl,omitempty"`

	// The protocol criteria specified for the DLP policy rule.
	Protocols []string `json:"protocols,omitempty"`

	// The admin rank of the admin who created the DLP policy rule.
	Rank int `json:"rank,omitempty"`

	// The DLP policy rule name.
	Name string `json:"name,omitempty"`

	// The description of the DLP policy rule.
	Description string `json:"description,omitempty"`

	// The list of file types to which the DLP policy rule must be applied.
	FileTypes []string `json:"fileTypes,omitempty"`

	// The list of cloud applications to which the DLP policy rule must be applied.
	CloudApplications []string `json:"cloudApplications,omitempty"`

	// The minimum file size (in KB) used for evaluation of the DLP policy rule.
	MinSize int `json:"minSize,omitempty"`

	// The action taken when traffic matches the DLP policy rule criteria.
	Action string `json:"action,omitempty"`

	// Enables or disables the DLP policy rule.
	State string `json:"state,omitempty"`

	// The match only criteria for DLP engines.
	MatchOnly bool `json:"matchOnly,omitempty"`

	// Timestamp when the DLP policy rule was last modified.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// Indicates a DLP policy rule without content inspection, when the value is set to true.
	WithoutContentInspection bool `json:"withoutContentInspection,omitempty"`

	// Enables or disables image file scanning.
	OcrEnabled bool `json:"ocrEnabled,omitempty"`

	DLPDownloadScanEnabled  bool `json:"dlpDownloadScanEnabled,omitempty"`
	ZCCNotificationsEnabled bool `json:"zccNotificationsEnabled,omitempty"`

	// Indicates whether a Zscaler Incident Receiver is associated to the DLP policy rule.
	ZscalerIncidentReceiver bool `json:"zscalerIncidentReceiver,omitempty"`

	// The email address of an external auditor to whom DLP email notifications are sent.
	ExternalAuditorEmail string `json:"externalAuditorEmail,omitempty"`

	// The auditor to which the DLP policy rule must be applied.
	Auditor *common.IDNameExtensions `json:"auditor,omitempty"`

	// The admin that modified the DLP policy rule last.
	LastModifiedBy *common.IDNameExtensions `json:"lastModifiedBy,omitempty"`

	// The template used for DLP notification emails.
	NotificationTemplate *common.IDNameExtensions `json:"notificationTemplate,omitempty"`

	// The DLP server, using ICAP, to which the transaction content is forwarded.
	IcapServer *common.IDNameExtensions `json:"icapServer,omitempty"`

	// The Name-ID pairs of locations to which the DLP policy rule must be applied.
	Locations []common.IDNameExtensions `json:"locations,omitempty"`

	// The Name-ID pairs of locations groups to which the DLP policy rule must be applied.
	LocationGroups []common.IDNameExtensions `json:"locationGroups,omitempty"`

	// The Name-ID pairs of groups to which the DLP policy rule must be applied.
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// The Name-ID pairs of departments to which the DLP policy rule must be applied.
	Departments []common.IDNameExtensions `json:"departments,omitempty"`

	// The Name-ID pairs of users to which the DLP policy rule must be applied.
	Users []common.IDNameExtensions `json:"users,omitempty"`

	// The list of URL categories to which the DLP policy rule must be applied.
	URLCategories []common.IDNameExtensions `json:"urlCategories,omitempty"`

	// The list of DLP engines to which the DLP policy rule must be applied.
	DLPEngines []common.IDNameExtensions `json:"dlpEngines,omitempty"`

	// The Name-ID pairs of time windows to which the DLP policy rule must be applied.
	TimeWindows []common.IDNameExtensions `json:"timeWindows,omitempty"`

	// The Name-ID pairs of rule labels associated to the DLP policy rule.
	Labels []common.IDNameExtensions `json:"labels,omitempty"`

	// The name-ID pairs of the groups that are excluded from the DLP policy rule.
	ExcludedGroups []common.IDNameExtensions `json:"excludedGroups,omitempty"`

	// The name-ID pairs of the departments that are excluded from the DLP policy rule.
	ExcludedDepartments []common.IDNameExtensions `json:"excludedDepartments,omitempty"`

	// The name-ID pairs of the users that are excluded from the DLP policy rule.
	ExcludedUsers []common.IDNameExtensions `json:"excludedUsers,omitempty"`
}

func (service *Service) Get(ruleID int) (*WebDLPRules, error) {
	var webDlpRules WebDLPRules
	err := service.Client.Read(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID), &webDlpRules)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning web dlp rule from Get: %d", webDlpRules.ID)
	return &webDlpRules, nil
}

func (service *Service) GetByName(ruleName string) (*WebDLPRules, error) {
	var webDlpRules []WebDLPRules
	err := common.ReadAllPages(service.Client, webDlpRulesEndpoint, &webDlpRules)
	if err != nil {
		return nil, err
	}
	for _, rule := range webDlpRules {
		if strings.EqualFold(rule.Name, ruleName) {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("no web dlp rule found with name: %s", ruleName)
}

func (service *Service) Create(ruleID *WebDLPRules) (*WebDLPRules, error) {
	resp, err := service.Client.Create(webDlpRulesEndpoint, *ruleID)
	if err != nil {
		return nil, err
	}

	createdWebDlpRules, ok := resp.(*WebDLPRules)
	if !ok {
		return nil, errors.New("object returned from api was not a web dlp rule pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning new web dlp rule from create: %d", createdWebDlpRules.ID)
	return createdWebDlpRules, nil
}

func (service *Service) Update(ruleID int, webDlpRules *WebDLPRules) (*WebDLPRules, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID), *webDlpRules)
	if err != nil {
		return nil, err
	}
	updatedWebDlpRules, _ := resp.(*WebDLPRules)

	service.Client.Logger.Printf("[DEBUG]returning updates from web dlp rule from update: %d", updatedWebDlpRules.ID)
	return updatedWebDlpRules, nil
}

func (service *Service) Delete(ruleID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", webDlpRulesEndpoint, ruleID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]WebDLPRules, error) {
	var webDlpRules []WebDLPRules
	err := common.ReadAllPages(service.Client, webDlpRulesEndpoint, &webDlpRules)
	return webDlpRules, err
}
