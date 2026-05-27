// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/logger"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
)

func TestZIACommon_Structure(t *testing.T) {
	t.Parallel()

	t.Run("IDNameExtensions JSON marshaling", func(t *testing.T) {
		idName := ziacommon.IDNameExtensions{
			ID:   12345,
			Name: "Test Resource",
			Extensions: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
			},
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Test Resource"`)
		assert.Contains(t, string(data), `"extensions"`)
	})

	t.Run("IDName JSON marshaling", func(t *testing.T) {
		idName := ziacommon.IDName{
			ID:     12345,
			Name:   "Test Resource",
			Parent: "Parent Resource",
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"parent":"Parent Resource"`)
	})

	t.Run("IDNameExternalID JSON marshaling", func(t *testing.T) {
		idName := ziacommon.IDNameExternalID{
			ID:         12345,
			Name:       "Test Resource",
			ExternalID: "ext-12345",
		}

		data, err := json.Marshal(idName)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"externalId":"ext-12345"`)
	})

	t.Run("UserGroups JSON marshaling", func(t *testing.T) {
		group := ziacommon.UserGroups{
			ID:              12345,
			Name:            "Engineering",
			IdpID:           100,
			Comments:        "Engineering team",
			IsSystemDefined: "false",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"idp_id":100`)
	})

	t.Run("UserDepartment JSON marshaling", func(t *testing.T) {
		dept := ziacommon.UserDepartment{
			ID:       12345,
			Name:     "Engineering",
			IdpID:    100,
			Comments: "Engineering department",
			Deleted:  false,
		}

		data, err := json.Marshal(dept)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Engineering"`)
	})

	t.Run("DeviceGroups JSON marshaling", func(t *testing.T) {
		dg := ziacommon.DeviceGroups{
			ID:   12345,
			Name: "Mobile Devices",
		}

		data, err := json.Marshal(dg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Mobile Devices"`)
	})

	t.Run("CommonNSS JSON marshaling", func(t *testing.T) {
		nss := ziacommon.CommonNSS{
			ID:          12345,
			PID:         100,
			Name:        "NSS Server",
			Description: "NSS server description",
			Deleted:     false,
			GetlID:      200,
		}

		data, err := json.Marshal(nss)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"pid":100`)
	})

	t.Run("ZPAAppSegments JSON marshaling", func(t *testing.T) {
		segment := ziacommon.ZPAAppSegments{
			ID:         12345,
			Name:       "ZPA Segment",
			ExternalID: "zpa-ext-123",
		}

		data, err := json.Marshal(segment)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"externalId":"zpa-ext-123"`)
	})

	t.Run("DatacenterSearchParameters structure", func(t *testing.T) {
		params := ziacommon.DatacenterSearchParameters{
			RoutableIP:                true,
			WithinCountryOnly:         true,
			IncludePrivateServiceEdge: false,
			IncludeCurrentVips:        true,
			SourceIp:                  "10.0.0.1",
			Latitude:                  37.7749,
			Longitude:                 -122.4194,
			Subcloud:                  "subcloud1",
		}

		assert.True(t, params.RoutableIP)
		assert.True(t, params.WithinCountryOnly)
		assert.Equal(t, "10.0.0.1", params.SourceIp)
	})

	t.Run("Order JSON marshaling", func(t *testing.T) {
		order := ziacommon.Order{
			On: "name",
			By: "asc",
		}

		data, err := json.Marshal(order)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"on":"name"`)
		assert.Contains(t, string(data), `"by":"asc"`)
	})

	t.Run("DataConsumed JSON marshaling", func(t *testing.T) {
		dc := ziacommon.DataConsumed{
			Min: 100,
			Max: 1000,
		}

		data, err := json.Marshal(dc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"min":100`)
		assert.Contains(t, string(data), `"max":1000`)
	})
}

func TestZIACommon_SortParams(t *testing.T) {
	t.Parallel()

	t.Run("GetSortParams with both parameters", func(t *testing.T) {
		params := ziacommon.GetSortParams(ziacommon.NameSortField, ziacommon.ASCSortOrder)
		assert.Contains(t, params, "sortBy=name")
		assert.Contains(t, params, "sortOrder=asc")
	})

	t.Run("GetSortParams with only sortBy", func(t *testing.T) {
		params := ziacommon.GetSortParams(ziacommon.IDSortField, "")
		assert.Equal(t, "sortBy=id", params)
	})

	t.Run("GetSortParams with only sortOrder", func(t *testing.T) {
		params := ziacommon.GetSortParams("", ziacommon.DESCSortOrder)
		assert.Equal(t, "sortOrder=desc", params)
	})

	t.Run("GetPageSize returns correct value", func(t *testing.T) {
		pageSize := ziacommon.GetPageSize()
		assert.Equal(t, 1000, pageSize)
	})

	t.Run("GetSortParams with neither parameter", func(t *testing.T) {
		params := ziacommon.GetSortParams("", "")
		assert.Empty(t, params)
	})
}

type pageItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestZIACommon_ReadAllPages(t *testing.T) {
	t.Parallel()

	t.Run("single page", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			assert.Equal(t, "1", r.URL.Query().Get("page"))
			assert.Equal(t, "1000", r.URL.Query().Get("pageSize"))
			return common.SuccessResponse([]pageItem{{ID: 1, Name: "one"}})
		})

		var list []pageItem
		err := ziacommon.ReadAllPages(context.Background(), api.Service.Client, path, &list)
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, "one", list[0].Name)
	})

	t.Run("multi-page with custom page size", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			switch r.URL.Query().Get("page") {
			case "1":
				assert.Equal(t, "2", r.URL.Query().Get("pageSize"))
				return common.SuccessResponse([]pageItem{
					{ID: 1, Name: "one"},
					{ID: 2, Name: "two"},
				})
			case "2":
				return common.SuccessResponse([]pageItem{{ID: 3, Name: "three"}})
			default:
				t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
				return common.NotFoundResponse()
			}
		})

		var list []pageItem
		err := ziacommon.ReadAllPages(context.Background(), api.Service.Client, path, &list, 2)
		require.NoError(t, err)
		require.Len(t, list, 3)
	})

	t.Run("endpoint already has query string", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")
		endpoint := path + "?search=foo"

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			assert.Equal(t, "foo", r.URL.Query().Get("search"))
			assert.Equal(t, "1", r.URL.Query().Get("page"))
			assert.Equal(t, "1000", r.URL.Query().Get("pageSize"))
			return common.SuccessResponse([]pageItem{{ID: 1, Name: "one"}})
		})

		var list []pageItem
		err := ziacommon.ReadAllPages(context.Background(), api.Service.Client, endpoint, &list)
		require.NoError(t, err)
		require.Len(t, list, 1)
	})

	t.Run("JMESPath filter from context", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")

		api.On("GET", path, common.SuccessResponse([]pageItem{
			{ID: 1, Name: "active"},
			{ID: 2, Name: "inactive"},
		}))

		ctx := zscaler.ContextWithJMESPath(context.Background(), "[?name=='active']")
		var list []pageItem
		err := ziacommon.ReadAllPages(ctx, api.Service.Client, path, &list)
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, "active", list[0].Name)
	})

	t.Run("nil list is no-op", func(t *testing.T) {
		api := common.NewZIATest(t)
		err := ziacommon.ReadAllPages[pageItem](context.Background(), api.Service.Client, common.ZIAPath("testResources"), nil)
		require.NoError(t, err)
	})

	t.Run("read error propagates", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")
		api.On("GET", path, common.NotFoundResponse())

		var list []pageItem
		err := ziacommon.ReadAllPages(context.Background(), api.Service.Client, path, &list)
		require.Error(t, err)
	})
}

func TestZIACommon_ReadPage(t *testing.T) {
	t.Parallel()

	t.Run("page without existing query", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			assert.Equal(t, "2", r.URL.Query().Get("page"))
			assert.Equal(t, "50", r.URL.Query().Get("pageSize"))
			return common.SuccessResponse([]pageItem{{ID: 10, Name: "ten"}})
		})

		var list []pageItem
		err := ziacommon.ReadPage(context.Background(), api.Service.Client, path, 2, &list, 50)
		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, 10, list[0].ID)
	})

	t.Run("page preserves and merges existing query", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")
		endpoint := path + "?search=foo&page=99"

		api.OnFunc("GET", path, func(r *http.Request, _ []byte) common.MockResponse {
			q := r.URL.Query()
			assert.Equal(t, "foo", q.Get("search"))
			assert.Equal(t, "3", q.Get("page"))
			assert.Equal(t, "1000", q.Get("pageSize"))
			return common.SuccessResponse([]pageItem{{ID: 1, Name: "one"}})
		})

		var list []pageItem
		err := ziacommon.ReadPage(context.Background(), api.Service.Client, endpoint, 3, &list)
		require.NoError(t, err)
		require.Len(t, list, 1)
	})

	t.Run("nil list is no-op", func(t *testing.T) {
		api := common.NewZIATest(t)
		err := ziacommon.ReadPage[pageItem](context.Background(), api.Service.Client, common.ZIAPath("testResources"), 1, nil)
		require.NoError(t, err)
	})

	t.Run("invalid query string returns error", func(t *testing.T) {
		api := common.NewZIATest(t)
		var list []pageItem
		err := ziacommon.ReadPage(context.Background(), api.Service.Client, common.ZIAPath("testResources")+"?a=%", 1, &list)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse query string")
	})

	t.Run("read error propagates", func(t *testing.T) {
		api := common.NewZIATest(t)
		path := common.ZIAPath("testResources")
		api.On("GET", path, common.NotFoundResponse())

		var list []pageItem
		err := ziacommon.ReadPage(context.Background(), api.Service.Client, path, 1, &list)
		require.Error(t, err)
	})
}

type scimResource struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

func newScimTestClient(t *testing.T, server *common.TestServer) *zia.ScimZiaClient {
	t.Helper()
	baseURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	return &zia.ScimZiaClient{
		ScimConfig: &zia.ZIAScimConfig{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Transport: &common.MockTransport{TestServerURL: server.URL},
			},
			AuthToken: "mock-scim-token",
			Logger:    logger.GetDefaultLogger("scim-test: "),
		},
	}
}

func TestZIACommon_GetAllPagesScimPostWithSearch(t *testing.T) {
	t.Parallel()

	const searchPath = "/Users/.search"

	t.Run("single page", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		server.OnFunc("POST", searchPath, func(r *http.Request, body []byte) common.MockResponse {
			assert.Equal(t, "application/scim+json", r.Header.Get("Content-Type"))
			var req map[string]any
			require.NoError(t, json.Unmarshal(body, &req))
			assert.Equal(t, float64(1), req["startIndex"])
			assert.Equal(t, float64(100), req["count"])
			return common.SuccessResponse(map[string]any{
				"Resources": []scimResource{
					{ID: "u-1", DisplayName: "Alice"},
					{ID: "u-2", DisplayName: "Bob"},
				},
				"totalResults": 2,
			})
		})

		client := newScimTestClient(t, server)
		got, resp, err := ziacommon.GetAllPagesScimPostWithSearch[scimResource](
			context.Background(), client, searchPath, 0, nil,
		)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Len(t, got, 2)
	})

	t.Run("multi-page aggregation", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		call := 0
		server.OnFunc("POST", searchPath, func(_ *http.Request, _ []byte) common.MockResponse {
			call++
			switch call {
			case 1:
				items := make([]scimResource, 100)
				for i := range items {
					items[i] = scimResource{ID: fmt.Sprintf("u-%d", i), DisplayName: "User"}
				}
				return common.SuccessResponse(map[string]any{
					"Resources":    items,
					"totalResults": 150,
				})
			case 2:
				items := make([]scimResource, 50)
				for i := range items {
					items[i] = scimResource{ID: fmt.Sprintf("u2-%d", i), DisplayName: "User"}
				}
				return common.SuccessResponse(map[string]any{
					"Resources":    items,
					"totalResults": 150,
				})
			default:
				return common.NotFoundResponse()
			}
		})

		client := newScimTestClient(t, server)
		got, _, err := ziacommon.GetAllPagesScimPostWithSearch[scimResource](
			context.Background(), client, searchPath, 100, nil,
		)
		require.NoError(t, err)
		assert.Len(t, got, 150)
	})

	t.Run("searchFunc short-circuits on first match", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		server.OnFunc("POST", searchPath, func(_ *http.Request, _ []byte) common.MockResponse {
			return common.SuccessResponse(map[string]any{
				"Resources": []scimResource{
					{ID: "u-1", DisplayName: "Alice"},
					{ID: "u-2", DisplayName: "Bob"},
				},
				"totalResults": 2,
			})
		})

		client := newScimTestClient(t, server)
		got, _, err := ziacommon.GetAllPagesScimPostWithSearch(
			context.Background(), client, searchPath, 100,
			func(u scimResource) bool { return u.DisplayName == "Bob" },
		)
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, "Bob", got[0].DisplayName)
	})

	t.Run("itemsPerPage above 100 is clamped", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)

		server.OnFunc("POST", searchPath, func(_ *http.Request, body []byte) common.MockResponse {
			var req map[string]any
			require.NoError(t, json.Unmarshal(body, &req))
			assert.Equal(t, float64(100), req["count"])
			return common.SuccessResponse(map[string]any{
				"Resources":    []scimResource{},
				"totalResults": 0,
			})
		})

		client := newScimTestClient(t, server)
		got, _, err := ziacommon.GetAllPagesScimPostWithSearch[scimResource](
			context.Background(), client, searchPath, 500, nil,
		)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("request error propagates", func(t *testing.T) {
		server := common.NewTestServer()
		t.Cleanup(server.Close)
		server.On("POST", searchPath, common.NotFoundResponse())

		client := newScimTestClient(t, server)
		_, _, err := ziacommon.GetAllPagesScimPostWithSearch[scimResource](
			context.Background(), client, searchPath, 100, nil,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "SCIM POST pagination failed")
	})
}

