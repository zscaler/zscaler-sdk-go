// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

const networkServicesPath = "/zia/api/v1/networkServices"

func makeNetworkServices(n int) []networkservices.NetworkServices {
	items := make([]networkservices.NetworkServices, n)
	for i := range items {
		items[i] = networkservices.NetworkServices{ID: i + 1, Name: fmt.Sprintf("svc-%04d", i+1), Type: "CUSTOM"}
	}
	return items
}

func TestNetworkServices_GetAll_Paged_SDK(t *testing.T) {
	for _, n := range []int{0, 126, 127, 128, 1200} {
		t.Run(fmt.Sprintf("%d services", n), func(t *testing.T) {
			server := common.NewTestServer()
			defer server.Close()

			items := makeNetworkServices(n)
			server.OnFunc("GET", networkServicesPath, pagedList(t, items, 127))

			service, err := common.CreateTestService(context.Background(), server, "123456")
			require.NoError(t, err)

			result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)
			require.NoError(t, err)
			require.Len(t, result, n)
			for i, svc := range result {
				assert.Equal(t, i+1, svc.ID)
			}

			queries := listQueries(t, server, networkServicesPath)
			assertPagedReads(t, queries, n/127+1)
			for _, q := range queries {
				assert.False(t, q.Has("protocol"))
				assert.False(t, q.Has("locale"))
			}
		})
	}
}

func TestNetworkServices_GetAll_FiltersOnEveryPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", networkServicesPath, pagedList(t, makeNetworkServices(200), 127))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	protocol, locale := "TCP", "en-US"
	result, err := networkservices.GetAllNetworkServices(context.Background(), service, &protocol, &locale)
	require.NoError(t, err)
	require.Len(t, result, 200)

	queries := listQueries(t, server, networkServicesPath)
	assertPagedReads(t, queries, 2)
	for _, q := range queries {
		assert.Equal(t, "TCP", q.Get("protocol"))
		assert.Equal(t, "en-US", q.Get("locale"))
	}
}

func TestNetworkServices_GetByName_LaterPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	server.OnFunc("GET", networkServicesPath, pagedList(t, makeNetworkServices(300), 127))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// On the third page, matched case-insensitively. A single unpaged read
	// would only see the first page.
	result, err := networkservices.GetByName(context.Background(), service, "SVC-0290", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 290, result.ID)
}

func TestNetworkServices_GetAll_APIIgnoresPageSize_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API returns the whole list whatever page and pageSize say.
	server.On("GET", networkServicesPath, common.SuccessResponse(makeNetworkServices(500)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result, 500)
	assert.Len(t, listQueries(t, server, networkServicesPath), 1)
}

func TestNetworkServices_GetAll_APIIgnoresPage_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// The API honours pageSize but always returns the first page.
	server.On("GET", networkServicesPath, common.SuccessResponse(makeNetworkServices(127)))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)
	require.NoError(t, err)
	assert.Len(t, result, 127)
	assert.Len(t, listQueries(t, server, networkServicesPath), 2)
}

func TestNetworkServices_GetAll_InvalidPageSizeReturned_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	// A limit below the SDK page size.
	server.OnFunc("GET", networkServicesPath, pagedList(t, makeNetworkServices(10), 100))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := networkservices.GetAllNetworkServices(context.Background(), service, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Invalid Page Size")
	assert.Nil(t, result)
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

func TestNetworkServices_GetAll_JMESPathAcrossPages_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	items := makeNetworkServices(200)
	items[5].Type = "PREDEFINED"
	items[150].Type = "PREDEFINED"
	server.OnFunc("GET", networkServicesPath, pagedList(t, items, 127))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	// Filter matches on both pages.
	ctx := zscaler.ContextWithJMESPath(context.Background(), "[?type=='PREDEFINED']")
	result, err := networkservices.GetAllNetworkServices(ctx, service, nil, nil)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, 6, result[0].ID)
	assert.Equal(t, 151, result[1].ID)
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
