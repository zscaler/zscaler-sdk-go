package adminuserrolemgmt

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	adminRolesEndpoint = "/adminRoles/lite"
)

type AdminRoles struct {
	ID              int      `json:"id"`
	Rank            int      `json:"rank,omitempty"`
	Name            string   `json:"name,omitempty"`
	PolicyAccess    string   `json:"policyAccess,omitempty"`
	DashboardAccess string   `json:"dashboardAccess"`
	ReportAccess    string   `json:"reportAccess,omitempty"`
	AnalysisAccess  string   `json:"analysisAccess,omitempty"`
	UsernameAccess  string   `json:"usernameAccess,omitempty"`
	AdminAcctAccess string   `json:"adminAcctAccess,omitempty"`
	IsAuditor       bool     `json:"isAuditor,omitempty"`
	Permissions     []string `json:"permissions,omitempty"`
	IsNonEditable   bool     `json:"isNonEditable,omitempty"`
	LogsLimit       string   `json:"logsLimit,omitempty"`
	RoleType        string   `json:"roleType,omitempty"`
}

func (service *Service) Get(adminRoleId int) (*AdminRoles, error) {
	v := new(AdminRoles)
	relativeURL := fmt.Sprintf("%s/%d", adminRolesEndpoint, adminRoleId)
	err := service.Client.Read(relativeURL, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (service *Service) GetByName(adminRoleName string) (*AdminRoles, error) {
	var adminRoles []AdminRoles
	err := service.Client.Read(adminRolesEndpoint, &adminRoles)
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

func (service *Service) GetAPIRole(apiRole string) (*AdminRoles, error) {
	var apiRoles []AdminRoles
	err := service.Client.Read(fmt.Sprintf("%s?includeApiRole=%s", adminRolesEndpoint, url.QueryEscape(apiRole)), &apiRoles)
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

func (service *Service) GetAuditorRole(auditorRole string) (*AdminRoles, error) {
	var auditorRoles []AdminRoles
	err := service.Client.Read(fmt.Sprintf("%s?includeAuditorRole=%s", adminRolesEndpoint, url.QueryEscape(auditorRole)), &auditorRoles)
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

func (service *Service) GetPartnerRole(partnerRole string) (*AdminRoles, error) {
	var partnerRoles []AdminRoles
	err := service.Client.Read(fmt.Sprintf("%s?includePartnerRole=%s", adminRolesEndpoint, url.QueryEscape(partnerRole)), &partnerRoles)
	if err != nil {
		return nil, err
	}
	for _, partnerRoleEnabled := range partnerRoles {
		if strings.EqualFold(partnerRoleEnabled.Name, partnerRole) {
			return &partnerRoleEnabled, nil
		}
	}
	return nil, fmt.Errorf("no auditor role found with name: %s", partnerRole)
}

func (service *Service) GetAllAdminRoles() ([]AdminRoles, error) {
	var adminRoles []AdminRoles
	err := service.Client.Read(adminRolesEndpoint, &adminRoles)
	return adminRoles, err
}
