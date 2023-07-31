package isolationprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa/services/common"
)

const (
	mgmtConfig               = "/mgmtconfig/v1/admin/customers/"
	isolationProfileEndpoint = "/isolation/profiles"
)

type IsolationProfile struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	Enabled            bool   `json:"enabled"`
	CreationTime       string `json:"creationTime,omitempty"`
	ModifiedBy         string `json:"modifiedBy,omitempty"`
	ModifiedTime       string `json:"modifiedTime,omitempty"`
	IsolationProfileID string `json:"isolationProfileId,omitempty"`
	IsolationTenantID  string `json:"isolationTenantId,omitempty"`
	IsolationURL       string `json:"isolationUrl"`
}

func (service *Service) Get(profileID string) (*IsolationProfile, *http.Response, error) {
	v := new(IsolationProfile)
	relativeURL := fmt.Sprintf("%s/%s", mgmtConfig+service.Client.Config.CustomerID+isolationProfileEndpoint, profileID)
	resp, err := service.Client.NewRequestDo("GET", relativeURL, common.Filter{MicroTenantID: service.microTenantID}, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

func (service *Service) GetByName(profileName string) (*IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + isolationProfileEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IsolationProfile](service.Client, relativeURL, common.Filter{Search: profileName, MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	for _, profile := range list {
		if strings.EqualFold(profile.Name, profileName) {
			return &profile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no isolation profile named '%s' was found", profileName)
}

func (service *Service) GetAll() ([]IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + isolationProfileEndpoint
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IsolationProfile](service.Client, relativeURL, common.Filter{MicroTenantID: service.microTenantID})
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
