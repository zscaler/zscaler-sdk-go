package policysetcontrollerv2

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfigV2 = "/mgmtconfig/v2/admin/customers/"
)

type PolicyRuleList struct {
	CreationTime       string               `json:"creationTime,omitempty"`
	Description        string               `json:"description,omitempty"`
	Enabled            bool                 `json:"enabled"`
	ID                 string               `json:"id,omitempty"`
	ModifiedBy         string               `json:"modifiedBy,omitempty"`
	ModifiedTime       string               `json:"modifiedTime,omitempty"`
	Name               string               `json:"name,omitempty"`
	Sorted             bool                 `json:"sorted"`
	PolicyType         string               `json:"policyType,omitempty"`
	MicroTenantID      string               `json:"microtenantId,omitempty"`
	MicroTenantName    string               `json:"microtenantName,omitempty"`
	PolicyRuleResource []PolicyRuleResource `json:"rules"`
}

type PolicyRuleResource struct {
	ID                       string                   `json:"id,omitempty"`
	Name                     string                   `json:"name,omitempty"`
	Description              string                   `json:"description,omitempty"`
	Action                   string                   `json:"action,omitempty"`
	ActionID                 string                   `json:"actionId,omitempty"`
	CreationTime             string                   `json:"creationTime,omitempty"`
	CustomMsg                string                   `json:"customMsg,omitempty"`
	DefaultRule              bool                     `json:"defaultRule,omitempty"`
	DefaultRuleName          string                   `json:"defaultRuleName,omitempty"`
	ModifiedBy               string                   `json:"modifiedBy,omitempty"`
	ModifiedTime             string                   `json:"modifiedTime,omitempty"`
	Operator                 string                   `json:"operator,omitempty"`
	PolicySetID              string                   `json:"policySetId"`
	PolicyType               string                   `json:"policyType,omitempty"`
	Priority                 string                   `json:"priority,omitempty"`
	ReauthIdleTimeout        string                   `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string                   `json:"reauthTimeout,omitempty"`
	RuleOrder                string                   `json:"ruleOrder"`
	ZpnIsolationProfileID    string                   `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string                   `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                   `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string                   `json:"microtenantId,omitempty"`
	MicroTenantName          string                   `json:"microtenantName,omitempty"`
	ConditionSetResourceV2   []ConditionSetResourceV2 `json:"conditions"`
	AppServerGroupsV2        []AppServerGroupsV2      `json:"appServerGroups"`
	AppConnectorGroupsV2     []AppConnectorGroupsV2   `json:"appConnectorGroups"`
	ServiceEdgeGroupsV2      []ServiceEdgeGroupsV2    `json:"serviceEdgeGroups"`
}

type ConditionSetResourceV2 struct {
	ID                string              `json:"id,omitempty"`
	CreationTime      string              `json:"creationTime,omitempty"`
	ModifiedBy        string              `json:"modifiedBy,omitempty"`
	ModifiedTime      string              `json:"modifiedTime,omitempty"`
	Negated           bool                `json:"negated,omitempty"`
	Operator          string              `json:"operator,omitempty"`
	MicroTenantID     string              `json:"microtenantId,omitempty"`
	OperandResourceV2 []OperandResourceV2 `json:"operands"`
}

type OperandResourceV2 struct {
	ID                          string                        `json:"id,omitempty"`
	CreationTime                string                        `json:"creationTime,omitempty"`
	ModifiedBy                  string                        `json:"modifiedBy,omitempty"`
	ModifiedTime                string                        `json:"modifiedTime,omitempty"`
	IdpID                       string                        `json:"idpId,omitempty"`
	ObjectType                  string                        `json:"objectType"`
	Values                      []string                      `json:"values"`
	OperandsResourceLHSRHSValue []OperandsResourceLHSRHSValue `json:"entryValues"`
}

type OperandsResourceLHSRHSValue struct {
	LHS string `json:"lhs,omitempty"`
	RHS string `json:"rhs,omitempty"`
}

type AppServerGroupsV2 struct {
	ID string `json:"id,omitempty"`
}

type AppConnectorGroupsV2 struct {
	ID string `json:"id,omitempty"`
}

type ServiceEdgeGroupsV2 struct {
	ID string `json:"id,omitempty"`
}

func (service *Service) GetByPolicyTypeV2(policyType string) (*PolicyRuleList, *http.Response, error) {
	v := new(PolicyRuleList)
	relativeURL := fmt.Sprintf(mgmtConfigV2 + service.Client.Config.CustomerID + "/policySet/policyType/" + policyType)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule/{ruleId}
func (service *Service) GetPolicyRuleV2(policySetID, ruleId string) (*PolicyRuleResource, *http.Response, error) {
	v := new(PolicyRuleResource)
	url := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("GET", url, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// POST --> mgmtconfig​/v2​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func (service *Service) CreateRuleV2(rule *PolicyRuleResource) (*PolicyRuleResource, *http.Response, error) {
	v := new(PolicyRuleResource)
	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo("POST", path, common.Filter{MicroTenantID: service.microTenantID}, &rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) UpdateV2(policySetID, ruleId string, policySetRule *PolicyRuleResource) (*http.Response, error) {
	if policySetRule != nil && len(policySetRule.ConditionSetResourceV2) == 0 {
		policySetRule.ConditionSetResourceV2 = []ConditionSetResourceV2{}
	} else {
		for i, condtion := range policySetRule.ConditionSetResourceV2 {
			if len(condtion.OperandResourceV2) == 0 {
				policySetRule.ConditionSetResourceV2[i].OperandResourceV2 = []OperandResourceV2{}
			} else {
				for i, operand := range condtion.OperandResourceV2 {
					if operand.ID != "" {
						condtion.OperandResourceV2[i].ID = ""
					}
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

// DELETE --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func (service *Service) DeleteV2(policySetID, ruleId string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetByNameAndTypeV2(policyType, ruleName string) (*PolicyRuleResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRuleResource](service.Client, relativeURL, common.Filter{Search: ruleName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}

	for _, p := range list {
		if strings.EqualFold(ruleName, p.Name) {
			return &p, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy rule named :%s found", ruleName)
}

func (service *Service) GetByNameAndTypesV2(policyTypes []string, ruleName string) (p *PolicyRuleResource, resp *http.Response, err error) {
	for _, policyType := range policyTypes {
		p, resp, err = service.GetByNameAndTypeV2(policyType, ruleName)
		if err != nil {
			continue
		} else {
			return
		}
	}
	return
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func (service *Service) ReorderV2(policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySet}/reorder
// ruleIdOrders is a map[ruleID]Order
func (service *Service) BulkReorderV2(policySetType string, ruleIdToOrder map[string]int) (*http.Response, error) {
	policySet, resp, err := service.GetByPolicyTypeV2(policySetType)
	if err != nil {
		return resp, err
	}
	all, resp, err := service.GetAllByTypeV2(policySetType)
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
	// Construct the URL path
	path := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/%s/reorder", policySet.ID)
	ruleIdsOrdered := []string{}
	for _, r := range all {
		ruleIdsOrdered = append(ruleIdsOrdered, r.ID)
	}

	// Create a new PUT request
	resp, err = service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, ruleIdsOrdered, nil)
	if err != nil {
		return nil, err
	}

	// Check for non-2xx status code and log response body for debugging
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close() // Ensure the body is always closed
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// Handle the error of reading the body (optional)
			log.Printf("Error reading response body: %s\n", err.Error())
		}
		log.Printf("Error response from API: %s\n", string(bodyBytes))
		return resp, fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}

func (service *Service) GetAllByTypeV2(policyType string) ([]PolicyRuleResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV2+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRuleResource](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
