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

package policy_group

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

type PolicyGroupResource struct {
	Name              string                           `json:"name,omitempty"`
	Description       string                           `json:"description,omitempty"`
	PolicyGroupSetGid string                           `json:"policyGroupSetGid,omitempty"`
	MicrotenantID     string                           `json:"microtenantId,omitempty"`
	MicrotenantName   string                           `json:"microtenantName,omitempty"`
	Type              string                           `json:"type,omitempty"`
	GroupCriteriaRule policysetcontrollerv2.PolicyRule `json:"groupCriteria,omitempty"`
}

// GET --> mgmtconfig​/v1​/admin​/customers​/{customerId}​/policySet​/{policySetId}​/rule/{ruleId}
func GetPolicyRule(ctx context.Context, service *zscaler.Service, groupSetID, groupID string) (*PolicyGroupResource, *http.Response, error) {
	v := new(PolicyGroupResource)
	url := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s", groupSetID, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", url, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetPolicyGroupByName(ctx context.Context, service *zscaler.Service, groupSetID, name string) (*PolicyGroupResource, *http.Response, error) {
	list, resp, err := GetAllPolicyGroups(ctx, service, groupSetID)
	if err != nil {
		return nil, resp, err
	}
	for _, group := range list {
		if strings.EqualFold(group.Name, name) {
			return &group, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no policy group named '%s' found", name)
}

func GetMicrotenantByName(ctx context.Context, service *zscaler.Service, groupSetID, groupName string) (*PolicyGroupResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/search", groupSetID)
	searchRequest := common.SearchRequest{
		FilterBy: &common.SearchFilterBy{
			FilterGroups: []common.SearchFilterGroup{
				{
					Filters: []common.SearchFilterItem{
						{
							FilterName: "name",
							Operator:   "EQ",
							Value:      groupName,
						},
					},
					Operator: "AND",
				},
			},
			Operator: "AND",
		},
	}
	list, resp, err := common.GetAllPagesGenericWithPostSearch[PolicyGroupResource](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, ns := range list {
		if strings.EqualFold(ns.Name, groupName) {
			return &ns, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no group named '%s' was found", groupName)
}

func CreateRule(ctx context.Context, service *zscaler.Service, group *PolicyGroupResource) (*PolicyGroupResource, *http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	v := new(PolicyGroupResource)
	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/rule", group.PolicyGroupSetGid)
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, common.Filter{MicroTenantID: service.MicroTenantID()}, group, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

// PUT --> /mgmtconfig/v1/admin/customers/{customerId}/policyGroupSet/{groupSetId}/group/{groupId}
func Update(ctx context.Context, service *zscaler.Service, groupSetID, groupID string, body *PolicyGroupResource) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s", groupSetID, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, body, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Delete(ctx context.Context, service *zscaler.Service, groupSetID, groupID string) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/%s", groupSetID, groupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ReorderGroup(ctx context.Context, service *zscaler.Service, groupSetID, groupID string, order int) (*http.Response, error) {
	ruleMutex.Lock()
	defer ruleMutex.Unlock()

	path := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/rule/%s/reorder/%d", groupSetID, groupID, order)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAllPolicyGroups(ctx context.Context, service *zscaler.Service, groupSetID string) ([]PolicyGroupResource, *http.Response, error) {
	relativeURL := fmt.Sprintf(mgmtConfigV1+service.Client.GetCustomerID()+"/policyGroupSet/%s/group/all", groupSetID)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PolicyGroupResource](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
