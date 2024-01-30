package privilegedapproval

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                 = "/mgmtconfig/v1/admin/customers/"
	privilegedApprovalEndpoint = "/privilegedApproval"
)

type PrivilegedApproval struct {
	ID           string         `json:"id,omitempty"`
	EmailIDs     []string       `json:"emailIds,omitempty"`
	StartTime    time.Time      `json:"startTime,omitempty"`
	EndTime      time.Time      `json:"endTime,omitempty"`
	Status       string         `json:"status,omitempty"`
	CreationTime string         `json:"creationTime,omitempty"`
	ModifiedBy   string         `json:"modifiedBy,omitempty"`
	ModifiedTime string         `json:"modifiedTime,omitempty"`
	WorkingHours *WorkingHours  `json:"workingHours"`
	Applications []Applications `json:"applications"`
}

type Applications struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkingHours struct {
	Days          []string       `json:"days,omitempty"`
	EndTime       time.Time      `json:"endTime,omitempty"`
	StartTime     time.Time      `json:"startTime,omitempty"`
	EndTimeCron   string         `json:"endTimeCron,omitempty"`
	StartTimeCron string         `json:"startTimeCron,omitempty"`
	TimeZone      *time.Location `json:"timeZone,omitempty"`
}

// UnmarshalJSON customizes the unmarshalling process to handle the epoch time.
func (p *PrivilegedApproval) UnmarshalJSON(data []byte) error {
	type Alias PrivilegedApproval
	auxiliary := &struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	startTimeInt, err := strconv.ParseInt(auxiliary.StartTime, 10, 64)
	if err == nil {
		p.StartTime = time.Unix(startTimeInt, 0)
	}
	endTimeInt, err := strconv.ParseInt(auxiliary.EndTime, 10, 64)
	if err == nil {
		p.EndTime = time.Unix(endTimeInt, 0)
	}
	return nil
}

func (w *WorkingHours) UnmarshalJSON(data []byte) error {
	type Alias WorkingHours
	auxiliary := &struct {
		TimeZoneStr string `json:"timeZone"`
		*Alias
	}{
		Alias: (*Alias)(w),
	}
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	loc, err := time.LoadLocation(auxiliary.TimeZoneStr)
	if err != nil {
		return err
	}
	w.TimeZone = loc
	return nil
}

func (service *Service) Get(approvalID string) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

/*
// Need to implement search by Email ID
func (service *Service) GetByEmailID(emailID string) (*PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGeneric[PrivilegedApproval](service.Client, relativeURL, emailID)
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.EmailIDs[], emailID) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", emailID)
}
*/

func (service *Service) Create(privilegedApproval *PrivilegedApproval) (*PrivilegedApproval, *http.Response, error) {
	v := new(PrivilegedApproval)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, nil, privilegedApproval, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(approvalID string, privilegedApproval *PrivilegedApproval) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, privilegedApproval, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(approvalID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+privilegedApprovalEndpoint, approvalID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]PrivilegedApproval, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + privilegedApprovalEndpoint
	list, resp, err := common.GetAllPagesGeneric[PrivilegedApproval](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
