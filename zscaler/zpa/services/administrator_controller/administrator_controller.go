package administrator_controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig            = "/zpa/mgmtconfig/v1/admin/customers/"
	administratorEndpoint = "/administrators"
)

type AdministratorController struct {
	ID                   string   `json:"id,omitempty"`
	Username             string   `json:"username,omitempty"`
	DisplayName          string   `json:"displayName,omitempty"`
	Email                string   `json:"email,omitempty"`
	Timezone             string   `json:"timezone,omitempty"`
	Password             string   `json:"password,omitempty"`
	TmpPassword          string   `json:"tmpPassword,omitempty"`
	RoleId               string   `json:"roleId,omitempty"`
	Comments             string   `json:"comments,omitempty"`
	LanguageCode         string   `json:"languageCode,omitempty"`
	Eula                 string   `json:"eula,omitempty"`
	IsEnabled            bool     `json:"isEnabled,omitempty"`
	ForcePwdChange       bool     `json:"forcePwdChange,omitempty"`
	TwoFactorAuthEnabled bool     `json:"twoFactorAuthEnabled,omitempty"`
	TwoFactorAuthType    string   `json:"twoFactorAuthType,omitempty"`
	TokenId              string   `json:"tokenId,omitempty"`
	PhoneNumber          string   `json:"phoneNumber,omitempty"`
	LocalLoginDisabled   bool     `json:"localLoginDisabled,omitempty"`
	PinSession           bool     `json:"pinSession,omitempty"`
	IsLocked             bool     `json:"isLocked,omitempty"`
	SyncVersion          string   `json:"syncVersion,omitempty"`
	DeliveryTag          string   `json:"deliveryTag,omitempty"`
	OperationType        string   `json:"operationType,omitempty"`
	GroupIds             []string `json:"groupIds,omitempty"`
	MicrotenantId        string   `json:"microtenantId,omitempty"`
	MicrotenantName      string   `json:"microtenantName,omitempty"`
	Role                 Role     `json:"role,omitempty"`
}

type Role struct {
	ID string `json:"id,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, adminID string) (*AdministratorController, *http.Response, error) {
	v := new(AdministratorController)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+administratorEndpoint, adminID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, adminName string) (*AdministratorController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + administratorEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AdministratorController](ctx, service.Client, relativeURL, common.Filter{Search: adminName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, admin := range list {
		if strings.EqualFold(admin.Username, adminName) {
			return &admin, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no administrator username named '%s' was found", adminName)
}

func Create(ctx context.Context, service *zscaler.Service, admin *AdministratorController) (*AdministratorController, *http.Response, error) {
	v := new(AdministratorController)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+administratorEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, admin, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, adminID string, admin *AdministratorController) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+administratorEndpoint, adminID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, admin, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, adminID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+administratorEndpoint, adminID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]AdministratorController, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + administratorEndpoint

	list, resp, err := common.GetAllPagesGenericWithCustomFilters[AdministratorController](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
