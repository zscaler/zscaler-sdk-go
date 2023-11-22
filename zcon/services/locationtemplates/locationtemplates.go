package locationtemplates

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v2/zcon/services/common"
)

const (
	locationTemplateEndpoint = "/locationTemplate"
)

type LocationTemplate struct {
	//ID of Cloud & Branch Connector location template
	ID int `json:"id,omitempty"`

	// Name of Cloud & Branch Connector location template
	Name string `json:"name,omitempty"`

	// Description of Cloud & Branch Connector location template
	Description string `json:"desc,omitempty"`

	// Details of Cloud & Branch Connector location template
	LocationTemplateDetails *LocationTemplateDetails `json:"template,omitempty"`

	// Whether Cloud & Branch Connector location template is editable
	Editable bool `json:"editable"`

	// Last time Cloud & Branch Connector location template was modified
	LastModTime int `json:"lastModTime,omitempty"`

	// User ID of last time Cloud & Branch Connector location template was modified
	LastModUid *common.UIDName `json:"lastModUid,omitempty"`
}

type LocationTemplateDetails struct {
	//Prefix of Cloud & Branch Connector location template
	TemplatePrefix string `json:"templatePrefix,omitempty"`

	// Enable if you want the Zscaler service to use the X-Forwarded-For (XFF) headers that your on-premise proxy server inserts in outbound HTTP requests
	XFFForwardEnabled bool `json:"xffForwardEnabled,omitempty"`

	// Indicates whether "Authentication Required" is enabled in the location template
	AuthRequired bool `json:"authRequired,omitempty"`

	// Enable to display an end user notification for unauthenticated traffic. If disabled, the action is treated as an allow-policy
	CautionEnabled bool `json:"cautionEnabled,omitempty"`

	// Enable this feature to display an Acceptable Use Policy (AUP) for unauthenticated traffic and require users to accept it
	AupEnabled bool `json:"aupEnabled,omitempty"`

	// If you enabled aupEnabled, specify in days how frequently the AUP is displayed to users
	AupTimeoutInDays int `json:"aupTimeoutInDays,omitempty"`

	// Enables the service's firewall controls
	OFWEnabled bool `json:"ofwEnabled"`

	// If you enabled ofwEnabled, enable this for IPS controls for the location template
	IPSControl bool `json:"ipsControl"`

	// Enable to specify the maximum bandwidth limits for download (Mbps) and upload (Mbps)
	EnforceBandwidthControl bool `json:"enforceBandwidthControl"`

	// If you enabled enforceBandwidthControl, specify the maximum bandwidth for upload (Mbps)
	UpBandwidth int `json:"upBandwidth,omitempty"`

	// If you enabled enforceBandwidthControl, specify the maximum bandwidth for download (Mbps)
	DnBandwidth int `json:"dnBandwidth,omitempty"`

	// Display Time Unit. The time unit to display for IP Surrogate idle time to disassociation
	// Support values are: "MINUTE", "HOUR", "DAY"
	DisplayTimeUnit string `json:"displayTimeUnit,omitempty"`

	// Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled
	IdleTimeInMinutes int `json:"idleTimeInMinutes,omitempty"`

	// Idle Time to Disassociation. The user mapping idle time (in minutes) is required if a Surrogate IP is enabled
	SurrogateIPEnforcedForKnownBrowsers bool `json:"surrogateIPEnforcedForKnownBrowsers,omitempty"`

	// Display Refresh Time Unit. The time unit to display for refresh time for re-validation of surrogacy
	// Support values are: "MINUTE", "HOUR", "DAY"
	SurrogateRefreshTimeUnit string `json:"surrogateRefreshTimeUnit,omitempty"`

	// Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates
	SurrogateRefreshTimeInMinutes int `json:"surrogateRefreshTimeInMinutes,omitempty"`

	// Refresh Time for re-validation of Surrogacy. The surrogate refresh time (in minutes) to re-validate the IP surrogates
	SurrogateIP bool `json:"surrogateIP,omitempty"`
}

func (service *Service) Get(locTemplateID int) (*LocationTemplate, error) {
	var location LocationTemplate
	err := service.Client.Read(fmt.Sprintf("%s/%d", locationTemplateEndpoint, locTemplateID), &location)
	if err != nil {
		return nil, err
	}

	log.Printf("Returning Location Template from Get: %d", location.ID)
	return &location, nil
}

func (service *Service) GetByName(templateName string) (*LocationTemplate, error) {
	var locations []LocationTemplate
	// We are assuming this location name will be in the firsy 1000 obejcts
	err := common.ReadAllPages(service.Client, locationTemplateEndpoint, &locations)
	if err != nil {
		return nil, err
	}
	for _, location := range locations {
		if strings.EqualFold(location.Name, templateName) {
			return &location, nil
		}
	}
	return nil, fmt.Errorf("no location template found with name: %s", templateName)
}

func (service *Service) Create(locations *LocationTemplate) (*LocationTemplate, error) {
	resp, err := service.Client.Create(locationTemplateEndpoint, *locations)
	if err != nil {
		return nil, err
	}

	createdLocations, ok := resp.(*LocationTemplate)
	if !ok {
		return nil, errors.New("object returned from api was not a location template pointer")
	}

	log.Printf("returning location template from create: %d", createdLocations.ID)
	return createdLocations, nil
}

func (service *Service) Update(locTemplateID int, locations *LocationTemplate) (*LocationTemplate, *http.Response, error) {
	resp, err := service.Client.UpdateWithPut(fmt.Sprintf("%s/%d", locationTemplateEndpoint, locTemplateID), *locations)
	if err != nil {
		return nil, nil, err
	}
	updatedLocations, _ := resp.(*LocationTemplate)

	log.Printf("returning location template from Update: %d", updatedLocations.ID)
	return updatedLocations, nil, nil
}

func (service *Service) Delete(locTemplateID int) (*http.Response, error) {
	err := service.Client.Delete(fmt.Sprintf("%s/%d", locationTemplateEndpoint, locTemplateID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (service *Service) GetAll() ([]LocationTemplate, error) {
	var templates []LocationTemplate
	err := common.ReadAllPages(service.Client, locationTemplateEndpoint, &templates)
	return templates, err
}
