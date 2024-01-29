package emergencyaccess

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig              = "/mgmtconfig/v1/admin/customers/"
	emergencyAccessEndpoint = "/emergencyAccess/user"
)

type EmergencyAccess struct {
	ActivatedOn       string `json:"activatedOn,omitempty"`
	AllowedActivate   bool   `json:"allowedActivate"`
	AllowedDeactivate bool   `json:"allowedDeactivate"`
	EmailId           string `json:"emailId,omitempty"`
	FirstName         string `json:"firstName,omitempty"`
	LastLoginTime     string `json:"lastLoginTime,omitempty"`
	LastName          string `json:"lastName,omitempty"`
	UpdateEnabled     bool   `json:"updateEnabled"`
	UserId            string `json:"userId,omitempty"`
	UserStatus        string `json:"userStatus,omitempty"`
}

func (service *Service) Get(userID string) (*EmergencyAccess, *http.Response, error) {
	v := new(EmergencyAccess)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// POST - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user
func (service *Service) Create(emergencyAccess EmergencyAccess, activateNow bool) (*EmergencyAccess, *http.Response, error) {
	v := new(EmergencyAccess)
	queryParams := url.Values{}
	queryParams.Set("activateNow", strconv.FormatBool(activateNow)) // Adding activateNow as a query parameter
	relativeURL := fmt.Sprintf("%s/%s%s?%s", mgmtConfig, service.Client.Config.CustomerID, emergencyAccessEndpoint, queryParams.Encode())
	resp, err := service.Client.NewRequestDo("POST", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, emergencyAccess, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) Update(userID string, emergencyAccess EmergencyAccess) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, emergencyAccess, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/activate
func (service *Service) Activate(userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/activate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PUT - /mgmtconfig/v1/admin/customers/{customerId}/emergencyAccess/user/{userId}/deactivate
func (service *Service) Deactivate(userID string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s/deactivate", mgmtConfig+service.Client.Config.CustomerID+emergencyAccessEndpoint, userID)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *Service) GetAll() ([]EmergencyAccess, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + emergencyAccessEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[EmergencyAccess](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
