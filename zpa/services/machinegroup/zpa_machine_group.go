package machinegroup

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
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

func (service *Service) Get(machineGroupID string) (*MachineGroup, *http.Response, error) {
	v := new(MachineGroup)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+machineGroupEndpoint, machineGroupID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(machineGroupName string) (*MachineGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + machineGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[MachineGroup](service.Client, relativeURL, common.Filter{Search: machineGroupName, MicroTenantID: service.microTenantID})
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

func (service *Service) GetAll() ([]MachineGroup, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + machineGroupEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[MachineGroup](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
