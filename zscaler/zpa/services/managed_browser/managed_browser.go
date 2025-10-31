package managed_browser

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig             = "/zpa/mgmtconfig/v1/admin/customers/"
	managedBrowserEndpoint = "/managedBrowserProfile/search"
)

type ManagedBrowserProfile struct {
	BrowserType          string               `json:"browserType,omitempty"`
	CreationTime         string               `json:"creationTime,omitempty"`
	CustomerID           string               `json:"customerId,omitempty"`
	Description          string               `json:"description,omitempty"`
	ID                   string               `json:"id,omitempty"`
	ModifiedBy           string               `json:"modifiedBy,omitempty"`
	ModifiedTime         string               `json:"modifiedTime,omitempty"`
	Name                 string               `json:"name,omitempty"`
	MicrotenantID        string               `json:"microtenantId,omitempty"`
	MicrotenantName      string               `json:"microtenantName,omitempty"`
	ChromePostureProfile ChromePostureProfile `json:"chromePostureProfile,omitempty"`
}

type ChromePostureProfile struct {
	ID               string `json:"id,omitempty"`
	BrowserType      string `json:"browserType,omitempty"`
	CrowdStrikeAgent bool   `json:"crowdStrikeAgent,omitempty"`
	CreationTime     string `json:"creationTime,omitempty"`
	ModifiedBy       string `json:"modifiedBy,omitempty"`
	ModifiedTime     string `json:"modifiedTime,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]ManagedBrowserProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + managedBrowserEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ManagedBrowserProfile](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, managedBrowserName string) (*ManagedBrowserProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + managedBrowserEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[ManagedBrowserProfile](ctx, service.Client, relativeURL, common.Filter{MicroTenantID: service.MicroTenantID()})
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, managedBrowserName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no managed browser profile named '%s' was found", managedBrowserName)
}
