// Package common — centralized OneAPI mock harness (factory layer).
//
// This file is the ergonomic top of the unit-test infrastructure. It
// composes the existing primitives in this package into a single
// per-test entry point so individual product tests can be written with
// zero boilerplate.
//
// File layout (all in package common):
//
//	mocks.go      → MockHandler / TestServer / response builders
//	testutils.go  → MockTransport + CreateTestService (raw factory)
//	paths.go      → URL builders per cloud (ZPAPath, ZIAPath, …)
//	envelopes.go  → typed list/pagination envelopes per cloud
//	oneapi.go     → THIS FILE: APITest struct + per-product NewXxxTest(t)
//
// Why factories?
// Every test file used to repeat the same five-line ritual:
//
//	server := common.NewTestServer()
//	defer server.Close()
//	service, err := common.CreateTestService(ctx, server, "123456")
//	require.NoError(t, err)
//
// …and ZPA tests additionally redeclared `const testCustomerID = "..."`
// in every file. The factories below collapse all of that into one
// call. Pick the factory matching your product:
//
//	api := common.NewZPATest(t)
//	api := common.NewZIATest(t)
//	api := common.NewZCCTest(t)
//	api := common.NewZDXTest(t)
//	api := common.NewZTWTest(t)
//	api := common.NewZIDTest(t)
//
// The mock server is closed automatically via t.Cleanup, so `defer
// server.Close()` is gone too.
//
// Quick-start
//
//	func TestFirewallFilteringRules_Get_SDK(t *testing.T) {
//	    api := common.NewZIATest(t)
//	    path := common.ZIAPath("firewallFilteringRules", "12345")
//
//	    api.On("GET", path, common.SuccessResponse(filteringrules.FirewallFilteringRules{
//	        ID: 12345, Name: "Block Bad IPs", Action: "BLOCK",
//	    }))
//
//	    got, err := filteringrules.Get(context.Background(), api.Service, 12345)
//	    require.NoError(t, err)
//	    assert.Equal(t, 12345, got.ID)
//	}
//
// What it does NOT do
//
//   - Path building → use ZPAPath / ZIAPath / ZCCPath / ZCCv2Path /
//     ZDXPath / ZTWPath / ZIDPath from paths.go.
//   - List envelopes → use ZPAList / ZIAList / ZCCList / ZCCv2List /
//     ZIDList / ZDXCursorList from envelopes.go.
//
// Existing tests that call NewTestServer() + CreateTestService() directly
// continue to work unchanged; this layer is purely additive.
package common

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
)

// =============================================================================
// Shared Constants
// =============================================================================

// TestCustomerID is the canonical fake ZPA customer ID used across the
// entire unit-test suite. ZPA tests should reference this constant
// instead of redeclaring `const testCustomerID = "..."` per file —
// keeping one source of truth means a future change (say, switching the
// ID format) is a one-line edit.
const TestCustomerID = "123456789"

// =============================================================================
// APITest — bundled (server, service) for one test
// =============================================================================

// APITest groups the mock server and a fully-configured OneAPI
// *zscaler.Service for a single unit test. Always construct via one of
// the per-product factories (NewZPATest, NewZIATest, …) — the server
// is auto-closed via t.Cleanup so callers never need
// `defer server.Close()`.
type APITest struct {
	// Server is the in-memory httptest.Server backing the mock
	// handler. Use it directly for advanced cases (Reset between
	// sub-tests, AssertRequestCount, LastRequest, …).
	Server *TestServer

	// Service is the OneAPI client all SDK functions take as their
	// second positional argument. The HTTP transport is already wired
	// to Server, the OAuth token is pre-populated (no real auth call
	// is made), and the cache is disabled.
	Service *zscaler.Service

	// CustomerID is the ZPA customer ID baked into the configured
	// service. For non-ZPA tests it's still set (the SDK keeps it on
	// Configuration regardless of product) but you'll only need to
	// reference it when building ZPA paths — and paths.go's ZPAPath
	// helper accepts it as its first argument.
	CustomerID string
}

// On is a shortcut for Server.On — saves a few keystrokes and keeps
// test bodies terse.
func (a *APITest) On(method, path string, response MockResponse) *APITest {
	a.Server.On(method, path, response)
	return a
}

// OnSequence is a shortcut for Server.OnSequence.
func (a *APITest) OnSequence(method, path string, responses ...MockResponse) *APITest {
	a.Server.OnSequence(method, path, responses...)
	return a
}

// OnFunc is a shortcut for Server.OnFunc.
func (a *APITest) OnFunc(method, path string, fn MockResponseFunc) *APITest {
	a.Server.OnFunc(method, path, fn)
	return a
}

// LastRequest returns the most-recent HTTP request received by the
// server, or nil. Convenience pass-through.
func (a *APITest) LastRequest() *RecordedRequest {
	return a.Server.LastRequest()
}

// newAPITest is the shared backbone for the per-product factories. It
// spins up a fresh server, builds a OneAPI client wired to that server,
// and registers Close via t.Cleanup so the caller doesn't have to.
//
// OneAPI-only: legacy client is intentionally not configured. The
// underlying CreateTestService nulls out cfg.UseLegacyClient and
// cfg.LegacyClient before constructing the OneAPI client, so this
// harness is safe to use even if the developer's shell has
// `ZSCALER_USE_LEGACY_CLIENT=true` or a stray ~/.zscaler/zscaler.yaml
// is on disk.
func newAPITest(t *testing.T, customerID string) *APITest {
	t.Helper()
	server := NewTestServer()
	t.Cleanup(server.Close)

	service, err := CreateTestService(context.Background(), server, customerID)
	require.NoError(t, err, "failed to create OneAPI test service")

	return &APITest{
		Server:     server,
		Service:    service,
		CustomerID: customerID,
	}
}

// NewZPATest builds an APITest pre-configured for ZPA — the customer
// ID (TestCustomerID) is baked into the service so paths.go's ZPAPath
// helper produces matching URLs out of the box.
func NewZPATest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}

// NewZIATest builds an APITest pre-configured for ZIA. ZIA URLs don't
// include a customer ID; the value on the service is harmless.
func NewZIATest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}

// NewZCCTest builds an APITest pre-configured for ZCC (works for v1
// and v2 endpoints — pick the right path builder via ZCCPath /
// ZCCv2Path).
func NewZCCTest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}

// NewZDXTest builds an APITest pre-configured for ZDX.
func NewZDXTest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}

// NewZTWTest builds an APITest pre-configured for ZTW (Cloud / Branch
// Connector).
func NewZTWTest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}

// NewZIDTest builds an APITest pre-configured for ZID (Zidentity).
func NewZIDTest(t *testing.T) *APITest {
	t.Helper()
	return newAPITest(t, TestCustomerID)
}
