package common

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/zpa"
)

const (
	DefaultPageSize = 100
)

type Pagination struct {
	PageSize int    `json:"pagesize,omitempty" url:"pagesize,omitempty"`
	Page     int    `json:"page,omitempty" url:"page,omitempty"`
	Search   string `json:"-" url:"-"`
	Search2  string `json:"search,omitempty" url:"search,omitempty"`
}

type DeleteApplicationQueryParams struct {
	ForceDelete bool `json:"forceDelete,omitempty" url:"forceDelete,omitempty"`
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

type AssociatedProfileNames struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
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

func getAllPagesGenericWithCustomFilters[T any](client *zpa.Client, relativeURL string, page, pageSize int, searchQuery string, buildFilters func(pageSize, page int, searchQuery string) interface{}) (int, []T, *http.Response, error) {
	var v struct {
		TotalPages interface{} `json:"totalPages"`
		List       []T         `json:"list"`
	}
	resp, err := client.NewRequestDo("GET", relativeURL, buildFilters(pageSize, page, searchQuery), nil, &v)
	if err != nil {
		return 0, nil, resp, err
	}

	pages := fmt.Sprintf("%v", v.TotalPages)
	totalPages, _ := strconv.Atoi(pages)

	return totalPages, v.List, resp, nil
}

func getAllPagesGeneric[T any](client *zpa.Client, relativeURL string, page, pageSize int, searchQuery string) (int, []T, *http.Response, error) {
	return getAllPagesGenericWithCustomFilters[T](
		client,
		relativeURL,
		page,
		pageSize,
		searchQuery,
		func(pageSize, page int, searchQuery string) interface{} {
			return Pagination{PageSize: pageSize, Page: page, Search2: searchQuery}
		},
	)
}

// GetAllPagesGeneric fetches all resources instead of just one single page
func GetAllPagesGeneric[T any](client *zpa.Client, relativeURL, searchQuery string) ([]T, *http.Response, error) {
	totalPages, result, resp, err := getAllPagesGeneric[T](client, relativeURL, 1, DefaultPageSize, searchQuery)
	if err != nil {
		return nil, resp, err
	}
	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGeneric[T](client, relativeURL, page, DefaultPageSize, searchQuery)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}

	return result, resp, nil
}

// GetAllPagesGenericWithCustomFilters fetches all resources instead of just one single page
func GetAllPagesGenericWithCustomFilters[T any](client *zpa.Client, relativeURL, searchQuery string, buildFilters func(pageSize, page int, searchQuery string) interface{}) ([]T, *http.Response, error) {
	totalPages, result, resp, err := getAllPagesGenericWithCustomFilters[T](client, relativeURL, 1, DefaultPageSize, searchQuery, buildFilters)
	if err != nil {
		return nil, resp, err
	}
	var l []T
	for page := 2; page <= totalPages; page++ {
		totalPages, l, resp, err = getAllPagesGenericWithCustomFilters[T](client, relativeURL, page, DefaultPageSize, searchQuery, buildFilters)
		if err != nil {
			return nil, resp, err
		}
		result = append(result, l...)
	}

	return result, resp, nil
}
