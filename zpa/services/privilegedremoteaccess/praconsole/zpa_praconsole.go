package praconsole

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig         = "/mgmtconfig/v1/admin/customers/"
	praConsoleEndpoint = "/praConsole"
)

type PRAConsole struct {
	ID             string          `json:"id"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	Enabled        bool            `json:"enabled"`
	IconText       string          `json:"iconText,omitempty"`
	CreationTime   string          `json:"creationTime,omitempty"`
	ModifiedBy     string          `json:"modifiedBy,omitempty"`
	ModifiedTime   string          `json:"modifiedTime,omitempty"`
	SRAApplication *SRAApplication `json:"praApplication"`
	SRAPortal      []SRAPortal     `json:"praPortals"`
}

type SRAApplication struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type SRAPortal struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (service *Service) Get(consoleID string) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, nil, nil, v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) GetByName(consoleName string) (*PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGeneric[PRAConsole](service.Client, relativeURL, consoleName)
	if err != nil {
		return nil, nil, err
	}
	for _, sra := range list {
		if strings.EqualFold(sra.Name, consoleName) {
			return &sra, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no pra console named '%s' was found", consoleName)
}

func (service *Service) Create(sraConsole *PRAConsole) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	resp, err := service.Client.NewRequestDo("POST", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, nil, sraConsole, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func (service *Service) Update(consoleID string, sraConsole *PRAConsole) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("PUT", path, nil, sraConsole, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) Delete(consoleID string) (*http.Response, error) {
	path := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.Config.CustomerID+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo("DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (service *Service) GetAll() ([]PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGeneric[PRAConsole](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
