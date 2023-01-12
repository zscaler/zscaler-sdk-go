package eventlogentryreport

import (
	"errors"
)

const (
	eventlogEntryReportEndpoint = "/eventlogEntryReport"
)

type EventLogEntryReportTaskInfo struct {
	Status                string `json:"status,omitempty"`
	ProgressItemsComplete int    `json:"progressItemsComplete,omitempty"`
	ProgressEndTime       int    `json:"progressEndTime,omitempty"`
	ErrorMessage          string `json:"errorMessage,omitempty"`
	ErrorCode             string `json:"errorCode,omitempty"`
}

type EventLogEntryReport struct {
	StartTime     int      `json:"startTime,omitempty"`
	EndTime       int      `json:"endTime,omitempty"`
	Page          int      `json:"page,omitempty"`
	PageSize      string   `json:"pageSize,omitempty"`
	Category      string   `json:"category,omitempty"`
	Subcategories []string `json:"subcategories,omitempty"`
	ActionResult  string   `json:"actionResult,omitempty"`
	Message       string   `json:"message,omitempty"`
	ErrorCode     string   `json:"errorCode,omitempty"`
	StatusCode    string   `json:"statusCode,omitempty"`
}

func (service *Service) GetAll() ([]EventLogEntryReportTaskInfo, error) {
	var eventLogEntryReport []EventLogEntryReportTaskInfo
	err := service.Client.Read(eventlogEntryReportEndpoint, &eventLogEntryReport)
	return eventLogEntryReport, err
}

func (service *Service) Create(eventLog *EventLogEntryReport) (*EventLogEntryReport, error) {
	resp, err := service.Client.Create(eventlogEntryReportEndpoint, eventLog)
	if err != nil {
		return nil, err
	}

	createdEventLogReport, ok := resp.(*EventLogEntryReport)
	if !ok {
		return nil, errors.New("object returned from api was not an event log entry report pointer")
	}

	service.Client.Logger.Printf("[DEBUG]returning event log entry report from create: %d", createdEventLogReport)
	return createdEventLogReport, nil
}
