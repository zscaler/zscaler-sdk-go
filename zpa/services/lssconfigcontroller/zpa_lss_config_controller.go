package lssconfigcontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                = "/mgmtconfig/v2/admin/customers/"
	mgmtConfigTypesAndFormats = "/mgmtconfig/v2/admin/"
	lssConfigEndpoint         = "/lssConfig"
)

type LSSResource struct {
	ID                 string              `json:"id,omitempty"`
	LSSConfig          *LSSConfig          `json:"config"`
	ConnectorGroups    []ConnectorGroups   `json:"connectorGroups,omitempty"`
	PolicyRule         *PolicyRule         `json:"policyRule,omitempty"`
	PolicyRuleResource *PolicyRuleResource `json:"policyRuleResource,omitempty"`
}

type LSSConfig struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	CreationTime    string   `json:"creationTime,omitempty"`
	ModifiedBy      string   `json:"modifiedBy,omitempty"`
	ModifiedTime    string   `json:"modifiedTime,omitempty"`
	Filter          []string `json:"filter,omitempty"`
	Format          string   `json:"format,omitempty"`
	AuditMessage    string   `json:"auditMessage,omitempty"`
	LSSHost         string   `json:"lssHost,omitempty"`
	LSSPort         string   `json:"lssPort,omitempty"`
	SourceLogType   string   `json:"sourceLogType,omitempty"`
	MicroTenantID   string   `json:"microtenantId,omitempty"`
	MicroTenantName string   `json:"microtenantName,omitempty"`
	UseTLS          bool     `json:"useTls,omitempty"`
}

type ConnectorGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AppServerGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PolicyRuleResource struct {
	ID                       string                         `json:"id,omitempty"`
	Name                     string                         `json:"name,omitempty"`
	Description              string                         `json:"description,omitempty"`
	Action                   string                         `json:"action,omitempty"`
	ActionID                 string                         `json:"actionId,omitempty"`
	ConnectorGroups          []ConnectorGroups              `json:"connectorGroups,omitempty"`
	AppServerGroups          []AppServerGroups              `json:"appServerGroups,omitempty"`
	CreationTime             string                         `json:"creationTime,omitempty"`
	ModifiedBy               string                         `json:"modifiedBy,omitempty"`
	ModifiedTime             string                         `json:"modifiedTime,omitempty"`
	AuditMessage             string                         `json:"auditMessage,omitempty"`
	CustomMsg                string                         `json:"customMsg,omitempty"`
	Operator                 string                         `json:"operator,omitempty"`
	PolicySetID              string                         `json:"policySetId,omitempty"`
	PolicyType               string                         `json:"policyType,omitempty"`
	Priority                 string                         `json:"priority,omitempty"`
	ReauthIdleTimeout        string                         `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string                         `json:"reauthTimeout,omitempty"`
	RuleOrder                string                         `json:"ruleOrder,omitempty"`
	ZpnCbiProfileID          string                         `json:"zpnCbiProfileId,omitempty"`
	ZpnInspectionProfileID   string                         `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                         `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string                         `json:"microtenantId,omitempty"`
	MicroTenantName          string                         `json:"microtenantName,omitempty"`
	Conditions               []PolicyRuleResourceConditions `json:"conditions,omitempty"`
}

type PolicyRule struct {
	Action                   string       `json:"action,omitempty"`
	ActionID                 string       `json:"actionId,omitempty"`
	BypassDefaultRule        bool         `json:"bypassDefaultRule,omitempty"`
	CreationTime             string       `json:"creationTime,omitempty"`
	CustomMsg                string       `json:"customMsg,omitempty"`
	DefaultRule              bool         `json:"defaultRule,omitempty"`
	Description              string       `json:"description,omitempty"`
	ID                       string       `json:"id,omitempty"`
	IsolationDefaultRule     bool         `json:"isolationDefaultRule,omitempty"`
	ModifiedBy               string       `json:"modifiedBy,omitempty"`
	ModifiedTime             string       `json:"modifiedTime,omitempty"`
	Name                     string       `json:"name,omitempty"`
	Operator                 string       `json:"operator,omitempty"`
	PolicySetID              string       `json:"policySetId,omitempty"`
	PolicyType               string       `json:"policyType,omitempty"`
	Priority                 string       `json:"priority,omitempty"`
	ReauthDefaultRule        bool         `json:"reauthDefaultRule,omitempty"`
	ReauthIdleTimeout        string       `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string       `json:"reauthTimeout,omitempty"`
	RuleOrder                string       `json:"ruleOrder,omitempty"`
	LssDefaultRule           bool         `json:"lssDefaultRule,omitempty"`
	ZpnCbiProfileID          string       `json:"zpnCbiProfileId,omitempty"`
	ZpnInspectionProfileID   string       `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string       `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string       `json:"microtenantId,omitempty"`
	MicroTenantName          string       `json:"microtenantName,omitempty"`
	Conditions               []Conditions `json:"conditions,omitempty"`
}

type Conditions struct {
	CreationTime string      `json:"creationTime,omitempty"`
	ID           string      `json:"id,omitempty"`
	ModifiedBy   string      `json:"modifiedBy,omitempty"`
	ModifiedTime string      `json:"modifiedTime,omitempty"`
	Negated      bool        `json:"negated"`
	Operands     *[]Operands `json:"operands,omitempty"`
	Operator     string      `json:"operator,omitempty"`
}

type PolicyRuleResourceConditions struct {
	ID           string                        `json:"id,omitempty"`
	CreationTime string                        `json:"creationTime,omitempty"`
	ModifiedBy   string                        `json:"modifiedBy,omitempty"`
	ModifiedTime string                        `json:"modifiedTime,omitempty"`
	Negated      bool                          `json:"negated"`
	Operands     *[]PolicyRuleResourceOperands `json:"operands,omitempty"`
	Operator     string                        `json:"operator,omitempty"`
}

type PolicyRuleResourceOperands struct {
	ID                          string                         `json:"id,omitempty"`
	CreationTime                string                         `json:"creationTime,omitempty"`
	ModifiedBy                  string                         `json:"modifiedBy,omitempty"`
	ModifiedTime                string                         `json:"modifiedTime,omitempty"`
	ObjectType                  string                         `json:"objectType,omitempty"`
	Values                      []string                       `json:"values,omitempty"`
	IDPID                       string                         `json:"idpId,omitempty"`
	OperandsResourceLHSRHSValue *[]OperandsResourceLHSRHSValue `json:"entryValues,omitempty"`
}

type OperandsResourceLHSRHSValue struct {
	RHS string `json:"rhs,omitempty"`
	LHS string `json:"lhs,omitempty"`
}

type Operands struct {
	CreationTime string `json:"creationTime,omitempty"`
	ID           string `json:"id,omitempty"`
	IdpID        string `json:"idpId,omitempty"`
	LHS          string `json:"lhs,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	ObjectType   string `json:"objectType,omitempty"`
	RHS          string `json:"rhs,omitempty"`
}

func (service *Service) Get(lssID string) (*LSSResource, *http.Response, error) {
	v := new(LSSResource)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(lssName string) (*LSSResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + lssConfigEndpoint
	list, resp, err := common.GetAllPagesGeneric[LSSResource](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, lss := range list {
		if strings.EqualFold(lss.LSSConfig.Name, lssName) {
			return &lss, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no lss controller named '%s' was found", lssName)
}

func (service *Service) Create(lssConfig *LSSResource) (*LSSResource, *http.Response, error) {
	v := new(LSSResource)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, nil, lssConfig, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(lssID string, lssConfig *LSSResource) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, lssConfig, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(lssID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+lssConfigEndpoint, lssID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]LSSResource, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + lssConfigEndpoint
	list, resp, err := common.GetAllPagesGeneric[LSSResource](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
