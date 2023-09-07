package dlp_idm_profiles

import (
	"fmt"
	"net/url"
	"strings"

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

func (service *Service) GetDLPProfileLiteID(ProfileLiteID int) (*DLPIDMProfileLite, error) {
	var profileLite DLPIDMProfileLite
	err := service.Client.Read(fmt.Sprintf("%s/%d", dlpIDMProfileLiteEndpoint, ProfileLiteID), &profileLite)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]returning idm profile template from Get: %d", profileLite.ProfileID)
	return &profileLite, nil
}

func (service *Service) GetDLPProfileLiteByName(profileLiteName string) (*DLPIDMProfileLite, error) {
	var profileLite []DLPIDMProfileLite
	err := common.ReadAllPages(service.Client, fmt.Sprintf("%s?name=%s", dlpIDMProfileLiteEndpoint, url.QueryEscape(profileLiteName)), &profileLite)
	if err != nil {
		return nil, err
	}
	for _, profileLite := range profileLite {
		if strings.EqualFold(profileLite.TemplateName, profileLiteName) {
			return &profileLite, nil
		}
	}
	return nil, fmt.Errorf("no idm profile template found with name: %s", profileLiteName)
}
