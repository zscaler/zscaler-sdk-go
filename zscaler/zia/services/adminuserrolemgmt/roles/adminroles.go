package roles

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const (
	adminRolesEndpoint = "/zia/api/v1/adminRoles/lite"
)

type AdminRoles struct {
	// Admin role Id
	ID int `json:"id"`

	// Admin rank of this admin role. This is applicable only when admin rank is enabled in the advanced settings. Default value is 7 (the lowest rank). The assigned admin rank determines the roles or admin users this user can manage, and which rule orders this admin can access.
	Rank int `json:"rank,omitempty"`

	// Name of the admin role
	Name string `json:"name,omitempty"`

	// Policy access permission
	PolicyAccess string `json:"policyAccess,omitempty"`

	// Dashboard access permission
	DashboardAccess string `json:"dashboardAccess"`

	// Report access permission
	ReportAccess string `json:"reportAccess,omitempty"`

	// Insights logs access permission
	AnalysisAccess string `json:"analysisAccess,omitempty"`

	// Username access permission. When set to NONE, the username will be obfuscated
	UsernameAccess string `json:"usernameAccess,omitempty"`

	// Admin and role management access permission
	AdminAcctAccess string `json:"adminAcctAccess,omitempty"`

	// Indicates whether this is an auditor role
	IsAuditor bool `json:"isAuditor,omitempty"`

	// List of functional areas to which this role has access. This attribute is subject to change
	Permissions []string `json:"permissions,omitempty"`

	// Indicates whether or not this admin user is editable/deletable
	IsNonEditable bool `json:"isNonEditable,omitempty"`

	// Log range limit
	LogsLimit string `json:"logsLimit,omitempty"`

	// The admin role type. ()This attribute is subject to change.)
	RoleType string `json:"roleType,omitempty"`
}

// func (service *Service) Get(adminRoleId int) (*AdminRoles, error) {
// 	v := new(AdminRoles)
// 	relativeURL := fmt.Sprintf("%s/%d", adminRolesEndpoint, adminRoleId)
// 	err := service.Client.Read(ctx, relativeURL, v)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return v, nil
// }

func GetByName(ctx context.Context, service *zscaler.Service, adminRoleName string) (*AdminRoles, error) {
	var adminRoles []AdminRoles
	err := service.Client.Read(ctx, adminRolesEndpoint, &adminRoles)
	if err != nil {
		return nil, err
	}
	for _, adminRole := range adminRoles {
		if strings.EqualFold(adminRole.Name, adminRoleName) {
			return &adminRole, nil
		}
	}
	return nil, fmt.Errorf("no admin role found with name: %s", adminRoleName)
}

func GetAPIRole(ctx context.Context, service *zscaler.Service, apiRole, includeApiRole string) (*AdminRoles, error) {
	var apiRoles []AdminRoles
	err := service.Client.Read(ctx, fmt.Sprintf("%s?includeApiRole=%s", adminRolesEndpoint, url.QueryEscape(includeApiRole)), &apiRoles)
	if err != nil {
		return nil, err
	}
	for _, apiRoleEnabled := range apiRoles {
		if strings.EqualFold(apiRoleEnabled.Name, apiRole) {
			return &apiRoleEnabled, nil
		}
	}
	return nil, fmt.Errorf("no api role found with name: %s", apiRole)
}

func GetAuditorRole(ctx context.Context, service *zscaler.Service, auditorRole, includeAuditorRole string) (*AdminRoles, error) {
	var auditorRoles []AdminRoles
	err := service.Client.Read(ctx, fmt.Sprintf("%s?includeAuditorRole=%s", adminRolesEndpoint, url.QueryEscape(includeAuditorRole)), &auditorRoles)
	if err != nil {
		return nil, err
	}
	for _, auditorRoleEnabled := range auditorRoles {
		if strings.EqualFold(auditorRoleEnabled.Name, auditorRole) {
			return &auditorRoleEnabled, nil
		}
	}
	return nil, fmt.Errorf("no auditor role found with name: %s", auditorRole)
}

func GetPartnerRole(ctx context.Context, service *zscaler.Service, partnerRole, includePartnerRole string) (*AdminRoles, error) {
	var partnerRoles []AdminRoles
	err := service.Client.Read(ctx, fmt.Sprintf("%s?includePartnerRole=%s", adminRolesEndpoint, url.QueryEscape(includePartnerRole)), &partnerRoles)
	if err != nil {
		return nil, err
	}
	for _, partnerRoleEnabled := range partnerRoles {
		if strings.EqualFold(partnerRoleEnabled.Name, partnerRole) {
			return &partnerRoleEnabled, nil
		}
	}
	return nil, fmt.Errorf("no partner role found with name: %s", partnerRole)
}

func GetAllAdminRoles(ctx context.Context, service *zscaler.Service) ([]AdminRoles, error) {
	var adminRoles []AdminRoles
	err := service.Client.Read(ctx, adminRolesEndpoint, &adminRoles)
	return adminRoles, err
}
