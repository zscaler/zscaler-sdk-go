package usermanagement

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zia/services/common"
)

const (
	usersEndpoint  = "/users"
	enrollEndpoint = "/enroll"
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

func (service *Service) Get(userID int) (*Users, error) {
	var user Users
	err := service.Client.Read(fmt.Sprintf("%s/%d", usersEndpoint, userID), &user)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning user from Get: %d", user.ID)
	return &user, nil
}

func (service *Service) GetUserByName(userName string) (*Users, error) {
	var users []Users
	err := service.Client.Read(fmt.Sprintf("%s?name=%s", usersEndpoint, url.QueryEscape(userName)), &users)
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

func (service *Service) Create(userID *Users) (*Users, error) {
	resp, err := service.Client.Create(usersEndpoint, *userID)
	if err != nil {
		return nil, err
	}

	createdUsers, ok := resp.(*Users)
	if !ok {
		return nil, errors.New("object returned from api was not a user pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning user from create: %v", createdUsers.ID)
	return createdUsers, nil
}

func (service *Service) Update(userID int, users *Users) (*Users, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", usersEndpoint, userID), *users)
	if err != nil {
		return nil, nil, err
	}
	updatedUser, _ := resp.(*Users)
	service.Client.Logger.Printf("[DEBUG]returning user from update: %d", updatedUser.ID)
	return updatedUser, nil, nil
}

func (service *Service) Delete(userID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", usersEndpoint, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAllUsers() ([]Users, error) {
	var users []Users
	err := common.ReadAllPages(service.Client, usersEndpoint, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
