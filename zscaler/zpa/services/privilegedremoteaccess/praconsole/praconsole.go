package praconsole

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
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
	PRAApplication  PRAApplication `json:"praApplication,omitempty"`
	PRAPortals      []PRAPortals   `json:"praPortals"`
}

type PRAApplication struct {
	// The unique identifier of the Privileged Remote Access-enabled application.
	ID string `json:"id,omitempty"`
	// The name of the Privileged Remote Access-enabled application.
	Name string `json:"name,omitempty"`
}

type PRAPortals struct {
	// The unique identifier of the privileged portal.
	ID string `json:"id,omitempty"`
	// The name of the privileged portal.
	Name string `json:"name,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, consoleID string) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetPraPortal(ctx context.Context, service *zscaler.Service, portalID string) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.GetCustomerID()+praConsoleEndpoint+"/praPortal", portalID)
	resp, err := service.Client.NewRequestDo(ctx, "GET", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, consoleName string) (*PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAConsole](ctx, service.Client, relativeURL, common.Filter{Search: consoleName, MicroTenantID: service.MicroTenantID()})
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

func Create(ctx context.Context, service *zscaler.Service, praConsole *PRAConsole) (*PRAConsole, *http.Response, error) {
	v := new(PRAConsole)
	resp, err := service.Client.NewRequestDo(ctx, "POST", mgmtConfig+service.Client.GetCustomerID()+praConsoleEndpoint, common.Filter{MicroTenantID: service.MicroTenantID()}, praConsole, &v)
	if err != nil {
		return nil, nil, err
	}
	return v, resp, nil
}

func CreatePraBulk(ctx context.Context, service *zscaler.Service, praConsoles []PRAConsole) ([]PRAConsole, *http.Response, error) {
	var responseConsoles []PRAConsole
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + praConsoleBulkEndpoint
	resp, err := service.Client.NewRequestDo(ctx, "POST", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, praConsoles, &responseConsoles)
	if err != nil {
		return nil, nil, err
	}
	return responseConsoles, resp, nil
}

func Update(ctx context.Context, service *zscaler.Service, consoleID string, praConsole *PRAConsole) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo(ctx, "PUT", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, praConsole, nil)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func Delete(ctx context.Context, service *zscaler.Service, consoleID string) (*http.Response, error) {
	relativeURL := fmt.Sprintf("%v/%v", mgmtConfig+service.Client.GetCustomerID()+praConsoleEndpoint, consoleID)
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()}, nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]PRAConsole, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + praConsoleEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[PRAConsole](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
