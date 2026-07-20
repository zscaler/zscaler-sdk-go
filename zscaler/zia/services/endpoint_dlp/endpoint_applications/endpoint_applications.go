package endpoint_applications

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

const (
	endPointApplicationsEndpoint = "/zia/api/v1/endPointApplications"
)

type ApplicationPolicies struct {
	RuleName string `json:"ruleName,omitempty"`
	RuleType string `json:"ruleType,omitempty"`
	ID       string `json:"id,omitempty"`
}

// GetApplicationCountFilterOptions holds the optional query parameters supported
// by GetApplicationCount and GetCloudAppsCount.
type GetApplicationCountFilterOptions struct {
	// Search is the search string used to match against application names. Optional.
	Search string

	// OsType filters the results by operating system (e.g., Windows OS and Mac OS). Optional.
	OsType string

	// ApplicationType filters the results by application type (e.g., Custom,
	// Discovered, and Well-Known). Optional.
	ApplicationType string
}

// GetApplicationCount retrieves the count of all endpoint applications. The
// search, osType, and applicationType parameters are optional. The API returns a
// simple integer count.
func GetApplicationCount(ctx context.Context, service *zscaler.Service, opts *GetApplicationCountFilterOptions) (int, error) {
	var count int
	endpoint := endPointApplicationsEndpoint + "/count"

	if q := applicationCountQuery(opts); q != "" {
		endpoint += "?" + q
	}

	if err := service.Client.Read(ctx, endpoint, &count); err != nil {
		return 0, fmt.Errorf("failed to retrieve endpoint application count: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning endpoint application count: %d", count)
	return count, nil
}

// GetCloudAppsCount retrieves the count of well-known and discovered endpoint
// applications as determined by the Zscaler service. The search, osType, and
// applicationType parameters are optional. The API returns a simple integer
// count.
func GetCloudAppsCount(ctx context.Context, service *zscaler.Service, opts *GetApplicationCountFilterOptions) (int, error) {
	var count int
	endpoint := endPointApplicationsEndpoint + "/cloudApps/count"

	if q := applicationCountQuery(opts); q != "" {
		endpoint += "?" + q
	}

	if err := service.Client.Read(ctx, endpoint, &count); err != nil {
		return 0, fmt.Errorf("failed to retrieve cloud endpoint application count: %w", err)
	}

	service.Client.GetLogger().Printf("[DEBUG] Returning cloud endpoint application count: %d", count)
	return count, nil
}

// GetApplicationPolicies retrieves the list of policy rules currently associated
// with the endpoint application(s) identified by the given resource IDs. The
// resourceIDs values correspond to the resourceId field returned by
// GetCustomApps and are sent as repeated query parameters.
func GetApplicationPolicies(ctx context.Context, service *zscaler.Service, resourceIDs []int) ([]ApplicationPolicies, error) {
	var policies []ApplicationPolicies
	endpoint := endPointApplicationsEndpoint + "/policies"

	queryParams := url.Values{}
	for _, id := range resourceIDs {
		queryParams.Add("resourceId", strconv.Itoa(id))
	}
	if len(queryParams) > 0 {
		endpoint += "?" + queryParams.Encode()
	}

	err := service.Client.Read(ctx, endpoint, &policies)
	return policies, err
}

// applicationCountQuery builds the shared query string used by the count
// endpoints from the optional filter options.
func applicationCountQuery(opts *GetApplicationCountFilterOptions) string {
	queryParams := url.Values{}
	if opts != nil {
		if opts.Search != "" {
			queryParams.Set("search", opts.Search)
		}
		if opts.OsType != "" {
			queryParams.Set("osType", opts.OsType)
		}
		if opts.ApplicationType != "" {
			queryParams.Set("applicationType", opts.ApplicationType)
		}
	}
	return queryParams.Encode()
}

// GetCategoriesWithNonEmptyAppsFilterOptions holds the optional query parameters
// supported by GetCategoriesWithNonEmptyApps.
type GetCategoriesWithNonEmptyAppsFilterOptions struct {
	// Search is the search string used to match against application names. Optional.
	Search string

	// OsType filters the results by operating system (e.g., Windows OS and Mac OS). Optional.
	OsType string
}

// GetCategoriesWithNonEmptyApps retrieves the categories that currently have
// endpoint applications grouped within them. The search and osType parameters
// are optional. The API returns a simple list of category name strings and
// supports page/pageSize pagination.
func GetCategoriesWithNonEmptyApps(ctx context.Context, service *zscaler.Service, opts *GetCategoriesWithNonEmptyAppsFilterOptions) ([]string, error) {
	var categories []string
	endpoint := endPointApplicationsEndpoint + "/getCategoriesWithNonEmptyApps"

	if opts != nil {
		queryParams := url.Values{}
		if opts.Search != "" {
			queryParams.Set("search", opts.Search)
		}
		if opts.OsType != "" {
			queryParams.Set("osType", opts.OsType)
		}
		if len(queryParams) > 0 {
			endpoint += "?" + queryParams.Encode()
		}
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &categories)
	return categories, err
}

// GetAllEndpointApplications retrieves the list of endpoint applications. The
// search, osType, and applicationType parameters are optional; search matches
// against endpoint application names. Pagination is handled automatically via
// ReadAllPages.
func GetAllEndpointApplications(ctx context.Context, service *zscaler.Service, opts *GetApplicationCountFilterOptions) ([]common.EndPointApplications, error) {
	var applications []common.EndPointApplications
	endpoint := endPointApplicationsEndpoint

	if q := applicationCountQuery(opts); q != "" {
		endpoint += "?" + q
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &applications, 5000)
	return applications, err
}

// GetAllEndpointApplications retrieves the list of endpoint applications. The
// search, osType, and applicationType parameters are optional; search matches
// against endpoint application names. Pagination is handled automatically via
// ReadAllPages.
func GetAllEndpointApplicationsLite(ctx context.Context, service *zscaler.Service, opts *GetApplicationCountFilterOptions) ([]common.EndPointApplications, error) {
	var applications []common.EndPointApplications
	endpoint := endPointApplicationsEndpoint + "/lite"

	if q := applicationCountQuery(opts); q != "" {
		endpoint += "?" + q
	}

	err := common.ReadAllPages(ctx, service.Client, endpoint, &applications, 5000)
	return applications, err
}
