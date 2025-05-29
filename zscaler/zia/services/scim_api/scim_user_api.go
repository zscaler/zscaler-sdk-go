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
	userScimConfigEndpoint = "/Users/"
)

type SCIMUserListResponse struct {
	Schemas      []string   `json:"schemas,omitempty"`
	TotalResults int        `json:"totalResults,omitempty"`
	StartIndex   int        `json:"startIndex,omitempty"`
	Resources    []SCIMUser `json:"Resources,omitempty"`
}

type SCIMUser struct {
	Schemas             []string        `json:"schemas,omitempty"`
	ID                  string          `json:"id,omitempty"`
	UserName            string          `json:"userName,omitempty"`
	DisplayName         string          `json:"displayName,omitempty"`
	EnterpriseExtension *EnterpriseUser `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	Meta                *Meta           `json:"meta,omitempty"`
}

type EnterpriseUser struct {
	Department string `json:"department,omitempty"`
}

type Meta struct {
	Created      string `json:"created,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Location     string `json:"location,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
}

func GetUser(ctx context.Context, service *zscaler.ScimZIAService, userID string) (*SCIMUser, *http.Response, error) {
	v := new(SCIMUser)
	relativeURL := fmt.Sprintf("%s/%s", userScimConfigEndpoint, userID)
	resp, err := service.Client.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetUserByName(ctx context.Context, service *zscaler.ScimZIAService, displayName string) (*SCIMUser, *http.Response, error) {
	list, resp, err := common.GetAllPagesScimPostWithSearch(
		ctx,
		service.Client,
		"/Users/.search",
		100,
		func(g SCIMUser) bool {
			return strings.EqualFold(g.DisplayName, displayName)
		},
	)
	if err != nil {
		return nil, resp, err
	}
	if len(list) == 0 {
		return nil, resp, fmt.Errorf("no SCIM user found with display name '%s'", displayName)
	}
	return &list[0], resp, nil
}

func CreateUser(ctx context.Context, service *zscaler.ScimZIAService, scimUser *SCIMUser) (*SCIMUser, *http.Response, error) {
	v := new(SCIMUser)
	relativeURL := userScimConfigEndpoint

	resp, err := service.Client.DoRequest(ctx, http.MethodPost, relativeURL, scimUser, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateUser(ctx context.Context, service *zscaler.ScimZIAService, userID string, scimUser *SCIMUser) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userScimConfigEndpoint, userID)

	resp, err := service.Client.DoRequest(ctx, http.MethodPut, relativeURL, scimUser, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, service *zscaler.ScimZIAService, userID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", userScimConfigEndpoint, userID)

	resp, err := service.Client.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAllUsers(ctx context.Context, service *zscaler.ScimZIAService) ([]SCIMGroup, *http.Response, error) {
	return common.GetAllPagesScimPostWithSearch[SCIMGroup](
		ctx,
		service.Client,
		"/Users/.search",
		100, // max per Zscaler SCIM API
		nil, // no filter
	)
}
