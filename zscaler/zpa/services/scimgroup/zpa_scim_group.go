package scimgroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	userConfig        = "/zpa/userconfig/v1/customers/"
	scimGroupEndpoint = "/scimgroup"
	idpIdPath         = "/idpId"
)

type ScimGroup struct {
	CreationTime int64  `json:"creationTime,omitempty"`
	ID           int64  `json:"id,omitempty"`
	IdpGroupID   string `json:"idpGroupId,omitempty"`
	IdpID        int64  `json:"idpId,omitempty"`
	IdpName      string `json:"idpName,omitempty"`
	ModifiedTime int64  `json:"modifiedTime,omitempty"`
	Name         string `json:"name,omitempty"`
	InternalID   string `json:"internalId,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, scimGroupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint, scimGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, scimName, idpId string) (*ScimGroup, *http.Response, error) {
	// Construct the API endpoint URL with query parameters
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint+idpIdPath, idpId)
	// Fetch the pages
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, common.Filter{
		Search:    scimName,
		SortBy:    string(service.SortBy),
		SortOrder: string(service.SortOrder),
	})
	if err != nil {
		return nil, resp, err
	}

	// Look for the group with the specified name
	for _, scim := range list {
		if strings.EqualFold(scim.Name, scimName) {
			return &scim, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("no SCIM group named '%s' was found", scimName)
}

func GetAllByIdpId(ctx context.Context, service *zscaler.Service, idpId string) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.GetCustomerID()+scimGroupEndpoint+idpIdPath, idpId)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](ctx, service.Client, relativeURL, common.Filter{
		SortBy:    string(service.SortBy),
		SortOrder: string(service.SortOrder),
	})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
