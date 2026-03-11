package tag_key_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig               = "/zpa/mgmtconfig/v1/admin/customers/"
	tagKeyPath               = "/tagKey"
	tagKeySearchPath         = "/tagKey/search"
	bulkUpdateStatusPath     = "/tagKey/bulkUpdateStatus"
)

type TagKey struct {
	ID              string     `json:"id,omitempty"`
	CustomerID      string     `json:"customerId,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	Enabled         bool       `json:"enabled"`
	NamespaceID     string     `json:"namespaceId,omitempty"`
	Origin          string     `json:"origin,omitempty"`
	Type            string     `json:"type,omitempty"`
	MicroTenantID   string     `json:"microtenantId,omitempty"`
	MicroTenantName string     `json:"microtenantName,omitempty"`
	SkipAudit       bool       `json:"skipAudit,omitempty"`
	TagValues       []TagValue `json:"tagValues"`
}

type TagValue struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type BulkUpdateStatusRequest struct {
	Enabled   bool     `json:"enabled"`
	TagKeyIDs []string `json:"tagKeyIds,omitempty"`
}

func namespacePath(customerID, namespaceID string) string {
	return mgmtConfig + customerID + "/namespace/" + namespaceID
}

func Get(ctx context.Context, service *zscaler.Service, namespaceID, tagKeyID string) (*TagKey, *http.Response, error) {
	v := new(TagKey)
	relativeURL := fmt.Sprintf("%s/%s", namespacePath(service.Client.GetCustomerID(), namespaceID)+tagKeyPath, tagKeyID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, namespaceID, tagKeyName string) (*TagKey, *http.Response, error) {
	relativeURL := namespacePath(service.Client.GetCustomerID(), namespaceID) + tagKeySearchPath
	searchRequest := common.SearchRequest{
		FilterBy: &common.SearchFilterBy{
			FilterGroups: []common.SearchFilterGroup{
				{
					Filters: []common.SearchFilterItem{
						{
							CommaSepValues: tagKeyName,
							FilterName:     "name",
							Operator:       "EQ",
							Value:          tagKeyName,
							Values:         []string{tagKeyName},
						},
					},
					Operator: "AND",
				},
			},
			Operator: "AND",
		},
		SortBy: &common.SearchSortBy{
			SortName:  "name",
			SortOrder: "ASC",
		},
	}
	list, resp, err := common.GetAllPagesGenericWithPostSearch[TagKey](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, tagKey := range list {
		if strings.EqualFold(tagKey.Name, tagKeyName) {
			return &tagKey, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no tag key named '%s' was found", tagKeyName)
}

func Create(ctx context.Context, service *zscaler.Service, namespaceID string, tagKey TagKey) (*TagKey, *http.Response, error) {
	v := new(TagKey)
	// API requires tagValues in the payload even when empty; ensure it's never nil
	payload := tagKey
	if payload.TagValues == nil {
		payload.TagValues = []TagValue{}
	}
	resp, err := service.Client.NewRequestDo(ctx, "POST", namespacePath(service.Client.GetCustomerID(), namespaceID)+tagKeyPath, common.Filter{MicroTenantID: service.MicroTenantID()}, payload, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, namespaceID, tagKeyID string, tagKey *TagKey) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", namespacePath(service.Client.GetCustomerID(), namespaceID)+tagKeyPath, tagKeyID)
	// API requires tagValues in the payload even when empty; ensure it's never nil
	payload := *tagKey
	if payload.TagValues == nil {
		payload.TagValues = []TagValue{}
	}
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, &payload, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, namespaceID, tagKeyID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", namespacePath(service.Client.GetCustomerID(), namespaceID)+tagKeyPath, tagKeyID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service, namespaceID string) ([]TagKey, *http.Response, error) {
	relativeURL := namespacePath(service.Client.GetCustomerID(), namespaceID) + tagKeySearchPath
	searchRequest := common.SearchRequest{
		SortBy: &common.SearchSortBy{
			SortName:  "name",
			SortOrder: "ASC",
		},
	}
	list, resp, err := common.GetAllPagesGenericWithPostSearch[TagKey](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func BulkUpdateStatus(ctx context.Context, service *zscaler.Service, namespaceID string, bulkUpdate BulkUpdateStatusRequest) (*http.Response, error) {
	path := namespacePath(service.Client.GetCustomerID(), namespaceID) + bulkUpdateStatusPath
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, bulkUpdate, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
