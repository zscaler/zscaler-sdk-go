package appservercontroller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                  = "/zpa/mgmtconfig/v1/admin/customers/"
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

func Get(service *zscaler.Service, id string) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func GetByName(service *zscaler.Service, appServerName string) (*ApplicationServer, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appServerControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationServer](service.Client, relativeURL, common.Filter{Search: appServerName, MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, appServerName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no application server named '%s' was found", appServerName)
}

func Create(service *zscaler.Service, server ApplicationServer) (*ApplicationServer, *http.Response, error) {
	v := new(ApplicationServer)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.GetCustomerID()+appServerControllerEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, server, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func Update(service *zscaler.Service, id string, appServer ApplicationServer) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("PUT", path, common.Filter{MicroTenantID: service.MicroTenantID()}, appServer, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func Delete(service *zscaler.Service, id string) (*http.Response, error) {
	path := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+appServerControllerEndpoint, id)
	resp, err := service.Client.NewRequestDo("DELETE", path, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(service *zscaler.Service) ([]ApplicationServer, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + appServerControllerEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ApplicationServer](service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
