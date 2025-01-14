package scim_api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
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
func GetGroup(ctx context.Context, service *zscaler.ScimService, groupID string) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetGroupByName(ctx context.Context, service *zscaler.ScimService, groupName string) (*ScimGroup, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint)

	// Use the pagination function with a search function
	list, resp, err := common.GetAllPagesScimGenericWithSearch[ScimGroup](ctx, service.ScimClient, relativeURL, 10, func(group ScimGroup) bool {
		return strings.EqualFold(group.DisplayName, groupName)
	})
	if err != nil {
		return nil, nil, err
	}

	// If no items were returned, the group was not found
	if len(list) == 0 {
		return nil, resp, fmt.Errorf("no SCIM group named '%s' was found", groupName)
	}

	return &list[0], resp, nil
}

func CreateGroup(ctx context.Context, service *zscaler.ScimService, scimGroup *ScimGroup) (*ScimGroup, *http.Response, error) {
	v := new(ScimGroup)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPost, groupScimConfigEndpoint, scimGroup, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateGroup(ctx context.Context, service *zscaler.ScimService, groupID string, scimGroup *ScimGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPut, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func PatchGroup(ctx context.Context, service *zscaler.ScimService, groupID string, scimGroup *ScimGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPatch, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func DeleteGroup(ctx context.Context, service *zscaler.ScimService, groupID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint, groupID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAllGroups(ctx context.Context, service *zscaler.ScimService, count ...int) ([]ScimGroup, *http.Response, error) {
	// Construct the base URL for SCIM groups
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint)

	// Extract count or pass 0 to let the pagination function handle defaults
	itemsPerPage := 0
	if len(count) > 0 && count[0] > 0 {
		itemsPerPage = count[0]
	}

	// Call the pagination function with nil as the searchFunc
	return common.GetAllPagesScimGenericWithSearch[ScimGroup](ctx, service.ScimClient, relativeURL, itemsPerPage, nil)
}
