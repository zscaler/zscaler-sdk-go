// Package common — URL builders for the OneAPI mock server.
//
// All Zscaler clouds the OneAPI client speaks to follow a fixed prefix
// convention; these helpers eliminate string concatenation in tests and
// stop typos like "/zpa/mgmtConfig/" (capital C) from creeping in.
//
// Reference (from CLAUDE.md / SDK source):
//
//	ZPA  → /zpa/mgmtconfig/v{1,2}/admin/customers/{customerID}/{resource}[/{id}]
//	     special: /zpa/waap-pra-config/v1/admin/customers/{customerID}/credential-pool
//	ZIA  → /zia/api/v1/{resource}[/{id}]
//	ZCC  → /zcc/papi/public/v{1,2}/{resource}[/{id}]
//	ZDX  → /zdx/v1/{resource}[/{id}]
//	ZTW  → /ztw/api/v1/{resource}[/{id}]
//	ZID  → /admin/api/v1/{resource}[/{id}]
package common

import (
	"net/url"
	"strings"
)

// ─── ZPA ─────────────────────────────────────────────────────────────────────

// ZPAPath builds a v1 ZPA mgmtconfig path.
//
//	ZPAPath("123456789", "appConnectorGroup")          → /zpa/mgmtconfig/v1/admin/customers/123456789/appConnectorGroup
//	ZPAPath("123456789", "appConnectorGroup", "abc-1") → /zpa/mgmtconfig/v1/admin/customers/123456789/appConnectorGroup/abc-1
func ZPAPath(customerID, resource string, parts ...string) string {
	return zpaPathWithBase("/zpa/mgmtconfig/v1/admin/customers/", customerID, resource, parts...)
}

// ZPAv2Path builds a v2 ZPA mgmtconfig path (used by lssconfigcontroller,
// userPortalLink, etc.).
func ZPAv2Path(customerID, resource string, parts ...string) string {
	return zpaPathWithBase("/zpa/mgmtconfig/v2/admin/customers/", customerID, resource, parts...)
}

// ZPAUserConfigPath builds a /zpa/userconfig/v1/customers/... path
// (used by SCIM attribute-header lookups, etc.).
func ZPAUserConfigPath(customerID, resource string, parts ...string) string {
	return zpaPathWithBase("/zpa/userconfig/v1/customers/", customerID, resource, parts...)
}

// ZPAWaapPRAPath builds a /zpa/waap-pra-config/v1/admin/customers/... path
// (used by pracredentialpool only).
func ZPAWaapPRAPath(customerID, resource string, parts ...string) string {
	return zpaPathWithBase("/zpa/waap-pra-config/v1/admin/customers/", customerID, resource, parts...)
}

func zpaPathWithBase(base, customerID, resource string, parts ...string) string {
	resource = strings.TrimPrefix(resource, "/")
	all := append([]string{base + customerID, resource}, parts...)
	return joinClean(all...)
}

// ─── ZIA ─────────────────────────────────────────────────────────────────────

// ZIAPath builds a /zia/api/v1/... path.
//
//	ZIAPath("urlFilteringRules")          → /zia/api/v1/urlFilteringRules
//	ZIAPath("urlFilteringRules", "12345") → /zia/api/v1/urlFilteringRules/12345
func ZIAPath(parts ...string) string {
	return joinClean(append([]string{"/zia/api/v1"}, parts...)...)
}

// ─── ZCC ─────────────────────────────────────────────────────────────────────

// ZCCPath builds a /zcc/papi/public/v1/... path.
func ZCCPath(parts ...string) string {
	return joinClean(append([]string{"/zcc/papi/public/v1"}, parts...)...)
}

// ZCCv2Path builds a /zcc/papi/public/v2/... path.
func ZCCv2Path(parts ...string) string {
	return joinClean(append([]string{"/zcc/papi/public/v2"}, parts...)...)
}

// ─── ZDX ─────────────────────────────────────────────────────────────────────

// ZDXPath builds a /zdx/v1/... path. Note: ZDX does NOT include "/api"
// in its base prefix — it's "/zdx/v1/...", not "/zdx/api/v1/...".
func ZDXPath(parts ...string) string {
	return joinClean(append([]string{"/zdx/v1"}, parts...)...)
}

// ─── ZTW ─────────────────────────────────────────────────────────────────────

// ZTWPath builds a /ztw/api/v1/... path.
func ZTWPath(parts ...string) string {
	return joinClean(append([]string{"/ztw/api/v1"}, parts...)...)
}

// ─── ZID ─────────────────────────────────────────────────────────────────────

// ZIDPath builds a /admin/api/v1/... path (the ZIdentity admin API).
func ZIDPath(parts ...string) string {
	return joinClean(append([]string{"/admin/api/v1"}, parts...)...)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

// joinClean joins URL segments with single slashes regardless of whether
// individual parts have leading/trailing slashes.
func joinClean(parts ...string) string {
	cleaned := make([]string, 0, len(parts))
	for i, p := range parts {
		p = strings.Trim(p, "/")
		if p == "" {
			continue
		}
		if i == 0 {
			// preserve absolute leading slash
			cleaned = append(cleaned, "/"+p)
		} else {
			cleaned = append(cleaned, p)
		}
	}
	return strings.Join(cleaned, "/")
}

// WithQuery appends URL-encoded query params to a path.
//
//	WithQuery(ZPAPath(cid, "appConnectorGroup"), map[string]string{"search": "foo"})
//	  → /zpa/mgmtconfig/v1/admin/customers/.../appConnectorGroup?search=foo
func WithQuery(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}
	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	return path + "?" + q.Encode()
}
