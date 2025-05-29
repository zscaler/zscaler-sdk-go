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
	userScimConfigEndpoint = "/v2/Users"
)

// ScimUsers represents the response from the /users endpoint
// ScimUser represents the SCIM User payload (both core + enterprise extension).
type ScimUser struct {
	Schemas    []string `json:"schemas"`
	ID         string   `json:"id,omitempty"`
	ExternalID string   `json:"externalId,omitempty"`

	// core attributes
	Division     string `json:"division,omitempty"`
	NickName     string `json:"nickName,omitempty"`
	Organization string `json:"organization,omitempty"`
	UserType     string `json:"userType,omitempty"`
	CostCenter   string `json:"costCenter,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Active       bool   `json:"active,omitempty"`
	DisplayName  string `json:"displayName,omitempty"`

	// enterprise extension
	Enterprise EnterpriseFields `json:"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,omitempty"`

	// multi‐valued / complex attrs
	Name   Name    `json:"name,omitempty"`
	Emails []Email `json:"emails,omitempty"`

	Meta common.Meta `json:"meta,omitempty"`
}

// EnterpriseFields holds everything under the “urn:enterprise” extension.
type EnterpriseFields struct {
	Division     string `json:"division,omitempty"`
	Organization string `json:"organization,omitempty"`
	CostCenter   string `json:"costCenter,omitempty"`
	Department   string `json:"department,omitempty"`
}

// Name is the SCIM “name” block.
type Name struct {
	Formatted  string `json:"formatted,omitempty"`
	FamilyName string `json:"familyName,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
}

// Email is a single entry in the “emails” array.
type Email struct {
	Value string `json:"value"`
}

func GetUser(ctx context.Context, service *zscaler.ScimZPAService, userID string) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	relativeURL := fmt.Sprintf("%s%s/%s", service.Client.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.Client.DoRequest(ctx, http.MethodGet, relativeURL, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetUserByName(ctx context.Context, service *zscaler.ScimZPAService, userName string) (*ScimUser, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s", service.Client.ScimConfig.IDPId, userScimConfigEndpoint)

	// Use the pagination function with a search function
	list, resp, err := common.GetAllPagesScimGenericWithSearch(ctx, service.Client, relativeURL, 10, func(group ScimUser) bool {
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

func CreateUser(ctx context.Context, service *zscaler.ScimZPAService, scimUser *ScimUser) (*ScimUser, *http.Response, error) {
	v := new(ScimUser)
	relativeURL := service.Client.ScimConfig.IDPId + userScimConfigEndpoint

	resp, err := service.Client.DoRequest(ctx, http.MethodPost, relativeURL, scimUser, v)
	if err != nil {
		return nil, resp, err
	}
	return v, resp, nil
}

func UpdateUser(ctx context.Context, service *zscaler.ScimZPAService, userID string, scimUser *ScimUser) (*http.Response, error) {
	relativeURL := service.Client.ScimConfig.IDPId + userScimConfigEndpoint + "/" + userID
	resp, err := service.Client.DoRequest(ctx, http.MethodPut, relativeURL, scimUser, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PatchUser(ctx context.Context, service *zscaler.ScimZPAService, userID string, scimUser *ScimUser) (*http.Response, error) {
	relativeURL := service.Client.ScimConfig.IDPId + userScimConfigEndpoint + "/" + userID
	resp, err := service.Client.DoRequest(ctx, http.MethodPatch, relativeURL, scimUser, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, service *zscaler.ScimZPAService, userID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s/%s", service.Client.ScimConfig.IDPId, userScimConfigEndpoint, userID)
	resp, err := service.Client.DoRequest(ctx, http.MethodDelete, relativeURL, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAllUsers(ctx context.Context, service *zscaler.ScimZPAService, count ...int) ([]ScimUser, *http.Response, error) {
	// Construct the base URL for SCIM groups
	relativeURL := fmt.Sprintf("%s%s", service.Client.ScimConfig.IDPId, userScimConfigEndpoint)

	// Extract count or pass 0 to let the pagination function handle defaults
	itemsPerPage := 0
	if len(count) > 0 && count[0] > 0 {
		itemsPerPage = count[0]
	}

	// Call the pagination function with nil as the searchFunc
	return common.GetAllPagesScimGenericWithSearch[ScimUser](ctx, service.Client, relativeURL, itemsPerPage, nil)
}
