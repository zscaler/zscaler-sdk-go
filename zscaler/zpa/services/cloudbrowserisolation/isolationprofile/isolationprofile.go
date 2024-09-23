package isolationprofile

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
)

const (
	mgmtConfig               = "/zpa/mgmtconfig/v1/admin/customers/"
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

func GetByName(service *zscaler.Service, profileName string) (*IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + isolationProfileEndpoint

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

func GetAll(service *zscaler.Service) ([]IsolationProfile, *http.Response, error) {
	relativeURL := mgmtConfig + service.Client.GetCustomerID() + isolationProfileEndpoint
	list, resp, err := common.GetAllPagesGeneric[IsolationProfile](service.Client, relativeURL, "")
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
