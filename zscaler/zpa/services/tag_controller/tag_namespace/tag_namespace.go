package tag_namespace

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
	namespaceEndpoint        = "/namespace"
	namespaceSearchEndpoint  = "/namespace/search"
)

type Namespace struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	Enabled         bool   `json:"enabled"`
	Origin          string `json:"origin,omitempty"`
	Type            string `json:"type,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
}

type UpdateStatusRequest struct {
	Enabled         bool   `json:"enabled"`
	NamespaceID     string `json:"namespaceId,omitempty"`
	MicroTenantID   string `json:"microtenantId,omitempty"`
	MicroTenantName string `json:"microtenantName,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, namespaceID string) (*Namespace, *http.Response, error) {
	v := new(Namespace)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+namespaceEndpoint, namespaceID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, namespaceName string) (*Namespace, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + namespaceSearchEndpoint
	searchRequest := common.SearchRequest{
		FilterBy: &common.SearchFilterBy{
			FilterGroups: []common.SearchFilterGroup{
				{
					Filters: []common.SearchFilterItem{
						{
							FilterName: "name",
							Operator:   "EQ",
							Value:      namespaceName,
						},
					},
					Operator: "AND",
				},
			},
			Operator: "AND",
		},
	}
	list, resp, err := common.GetAllPagesGenericWithPostSearch[Namespace](ctx, service.Client, relativeURL, searchRequest, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, ns := range list {
		if strings.EqualFold(ns.Name, namespaceName) {
			return &ns, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no namespace named '%s' was found", namespaceName)
}

func Create(ctx context.Context, service *zscaler.Service, namespace Namespace) (*Namespace, *http.Response, error) {
	v := new(Namespace)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+namespaceEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, namespace, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, namespaceID string, namespace *Namespace) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+namespaceEndpoint, namespaceID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, namespace, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, namespaceID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+namespaceEndpoint, namespaceID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]Namespace, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + namespaceSearchEndpoint
	list, resp, err := common.GetAllPagesGenericWithPostSearch[Namespace](ctx, service.Client, relativeURL, common.SearchRequest{}, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func UpdateStatus(ctx context.Context, service *zscaler.Service, namespaceID string, statusUpdate UpdateStatusRequest) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v/status", mgmtConfig+service.Client.GetCustomerID()+namespaceEndpoint, namespaceID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, statusUpdate, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
