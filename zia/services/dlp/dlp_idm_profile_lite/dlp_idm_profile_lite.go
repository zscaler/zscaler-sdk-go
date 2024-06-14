package dlp_idm_profile_lite

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	dlpIDMProfileLiteEndpoint = "/idmprofile/lite"
)

// Gets a list of active IDM templates (or IDM profiles) and their criteria, only.
type DLPIDMProfileLite struct {
	// The identifier (1-64) for the IDM template (i.e., IDM profile) that is unique within the organization.
	ProfileID int `json:"profileId,omitempty"`

	// The IDM template name.
	TemplateName string `json:"templateName,omitempty"`

	// The name of the Index Tool virtual machine (VM) that the IDM template belongs to.
	ClientVM *common.IDNameExtensions `json:"clientVm,omitempty"`

	// The name of the Index Tool virtual machine (VM) that the IDM template belongs to.
	NumDocuments int `json:"numDocuments,omitempty"`

	// The date and time the IDM template was last modified.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// The admin that modified the IDM template last.
	ModifiedBy *common.IDNameExtensions `json:"modifiedBy,omitempty"`
}

func GetDLPProfileLiteID(service *services.Service, ProfileLiteID int, activeOnly bool) (*DLPIDMProfileLite, error) {
	endpoint := fmt.Sprintf("%s/%d", dlpIDMProfileLiteEndpoint, ProfileLiteID)
	if activeOnly {
		endpoint += "?activeOnly=true"
	}

	var profileLite DLPIDMProfileLite
	err := service.Client.Read(endpoint, &profileLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning idm profile template from Get: %d", profileLite.ProfileID)
	return &profileLite, nil
}

func GetDLPProfileLiteByName(service *services.Service, profileLiteName string, activeOnly bool) (*DLPIDMProfileLite, error) {
	queryParameters := url.Values{}
	queryParameters.Set("name", profileLiteName)
	if activeOnly {
		queryParameters.Set("activeOnly", "true")
	}

	endpoint := fmt.Sprintf("%s?%s", dlpIDMProfileLiteEndpoint, queryParameters.Encode())
	var profileLite []DLPIDMProfileLite
	err := common.ReadAllPages(service.Client, endpoint, &profileLite)
	if err != nil {
		return nil, err
	}
	for _, profile := range profileLite {
		if strings.EqualFold(profile.TemplateName, profileLiteName) {
			return &profile, nil
		}
	}
	return nil, fmt.Errorf("no idm profile template found with name: %s", profileLiteName)
}

func GetAll(service *services.Service, activeOnly bool) ([]DLPIDMProfileLite, error) {
	endpoint := dlpIDMProfileLiteEndpoint
	if activeOnly {
		endpoint += "?activeOnly=true"
	}

	var idmpProfile []DLPIDMProfileLite
	err := common.ReadAllPages(service.Client, endpoint, &idmpProfile)
	return idmpProfile, err
}
