package shadow_it_report

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/common"
)

const (
	customTagsEndpoint                  = "/customTags"
	cloudApplicationsLiteEndpoint       = "/cloudApplications/lite"
	cloudApplicationsBulkDeleteEndpoint = "/cloudApplications/bulkUpdate"
)

type CloudApplications struct {
	// The cloud application status that indicates whether it is sanctioned or unsanctioned
	SanctionedState string `json:"sanctionedState,omitempty"`

	// The list of cloud application IDs for which the status (sanctioned or unsanctioned) and tags have to be updated
	ApplicationIDs []int `json:"applicationIds,omitempty"`

	// The list of custom tags that must be assigned to the cloud applications
	CustomTags []ShadowITReportNameID `json:"customTags,omitempty"`
}

// Gets the list of predefined and custom cloud applications or custom tags
type ShadowITReportNameID struct {
	// Unique identifier of the cloud application or custom tags
	ID int `json:"id,omitempty"`

	// The name of the cloud application or custom tags
	Name string `json:"name,omitempty"`
}

func (service *Service) GetCloudAppLite(cloudAppID int) (*ShadowITReportNameID, error) {
	var cloudApp ShadowITReportNameID
	err := service.Client.Read(fmt.Sprintf("%s/%d", cloudApplicationsLiteEndpoint, cloudAppID), &cloudApp)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning cloud application from Get: %d", cloudApp.ID)
	return &cloudApp, nil
}

func (service *Service) GetCloudAppByName(cloudAppName string) (*ShadowITReportNameID, error) {
	var cloudApps []ShadowITReportNameID
	err := common.ReadAllPages(service.Client, cloudApplicationsLiteEndpoint, &cloudApps)
	if err != nil {
		return nil, err
	}
	for _, cloudApp := range cloudApps {
		if strings.EqualFold(cloudApp.Name, cloudAppName) {
			return &cloudApp, nil
		}
	}
	return nil, fmt.Errorf("no cloud application found with name: %s", cloudAppName)
}

func (service *Service) GetAllCloudApp() ([]ShadowITReportNameID, error) {
	var cloudApps []ShadowITReportNameID
	err := common.ReadAllPages(service.Client, cloudApplicationsLiteEndpoint, &cloudApps)
	return cloudApps, err
}

func (service *Service) GetCustomTags(cloudAppID int) (*ShadowITReportNameID, error) {
	var cloudApp ShadowITReportNameID
	err := service.Client.Read(fmt.Sprintf("%s/%d", customTagsEndpoint, cloudAppID), &cloudApp)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG] Returning custom tag from Get: %d", cloudApp.ID)
	return &cloudApp, nil
}

func (service *Service) GetCustomTagsByName(tagName string) (*ShadowITReportNameID, error) {
	var customTags []ShadowITReportNameID
	err := common.ReadAllPages(service.Client, customTagsEndpoint, &customTags)
	if err != nil {
		return nil, err
	}
	for _, customTag := range customTags {
		if strings.EqualFold(customTag.Name, tagName) {
			return &customTag, nil
		}
	}
	return nil, fmt.Errorf("no custom tag found with name: %s", tagName)
}

func (service *Service) GetAllCustomTags() ([]ShadowITReportNameID, error) {
	var cloudApps []ShadowITReportNameID
	err := common.ReadAllPages(service.Client, cloudApplicationsLiteEndpoint, &cloudApps)
	return cloudApps, err
}

func (service *Service) CloudAppBulkUpdate(cloudApps *CloudApplications) (*CloudApplications, *http.Response, error) {
	respContent, err := service.Client.UpdateWithPut(cloudApplicationsBulkDeleteEndpoint, cloudApps)
	if err != nil {
		// We can't get an *http.Response from the error, so just return the error.
		return nil, nil, err
	}

	// Attempt to unmarshal the response content into a CloudApplications struct
	var updatedCloudApps CloudApplications
	err = json.Unmarshal(respContent.([]byte), &updatedCloudApps) // Assuming respContent is a []byte containing JSON
	if err != nil {
		// Handle the case where unmarshalling fails
		service.Client.Logger.Printf("[ERROR] Error unmarshalling response: %v", err)
		return nil, nil, err
	}

	// Success case: we return the unmarshalled *CloudApplications, no HTTP response is available here
	return &updatedCloudApps, nil, nil
}
