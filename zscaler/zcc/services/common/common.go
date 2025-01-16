package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	DefaultPageSize = 30
	MaxPageSize     = 5000
)

type Pagination struct {
	PageSize int `json:"pagesize,omitempty" url:"pagesize,omitempty"`
	Page     int `json:"page,omitempty" url:"page,omitempty"`
}

type QueryParams struct {
	Page       int    `url:"page,omitempty"`
	PageSize   int    `url:"pageSize,omitempty"`
	Search     string `url:"search,omitempty"`
	SearchType string `url:"searchType,omitempty"`
	DeviceType int    `url:"deviceType,omitempty"`
}

// NewPagination creates a new Pagination struct with provided page size
// If page size is less than or equal to 0, it uses the default page size
// If page size is greater than the max page size, it uses the max page size
func NewPagination(pageSize int) Pagination {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{PageSize: pageSize}
}

func queryParamsToURLValues(params interface{}) (url.Values, error) {
	values := url.Values{}
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	var mapParams map[string]interface{}
	if err := json.Unmarshal(data, &mapParams); err != nil {
		return nil, err
	}
	for key, value := range mapParams {
		if value != nil {
			values.Set(key, fmt.Sprintf("%v", value))
		}
	}
	return values, nil
}

func ReadAllPages[T any](ctx context.Context, client *zscaler.Client, endpoint string, queryParams interface{}, pageSize int) ([]T, error) {
	pagination := NewPagination(pageSize)
	var allResults []T
	page := 1

	for {
		var pageResults []T

		// Update the query parameters with pagination details
		q := url.Values{}
		if queryParams != nil {
			// Convert queryParams to URL values if needed
			queryString, err := queryParamsToURLValues(queryParams)
			if err != nil {
				return nil, fmt.Errorf("failed to parse query params: %w", err)
			}
			q = queryString
		}

		// Add pagination parameters
		q.Set("pageSize", fmt.Sprintf("%d", pagination.PageSize))
		q.Set("page", fmt.Sprintf("%d", page))

		// Include the optional `search` parameter if present in queryParams
		if searchValue, ok := queryParams.(map[string]interface{})["search"]; ok {
			q.Set("search", fmt.Sprintf("%v", searchValue))
		}

		// Build the final endpoint URL
		fullURL := fmt.Sprintf("%s?%s", endpoint, q.Encode())

		// Fetch the current page
		_, err := client.NewRequestDo(ctx, "GET", fullURL, nil, nil, &pageResults)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, pageResults...)

		// Break if the number of results is less than the page size (last page)
		if len(pageResults) < pagination.PageSize {
			break
		}
		page++
	}

	return allResults, nil
}
