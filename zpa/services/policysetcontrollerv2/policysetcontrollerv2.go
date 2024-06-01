package policysetcontrollerv2

import (
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV2 = "/mgmtconfig/v2/admin/customers/"
)

type PolicySet struct {
	CreationTime    string       `json:"creationTime,omitempty"`
	Description     string       `json:"description,omitempty"`
	Enabled         bool         `json:"enabled"`
	ID              string       `json:"id,omitempty"`
	ModifiedBy      string       `json:"modifiedBy,omitempty"`
	ModifiedTime    string       `json:"modifiedTime,omitempty"`
	Name            string       `json:"name,omitempty"`
	Sorted          bool         `json:"sorted"`
	PolicyType      string       `json:"policyType,omitempty"`
	MicroTenantID   string       `json:"microtenantId,omitempty"`
	MicroTenantName string       `json:"microtenantName,omitempty"`
	Rules           []PolicyRule `json:"rules"`
}

// ######################################################################################################
// ########################################## API V1 Structure ##########################################
// ################################### Used to process the API Response #################################
// ######################################################################################################
type PolicyRuleResource struct {
	ID                       string                         `json:"id,omitempty"`
	Name                     string                         `json:"name,omitempty"`
	Description              string                         `json:"description,omitempty"`
	Action                   string                         `json:"action,omitempty"`
	ActionID                 string                         `json:"actionId,omitempty"`
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
	ZpnIsolationProfileID    string                         `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string                         `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                         `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string                         `json:"microtenantId,omitempty"`
	MicroTenantName          string                         `json:"microtenantName,omitempty"`
	Conditions               []PolicyRuleResourceConditions `json:"conditions,omitempty"`
	AppConnectorGroups       []AppConnectorGroups           `json:"connectorGroups,omitempty"`
	AppServerGroups          []AppServerGroups              `json:"appServerGroups,omitempty"`
	ServiceEdgeGroups        []ServiceEdgeGroups            `json:"serviceEdgeGroups,omitempty"`
	Credential               *Credential                    `json:"credential,omitempty"`
	PrivilegedCapabilities   PrivilegedCapabilities         `json:"privilegedCapabilities,omitempty"`
}

type Conditions struct {
	CreationTime string     `json:"creationTime,omitempty"`
	ID           string     `json:"id,omitempty"`
	ModifiedBy   string     `json:"modifiedBy,omitempty"`
	ModifiedTime string     `json:"modifiedTime,omitempty"`
	Negated      bool       `json:"negated"`
	Operands     []Operands `json:"operands,omitempty"`
	Operator     string     `json:"operator,omitempty"`
}

type Operands struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	CreationTime string `json:"creationTime,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	IdpID        string `json:"idpId,omitempty"`
	LHS          string `json:"lhs,omitempty"`
	RHS          string `json:"rhs,omitempty"`
	ObjectType   string `json:"objectType,omitempty"`
}

// ######################################################################################################
// ########################################## API V2 Structure ##########################################
// ################################### Used to process the API Request ##################################
// ######################################################################################################

type PolicyRule struct {
	ID                       string                         `json:"id,omitempty"`
	Name                     string                         `json:"name,omitempty"`
	Action                   string                         `json:"action,omitempty"`
	ActionID                 string                         `json:"actionId,omitempty"`
	CustomMsg                string                         `json:"customMsg,omitempty"`
	Description              string                         `json:"description,omitempty"`
	CreationTime             string                         `json:"creationTime,omitempty"`
	ModifiedBy               string                         `json:"modifiedBy,omitempty"`
	ModifiedTime             string                         `json:"modifiedTime,omitempty"`
	Operator                 string                         `json:"operator,omitempty"`
	PolicySetID              string                         `json:"policySetId,omitempty"`
	PolicyType               string                         `json:"policyType,omitempty"`
	Priority                 string                         `json:"priority,omitempty"`
	ReauthIdleTimeout        string                         `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string                         `json:"reauthTimeout,omitempty"`
	RuleOrder                string                         `json:"ruleOrder,omitempty"`
	ZpnIsolationProfileID    string                         `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string                         `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                         `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string                         `json:"microtenantId,omitempty"`
	MicroTenantName          string                         `json:"microtenantName,omitempty"`
	Version                  string                         `json:"version,omitempty"`
	AppConnectorGroups       []AppConnectorGroups           `json:"connectorGroups,omitempty"`
	AppServerGroups          []AppServerGroups              `json:"appServerGroups,omitempty"`
	ServiceEdgeGroups        []ServiceEdgeGroups            `json:"serviceEdgeGroups,omitempty"`
	Conditions               []PolicyRuleResourceConditions `json:"conditions,omitempty"`
	Credential               *Credential                    `json:"credential,omitempty"`
	PrivilegedCapabilities   PrivilegedCapabilities         `json:"privilegedCapabilities,omitempty"`
}

type PolicyRuleResourceConditions struct {
	ID           string                       `json:"id,omitempty"`
	CreationTime string                       `json:"creationTime,omitempty"`
	ModifiedBy   string                       `json:"modifiedBy,omitempty"`
	ModifiedTime string                       `json:"modifiedTime,omitempty"`
	Negated      bool                         `json:"negated"`
	Operator     string                       `json:"operator,omitempty"`
	Operands     []PolicyRuleResourceOperands `json:"operands,omitempty"`
}

type PolicyRuleResourceOperands struct {
	ID                string                        `json:"id,omitempty"`
	CreationTime      string                        `json:"creationTime,omitempty"`
	ModifiedBy        string                        `json:"modifiedBy,omitempty"`
	ModifiedTime      string                        `json:"modifiedTime,omitempty"`
	ObjectType        string                        `json:"objectType,omitempty"`
	Values            []string                      `json:"values,omitempty"`
	IDPID             string                        `json:"idpId,omitempty"`
	LHS               string                        `json:"lhs,omitempty"`
	RHS               string                        `json:"rhs,omitempty"`
	EntryValuesLHSRHS []OperandsResourceLHSRHSValue `json:"entryValues,omitempty"`
}

type OperandsResourceLHSRHSValue struct {
	RHS string `json:"rhs,omitempty"`
	LHS string `json:"lhs,omitempty"`
}

type AppServerGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AppConnectorGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ServiceEdgeGroups struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Credential struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type PrivilegedCapabilities struct {
	ID            string   `json:"id"`
	CreationTime  string   `json:"creationTime,omitempty"`
	ModifiedBy    string   `json:"modifiedBy,omitempty"`
	ModifiedTime  string   `json:"modifiedTime,omitempty"`
	MicroTenantID string   `json:"microtenantId,omitempty"`
	Capabilities  []string `json:"capabilities,omitempty"`
}

// POST --> mgmtconfig​/v2​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func (service *Service) CreateRule(rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo("POST", path, common.Filter{MicroTenantID: service.microTenantID}, &rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) UpdateRule(policySetID, ruleId string, policySetRule *PolicyRule) (*http.Response, error) {
	// Correct the initialization of Conditions slice with the correct type
	if policySetRule != nil && len(policySetRule.Conditions) == 0 {
		policySetRule.Conditions = []PolicyRuleResourceConditions{}
	} else {
		for i, condition := range policySetRule.Conditions {
			if len(condition.Operands) == 0 {
				policySetRule.Conditions[i].Operands = []PolicyRuleResourceOperands{}
			} else {
				for j, operand := range condition.Operands {
					// Clearing the ID if present, assuming you want to ensure IDs are not sent in updates
					if operand.ID != "" {
						condition.Operands[j].ID = ""
					}
					// If there's more logic to be added for handling Operands, do so here
				}
			}
		}
	}

	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, policySetRule, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// ConvertV1ResponseToV2Request converts a PolicyRuleResource (API v1 response) to a PolicyRule (API v2 request) with aggregated values.
func ConvertV1ResponseToV2Request(v1Response PolicyRuleResource) PolicyRule {
	v2Request := PolicyRule{
		ID:                    v1Response.ID,
		Name:                  v1Response.Name,
		Description:           v1Response.Description,
		Action:                v1Response.Action,
		PolicySetID:           v1Response.PolicySetID,
		Operator:              v1Response.Operator,
		CustomMsg:             v1Response.CustomMsg,
		ZpnIsolationProfileID: v1Response.ZpnIsolationProfileID,
		Conditions:            make([]PolicyRuleResourceConditions, 0),
	}

	for _, condition := range v1Response.Conditions {
		newCondition := PolicyRuleResourceConditions{
			Operator: condition.Operator,
			Operands: make([]PolicyRuleResourceOperands, 0),
		}

		// Use a map to aggregate RHS values by ObjectType
		operandMap := make(map[string][]string)
		entryValuesMap := make(map[string][]OperandsResourceLHSRHSValue)

		for _, operand := range condition.Operands {
			switch operand.ObjectType {
			case "APP", "APP_GROUP", "CONSOLE", "MACHINE_GRP", "LOCATION", "BRANCH_CONNECTOR_GROUP", "EDGE_CONNECTOR_GROUP", "CLIENT_TYPE":
				operandMap[operand.ObjectType] = append(operandMap[operand.ObjectType], operand.RHS)
			case "PLATFORM", "POSTURE", "TRUSTED_NETWORK", "SAML", "SCIM", "SCIM_GROUP", "COUNTRY_CODE":
				entryValuesMap[operand.ObjectType] = append(entryValuesMap[operand.ObjectType], OperandsResourceLHSRHSValue{
					LHS: operand.LHS,
					RHS: operand.RHS,
				})
			}
		}

		// Create operand blocks from the aggregated data
		for objectType, values := range operandMap {
			newCondition.Operands = append(newCondition.Operands, PolicyRuleResourceOperands{
				ObjectType: objectType,
				Values:     values,
			})
		}

		for objectType, entryValues := range entryValuesMap {
			newCondition.Operands = append(newCondition.Operands, PolicyRuleResourceOperands{
				ObjectType:        objectType,
				EntryValuesLHSRHS: entryValues,
			})
		}
		v2Request.Conditions = append(v2Request.Conditions, newCondition)
	}
	return v2Request
}
