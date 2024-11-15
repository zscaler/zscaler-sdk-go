package scim_api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
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

func GetUser(ctx context.Context, service *services.ScimService, userID string) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func CreateUser(ctx context.Context, service *services.ScimService, scimUser *ScimUser) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPost, userScimConfigEndpoint, scimUser, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func UpdateUser(ctx context.Context, service *services.ScimService, userID string, ScimUser *ScimUser) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPut, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func PatchUser(ctx context.Context, service *services.ScimService, userID string, ScimUser *ScimUser) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodPatch, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func DeleteUser(ctx context.Context, service *services.ScimService, userID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.ScimClient.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAllUsers(ctx context.Context, service *services.ScimService, count ...int) ([]ScimUser, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s", service.ScimClient.ScimConfig.IDPId, userScimConfigEndpoint)

	var itemsPerPage int
	if len(count) > 0 && count[0] > 0 {
		itemsPerPage = count[0]
	}

	list, resp, err := common.GetAllPagesScimGeneric[ScimUser](ctx, service.ScimClient, relativeURL, itemsPerPage)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
