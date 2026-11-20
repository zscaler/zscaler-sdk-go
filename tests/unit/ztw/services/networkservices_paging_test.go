package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/ztw/services/policyresources/networkservices"
)

const (
	ztwNetworkServicesPath      = "/ztw/api/v1/networkServices"
	ztwNetworkServiceGroupsPath = "/ztw/api/v1/networkServiceGroups"

	// ztwMaxPageSize is the page size limit the API enforces on a tenant
	// without the extended limit.
	ztwMaxPageSize = 127
)

func makeZTWNetworkServices(n int) []networkservices.NetworkServices {
	items := make([]networkservices.NetworkServices, n)
	for i := range items {
		items[i] = networkservices.NetworkServices{ID: i + 1, Name: fmt.Sprintf("service-%04d", i+1)}
	}
	return items
}

func makeZTWNetworkServiceGroups(n int) []networkservicegroups.NetworkServiceGroups {
	items := make([]networkservicegroups.NetworkServiceGroups, n)
	for i := range items {
		items[i] = networkservicegroups.NetworkServiceGroups{ID: i + 1, Name: fmt.Sprintf("group-%04d", i+1)}
	}
	return items
}

// ztwPagedList serves a list endpoint like the API after the page size
// validation change: a pageSize above maxPageSize is rejected with
// "Invalid Page Size", and page/pageSize select a slice of the list.
func ztwPagedList[T any](items []T, maxPageSize int) common.MockResponseFunc {
	return func(req *http.Request, _ []byte) common.MockResponse {
		q := req.URL.Query()
		size, page := maxPageSize, 1
		if v := q.Get("pageSize"); v != "" {
			size, _ = strconv.Atoi(v)
		}
		if v := q.Get("page"); v != "" {
			page, _ = strconv.Atoi(v)
		}
		if size < 1 || size > maxPageSize {
			return common.MockResponse{
				StatusCode: http.StatusBadRequest,
				Body: map[string]any{
					"code":    "INVALID_INPUT_ARGUMENT",
					"message": "Invalid Page Size " + q.Get("pageSize"),
				},
			}
		}
		start := (page - 1) * size
		if start >= len(items) {
			return common.SuccessResponse([]T{})
		}
		end := min(start+size, len(items))
		return common.SuccessResponse(items[start:end])
	}
}

// ztwListQueries returns the query of every GET sent to path.
func ztwListQueries(t *testing.T, server *common.TestServer, path string) []url.Values {
	t.Helper()
	var out []url.Values
	for _, r := range server.AllRequests() {
		if r.Method == http.MethodGet && r.Path == path {
			q, err := url.ParseQuery(r.Query)
			require.NoError(t, err)
			out = append(out, q)
		}
	}
	return out
}

// ztwAssertPagedReads checks that the reads asked for pages 1..n of size 127.
func ztwAssertPagedReads(t *testing.T, queries []url.Values, n int) {
	t.Helper()
	require.Len(t, queries, n)
	for i, q := range queries {
		assert.Equal(t, "127", q.Get("pageSize"), "request %d", i+1)
		assert.Equal(t, strconv.Itoa(i+1), q.Get("page"), "request %d", i+1)
	}
}

// =====================================================
// Network Services paging
// =====================================================

func TestZTWNetworkServices_GetAll_PageSizeWithinLimit_SDK(t *testing.T) {
	for _, n := range []int{0, 1, 126, 127, 128, 254, 1200} {
		t.Run(fmt.Sprintf("%d services", n), func(t *testing.T) {
			server := common.NewTestServer()
			defer server.Close()

			server.OnFunc("GET", ztwNetworkServicesPath, ztwPagedList(makeZTWNetworkServices(n), ztwMaxPageSize))

			service, err := common.CreateTestService(context.Background(), server, "123456")
			require.NoError(t, err)

			result, err := networkservices.GetAllNetworkServices(context.Background(), service)
			require.NoError(t, err)
			require.Len(t, result, n)
			for i, s := range result {
				assert.Equal(t, i+1, s.ID)
			}

			// One request per full page, plus the short (possibly empty) last page.
			ztwAssertPagedReads(t, ztwListQueries(t, server, ztwNetworkServicesPath), n/ztwMaxPageSize+1)
		})
	}
}

func TestZTWNetworkServices_GetByName_LaterPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", ztwNetworkServicesPath, ztwPagedList(makeZTWNetworkServices(300), ztwMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// On the third page, matched case-insensitively.
	result, err := networkservices.GetByName(context.Background(), service, "SERVICE-0290")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 290, result.ID)
}

func TestZTWNetworkServices_GetAll_APIIgnoresPageSize_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API returns the whole list whatever page and pageSize say.
	server.On("GET", ztwNetworkServicesPath, common.SuccessResponse(makeZTWNetworkServices(500)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 500)
	assert.Len(t, ztwListQueries(t, server, ztwNetworkServicesPath), 1)
}

func TestZTWNetworkServices_GetAll_APIIgnoresPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API honours pageSize but always returns the first page.
	server.On("GET", ztwNetworkServicesPath, common.SuccessResponse(makeZTWNetworkServices(ztwMaxPageSize)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, ztwMaxPageSize)
	assert.Len(t, ztwListQueries(t, server, ztwNetworkServicesPath), 2)
}

func TestZTWNetworkServices_GetAll_APIErrorReturned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// A tenant whose limit is below the SDK page size.
	server.OnFunc("GET", ztwNetworkServicesPath, ztwPagedList(makeZTWNetworkServices(10), 100))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Page Size")
	assert.Nil(t, result)
}

func TestZTWNetworkServices_GetAll_JMESPathAcrossPages_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	items := makeZTWNetworkServices(200)
	items[10].Type = "CUSTOM"
	items[150].Type = "CUSTOM"
	server.OnFunc("GET", ztwNetworkServicesPath, ztwPagedList(items, ztwMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ctx := zscaler.ContextWithJMESPath(context.Background(), "[?type=='CUSTOM']")
	result, err := networkservices.GetAllNetworkServices(ctx, service)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 11, result[0].ID)
	assert.Equal(t, 151, result[1].ID)
}

// =====================================================
// Network Service Groups paging
// =====================================================

func TestZTWNetworkServiceGroups_GetAll_PageSizeWithinLimit_SDK(t *testing.T) {
	for _, n := range []int{0, 1, 126, 127, 128, 254, 1000} {
		t.Run(fmt.Sprintf("%d groups", n), func(t *testing.T) {
			server := common.NewTestServer()
			defer server.Close()

			server.OnFunc("GET", ztwNetworkServiceGroupsPath, ztwPagedList(makeZTWNetworkServiceGroups(n), ztwMaxPageSize))

			service, err := common.CreateTestService(context.Background(), server, "123456")
			require.NoError(t, err)

			result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
			require.NoError(t, err)
			require.Len(t, result, n)
			for i, g := range result {
				assert.Equal(t, i+1, g.ID)
			}

			ztwAssertPagedReads(t, ztwListQueries(t, server, ztwNetworkServiceGroupsPath), n/ztwMaxPageSize+1)
		})
	}
}

func TestZTWNetworkServiceGroups_GetByName_LaterPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", ztwNetworkServiceGroupsPath, ztwPagedList(makeZTWNetworkServiceGroups(300), ztwMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetNetworkServiceGroupsByName(context.Background(), service, "GROUP-0290")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 290, result.ID)
}

func TestZTWNetworkServiceGroups_GetAll_APIIgnoresPageSize_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", ztwNetworkServiceGroupsPath, common.SuccessResponse(makeZTWNetworkServiceGroups(500)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 500)
	assert.Len(t, ztwListQueries(t, server, ztwNetworkServiceGroupsPath), 1)
}

func TestZTWNetworkServiceGroups_GetAll_APIIgnoresPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", ztwNetworkServiceGroupsPath, common.SuccessResponse(makeZTWNetworkServiceGroups(ztwMaxPageSize)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, ztwMaxPageSize)
	assert.Len(t, ztwListQueries(t, server, ztwNetworkServiceGroupsPath), 2)
}

func TestZTWNetworkServiceGroups_GetAll_APIErrorReturned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", ztwNetworkServiceGroupsPath, ztwPagedList(makeZTWNetworkServiceGroups(10), 100))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Page Size")
	assert.Nil(t, result)
}

func TestZTWNetworkServiceGroups_GetAll_JMESPathAcrossPages_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	items := makeZTWNetworkServiceGroups(200)
	items[10].Description = "match"
	items[150].Description = "match"
	server.OnFunc("GET", ztwNetworkServiceGroupsPath, ztwPagedList(items, ztwMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ctx := zscaler.ContextWithJMESPath(context.Background(), "[?description=='match']")
	result, err := networkservicegroups.GetAllNetworkServiceGroups(ctx, service)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 11, result[0].ID)
	assert.Equal(t, 151, result[1].ID)
}
