package common

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

const (
	DefaultPageSize = 50
	MaxPageSize     = 5000
)

const (
	DeviceTypeIOS     = 1
	DeviceTypeAndroid = 2
	DeviceTypeWindows = 3
	DeviceTypeMacOS   = 4
	DeviceTypeLinux   = 5
)

var deviceTypeNameToID = map[string]int{
	"ios":     DeviceTypeIOS,
	"android": DeviceTypeAndroid,
	"windows": DeviceTypeWindows,
	"macos":   DeviceTypeMacOS,
	"linux":   DeviceTypeLinux,
}

var deviceTypeIDToName = map[int]string{
	DeviceTypeIOS:     "iOS",
	DeviceTypeAndroid: "Android",
	DeviceTypeWindows: "Windows",
	DeviceTypeMacOS:   "macOS",
	DeviceTypeLinux:   "Linux",
}

// GetDeviceTypeByName converts a device type name (e.g., "windows", "macOS")
// to its corresponding API integer value. It also accepts numeric strings
// (e.g., "3") for pass-through. Returns 0 and nil for empty input.
func GetDeviceTypeByName(name string) (int, error) {
	if name == "" {
		return 0, nil
	}
	if val, ok := deviceTypeNameToID[strings.ToLower(name)]; ok {
		return val, nil
	}
	if val, err := strconv.Atoi(name); err == nil {
		if _, ok := deviceTypeIDToName[val]; ok {
			return val, nil
		}
		return 0, fmt.Errorf("invalid device type ID: %d. Valid IDs: 1 (iOS), 2 (Android), 3 (Windows), 4 (macOS), 5 (Linux)", val)
	}
	return 0, fmt.Errorf("invalid device type: %q. Valid values: ios, android, windows, macos, linux (or 1-5)", name)
}

// GetDeviceTypeName converts a device type integer to its display name.
func GetDeviceTypeName(id int) string {
	if name, ok := deviceTypeIDToName[id]; ok {
		return name
	}
	return ""
}

type Pagination struct {
	PageSize int `json:"pagesize,omitempty" url:"pagesize,omitempty"`
	Page     int `json:"page,omitempty" url:"page,omitempty"`
}

// QueryParams is the unified query parameter struct for all ZCC GET list
// endpoints. Fields are serialized via go-querystring url tags; zero-value
// fields are omitted automatically.
type QueryParams struct {
	Page       int    `url:"page,omitempty"`
	PageSize   int    `url:"pageSize,omitempty"`
	Search     string `url:"search,omitempty"`
	SearchType string `url:"searchType,omitempty"`
	DeviceType int    `url:"deviceType,omitempty"`
	Username   string `url:"username,omitempty"`
	OsType     string `url:"osType,omitempty"`
	UserType   string `url:"userType,omitempty"`
}

// NewPagination creates a Pagination with bounded page size.
func NewPagination(pageSize int) Pagination {
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Pagination{PageSize: pageSize}
}

func clampPageSize(pageSize int) int {
	if pageSize <= 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

// ReadAllPages iterates through all pages of a paginated ZCC GET endpoint
// and returns the aggregated results. The QueryParams struct is passed
// directly to NewZccRequestDo so go-querystring handles serialization.
func ReadAllPages[T any](ctx context.Context, client *zscaler.Client, endpoint string, params QueryParams, pageSize int) ([]T, error) {
	pageSize = clampPageSize(pageSize)
	params.PageSize = pageSize
	params.Page = 1

	var allResults []T
	for {
		var pageResults []T
		_, err := client.NewZccRequestDo(ctx, "GET", endpoint, params, nil, &pageResults)
		if err != nil {
			return nil, err
		}

		allResults = append(allResults, pageResults...)

		if len(pageResults) < pageSize {
			break
		}
		params.Page++
	}

	filtered, err := zscaler.ApplyJMESPathFromContext(ctx, allResults)
	if err != nil {
		return nil, err
	}

	return filtered, nil
}

// ReadPage fetches a single page from a paginated ZCC GET endpoint.
// Useful when the caller manages pagination externally or only needs one page.
func ReadPage[T any](ctx context.Context, client *zscaler.Client, endpoint string, params QueryParams) ([]T, error) {
	if params.PageSize <= 0 {
		params.PageSize = DefaultPageSize
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	var results []T
	_, err := client.NewZccRequestDo(ctx, "GET", endpoint, params, nil, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
