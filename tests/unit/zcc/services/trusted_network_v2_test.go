// Package services provides unit tests for ZCC services
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	commontests "github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/common"
	tn "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zcc/services/trusted_network_v2"
)

// =====================================================
// SDK Function Tests — exercise actual SDK code paths
// against a mock HTTP server.
// =====================================================

func TestTrustedNetworkV2_Get_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/trusted-networks/{id}"
	server.On("GET", itemPath, commontests.SuccessResponse(tn.TrustedNetworkV2{
		ID:                42,
		CompanyID:         12345,
		Name:              "Corporate HQ",
		ConditionType:     "AND",
		Active:            true,
		DNSServerIPs:      []string{"8.8.8.8", "1.1.1.1"},
		TrustedSubnetIPs:  []string{"10.0.0.0/8"},
		TrustedGatewayIPs: []string{"10.0.0.1"},
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	got, err := tn.Get(context.Background(), service, 42)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 42, got.ID)
	assert.Equal(t, "Corporate HQ", got.Name)
	assert.Equal(t, "AND", got.ConditionType)
	assert.True(t, got.Active)
	assert.Equal(t, []string{"8.8.8.8", "1.1.1.1"}, got.DNSServerIPs)
}

func TestTrustedNetworkV2_Create_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/trusted-networks"
	server.On("POST", listPath, commontests.SuccessResponse(tn.TrustedNetworkV2{
		ID:               99,
		CompanyID:        12345,
		Name:             "New Network",
		ConditionType:    "AND",
		Active:           true,
		TrustedSubnetIPs: []string{"192.168.0.0/16"},
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	created, _, err := tn.Create(context.Background(), service, &tn.TrustedNetworkV2{
		Name:             "New Network",
		ConditionType:    "AND",
		Active:           true,
		TrustedSubnetIPs: []string{"192.168.0.0/16"},
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, 99, created.ID)
	assert.Equal(t, "New Network", created.Name)
}

func TestTrustedNetworkV2_Create_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = tn.Create(context.Background(), service, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "trusted network is required")
}

func TestTrustedNetworkV2_Update_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/trusted-networks/{id}"
	server.On("PUT", itemPath, commontests.SuccessResponse(tn.TrustedNetworkV2{
		ID:               42,
		Name:             "Updated Network",
		ConditionType:    "OR",
		Active:           true,
		TrustedSubnetIPs: []string{"172.16.0.0/12"},
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updated, _, err := tn.Update(context.Background(), service, 42, &tn.TrustedNetworkV2{
		Name:             "Updated Network",
		ConditionType:    "OR",
		Active:           true,
		TrustedSubnetIPs: []string{"172.16.0.0/12"},
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, 42, updated.ID)
	assert.Equal(t, "Updated Network", updated.Name)
	assert.Equal(t, "OR", updated.ConditionType)
}

func TestTrustedNetworkV2_PartialUpdate_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/trusted-networks/{id}"
	server.On("PATCH", itemPath, commontests.SuccessResponse(tn.TrustedNetworkV2{
		ID:            42,
		Name:          "Renamed",
		ConditionType: "AND",
		Active:        true,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	patched, _, err := tn.PartialUpdate(context.Background(), service, 42, &tn.TrustedNetworkV2{
		Name: "Renamed",
	})
	require.NoError(t, err)
	require.NotNil(t, patched)
	assert.Equal(t, 42, patched.ID)
	assert.Equal(t, "Renamed", patched.Name)
}

func TestTrustedNetworkV2_PartialUpdate_NilBody_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, _, err = tn.PartialUpdate(context.Background(), service, 42, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "trusted network is required")
}

func TestTrustedNetworkV2_Delete_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	itemPath := "/zcc/papi/public/v2/trusted-networks/{id}"
	server.On("DELETE", itemPath, commontests.NoContentResponse())

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = tn.Delete(context.Background(), service, 42)
	require.NoError(t, err)
}

// GetAll on a single-page response: server returns the entire dataset on
// page one and the pagination loop terminates after one iteration. This is
// the most common case in practice.
func TestTrustedNetworkV2_GetAll_SinglePage_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/trusted-networks"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items: []tn.TrustedNetworkV2{
			{ID: 1, Name: "alpha"},
			{ID: 2, Name: "beta"},
			{ID: 3, Name: "gamma"},
		},
		Total:  3,
		Offset: 0,
		Limit:  50,
		Count:  3,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	all, err := tn.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	require.Len(t, all, 3)
	assert.Equal(t, "alpha", all[0].Name)
	assert.Equal(t, "gamma", all[2].Name)
}

// GetAll must send `skip` / `perPage` (the documented v2 names) — not the
// older `page` / `pageSize`. This is the regression guard for the bug
// where the SDK was sending unknown params, the server ignored them, and
// callers silently looped over the default first page.
func TestTrustedNetworkV2_GetAll_UsesDocumentedQueryParams_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/trusted-networks"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items: []tn.TrustedNetworkV2{{ID: 1, Name: "alpha"}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = tn.GetAll(context.Background(), service, &tn.GetAllFilterOptions{
		Keyword: "alpha",
		Type:    tn.FilterTypeName,
	})
	require.NoError(t, err)

	last := server.LastRequest()
	require.NotNil(t, last)
	q := last.Query

	assert.Contains(t, q, "perPage=50", "v2 endpoints take perPage, not pageSize")
	assert.NotContains(t, q, "pageSize=", "pageSize is the v1 param; v2 must use perPage")
	assert.NotContains(t, q, "page=", "v2 endpoints take skip, not page")

	assert.Contains(t, q, "keyword=alpha")
	assert.Contains(t, q, "type=NAME")
	// skip=0 is omitted by omitempty — that's intentional and matches the
	// documented default. Just ensure no skip param leaks with a wrong value.
	assert.NotContains(t, q, "skip=50")
}

// GetAll over two pages exercises the skip += perPage iteration. The
// canned httptest handler returns a full page first (count == limit) and a
// short page second (count < limit) so the loop terminates on the
// `last-page heuristic` branch.
func TestTrustedNetworkV2_GetAll_MultiPage_SDK(t *testing.T) {
	var calls int32
	queries := make([]string, 0, 2)

	first := common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items:  buildItems(1, 50),
		Total:  52,
		Offset: 0,
		Limit:  50,
		Count:  50,
	}
	second := common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items:  buildItems(51, 2),
		Total:  52,
		Offset: 50,
		Limit:  50,
		Count:  2,
	}

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if atomic.AddInt32(&calls, 1) == 1 {
			_ = json.NewEncoder(w).Encode(first)
			return
		}
		_ = json.NewEncoder(w).Encode(second)
	}))
	defer upstream.Close()

	service, err := commontests.CreateTestService(context.Background(),
		&commontests.TestServer{Server: upstream, Handler: commontests.NewMockHandler()},
		"123456")
	require.NoError(t, err)

	all, err := tn.GetAll(context.Background(), service, nil)
	require.NoError(t, err)
	assert.Len(t, all, 52)

	require.Len(t, queries, 2)
	assert.Contains(t, queries[0], "perPage=50")
	assert.NotContains(t, queries[0], "skip=")
	assert.Contains(t, queries[1], "skip=50")
	assert.Contains(t, queries[1], "perPage=50")
}

func TestTrustedNetworkV2_GetByName_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/trusted-networks"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items: []tn.TrustedNetworkV2{
			{ID: 7, Name: "Corporate HQ"},
			{ID: 8, Name: "Branch Office"},
		},
		Total: 2, Offset: 0, Limit: 50, Count: 2,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// Case-insensitive exact match.
	got, err := tn.GetByName(context.Background(), service, "corporate hq")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 7, got.ID)
	assert.Equal(t, "Corporate HQ", got.Name)

	// Confirm the server-side narrow was applied: keyword + type=NAME.
	last := server.LastRequest()
	require.NotNil(t, last)
	assert.Contains(t, last.Query, "keyword=corporate")
	assert.Contains(t, last.Query, "type=NAME")
}

func TestTrustedNetworkV2_GetByName_NotFound_SDK(t *testing.T) {
	server := commontests.NewTestServer()
	defer server.Close()

	listPath := "/zcc/papi/public/v2/trusted-networks"
	server.On("GET", listPath, commontests.SuccessResponse(common.PaginatedResponseV2[tn.TrustedNetworkV2]{
		Items: []tn.TrustedNetworkV2{{ID: 1, Name: "alpha"}},
		Total: 1, Offset: 0, Limit: 50, Count: 1,
	}))

	service, err := commontests.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = tn.GetByName(context.Background(), service, "does-not-exist")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no trusted network found with name: does-not-exist")
}

// =====================================================
// Structure tests — JSON marshaling / unmarshaling
// =====================================================

func TestTrustedNetworkV2_Structure(t *testing.T) {
	t.Parallel()

	t.Run("TrustedNetworkV2 marshals expected JSON", func(t *testing.T) {
		network := tn.TrustedNetworkV2{
			ID:                42,
			Name:              "Corporate",
			ConditionType:     "AND",
			Active:            true,
			DNSServerIPs:      []string{"8.8.8.8"},
			TrustedSubnetIPs:  []string{"10.0.0.0/8"},
			TrustedGatewayIPs: []string{"10.0.0.1"},
		}

		data, err := json.Marshal(network)
		require.NoError(t, err)
		body := string(data)

		assert.Contains(t, body, `"id":42`)
		assert.Contains(t, body, `"name":"Corporate"`)
		assert.Contains(t, body, `"conditionType":"AND"`)
		assert.Contains(t, body, `"active":true`)
		assert.Contains(t, body, `"dnsServerIps":["8.8.8.8"]`)
	})

	t.Run("TrustedNetworkV2 unmarshals envelope items", func(t *testing.T) {
		jsonData := `{
			"items": [
				{
					"id": 99,
					"companyId": 12345,
					"name": "Branch",
					"conditionType": "OR",
					"active": true,
					"trustedSubnetIps": ["10.1.0.0/16"]
				}
			],
			"total": 1,
			"offset": 0,
			"limit": 50,
			"count": 1
		}`

		var page common.PaginatedResponseV2[tn.TrustedNetworkV2]
		err := json.Unmarshal([]byte(jsonData), &page)
		require.NoError(t, err)
		require.Len(t, page.Items, 1)
		assert.Equal(t, 99, page.Items[0].ID)
		assert.Equal(t, "Branch", page.Items[0].Name)
		assert.Equal(t, []string{"10.1.0.0/16"}, page.Items[0].TrustedSubnetIPs)
		assert.Equal(t, 1, page.Count)
		assert.Equal(t, 50, page.Limit)
	})

	t.Run("FilterType constants cover documented enum values", func(t *testing.T) {
		expected := map[string]bool{
			"NAME":                 true,
			"DNS_SERVERS":          true,
			"DNS_SEARCH_DOMAINS":   true,
			"HOST_NAME_IP":         true,
			"TRUSTED_SUBNETS":      true,
			"TRUSTED_GATEWAYS":     true,
			"TRUSTED_DHCP_SERVERS": true,
			"TRUSTED_EGRESS_IPS":   true,
			"SSID":                 true,
		}
		got := map[string]bool{
			tn.FilterTypeName:               true,
			tn.FilterTypeDNSServers:         true,
			tn.FilterTypeDNSSearchDomains:   true,
			tn.FilterTypeHostnameIP:         true,
			tn.FilterTypeTrustedSubnets:     true,
			tn.FilterTypeTrustedGateways:    true,
			tn.FilterTypeTrustedDHCPServers: true,
			tn.FilterTypeTrustedEgressIPs:   true,
			tn.FilterTypeSSID:               true,
		}
		assert.Equal(t, expected, got)
	})
}

// buildItems returns a slice of n TrustedNetworkV2 records with sequential
// IDs starting at startID, used to fabricate paginated test responses.
func buildItems(startID, n int) []tn.TrustedNetworkV2 {
	items := make([]tn.TrustedNetworkV2, n)
	for i := 0; i < n; i++ {
		id := startID + i
		items[i] = tn.TrustedNetworkV2{ID: id, Name: fmt.Sprintf("tn-%d", id)}
	}
	return items
}
