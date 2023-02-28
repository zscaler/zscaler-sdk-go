package eventlogentryreport

import (
	"errors"
)

const (
	eventlogEntryReportEndpoint = "/eventlogEntryReport"
)

type EventLogEntryReportTaskInfo struct {
	// Status of running task
	Status string `json:"status,omitempty"`

	// Number of items processed
	ProgressItemsComplete int `json:"progressItemsComplete,omitempty"`

	// End time
	ProgressEndTime int `json:"progressEndTime,omitempty"`

	// Error message
	ErrorMessage string `json:"errorMessage,omitempty"`
	ErrorCode    string `json:"errorCode,omitempty"`
}

type EventLogEntryReport struct {
	// The start time in the time range used to generate the event log report
	StartTime int `json:"startTime,omitempty"`

	// The end time in the time range used to generate the event log report
	EndTime  int    `json:"endTime,omitempty"`
	Page     int    `json:"page,omitempty"`
	PageSize string `json:"pageSize,omitempty"`

	// Filters the list based on the category for which the events were recorded.
	Category string `json:"category,omitempty"`

	// Filters the list based on areas within a category where the events were recorded
	Subcategories []string `json:"subcategories,omitempty"`

	// Filters the list based on the outcome (i.e., Failure or Success) of the events recorded
	ActionResult string `json:"actionResult,omitempty"`

	// The search string used to match against the event log message
	Message string `json:"message,omitempty"`

	// The search string used to match against the error code in event log entries
	ErrorCode string `json:"errorCode,omitempty"`

	// The search string used to match against the status code in event log entries
	StatusCode string `json:"statusCode,omitempty"`
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
