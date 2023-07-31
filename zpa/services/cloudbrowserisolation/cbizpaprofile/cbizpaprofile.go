package cbizpaprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	cbiConfig                 = "/cbiconfig/cbi/api/customers/"
	zpaProfileEndpoint string = "/zpaprofiles"
)

type ZPAProfiles struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	Enabled      bool   `json:"enabled"`
	CreationTime string `json:"creationTime,omitempty"`
	ModifiedBy   string `json:"modifiedBy,omitempty"`
	ModifiedTime string `json:"modifiedTime,omitempty"`
	CBITenantID  string `json:"cbiTenantId,omitempty"`
	CBIProfileID string `json:"cbiProfileId,omitempty"`
	CBIURL       string `json:"cbiUrl"`
}

// The current API does not seem to support search by ID
func (service *Service) Get(profileID string) (*ZPAProfiles, *http.Response, error) {
	v := new(ZPAProfiles)
	relativeURL := fmt.Sprintf("%s/%s", cbiConfig+service.Client.Config.CustomerID+zpaProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// The current API does not seem to support search by Name
func (service *Service) GetByName(profileName string) (*ZPAProfiles, *http.Response, error) {
	list, resp, err := service.GetAll()
	if err != nil {
		return nil, nil, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, profileName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no zpa profile named '%s' was found", profileName)
}

func (service *Service) GetAll() ([]ZPAProfiles, *http.Response, error) {
	relativeURL := cbiConfig + service.Client.Config.CustomerID + zpaProfileEndpoint
	var list []ZPAProfiles
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
