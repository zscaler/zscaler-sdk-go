package scim_api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	groupScimConfigEndpoint = "/v2/Groups"
)

// ScimGroups represents the response from the /groups endpoint
type ScimGroup struct {
	Schemas     []string    `json:"schemas"`
	ID          string      `json:"id"`
	ExternalID  *string     `json:"externalId,omitempty"`
	DisplayName string      `json:"displayName"`
	Meta        common.Meta `json:"meta"`
}

// Group represents an individual group within the Resources array in the response
type Group struct {
	Schemas     []string    `json:"schemas"`
	ID          string      `json:"id"`
	ExternalID  *string     `json:"externalId,omitempty"`
	DisplayName string      `json:"displayName"`
	Meta        common.Meta `json:"meta"`
}

// GetGroup retrieves a specific SCIM group by groupID
func GetGroup(ctx context.Context, service *services.ScimService, groupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func CreateGroup(ctx context.Context, service *services.ScimService, scimGroup *ScimGroup) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPost, groupScimConfigEndpoint, scimGroup, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateGroup(ctx context.Context, service *services.ScimService, groupID string, scimGroup *ScimGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPut, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func PatchGroup(ctx context.Context, service *services.ScimService, groupID string, scimGroup *ScimGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPatch, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func DeleteGroup(ctx context.Context, service *services.ScimService, groupID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAllGroups(ctx context.Context, service *services.ScimService, count ...int) ([]ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint)

	var itemsPerPage int
	if len(count) > 0 && count[0] > 0 {
		itemsPerPage = count[0]
	}

	list, resp, err := common.GetAllPagesScimGeneric[ScimGroup](ctx, service.ScimClient, relativeURL, itemsPerPage)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
