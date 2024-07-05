package cbizpaprofile

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services"
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

type QueryParams struct {
	ShowDisabled string `url:"showDisabled,omitempty"`
	ScopeId      string `url:"scopeId,omitempty"`
}

// Get retrieves a profile by its ID. This function now uses GetAll with optional parameters correctly.
func Get(service *services.Service, profileID string) (*ZPAProfiles, *http.Response, error) {
	// Using nil for optional parameters as defaults
	profiles, resp, err := GetAll(service, nil, nil)
	if err != nil {
		return nil, resp, err
	}

	for _, profile := range profiles {
		if profile.ID == profileID {
			return &profile, resp, nil
		}
	}

	return nil, resp, fmt.Errorf("no isolation profile with ID '%s' was found", profileID)
}

// GetByName retrieves a profile by name. This function now uses GetAll with optional parameters correctly.
func GetByName(service *services.Service, profileName string) (*ZPAProfiles, *http.Response, error) {
	// Using nil for optional parameters as defaults
	list, resp, err := GetAll(service, nil, nil)
	if err != nil {
		return nil, nil, err
	}
	for _, profile := range list {
		if strings.EqualFold(profile.Name, profileName) {
			return &profile, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no zpa profile named '%s' was found", profileName)
}

// GetAll retrieves all profiles, with optional parameters to show disabled profiles and filter by scopeId.
func GetAll(service *services.Service, showDisabled *bool, scopeId *int) ([]ZPAProfiles, *http.Response, error) {
	relativeURL := fmt.Sprintf("%s%s%s", cbiConfig, service.Client.Config.CustomerID, zpaProfileEndpoint)

	// Prepare query parameters using a struct
	params := QueryParams{}
	if showDisabled != nil {
		params.ShowDisabled = strconv.FormatBool(*showDisabled)
	}
	if scopeId != nil {
		params.ScopeId = strconv.Itoa(*scopeId)
	}

	var list []ZPAProfiles
	resp, err := service.Client.NewRequestDo("GET", relativeURL, params, nil, &list)
	if err != nil {
		return nil, nil, err
	}
	return list, resp, nil
}
