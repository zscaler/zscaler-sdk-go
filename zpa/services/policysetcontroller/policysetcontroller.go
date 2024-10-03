package policysetcontroller

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/appconnectorgroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/servergroup"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/serviceedgegroup"
)

const (
	mgmtConfig = "/mgmtconfig/v1/admin/customers/"
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

type PolicyRule struct {
	Action                   string                                `json:"action,omitempty"`
	ActionID                 string                                `json:"actionId,omitempty"`
	BypassDefaultRule        bool                                  `json:"bypassDefaultRule"`
	CreationTime             string                                `json:"creationTime,omitempty"`
	CustomMsg                string                                `json:"customMsg,omitempty"`
	DefaultRule              bool                                  `json:"defaultRule,omitempty"`
	DefaultRuleName          string                                `json:"defaultRuleName,omitempty"`
	Description              string                                `json:"description,omitempty"`
	ID                       string                                `json:"id,omitempty"`
	ModifiedBy               string                                `json:"modifiedBy,omitempty"`
	ModifiedTime             string                                `json:"modifiedTime,omitempty"`
	Name                     string                                `json:"name,omitempty"`
	Operator                 string                                `json:"operator,omitempty"`
	PolicySetID              string                                `json:"policySetId"`
	PolicyType               string                                `json:"policyType,omitempty"`
	Priority                 string                                `json:"priority,omitempty"`
	ReauthDefaultRule        bool                                  `json:"reauthDefaultRule"`
	ReauthIdleTimeout        string                                `json:"reauthIdleTimeout,omitempty"`
	ReauthTimeout            string                                `json:"reauthTimeout,omitempty"`
	RuleOrder                string                                `json:"ruleOrder"`
	LSSDefaultRule           bool                                  `json:"lssDefaultRule"`
	ZpnCbiProfileID          string                                `json:"zpnCbiProfileId,omitempty"`
	ZpnIsolationProfileID    string                                `json:"zpnIsolationProfileId,omitempty"`
	ZpnInspectionProfileID   string                                `json:"zpnInspectionProfileId,omitempty"`
	ZpnInspectionProfileName string                                `json:"zpnInspectionProfileName,omitempty"`
	MicroTenantID            string                                `json:"microtenantId,omitempty"`
	MicroTenantName          string                                `json:"microtenantName,omitempty"`
	Conditions               []Conditions                          `json:"conditions"`
	AppServerGroups          []servergroup.ServerGroup             `json:"appServerGroups"`
	AppConnectorGroups       []appconnectorgroup.AppConnectorGroup `json:"appConnectorGroups"`
	ServiceEdgeGroups        []serviceedgegroup.ServiceEdgeGroup   `json:"serviceEdgeGroups"`
	Credential               *Credential                           `json:"credential,omitempty"`
	PrivilegedCapabilities   PrivilegedCapabilities                `json:"privilegedCapabilities,omitempty"`
}

type Conditions struct {
	CreationTime  string     `json:"creationTime,omitempty"`
	ID            string     `json:"id,omitempty"`
	ModifiedBy    string     `json:"modifiedBy,omitempty"`
	ModifiedTime  string     `json:"modifiedTime,omitempty"`
	Negated       bool       `json:"negated"`
	Operands      []Operands `json:"operands"`
	Operator      string     `json:"operator,omitempty"`
	MicroTenantID string     `json:"microtenantId,omitempty"`
}

type Operands struct {
	CreationTime  string `json:"creationTime,omitempty"`
	ID            string `json:"id,omitempty"`
	IdpID         string `json:"idpId,omitempty"`
	LHS           string `json:"lhs,omitempty"`
	ModifiedBy    string `json:"modifiedBy,omitempty"`
	ModifiedTime  string `json:"modifiedTime,omitempty"`
	Name          string `json:"name,omitempty"`
	ObjectType    string `json:"objectType,omitempty"`
	RHS           string `json:"rhs,omitempty"`
	MicroTenantID string `json:"microtenantId,omitempty"`
}

/*
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
*/
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

type Count struct {
	Count string `json:"count"`
}

func GetByPolicyType(service *services.Service, policyType string) (*PolicySet, *http.Response, error) {
	v := new(PolicySet)
	relativeURL := fmt.Sprintf(mgmtConfig + service.Client.Config.CustomerID + "/policySet/policyType/" + policyType)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetPolicyRule(service *services.Service, policySetID, ruleId string) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	url := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("GET", url, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// POST --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule
func CreateRule(service *services.Service, rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule", rule.PolicySetID)
	resp, err := service.Client.NewRequestDo("POST", path, common.Filter{MicroTenantID: service.MicroTenantID()}, rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func UpdateRule(service *services.Service, policySetID, ruleId string, policySetRule *PolicyRule) (*http.Response, error) {
	if policySetRule != nil && len(policySetRule.Conditions) == 0 {
		policySetRule.Conditions = []Conditions{}
	} else {
		for i, condtion := range policySetRule.Conditions {
			if len(condtion.Operands) == 0 {
				policySetRule.Conditions[i].Operands = []Operands{}
			} else {
				for i, operand := range condtion.Operands {
					if operand.Name != "" {
						condtion.Operands[i].Name = ""
					}
				}
			}
		}
	}
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, policySetRule, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// DELETE --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule​/{ruleId}
func Delete(service *services.Service, policySetID, ruleId string) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s", policySetID, ruleId)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetByNameAndType(service *services.Service, policyType, ruleName string) (*PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](service.Client, relativeURL, common.Filter{Search: ruleName, MicroTenantID: service.MicroTenantID()})
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

func GetByNameAndTypes(service *services.Service, policyTypes []string, ruleName string) (p *PolicyRule, resp *http.Response, err error) {
	for _, policyType := range policyTypes {
		p, resp, err = GetByNameAndType(service, policyType, ruleName)
		if err == nil {
			return p, resp, nil
		}
	}
	return nil, nil, err
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySetId}/rule/{ruleId}/reorder/{newOrder}
func Reorder(service *services.Service, policySetID, ruleId string, order int) (*http.Response, error) {
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/rule/%s/reorder/%d", policySetID, ruleId, order)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policySet/{policySet}/reorder
// ruleIdOrders is a map[ruleID]Order
func BulkReorder(service *services.Service, policySetType string, ruleIdToOrder map[string]int) (*http.Response, error) {
	policySet, resp, err := GetByPolicyType(service, policySetType)
	if err != nil {
		return resp, err
	}
	all, resp, err := GetAllByType(service, policySetType)
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
	path := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/%s/reorder", policySet.ID)
	ruleIdsOrdered := []string{}
	for _, r := range all {
		ruleIdsOrdered = append(ruleIdsOrdered, r.ID)
	}

	// Create a new PUT request
	resp, err = service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, ruleIdsOrdered, nil)
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

func GetAllByType(service *services.Service, policyType string) ([]PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfig+service.Client.Config.CustomerID+"/policySet/rules/policyType/%s", policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
