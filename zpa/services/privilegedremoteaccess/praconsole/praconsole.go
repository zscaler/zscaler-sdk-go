package praconsole

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig             = "/mgmtconfig/v1/admin/customers/"
	praConsoleEndpoint     = "/praConsole"
	praConsoleBulkEndpoint = "/praConsole/bulk"
)

type PRAConsole struct {
	// The unique identifier of the privileged console
	ID string `json:"id,omitempty"`

	// The name of the privileged console.
	Name string `json:"name,omitempty"`

	// The description of the privileged console.
	Description string `json:"description,omitempty"`

	// Whether or not the privileged console is enabled.
	Enabled bool `json:"enabled"`

	// The privileged console icon. The icon image is converted to base64 encoded text format.
	IconText string `json:"iconText,omitempty"`

	// The time the privileged console is created.
	CreationTime string `json:"creationTime,omitempty"`

	// The tenant who modified the privileged console.
	ModifiedBy string `json:"modifiedBy,omitempty"`

	// The time the privileged console is modified.
	ModifiedTime    string         `json:"modifiedTime,omitempty"`
	MicroTenantID   string         `json:"microtenantId,omitempty"`
	MicroTenantName string         `json:"microtenantName,omitempty"`
	PRAApplication  PRAApplication `json:"praApplication"`
	PRAPortals      []PRAPortals   `json:"praPortals"`
}

type PRAApplication struct {
	// The unique identifier of the Privileged Remote Access-enabled application.
	ID string `json:"id"`
	// The name of the Privileged Remote Access-enabled application.
	Name string `json:"name"`
}

type PRAPortals struct {
	// The unique identifier of the privileged portal.
	ID string `json:"id"`
	// The name of the privileged portal.
	Name string `json:"name"`
}

func (service *Service) Get(consoleID string) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(consoleName string) (*PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAConsole](service.Client, relativeURL, common.Filter{Search: consoleName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, cred := range list {
		if strings.EqualFold(cred.Name, consoleName) {
			return &cred, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no pra  console named '%s' was found", consoleName)
}

func (service *Service) Create(praConsole *PRAConsole) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, common.Filter{MicroTenantID: service.microTenantID}, praConsole, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) CreatePraBulk(praConsoles []PRAConsole) ([]PRAConsole, *http.Response, error) {
	var responseConsoles []PRAConsole
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praConsoleBulkEndpoint
	resp, err := service.Client.NewRequestDo("POST", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, praConsoles, &responseConsoles)
	if err != nil {
		return nil, nil, err
	}
	return responseConsoles, resp, nil
}

func (service *Service) Update(consoleID string, praConsole *PRAConsole) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("PUT", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, praConsole, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *Service) Delete(consoleID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("DELETE", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAConsole](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
