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
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestNetworkServices_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/zia/api/v1/networkServices/12345"

	server.On("GET", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:          serviceID,
		Name:        "HTTPS",
		Description: "HTTPS Traffic",
		Type:        "STANDARD",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.Get(context.Background(), service, serviceID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceID, result.ID)
	assert.Equal(t, "HTTPS", result.Name)
}

func TestNetworkServices_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/networkServices"

	server.On("GET", path, common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "HTTP", Type: "STANDARD"},
		{ID: 2, Name: "HTTPS", Type: "STANDARD"},
		{ID: 3, Name: "Custom Service", Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)

	require.NoError(t, err)
	assert.Len(t, result, 3)
}

func TestNetworkServices_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/networkServices"

	server.On("POST", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:   99999,
		Name: "New Service",
		Type: "CUSTOM",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newSvc := &networkservices.NetworkServices{
		Name: "New Service",
		Type: "CUSTOM",
	}

	result, err := networkservices.Create(context.Background(), service, newSvc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestNetworkServices_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/zia/api/v1/networkServices/12345"

	server.On("PUT", path, common.SuccessResponse(networkservices.NetworkServices{
		ID:   serviceID,
		Name: "Updated Service",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateSvc := &networkservices.NetworkServices{
		ID:   serviceID,
		Name: "Updated Service",
	}

	result, _, err := networkservices.Update(context.Background(), service, serviceID, updateSvc)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Service", result.Name)
}

func TestNetworkServices_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceID := 12345
	path := "/zia/api/v1/networkServices/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = networkservices.Delete(context.Background(), service, serviceID)

	require.NoError(t, err)
}

func TestNetworkServices_GetByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	serviceName := "Custom HTTPS"
	path := "/zia/api/v1/networkServices"

	server.On("GET", path, common.SuccessResponse([]networkservices.NetworkServices{
		{ID: 1, Name: "HTTP", Type: "STANDARD"},
		{ID: 2, Name: serviceName, Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetByName(context.Background(), service, serviceName, nil, nil)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, serviceName, result.Name)
}

// rejectPagingNetworkServices serves /networkServices like the beta cloud
// (terraform-provider-zia issue #617): the endpoint is not paginated, and a
// request carrying page or pageSize is rejected with "Invalid Page Size".
func rejectPagingNetworkServices(t *testing.T, items []networkservices.NetworkServices) common.MockResponseFunc {
	t.Helper()
	return func(req *http.Request, _ []byte) common.MockResponse {
		q := req.URL.Query()
		if q.Has("pageSize") || q.Has("page") {
			return common.MockResponse{
				StatusCode: http.StatusBadRequest,
				Body: map[string]any{
					"code":    "INVALID_INPUT_ARGUMENT",
					"message": "Invalid Page Size " + q.Get("pageSize"),
				},
			}
		}
		return common.SuccessResponse(items)
	}
}

// assertSingleUnpagedRead checks that exactly one list request reached the
// API and that it carried no page or pageSize parameter.
func assertSingleUnpagedRead(t *testing.T, server *common.TestServer) url.Values {
	t.Helper()
	reqs := server.AllRequests()
	var lists []common.RecordedRequest
	for _, r := range reqs {
		if r.Method == http.MethodGet && r.Path == "/zia/api/v1/networkServices" {
			lists = append(lists, r)
		}
	}
	require.Len(t, lists, 1, "network services must be read in a single request")
	q, err := url.ParseQuery(lists[0].Query)
	require.NoError(t, err)
	assert.False(t, q.Has("pageSize"), "pageSize must not be sent: %q", lists[0].Query)
	assert.False(t, q.Has("page"), "page must not be sent: %q", lists[0].Query)
	return q
}

func TestNetworkServices_GetAll_NoPagination_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// More than any page size the SDK has used, so a paged read would either be
	// rejected or stop after the first page.
	items := make([]networkservices.NetworkServices, 1200)
	for i := range items {
		items[i] = networkservices.NetworkServices{ID: i + 1, Name: fmt.Sprintf("svc-%04d", i+1)}
	}
	server.OnFunc("GET", "/zia/api/v1/networkServices", rejectPagingNetworkServices(t, items))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result, len(items))

	q := assertSingleUnpagedRead(t, server)
	assert.Empty(t, q, "no query parameters expected without filters")
}

func TestNetworkServices_GetAll_FiltersWithoutPagination_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", "/zia/api/v1/networkServices", rejectPagingNetworkServices(t, []networkservices.NetworkServices{
		{ID: 1, Name: "HTTPS", Protocol: "TCP"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	protocol, locale := "TCP", "en-US"
	result, err := networkservices.GetAllNetworkServices(context.Background(), service, &protocol, &locale)
	require.NoError(t, err)
	require.Len(t, result, 1)

	q := assertSingleUnpagedRead(t, server)
	assert.Equal(t, "TCP", q.Get("protocol"))
	assert.Equal(t, "en-US", q.Get("locale"))
}

func TestNetworkServices_GetByName_NoPagination_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", "/zia/api/v1/networkServices", rejectPagingNetworkServices(t, []networkservices.NetworkServices{
		{ID: 1, Name: "HTTP", Type: "PREDEFINED"},
		{ID: 2, Name: "HTTPS", Type: "PREDEFINED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// Same lookup as the issue: name "https", matched case-insensitively.
	result, err := networkservices.GetByName(context.Background(), service, "https", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 2, result.ID)

	assertSingleUnpagedRead(t, server)
}

func TestNetworkServices_GetByName_APIErrorReturned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.On("GET", "/zia/api/v1/networkServices", common.MockResponse{
		StatusCode: http.StatusBadRequest,
		Body:       map[string]any{"code": "INVALID_INPUT_ARGUMENT", "message": "bad request"},
	})

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetByName(context.Background(), service, "HTTPS", nil, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestNetworkServices_GetAll_JMESPathStillApplied_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", "/zia/api/v1/networkServices", rejectPagingNetworkServices(t, []networkservices.NetworkServices{
		{ID: 1, Name: "HTTP", Type: "PREDEFINED"},
		{ID: 2, Name: "Custom", Type: "CUSTOM"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	ctx := zscaler.ContextWithJMESPath(context.Background(), "[?type=='CUSTOM']")
	result, err := networkservices.GetAllNetworkServices(ctx, service, nil, nil)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "Custom", result[0].Name)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestNetworkServices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkServices JSON marshaling", func(t *testing.T) {
		svc := networkservices.NetworkServices{
			ID:           12345,
			Name:         "Custom HTTPS",
			Description:  "Custom HTTPS service",
			Type:         "CUSTOM",
			SrcTCPPorts:  []networkservices.NetworkPorts{{Start: 443, End: 443}},
			DestTCPPorts: []networkservices.NetworkPorts{{Start: 443, End: 443}},
		}

		data, err := json.Marshal(svc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Custom HTTPS"`)
	})

	t.Run("NetworkServices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "DNS",
			"type": "STANDARD",
			"description": "DNS service",
			"destTcpPorts": [{"start": 53, "end": 53}],
			"destUdpPorts": [{"start": 53, "end": 53}]
		}`

		var svc networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonData), &svc)
		require.NoError(t, err)

		assert.Equal(t, 54321, svc.ID)
		assert.Equal(t, "DNS", svc.Name)
	})
}
