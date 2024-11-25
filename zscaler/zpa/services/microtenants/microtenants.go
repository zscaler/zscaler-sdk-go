package microtenants

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig           = "/zpa/mgmtconfig/v1/admin/customers/"
	microtenantsEndpoint = "/microtenants"
)

type MicroTenant struct {
	ID                         string        `json:"id,omitempty"`
	Name                       string        `json:"name,omitempty"`
	Description                string        `json:"description,omitempty"`
	Enabled                    bool          `json:"enabled"`
	CriteriaAttribute          string        `json:"criteriaAttribute,omitempty"`
	CriteriaAttributeValues    []string      `json:"criteriaAttributeValues,omitempty"`
	PrivilegedApprovalsEnabled bool          `json:"privilegedApprovalsEnabled"`
	Operator                   string        `json:"operator,omitempty"`
	Priority                   string        `json:"priority,omitempty"`
	CreationTime               string        `json:"creationTime,omitempty"`
	ModifiedBy                 string        `json:"modifiedBy,omitempty"`
	ModifiedTime               string        `json:"modifiedTime,omitempty"`
	Roles                      []Roles       `json:"roles,omitempty"`
	UserResource               *UserResource `json:"user,omitempty"`
}

type Roles struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	CustomRole bool   `json:"customRole,omitempty"`
}

type UserResource struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	Enabled            bool     `json:"enabled,omitempty"`
	Comments           string   `json:"comments,omitempty"`
	CustomerID         string   `json:"customerId,omitempty"`
	DeliveryTag        string   `json:"deliveryTag,omitempty"`
	DisplayName        string   `json:"displayName,omitempty"`
	Email              string   `json:"email,omitempty"`
	Eula               string   `json:"eula,omitempty"`
	ForcePwdChange     bool     `json:"forcePwdChange,omitempty"`
	GroupIDs           []string `json:"groupIds,omitempty"`
	IAMUserID          string   `json:"iamUserId,omitempty"`
	IsEnabled          bool     `json:"isEnabled,omitempty"`
	IsLocked           bool     `json:"isLocked,omitempty"`
	LanguageCode       string   `json:"languageCode,omitempty"`
	LocalLoginDisabled bool     `json:"localLoginDisabled,omitempty"`
	OneIdentityUser    bool     `json:"oneIdentityUser,omitempty"`
	OperationType      string   `json:"operationType,omitempty"`
	Password           string   `json:"password,omitempty"`
	PhoneNumber        string   `json:"phoneNumber,omitempty"`
	PinSession         bool     `json:"pinSession,omitempty"`
	RoleID             string   `json:"roleId,omitempty"`
	MicrotenantID      string   `json:"microtenantId,omitempty"`
	MicrotenantName    string   `json:"microtenantName,omitempty"`
	SyncVersion        string   `json:"syncVersion,omitempty"`
	Timezone           string   `json:"timezone,omitempty"`
	TmpPassword        string   `json:"tmpPassword,omitempty"`

	// This field is mandatory if twoFactorAuthEnabled is set.
	TokenID string `json:"tokenId,omitempty"`

	TwoFactorAuthEnabled bool `json:"twoFactorAuthEnabled,omitempty"`

	// This field is mandatory if twoFactorAuthEnabled is set. Accepted values: YUBIKEY/TOTP
	TwoFactorAuthType string `json:"twoFactorAuthType,omitempty"`

	// Mandatory only for POST. Not mandatory for PUT/DELETE requests.
	Username string `json:"username,omitempty"`

	// Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
	CreationTime string `json:"creationTime,omitempty"`

	// Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
	ModifiedTime string `json:"modifiedTime,omitempty"`
}

type MicroTenantSummary struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Get(ctx context.Context, service *zscaler.Service, microtenantID string) (*MicroTenant, *http.Response, error) {
	v := new(MicroTenant)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+microtenantsEndpoint, microtenantID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, microTenantName string) (*MicroTenant, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + microtenantsEndpoint
	list, resp, err := common.GetAllPagesGeneric[MicroTenant](ctx, service.Client, relativeURL, microTenantName)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, microTenantName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no microtenant named '%s' was found", microTenantName)
}

func Create(ctx context.Context, service *zscaler.Service, microTenant MicroTenant) (*MicroTenant, *http.Response, error) {
	v := new(MicroTenant)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+microtenantsEndpoint, nil, microTenant, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, microTenantID string, microTenant *MicroTenant) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+microtenantsEndpoint, microTenantID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, nil, microTenant, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, microTenantID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+microtenantsEndpoint, microTenantID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]MicroTenant, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + microtenantsEndpoint
	list, resp, err := common.GetAllPagesGeneric[MicroTenant](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
