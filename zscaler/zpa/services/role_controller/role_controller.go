package role_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig               = "/zpa/mgmtconfig/v1/admin/customers/"
	permissionGroupsEndpoint = "/permissionGroups"
	rolesEndpoint            = "/roles"
)

type RoleController struct {
	ID                        string                 `json:"id,omitempty"`
	ModifiedTime              string                 `json:"modifiedTime,omitempty"`
	CreationTime              string                 `json:"creationTime,omitempty"`
	ModifiedBy                string                 `json:"modifiedBy,omitempty"`
	Name                      string                 `json:"name,omitempty"`
	MicrotenantID             string                 `json:"microtenantId,omitempty"`
	MicrotenantName           string                 `json:"microtenantName,omitempty"`
	Description               string                 `json:"description,omitempty"`
	BypassAccestorAccessCheck bool                   `json:"bypassRemoteAssistanceCheck,omitempty"`
	CustomRole                bool                   `json:"customRole,omitempty"`
	SystemRole                bool                   `json:"systemRole,omitempty"`
	RestrictedRole            bool                   `json:"restrictedRole,omitempty"`
	Users                     string                 `json:"users,omitempty"`
	APIKeys                   string                 `json:"apiKeys,omitempty"`
	NewAuditMessage           string                 `json:"newAuditMessage,omitempty"`
	OldAuditMessage           string                 `json:"oldAuditMessage,omitempty"`
	ClassPermissionGroups     []ClassPermissionGroup `json:"classPermissionGroups,omitempty"`
	Permissions               []Permission           `json:"permissions,omitempty"`
}

type Permission struct {
	ID             string    `json:"id,omitempty"`
	PermissionMask string    `json:"permissionMask,omitempty"`
	Role           string    `json:"role,omitempty"`
	CustomerID     string    `json:"customerId,omitempty"`
	ModifiedTime   string    `json:"modifiedTime,omitempty"`
	CreationTime   string    `json:"creationTime,omitempty"`
	ModifiedBy     string    `json:"modifiedBy,omitempty"`
	ClassType      ClassType `json:"classType,omitempty"`
}

type ClassPermission struct {
	Permission PermissionDetail `json:"permission,omitempty"`
	ClassType  ClassType        `json:"classType,omitempty"`

	ID           string `json:"id,omitempty"`
	CreationTime string `json:"creationTime,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
}

type PermissionDetail struct {
	Mask    string `json:"mask,omitempty"`
	MaxMask string `json:"maxMask,omitempty"`
	Type    string `json:"type,omitempty"` // FULL or VIEW_ONLY
}

type ClassType struct {
	ID             string `json:"id,omitempty"`
	ACLClass       string `json:"aclClass,omitempty"`
	FriendlyName   string `json:"friendlyName,omitempty"`
	CustomerID     string `json:"customerId,omitempty"`
	LocalScopeMask string `json:"localScopeMask,omitempty"`
	CreationTime   string `json:"creationTime,omitempty"`
	ModifiedTime   string `json:"modifiedTime,omitempty"`
	ModifiedBy     string `json:"modifiedBy,omitempty"`
}

type ClassPermissionGroup struct {
	ID                        string            `json:"id,omitempty"`
	Name                      string            `json:"name,omitempty"`
	CreationTime              string            `json:"creationTime,omitempty"`
	ModifiedTime              string            `json:"modifiedTime,omitempty"`
	ModifiedBy                string            `json:"modifiedBy,omitempty"`
	Hidden                    bool              `json:"hidden,omitempty"`
	Internal                  bool              `json:"internal,omitempty"`
	LocalScopePermissionGroup bool              `json:"localScopePermissionGroup,omitempty"`
	ClassPermissions          []ClassPermission `json:"classPermissions,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, roleID string) (*RoleController, *http.Response, error) {
	v := new(RoleController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+rolesEndpoint, roleID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, roleName string) (*RoleController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + rolesEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[RoleController](ctx, service.Client, relativeURL, common.Filter{Search: roleName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, role := range list {
		if strings.EqualFold(role.Name, roleName) {
			return &role, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no role named '%s' was found", roleName)
}

func Create(ctx context.Context, service *zscaler.Service, role *RoleController) (*RoleController, *http.Response, error) {
	v := new(RoleController)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+rolesEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, role, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, roleID string, role *RoleController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+rolesEndpoint, roleID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, role, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, roleID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+rolesEndpoint, roleID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]RoleController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + rolesEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[RoleController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetPermissionGroups(ctx context.Context, service *zscaler.Service) ([]ClassPermissionGroup, *http.Response, error) {
	var groups []ClassPermissionGroup
	url := mgmtConfig + service.Client.GetCustomerID() + permissionGroupsEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "GET", url, nil, nil, &groups)
	if err != nil {
		return nil, nil, err
	}
	return groups, resp, nil
}
