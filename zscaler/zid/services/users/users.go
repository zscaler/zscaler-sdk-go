package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/services/groups"
)

const (
	usersEndpoint = "/admin/api/v1/users"
)

type Users struct {
	ID              string                    `json:"id,omitempty"`
	Source          string                    `json:"source,omitempty"`
	LoginName       string                    `json:"loginName,omitempty"`
	DisplayName     string                    `json:"displayName,omitempty"`
	FirstName       string                    `json:"firstName,omitempty"`
	LastName        string                    `json:"lastName,omitempty"`
	PrimaryEmail    string                    `json:"primaryEmail,omitempty"`
	SecondaryEmail  string                    `json:"secondaryEmail,omitempty"`
	Status          bool                      `json:"status,omitempty"`
	Department      *common.IDNameDisplayName `json:"department,omitempty"`
	IDP             *common.IDNameDisplayName `json:"idp,omitempty"`
	CustomAttrsInfo map[string]interface{}    `json:"customAttrsInfo,omitempty"`
}

type UsersResponse = common.PaginationResponse[Users]

func GetUser(ctx context.Context, service *zscaler.Service, userID string) (*Users, error) {
	var user Users
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning user from Get: %s", user.ID)
	return &user, nil
}

// GetAll retrieves all users with optional pagination and filtering parameters
func GetAll(ctx context.Context, service *zscaler.Service, queryParams *common.PaginationQueryParams) ([]Users, error) {
	return common.ReadAllPagesWithPagination[Users](ctx, service.Client, usersEndpoint, queryParams)
}

// GetByName retrieves users by searching through paginated data for the specified name
func GetByName(ctx context.Context, service *zscaler.Service, name string) ([]Users, error) {
	var allUsers []Users
	var currentOffset int
	pageSize := 100 // Use a reasonable page size for searching

	for {
		// Create query params for current page
		queryParams := common.NewPaginationQueryParams(pageSize)
		queryParams.WithOffset(currentOffset)

		// Get current page
		pageResponse, err := common.ReadPageWithPagination[Users](ctx, service.Client, usersEndpoint, &queryParams)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Search through records in this page
		for _, user := range pageResponse.Records {
			if strings.Contains(strings.ToLower(user.DisplayName), strings.ToLower(name)) {
				allUsers = append(allUsers, user)
			}
		}

		// Check if we've reached the end
		if len(pageResponse.Records) < pageSize || pageResponse.NextLink == "" {
			break
		}

		// Move to next page
		currentOffset += len(pageResponse.Records)
	}

	return allUsers, nil
}

// GetUsers retrieves users within a specific group
func GetUsers(ctx context.Context, service *zscaler.Service, userID string, queryParams *common.PaginationQueryParams) ([]interface{}, error) {
	usersEndpoint := fmt.Sprintf("%s/%s/users", usersEndpoint, userID)
	return common.ReadAllPagesWithPagination[interface{}](ctx, service.Client, usersEndpoint, queryParams)
}

func Create(ctx context.Context, service *zscaler.Service, user *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.Create(ctx, usersEndpoint, *user)
	if err != nil {
		return nil, nil, err
	}

	createdUser, ok := resp.(*Users)
	if !ok {
		return nil, nil, errors.New("object returned from api was not a user pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning new user from create: %s", createdUser.ID)
	return createdUser, nil, nil
}

func Update(ctx context.Context, service *zscaler.Service, userID string, user *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID), *user)
	if err != nil {
		return nil, nil, err
	}
	updatedUser, _ := resp.(*Users)

	service.Client.GetLogger().Printf("[DEBUG]returning updated user from update: %s", updatedUser.ID)
	return updatedUser, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, userID string) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func GetGroupsByUser(ctx context.Context, service *zscaler.Service, userID string, queryParams *common.PaginationQueryParams) (*common.PaginationResponse[groups.Groups], error) {
	groupsEndpoint := fmt.Sprintf("%s/%s/groups", usersEndpoint, userID)
	return common.ReadPageWithPagination[groups.Groups](ctx, service.Client, groupsEndpoint, queryParams)
}

// func ResetPassword(ctx context.Context, service *zscaler.Service, userID string) (*http.Response, error) {
// 	err, _ := service.Client.Create(ctx, fmt.Sprintf("%s/%s", usersEndpoint, userID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }

// func SetSkipMFA(ctx context.Context, service *zscaler.Service, userID string) (*http.Response, error) {
// 	err, _ := service.Client.Create(ctx, fmt.Sprintf("%s/%s", usersEndpoint+"/resetpassword", userID))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
