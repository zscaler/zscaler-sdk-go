package machinegroup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                  = "/zpa/mgmtconfig/v1/admin/customers/"
	machineGroupEndpoint string = "/machineGroup"
)

type MachineGroup struct {
	ID              string     `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description,omitempty"`
	Enabled         bool       `json:"enabled,omitempty"`
	CreationTime    string     `json:"creationTime,omitempty"`
	Machines        []Machines `json:"machines,omitempty"`
	ModifiedBy      string     `json:"modifiedBy,omitempty"`
	ModifiedTime    string     `json:"modifiedTime,omitempty"`
	MicroTenantID   string     `json:"microtenantId,omitempty"`
	MicroTenantName string     `json:"microtenantName,omitempty"`
}

type Machines struct {
	ID               string                 `json:"id,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Description      string                 `json:"description,omitempty"`
	CreationTime     string                 `json:"creationTime,omitempty"`
	Fingerprint      string                 `json:"fingerprint,omitempty"`
	IssuedCertID     string                 `json:"issuedCertId,omitempty"`
	MachineGroupID   string                 `json:"machineGroupId,omitempty"`
	MachineGroupName string                 `json:"machineGroupName,omitempty"`
	MachineTokenID   string                 `json:"machineTokenId,omitempty"`
	ModifiedBy       string                 `json:"modifiedBy,omitempty"`
	ModifiedTime     string                 `json:"modifiedTime,omitempty"`
	MicroTenantID    string                 `json:"microtenantId,omitempty"`
	MicroTenantName  string                 `json:"microtenantName,omitempty"`
	SigningCert      map[string]interface{} `json:"signingCert,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, machineGroupID string) (*MachineGroup, *http.Response, error) {
	v := new(MachineGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+machineGroupEndpoint, machineGroupID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, machineGroupName string) (*MachineGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + machineGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[MachineGroup](ctx, service.Client, relativeURL, common.Filter{Search: machineGroupName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, machineGroupName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no machine group named '%s' was found", machineGroupName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]MachineGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + machineGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[MachineGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetMachineGroupSummary(ctx context.Context, service *zscaler.Service) ([]MachineGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + machineGroupEndpoint + "/summary"
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[MachineGroup](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
