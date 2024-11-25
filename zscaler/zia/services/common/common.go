package common

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/SecurityGeekIO/zscaler-sdk-go/v3/zscaler"
)

const pageSize = 1000

type IDNameExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type IDExtensions struct {
	ID         int                    `json:"id,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type IDName struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type IDCustom struct {
	ID int `json:"id,omitempty"`
}

type ZPAAppSegments struct {
	// A unique identifier assigned to the Application Segment
	ID int `json:"id"`

	// The name of the Application Segment
	Name string `json:"name,omitempty"`

	// Indicates the external ID. Applicable only when this reference is of an external entity.
	ExternalID string `json:"externalId"`
}

type UserGroups struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
}

type UserDepartment struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	IdpID    int    `json:"idp_id,omitempty"`
	Comments string `json:"comments,omitempty"`
	Deleted  bool   `json:"deleted,omitempty"`
}

type DeviceGroups struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Devices struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type IDNameWorkloadGroup struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DatacenterSearchParameters struct {
	RoutableIP                bool
	WithinCountryOnly         bool
	IncludePrivateServiceEdge bool
	IncludeCurrentVips        bool
	SourceIp                  string
	Latitude                  float64
	Longitude                 float64
	Subcloud                  string
}

type SandboxRSS struct {
	Risk             string `json:"Risk,omitempty"`
	Signature        string `json:"Signature,omitempty"`
	SignatureSources string `json:"SignatureSources,omitempty"`
}

type Order struct {
	On string `json:"on,omitempty"`
	By string `json:"by,omitempty"`
}

type DataConsumed struct {
	Min int `json:"min,omitempty"`
	Max int `json:"max,omitempty"`
}

// GetPageSize returns the page size.
func GetPageSize() int {
	return pageSize
}

func ReadAllPages[T any](ctx context.Context, client *zscaler.Client, endpoint string, list *[]T) error {
	if list == nil {
		return nil
	}
	page := 1
	if !strings.Contains(endpoint, "?") {
		endpoint += "?"
	}

	for {
		pageItems := []T{}
		err := client.Read(ctx, fmt.Sprintf("%s&pageSize=%d&page=%d", endpoint, pageSize, page), &pageItems)
		if err != nil {
			return err
		}
		*list = append(*list, pageItems...)
		if len(pageItems) < pageSize {
			break
		}
		page++
	}
	return nil
}

func ReadPage[T any](ctx context.Context, client *zscaler.Client, endpoint string, page int, list *[]T) error {
	if list == nil {
		return nil
	}

	// Parse the endpoint into a URL.
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("could not parse endpoint URL: %w", err)
	}

	// Get the existing query parameters and add new ones.
	q := u.Query()
	q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	q.Set("page", fmt.Sprintf("%d", page))

	// Set the URL's RawQuery to the encoded query parameters.
	u.RawQuery = q.Encode()

	// Convert the URL back to a string and read the page.
	pageItems := []T{}
	err = client.Read(ctx, u.String(), &pageItems)
	if err != nil {
		return err
	}
	*list = pageItems
	return nil
}

type (
	SortOrder string
	SortField string
)

const (
	ASCSortOrder          SortOrder = "ASC"
	DESCSortOrder                   = "DESC"
	IDSortField           SortField = "id"
	NameSortField                   = "name"
	CreationTimeSortField           = "creationTime"
	ModifiedTimeSortField           = "modifiedTime"
)

func GetSortParams(sortBy SortField, sortOrder SortOrder) string {
	params := ""
	if sortBy != "" {
		params = "sortBy=" + string(sortBy)
	}
	if sortOrder != "" {
		if params != "" {
			params += "&"
		}
		params += "sortOrder=" + string(sortOrder)
	}
	return params
}
