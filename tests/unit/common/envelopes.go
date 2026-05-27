// Package common — list / pagination envelope builders for the OneAPI mock server.
//
// Each cloud has a different list-response shape. Tests should never hand-write
// these envelopes; use the helpers below so the shape stays in sync with the
// real pagination engines under zscaler/<cloud>/services/common/.
//
// Envelope reference (from zscaler/<cloud>/services/common/common.go):
//
//	ZPA  → {"list": [...], "totalPages": <int|string>}
//	ZIA  → bare JSON array (or for some endpoints, {"list": [...]})
//	ZCC v1 → bare JSON array
//	ZCC v2 → {"items": [...], "total": N, "offset": 0, "limit": N, "count": N}
//	ZDX  → varies by endpoint; cursor-based with "next_offset"
//	ZTW  → bare JSON array (with fixed pageSize=1000)
//	ZID  → {"records": [...], "results_total": N, "next_link": ""}
package common

// ─── ZPA ─────────────────────────────────────────────────────────────────────

// ZPAList wraps items in the ZPA single-page envelope. Use for tests where
// the SDK consumer doesn't iterate pages (totalPages defaults to 1).
//
// Example:
//
//	server.On("GET", path, common.SuccessResponse(common.ZPAList(groups)))
func ZPAList[T any](items []T) map[string]any {
	return ZPAListPaged(items, 1)
}

// ZPAListPaged wraps items in the ZPA pagination envelope with explicit
// totalPages. Use with OnSequence to model multi-page responses.
//
// Example:
//
//	server.OnSequence("GET", path,
//	    common.SuccessResponse(common.ZPAListPaged(page1, 2)),
//	    common.SuccessResponse(common.ZPAListPaged(page2, 2)),
//	)
func ZPAListPaged[T any](items []T, totalPages int) map[string]any {
	if items == nil {
		items = []T{}
	}
	return map[string]any{
		"list":       items,
		"totalPages": totalPages,
	}
}

// ─── ZIA ─────────────────────────────────────────────────────────────────────

// ZIAList is the identity helper for ZIA list responses (bare JSON arrays).
// Provided for symmetry with the other clouds and to make test intent
// explicit ("this is a ZIA list response").
//
// Example:
//
//	server.On("GET", path, common.SuccessResponse(common.ZIAList(rules)))
func ZIAList[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}

// ─── ZCC v1 ──────────────────────────────────────────────────────────────────

// ZCCList is the identity helper for ZCC v1 list responses (bare JSON arrays).
func ZCCList[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}

// ─── ZCC v2 ──────────────────────────────────────────────────────────────────

// ZCCv2List wraps items in the ZCC v2 pagination envelope with sensible
// defaults — total = len(items), offset = 0, limit = len(items), count = len(items).
// For multi-page tests use ZCCv2ListPaged.
func ZCCv2List[T any](items []T) map[string]any {
	n := len(items)
	return ZCCv2ListPaged(items, n, 0, n)
}

// ZCCv2ListPaged wraps items in the explicit ZCC v2 envelope.
//
//	total  — total records across all pages
//	offset — zero-based starting index of this page
//	limit  — max records per page
//	count  — actual number of records on this page (auto-derived)
func ZCCv2ListPaged[T any](items []T, total, offset, limit int) map[string]any {
	if items == nil {
		items = []T{}
	}
	return map[string]any{
		"items":  items,
		"total":  total,
		"offset": offset,
		"limit":  limit,
		"count":  len(items),
	}
}

// ─── ZDX ─────────────────────────────────────────────────────────────────────

// ZDXCursorList wraps items in the ZDX cursor-pagination envelope used by
// most reporting endpoints. Pass an empty nextOffset on the last page so
// the SDK loop terminates.
//
// Example:
//
//	server.OnSequence("GET", path,
//	    common.SuccessResponse(common.ZDXCursorList(page1, "abc123")),
//	    common.SuccessResponse(common.ZDXCursorList(page2, "")),
//	)
func ZDXCursorList[T any](items []T, nextOffset string) map[string]any {
	if items == nil {
		items = []T{}
	}
	out := map[string]any{
		"items": items,
	}
	if nextOffset != "" {
		out["next_offset"] = nextOffset
	}
	return out
}

// ─── ZTW ─────────────────────────────────────────────────────────────────────

// ZTWList is the identity helper for ZTW list responses (bare JSON arrays,
// fixed pageSize=1000 in the real engine).
func ZTWList[T any](items []T) []T {
	if items == nil {
		return []T{}
	}
	return items
}

// ─── ZID ─────────────────────────────────────────────────────────────────────

// ZIDList wraps items in the ZID pagination envelope (single page —
// next_link defaults to "" so the loop terminates immediately).
func ZIDList[T any](items []T) map[string]any {
	return ZIDListPaged(items, "")
}

// ZIDListPaged wraps items in the ZID envelope with an explicit next_link.
// Pass nextLink="" on the final page so the SDK pagination loop stops.
func ZIDListPaged[T any](items []T, nextLink string) map[string]any {
	if items == nil {
		items = []T{}
	}
	return map[string]any{
		"records":       items,
		"results_total": len(items),
		"next_link":     nextLink,
	}
}

// ─── Common error bodies ─────────────────────────────────────────────────────

// ZPANotFoundBody returns the JSON body ZPA emits for resource.not.found.
func ZPANotFoundBody() string {
	return `{"id": "resource.not.found", "message": "Resource not found"}`
}

// ZIANotFoundBody returns the JSON body ZIA emits for missing resources.
func ZIANotFoundBody() string {
	return `{"code": "RESOURCE_NOT_FOUND", "message": "Resource not found"}`
}
