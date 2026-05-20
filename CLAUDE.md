# Zscaler Go SDK

Go SDK for the Zscaler Zero Trust Exchange. Module: `github.com/zscaler/zscaler-sdk-go/v3`. Services: ZIA, ZPA, ZCC, ZDX, ZTW (Cloud & Branch Connector), ZID (ZIdentity).

## Architecture Overview

```
zscaler/
├── oneapiclient.go         # Client struct, OAuth2 token management, VERSION constant
├── oneapiconfig.go         # Configuration struct, functional options (WithClientID, WithCache, etc.)
├── service.go              # Service struct (shared across all clouds), microtenant, sort
├── jmespath.go             # Cross-service JMESPath client-side filtering utility
├── user_agent.go           # User-Agent header builder
├── utils.go                # Shared utilities
├── ziarequests.go          # ZIA HTTP methods: Read, Create, UpdateWithPut, Delete
├── zparequests.go          # ZPA HTTP methods: NewRequestDo (generic 6-param)
├── zccrequests.go          # ZCC HTTP methods: NewZccRequestDo (manual response handling)
├── zdxrequests.go          # ZDX HTTP methods: NewRequestDo (same as ZPA)
├── ztwrequests.go          # ZTW HTTP methods: ReadResource, CreateResource, UpdateWithPutResource, DeleteResource
├── errorx/                 # Error types: ErrorResponse, IsObjectNotFound(), IsLimitExceeded()
├── zia/
│   ├── v2_client.go        # ZIA client initialization, base URL routing
│   ├── v2_config.go        # ZIA-specific config
│   └── services/
│       ├── common/common.go  # ReadAllPages, ReadPage, SCIM pagination, IDNameExtensions types
│       └── <service>/        # One package per resource (e.g., firewallpolicies/, urlcategories/)
├── zpa/
│   ├── v2_client.go
│   ├── v2_config.go
│   └── services/
│       ├── common/common.go  # GetAllPagesGenericWithCustomFilters, Filter, Pagination, CommonIDName types
│       └── <service>/
├── zcc/
│   ├── v2_client.go
│   ├── v2_config.go
│   └── services/
│       ├── common/common.go  # ReadAllPages[T] / ReadPage[T] (v1 bare arrays), ReadAllPagesV2[T] / ReadPageV2[T] / PaginatedResponseV2[T] (v2 envelope), QueryParams
│       └── <service>/
├── zdx/
│   └── services/
│       ├── common/common.go  # GetFromToFilters (time range, offset/limit, cursor)
│       └── <service>/        # Per-domain: reports/, alerts/, inventory/
├── ztw/
│   ├── v2_client.go
│   ├── v2_config.go
│   └── services/
│       ├── common/common.go  # ReadAllPages (fixed 1000 pageSize)
│       └── <service>/
├── zid/
│   └── services/
│       ├── common/common.go  # ReadAllPagesWithPagination, ReadAllPagesWithCursor, PaginationResponse[T]
│       └── <service>/
└── zwa/
    └── services/
        ├── common/common.go  # ReadAllPages (cursor-based with TotalPages)
        └── <service>/
```

### Request Flow

1. Caller creates `Configuration` via functional options (`WithClientID`, `WithClientSecret`, etc.)
2. `NewOneAPIClient(config)` → `*Service` with lazy per-cloud client initialization
3. Service function calls `service.Client.<Method>()` which routes to the correct cloud HTTP client
4. OAuth2 token is obtained on first request, auto-refreshed on 401
5. Rate limiter enforces per-cloud limits (e.g., ZIA: 20 GET/10s, 10 POST/10s)
6. Response is unmarshaled into the target struct; errors propagate as `errorx.ErrorResponse`
7. GET responses are cached when `WithCache(true)` is set; mutations auto-invalidate

## Cloud Service Matrix

| Cloud | Package | ID Type | Request Methods | Pagination | Endpoint Pattern |
|-------|---------|---------|-----------------|------------|-----------------|
| **ZIA** | `zscaler/zia/services/` | `int` | `Read/Create/UpdateWithPut/Delete` | `common.ReadAllPages` (page/pageSize, stop at `len < pageSize`) | `/zia/api/v1/<resource>` |
| **ZPA** | `zscaler/zpa/services/` | `string` | `NewRequestDo` (6 params) | `common.GetAllPagesGenericWithCustomFilters` (totalPages envelope) | `/zpa/mgmtconfig/v1/admin/customers/<customerID>/<resource>` |
| **ZCC** | `zscaler/zcc/services/` | varies | `NewZccRequestDo` (manual response) | v1: `common.ReadAllPages[T]` (bare array, page/pageSize) — v2: `common.ReadAllPagesV2[T]` (items envelope with total/offset/limit/count) | `/zcc/papi/public/v1/<resource>` and `/zcc/papi/public/v2/<resource>` |
| **ZDX** | `zscaler/zdx/services/` | `int` | `NewRequestDo` (same as ZPA) | cursor-based (`next_offset` token) | `/zdx/api/v1/<resource>` |
| **ZTW** | `zscaler/ztw/services/` | `int` | `ReadResource/CreateResource/UpdateWithPutResource/DeleteResource` | `common.ReadAllPages` (fixed 1000) | `/ztw/api/v1/<resource>` |
| **ZID** | `zscaler/zid/services/` | `string` | `Read/Create/UpdateWithPut/Delete` (same as ZIA) | `common.ReadAllPagesWithPagination` (offset/limit + `next_link`) | `/admin/api/v1/<resource>` |

## CRUD Patterns by Cloud

### Function Signature Convention

All SDK service functions are **package-level functions** (not methods on a struct). They always take `ctx context.Context` first and `service *zscaler.Service` second.

### ZIA

```go
func Get(ctx context.Context, service *zscaler.Service, id int) (*T, error)
func Create(ctx context.Context, service *zscaler.Service, resource *T) (*T, *http.Response, error)
func Update(ctx context.Context, service *zscaler.Service, id int, resource *T) (*T, error)
func Delete(ctx context.Context, service *zscaler.Service, id int) (*http.Response, error)
func GetAll(ctx context.Context, service *zscaler.Service, opts *GetAllFilterOptions) ([]T, error)
func GetByName(ctx context.Context, service *zscaler.Service, name string) (*T, error)
```

- `Create` returns `(interface{}, error)` from `Client.Create` — type-assert the result
- Endpoint: `fmt.Sprintf("%s/%d", endpoint, resourceID)` with int IDs

### ZPA

```go
func Get(ctx context.Context, service *zscaler.Service, id string) (*T, *http.Response, error)
func Create(ctx context.Context, service *zscaler.Service, resource T) (*T, *http.Response, error)
func Update(ctx context.Context, service *zscaler.Service, id string, resource *T) (*http.Response, error)
func Delete(ctx context.Context, service *zscaler.Service, id string) (*http.Response, error)
func GetAll(ctx context.Context, service *zscaler.Service) ([]T, *http.Response, error)
func GetByName(ctx context.Context, service *zscaler.Service, name string) (*T, *http.Response, error)
```

- **Every** call must pass `common.Filter{MicroTenantID: service.MicroTenantID()}`
- URL: `mgmtConfig + service.Client.GetCustomerID() + resourceEndpoint`
- Create takes struct by value, not pointer

### ZTW

Same as ZIA but client methods are **`Resource`-suffixed**: `ReadResource`, `CreateResource`, `UpdateWithPutResource`, `DeleteResource`. This is the most common mistake when implementing ZTW services.

### ZID

Same method names as ZIA (`Read/Create/UpdateWithPut/Delete`) but uses `string` IDs like ZPA. No activation. `GetByName` returns `[]T` with partial matching (`strings.Contains`), not a single item.

### ZCC

Uses `NewZccRequestDo` with manual `resp.Body.Close()` and `json.NewDecoder`. Some endpoints also support `service.Client.Read`.

### ZDX

Uses `NewRequestDo` like ZPA. Cursor-based pagination via `next_offset` token. Most endpoints are read-only.

## HTTP Method Discipline (do NOT call ExecuteRequest from a service)

**Hard rule:** code under `zscaler/<cloud>/services/**` must perform every HTTP call through a `service.Client.<Helper>` method. Direct calls to `service.Client.ExecuteRequest`, `http.MethodGet`, `http.MethodPost`, etc. inside a service file are forbidden. The helper layer is the contract; services must not bypass it.

If no existing helper fits the API contract, **add a new helper** to the appropriate request file (e.g. `zscaler/ziarequests.go` for ZIA-style endpoints), then call the new helper from the service. Do not inline `ExecuteRequest` in the service "just for this one endpoint".

### ZIA helper catalog (`zscaler/ziarequests.go`)

| Helper | When to use |
|---|---|
| `Client.Read(ctx, endpoint, &out)` | GET that returns JSON; decoded into `out` |
| `Client.Create(ctx, endpoint, struct)` | POST a struct, response decoded into the **same** struct type |
| `Client.UpdateWithPut(ctx, endpoint, struct)` | PUT a struct |
| `Client.Update(ctx, endpoint, struct)` | PATCH a struct (`application/merge-patch+json`) |
| `Client.Delete(ctx, endpoint)` | DELETE, no body |
| `Client.BulkDelete(ctx, endpoint, payload)` | POST a delete batch |
| `Client.CreateWithSlicePayload(ctx, endpoint, slice)` | POST a JSON array, returns raw bytes |
| `Client.UpdateWithSlicePayload(ctx, endpoint, slice)` | PUT a JSON array, returns raw bytes |
| `Client.CreateWithRawPayload(ctx, endpoint, payload string)` | POST a raw string body as `application/json` (e.g. PAC file content) |
| `Client.CreateWithNoContent(ctx, endpoint, struct)` | POST that returns `204 No Content` (e.g. async export trigger) |
| `Client.ReadRaw(ctx, endpoint, requestContentType)` | GET that returns a non-JSON body (CSV, binary). Pass `""` for ZIA OneAPI file downloads — the request defaults to `application/json`; using `text/csv` on a body-less GET can yield 415 from the gateway. |
| `Client.CreateWithRawPayloadAndContentType(ctx, endpoint, []byte, contentType)` | POST a raw body with a non-JSON Content-Type (e.g. `text/csv` upload) |
| `Client.CreateWithJSONResponse(ctx, endpoint, requestStruct, &responseStruct)` | POST where the **request and response Go types differ** (e.g. validate / preview / dry-run endpoints). |

### Reference patterns

```go
// ZIA — CSV export (GET, no body; response is CSV). Request must use
// application/json (ReadRaw with ""), not text/csv, or OneAPI returns 415.
func ExportFoo(ctx context.Context, service *zscaler.Service) ([]byte, error) {
    return service.Client.ReadRaw(ctx, fooExportEndpoint, "")
}

// ZIA — CSV import (POST body is CSV; Content-Type must be text/csv)
func ImportFoo(ctx context.Context, service *zscaler.Service, csv []byte) (*http.Response, error) {
    _, resp, err := service.Client.CreateWithRawPayloadAndContentType(ctx, fooImportEndpoint, csv, "text/csv")
    return resp, err
}

// ZIA — validate / dry-run (request and response types differ)
func ValidateFoo(ctx context.Context, service *zscaler.Service, body string) (*FooValidation, error) {
    var v FooValidation
    err := service.Client.CreateWithJSONResponse(ctx, fooValidateEndpoint, FooValidateRequest{Body: body}, &v)
    return &v, err
}
```

### Other clouds

- **ZPA** — use `service.Client.NewRequestDo(...)`; never inline `http.NewRequest`.
- **ZCC** — use `service.Client.NewZccRequestDo(...)` or `service.Client.Read(...)`; never inline `http.NewRequest`.
- **ZDX** — use `service.Client.NewRequestDo(...)`.
- **ZTW** — use `Resource`-suffixed helpers (`ReadResource`, `CreateResource`, `UpdateWithPutResource`, `DeleteResource`).
- **ZID** — same helper names as ZIA.

### Code review checklist

- ❌ `http.MethodGet` / `http.MethodPost` / `http.MethodPut` / `http.MethodDelete` literal in a service file → add or use a `Client.*` helper instead.
- ❌ `service.Client.ExecuteRequest(...)` in a service file → add or use a `Client.*` helper instead.
- ❌ `http.NewRequest(...)` in a service file → add or use a `Client.*` helper instead.
- ❌ Inline `json.Marshal` of a request body in a service file → the helper handles marshaling.
- ✅ The only file in `zscaler/<cloud>/...` that calls `ExecuteRequest` is the cloud's request file (`ziarequests.go`, etc.). The single carve-out is `zscaler/zia/services/sandbox/sandbox_submission` — it talks to a different host (`/zscsb`), needs `url.Values` query params, and computes `Content-Type` from the uploaded filename's extension. None of those concerns fit a generic helper signature, so inline `ExecuteRequest` is acceptable **there only**. Do not cite sandbox as precedent for new services; if a new endpoint needs a one-off behavior, add a new helper to `ziarequests.go` first.

## Pagination Engines

| Engine | Location | Page Detection | Used By |
|--------|----------|---------------|---------|
| `ReadAllPages` | `zia/services/common/common.go` | `len(items) < pageSize` | ZIA (configurable pageSize) |
| `ReadAllPages` | `ztw/services/common/common.go` | `len(items) < pageSize` | ZTW (fixed 1000) |
| `ReadAllPages[T]` | `zcc/services/common/common.go` | `len(items) < pageSize` | ZCC v1 (bare JSON array response) |
| `ReadAllPagesV2[T]` | `zcc/services/common/common.go` | `len(allResults) >= total` OR `count < limit` OR `count == 0` | ZCC v2 (`{items, total, offset, limit, count}` envelope) |
| `GetAllPagesGenericWithCustomFilters[T]` | `zpa/services/common/common.go` | `page <= totalPages` from envelope | ZPA |
| `ReadAllPagesWithPagination` | `zid/services/common/common.go` | `next_link == ""` or `len < limit` | ZID (offset/limit) |
| `ReadAllPagesWithCursor` | `zid/services/common/common.go` | Chase `next_link` until empty | ZID (cursor) |
| Per-domain functions | `zdx/services/*/` | `next_offset` token is empty | ZDX |
| `ReadAllPages` | `zwa/services/common/common.go` | `currentPageSize < pageSize` or `page >= totalPages` | ZWA (cursor with TotalPages) |

All pagination engines ultimately return `[]T` to the caller. Pagination metadata is consumed internally.

## JMESPath Client-Side Filtering

JMESPath is integrated into all pagination engines via `context.Context`. Every `GetAll`, `GetByName`, and list function that uses a pagination helper automatically applies JMESPath filtering when an expression is present in the context.

### Usage (Context-based — automatic)

```go
ctx := zscaler.ContextWithJMESPath(ctx, "[?osType=='Windows']")
devices, err := devices.GetAll(ctx, service, nil)
// devices is already filtered — no additional code needed
```

The expression flows through `context.Context` to the pagination engine, which applies it after aggregating all pages. **Zero changes to any service function or caller signature are required.**

### Integrated Pagination Engines

| Engine | Package | Context Check |
|--------|---------|---------------|
| `ReadAllPages` | `zia/services/common` | Yes |
| `GetAllPagesGenericWithCustomFilters` | `zpa/services/common` | Yes |
| `GetAllPagesGeneric` | `zpa/services/common` | Yes |
| `GetAllPagesGenericWithPostSearch` | `zpa/services/common` | Yes |
| `ReadAllPages` | `zcc/services/common` | Yes |
| `ReadAllPagesV2` | `zcc/services/common` | Yes |
| `ReadAllPages` | `ztw/services/common` | Yes |
| `ReadAllPagesWithPagination` | `zid/services/common` | Yes |
| `ReadAllPagesWithCursor` | `zid/services/common` | Yes |
| `ReadAllPages` | `zwa/services/common` | Yes |

### API Reference

**Context helpers** (`zscaler/jmespath.go`):
- `ContextWithJMESPath(ctx, expression) context.Context` — attach expression to context
- `JMESPathFromContext(ctx) string` — extract expression (returns "" if none)

**Generic filter** (`zscaler/jmespath.go`):
- `ApplyJMESPathFilter[T any](items []T, expression string) ([]T, error)` — typed filter, marshals through JSON, returns `[]T`
- `ApplyJMESPathFromContext[T any](ctx, items []T) ([]T, error)` — checks context and applies filter

**Untyped search** (`zscaler/jmespath.go`):
- `SearchJMESPath(data interface{}, expression string) (interface{}, error)` — works with any shape, returns `interface{}` (for projections that reshape data)

### Expression Examples

```go
// Filter: keep only active Windows devices
"[?osType=='Windows' && status=='active']"

// Projection: extract names only (returns []map, not []T)
"[*].{name: name, id: id}"

// Combined filter + field extraction
"[?enabled==`true`].name"

// Count
"length([?status=='inactive'])"
```

### Key Details

- Expressions use **camelCase JSON field names** from struct tags, not Go PascalCase field names
- When no expression is set in context, pagination returns results unchanged (zero overhead)
- `ApplyJMESPathFilter[T]` round-trips through JSON, so works only for filter expressions that preserve object shape
- `SearchJMESPath` returns `interface{}` and supports projections that reshape data
- ZDX has no centralized pagination helper — use `SearchJMESPath` or `ApplyJMESPathFilter` manually after getting results

## Common Types by Cloud

**ZIA** (`zscaler/zia/services/common`): `IDNameExtensions` (`{ID int, Name string, Extensions map}`), `IDName`, `IDNameExternalID`, `ZPAAppSegments`

**ZPA** (`zscaler/zpa/services/common`): `CommonIDName` (`{ID string, Name string}`), `CommonSummary`, `Filter` (includes `MicroTenantID`), `Pagination`

**ZTW** (`zscaler/ztw/services/common`): `CommonIDNameExternalID`, `CommonIDName` (`{ID int, Name string}`), `IDNameExtensions`, `ECVMs`, `ManagementNw`

**ZID** (`zscaler/zid/services/common`): `IDNameDisplayName`, `PaginationResponse[T]` (`{ResultsTotal, Records, NextLink}`), `PaginationQueryParams` (fluent builder)

Each cloud has its own `common` package — never import ZIA's common in a ZPA service or vice versa.

## JSON Tag & Struct Rules

1. **Always use exact JSON key** from the API as the tag
2. **Add `omitempty` to most fields** — prevents sending zero values
3. **Do NOT add `omitempty` to meaningful booleans** where `false` is an explicit value (`enabled`, `defaultRule`, `predefined`)
4. **Do NOT add `omitempty`** where zero value is valid (`"order": 0` means first position)
5. **ID types**: `int` for ZIA/ZCC/ZDX/ZTW, `string` for ZPA/ZID — never mix
6. **Read-only fields** (timestamps, modifiedBy) always use `omitempty`

## Client Initialization

```go
config, err := zscaler.NewConfiguration(
    zscaler.WithClientID("client-id"),
    zscaler.WithClientSecret("secret"),
    zscaler.WithVanityDomain("acme"),
    zscaler.WithCloud("zscloud"),
    zscaler.WithCache(true),
    zscaler.WithCacheTtl(10 * time.Minute),
)
service, err := zscaler.NewOneAPIClient(config)
```

### Service Structure

```go
type Service struct {
    Client        *Client
    LegacyClient  *LegacyClient
    microTenantID *string
    SortOrder     SortOrder
    SortBy        SortField
}
```

- `WithMicroTenant(id) *Service` — scope to a microtenant (ZPA only)
- `MicroTenantID() *string` — current microtenant (may be nil)
- `WithSort(sortBy, sortOrder) *Service` — custom sort

## Error Handling

The `errorx` package provides:
- `ErrorResponse` — wraps HTTP errors with parsed API details
- `IsObjectNotFound()` — `true` for 404 / `resource.not.found`
- `IsLimitExceeded()` — `true` for 403 tenant limit errors

Service functions should NOT catch/wrap these — let them propagate to the caller.

## Rate Limiting & Retries

Automatic, per-cloud:
- **Rate limiting** — e.g., ZIA: 20 GET/10s, 10 POST/10s
- **429 / 503 / 401** — `Retry-After` is honoured as the FLOOR for the first retry, then grown exponentially per consecutive retry of the same call (`base · 2^attempt`, capped at `RetryWaitMax`, default 10s) and perturbed by ±25% jitter so parallel goroutines do not stampede the per-endpoint limiter (v3.8.33+, `oneapiconfig.go` `Backoff` closure + `jitter` helper). The same policy is mirrored in `ExecuteRequest`'s outer 429 fallback.
- **`MaxNumOfRetries`** — default `10` (v3.8.33+, was `100`). Override via `cfg.Zscaler.Client.RateLimit.MaxRetries` or env `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES`. Lowered because a single stuck call no longer needs to monopolise a goroutine for ~100s; ten attempts across a jittered exponential window cover all observed real-world recoveries.
- **401 SESSION_NOT_VALID** — auto-refreshes OAuth2 token; bounded by `MaxSessionNotValidRetries` (default 3).
- **409 / 412 EDIT_LOCK_NOT_AVAILABLE / `Failed during enter Org barrier`** — exponential backoff in both the retryablehttp `CheckRetry` (`errorx.IsEditLockError`) and the `ExecuteRequest` outer loop.

Service implementations must NOT implement their own retry logic. The retry policy is locked in by `zscaler/oneapiconfig_retry_test.go` (`TestJitter`, `TestRetryBackoffPolicy`, `TestRetryMaxDefault`) — those tests fail loudly if a future change regresses the contract.

## Caching

When `WithCache(true)`:
- GET responses cached by URL key
- POST/PUT/DELETE auto-invalidate related entries via parent-collection invalidation (v3.8.32+ for the OneAPI client; previously only the legacy v2 clients did this, which manifested as stale `GetAll` snapshots in the Terraform provider's rule reorder loop — see SUP-3988).
- No action needed in service implementations

## Critical Gotchas

- **ZIA/ZTW require activation.** Changes are staged until `activation.UpdateActivationStatus` is called. The SDK does NOT auto-activate.
- **ZPA always needs MicroTenantID.** Every `NewRequestDo` call must include `common.Filter{MicroTenantID: service.MicroTenantID()}`, even if nil.
- **ZTW uses `Resource`-suffixed methods.** `ReadResource`, `CreateResource`, `UpdateWithPutResource`, `DeleteResource` — NOT `Read`, `Create`, etc. Most common ZTW mistake.
- **ZIA has a global edit lock.** Only one mutation at a time across the entire tenant. SDK retries automatically.
- **ZPA IDs are strings**, even though they look numeric. ZIA/ZTW IDs are `int`.
- **ZID `GetByName` returns `[]T`** with partial matching (`strings.Contains`), not a single item.
- **ZCC requires manual response handling.** `defer resp.Body.Close()`, check `resp.StatusCode`, decode with `json.NewDecoder`.
- **ZDX is mostly read-only.** Cursor-based pagination with `next_offset` — no centralized `ReadAllPages`.
- **Boolean omitempty.** Adding `omitempty` to `Enabled` means `false` is never sent — use `json:"enabled"` without `omitempty` for meaningful booleans.

## Directory Layout for New Services

```
zscaler/<cloud>/services/<service_name>/
├── <service_name>.go           # Structs + CRUD functions
└── <service_name>_test.go      # Integration tests

tests/unit/<cloud>/services/
└── <service_name>_test.go      # Unit tests (mock HTTP, no live API)
```

Package name: lowercase, no underscores. File name matches package name.

## Adding a New Service

1. **Determine the cloud** — ZIA, ZPA, ZCC, ZDX, ZTW, or ZID
2. **Get the JSON payload** — never guess field names; require the API JSON body
3. **Create package** under `zscaler/<cloud>/services/<service_name>/`
4. **Define structs** using JSON-to-Go mapping rules and cloud-specific common types
5. **Implement CRUD** using the correct cloud pattern (see CRUD Patterns section)
6. **Implement GetAll** using the correct pagination helper
7. **Implement GetByName** using `strings.EqualFold` (ZIA/ZPA/ZTW) or `strings.Contains` (ZID)
8. **Write integration tests** with `tests.NewOneAPIClient()` and `defer` cleanup
9. **Write unit tests** in `tests/unit/<cloud>/services/` with mock HTTP server

## Adding Fields to an Existing Service

1. Get the updated JSON payload with new fields
2. Add fields to the struct with correct JSON tags and `omitempty` rules
3. No CRUD changes needed — serialization is automatic
4. Update tests if new fields affect behavior

## Test Patterns

### Integration Tests

```go
service, err := tests.NewOneAPIClient()
ctx := context.Background()
resource := ResourceName{Name: fmt.Sprintf("tests-sdk-go-%d", time.Now().UnixNano())}
created, _, err := Create(ctx, service, &resource)
defer func() { Delete(ctx, service, created.ID) }()
```

### Unit Tests

Located in `tests/unit/<cloud>/services/`. Use `package unit`, mock server via `common.NewTestServer()`, register responses with `server.On(method, path, response)`.

Required coverage: Get, GetByName, Create, Update, Delete, GetAll, GetByName_NotFound, Get_NotFound.

## Development

- **Go version**: 1.24+
- **Module**: `github.com/zscaler/zscaler-sdk-go/v3`
- **Vendor mode**: `go mod vendor` after dependency changes
- **Lint**: `go vet ./...`
- **Test (unit)**: `go test ./tests/unit/... -v`
- **Test (integration)**: `go test ./zscaler/<cloud>/services/<service>/ -v` (requires env vars)
- **Format**: `gofmt -w .`

### Environment Variables

- `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET` — OneAPI credentials
- `ZSCALER_VANITY_DOMAIN` — Tenant vanity domain
- `ZSCALER_CLOUD` — Cloud instance (e.g., `zscloud`, `zscalerbeta`)
- `ZSCALER_CUSTOMER_ID` — Required for ZPA

## Release Versioning

All three must be updated in sync:

1. **`CHANGELOG.md`** — new entry at top with `[PR #NNN](https://github.com/zscaler/zscaler-sdk-go/pull/NNN)` links
2. **`docs/guides/release-notes.md`** — identical content, bump `Last updated: vX.Y.Z`
3. **`zscaler/oneapiclient.go`** — update `VERSION = "X.Y.Z"` constant

## Downstream Consumers

This SDK is consumed by:
- **Terraform Provider** (`zscaler/terraform-provider-zscaler`) — primary consumer, handles activation and error mapping
- **Zscaler MCP Server** (`zscaler/zscaler-mcp-server`) — MCP tool server wrapping SDK operations
- **Direct Go applications** — any Go code using the SDK library
