package adminauditlogs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	auditLogEntryReportEndpoint = "/zia/api/v1/auditlogEntryReport"
)

type AuditLogEntryReportTaskInfo struct {
	// Status of running task
	Status string `json:"status,omitempty"`

	// Number of items processed
	ProgressItemsComplete int `json:"progressItemsComplete,omitempty"`

	// End time
	ProgressEndTime int `json:"progressEndTime,omitempty"`

	// Error message
	ErrorMessage string `json:"errorMessage,omitempty"`

	ErrorCode string `json:"errorCode,omitempty"`
}

type AuditLogEntryRequest struct {

	// Action type for audit log entry
	TargetOrgId int `json:"targetOrgId,omitempty"`

	// Action type for audit log entry
	TraceId int `json:"traceId,omitempty"`

	// The start time in the time range used to generate the event log report
	StartTime int `json:"startTime,omitempty"`

	// The end time in the time range used to generate the event log report
	EndTime int `json:"endTime,omitempty"`

	Page int `json:"page,omitempty"`

	PageSize string `json:"pageSize,omitempty"`

	// Admin name for audit log entry
	AdminName string `json:"adminName,omitempty"`

	// Object name for audit log entry
	ObjectName string `json:"objectName,omitempty"`

	// Interface for audit log entry
	ActionInterface string `json:"actionInterface,omitempty"`

	// Filters the list based on the category for which the events were recorded.
	Category string `json:"category,omitempty"`

	// Filters the list based on areas within a category where the events were recorded
	Subcategories []string `json:"subcategories,omitempty"`

	// Filters the list based on the outcome (i.e., Failure or Success) of the events recorded
	ActionResult string `json:"actionResult,omitempty"`

	// Action type for audit log entry
	ActionTypes []string `json:"actionTypes,omitempty"`

	// Client IP for audit log entry
	ClientIP int `json:"clientIP,omitempty"`
}

func GetAll(ctx context.Context, service *zscaler.Service) (AuditLogEntryReportTaskInfo, error) {
	var auditLogEntryReport AuditLogEntryReportTaskInfo
	err := service.Client.Read(ctx, auditLogEntryReportEndpoint, &auditLogEntryReport)
	return auditLogEntryReport, err
}

func GetAdminAuditLogsDownload(ctx context.Context, service *zscaler.Service) ([]byte, error) {
	var csvData []byte

	// Perform a GET request to download the CSV file
	err := service.Client.Read(ctx, auditLogEntryReportEndpoint+"/download", &csvData)
	if err != nil {
		return nil, fmt.Errorf("failed to download audit log report: %w", err)
	}

	return csvData, nil
}

func CreateAdminAuditLogsExport(ctx context.Context, service *zscaler.Service, exportRequest AuditLogEntryRequest) (*http.Response, error) {
	// Call CreateWithNoContent and directly assign the *http.Response
	httpResp, err := service.Client.CreateWithNoContent(ctx, auditLogEntryReportEndpoint, exportRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to export audit log entry report: %w", err)
	}

	// Ensure the response is 204 No Content as expected
	if httpResp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("unexpected response code: %d, expected 204 No Content", httpResp.StatusCode)
	}

	// Log successful report creation
	service.Client.GetLogger().Printf("[DEBUG] Successfully triggered audit log entry report export with payload: %+v", exportRequest)

	return httpResp, nil
}

func Delete(ctx context.Context, service *zscaler.Service) (*http.Response, error) {
	err := service.Client.Delete(ctx, auditLogEntryReportEndpoint)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
