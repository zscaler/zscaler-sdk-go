package scim_api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	groupScimConfigEndpoint = "/Groups/"
)

type SCIMGroup struct {
	Schemas     []string          `json:"schemas,omitempty"`
	ID          string            `json:"id,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	ExternalID  string            `json:"externalId,omitempty"`
	Members     []SCIMGroupMember `json:"members,omitempty"`
	Meta        *Meta             `json:"meta,omitempty"`
}

type SCIMGroupMember struct {
	Value string `json:"value,omitempty"`
	Ref   string `json:"$ref,omitempty"` // shown in response, not needed in requests
}

type SCIMGroupListResponse struct {
	Schemas      []string    `json:"schemas,omitempty"`
	TotalResults int         `json:"totalResults,omitempty"`
	StartIndex   int         `json:"startIndex,omitempty"`
	ItemsPerPage int         `json:"itemsPerPage,omitempty"`
	Resources    []SCIMGroup `json:"Resources,omitempty"`
}

func GetGroup(ctx context.Context, service *zscaler.ScimZIAService, groupID string) (*SCIMGroup, *http.Response, error) {
	v := new(SCIMGroup)
	relativeURL := fmt.Sprintf("%s/%s", groupScimConfigEndpoint, groupID)
	resp, err := service.Client.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetGroupByName(ctx context.Context, service *zscaler.ScimZIAService, displayName string) (*SCIMGroup, *http.Response, error) {
	list, resp, err := common.GetAllPagesScimPostWithSearch(
		ctx,
		service.Client,
		"/Groups/.search",
		100,
		func(g SCIMGroup) bool {
			return strings.EqualFold(g.DisplayName, displayName)
		},
	)
	if err != nil {
		return nil, resp, err
	}
	if len(list) == 0 {
		return nil, resp, fmt.Errorf("no SCIM group found with display name '%s'", displayName)
	}
	return &list[0], resp, nil
}

func CreateGroup(ctx context.Context, service *zscaler.ScimZIAService, group *SCIMGroup) (*SCIMGroup, *http.Response, error) {
	v := new(SCIMGroup)
	relativeURL := groupScimConfigEndpoint

	resp, err := service.Client.DoRequest(ctx, http.MethodPost, relativeURL, group, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateGroup(ctx context.Context, service *zscaler.ScimZIAService, groupID string, scimGroup *SCIMGroup) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", groupScimConfigEndpoint, groupID)

	resp, err := service.Client.DoRequest(ctx, http.MethodPut, relativeURL, scimGroup, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteGroup(ctx context.Context, service *zscaler.ScimZIAService, groupID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", groupScimConfigEndpoint, groupID)

	resp, err := service.Client.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAllGroups(ctx context.Context, service *zscaler.ScimZIAService) ([]SCIMGroup, *http.Response, error) {
	return common.GetAllPagesScimPostWithSearch[SCIMGroup](
		ctx,
		service.Client,
		"/Groups/.search",
		100, // max per Zscaler SCIM API
		nil, // no filter
	)
}
