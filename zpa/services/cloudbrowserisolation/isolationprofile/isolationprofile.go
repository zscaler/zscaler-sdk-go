package isolationprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
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
	// First get all the profiles
	profiles, resp, err := service.GetAll()
	if err != nil {
		return nil, resp, err
	}

	// Loop through the profiles and find the one with the matching ID
	for _, profile := range profiles {
		if profile.ID == profileID {
			return &profile, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("no isolation profile with ID '%s' was found", profileID)
}

func (service *Service) GetByName(profileName string) (*IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + isolationProfileEndpoint

	// Set up custom filters for pagination
	filters := common.Filter{Search: profileName} // We only have the Search filter as per your example. You can add more filters if required.

	// Use the custom pagination function
	list, resp, err := common.GetAllPagesGenericWithCustomFilters[IsolationProfile](service.Client, relativeURL, filters)
	if err != nil {
		return nil, nil, err
	}

	// The rest remains the same as your logic for finding the profile by its name
	for _, profile := range list {
		if strings.EqualFold(profile.Name, profileName) {
			return &profile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no isolation profile named '%s' was found", profileName)
}

func (service *Service) GetAll() ([]IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.Config.CustomerID + isolationProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[IsolationProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
