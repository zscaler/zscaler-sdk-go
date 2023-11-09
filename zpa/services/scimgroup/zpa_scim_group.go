package scimgroup

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	userConfig        = "/userconfig/v1/customers/"
	scimGroupEndpoint = "/scimgroup"
	idpId             = "/idpId"
	SortOrderAsc      = "ASC"
	SortOrderDesc     = "DSC"
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

func (service *Service) Get(scimGroupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint, scimGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// func (service *Service) GetByName(scimName, IdpId string) (*ScimGroup, *http.Response, error) {
// 	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpId, IdpId)
// 	list, resp, err := common.GetAllPagesGeneric[ScimGroup](service.Client, relativeURL, "")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	for _, scim := range list {
// 		if strings.EqualFold(scim.Name, scimName) {
// 			return &scim, resp, nil
// 		}
// 	}
// 	return nil, resp, fmt.Errorf("no scim named '%s' was found", scimName)
// }

func paginationQuery(p common.Pagination) string {
	v := url.Values{}
	v.Set("pagesize", strconv.Itoa(p.PageSize))
	v.Set("page", strconv.Itoa(p.Page))
	if p.SortBy != "" {
		v.Set("sortBy", p.SortBy)
	}
	if p.SortOrder != "" {
		v.Set("sortOrder", p.SortOrder)
	}
	// add other fields as needed...
	return v.Encode()
}

func (service *Service) GetByName(scimName, idpId, sortOrder string) (*ScimGroup, *http.Response, error) {
	// Validate the sortOrder input
	if sortOrder != SortOrderAsc && sortOrder != SortOrderDesc {
		return nil, nil, fmt.Errorf("invalid sort order: %s", sortOrder)
	}

	// Set up pagination with sort options
	pagination := common.Pagination{
		PageSize:  common.DefaultPageSize,
		Page:      1,
		SortBy:    "name",
		SortOrder: sortOrder, // Using the user-specified sort order
	}

	// Construct the API endpoint URL with query parameters
	relativeURL := fmt.Sprintf("%s%s%s%s", userConfig, service.Client.Config.CustomerID, scimGroupEndpoint, idpId)
	query := paginationQuery(pagination) // Use the helper function to get the query string.
	if query != "" {
		relativeURL = fmt.Sprintf("%s?%s", relativeURL, query)
	}

	// Fetch the pages
	list, resp, err := common.GetAllPagesGeneric[ScimGroup](service.Client, relativeURL, "")
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

func (service *Service) GetAllByIdpId(IdpId string) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpId, IdpId)
	list, resp, err := common.GetAllPagesGeneric[ScimGroup](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
