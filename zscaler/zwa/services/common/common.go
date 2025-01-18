package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zwa"
)

const pageSize = 1000

type IncidentDetails struct {
	InternalID            string              `json:"internalId"`
	IntegrationType       string              `json:"integrationType"`
	TransactionID         string              `json:"transactionId"`
	SourceType            string              `json:"sourceType"`
	SourceSubType         string              `json:"sourceSubType"`
	SourceActions         []string            `json:"sourceActions"`
	Severity              string              `json:"severity"`
	Priority              string              `json:"priority"`
	MatchingPolicies      MatchingPolicies    `json:"matchingPolicies"`
	MatchCount            int                 `json:"matchCount"`
	CreatedAt             string              `json:"createdAt"`
	LastUpdatedAt         string              `json:"lastUpdatedAt"`
	SourceFirstObservedAt string              `json:"sourceFirstObservedAt"`
	SourceLastObservedAt  string              `json:"sourceLastObservedAt"`
	UserInfo              UserInfo            `json:"userInfo"`
	ApplicationInfo       ApplicationInfo     `json:"applicationInfo"`
	ContentInfo           ContentInfo         `json:"contentInfo"`
	NetworkInfo           NetworkInfo         `json:"networkInfo"`
	MetadataFileURL       string              `json:"metadataFileUrl"`
	Status                string              `json:"status"`
	Resolution            string              `json:"resolution"`
	AssignedAdmin         AssignedAdmin       `json:"assignedAdmin"`
	LastNotifiedUser      LastNotifiedUser    `json:"lastNotifiedUser"`
	Notes                 []Note              `json:"notes"`
	ClosedCode            string              `json:"closedCode"`
	IncidentGroupIDs      []int               `json:"incidentGroupIds"`
	IncidentGroups        []IncidentGroup     `json:"incidentGroups"`
	DLPIncidentTickets    []DLPIncidentTicket `json:"dlpIncidentTickets"`
	Labels                []Label             `json:"labels"`
}

type MatchingPolicies struct {
	Engines      []Engine     `json:"engines"`
	Rules        []Rule       `json:"rules"`
	Dictionaries []Dictionary `json:"dictionaries"`
}

type Engine struct {
	Name string `json:"name"`
	Rule string `json:"rule"`
}

type Rule struct {
	Name string `json:"name"`
}

type Dictionary struct {
	Name           string `json:"name"`
	MatchCount     int    `json:"matchCount"`
	NameMatchCount string `json:"nameMatchCount"`
}

type UserInfo struct {
	Name             string      `json:"name"`
	Email            string      `json:"email"`
	ClientIP         string      `json:"clientIP"`
	UniqueIdentifier string      `json:"uniqueIdentifier"`
	UserID           int         `json:"userId"`
	Department       string      `json:"department"`
	HomeCountry      string      `json:"homeCountry"`
	ManagerInfo      ManagerInfo `json:"managerInfo"`
}

type ManagerInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ApplicationInfo struct {
	URL                   string `json:"url"`
	Category              string `json:"category"`
	Name                  string `json:"name"`
	HostnameOrApplication string `json:"hostnameOrApplication"`
	AdditionalInfo        string `json:"additionalInfo"`
}

type ContentInfo struct {
	FileName       string `json:"fileName"`
	FileType       string `json:"fileType"`
	AdditionalInfo string `json:"additionalInfo"`
}

type NetworkInfo struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type AssignedAdmin struct {
	Email string `json:"email"`
}

type LastNotifiedUser struct {
	Role  string `json:"role"`
	Email string `json:"email"`
}

type IncidentGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DLPIncidentTicket struct {
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

// Common Struct
type Note struct {
	Body          string `json:"body"`
	CreatedAt     string `json:"createdAt"`
	LastUpdatedAt string `json:"lastUpdatedAt"`
	CreatedBy     int    `json:"createdBy"`
	LastUpdatedBy int    `json:"lastUpdatedBy"`
}

// Common Struct
type Label struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ##################### FILTERS DLP INCIDENTS BASED ON RANGE #####################
type CommonDLPIncidentFiltering struct {
	Fields    []Fields  `json:"fields"`
	TimeRange TimeRange `json:"timeRange"`
}

type Fields struct {
	Name  string   `json:"name"`
	Value []string `json:"value"`
}

type TimeRange struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// Cursor represents the pagination details returned in responses.
type Cursor struct {
	TotalPages        int    `json:"totalPages"`
	CurrentPageNumber int    `json:"currentPageNumber"`
	CurrentPageSize   int    `json:"currentPageSize"`
	PageID            string `json:"pageId"`
	TotalElements     int    `json:"totalElements"`
}

// PaginationParams represents the pagination parameters used in requests.
type PaginationParams struct {
	Page     *int    `url:"page,omitempty"`
	PageSize *int    `url:"pageSize,omitempty"`
	PageID   *string `url:"pageId,omitempty"`
}

// IntPtr converts an int value to a pointer to an int.
func IntPtr(i int) *int {
	return &i
}

// GetPageSize returns the page size.
func GetPageSize() int {
	return pageSize
}

func ReadAllPages[T any](ctx context.Context, client *zwa.Client, method, endpoint string, params *PaginationParams, requestBody interface{}) ([]T, *Cursor, error) {
	var allResults []T
	page := 1
	pageSize := 1000 // Default page size
	var cursor Cursor

	// Override default params if provided
	if params != nil {
		if params.Page != nil {
			page = *params.Page
		}
		if params.PageSize != nil {
			pageSize = *params.PageSize
		}
	}

	for {
		// Add pagination parameters dynamically
		queryParams := url.Values{}
		queryParams.Set("page", fmt.Sprintf("%d", page))
		queryParams.Set("pageSize", fmt.Sprintf("%d", pageSize))
		if params != nil && params.PageID != nil {
			queryParams.Set("pageId", *params.PageID)
		}

		// Parse the endpoint into a URL and append query parameters
		baseURL, err := url.Parse(endpoint)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid endpoint URL: %w", err)
		}
		baseURL.RawQuery = queryParams.Encode()

		// Prepare the response container
		var pageResults struct {
			Items  []T    `json:"logs"`
			Cursor Cursor `json:"cursor"`
		}

		// Execute the request
		var resp *http.Response
		if method == http.MethodGet {
			resp, err = client.NewRequestDo(ctx, method, baseURL.String(), nil, nil, &pageResults)
		} else if method == http.MethodPost {
			resp, err = client.NewRequestDo(ctx, method, baseURL.String(), nil, requestBody, &pageResults)
		} else {
			return nil, nil, fmt.Errorf("unsupported HTTP method: %s", method)
		}

		if err != nil {
			return nil, nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}
		defer resp.Body.Close()

		// Append items and update the cursor
		allResults = append(allResults, pageResults.Items...)
		cursor = pageResults.Cursor

		// Break if no more pages
		if cursor.CurrentPageSize < pageSize || page >= cursor.TotalPages-1 {
			break
		}
		page++
	}

	return allResults, &cursor, nil
}

func ReadPage[T any](ctx context.Context, client *zwa.Client, endpoint string, params PaginationParams) ([]T, *Cursor, error) {
	// Add pagination parameters dynamically
	queryParams := url.Values{}
	if params.Page != nil {
		queryParams.Set("page", fmt.Sprintf("%d", *params.Page))
	}
	if params.PageSize != nil {
		queryParams.Set("pageSize", fmt.Sprintf("%d", *params.PageSize))
	}
	if params.PageID != nil {
		queryParams.Set("pageId", *params.PageID)
	}

	// Build the full URL
	fullURL := fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())

	// Read a single page of results
	var pageResults struct {
		Items  []T    `json:"items"`
		Cursor Cursor `json:"cursor"`
	}
	resp, err := client.NewRequestDo(ctx, "GET", fullURL, nil, nil, &pageResults)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	return pageResults.Items, &pageResults.Cursor, nil
}
