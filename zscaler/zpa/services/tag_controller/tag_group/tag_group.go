package tag_group

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	tagGroupEndpoint       = "/tagGroup"
	tagGroupSearchEndpoint = "/tagGroup/search"
)

type TagGroup struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
	Tags            []Tag  `json:"tags"`
}

type Tag struct {
	Namespace *TagNamespace `json:"namespace,omitempty"`
	Origin    string        `json:"origin,omitempty"`
	TagKey    *TagKey       `json:"tagKey,omitempty"`
	TagValue  *TagValue     `json:"tagValue,omitempty"`
}

type TagNamespace struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
}

type TagKey struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled"`
}

type TagValue struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, tagGroupID string) (*TagGroup, *http.Response, error) {
	v := new(TagGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+tagGroupEndpoint, tagGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, tagGroupName string) (*TagGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + tagGroupSearchEndpoint
	searchRequest := common.SearchRequest{
		FilterBy: &common.SearchFilterBy{
			FilterGroups: []common.SearchFilterGroup{
				{
					Filters: []common.SearchFilterItem{
						{
							CommaSepValues: tagGroupName,
							FilterName:     "name",
							Operator:       "EQ",
							Value:          tagGroupName,
							Values:         []string{tagGroupName},
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
	list, resp, err := common.GetAllPagesGenericWithPostSearch[TagGroup](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, tagGroup := range list {
		if strings.EqualFold(tagGroup.Name, tagGroupName) {
			return &tagGroup, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no tag group named '%s' was found", tagGroupName)
}

func Create(ctx context.Context, service *zscaler.Service, tagGroup TagGroup) (*TagGroup, *http.Response, error) {
	v := new(TagGroup)
	// API requires tags in the payload even when empty; ensure it's never nil
	payload := tagGroup
	if payload.Tags == nil {
		payload.Tags = []Tag{}
	}
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+tagGroupEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, payload, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, tagGroupID string, tagGroup *TagGroup) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+tagGroupEndpoint, tagGroupID)
	// API requires tags in the payload even when empty; ensure it's never nil
	payload := *tagGroup
	if payload.Tags == nil {
		payload.Tags = []Tag{}
	}
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, &payload, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, tagGroupID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+tagGroupEndpoint, tagGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]TagGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + tagGroupSearchEndpoint
	searchRequest := common.SearchRequest{
		SortBy: &common.SearchSortBy{
			SortName:  "name",
			SortOrder: "ASC",
		},
	}
	list, resp, err := common.GetAllPagesGenericWithPostSearch[TagGroup](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
