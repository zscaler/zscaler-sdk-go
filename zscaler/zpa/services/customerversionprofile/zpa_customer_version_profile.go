package customerversionprofile

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig                     = "/zpa/mgmtconfig/v1/admin/customers/"
	customerVersionProfileEndpoint = "/visible/versionProfiles"
)

type CustomerVersionProfile struct {
	CreationTime                   string                        `json:"creationTime,omitempty"`
	CustomScopeCustomerIDs         []CustomScopeCustomerIDs      `json:"customScopeCustomerIds"`
	CustomScopeRequestCustomerIDs  CustomScopeRequestCustomerIDs `json:"customScopeRequestCustomerIds"`
	CustomerID                     string                        `json:"customerId"`
	Description                    string                        `json:"description"`
	ID                             string                        `json:"id,omitempty"`
	ModifiedBy                     string                        `json:"modifiedBy"`
	ModifiedTime                   string                        `json:"modifiedTime"`
	Name                           string                        `json:"name"`
	Versions                       []Versions                    `json:"versions"`
	VisibilityScope                string                        `json:"visibilityScope"`
	UpgradePriority                string                        `json:"upgradePriority"`
	NumberOfAssistants             string                        `json:"numberOfAssistants"`
	NumberOfCustomers              string                        `json:"numberOfCustomers"`
	NumberOfPrivateBrokers         string                        `json:"numberOfPrivateBrokers"`
	NumberOfSiteControllers        string                        `json:"numberOfSiteControllers"`
	NumberOfUpdatedAssistants      string                        `json:"numberOfUpdatedAssistants"`
	NumberOfUpdatedPrivateBrokers  string                        `json:"numberOfUpdatedPrivateBrokers"`
	NumberOfUpdatedSiteControllers string                        `json:"numberOfUpdatedSiteControllers"`
}

type CustomScopeCustomerIDs struct {
	Name                 string `json:"name"`
	CustomerID           string `json:"customerId"`
	ExcludeConstellation bool   `json:"excludeConstellation"`
	IsPartner            bool   `json:"isPartner"`
}

type CustomScopeRequestCustomerIDs struct {
	AddCustomerIDs    string `json:"addCustomerIds"`
	DeletecustomerIDs string `json:"deleteCustomerIds"`
}

type Versions struct {
	ID                       string `json:"id,omitempty"`
	CustomerID               string `json:"customerId"`
	CreationTime             string `json:"creationTime,omitempty"`
	ModifiedBy               string `json:"modifiedBy"`
	ModifiedTime             string `json:"modifiedTime"`
	Platform                 string `json:"platform"`
	RestartAfterUptimeInDays string `json:"restartAfterUptimeInDays"`
	Role                     string `json:"role"`
	Version                  string `json:"version"`
	VersionProfileGID        string `json:"version_profile_gid"`
}

func GetByName(ctx context.Context, service *zscaler.Service, versionProfileName string) (*CustomerVersionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + customerVersionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[CustomerVersionProfile](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, versionProfileName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no version profile named '%s' was found", versionProfileName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]CustomerVersionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + customerVersionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[CustomerVersionProfile](ctx, service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
