package praapproval

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                 = "/zpa/mgmtconfig/v1/admin/customers/"
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

func Get(ctx context.Context, service *zscaler.Service, approvalID string) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByEmailID(ctx context.Context, service *zscaler.Service, emailID string) (*PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivilegedApproval](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, privilegedApproval *PrivilegedApproval) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+privilegedApprovalEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, privilegedApproval, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, approvalID string, privilegedApproval *PrivilegedApproval) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, privilegedApproval, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, approvalID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func DeleteExpired(ctx context.Context, service *zscaler.Service) (*http.Response, error) {
	path := fmt.Sprintf("%s%s%s/expired", mgmtConfig, service.Client.GetCustomerID(), privilegedApprovalEndpoint)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PrivilegedApproval](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
