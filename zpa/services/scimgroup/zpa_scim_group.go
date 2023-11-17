package scimgroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	userConfig        = "/userconfig/v1/customers/"
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

func (service *Service) Get(scimGroupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint, scimGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(scimName, idpId string) (*ScimGroup, *http.Response, error) {
	// Construct the API endpoint URL with query parameters
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpIdPath, idpId)
	// Fetch the pages
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](service.Client, relativeURL, common.Filter{
		Search:    scimName,
		SortBy:    string(service.sortBy),
		SortOrder: string(service.sortOrder),
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

func (service *Service) GetAllByIdpId(idpId string) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userConfig+service.Client.Config.CustomerID+scimGroupEndpoint+idpIdPath, idpId)
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ScimGroup](service.Client, relativeURL, common.Filter{
		SortBy:    string(service.sortBy),
		SortOrder: string(service.sortOrder),
	})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
