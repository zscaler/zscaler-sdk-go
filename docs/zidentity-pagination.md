# ZIdentity Pagination System

This document describes the centralized pagination system implemented for the zidentity service in the Zscaler Go SDK.

## Overview

The zidentity service uses a standardized pagination response format across all endpoints. The pagination system is designed to handle:

- Standard offset/limit pagination
- Cursor-based pagination using `next_link` and `prev_link`
- Common query parameters like filtering by name and excluding dynamic groups
- Automatic handling of pagination across multiple pages

## Pagination Response Structure

All zidentity endpoints return responses in this standardized format:

```json
{
  "results_total": 0,
  "pageOffset": 0,
  "pageSize": 0,
  "next_link": "string",
  "prev_link": "string",
  "records": [
    {
      "name": "string",
      "description": "string",
      "id": "string",
      "source": "UI",
      "idp": {
        "id": "string",
        "name": "string",
        "displayName": "string"
      },
      "isDynamicGroup": true,
      "adminEntitlementEnabled": true,
      "serviceEntitlementEnabled": true
    }
  ]
}
```

## Common Query Parameters

The system supports these common query parameters:

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `offset` | integer | Starting point for pagination | `offset=50` |
| `limit` | integer | Maximum records per request (0-1000) | `limit=100` |
| `name[like]` | string | Case-insensitive partial name match | `name[like]=admin` |
| `excludedynamicgroups` | boolean | Exclude dynamic groups from results | `excludedynamicgroups=true` |

## Usage Examples

### Basic Usage

```go
import (
    "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/service/common"
    "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zidentity/service/groups"
)

// Get all groups with default pagination
allGroups, err := groups.GetAll(ctx, service, nil)

// Get groups with custom pagination
queryParams := common.NewPaginationQueryParams(50)
queryParams.WithOffset(0)
customGroups, err := groups.GetAll(ctx, service, &queryParams)
```

### Simple Function Usage

The API provides clear, simple functions for different use cases:

```go
// Get all groups (no parameters)
allGroups, err := groups.GetAll(ctx, service, nil)

// Get groups with pagination parameters
queryParams := common.NewPaginationQueryParams(50)
queryParams.WithOffset(0)
customGroups, err := groups.GetAll(ctx, service, &queryParams)

// Get groups by name (searches through all paginated data)
groupsByName, err := groups.GetByName(ctx, service, "admin")

// Get groups with filtering parameters
queryParams := common.NewPaginationQueryParams(100)
queryParams.WithExcludeDynamicGroups(true)
nonDynamicGroups, err := groups.GetAll(ctx, service, &queryParams)

// Advanced filtering with multiple parameters
advancedParams := common.NewPaginationQueryParams(25)
advancedParams.WithExcludeDynamicGroups(true)
advancedParams.WithOffset(50)
advancedGroups, err := groups.GetAll(ctx, service, &advancedParams)
```

### Single Page Retrieval

```go
// Get a single page of results
pageParams := common.NewPaginationQueryParams(10)
pageParams.WithOffset(0)
pageResponse, err := groups.GetPage(ctx, service, &pageParams)

// Access pagination metadata
fmt.Printf("Total: %d, Offset: %d, PageSize: %d\n", 
    pageResponse.ResultsTotal, pageResponse.PageOffset, pageResponse.PageSize)
fmt.Printf("Records: %d\n", len(pageResponse.Records))
if pageResponse.NextLink != "" {
    fmt.Printf("Next page: %s\n", pageResponse.NextLink)
}
```

### Cursor-Based Pagination

```go
// Use cursor-based pagination (follows next_link automatically)
cursorGroups, err := groups.GetAllWithCursor(ctx, service, nil)
```

## Key Components

### PaginationResponse[T]

Generic structure for paginated responses:

```go
type PaginationResponse[T any] struct {
    ResultsTotal int    `json:"results_total,omitempty"`
    PageOffset   int    `json:"pageOffset,omitempty"`
    PageSize     int    `json:"pageSize,omitempty"`
    NextLink     string `json:"next_link,omitempty"`
    PrevLink     string `json:"prev_link,omitempty"`
    Records      []T    `json:"records,omitempty"`
}
```

### PaginationQueryParams

Structure for query parameters with fluent interface:

```go
type PaginationQueryParams struct {
    Offset              int    `url:"offset,omitempty"`
    Limit               int    `url:"limit,omitempty"`
    NameLike            string `url:"name[like],omitempty"`
    ExcludeDynamicGroups bool   `url:"excludedynamicgroups,omitempty"`
}
```

### Helper Functions

- `NewPaginationQueryParams(pageSize int)` - Create new query params
- `WithNameFilter(name string)` - Add name filter
- `WithExcludeDynamicGroups(exclude bool)` - Exclude dynamic groups
- `WithOffset(offset int)` - Set pagination offset
- `WithLimit(limit int)` - Set pagination limit
- `ToURLValues()` - Convert to URL query parameters

### Pagination Functions

- `ReadAllPagesWithPagination[T]()` - Read all pages using offset/limit
- `ReadPageWithPagination[T]()` - Read single page
- `ReadAllPagesWithCursor[T]()` - Read all pages using cursor links
- `BuildEndpointWithParams()` - Build URL with query parameters
- `ParsePaginationResponse[T]()` - Parse JSON response

## Configuration

Default pagination options:

```go
var DefaultPaginationOptions = PaginationOptions{
    DefaultPageSize: 100,
    MaxPageSize:     1000,
    UseCursor:       false,
}
```

## Best Practices

1. **Use appropriate page sizes**: Default is 100, max is 1000
2. **Handle errors gracefully**: Always check for errors from pagination functions
3. **Use filtering when possible**: Reduce data transfer with appropriate filters
4. **Combine parameters flexibly**: All parameters are optional and can be mixed and matched
5. **Consider cursor pagination**: For large datasets, cursor-based pagination may be more efficient
6. **Cache results when appropriate**: For frequently accessed data, consider caching
7. **Start simple**: Begin with minimal parameters and add complexity as needed

## Extending the System

To add new query parameters:

1. Add the field to `PaginationQueryParams`
2. Add a fluent method (e.g., `WithNewFilter()`)
3. Update `ToURLValues()` to include the new parameter
4. Update documentation

To add new pagination functions:

1. Follow the existing pattern in `common/common.go`
2. Use generics for type safety
3. Include proper error handling
4. Add comprehensive documentation

## Migration from Other Services

If migrating from other Zscaler services (ZIA, ZPA, etc.):

1. Replace service-specific pagination with `PaginationQueryParams`
2. Use `PaginationResponse[T]` instead of custom response structures
3. Replace manual pagination loops with `ReadAllPagesWithPagination[T]()`
4. Update query parameter building to use the fluent interface

## Error Handling

The pagination system includes comprehensive error handling:

- Invalid page sizes are automatically clamped to valid ranges
- Network errors are wrapped with context
- JSON parsing errors are handled gracefully
- Empty responses are handled properly

## Performance Considerations

- Default page size of 100 provides good balance between performance and memory usage
- Cursor-based pagination can be more efficient for large datasets
- Filtering reduces data transfer and processing time
- Consider using goroutines for parallel page fetching in some scenarios 