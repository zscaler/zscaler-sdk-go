package praapproval

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                 = "/mgmtconfig/v1/admin/customers/"
	privilegedApprovalEndpoint = "/approval"
)

type PrivilegedApproval struct {
	// The unique identifier of the privileged approval.
	ID string `json:"id,omitempty"`

	// The email address of the user that you are assigning the privileged approval to.
	EmailIDs []string `json:"emailIds,omitempty"`

	// The start date that the user has access to the privileged approval.
	StartTime string `json:"startTime,omitempty"`

	// StartTime    time.Time      `json:"startTime,omitempty"`
	// EndTime      time.Time      `json:"endTime,omitempty"`
	// The end date that the user no longer has access to the privileged approval.
	EndTime string `json:"endTime,omitempty"`

	// The status of the privileged approval. The supported values are:
	// INVALID: The privileged approval is invalid.
	// ACTIVE: The privileged approval is currently available for the user.
	// FUTURE: The privileged approval is available for a user at a set time in the future.
	// EXPIRED: The privileged approval is no longer available for the user.
	Status string `json:"status,omitempty"`

	// The time the privileged approval is created.
	CreationTime string `json:"creationTime,omitempty"`

	// The unique identifier of the tenant who modified the privileged approval.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// The time the privileged approval is modified.
	ModifiedTime string `json:"modifiedTime,omitempty"`

	// The unique identifier of the Microtenant for the ZPA tenant.
	// If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant.
	// Pass microtenantId as null to retrieve data from all customers associated with the tenant.
	MicroTenantID string `json:"microtenantId,omitempty"`

	// The name of the Microtenant.
	MicroTenantName string `json:"microtenantName,omitempty"`

	WorkingHours *WorkingHours `json:"workingHours"`
	// The List of application segments
	Applications []Applications `json:"applications"`
}

// The List of application segments
type Applications struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkingHours struct {
	// The days of the week that you want to enable the privileged approval.
	Days []string `json:"days,omitempty"`
	// EndTime       time.Time `json:"endTime,omitempty"`
	// StartTime     time.Time `json:"startTime,omitempty"`

	// The start time that the user has access to the privileged approval.
	StartTime string `json:"startTime,omitempty"`

	// The end time that the user no longer has access to the privileged approval.
	EndTime string `json:"endTime,omitempty"`

	//The cron expression provided to configure the privileged approval start time working hours.
	// The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year].
	// For example, 0 15 10 ? * MON-FRI represents the start time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
	StartTimeCron string `json:"startTimeCron,omitempty"`

	// The cron expression provided to configure the privileged approval end time working hours.
	// The standard cron expression format is [Seconds][Minutes][Hours][Day of the Month][Month][Day of the Week][Year].
	// For example, 0 15 10 ? * MON-FRI represents the end time working hours for 10:15 AM every Monday, Tuesday, Wednesday, Thursday and Friday.
	EndTimeCron string `json:"endTimeCron,omitempty"`

	// The time zone for the time window of a privileged approval.
	TimeZone string `json:"timeZone,omitempty"`
	// TimeZone *time.Location `json:"timeZone,omitempty"`
}

func (service *Service) Get(approvalID string) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByEmailID(emailID string) (*PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGeneric[PrivilegedApproval](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		for _, appEmailID := range app.EmailIDs {
			if strings.EqualFold(appEmailID, emailID) {
				return &app, resp, nil
			}
		}
	}
	return nil, resp, fmt.Errorf("no privileged approval with emailID '%s' was found", emailID)
}

func (service *Service) Create(privilegedApproval *PrivilegedApproval) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, common.Filter{MicroTenantID: service.microTenantID}, privilegedApproval, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(approvalID string, privilegedApproval *PrivilegedApproval) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, privilegedApproval, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(approvalID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) DeleteExpired() (*http.Response, error) {
	path := fmt.Sprintf("%s%s%s/expired", mgmtConfig, service.Client.Config.CustomerID, privilegedApprovalEndpoint)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (service *Service) GetAll() ([]PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivilegedApproval](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
