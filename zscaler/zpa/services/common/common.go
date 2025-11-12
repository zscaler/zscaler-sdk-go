package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
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

type CommonIDName struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ZPA Common Structs to Avoid Repetion
type CommonConfigDetails struct {
	Name   string `json:"name,omitempty"`
	Reason string `json:"reason,omitempty"`
}

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

type ExtranetDTO struct {
	LocationDTO      []LocationDTO      `json:"locationDTO,omitempty"`
	LocationGroupDTO []LocationGroupDTO `json:"locationGroupDTO,omitempty"`
	ZiaErName        string             `json:"ziaErName,omitempty"`
	ZpnErID          string             `json:"zpnErId,omitempty"`
}

type LocationDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type LocationGroupDTO struct {
	ID           string          `json:"id,omitempty"`
	Name         string          `json:"name,omitempty"`
	ZiaLocations []CommonSummary `json:"ziaLocations,omitempty"`
}

type ZPNERID struct {
	ID              string `json:"id,omitempty"`
	CreationTime    string `json:"creationTime,omitempty"`
	ModifiedBy      string `json:"modifiedBy,omitempty"`
	ModifiedTime    string `json:"modifiedTime,omitempty"`
	ZIACloud        string `json:"ziaCloud,omitempty"`
	ZIAErID         string `json:"ziaErId,omitempty"`
	ZIAErName       string `json:"ziaErName,omitempty"`
	ZIAModifiedTime string `json:"ziaModifiedTime,omitempty"`
	ZIAOrgID        string `json:"ziaOrgId,omitempty"`
}

type CommonSummary struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
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

// func getAllPagesGenericWithCustomFilters[T any](ctx context.Context, client *zscaler.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
// 	var paged struct {
// 		TotalPages interface{} `json:"totalPages"`
// 		List       []T         `json:"list"`
// 	}

// 	var rawList []T

// 	resp, err := client.NewRequestDo(ctx, "GET", relativeURL, Pagination{
// 		Search2:       filters.Search,
// 		MicroTenantID: filters.MicroTenantID,
// 		PageSize:      pageSize,
// 		Page:          page,
// 		SortBy:        filters.SortBy,
// 		SortOrder:     filters.SortOrder,
// 	}, nil, &paged)

// 	if err == nil && len(paged.List) > 0 {
// 		// ✅ Standard paginated response
// 		pages := fmt.Sprintf("%v", paged.TotalPages)
// 		totalPages, _ := strconv.Atoi(pages)
// 		return totalPages, paged.List, resp, nil
// 	}

// 	// 🔄 Retry as raw array (non-paginated)
// 	resp, err = client.NewRequestDo(ctx, "GET", relativeURL, Pagination{
// 		Search2:       filters.Search,
// 		MicroTenantID: filters.MicroTenantID,
// 		PageSize:      pageSize,
// 		Page:          page,
// 		SortBy:        filters.SortBy,
// 		SortOrder:     filters.SortOrder,
// 	}, nil, &rawList)
// 	if err != nil {
// 		return 0, nil, resp, err
// 	}

// 	// ✅ API returned a raw array: treat as single page
// 	return 1, rawList, resp, nil
// }

func getAllPagesGenericWithCustomFilters[T any](ctx context.Context, client *zscaler.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
	var paged struct {
		TotalPages interface{} `json:"totalPages"`
		List       []T         `json:"list"`
	}

	// Build pagination struct with search filter (already converted to filter format)
	// The encoding (%20 vs +) is handled in zparequests.go based on endpoint type
	pagination := Pagination{
		MicroTenantID: filters.MicroTenantID,
		PageSize:      pageSize,
		Page:          page,
		SortBy:        filters.SortBy,
		SortOrder:     filters.SortOrder,
	}
	if filters.Search != "" {
		pagination.Search2 = filters.Search
	}

	resp, err := client.NewRequestDo(ctx, "GET", relativeURL, pagination, nil, &paged)

	if err != nil {
		return 0, nil, resp, err
	}

	pages := fmt.Sprintf("%v", paged.TotalPages)
	totalPages, _ := strconv.Atoi(pages)

	// Even if totalPages == 0, return list to prevent fallback to raw array
	return totalPages, paged.List, resp, nil
}

func getAllPagesGeneric[T any](ctx context.Context, client *zscaler.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
	return getAllPagesGenericWithCustomFilters[T](
		ctx,
		client,
		relativeURL,
		page,
		pageSize,
		filters,
	)
}

// GetAllPagesGeneric fetches all resources instead of just one single page
func GetAllPagesGeneric[T any](ctx context.Context, client *zscaler.Client, relativeURL, searchQuery string) ([]T, *http.Response, error) {
	// Convert search query to filter format for ZPA endpoints (except SCIM endpoints)
	// Don't pre-encode here - let the query parameter encoding handle it
	if isZPAEndpoint(relativeURL) && !isSCIMEndpoint(relativeURL) {
		searchQuery = convertZPASearchToFilter(searchQuery)
	} else {
		// For non-ZPA endpoints or SCIM endpoints, sanitize the query
		searchQuery = sanitizeSearchQuery(searchQuery)
	}

	totalPages, result, resp, err := getAllPagesGeneric[T](ctx, client, relativeURL, 1, DefaultPageSize, Filter{Search: searchQuery})
	if err != nil {
		return nil, resp, err
	}
	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGeneric[T](ctx, client, relativeURL, page, DefaultPageSize, Filter{Search: searchQuery})
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

func getMicroTenantByName(ctx context.Context, client *zscaler.Client, microTenantName string) (*microTenantSample, *http.Response, error) {
	relativeURL := "/zpa/mgmtconfig/v1/admin/customers/" + client.GetCustomerID() + "/microtenants"
	list, resp, err := GetAllPagesGeneric[microTenantSample](ctx, client, relativeURL, microTenantName)
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
func GetAllPagesGenericWithCustomFilters[T any](ctx context.Context, client *zscaler.Client, relativeURL string, filters Filter) ([]T, *http.Response, error) {
	if (filters.MicroTenantID == nil || *filters.MicroTenantID == "") && filters.MicroTenantName != nil && *filters.MicroTenantName != "" {
		// get microtenant id by name
		mt, resp, err := getMicroTenantByName(ctx, client, *filters.MicroTenantName)
		if err == nil {
			return nil, resp, err
		}
		if mt != nil {
			filters.MicroTenantID = &mt.ID
		}
	}

	// Ensure Search term is sanitized correctly (prevent double encoding)
	// For ZPA, convert simple search strings to filter format (name+EQ+<value>)
	// Exception: SCIM endpoints use plain search strings
	if filters.Search != "" {
		// Check if this is a ZPA endpoint (but not SCIM)
		if isZPAEndpoint(relativeURL) && !isSCIMEndpoint(relativeURL) {
			filters.Search = convertZPASearchToFilter(filters.Search)
		} else {
			filters.Search = sanitizeSearchQuery(filters.Search)
		}
	}

	// Attempt full search first.
	totalPages, result, resp, err := getAllPagesGenericWithCustomFilters[T](ctx, client, relativeURL, 1, DefaultPageSize, filters)
	// If the full search fails and the query contains multiple words, try a partial search.
	if err != nil && strings.Count(filters.Search, " ") > 0 {
		// Extract the value part if search is in filter format (e.g., name+EQ+value or name%2BEQ%2Bvalue)
		searchValue := filters.Search
		if isZPAEndpoint(relativeURL) {
			// Check for unencoded filter format (name+EQ+value)
			if strings.Contains(filters.Search, "+EQ+") {
				parts := strings.SplitN(filters.Search, "+EQ+", 2)
				if len(parts) == 2 {
					searchValue = parts[1]
				}
			} else if strings.Contains(filters.Search, "%2BEQ%2B") {
				// Check for URL-encoded filter format (name%2BEQ%2Bvalue)
				parts := strings.SplitN(filters.Search, "%2BEQ%2B", 2)
				if len(parts) == 2 {
					// URL decode the value part
					if decoded, err := url.QueryUnescape(parts[1]); err == nil {
						searchValue = decoded
					} else {
						searchValue = parts[1]
					}
				}
			}
		}

		// Only try partial search if we have multiple words in the value
		if strings.Count(searchValue, " ") > 0 {
			tokens := strings.Split(searchValue, " ")
			if len(tokens) >= 2 {
				// Use only the first two words for partial search
				partialValue := strings.Join(tokens[:2], " ")
				// Reconstruct filter format if needed (but not for SCIM endpoints)
				if isZPAEndpoint(relativeURL) && !isSCIMEndpoint(relativeURL) && (strings.Contains(filters.Search, "+EQ+") || strings.Contains(filters.Search, "%2BEQ%2B")) {
					filters.Search = convertZPASearchToFilter(partialValue)
				} else {
					filters.Search = partialValue
				}
				totalPages, result, resp, err = getAllPagesGenericWithCustomFilters[T](ctx, client, relativeURL, 1, DefaultPageSize, filters)
			}
		}
	}
	if err != nil {
		return nil, resp, err
	}

	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGenericWithCustomFilters[T](ctx, client, relativeURL, page, DefaultPageSize, filters)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}

	return result, resp, nil
}

func GetAllPagesScimGenericWithSearch[T any](
	ctx context.Context,
	client *zpa.ScimZpaClient,
	baseURL string,
	itemsPerPage int,
	searchFunc func(T) bool,
) ([]T, *http.Response, error) {
	// Enforce default and maximum limits for itemsPerPage
	if itemsPerPage <= 0 {
		itemsPerPage = 10 // Default to 10 if not specified
	} else if itemsPerPage > 100 {
		itemsPerPage = 100 // Enforce maximum limit of 100
	}

	var allResources []T
	startIndex := 1
	var lastResp *http.Response

	for {
		// Construct the paginated URL
		paginatedURL := fmt.Sprintf("%s?startIndex=%d&count=%d", baseURL, startIndex, itemsPerPage)

		// Define the structure for the paginated response
		var paginatedResponse struct {
			Resources    []T `json:"Resources"`
			TotalResults int `json:"totalResults"`
		}

		// Perform the HTTP request and parse the response
		resp, err := client.DoRequest(ctx, http.MethodGet, paginatedURL, nil, &paginatedResponse)
		if err != nil {
			return nil, resp, fmt.Errorf("error fetching paginated data: %w", err)
		}
		lastResp = resp // Track the last HTTP response

		// Iterate through the resources to search for the specific item
		for _, resource := range paginatedResponse.Resources {
			if searchFunc != nil && searchFunc(resource) {
				// Return immediately if the desired item is found
				return []T{resource}, resp, nil
			}
		}

		// Append resources to the result set if not searching
		if searchFunc == nil {
			allResources = append(allResources, paginatedResponse.Resources...)
		}

		// Check if all records have been fetched
		if startIndex+itemsPerPage > paginatedResponse.TotalResults || len(paginatedResponse.Resources) == 0 {
			break
		}

		// Move to the next page
		startIndex += itemsPerPage
	}

	// Return all resources if no specific item was found
	return allResources, lastResp, nil
}

// isZPAEndpoint checks if the given URL is a ZPA endpoint
func isZPAEndpoint(url string) bool {
	return strings.Contains(url, "/zpa/") || strings.Contains(url, "/mgmtconfig/")
}

// isSCIMEndpoint checks if the given URL is a SCIM endpoint
// SCIM endpoints use plain search strings, not the name+EQ+value filter format
func isSCIMEndpoint(url string) bool {
	return strings.Contains(url, "/scimgroup") || strings.Contains(url, "/scimattributeheader")
}

// convertZPASearchToFilter converts a simple search string to ZPA filter format
// If the search string already contains filter operators, it returns it as-is
// Otherwise, it converts to name+EQ+<value> format for exact name matching
// Note: SCIM endpoints should not use this function as they expect plain search strings
func convertZPASearchToFilter(search string) string {
	// Trim whitespace
	search = strings.TrimSpace(search)
	if search == "" {
		return search
	}

	// Remove any existing quotes
	if strings.HasPrefix(search, `"`) && strings.HasSuffix(search, `"`) {
		search = search[1 : len(search)-1]
		search = strings.TrimSpace(search)
	}

	// Check if the search string already contains filter operators
	// Common ZPA filter operators: EQ, NE, GT, LT, GE, LE, CONTAINS, STARTSWITH, ENDSWITH
	// These can appear as +EQ+, +NE+, etc. or URL-encoded as %2BEQ%2B, etc.
	filterOperators := []string{
		"+EQ+", "+NE+", "+GT+", "+LT+", "+GE+", "+LE+",
		"+CONTAINS+", "+STARTSWITH+", "+ENDSWITH+",
		"%2BEQ%2B", "%2BNE%2B", "%2BGT%2B", "%2BLT%2B", "%2BGE%2B", "%2BLE%2B",
		"%2BCONTAINS%2B", "%2BSTARTSWITH%2B", "%2BENDSWITH%2B",
	}

	// Check if search already contains any filter operator pattern
	searchUpper := strings.ToUpper(search)
	for _, op := range filterOperators {
		if strings.Contains(searchUpper, strings.ToUpper(op)) {
			// Already in filter format - return as-is (don't sanitize, just trim)
			return strings.TrimSpace(search)
		}
	}

	// Simple search string - convert to name+EQ+<value> format for exact name matching
	// For ZPA filter format, preserve multiple consecutive spaces (don't collapse them)
	// Only trim leading/trailing spaces, but preserve internal spacing
	sanitizedValue := strings.TrimSpace(search)
	// Remove special characters except spaces, alphanumeric characters, dashes, underscores, slashes, and dots
	// But preserve multiple consecutive spaces
	re := regexp.MustCompile(`[^a-zA-Z0-9\s_/\-\.]`)
	sanitizedValue = re.ReplaceAllString(sanitizedValue, "")
	// Convert to filter format: name+EQ+<value>
	// Note: The + symbols in +EQ+ will be URL-encoded as %2B when the query is built
	// Spaces in the value will be encoded as + by url.Values.Encode(), which works for the API
	filterFormat := fmt.Sprintf("name+EQ+%s", sanitizedValue)
	return filterFormat
}

func sanitizeSearchQuery(query string) string {
	// Just trim and return - let URL encoding handle special characters
	return strings.TrimSpace(query)
}
