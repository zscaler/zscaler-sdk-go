package microtenants

/*
import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig           = "/mgmtconfig/v1/admin/customers/"
	microtenantsEndpoint = "/microtenants"
)

type MicroTenant struct {
	ID                      string        `json:"id"`
	Name                    string        `json:"name,omitempty"`
	Description             string        `json:"description,omitempty"`
	Enabled                 bool          `json:"enabled"`
	CriteriaAttribute       string        `json:"criteriaAttribute,omitempty"`
	CriteriaAttributeValues string        `json:"criteriaAttributeValues,omitempty"`
	Operator                string        `json:"operator,omitempty"`
	Priority                string        `json:"priority,omitempty"`
	CreationTime            string        `json:"creationTime,omitempty"`
	ModifiedBy              string        `json:"modifiedBy,omitempty"`
	ModifiedTime            string        `json:"modifiedTime,omitempty"`
	Roles                   []Roles       `json:"roles,omitempty"`
	UserResource            *UserResource `json:"user,omitempty"`
}

type Roles struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	CustomRole bool   `json:"customRole,omitempty"`
}

type UserResource struct {
	// Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Comments    string `json:"comments,omitempty"`

	// Applicable only for GET request. Ignored in PUT/POST/DELETE requests.
	CustomerID         string `json:"customerId,omitempty"`
	DisplayName        string `json:"displayName,omitempty"`
	Email              string `json:"email,omitempty"`
	Eula               string `json:"eula,omitempty"`
	ForcePwdChange     bool   `json:"forcePwdChange,omitempty"`
	GroupIDs           string `json:"groupIds,omitempty"`
	IsEnabled          bool   `json:"isEnabled,omitempty"`
	IsLocked           bool   `json:"isLocked,omitempty"`
	LanguageCode       string `json:"languageCode,omitempty"`
	LocalLoginDisabled bool   `json:"localLoginDisabled,omitempty"`

	// Mandatory only for POST. Not mandatory for PUT/DELETE requests.
	Password        string `json:"password,omitempty"`
	PhoneNumber     bool   `json:"phoneNumber,omitempty"`
	PinSession      string `json:"pinSession,omitempty"`
	RoleID          bool   `json:"roleId,omitempty"`
	MicrotenantID   string `json:"microtenantId,omitempty"`
	MicrotenantName string `json:"microtenantName,omitempty"`
	Timezone        string `json:"timezone,omitempty"`
	TmpPassword     string `json:"tmpPassword,omitempty"`

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

func (service *Service) Get(id string) (*MicroTenant, *http.Response, error) {
	v := new(MicroTenant)
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+microtenantsEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", path, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(microtenantName string) (*MicroTenant, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + microtenantsEndpoint
	list, resp, err := common.GetAllPagesGeneric[MicroTenant](service.Client, relativeURL, common.IncludeRolesQueryParams{IncludeRoles: true}, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	for _, micro := range list {
		if strings.EqualFold(micro.Name, microtenantName) {
			return &micro, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no microtenant named '%s' was found", microtenantName)
}

func (service *Service) Create(microTenant MicroTenant) (*MicroTenant, *http.Response, error) {
	v := new(MicroTenant)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+microtenantsEndpoint, nil, microTenant, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(microTenantID string, microTenant *MicroTenant) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+microtenantsEndpoint, microTenantID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, nil, microTenant, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(microTenantID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+microtenantsEndpoint, microTenantID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (service *Service) GetAll() ([]MicroTenant, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + microtenantsEndpoint
	list, resp, err := common.GetAllPagesGeneric[MicroTenant](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
*/
