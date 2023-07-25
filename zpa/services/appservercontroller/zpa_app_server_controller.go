package appservercontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig                  = "/mgmtconfig/v1/admin/customers/"
	appServerControllerEndpoint = "/server"
)

type ApplicationServer struct {
	Address           string   `json:"address"`
	AppServerGroupIds []string `json:"appServerGroupIds"`
	ConfigSpace       string   `json:"configSpace,omitempty"`
	CreationTime      string   `json:"creationTime,"`
	Description       string   `json:"description"`
	Enabled           bool     `json:"enabled"`
	ID                string   `json:"id,omitempty"`
	ModifiedBy        string   `json:"modifiedBy"`
	ModifiedTime      string   `json:"modifiedTime"`
	Name              string   `json:"name"`
	MicroTenantID     string   `json:"microtenantId,omitempty"`
	MicroTenantName   string   `json:"microtenantName,omitempty"`
}

func (service *Service) Get(id string) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(appServerName string) (*ApplicationServer, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appServerControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationServer](service.Client, relativeURL, common.Filter{Search: appServerName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appServerName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application named '%s' was found", appServerName)
}

func (service *Service) Create(server ApplicationServer) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, common.Filter{MicroTenantID: service.microTenantID}, server, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(id string, appServer ApplicationServer) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.microTenantID}, appServer, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]ApplicationServer, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + appServerControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationServer](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
