package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	usersEndpoint    = "/zia/api/v1/users"
	enrollEndpoint   = "/zia/api/v1/enroll"
	maxBulkDeleteIDs = 500
)

type Users struct {
	// User ID
	ID int `json:"id"`

	// User name. This appears when choosing users for policies.
	Name string `json:"name,omitempty"`

	// User email consists of a user name and domain name. It does not have to be a valid email address, but it must be unique and its domain must belong to the organization.
	Email string `json:"email,omitempty"`

	// List of Groups a user belongs to. Groups are used in policies.
	Groups []common.IDNameExtensions `json:"groups,omitempty"`

	// Department a user belongs to
	Department *common.UserDepartment `json:"department,omitempty"`

	// Additional information about this user.
	Comments string `json:"comments,omitempty"`

	// Temporary Authentication Email. If you enabled one-time tokens or links, enter the email address to which the Zscaler service sends the tokens or links. If this is empty, the service sends the email to the User email.
	TempAuthEmail string `json:"tempAuthEmail,omitempty"`

	// Accepted Authentication Methods. Support values are "BASIC" and "DIGEST"
	AuthMethods []string `json:"authMethods,omitempty"`

	// User's password. Applicable only when authentication type is Hosted DB. Password strength must follow what is defined in the auth settings.
	Password string `json:"password,omitempty"`

	// True if this user is an Admin user
	AdminUser bool `json:"adminUser"`

	// User type. Provided only if this user is not an end user.
	Type string `json:"type,omitempty"`

	Deleted bool `json:"deleted"`
}

type EnrollResult struct {
	UserID int `json:"userId"`
}

type EnrollUserRequest struct {
	AuthMethods []string `json:"authMethods"`
	Password    string   `json:"password"`
}

func Get(ctx context.Context, service *zscaler.Service, userID int) (*Users, error) {
	var user Users
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning user from Get: %d", user.ID)
	return &user, nil
}

func GetUserByName(ctx context.Context, service *zscaler.Service, userName string) (*Users, error) {
	var users []Users
	err := service.Client.Read(ctx, fmt.Sprintf("%s?name=%s&%s", usersEndpoint, url.QueryEscape(userName), common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder))), &users)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if strings.EqualFold(user.Name, userName) {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("no user found with name: %s", userName)
}

func EnrollUser(ctx context.Context, service *zscaler.Service, userID int, request EnrollUserRequest) (*EnrollResult, error) {
	if len(request.AuthMethods) == 0 {
		err := errors.New("authMethods is required")
		service.Client.GetLogger().Printf("[ERROR] enroll user failed: %v", err)
	}
	for _, method := range request.AuthMethods {
		// method most be one of the following: BASIC, DIGEST
		if method != "BASIC" && method != "DIGEST" {
			err := fmt.Errorf("authMethods must be one of the following: BASIC, DIGEST. Found: %s", method)
			service.Client.GetLogger().Printf("[ERROR] enroll user failed: %v", err)
			return nil, err
		}
	}
	_, err := service.Client.Create(ctx, fmt.Sprintf("%s/%d%s", usersEndpoint, userID, enrollEndpoint), request)
	if err != nil {
		return nil, err
	}
	return &EnrollResult{UserID: userID}, nil
}

func Create(ctx context.Context, service *zscaler.Service, user *Users) (*Users, error) {
	resp, err := service.Client.Create(ctx, usersEndpoint, *user)
	if err != nil {
		return nil, err
	}

	createdUsers, ok := resp.(*Users)
	if !ok {
		return nil, errors.New("object returned from api was not a user pointer")
	}

	service.Client.GetLogger().Printf("[DEBUG]returning user from create: %v", createdUsers.ID)
	return createdUsers, nil
}

func Update(ctx context.Context, service *zscaler.Service, userID int, users *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(ctx, fmt.Sprintf("%s/%d", usersEndpoint, userID), *users)
	if err != nil {
		return nil, nil, err
	}
	updatedUser, _ := resp.(*Users)
	service.Client.GetLogger().Printf("[DEBUG]returning user from update: %d", updatedUser.ID)
	return updatedUser, nil, nil
}

func Delete(ctx context.Context, service *zscaler.Service, userID int) (*http.Response, error) {
	err := service.Client.Delete(ctx, fmt.Sprintf("%s/%d", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func BulkDelete(ctx context.Context, service *zscaler.Service, ids []int) (*http.Response, error) {
	if len(ids) > maxBulkDeleteIDs {
		// Truncate the list to the first 100 IDs
		ids = ids[:maxBulkDeleteIDs]
		service.Client.GetLogger().Printf("[INFO] Truncating IDs list to the first %d items", maxBulkDeleteIDs)
	}

	// Define the payload
	payload := map[string][]int{
		"ids": ids,
	}
	return service.Client.BulkDelete(ctx, usersEndpoint+"/bulkDelete", payload)
}

func GetAllUsers(ctx context.Context, service *zscaler.Service) ([]Users, error) {
	var users []Users
	err := common.ReadAllPages(ctx, service.Client, usersEndpoint+"?"+common.GetSortParams(common.SortField(service.SortBy), common.SortOrder(service.SortOrder)), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
