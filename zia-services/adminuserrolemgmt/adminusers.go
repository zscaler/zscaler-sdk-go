package adminuserrolemgmt

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/terraform-provider-zia/gozscaler/common"
)

const (
	adminUsersEndpoint = "/adminUsers"
)

type AdminUsers struct {
	ID                            int                       `json:"id,omitempty"`
	LoginName                     string                    `json:"loginName,omitempty"`
	UserName                      string                    `json:"userName,omitempty"`
	Email                         string                    `json:"email,omitempty"`
	Comments                      string                    `json:"comments,omitempty"`
	Disabled                      bool                      `json:"disabled,omitempty"`
	Password                      string                    `json:"password,omitempty"`
	PasswordLastModifiedTime      int                       `json:"pwdLastModifiedTime,omitempty"`
	IsNonEditable                 bool                      `json:"isNonEditable,omitempty"`
	IsPasswordLoginAllowed        bool                      `json:"isPasswordLoginAllowed,omitempty"`
	IsPasswordExpired             bool                      `json:"isPasswordExpired,omitempty"`
	IsAuditor                     bool                      `json:"isAuditor,omitempty"`
	IsSecurityReportCommEnabled   bool                      `json:"isSecurityReportCommEnabled,omitempty"`
	IsServiceUpdateCommEnabled    bool                      `json:"isServiceUpdateCommEnabled,omitempty"`
	IsProductUpdateCommEnabled    bool                      `json:"isProductUpdateCommEnabled,omitempty"`
	IsExecMobileAppEnabled        bool                      `json:"isExecMobileAppEnabled,omitempty"`
	AdminScopeGroupMemberEntities []common.IDNameExtensions `json:"adminScopescopeGroupMemberEntities,omitempty"`
	AdminScopeEntities            []common.IDNameExtensions `json:"adminScopeScopeEntities,omitempty"`
	AdminScopeType                string                    `json:"adminScopeType,omitempty"`
	Role                          *Role                     `json:"role,omitempty"`
	ExecMobileAppTokens           []ExecMobileAppTokens     `json:"execMobileAppTokens,omitempty"`
}
type Role struct {
	ID           int                    `json:"id,omitempty"`
	Name         string                 `json:"name,omitempty"`
	IsNameL10Tag bool                   `json:"isNameL10nTag,omitempty"`
	Extensions   map[string]interface{} `json:"extensions,omitempty"`
}
type ExecMobileAppTokens struct {
	Cloud       string `json:"cloud,omitempty"`
	OrgId       int    `json:"orgId,omitempty"`
	Name        string `json:"name,omitempty"`
	TokenId     string `json:"tokenId,omitempty"`
	Token       string `json:"token,omitempty"`
	TokenExpiry int    `json:"tokenExpiry,omitempty"`
	CreateTime  int    `json:"createTime,omitempty"`
	DeviceId    string `json:"deviceId,omitempty"`
	DeviceName  string `json:"deviceName,omitempty"`
}

func (service *Service) GetAdminUsers(adminUserId int) (*AdminUsers, error) {
	v := new(AdminUsers)
	relativeURL := fmt.Sprintf("%s/%d", adminUsersEndpoint, adminUserId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetAdminUsersByLoginName(adminUsersLoginName string) (*AdminUsers, error) {
	var adminUsers []AdminUsers
	err := service.Client.Read(adminUsersEndpoint, &adminUsers)
	if err != nil {
		return nil, err
	}
	for _, adminUser := range adminUsers {
		if strings.EqualFold(adminUser.LoginName, adminUsersLoginName) {
			return &adminUser, nil
		}
	}
	return nil, fmt.Errorf("no admin login found with name: %s", adminUsersLoginName)
}

func (service *Service) GetAdminByUsername(adminUsername string) (*AdminUsers, error) {
	var adminUsers []AdminUsers
	err := service.Client.Read(adminUsersEndpoint, &adminUsers)
	if err != nil {
		return nil, err
	}
	for _, adminUser := range adminUsers {
		if strings.EqualFold(adminUser.UserName, adminUsername) {
			return &adminUser, nil
		}
	}
	return nil, fmt.Errorf("no admin found with username: %s", adminUsername)
}

func (service *Service) CreateAdminUser(adminUser AdminUsers) (*AdminUsers, error) {
	resp, err := service.Client.Create(adminUsersEndpoint, adminUser)
	if err != nil {
		return nil, err
	}
	res, ok := resp.(*AdminUsers)
	if !ok {
		return nil, fmt.Errorf("couldn't marshal response to a valid objectm: %#v", resp)
	}
	return res, nil
}

func (service *Service) UpdateAdminUser(adminUserID int, adminUser AdminUsers) (*AdminUsers, error) {
	path := fmt.Sprintf("%s/%d", adminUsersEndpoint, adminUserID)
	resp, err := service.Client.UpdateWithPut(path, adminUser)
	if err != nil {
		return nil, err
	}
	res, _ := resp.(AdminUsers)
	return &res, err
}

func (service *Service) DeleteAdminUser(adminUserID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", adminUsersEndpoint, adminUserID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
