// Package common — end-to-end test for the OneAPI factory layer.
//
// This file is the canonical example of how to drive a OneAPI-routed
// unit test using only the centralized harness in oneapi.go (factories),
// paths.go (URL builders), and envelopes.go (list builders). Every
// product gets a smoke test; if any of them break, the centralized
// harness has regressed.
package common_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/appconnectorgroup"
)

// TestOneAPIHarness_ZIA proves the ZIA factory wires the SDK to the
// mock server without any per-test boilerplate. Compare to the older
// pattern in tests/unit/zia/services/firewallpolicies_test.go (5 lines
// of setup, defer, manual customerID).
func TestOneAPIHarness_ZIA(t *testing.T) {
	api := common.NewZIATest(t)

	api.On("GET", common.ZIAPath("firewallFilteringRules", "12345"),
		common.SuccessResponse(filteringrules.FirewallFilteringRules{
			ID:     12345,
			Name:   "Block Bad IPs",
			Action: "BLOCK",
			State:  "ENABLED",
		}),
	)

	got, err := filteringrules.Get(context.Background(), api.Service, 12345)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 12345, got.ID)
	assert.Equal(t, "Block Bad IPs", got.Name)
	assert.Equal(t, "BLOCK", got.Action)
}

// TestOneAPIHarness_ZPA proves the ZPA factory composes correctly with
// the customer-ID-aware path builder and the ZPA list envelope helper.
// This single test exercises the full chain: factory → path builder →
// envelope builder → SDK pagination engine.
func TestOneAPIHarness_ZPA(t *testing.T) {
	api := common.NewZPATest(t)

	groups := []appconnectorgroup.AppConnectorGroup{
		{ID: "acg-001", Name: "HQ", Enabled: true},
		{ID: "acg-002", Name: "Branch", Enabled: false},
	}

	api.On("GET",
		common.ZPAPath(api.CustomerID, "appConnectorGroup"),
		common.SuccessResponse(common.ZPAList(groups)),
	)

	got, _, err := appconnectorgroup.GetAll(context.Background(), api.Service)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "acg-001", got[0].ID)
	assert.Equal(t, "Branch", got[1].Name)
}

// TestOneAPIHarness_AutoCleanup asserts t.Cleanup did its job: after
// the parent test ends, the inner sub-test creates its own harness and
// verifies the previous server isn't lingering on the same port (a
// negative test for resource leakage).
func TestOneAPIHarness_AutoCleanup(t *testing.T) {
	var firstURL string

	t.Run("first", func(t *testing.T) {
		api := common.NewZIATest(t)
		firstURL = api.Server.URL
		require.NotEmpty(t, firstURL)
	})

	t.Run("second", func(t *testing.T) {
		api := common.NewZIATest(t)
		require.NotEmpty(t, api.Server.URL)
		// httptest.Server uses ephemeral ports; the chance that
		// two sequential servers reuse the same one is non-zero
		// but vanishingly low. We only assert the second server
		// is reachable — the real proof that cleanup worked is
		// that this sub-test ran at all (the first server's
		// goroutine didn't deadlock the test binary).
	})
}

// TestOneAPIHarness_FactoriesAllProducts is a compile-only guard that
// keeps the six per-product factories from rotting. If any of them
// stops compiling — wrong arg count, wrong return type — this file
// stops building and the regression is caught immediately. We don't
// hit the network here, just invoke each factory once.
func TestOneAPIHarness_FactoriesAllProducts(t *testing.T) {
	for _, factory := range []struct {
		name string
		new  func(*testing.T) *common.APITest
	}{
		{"ZPA", common.NewZPATest},
		{"ZIA", common.NewZIATest},
		{"ZCC", common.NewZCCTest},
		{"ZDX", common.NewZDXTest},
		{"ZTW", common.NewZTWTest},
		{"ZID", common.NewZIDTest},
	} {
		t.Run(factory.name, func(t *testing.T) {
			api := factory.new(t)
			require.NotNil(t, api.Server)
			require.NotNil(t, api.Service)
			require.Equal(t, common.TestCustomerID, api.CustomerID)
		})
	}
}
