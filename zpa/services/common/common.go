package common

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v2/zpa"
)

const (
	DefaultPageSize = 500
)

type Pagination struct {
	PageSize        int     `json:"pagesize,omitempty" url:"pagesize,omitempty"`
	Page            int     `json:"page,omitempty" url:"page,omitempty"`
	Search          string  `json:"-" url:"-"`
	Search2         string  `json:"search,omitempty" url:"search,omitempty"`
	MicroTenantID   *string `url:"microtenantId,omitempty"`
	MicroTenantName *string `url:"-,omitempty"`
	SortBy          string  `json:"sortBy,omitempty" url:"sortBy,omitempty"`       // New field for sorting by attribute
	SortOrder       string  `json:"sortOrder,omitempty" url:"sortOrder,omitempty"` // New field for the sort order (ASC or DESC)
}

type Filter struct {
	Search          string  `url:"search,omitempty"`
	MicroTenantID   *string `url:"microtenantId,omitempty"`
	MicroTenantName *string `url:"-,omitempty"`
	SortBy          string  `url:"sortBy,omitempty"`          // New field for sorting by attribute
	SortOrder       string  `url:"sortOrder,omitempty"`       // New field for the sort order (ASC or DESC)
	ApplicationType string  `url:"applicationType,omitempty"` // New field for filtering by application type
	ExpandAll       bool    `url:"expandAll,omitempty"`       // New field for deciding whether to expand all attributes
}

type DeleteApplicationQueryParams struct {
	ForceDelete   bool    `json:"forceDelete,omitempty" url:"forceDelete,omitempty"`
	MicroTenantID *string `url:"microtenantId,omitempty"`
}
type NetworkPorts struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}

// ZPA Inspection Rules
type Rules struct {
	Conditions []Conditions `json:"conditions,omitempty"`
	Names      string       `json:"names,omitempty"`
	Type       string       `json:"type,omitempty"`
	Version    string       `json:"version,omitempty"`
}

type Conditions struct {
	LHS string `json:"lhs,omitempty"`
	OP  string `json:"op,omitempty"`
	RHS string `json:"rhs,omitempty"`
}

type CustomCommonControls struct {
	ID                               string                   `json:"id,omitempty"`
	Name                             string                   `json:"name,omitempty"`
	Action                           string                   `json:"action,omitempty"`
	ActionValue                      string                   `json:"actionValue,omitempty"`
	AssociatedInspectionProfileNames []AssociatedProfileNames `json:"associatedInspectionProfileNames,omitempty"`
	Attachment                       string                   `json:"attachment,omitempty"`
	ControlGroup                     string                   `json:"controlGroup,omitempty"`
	ControlNumber                    string                   `json:"controlNumber,omitempty"`
	ControlType                      string                   `json:"controlType,omitempty"`
	CreationTime                     string                   `json:"creationTime,omitempty"`
	DefaultAction                    string                   `json:"defaultAction,omitempty"`
	DefaultActionValue               string                   `json:"defaultActionValue,omitempty"`
	Description                      string                   `json:"description,omitempty"`
	ModifiedBy                       string                   `json:"modifiedBy,omitempty"`
	ModifiedTime                     string                   `json:"modifiedTime,omitempty"`
	ParanoiaLevel                    string                   `json:"paranoiaLevel,omitempty"`
	ProtocolType                     string                   `json:"protocolType,omitempty"`
	Severity                         string                   `json:"severity,omitempty"`
	Version                          string                   `json:"version,omitempty"`
}

type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ZPA Common Structs to Avoid Repetion
type CommonConfigDetails struct {
	Name   string `json:"name,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// Shared Between App Connector and Service Edge Controllers and Groups
type ZPNSubModuleUpgrade struct {
	ID              string `json:"id,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	CurrentVersion  string `json:"currentVersion,omitempty"`
	EntityGid       string `json:"entityGid,omitempty"`
	EntityType      string `json:"entityType,omitempty"`
	ExpectedVersion string `json:"expectedVersion,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	PreviousVersion string `json:"previousVersion,omitempty"`
	Role            string `json:"role,omitempty"`
	UpgradeStatus   string `json:"upgradeStatus,omitempty"`
	UpgradeTime     string `json:"upgradeTime,omitempty"`
}

type Meta struct {
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Location     string    `json:"location"`
	ResourceType string    `json:"resourceType"`
}

// RemoveCloudSuffix removes appended cloud name (zscalerthree.net) i.e "CrowdStrike_ZPA_Pre-ZTA (zscalerthree.net)"
func RemoveCloudSuffix(str string) string {
	reg := regexp.MustCompile(`(.*)[\s]+\([a-zA-Z0-9\-_\.]*\)[\s]*$`)
	res := reg.ReplaceAllString(str, "${1}")
	return strings.Trim(res, " ")
}

func InList(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func getAllPagesGenericWithCustomFilters[T any](client *zpa.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
	var v struct {
		TotalPages interface{} `json:"totalPages"`
		List       []T         `json:"list"`
	}
	resp, err := client.NewRequestDo("GET", relativeURL, Pagination{
		Search2:       filters.Search,
		MicroTenantID: filters.MicroTenantID,
		PageSize:      pageSize,
		Page:          page,
		SortBy:        filters.SortBy,
		SortOrder:     filters.SortOrder,
	}, nil, &v)
	if err != nil {
		return 0, nil, resp, err
	}

	pages := fmt.Sprintf("%v", v.TotalPages)
	totalPages, _ := strconv.Atoi(pages)

	return totalPages, v.List, resp, nil
}

func getAllPagesGeneric[T any](client *zpa.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
	return getAllPagesGenericWithCustomFilters[T](
		client,
		relativeURL,
		page,
		pageSize,
		filters,
	)
}

// GetAllPagesGeneric fetches all resources instead of just one single page
// func GetAllPagesGeneric[T any](client *zpa.Client, relativeURL, searchQuery string) ([]T, *http.Response, error) {
// 	searchQuery = url.QueryEscape(searchQuery)
// 	totalPages, result, resp, err := getAllPagesGeneric[T](client, relativeURL, 1, DefaultPageSize, Filter{Search: searchQuery})
// 	if err != nil {
// 		return nil, resp, err
// 	}
// 	var l []T
// 	for page := 2; page <= totalPages; page++ {
// 		totalPages, l, resp, err = getAllPagesGeneric[T](client, relativeURL, page, DefaultPageSize, Filter{Search: searchQuery})
// 		if err != nil {
// 			return nil, resp, err
// 		}
// 		result = append(result, l...)
// 	}

// 	return result, resp, nil
// }

func GetAllPagesGeneric[T any](client *zpa.Client, relativeURL, searchQuery string) ([]T, *http.Response, error) {
	// Remove manual url.QueryEscape(searchQuery)
	// Optionally, you can still do searchQuery = strings.TrimSpace(searchQuery) if needed

	totalPages, result, resp, err := getAllPagesGeneric[T](
		client,
		relativeURL,
		1,
		DefaultPageSize,
		Filter{Search: searchQuery},
	)
	if err != nil {
		return nil, resp, err
	}

	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGeneric[T](
			client,
			relativeURL,
			page,
			DefaultPageSize,
			Filter{Search: searchQuery},
		)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}
	return result, resp, nil
}

type microTenantSample struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func getMicroTenantByName(client *zpa.Client, microTenantName string) (*microTenantSample, *http.Response, error) {
	relativeURL := "/mgmtconfig/v1/admin/customers/" + client.Config.CustomerID + "/microtenants"
	list, resp, err := GetAllPagesGeneric[microTenantSample](client, relativeURL, microTenantName)
	if err != nil {
		return nil, resp, err
	}
	for _, app := range list {
		if strings.EqualFold(app.Name, microTenantName) {
			return &app, resp, nil
		}
	}
	return nil, resp, fmt.Errorf("no microtenant named '%s' was found", microTenantName)
}

// GetAllPagesGenericWithCustomFilters fetches all resources instead of just one single page
func GetAllPagesGenericWithCustomFilters[T any](client *zpa.Client, relativeURL string, filters Filter) ([]T, *http.Response, error) {
	if (filters.MicroTenantID == nil || *filters.MicroTenantID == "") && filters.MicroTenantName != nil && *filters.MicroTenantName != "" {
		// get microtenant id by name
		mt, resp, err := getMicroTenantByName(client, *filters.MicroTenantName)
		if err == nil {
			return nil, resp, err
		}
		if mt != nil {
			filters.MicroTenantID = &mt.ID
		}
	}

	// Updated filter search: replace spaces with '&' for the API's search query format
	filters.Search = strings.ReplaceAll(filters.Search, " ", "&")

	totalPages, result, resp, err := getAllPagesGenericWithCustomFilters[T](client, relativeURL, 1, DefaultPageSize, filters)
	if err != nil {
		return nil, resp, err
	}
	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGenericWithCustomFilters[T](client, relativeURL, page, DefaultPageSize, filters)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}

	return result, resp, nil
}

func GetAllPagesScimGeneric[T any](ctx context.Context, client *zpa.ScimClient, baseURL string, itemsPerPage int) ([]T, *http.Response, error) {
	var allResources []T
	startIndex := 1
	var lastResp *http.Response

	for {
		// Construct URL with optional pagination parameters
		paginatedURL := baseURL
		if itemsPerPage > 0 {
			paginatedURL = fmt.Sprintf("%s?startIndex=%d&count=%d", baseURL, startIndex, itemsPerPage)
		}

		// Define the structure for paginated response
		var paginatedResponse struct {
			Resources    []T `json:"Resources"`
			TotalResults int `json:"totalResults"`
		}

		// Perform the request and capture the response
		resp, err := client.DoRequest(ctx, http.MethodGet, paginatedURL, nil, &paginatedResponse)
		if err != nil {
			return nil, resp, fmt.Errorf("error fetching paginated data: %w", err)
		}
		lastResp = resp // Track last response for return

		// Append resources and break if all results are retrieved
		allResources = append(allResources, paginatedResponse.Resources...)
		if len(allResources) >= paginatedResponse.TotalResults || len(paginatedResponse.Resources) == 0 {
			break
		}

		// Update startIndex for next page
		startIndex += itemsPerPage
	}

	return allResources, lastResp, nil
}
