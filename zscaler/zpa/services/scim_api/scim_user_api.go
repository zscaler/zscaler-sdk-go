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
	userScimConfigEndpoint = "/v2/Users/"
)

// ScimUsers represents the response from the /users endpoint
type ScimUser struct {
	Schemas     []string    `json:"schemas"`
	ID          string      `json:"id"`
	ExternalID  *string     `json:"externalId,omitempty"`
	DisplayName string      `json:"displayName"`
	Meta        common.Meta `json:"meta"`
}

// User represents an individual user within the Resources array
type User struct {
	Schemas     []string       `json:"schemas"`
	ID          string         `json:"id"`
	ExternalID  string         `json:"externalId"`
	Active      bool           `json:"active"`
	UserName    string         `json:"userName"`
	DisplayName string         `json:"displayName,omitempty"`
	Name        UserName       `json:"name"`
	Groups      []UserGroup    `json:"groups,omitempty"`
	Enterprise  EnterpriseUser `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`
	Meta        common.Meta    `json:"meta"`
}

// UserName represents the user's name information
type UserName struct {
	FamilyName string `json:"familyName"`
	GivenName  string `json:"givenName"`
}

// UserGroup represents a group that the user belongs to
type UserGroup struct {
	Display string `json:"display"`
	Value   string `json:"value"`
	Ref     string `json:"$ref"`
}

// EnterpriseUser represents the enterprise-specific extension for user details
type EnterpriseUser struct {
	Department string `json:"department"`
}

func GetUser(ctx context.Context, service *zscaler.ScimService, userID string) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetUserByName(ctx context.Context, service *zscaler.ScimService, userName string) (*ScimUser, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, groupScimConfigEndpoint)

	// Use the pagination function with a search function
	list, resp, err := common.GetAllPagesScimGenericWithSearch[ScimUser](ctx, service.ScimClient, relativeURL, 10, func(group ScimUser) bool {
		return strings.EqualFold(group.DisplayName, userName)
	})
	if err != nil {
		return nil, nil, err
	}

	// If no items were returned, the user was not found
	if len(list) == 0 {
		return nil, resp, fmt.Errorf("no SCIM user named '%s' was found", userName)
	}

	return &list[0], resp, nil
}

func CreateUser(ctx context.Context, service *zscaler.ScimService, scimUser *ScimUser) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPost, userScimConfigEndpoint, scimUser, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateUser(ctx context.Context, service *zscaler.ScimService, userID string, ScimUser *ScimUser) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPut, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func PatchUser(ctx context.Context, service *zscaler.ScimService, userID string, ScimUser *ScimUser) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPatch, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func DeleteUser(ctx context.Context, service *zscaler.ScimService, userID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAllUsers(ctx context.Context, service *zscaler.ScimService, count ...int) ([]ScimUser, *http.Response, error) {
	// Construct the base URL for SCIM groups
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint)

	// Extract count or pass 0 to let the pagination function handle defaults
	itemsPerPage := 0
	if len(count) > 0 && count[0] > 0 {
		itemsPerPage = count[0]
	}

	// Call the pagination function with nil as the searchFunc
	return common.GetAllPagesScimGenericWithSearch[ScimUser](ctx, service.ScimClient, relativeURL, itemsPerPage, nil)
}
