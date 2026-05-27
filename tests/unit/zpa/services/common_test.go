// Package unit provides unit tests for ZPA services common utilities
package unit

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa"
	zpacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/tag_controller/tag_group"
)

func TestCommon_RemoveCloudSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes cloud suffix with parentheses",
			input:    "CrowdStrike_ZPA_Pre-ZTA (zscalerthree.net)",
			expected: "CrowdStrike_ZPA_Pre-ZTA",
		},
		{
			name:     "removes cloud suffix with different domain",
			input:    "My App (zscaler.net)",
			expected: "My App",
		},
		{
			name:     "removes cloud suffix with underscores and hyphens",
			input:    "Test_App-Name (test-cloud_123.net)",
			expected: "Test_App-Name",
		},
		{
			name:     "no change when no suffix",
			input:    "Simple Name",
			expected: "Simple Name",
		},
		{
			name:     "handles empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "handles string with only spaces",
			input:    "   ",
			expected: "",
		},
		{
			name:     "removes trailing spaces after suffix removal",
			input:    "App Name   (cloud.net)  ",
			expected: "App Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zpacommon.RemoveCloudSuffix(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommon_InList(t *testing.T) {
	tests := []struct {
		name     string
		list     []string
		item     string
		expected bool
	}{
		{
			name:     "item exists in list",
			list:     []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "item not in list",
			list:     []string{"apple", "banana", "cherry"},
			item:     "grape",
			expected: false,
		},
		{
			name:     "empty list",
			list:     []string{},
			item:     "apple",
			expected: false,
		},
		{
			name:     "nil list",
			list:     nil,
			item:     "apple",
			expected: false,
		},
		{
			name:     "empty string item exists",
			list:     []string{"", "test"},
			item:     "",
			expected: true,
		},
		{
			name:     "case sensitive match",
			list:     []string{"Apple", "Banana"},
			item:     "apple",
			expected: false,
		},
		{
			name:     "exact case match",
			list:     []string{"Apple", "Banana"},
			item:     "Apple",
			expected: true,
		},
		{
			name:     "single item list - match",
			list:     []string{"only"},
			item:     "only",
			expected: true,
		},
		{
			name:     "single item list - no match",
			list:     []string{"only"},
			item:     "other",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zpacommon.InList(tt.list, tt.item)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCommon_Filter_Structure(t *testing.T) {
	t.Run("Filter with all fields", func(t *testing.T) {
		tenantID := "tenant-123"
		filter := zpacommon.Filter{
			Search:        "test-search",
			MicroTenantID: &tenantID,
			SortBy:        "name",
			SortOrder:     "ASC",
		}

		assert.Equal(t, "test-search", filter.Search)
		assert.Equal(t, "tenant-123", *filter.MicroTenantID)
		assert.Equal(t, "name", filter.SortBy)
		assert.Equal(t, "ASC", filter.SortOrder)
	})

	t.Run("Empty filter", func(t *testing.T) {
		filter := zpacommon.Filter{}

		assert.Empty(t, filter.Search)
		assert.Nil(t, filter.MicroTenantID)
		assert.Empty(t, filter.SortBy)
		assert.Empty(t, filter.SortOrder)
	})
}

func TestCommon_Pagination_Structure(t *testing.T) {
	t.Run("Pagination with all fields", func(t *testing.T) {
		tenantID := "tenant-456"
		pagination := zpacommon.Pagination{
			Page:          1,
			PageSize:      100,
			Search:        "test",
			Search2:       "search2",
			MicroTenantID: &tenantID,
			SortBy:        "creationTime",
			SortOrder:     "DESC",
		}

		assert.Equal(t, 1, pagination.Page)
		assert.Equal(t, 100, pagination.PageSize)
		assert.Equal(t, "test", pagination.Search)
		assert.Equal(t, "search2", pagination.Search2)
		assert.Equal(t, "tenant-456", *pagination.MicroTenantID)
		assert.Equal(t, "creationTime", pagination.SortBy)
		assert.Equal(t, "DESC", pagination.SortOrder)
	})
}

type zpaListItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestZPACommon_GetAllPagesGeneric(t *testing.T) {
	t.Parallel()

	t.Run("single page with search filter conversion", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			q := r.URL.Query()
			assert.Equal(t, "1", q.Get("page"))
			assert.Equal(t, "500", q.Get("pagesize"))
			search := q.Get("search")
			assert.Contains(t, search, "EQ")
			return common.SuccessResponse(common.ZPAList([]zpaListItem{
				{ID: "app-1", Name: "My App"},
			}))
		})

		got, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](context.Background(), api.Service.Client, path, "My App")
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("single page empty search", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			q := r.URL.Query()
			assert.Equal(t, "1", q.Get("page"))
			assert.Equal(t, "500", q.Get("pagesize"))
			return common.SuccessResponse(common.ZPAListPaged([]zpaListItem{
				{ID: "only", Name: "Only"},
			}, 1))
		})

		got, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](context.Background(), api.Service.Client, path, "")
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "only", got[0].ID)
	})

	t.Run("multi-page aggregation", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			switch r.URL.Query().Get("page") {
			case "1":
				assert.Equal(t, "500", r.URL.Query().Get("pagesize"))
				return common.SuccessResponse(common.ZPAListPaged([]zpaListItem{
					{ID: "app-1", Name: "One"},
				}, 2))
			case "2":
				return common.SuccessResponse(common.ZPAListPaged([]zpaListItem{
					{ID: "app-2", Name: "Two"},
				}, 2))
			default:
				t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
				return common.NotFoundResponse()
			}
		})

		got, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](context.Background(), api.Service.Client, path, "")
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("JMESPath filter from context", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.On("GET", path, common.SuccessResponse(common.ZPAList([]zpaListItem{
			{ID: "1", Name: "active"},
			{ID: "2", Name: "inactive"},
		})))

		ctx := zscaler.ContextWithJMESPath(context.Background(), "[?name=='active']")
		got, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](ctx, api.Service.Client, path, "")
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("read error propagates", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")
		api.On("GET", path, common.NotFoundResponse())

		_, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](context.Background(), api.Service.Client, path, "")
		require.Error(t, err)
	})

	t.Run("scimgroup URL uses plain search not name+EQ filter", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAUserConfigPath(api.CustomerID, "scimgroup")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			q := r.URL.Query()
			assert.Equal(t, "Alice Group", q.Get("search"))
			assert.NotContains(t, q.Get("search"), "EQ")
			return common.SuccessResponse(common.ZPAList([]zpaListItem{{ID: "sg-1", Name: "Alice Group"}}))
		})

		got, _, err := zpacommon.GetAllPagesGeneric[zpaListItem](
			context.Background(), api.Service.Client, path, "Alice Group",
		)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "sg-1", got[0].ID)
	})
}

func TestZPACommon_GetAllPagesGenericWithCustomFilters(t *testing.T) {
	t.Parallel()

	t.Run("custom filters with sort", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			assert.Equal(t, "name", r.URL.Query().Get("sortBy"))
			assert.Equal(t, "ASC", r.URL.Query().Get("sortOrder"))
			return common.SuccessResponse(common.ZPAList([]zpaListItem{{ID: "1", Name: "App"}}))
		})

		got, _, err := zpacommon.GetAllPagesGenericWithCustomFilters[zpaListItem](context.Background(), api.Service.Client, path, zpacommon.Filter{
			SortBy:    "name",
			SortOrder: "ASC",
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("microTenantID passed as query param", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")
		mt := "micro-tenant-1"

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			assert.Equal(t, mt, r.URL.Query().Get("microtenantId"))
			assert.Equal(t, "1", r.URL.Query().Get("page"))
			return common.SuccessResponse(common.ZPAList([]zpaListItem{{ID: "1", Name: "Scoped"}}))
		})

		got, _, err := zpacommon.GetAllPagesGenericWithCustomFilters[zpaListItem](context.Background(), api.Service.Client, path, zpacommon.Filter{
			MicroTenantID: &mt,
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("multi-page", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			switch r.URL.Query().Get("page") {
			case "1":
				return common.SuccessResponse(common.ZPAListPaged([]zpaListItem{{ID: "a", Name: "A"}}, 2))
			case "2":
				return common.SuccessResponse(common.ZPAListPaged([]zpaListItem{{ID: "b", Name: "B"}}, 2))
			default:
				t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
				return common.NotFoundResponse()
			}
		})

		got, _, err := zpacommon.GetAllPagesGenericWithCustomFilters[zpaListItem](context.Background(), api.Service.Client, path, zpacommon.Filter{})
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("partial search fallback on multi-word query", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "application")
		call := 0

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			call++
			if call == 1 {
				return common.NotFoundResponse()
			}
			return common.SuccessResponse(common.ZPAList([]zpaListItem{
				{ID: "1", Name: "My Production App"},
			}))
		})

		got, _, err := zpacommon.GetAllPagesGenericWithCustomFilters[zpaListItem](context.Background(), api.Service.Client, path, zpacommon.Filter{
			Search: "My Production App Extra",
		})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.GreaterOrEqual(t, call, 2)
	})
}

func TestZPACommon_GetAllPagesGenericWithPostSearch(t *testing.T) {
	t.Parallel()

	t.Run("single page", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "tagGroup", "search")

		api.OnFunc("POST", path, func(_ *http.Request, _ []byte) common.MockResponse {
			return common.SuccessResponse(common.ZPAListPaged([]tag_group.TagGroup{
				{ID: "tg-1", Name: "Only"},
			}, 1))
		})

		got, _, err := zpacommon.GetAllPagesGenericWithPostSearch[tag_group.TagGroup](
			context.Background(), api.Service.Client, path, zpacommon.SearchRequest{}, zpacommon.Filter{},
		)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("multi-page", func(t *testing.T) {
		api := common.NewZPATest(t)
		path := common.ZPAPath(api.CustomerID, "tagGroup", "search")
		call := 0

		api.OnFunc("POST", path, func(_ *http.Request, body []byte) common.MockResponse {
			call++
			switch call {
			case 1:
				return common.SuccessResponse(common.ZPAListPaged([]tag_group.TagGroup{
					{ID: "tg-1", Name: "Group 1"},
				}, 2))
			case 2:
				return common.SuccessResponse(common.ZPAListPaged([]tag_group.TagGroup{
					{ID: "tg-2", Name: "Group 2"},
				}, 2))
			default:
				t.Fatalf("unexpected POST call %d", call)
				return common.NotFoundResponse()
			}
		})

		got, _, err := zpacommon.GetAllPagesGenericWithPostSearch[tag_group.TagGroup](
			context.Background(), api.Service.Client, path, zpacommon.SearchRequest{}, zpacommon.Filter{},
		)
		require.NoError(t, err)
		assert.Len(t, got, 2)
	})
}

type scimZpaResource struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func newScimZpaTestClient(t *testing.T, server *common.TestServer) *zpa.ScimZpaClient {
	t.Helper()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	return &zpa.ScimZpaClient{
		ScimConfig: &zpa.ZPAScimConfig{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Transport: &common.MockTransport{TestServerURL: server.URL},
			},
			AuthToken: "mock-scim-token",
			Logger:    logger.GetDefaultLogger("zpa-scim-test: "),
		},
	}
}

func TestZPACommon_GetAllPagesScimGenericWithSearch(t *testing.T) {
	t.Parallel()

	const basePath = "/Users"

	t.Run("aggregate all pages", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)
		call := 0

		server.OnFunc("GET", basePath, func(r *http.Request, _ []byte) common.MockResponse {
			call++
			start := r.URL.Query().Get("startIndex")
			if start == "1" || start == "" {
				items := make([]scimZpaResource, 10)
				for i := range items {
					items[i] = scimZpaResource{ID: "u-1", DisplayName: "User"}
				}
				return common.SuccessResponse(map[string]any{
					"Resources":    items,
					"totalResults": 15,
				})
			}
			return common.SuccessResponse(map[string]any{
				"Resources":    []scimZpaResource{{ID: "u-2", DisplayName: "User"}},
				"totalResults": 15,
			})
		})

		client := newScimZpaTestClient(t, server)
		got, _, err := zpacommon.GetAllPagesScimGenericWithSearch[scimZpaResource](context.Background(), client, basePath, 10, nil)
		require.NoError(t, err)
		assert.Len(t, got, 11)
	})

	t.Run("search short-circuit", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		server.On("GET", basePath, common.SuccessResponse(map[string]any{
			"Resources": []scimZpaResource{
				{ID: "1", DisplayName: "Alice"},
				{ID: "2", DisplayName: "Bob"},
			},
			"totalResults": 2,
		}))

		client := newScimZpaTestClient(t, server)
		got, _, err := zpacommon.GetAllPagesScimGenericWithSearch(context.Background(), client, basePath, 100,
			func(u scimZpaResource) bool { return u.DisplayName == "Bob" },
		)
		require.NoError(t, err)
		require.Len(t, got, 1)
	})

	t.Run("do request failure", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		baseURL, err := url.Parse(server.URL)
		require.NoError(t, err)

		client := &zpa.ScimZpaClient{
			ScimConfig: &zpa.ZPAScimConfig{
				BaseURL: baseURL,
				HTTPClient: &http.Client{
					Transport: errRoundTripper{err: errors.New("scim transport failure")},
				},
				AuthToken: "mock-scim-token",
				Logger:    logger.GetDefaultLogger("zpa-scim-test: "),
			},
		}

		_, _, err = zpacommon.GetAllPagesScimGenericWithSearch[scimZpaResource](context.Background(), client, "/Users", 10, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching paginated data")
	})
}

type errRoundTripper struct {
	err error
}

func (e errRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, e.err
}
