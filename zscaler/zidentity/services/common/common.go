package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

// Common structures used across zidentity services
type IDNameDisplayName struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// PaginationResponse represents the standard pagination response structure
// used across zidentity API endpoints
type PaginationResponse[T any] struct {
	ResultsTotal int    `json:"results_total,omitempty"`
	PageOffset   int    `json:"pageOffset,omitempty"`
	PageSize     int    `json:"pageSize,omitempty"`
	NextLink     string `json:"next_link,omitempty"`
	PrevLink     string `json:"prev_link,omitempty"`
	Records      []T    `json:"records,omitempty"`
}

// PaginationQueryParams represents common query parameters for pagination
type PaginationQueryParams struct {
	Offset               int    `url:"offset,omitempty"`
	Limit                int    `url:"limit,omitempty"`
	NameLike             string `url:"name[like],omitempty"`
	ExcludeDynamicGroups bool   `url:"excludedynamicgroups,omitempty"`
	// User-specific parameters for /groups/{id}/users endpoint
	LoginName        []string `url:"loginname,omitempty"`
	LoginNameLike    string   `url:"loginname[like],omitempty"`
	DisplayNameLike  string   `url:"displayname[like],omitempty"`
	PrimaryEmailLike string   `url:"primaryemail[like],omitempty"`
	DomainName       []string `url:"domainname,omitempty"`
	IDPName          []string `url:"idpname,omitempty"`
}

// PaginationOptions provides configuration for pagination behavior
type PaginationOptions struct {
	DefaultPageSize int
	MaxPageSize     int
	UseCursor       bool // If true, use cursor-based pagination (next_link/prev_link)
}

// Default pagination options for zidentity service
var DefaultPaginationOptions = PaginationOptions{
	DefaultPageSize: 100,
	MaxPageSize:     1000,
	UseCursor:       false,
}

// NewPaginationQueryParams creates a new PaginationQueryParams with sensible defaults
func NewPaginationQueryParams(pageSize int) PaginationQueryParams {
	if pageSize <= 0 {
		pageSize = DefaultPaginationOptions.DefaultPageSize
	}
	if pageSize > DefaultPaginationOptions.MaxPageSize {
		pageSize = DefaultPaginationOptions.MaxPageSize
	}

	return PaginationQueryParams{
		Limit: pageSize,
	}
}

// WithNameFilter adds a name filter to the query parameters
func (p *PaginationQueryParams) WithNameFilter(name string) *PaginationQueryParams {
	p.NameLike = name
	return p
}

// WithExcludeDynamicGroups excludes dynamic groups from results
func (p *PaginationQueryParams) WithExcludeDynamicGroups(exclude bool) *PaginationQueryParams {
	p.ExcludeDynamicGroups = exclude
	return p
}

// WithOffset sets the offset for pagination
func (p *PaginationQueryParams) WithOffset(offset int) *PaginationQueryParams {
	if offset >= 0 {
		p.Offset = offset
	}
	return p
}

// WithLimit sets the limit for pagination
func (p *PaginationQueryParams) WithLimit(limit int) *PaginationQueryParams {
	if limit > 0 && limit <= DefaultPaginationOptions.MaxPageSize {
		p.Limit = limit
	}
	return p
}

// ToURLValues converts PaginationQueryParams to url.Values
func (p *PaginationQueryParams) ToURLValues() url.Values {
	values := url.Values{}

	if p.Offset > 0 {
		values.Set("offset", strconv.Itoa(p.Offset))
	}
	if p.Limit > 0 {
		values.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.NameLike != "" {
		values.Set("name[like]", p.NameLike)
	}
	if p.ExcludeDynamicGroups {
		values.Set("excludedynamicgroups", "true")
	}
	// User-specific parameters
	if len(p.LoginName) > 0 {
		for _, loginName := range p.LoginName {
			values.Add("loginname", loginName)
		}
	}
	if p.LoginNameLike != "" {
		values.Set("loginname[like]", p.LoginNameLike)
	}
	if p.DisplayNameLike != "" {
		values.Set("displayname[like]", p.DisplayNameLike)
	}
	if p.PrimaryEmailLike != "" {
		values.Set("primaryemail[like]", p.PrimaryEmailLike)
	}
	if len(p.DomainName) > 0 {
		for _, domainName := range p.DomainName {
			values.Add("domainname", domainName)
		}
	}
	if len(p.IDPName) > 0 {
		for _, idpName := range p.IDPName {
			values.Add("idpname", idpName)
		}
	}

	return values
}

// ReadAllPagesWithPagination reads all pages using the standard zidentity pagination response format
func ReadAllPagesWithPagination[T any](ctx context.Context, client *zscaler.Client, endpoint string, queryParams *PaginationQueryParams) ([]T, error) {
	var allRecords []T
	var currentOffset int

	if queryParams == nil {
		queryParams = &PaginationQueryParams{
			Limit: DefaultPaginationOptions.DefaultPageSize,
		}
	}

	for {
		// Set current offset
		queryParams.Offset = currentOffset

		// Build URL with query parameters
		urlValues := queryParams.ToURLValues()
		fullURL := endpoint
		if len(urlValues) > 0 {
			fullURL = fmt.Sprintf("%s?%s", endpoint, urlValues.Encode())
		}

		// Fetch current page
		var response PaginationResponse[T]
		err := client.Read(ctx, fullURL, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page at offset %d: %w", currentOffset, err)
		}

		// Add records to results
		allRecords = append(allRecords, response.Records...)

		// Check if we've reached the end
		if len(response.Records) < queryParams.Limit || response.NextLink == "" {
			break
		}

		// Update offset for next iteration
		currentOffset += len(response.Records)
	}

	return allRecords, nil
}

// ReadPageWithPagination reads a single page using the standard zidentity pagination response format
func ReadPageWithPagination[T any](ctx context.Context, client *zscaler.Client, endpoint string, queryParams *PaginationQueryParams) (*PaginationResponse[T], error) {
	if queryParams == nil {
		queryParams = &PaginationQueryParams{
			Limit: DefaultPaginationOptions.DefaultPageSize,
		}
	}

	// Build URL with query parameters
	urlValues := queryParams.ToURLValues()
	fullURL := endpoint
	if len(urlValues) > 0 {
		fullURL = fmt.Sprintf("%s?%s", endpoint, urlValues.Encode())
	}

	// Fetch page
	var response PaginationResponse[T]
	err := client.Read(ctx, fullURL, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}

	return &response, nil
}

// ReadAllPagesWithCursor reads all pages using cursor-based pagination (next_link/prev_link)
func ReadAllPagesWithCursor[T any](ctx context.Context, client *zscaler.Client, endpoint string, queryParams *PaginationQueryParams) ([]T, error) {
	var allRecords []T

	if queryParams == nil {
		queryParams = &PaginationQueryParams{
			Limit: DefaultPaginationOptions.DefaultPageSize,
		}
	}

	// Build initial URL with query parameters
	urlValues := queryParams.ToURLValues()
	currentURL := endpoint
	if len(urlValues) > 0 {
		currentURL = fmt.Sprintf("%s?%s", endpoint, urlValues.Encode())
	}

	for {
		// Fetch current page
		var response PaginationResponse[T]
		err := client.Read(ctx, currentURL, &response)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page: %w", err)
		}

		// Add records to results
		allRecords = append(allRecords, response.Records...)

		// Check if we've reached the end
		if response.NextLink == "" {
			break
		}

		// Use next_link for next iteration
		currentURL = response.NextLink
	}

	return allRecords, nil
}

// BuildEndpointWithParams builds an endpoint URL with query parameters
func BuildEndpointWithParams(endpoint string, queryParams *PaginationQueryParams) string {
	if queryParams == nil {
		return endpoint
	}

	urlValues := queryParams.ToURLValues()
	if len(urlValues) == 0 {
		return endpoint
	}

	return fmt.Sprintf("%s?%s", endpoint, urlValues.Encode())
}

// ParsePaginationResponse parses a JSON response into PaginationResponse
func ParsePaginationResponse[T any](data []byte) (*PaginationResponse[T], error) {
	var response PaginationResponse[T]
	err := json.Unmarshal(data, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pagination response: %w", err)
	}
	return &response, nil
}

// WithLoginName adds login name filter(s) to the query parameters
func (p *PaginationQueryParams) WithLoginName(loginNames []string) *PaginationQueryParams {
	p.LoginName = loginNames
	return p
}

// WithLoginNameLike adds a login name like filter to the query parameters
func (p *PaginationQueryParams) WithLoginNameLike(loginName string) *PaginationQueryParams {
	p.LoginNameLike = loginName
	return p
}

// WithDisplayNameLike adds a display name like filter to the query parameters
func (p *PaginationQueryParams) WithDisplayNameLike(displayName string) *PaginationQueryParams {
	p.DisplayNameLike = displayName
	return p
}

// WithPrimaryEmailLike adds a primary email like filter to the query parameters
func (p *PaginationQueryParams) WithPrimaryEmailLike(email string) *PaginationQueryParams {
	p.PrimaryEmailLike = email
	return p
}

// WithDomainName adds domain name filter(s) to the query parameters
func (p *PaginationQueryParams) WithDomainName(domainNames []string) *PaginationQueryParams {
	p.DomainName = domainNames
	return p
}

// WithIDPName adds IDP name filter(s) to the query parameters
func (p *PaginationQueryParams) WithIDPName(idpNames []string) *PaginationQueryParams {
	p.IDPName = idpNames
	return p
}
