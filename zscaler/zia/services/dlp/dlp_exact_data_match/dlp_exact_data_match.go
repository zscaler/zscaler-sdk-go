package dlp_exact_data_match

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	dlpEDMSchemaEndpoint = "/zia/api/v1/dlpExactDataMatchSchemas"
)

type DLPEDMSchema struct {
	// The identifier (1-65519) for the EDM schema (i.e., EDM template) that is unique within the organization.
	SchemaID int `json:"schemaId,omitempty"`

	// The unique identifer for the Index Tool that was used to create the EDM template. This attribute is ignored by PUT requests, but required for POST requests.
	EDMClient *common.IDNameExtensions `json:"edmClient,omitempty"`

	// The EDM schema (i.e., EDM template) name. This attribute is ignored by PUT requests, but required for POST requests.
	ProjectName string `json:"projectName,omitempty"`

	// The revision number of the CSV file upload to the Index Tool. This attribute is required by PUT requests.
	Revision int `json:"revision,omitempty"`

	// The generated filename, excluding the extention.
	Filename string `json:"filename,omitempty"`

	// The generated filename, excluding the extention.
	OriginalFileName string `json:"originalFileName,omitempty"`

	// The status of the EDM template's CSV file upload to the Index Tool. This attribute is required by PUT and POST requests.
	FileUploadStatus string `json:"fileUploadStatus,omitempty"`

	// The status of the EDM template's CSV file upload to the Index Tool. This attribute is required by PUT and POST requests.
	SchemaStatus string `json:"schemaStatus,omitempty"`

	// The total count of actual columns selected from the CSV file. This attribute is required by PUT and POST requests.
	OrigColCount int `json:"origColCount,omitempty"`

	// Timestamp when the EDM schema (i.e., EDM template) was last modified. Ignored if the request is PUT, POST, or DELETE.
	LastModifiedTime int `json:"lastModifiedTime,omitempty"`

	// The admin that modified the EDM template's schema last.
	ModifiedBy *common.IDNameExtensions `json:"modifiedBy,omitempty"`

	// The login name (or userid) the admin who created the EDM schema (i.e., EDM template). Ignored if the request is PUT, POST, or DELETE.
	CreatedBy *common.IDNameExtensions `json:"createdBy,omitempty"`

	// The total number of cells used by the EDM schema (i.e., EDM template).
	CellsUsed int `json:"cellsUsed,omitempty"`

	// Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.
	SchemaActive bool `json:"schemaActive,omitempty"`

	// The total number of cells used by the EDM schema (i.e., EDM template).
	SchedulePresent bool `json:"schedulePresent,omitempty"`

	// Indicates the status of a specified EDM schema (i.e., EDM template). If this value is set to true, the schema is active and can be used by DLP dictionaries.
	TokenList []TokenList `json:"tokenList,omitempty"`

	// The schedule details, if present for the EDM schema (i.e., EDM template). Ignored if the request is PUT, POST, or DELETE.
	Schedule Schedule `json:"schedule,omitempty"`
}

type TokenList struct {
	// The token (i.e., criteria) name. This attribute is required by PUT and POST requests.
	Name string `json:"name,omitempty"`

	// The token (i.e., criteria) name. This attribute is required by PUT and POST requests.
	Type string `json:"type,omitempty"`

	// Indicates whether the token is a primary key.
	PrimaryKey bool `json:"primaryKey,omitempty"`

	// The column position for the token in the original CSV file uploaded to the Index Tool, starting from 1. This attribue required by PUT and POST requests.
	OriginalColumn int `json:"originalColumn,omitempty"`

	// The column position for the token in the hashed file, starting from 1.
	HashfileColumnOrder int `json:"hashfileColumnOrder,omitempty"`

	// The length of the column bitmap in the hashed file.
	ColLengthBitmap int `json:"colLengthBitmap,omitempty"`
}

type Schedule struct {
	// The schedule type for the EDM schema (i.e., EDM template), Monthly, Weekly, Daily, or None. This attribute is required by PUT and POST requests.
	ScheduleType string `json:"scheduleType,omitempty"`

	// The day of the month the EDM schema (i.e., EDM template) is scheduled for. This attribute is required by PUT and POST requests, and if the scheduleType is set to MONTHLY.
	ScheduleDayOfMonth []string `json:"scheduleDayOfMonth,omitempty"`

	// The day of the week the EDM schema (i.e., EDM template) is scheduled for. This attribute is required by PUT and POST requests, and if the scheduleType is set to WEEKLY.
	ScheduleDayOfWeek []string `json:"scheduleDayOfWeek,omitempty"`

	// The time of the day (in minutes) that the EDM schema (i.e., EDM template) is scheduled for. For example: at 3am= 180 mins. This attribute is required by PUT and POST requests.
	ScheduleTime int `json:"scheduleTime,omitempty"`

	// If set to true, the schedule for the EDM schema (i.e., EDM template) is temporarily in a disabled state. This attribute is required by PUT requests in order to disable or enable a schedule.
	ScheduleDisabled bool `json:"scheduleDisabled,omitempty"`
}

func GetDLPEDMSchemaID(ctx context.Context, service *zscaler.Service, edmSchemaID int) (*DLPEDMSchema, error) {
	var edmSchema DLPEDMSchema
	err := service.Client.Read(ctx, fmt.Sprintf("%s/%d", dlpEDMSchemaEndpoint, edmSchemaID), &edmSchema)
	if err != nil {
		return nil, err
	}

	service.Client.GetLogger().Printf("[DEBUG]returning edm schema from Get: %d", edmSchema.SchemaID)
	return &edmSchema, nil
}

func GetDLPEDMByName(ctx context.Context, service *zscaler.Service, edmSchemaName string) (*DLPEDMSchema, error) {
	var edmSchema []DLPEDMSchema
	err := common.ReadAllPages(ctx, service.Client, fmt.Sprintf("%s?name=%s", dlpEDMSchemaEndpoint, url.QueryEscape(edmSchemaName)), &edmSchema)
	if err != nil {
		return nil, err
	}
	for _, edmSchema := range edmSchema {
		if strings.EqualFold(edmSchema.ProjectName, edmSchemaName) {
			return &edmSchema, nil
		}
	}
	return nil, fmt.Errorf("no edm schema found with name: %s", edmSchemaName)
}

func GetAll(ctx context.Context, service *zscaler.Service) ([]DLPEDMSchema, error) {
	var edmData []DLPEDMSchema
	err := common.ReadAllPages(ctx, service.Client, dlpEDMSchemaEndpoint, &edmData)
	return edmData, err
}
