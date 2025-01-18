package dlp_incidents

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa/services/common"
)

const (
	baseIncidentEndpoint = "/dlp/v1/incidents"
)

// ##################### ADDS NOTES TO DLP INCIDENTS #####################
type CreateNoteRequest struct {
	Notes string `json:"notes"`
}

type ResolutionDetailsRequest struct {
	ResolutionLabel common.Label `json:"resolutionLabel"`
	ResolutionCode  string       `json:"resolutionCode"`
	Notes           string       `json:"notes"`
}

// ##################### ASSIGNS LABELS TO DLP INCIDENTS #####################
type LabelsRequest struct {
	Labels []common.Label `json:"labels"`
}

// ##################### INCIDENT GROUP SEARCH #####################
type IncidentGroupRequest struct {
	IncidentGroupIDs []int `json:"incidentGroupIds"`
}

type IncidentGroupsResponse struct {
	IncidentGroups []IncidentGroup `json:"incidentGroups"`
}

type IncidentGroup struct {
	ID                              int    `json:"id"`
	Name                            string `json:"name"`
	Description                     string `json:"description"`
	Status                          string `json:"status"`
	IncidentGroupType               string `json:"incidentGroupType"`
	IsDLPIncidentGroupAlreadyMapped bool   `json:"isDLPIncidentGroupAlreadyMapped"`
	IsDLPAdminConfigAlreadyMapped   bool   `json:"isDLPAdminConfigAlreadyMapped"`
}

// ##################### DLP INCIDENT CHANGE HISTORY #####################
type IncidentHistoryResponse struct {
	IncidentID    string          `json:"incidentId"`
	StartDate     string          `json:"startDate"`
	EndDate       string          `json:"endDate"`
	ChangeHistory []ChangeHistory `json:"changeHistory"`
}

type ChangeHistory struct {
	ChangeType    string     `json:"changeType"`
	ChangedAt     string     `json:"changedAt"`
	ChangedByName string     `json:"changedByName"`
	ChangeData    ChangeData `json:"changeData"`
	Comment       string     `json:"comment"`
}

type ChangeData struct {
	Before string `json:"before"`
	After  string `json:"after"`
}

// ##################### DLP INCIDENT GENERATED TICKETS #####################
type DLPIncidentTicketsResponse struct {
	Tickets []Ticket `json:"tickets"`
}

type Ticket struct {
	TicketType          string     `json:"ticketType"`
	TicketingSystemName string     `json:"ticketingSystemName"`
	ProjectID           string     `json:"projectId"`
	ProjectName         string     `json:"projectName"`
	TicketInfo          TicketInfo `json:"ticketInfo"`
}

type TicketInfo struct {
	TicketID  string `json:"ticketId"`
	TicketURL string `json:"ticketUrl"`
	State     string `json:"state"`
}

type DLPIncidentTriggerData map[string]string

type DLPIncidentEvidence struct {
	FileName       string `json:"fileName"`
	FileType       string `json:"fileType"`
	AdditionalInfo string `json:"additionalInfo"`
	EvidenceURL    string `json:"evidenceURL"`
}

func CreateNotes(ctx context.Context, service *services.Service, dlpIncidentID int, note string) (*common.IncidentDetails, *http.Response, error) {
	if note == "" {
		return nil, nil, errors.New("note is required")
	}

	path := fmt.Sprintf("%s/notes/%d", baseIncidentEndpoint, dlpIncidentID)

	requestPayload := CreateNoteRequest{Notes: note}

	var response common.IncidentDetails

	resp, err := service.Client.NewRequestDo(ctx, "POST", path, nil, requestPayload, &response)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create note for DLP incident %d: %w", dlpIncidentID, err)
	}
	return &response, resp, nil
}

func UpdateIncidentStatus(ctx context.Context, service *services.Service, dlpIncidentID string, close string) (*common.IncidentDetails, *http.Response, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}
	if close == "" {
		return nil, nil, errors.New("note is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s/close", baseIncidentEndpoint, dlpIncidentID)

	// Create the request payload
	requestPayload := ResolutionDetailsRequest{Notes: close}

	// Initialize a variable to hold the response
	var response common.IncidentDetails

	// Make the POST request
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, nil, requestPayload, &response)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to close DLP incident %s: %w", dlpIncidentID, err)
	}

	return &response, resp, nil
}

func AssignLabels(ctx context.Context, service *services.Service, dlpIncidentID string, labels []common.Label) (*common.IncidentDetails, *http.Response, error) {
	if len(labels) == 0 {
		return nil, nil, errors.New("labels are required")
	}

	path := fmt.Sprintf("%s/%s/labels", baseIncidentEndpoint, dlpIncidentID)

	requestPayload := LabelsRequest{
		Labels: labels,
	}

	var response common.IncidentDetails

	// Make the POST request
	resp, err := service.Client.NewRequestDo(ctx, "POST", path, nil, requestPayload, &response)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to assign labels to DLP incident %s: %w", dlpIncidentID, err)
	}

	return &response, resp, nil
}

func FilterIncidentSearch(ctx context.Context, service *services.Service, filters common.CommonDLPIncidentFiltering, paginationParams *common.PaginationParams) ([]common.IncidentDetails, *common.Cursor, error) {
	// Construct the endpoint URL
	path := fmt.Sprintf("%s/search", baseIncidentEndpoint)

	// Read all pages of audit logs using POST
	allResults, cursor, err := common.ReadAllPages[common.IncidentDetails](ctx, service.Client, http.MethodPost, path, paginationParams, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch DLP incident search: %w", err)
	}

	return allResults, cursor, nil
}

func AssignIncidentGroups(ctx context.Context, service *services.Service, dlpIncidentID int, groupIDs []int) (*IncidentGroupsResponse, *http.Response, error) {
	if len(groupIDs) == 0 {
		return nil, nil, errors.New("incident group IDs are required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%d/incident-groups/search", baseIncidentEndpoint, dlpIncidentID)

	// Prepare the request payload
	requestPayload := IncidentGroupRequest{
		IncidentGroupIDs: groupIDs,
	}

	var response IncidentGroupsResponse

	resp, err := service.Client.NewRequestDo(ctx, "POST", path, nil, requestPayload, &response)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to assign incident groups to DLP incident %d: %w", dlpIncidentID, err)
	}

	return &response, resp, nil
}

// ############################ DLP INCIDENT TRANSACTIONS ############################
// Gets the list of all DLP incidents associated with the transaction ID. A transaction ID can contain one or more DLP incidents.
func GetIncidentTransactions(ctx context.Context, service *services.Service, transactionID string, paginationParams *common.PaginationParams) ([]common.IncidentDetails, *common.Cursor, error) {
	if transactionID == "" {
		return nil, nil, errors.New("transaction ID is required")
	}

	// Construct the endpoint URL with transaction ID after /transactions
	endpoint := fmt.Sprintf("%s/transactions/%s", baseIncidentEndpoint, transactionID)

	// Use ReadAllPages to handle pagination
	allResults, cursor, err := common.ReadAllPages[common.IncidentDetails](ctx, service.Client, http.MethodGet, endpoint, paginationParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch transactions for transaction ID %s: %w", transactionID, err)
	}

	return allResults, cursor, nil
}

// Gets the DLP incident details based on the incident ID.
func GetDLPIncident(ctx context.Context, service *services.Service, dlpIncidentID string, fields []string) (*common.IncidentDetails, *http.Response, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s", baseIncidentEndpoint, dlpIncidentID)

	// Prepare the query parameters
	queryParams := url.Values{}
	if len(fields) > 0 {
		for _, field := range fields {
			queryParams.Add("fields", field)
		}
	}

	// Append query parameters to the URL
	fullURL := fmt.Sprintf("%s?%s", path, queryParams.Encode())

	// Initialize a variable to hold the response
	var incidentDetails common.IncidentDetails

	// Make the GET request
	resp, err := service.Client.NewRequestDo(ctx, "GET", fullURL, nil, nil, &incidentDetails)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve DLP incident %s: %w", dlpIncidentID, err)
	}

	return &incidentDetails, resp, nil
}

// Deletes the DLP incident for the specified incident ID.
func DeleteDLPIncident(ctx context.Context, service *services.Service, dlpIncidentID string) (*http.Response, error) {
	if dlpIncidentID == "" {
		return nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s", baseIncidentEndpoint, dlpIncidentID)

	// Make the DELETE request
	resp, err := service.Client.NewRequestDo(ctx, "DELETE", path, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to delete DLP incident %s: %w", dlpIncidentID, err)
	}

	return resp, nil
}

// Gets the details of updates made to an incident based on the given ID and timeline.
func HistoryDLPIncident(ctx context.Context, service *services.Service, dlpIncidentID string) (*IncidentHistoryResponse, *http.Response, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s/change-history", baseIncidentEndpoint, dlpIncidentID)

	// Initialize a variable to hold the response
	var response IncidentHistoryResponse

	// Make the GET request
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, &response)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieve DLP incident history %s: %w", dlpIncidentID, err)
	}

	return &response, resp, nil
}

// Gets the information of the ticket generated for the incident. For example, ticket type, ticket ID, ticket status, etc.
func GetDLPIncidentTickets(ctx context.Context, service *services.Service, dlpIncidentID string, paginationParams *common.PaginationParams) ([]Ticket, *common.Cursor, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/tickets/%s", baseIncidentEndpoint, dlpIncidentID)

	// Use ReadAllPages to fetch paginated results
	allResults, cursor, err := common.ReadAllPages[Ticket](ctx, service.Client, http.MethodGet, path, paginationParams, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch tickets for DLP incident %s: %w", dlpIncidentID, err)
	}

	return allResults, cursor, nil
}

func GetDLPIncidentTriggers(ctx context.Context, service *services.Service, dlpIncidentID string) (DLPIncidentTriggerData, *http.Response, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s/triggers", baseIncidentEndpoint, dlpIncidentID)

	// Initialize a variable to hold the response
	var triggers DLPIncidentTriggerData

	// Make the GET request
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, &triggers)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch triggers for DLP incident %s: %w", dlpIncidentID, err)
	}

	return triggers, resp, nil
}

func GetDLPIncidentEvidence(ctx context.Context, service *services.Service, dlpIncidentID string) (*DLPIncidentEvidence, *http.Response, error) {
	if dlpIncidentID == "" {
		return nil, nil, errors.New("valid DLP incident ID is required")
	}

	// Construct the endpoint URL
	path := fmt.Sprintf("%s/%s/evidence", baseIncidentEndpoint, dlpIncidentID)

	// Initialize a variable to hold the response
	var evidence DLPIncidentEvidence

	// Make the GET request
	resp, err := service.Client.NewRequestDo(ctx, "GET", path, nil, nil, &evidence)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch evidence for DLP incident %s: %w", dlpIncidentID, err)
	}

	return &evidence, resp, nil
}
