package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

// IntOrString is a numeric field that the ZCC API is inconsistent about on
// the wire: some endpoints accept and return integers, others return the
// same value quoted as a string ("0", "12", ...), and a few fields come
// back as empty strings when the policy has never been touched. Modelling
// these as IntOrString lets the SDK:
//
//   - serialize as a JSON number (the form every ZCC POST/PUT endpoint
//     actually accepts; quoted-int strings are rejected with a silent
//     {"success":"false","id":0} on the web_policy /edit endpoint),
//   - unmarshal from either a JSON number, a numeric string, a null, or
//     an empty string (treated as 0).
//
// The underlying type is int, so callers operate on it like a normal int
// (assign with IntOrString(n), read with int(v), etc.).
type IntOrString int

// MarshalJSON emits the value as a JSON number. The default int marshal
// would do the same; the explicit method documents the contract.
func (i IntOrString) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(i))
}

// UnmarshalJSON accepts a JSON integer (12), a JSON float that happens to
// carry an integer value (1.0, -1.0 — the shape json.Marshal produces for a
// float64 that the ZCC listByCompany endpoint uses for fields like
// `ruleOrder` and `logMode`), a quoted numeric string ("12"), null, or
// the empty string. Anything else returns an error.
func (i *IntOrString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) || bytes.Equal(trimmed, []byte(`""`)) {
		*i = 0
		return nil
	}
	// First try a strict int decode — handles "12" and 12.
	var num int
	if err := json.Unmarshal(trimmed, &num); err == nil {
		*i = IntOrString(num)
		return nil
	}
	// Some ZCC endpoints serialize integer values as JSON floats with a
	// trailing .0 (e.g. `1.0`, `-1.0`); Go's json package refuses to
	// decode those into int, so retry through float64 and truncate.
	var f float64
	if err := json.Unmarshal(trimmed, &f); err == nil {
		*i = IntOrString(int(f))
		return nil
	}
	var str string
	if err := json.Unmarshal(trimmed, &str); err == nil {
		if str == "" {
			*i = 0
			return nil
		}
		if parsed, err := strconv.Atoi(str); err == nil {
			*i = IntOrString(parsed)
			return nil
		}
		// Quoted floats like "1.0" — accept those too for symmetry with
		// the bare-number path above.
		if parsed, err := strconv.ParseFloat(str, 64); err == nil {
			*i = IntOrString(int(parsed))
			return nil
		}
	}
	return fmt.Errorf("invalid value for IntOrString: %s", string(data))
}

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

// QueryParams is the unified query parameter struct for ZCC v1 GET list
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

// QueryParamsV2 is the unified query parameter struct for ZCC v2 GET list
// endpoints. v2 endpoints differ from v1 in three ways: pagination is
// offset-based (skip + perPage) instead of page-based, the search field
// is named "keyword" instead of "search", and most endpoints layer one or
// more endpoint-specific filters (e.g. type for trusted-networks,
// platformType for zia-posture-profiles). Fields are serialized via
// go-querystring url tags; zero-value fields are omitted automatically.
//
// The union of optional filters is modelled here so the pagination helper
// has a single signature. Per-endpoint convenience structs in each service
// translate caller-facing options into the relevant subset.
type QueryParamsV2 struct {
	Skip    int    `url:"skip,omitempty"`
	PerPage int    `url:"perPage,omitempty"`
	Keyword string `url:"keyword,omitempty"`

	// Endpoint-specific filters. Each is honoured only by the endpoints
	// documented to accept it; they are inert on others.
	Type         string `url:"type,omitempty"`         // /trusted-networks
	PlatformType int    `url:"platformType,omitempty"` // /zia-posture-profiles
}

type ZCCResponse struct {
	Success string `json:"success"`
	Error   string `json:"error,omitempty"`
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

// PaginatedResponseV2 is the envelope shape used by the ZCC v2 API. Unlike
// the v1 endpoints — which return a bare JSON array — every v2 list endpoint
// wraps the records in an object that also carries pagination metadata:
//
//	{
//	  "items":  [ ... ],   // records on this page
//	  "total":  <int>,     // total records across all pages
//	  "offset": <int>,     // zero-based starting index of this page
//	  "limit":  <int>,     // max records per page (i.e. page size)
//	  "count":  <int>      // actual number of records in items on this page
//	}
//
// T is the per-item element type (e.g. TrustedNetworkV2). The envelope is
// public so callers of ReadPageV2 can inspect the metadata directly.
type PaginatedResponseV2[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
	Count  int `json:"count"`
}

// ReadAllPagesV2 iterates through every page of a paginated ZCC v2 GET
// endpoint and returns the aggregated items. Pagination is offset-based:
// Skip is advanced by PerPage after each successful page. To avoid hanging
// on a malformed envelope, termination uses whichever signal arrives first:
//
//  1. server-authoritative — collected results meet or exceed total
//  2. last-page heuristic  — count < limit (server returned a short page)
//  3. empty-page safety    — count == 0 or items is empty
//
// perPage is clamped to [DefaultPageSize, MaxPageSize]. JMESPath filtering
// from context is applied after aggregation, matching the v1 helper.
//
// DO NOT use this for v1 endpoints — those return a bare JSON array and the
// existing ReadAllPages helper is the correct call.
func ReadAllPagesV2[T any](ctx context.Context, client *zscaler.Client, endpoint string, params QueryParamsV2, perPage int) ([]T, error) {
	perPage = clampPageSize(perPage)
	params.PerPage = perPage
	if params.Skip < 0 {
		params.Skip = 0
	}

	var allResults []T
	for {
		var page PaginatedResponseV2[T]
		if _, err := client.NewZccRequestDo(ctx, "GET", endpoint, params, nil, &page); err != nil {
			return nil, err
		}

		allResults = append(allResults, page.Items...)

		// Stop on any termination signal. Each guard is independent so a
		// missing field (e.g. Total == 0) never traps us in a loop.
		if page.Count == 0 || len(page.Items) == 0 {
			break
		}
		if page.Limit > 0 && page.Count < page.Limit {
			break
		}
		if page.Total > 0 && len(allResults) >= page.Total {
			break
		}

		params.Skip += perPage
	}

	filtered, err := zscaler.ApplyJMESPathFromContext(ctx, allResults)
	if err != nil {
		return nil, err
	}
	return filtered, nil
}

// ReadPageV2 fetches a single page from a paginated ZCC v2 GET endpoint and
// returns the full envelope (items + total/offset/limit/count). Use this
// when the caller manages pagination externally, only needs one page, or
// wants the server-reported totals for display purposes.
func ReadPageV2[T any](ctx context.Context, client *zscaler.Client, endpoint string, params QueryParamsV2) (PaginatedResponseV2[T], error) {
	if params.PerPage <= 0 {
		params.PerPage = DefaultPageSize
	}
	if params.Skip < 0 {
		params.Skip = 0
	}

	var page PaginatedResponseV2[T]
	if _, err := client.NewZccRequestDo(ctx, "GET", endpoint, params, nil, &page); err != nil {
		return PaginatedResponseV2[T]{}, err
	}
	return page, nil
}
