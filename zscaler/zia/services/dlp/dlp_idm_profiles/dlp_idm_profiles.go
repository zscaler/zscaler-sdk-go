package dlp_idm_profiles

import (
	"context"
	"fmt"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpIDMProfileEndpoint = "/zia/api/v1/idmprofile"
)

type DLPIDMProfile struct {
	// The identifier (1-64) for the IDM template (i.e., IDM profile) that is unique within the organization.
	ProfileID int `json:"profileId,omitempty"`

	// The IDM template name, which is unique per Index Tool.
	ProfileName string `json:"profileName,omitempty"`

	// The IDM template's description.
	ProfileDesc string `json:"profileDesc,omitempty"`

	// The IDM template's type. Supported values are: "LOCAL", "REMOTECRON", and "REMOTE"
	ProfileType string `json:"profileType,omitempty"`

	// The fully qualified domain name (FQDN) of the IDM template's host machine.
	Host string `json:"host,omitempty"`

	// The port number of the IDM template's host machine.
	Port int `json:"port,omitempty"`

	// The IDM template's directory file path, where all files are present.
	ProfileDirPath string `json:"profileDirPath,omitempty"`

	// The schedule type for the IDM template's schedule (i.e., Monthly, Weekly, Daily, or None). This attribute is required by PUT and POST requests.
	ScheduleType string `json:"scheduleType,omitempty"`

	// The day the IDM template is scheduled for. This attribute is required by PUT and POST requests.
	ScheduleDay int `json:"scheduleDay,omitempty"`

	// The day of the month that the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to MONTHLY.
	ScheduleDayOfMonth []string `json:"scheduleDayOfMonth,omitempty"`

	// The day of the week the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to WEEKLY.
	ScheduleDayOfWeek []string `json:"scheduleDayOfWeek,omitempty"`

	// The time of the day (in minutes) that the IDM template is scheduled for. For example: at 3am= 180 mins. This attribute is required by PUT and POST requests.
	ScheduleTime int `json:"scheduleTime,omitempty"`

	// If set to true, the schedule for the IDM template is temporarily in a disabled state. This attribute is required by PUT requests in order to disable or enable a schedule.
	ScheduleDisabled bool `json:"scheduleDisabled,omitempty"`

	// The status of the file uploaded to the Index Tool for the IDM template.
	UploadStatus string `json:"uploadStatus"`

	// The username to be used on the IDM template's host machine.
	UserName string `json:"userName,omitempty"`

	// The version number for the IDM template.
	Version int `json:"version,omitempty"`

	// The unique identifer for the Index Tool that was used to create the IDM template. This attribute is required by POST requests, but ignored if provided in PUT requests.
	IDMClient *common.IDNameExtensions `json:"idmClient,omitempty"`

	// The total volume of all the documents associated to the IDM template.
	VolumeOfDocuments int `json:"volumeOfDocuments,omitempty"`

	// The number of documents associated to the IDM template.
	NumDocuments int `json:"numDocuments,omitempty"`

	// The date and time the IDM template was last modified.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// The admin that modified the IDM template last.
	ModifiedBy *common.IDNameExtensions `json:"modifiedBy,omitempty"`
}

func Get(ctx context.Context, service *zscaler.Service, idmProfileID int) (*DLPIDMProfile, error) {
	var idmpProfile DLPIDMProfile
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpIDMProfileEndpoint, idmProfileID), &idmpProfile)
	if err != nil {
		return nil, err
	}

	service.Client.Logger.Printf("[DEBUG]Returning dlp icap server from Get: %d", idmpProfile.ProfileID)
	return &idmpProfile, nil
}

func GetByName(ctx context.Context, service *zscaler.Service, idmProfileName string) (*DLPIDMProfile, error) {
	var idmpProfile []DLPIDMProfile
	err := common.ReadAllPages(ctx, service.Client, dlpIDMProfileEndpoint, &idmpProfile)
	if err != nil {
		return nil, err
	}
	for _, icap := range idmpProfile {
		if strings.EqualFold(icap.ProfileName, idmProfileName) {
			return &icap, nil
		}
	}
	return nil, fmt.Errorf("no dlp icap server found with name: %s", idmProfileName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DLPIDMProfile, error) {
	var idmpProfile []DLPIDMProfile
	err := common.ReadAllPages(ctx, service.Client, dlpIDMProfileEndpoint, &idmpProfile)
	return idmpProfile, err
}
