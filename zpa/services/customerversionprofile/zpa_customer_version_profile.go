package customerversionprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

const (
	mgmtConfig                     = "/mgmtconfig/v1/admin/customers/"
	customerVersionProfileEndpoint = "/visible/versionProfiles"
)

type CustomerVersionProfile struct {
	CreationTime                  string                        `json:"creationTime,omitempty"`
	CustomScopeCustomerIDs        []CustomScopeCustomerIDs      `json:"customScopeCustomerIds"`
	CustomScopeRequestCustomerIDs CustomScopeRequestCustomerIDs `json:"customScopeRequestCustomerIds"`
	CustomerID                    string                        `json:"customerId"`
	Description                   string                        `json:"description"`
	ID                            string                        `json:"id,omitempty"`
	ModifiedBy                    string                        `json:"modifiedBy"`
	ModifiedTime                  string                        `json:"modifiedTime"`
	Name                          string                        `json:"name"`
	UpgradePriority               string                        `json:"upgradePriority"`
	Versions                      []Versions                    `json:"versions"`
	VisibilityScope               string                        `json:"visibilityScope"`
}

type CustomScopeCustomerIDs struct {
	CustomerID           string `json:"customerId"`
	ExcludeConstellation bool   `json:"excludeConstellation"`
	Name                 string `json:"name"`
}

type CustomScopeRequestCustomerIDs struct {
	AddCustomerIDs    string `json:"addCustomerIds"`
	DeletecustomerIDs string `json:"deleteCustomerIds"`
}

type Versions struct {
	CreationTime             string `json:"creationTime,omitempty"`
	CustomerID               string `json:"customerId"`
	ID                       string `json:"id,omitempty"`
	ModifiedBy               string `json:"modifiedBy"`
	ModifiedTime             string `json:"modifiedTime"`
	Platform                 string `json:"platform"`
	RestartAfterUptimeInDays string `json:"restartAfterUptimeInDays"`
	Role                     string `json:"role"`
	Version                  string `json:"version"`
	VersionProfileGID        string `json:"version_profile_gid"`
}

func GetByName(service *services.Service, versionProfileName string) (*CustomerVersionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + customerVersionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[CustomerVersionProfile](service.Client, relativeURL, "")
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

func GetAll(service *services.Service) ([]CustomerVersionProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + customerVersionProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[CustomerVersionProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
