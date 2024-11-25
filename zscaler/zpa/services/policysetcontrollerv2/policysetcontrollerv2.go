package policysetcontrollerv2

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/servergroup"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/serviceedgegroup"
)

const (
	mgmtConfigV1 = "/zpa/mgmtconfig/v1/admin/customers/"
	mgmtConfigV2 = "/zpa/mgmtconfig/v2/admin/customers/"
)

var ruleMutex sync.Mutex

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
	MicroTenantID   string       `json:"microtenantId"`
	MicroTenantName string       `json:"microtenantName,omitempty"`
	Rules           []PolicyRule `json:"rules"`
}

// ######################################################################################################
// ########################################## API V1 Structure ##########################################
// ################################### Used to process the API Response #################################
// ######################################################################################################
type PolicyRuleResource struct {
	ID                           string                                `json:"id,omitempty"`
	Name                         string                                `json:"name,omitempty"`
	Description                  string                                `json:"description,omitempty"`
	Disabled                     string                                `json:"disabled,omitempty"`
	ExtranetEnabled              bool                                  `json:"extranetEnabled,omitempty"`
	Action                       string                                `json:"action,omitempty"`
	ActionID                     string                                `json:"actionId,omitempty"`
	CreationTime                 string                                `json:"creationTime,omitempty"`
	ModifiedBy                   string                                `json:"modifiedBy,omitempty"`
	ModifiedTime                 string                                `json:"modifiedTime,omitempty"`
	AuditMessage                 string                                `json:"auditMessage,omitempty"`
	CustomMsg                    string                                `json:"customMsg,omitempty"`
	DefaultRule                  bool                                  `json:"defaultRule,omitempty"`
	DefaultRuleName              string                                `json:"defaultRuleName,omitempty"`
	Operator                     string                                `json:"operator,omitempty"`
	PolicySetID                  string                                `json:"policySetId,omitempty"`
	PolicyType                   string                                `json:"policyType,omitempty"`
	Priority                     string                                `json:"priority,omitempty"`
	ReauthIdleTimeout            string                                `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout                string                                `json:"reauthTimeout,omitempty"`
	RuleOrder                    string                                `json:"ruleOrder,omitempty"`
	ZpnIsolationProfileID        string                                `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID       string                                `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName     string                                `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID                string                                `json:"microtenantId"`
	MicroTenantName              string                                `json:"microtenantName,omitempty"`
	Conditions                   []PolicyRuleResourceConditions        `json:"conditions"`
	AppServerGroups              []servergroup.ServerGroup             `json:"appServerGroups"`
	AppConnectorGroups           []appconnectorgroup.AppConnectorGroup `json:"appConnectorGroups"`
	ServiceEdgeGroups            []serviceedgegroup.ServiceEdgeGroup   `json:"serviceEdgeGroups"`
	Credential                   *Credential                           `json:"credential,omitempty"`
	PrivilegedCapabilities       PrivilegedCapabilities                `json:"privilegedCapabilities,omitempty"`
	ExtranetDTO                  ExtranetDTO                           `json:"extranetDTO,omitempty"`
	PrivilegedPortalCapabilities PrivilegedPortalCapabilities          `json:"privilegedPortalCapabilities,omitempty"`
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
	ID                           string                                `json:"id,omitempty"`
	Name                         string                                `json:"name,omitempty"`
	Action                       string                                `json:"action,omitempty"`
	ActionID                     string                                `json:"actionId,omitempty"`
	CustomMsg                    string                                `json:"customMsg,omitempty"`
	Description                  string                                `json:"description,omitempty"`
	Disabled                     string                                `json:"disabled,omitempty"`
	ExtranetEnabled              bool                                  `json:"extranetEnabled,omitempty"`
	CreationTime                 string                                `json:"creationTime,omitempty"`
	ModifiedBy                   string                                `json:"modifiedBy,omitempty"`
	ModifiedTime                 string                                `json:"modifiedTime,omitempty"`
	Operator                     string                                `json:"operator,omitempty"`
	PolicySetID                  string                                `json:"policySetId,omitempty"`
	PolicyType                   string                                `json:"policyType,omitempty"`
	Priority                     string                                `json:"priority,omitempty"`
	ReauthIdleTimeout            string                                `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout                string                                `json:"reauthTimeout,omitempty"`
	RuleOrder                    string                                `json:"ruleOrder,omitempty"`
	ZpnIsolationProfileID        string                                `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID       string                                `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName     string                                `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID                string                                `json:"microtenantId"`
	MicroTenantName              string                                `json:"microtenantName,omitempty"`
	AppServerGroups              []servergroup.ServerGroup             `json:"appServerGroups"`
	AppConnectorGroups           []appconnectorgroup.AppConnectorGroup `json:"appConnectorGroups"`
	ServiceEdgeGroups            []serviceedgegroup.ServiceEdgeGroup   `json:"serviceEdgeGroups"`
	Conditions                   []PolicyRuleResourceConditions        `json:"conditions,omitempty"`
	Credential                   *Credential                           `json:"credential,omitempty"`
	PrivilegedCapabilities       PrivilegedCapabilities                `json:"privilegedCapabilities,omitempty"`
	ExtranetDTO                  ExtranetDTO                           `json:"extranetDTO,omitempty"`
	PrivilegedPortalCapabilities PrivilegedPortalCapabilities          `json:"privilegedPortalCapabilities,omitempty"`
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

type Credential struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type PrivilegedCapabilities struct {
	ID            string   `json:"id"`
	CreationTime  string   `json:"creationTime,omitempty"`
	ModifiedBy    string   `json:"modifiedBy,omitempty"`
	ModifiedTime  string   `json:"modifiedTime,omitempty"`
	MicroTenantID string   `json:"microtenantId"`
	Capabilities  []string `json:"capabilities,omitempty"`
}

type ExtranetDTO struct {
	LocationDTO      []LocationDTO      `json:"locationDTO,omitempty"`
	LocationGroupDTO []LocationGroupDTO `json:"locationGroupDTO,omitempty"`
	ZiaErName        string             `json:"ziaErName,omitempty"`
	ZpnErID          string             `json:"zpnErId,omitempty"`
}

type LocationDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type LocationGroupDTO struct {
	ID           string        `json:"id,omitempty"`
	Name         string        `json:"name,omitempty"`
	ZiaLocations []ZiaLocation `json:"ziaLocations,omitempty"`
}

type ZiaLocation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type PrivilegedPortalCapabilities struct {
	Capabilities  []string `json:"capabilities,omitempty"`
	MicroTenantID string   `json:"microtenantId,omitempty"`
}

func GetByPolicyType(ctx context.Context, service *zscaler.Service, policyType string) (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfigV1 + service.Client.GetCustomerID() + "/policySet/policyType/" + policyType)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule/{ruleId}
func GetPolicyRule(ctx context.Context, service *zscaler.Service, policySetID, ruleId string) (*PolicyRuleResource, *http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	v := new(PolicyRuleResource)
	url := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo(ctx, "GET", url, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// POST --> mgmtconfig​/v2​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func CreateRule(ctx context.Context, service *zscaler.Service, rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	v := new(PolicyRule)
	path := fmt.Sprintf(mgmtConfigV2+service.Client.GetCustomerID()+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, common.Filter{MicroTenantID: service.MicroTenantID()}, rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func UpdateRule(ctx context.Context, service *zscaler.Service, policySetID, ruleId string, policySetRule *PolicyRule) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

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

	path := fmt.Sprintf(mgmtConfigV2+service.Client.GetCustomerID()+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, policySetRule, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DELETE --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func Delete(ctx context.Context, service *zscaler.Service, policySetID, ruleId string) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetByNameAndType(ctx context.Context, service *zscaler.Service, policyType, ruleName string) (*PolicyRuleResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRuleResource](ctx, service.Client, relativeURL, common.Filter{Search: ruleName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}

	for _, p := range list {
		if strings.EqualFold(ruleName, p.Name) {
			return &p, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy rule named '%s' found", ruleName)
}

func GetByNameAndTypes(ctx context.Context, service *zscaler.Service, policyTypes []string, ruleName string) (*PolicyRuleResource, *http.Response, error) {
	for _, policyType := range policyTypes {
		p, resp, err := GetByNameAndType(context.Background(), service, policyType, ruleName)
		if err == nil {
			return p, resp, nil
		}
	}
	return nil, nil, fmt.Errorf("no policy rule named '%s' found in any policy type", ruleName)
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func Reorder(ctx context.Context, service *zscaler.Service, policySetID, ruleId string, order int) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySet}/reorder
// ruleIdOrders is a map[ruleID]Order
func BulkReorder(ctx context.Context, service *zscaler.Service, policySetType string, ruleIdToOrder map[string]int) (*http.Response, error) {
	policySet, resp, err := GetByPolicyType(context.Background(), service, policySetType)
	if err != nil {
		return resp, err
	}
	all, resp, err := GetAllByType(context.Background(), service, policySetType)
	if err != nil {
		return resp, err
	}
	sort.SliceStable(all, func(i, j int) bool {
		ruleIDi := all[i].ID
		ruleIDj := all[j].ID

		// Check if ruleIDi and ruleIDj exist in the ruleIdToOrder map
		orderi, existsi := ruleIdToOrder[ruleIDi]
		orderj, existsj := ruleIdToOrder[ruleIDj]

		// If both rules exist in the map, compare their orders
		if existsi && existsj {
			return orderi <= orderj
		}

		// If only one of the rules exists in the map, prioritize it
		if existsi {
			return true
		} else if existsj {
			return false
		}

		// If neither rule exists in the map, maintain their relative order
		return i <= j
	})
	// Prepare the rule order list, ensuring Default_Rule is placed at the end
	var ruleIdsOrdered []string
	var defaultRuleID string

	for _, rule := range all {
		// Check if this is the Default_Rule
		if rule.Name == "Default_Rule" {
			defaultRuleID = rule.ID
			continue
		}
		// Add all other rules to the ordered list
		ruleIdsOrdered = append(ruleIdsOrdered, rule.ID)
	}

	// Append the Default_Rule to the end if it exists
	if defaultRuleID != "" {
		ruleIdsOrdered = append(ruleIdsOrdered, defaultRuleID)
	}

	// Log the final order for debugging
	log.Printf("Final rule order: %v", ruleIdsOrdered)

	// Construct the URL path
	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/%s/reorder", policySet.ID)

	// Create a new PUT request
	resp, err = service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, ruleIdsOrdered, nil)
	if err != nil {
		return nil, err
	}

	// Check for non-2xx status code and log response body for debugging
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close() // Ensure the body is always closed
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body: %s\n", err.Error())
		}
		log.Printf("Error response from API: %s\n", string(bodyBytes))
		return resp, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func GetAllByType(ctx context.Context, service *zscaler.Service, policyType string) ([]PolicyRuleResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRuleResource](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
