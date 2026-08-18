/*
Copyright (c) 2023, Zscaler Inc.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/

package policy_group_set

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

const (
	mgmtConfigV1                     = "/zpa/mgmtconfig/v1/admin/customers/"
	policyGroupSetEndpoint           = "/policyGroupSet"
	policyGroupSetPolicyTypeEndpoint = "/policyGroupSet/policyType"
)

type PolicyRule struct {
	CreationTime         string                           `json:"creationTime,omitempty"`
	Description          string                           `json:"description,omitempty"`
	GroupCriteriaRuleGid string                           `json:"groupCriteriaRuleGid,omitempty"`
	GroupOrder           string                           `json:"groupOrder,omitempty"`
	ID                   string                           `json:"id"`
	ModifiedBy           string                           `json:"modifiedBy,omitempty"`
	ModifiedTime         string                           `json:"modifiedTime,omitempty"`
	Name                 string                           `json:"name,omitempty"`
	PolicyGroupSetGid    string                           `json:"policyGroupSetGid,omitempty"`
	MicrotenantID        string                           `json:"microtenantId,omitempty"`
	MicrotenantName      string                           `json:"microtenantName,omitempty"`
	Type                 string                           `json:"type,omitempty"`
	GroupCriteriaRule    policysetcontrollerv2.PolicyRule `json:"groupCriteriaRule,omitempty"`
}

type PolicyGroupSet struct {
	CreationTime          string   `json:"creationTime,omitempty"`
	CustomPolicyGroupGids []string `json:"customPolicyGroupGids,omitempty"`
	DefaultPolicyGroupGid string   `json:"defaultPolicyGroupGid,omitempty"`
	GlobalPolicyGroupGid  string   `json:"globalPolicyGroupGid,omitempty"`
	ID                    string   `json:"id,omitempty"`
	ModifiedBy            string   `json:"modifiedBy,omitempty"`
	ModifiedTime          string   `json:"modifiedTime,omitempty"`
	Name                  string   `json:"name,omitempty"`
	PolicyType            string   `json:"policyType,omitempty"`
	MicrotenantId         string   `json:"microtenantId,omitempty"`
	MicrotenantName       string   `json:"microtenantName,omitempty"`
}

type PolicyGroupSetSummary struct {
	GroupCountExcludingGlobal string                   `json:"groupCountExcludingGlobal,omitempty"`
	ID                        string                   `json:"id,omitempty"`
	Name                      string                   `json:"name,omitempty"`
	PolicyType                string                   `json:"policyType,omitempty"`
	PolicyGroupSummaryList    []PolicyGroupSummaryList `json:"policyGroupSummaryList,omitempty"`
}

type PolicyGroupSummaryList struct {
	GroupCriteriaCount string `json:"groupCriteriaCount,omitempty"`
	GroupOrder         string `json:"groupOrder,omitempty"`
	Name               string `json:"name,omitempty"`
	ID                 string `json:"id,omitempty"`
	RuleCount          string `json:"ruleCount,omitempty"`
	Type               string `json:"type,omitempty"`
}

type PolicySummaryStats struct {
	DisabledRules     string `json:"disabledRules,omitempty"`
	EnabledRules      string `json:"enabledRules,omitempty"`
	TotalPolicyGroups string `json:"totalPolicyGroups,omitempty"`
	TotalRules        string `json:"totalRules,omitempty"`
}

// GetPolicyGroupSetByID returns a specific Policy Group Set by ID.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/{groupSetId}
// Single (non-paginated) object; supports the microtenant filter.
func GetPolicyGroupSetByID(ctx context.Context, service *zscaler.Service, groupSetID string) (*PolicyGroupSet, *http.Response, error) {
	v := new(PolicyGroupSet)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+policyGroupSetEndpoint, groupSetID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GetAllPolicyGroupSets returns all Policy Group Sets for a customer.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet
// The endpoint returns a bare JSON array (non-paginated).
func GetAllPolicyGroupSets(ctx context.Context, service *zscaler.Service) ([]PolicyGroupSetSummary, *http.Response, error) {
	var v []PolicyGroupSetSummary
	relativeURL := mgmtConfigV1 + service.Client.GetCustomerID() + policyGroupSetEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GetAllPolicyTypeRules returns the paginated rules across groups within a
// Policy Group Set for the given policy type.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/policyType/{policyType}/rules
func GetAllPolicyTypeRules(ctx context.Context, service *zscaler.Service, policyType string) ([]PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s/rules", mgmtConfigV1+service.Client.GetCustomerID()+policyGroupSetPolicyTypeEndpoint, policyType)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

// GetPolicyGroupSetSummary returns the Policy Group Set summary for a customer
// for the given policy type.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/policyType/{policyType}/summary
// Single (non-paginated) object.
func GetPolicyGroupSetSummary(ctx context.Context, service *zscaler.Service, policyType string) (*PolicyGroupSetSummary, *http.Response, error) {
	v := new(PolicyGroupSetSummary)
	relativeURL := fmt.Sprintf("%s/%s/summary", mgmtConfigV1+service.Client.GetCustomerID()+policyGroupSetPolicyTypeEndpoint, policyType)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GetPolicyType returns the Policy Group Set for a customer for the given policy type.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/policyType/{policyType}
// Single (non-paginated) object.
func GetPolicyType(ctx context.Context, service *zscaler.Service, policyType string) (*PolicyGroupSet, *http.Response, error) {
	v := new(PolicyGroupSet)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfigV1+service.Client.GetCustomerID()+policyGroupSetPolicyTypeEndpoint, policyType)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GetPolicyTypeSummaryStats returns the summary stats for groups and rules
// within a Policy Group Set for the given policy type.
// GET /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/policyType/{policyType}/summaryStats
// Single (non-paginated) object.
func GetPolicyTypeSummaryStats(ctx context.Context, service *zscaler.Service, policyType string) (*PolicySummaryStats, *http.Response, error) {
	v := new(PolicySummaryStats)
	relativeURL := fmt.Sprintf("%s/%s/summaryStats", mgmtConfigV1+service.Client.GetCustomerID()+policyGroupSetPolicyTypeEndpoint, policyType)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}
