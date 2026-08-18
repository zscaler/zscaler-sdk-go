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

package policy_group_rule

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontrollerv2"
)

const (
	mgmtConfigV1 = "/zpa/mgmtconfig/v1/admin/customers/"
)

var ruleMutex sync.Mutex

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

// Get returns a policy rule within a policy group.
// GET .../policyGroupSet/{groupSetId}/group/{groupId}/rule/{ruleId}
func Get(ctx context.Context, service *zscaler.Service, groupSetID, groupID, ruleID string) (*PolicyRule, *http.Response, error) {
	v := new(PolicyRule)
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s/rule/%s", groupSetID, groupID, ruleID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// GetPolicyGroupRuleByName returns the first rule whose name matches within the
// given policy group.
func GetPolicyGroupRuleByName(ctx context.Context, service *zscaler.Service, groupSetID, groupID, name string) (*PolicyRule, *http.Response, error) {
	list, resp, err := GetAll(ctx, service, groupSetID, groupID)
	if err != nil {
		return nil, resp, err
	}
	for _, rule := range list {
		if strings.EqualFold(rule.Name, name) {
			return &rule, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy group rule named '%s' found", name)
}

// Create adds a new policy rule for a given policy group.
// POST .../policyGroupSet/{groupSetId}/group/{groupId}/rule
func Create(ctx context.Context, service *zscaler.Service, groupSetID, groupID string, rule *PolicyRule) (*PolicyRule, *http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	v := new(PolicyRule)
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s/rule", groupSetID, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, rule, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// Delete deletes a policy rule within a policy group.
// DELETE .../policyGroupSet/{groupSetId}/group/{groupId}/rule/{ruleId}
func Delete(ctx context.Context, service *zscaler.Service, groupSetID, groupID, ruleID string) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s/rule/%s", groupSetID, groupID, ruleID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Reorder updates the rule order of a rule within a policy group.
// PUT .../policyGroupSet/{groupSetId}/group/{groupId}/rule/{ruleId}/reorder/{newOrder}
func Reorder(ctx context.Context, service *zscaler.Service, groupSetID, groupID, ruleID string, newOrder int) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s/rule/%s/reorder/%d", groupSetID, groupID, ruleID, newOrder)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAll returns all policy group rules within a policy group. The endpoint
// supports advanced search and pagination, handled by the ReadAll paginator.
// GET .../policyGroupSet/{groupSetId}/group/{groupId}/rule
func GetAll(ctx context.Context, service *zscaler.Service, groupSetID, groupID string) ([]PolicyRule, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s/rule", groupSetID, groupID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyRule](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
