package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
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

func getAllPagesGenericWithCustomFilters[T any](ctx context.Context, client *zscaler.Client, relativeURL string, page, pageSize int, filters Filter) (int, []T, *http.Response, error) {
	var v struct {
		TotalPages interface{} `json:"totalPages"`
		List       []T         `json:"list"`
	}
	resp, err := client.NewRequestDo(ctx, "GET", relativeURL, Pagination{
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
	searchQuery = url.QueryEscape(searchQuery)
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

	// Updated filter search: replace spaces with '&' for the API's search query format
	filters.Search = strings.ReplaceAll(filters.Search, " ", "&")

	totalPages, result, resp, err := getAllPagesGenericWithCustomFilters[T](ctx, client, relativeURL, 1, DefaultPageSize, filters)
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
