package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservicegroups"
)

const networkServiceGroupsPath = "/zia/api/v1/networkServiceGroups"

func sampleNetworkServiceGroup() networkservicegroups.NetworkServiceGroups {
	return networkservicegroups.NetworkServiceGroups{
		Name:        "tests-nw-svc-group",
		Description: "tests-nw-svc-group",
		Services: []networkservicegroups.Services{
			{ID: 100},
			{ID: 101},
			{ID: 102},
		},
	}
}

func TestNetworkServiceGroups_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	group := sampleNetworkServiceGroup()
	group.ID = groupID

	server.On("GET", networkServiceGroupsPath+"/100", common.SuccessResponse(group))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetNetworkServiceGroups(context.Background(), service, groupID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupID, result.ID)
	assert.Len(t, result.Services, 3)
}

func TestNetworkServiceGroups_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupName := "tests-nw-svc-group"
	server.On("GET", networkServiceGroupsPath, common.SuccessResponse([]networkservicegroups.NetworkServiceGroups{
		{ID: 1, Name: "Other Group"},
		func() networkservicegroups.NetworkServiceGroups {
			g := sampleNetworkServiceGroup()
			g.ID = 100
			return g
		}(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetNetworkServiceGroupsByName(context.Background(), service, groupName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, groupName, result.Name)
}

func TestNetworkServiceGroups_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	created := sampleNetworkServiceGroup()
	created.ID = 99999

	server.On("POST", networkServiceGroupsPath, common.SuccessResponse(created))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	input := sampleNetworkServiceGroup()
	result, err := networkservicegroups.CreateNetworkServiceGroups(context.Background(), service, &input)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestNetworkServiceGroups_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	groupID := 100
	updated := sampleNetworkServiceGroup()
	updated.ID = groupID
	updated.Name = "updated-nw-svc-group"

	server.On("PUT", networkServiceGroupsPath+"/100", common.SuccessResponse(updated))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := networkservicegroups.UpdateNetworkServiceGroups(context.Background(), service, groupID, &updated)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "updated-nw-svc-group", result.Name)
}

func TestNetworkServiceGroups_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("DELETE", networkServiceGroupsPath+"/100", common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkservicegroups.DeleteNetworkServiceGroups(context.Background(), service, 100)

	require.NoError(t, err)
}

func TestNetworkServiceGroups_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", networkServiceGroupsPath, common.SuccessResponse([]networkservicegroups.NetworkServiceGroups{
		sampleNetworkServiceGroup(),
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 1)
}

// networkServiceGroupsMaxPageSize is the page size limit the API enforces on a
// tenant without the extended service group limit.
const networkServiceGroupsMaxPageSize = 127

func makeNetworkServiceGroups(n int) []networkservicegroups.NetworkServiceGroups {
	items := make([]networkservicegroups.NetworkServiceGroups, n)
	for i := range items {
		items[i] = networkservicegroups.NetworkServiceGroups{ID: i + 1, Name: fmt.Sprintf("group-%04d", i+1)}
	}
	return items
}

// pagedList serves a list endpoint like the ZIA API after the page size
// validation change: a pageSize above maxPageSize is rejected with
// "Invalid Page Size", and page/pageSize select a slice of the list.
func pagedList[T any](t *testing.T, items []T, maxPageSize int) common.MockResponseFunc {
	t.Helper()
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

// listQueries returns the query of every GET sent to path.
func listQueries(t *testing.T, server *common.TestServer, path string) []url.Values {
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

// assertPagedReads checks that the reads asked for pages 1..n of size 127.
func assertPagedReads(t *testing.T, queries []url.Values, n int) {
	t.Helper()
	require.Len(t, queries, n)
	for i, q := range queries {
		assert.Equal(t, "127", q.Get("pageSize"), "request %d", i+1)
		assert.Equal(t, strconv.Itoa(i+1), q.Get("page"), "request %d", i+1)
	}
}

func TestNetworkServiceGroups_GetAll_PageSizeWithinLimit_SDK(t *testing.T) {
	for _, n := range []int{0, 1, 126, 127, 128, 254, 1000} {
		t.Run(fmt.Sprintf("%d groups", n), func(t *testing.T) {
			server := common.NewTestServer()
			defer server.Close()

			items := makeNetworkServiceGroups(n)
			server.OnFunc("GET", networkServiceGroupsPath, pagedList(t, items, networkServiceGroupsMaxPageSize))

			service, err := common.CreateTestService(context.Background(), server, "123456")
			require.NoError(t, err)

			result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
			require.NoError(t, err)
			require.Len(t, result, n)
			for i, g := range result {
				assert.Equal(t, i+1, g.ID)
			}

			// One request per full page, plus the short (possibly empty) last page.
			assertPagedReads(t, listQueries(t, server, networkServiceGroupsPath), n/networkServiceGroupsMaxPageSize+1)
		})
	}
}

func TestNetworkServiceGroups_GetByName_PageSizeWithinLimit_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	items := makeNetworkServiceGroups(300)
	server.OnFunc("GET", networkServiceGroupsPath, pagedList(t, items, networkServiceGroupsMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// On the third page, matched case-insensitively.
	result, err := networkservicegroups.GetNetworkServiceGroupsByName(context.Background(), service, "GROUP-0290")
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 290, result.ID)
}

func TestNetworkServiceGroups_GetAll_APIIgnoresPageSize_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API returns the whole list whatever page and pageSize say.
	items := makeNetworkServiceGroups(500)
	server.On("GET", networkServiceGroupsPath, common.SuccessResponse(items))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, 500)
	assert.Len(t, listQueries(t, server, networkServiceGroupsPath), 1)
}

func TestNetworkServiceGroups_GetAll_APIIgnoresPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API honours pageSize but always returns the first page.
	items := makeNetworkServiceGroups(networkServiceGroupsMaxPageSize)
	server.On("GET", networkServiceGroupsPath, common.SuccessResponse(items))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.NoError(t, err)
	assert.Len(t, result, networkServiceGroupsMaxPageSize)
	assert.Len(t, listQueries(t, server, networkServiceGroupsPath), 2)
}

func TestNetworkServiceGroups_GetAll_APIErrorReturned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// A tenant whose limit is below the SDK page size.
	server.OnFunc("GET", networkServiceGroupsPath, pagedList(t, makeNetworkServiceGroups(10), 100))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservicegroups.GetAllNetworkServiceGroups(context.Background(), service)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Page Size")
	assert.Nil(t, result)
}

func TestNetworkServiceGroups_GetAll_JMESPathAcrossPages_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	items := makeNetworkServiceGroups(200)
	items[10].Description = "match"
	items[150].Description = "match"
	server.OnFunc("GET", networkServiceGroupsPath, pagedList(t, items, networkServiceGroupsMaxPageSize))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ctx := zscaler.ContextWithJMESPath(context.Background(), "[?description=='match']")
	result, err := networkservicegroups.GetAllNetworkServiceGroups(ctx, service)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 11, result[0].ID)
	assert.Equal(t, 151, result[1].ID)
}

func TestNetworkServiceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("JSON marshaling", func(t *testing.T) {
		group := sampleNetworkServiceGroup()
		group.ID = 100

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"name":"tests-nw-svc-group"`)
	})
}
