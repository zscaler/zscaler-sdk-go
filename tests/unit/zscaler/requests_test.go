// Package zscaler provides unit tests for core zscaler SDK request functions
package zscaler

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
)

// =====================================================
// ZDX Request Tests
// =====================================================

func TestNewZdxRequestDo(t *testing.T) {
	t.Run("GET request with response parsing", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		// Mock a ZDX endpoint
		server.On("GET", "/zdx/v1/devices", common.SuccessResponse(map[string]interface{}{
			"devices": []map[string]interface{}{
				{"id": "device-1", "name": "Device 1"},
				{"id": "device-2", "name": "Device 2"},
			},
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result map[string]interface{}
		resp, err := service.Client.NewZdxRequestDo(context.Background(), "GET", "/zdx/v1/devices", nil, nil, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, result["devices"])
	})

	t.Run("GET request with query parameters", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		type QueryParams struct {
			Limit int `url:"limit"`
			From  int `url:"from"`
		}

		server.On("GET", "/zdx/v1/devices", common.SuccessResponse(map[string]interface{}{
			"count": 10,
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		params := QueryParams{Limit: 100, From: 1699900000}
		var result map[string]interface{}
		resp, err := service.Client.NewZdxRequestDo(context.Background(), "GET", "/zdx/v1/devices", params, nil, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("POST request with body", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/zdx/v1/analysis", common.SuccessResponse(map[string]interface{}{
			"id":     "analysis-123",
			"status": "pending",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		requestBody := map[string]interface{}{
			"device_id": "device-123",
			"type":      "network",
		}

		var result map[string]interface{}
		resp, err := service.Client.NewZdxRequestDo(context.Background(), "POST", "/zdx/v1/analysis", nil, requestBody, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "analysis-123", result["id"])
	})

	t.Run("Request with nil response target", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("DELETE", "/zdx/v1/analysis/123", common.SuccessResponse(nil))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		// nil response target - should not try to decode
		resp, err := service.Client.NewZdxRequestDo(context.Background(), "DELETE", "/zdx/v1/analysis/123", nil, nil, nil)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

// =====================================================
// ZIA Request Tests (uses Create/Read/Update/Delete methods)
// =====================================================

func TestZiaRequests(t *testing.T) {
	t.Run("Read - GET request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zia/api/v1/ruleLabels/12345", common.SuccessResponse(map[string]interface{}{
			"id":   12345,
			"name": "Test Label",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result map[string]interface{}
		err = service.Client.Read(context.Background(), "/zia/api/v1/ruleLabels/12345", &result)

		require.NoError(t, err)
		assert.Equal(t, float64(12345), result["id"])
	})

	t.Run("Create - POST request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/zia/api/v1/ruleLabels", common.SuccessResponse(map[string]interface{}{
			"ID":   12345,
			"Name": "Test Label",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		type RuleLabel struct {
			ID          int    `json:"id,omitempty"`
			Name        string `json:"name"`
			Description string `json:"description,omitempty"`
		}

		requestBody := RuleLabel{Name: "Test Label", Description: "Test description"}
		result, err := service.Client.Create(context.Background(), "/zia/api/v1/ruleLabels", requestBody)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("Update - PATCH request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		// Update uses PATCH method
		server.On("PATCH", "/zia/api/v1/ruleLabels/12345", common.SuccessResponse(map[string]interface{}{
			"ID":   12345,
			"Name": "Updated Label",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		type RuleLabel struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name"`
		}

		requestBody := RuleLabel{ID: 12345, Name: "Updated Label"}
		result, err := service.Client.Update(context.Background(), "/zia/api/v1/ruleLabels/12345", requestBody)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("UpdateWithPut - PUT request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		// UpdateWithPut uses PUT method
		server.On("PUT", "/zia/api/v1/ruleLabels/12345", common.SuccessResponse(map[string]interface{}{
			"ID":   12345,
			"Name": "Updated Label",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		type RuleLabel struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name"`
		}

		requestBody := RuleLabel{ID: 12345, Name: "Updated Label"}
		result, err := service.Client.UpdateWithPut(context.Background(), "/zia/api/v1/ruleLabels/12345", requestBody)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("Delete - DELETE request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("DELETE", "/zia/api/v1/ruleLabels/12345", common.SuccessResponseWithStatus(http.StatusNoContent, nil))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		err = service.Client.Delete(context.Background(), "/zia/api/v1/ruleLabels/12345")

		require.NoError(t, err)
	})

	t.Run("BulkDelete request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/zia/api/v1/ruleLabels/bulkDelete", common.SuccessResponse(nil))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		ids := []int{12345, 12346, 12347}
		_, err = service.Client.BulkDelete(context.Background(), "/zia/api/v1/ruleLabels/bulkDelete", ids)

		require.NoError(t, err)
	})
}

// =====================================================
// ZPA Request Tests (uses NewRequestDo)
// =====================================================

func TestZpaRequests(t *testing.T) {
	t.Run("GET request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123456/application", common.SuccessResponse(map[string]interface{}{
			"list": []map[string]interface{}{
				{"id": "app-1", "name": "App 1"},
			},
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result map[string]interface{}
		resp, err := service.Client.NewRequestDo(context.Background(), "GET", "/zpa/mgmtconfig/v1/admin/customers/123456/application", nil, nil, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotNil(t, result["list"])
	})

	t.Run("POST request with body", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/zpa/mgmtconfig/v1/admin/customers/123456/application", common.SuccessResponse(map[string]interface{}{
			"id":   "new-app-123",
			"name": "New Application",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		requestBody := map[string]interface{}{
			"name":        "New Application",
			"domainNames": []string{"example.com"},
		}

		var result map[string]interface{}
		resp, err := service.Client.NewRequestDo(context.Background(), "POST", "/zpa/mgmtconfig/v1/admin/customers/123456/application", nil, requestBody, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

// =====================================================
// ZCC Request Tests
// =====================================================

func TestNewZccRequestDo(t *testing.T) {
	t.Run("GET request", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zcc/papi/public/v1/getDevices", common.SuccessResponse(map[string]interface{}{
			"devices": []map[string]interface{}{
				{"udid": "device-1", "machineHostname": "host-1"},
			},
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result map[string]interface{}
		resp, err := service.Client.NewZccRequestDo(context.Background(), "GET", "/zcc/papi/public/v1/getDevices", nil, nil, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	t.Run("PUT request for update", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("PUT", "/zcc/papi/public/v1/webPolicy", common.SuccessResponse(map[string]interface{}{
			"id":     "policy-1",
			"status": "updated",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		requestBody := map[string]interface{}{
			"policyType": "web",
			"enabled":    true,
		}

		var result map[string]interface{}
		resp, err := service.Client.NewZccRequestDo(context.Background(), "PUT", "/zcc/papi/public/v1/webPolicy", nil, requestBody, &result)

		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

// =====================================================
// ZTW Resource Request Tests
// =====================================================

func TestZtwResourceRequests(t *testing.T) {
	t.Run("ReadResource", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/ztw/api/v1/ecgroup/123", common.SuccessResponse(map[string]interface{}{
			"id":   123,
			"name": "EC Group 1",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result map[string]interface{}
		err = service.Client.ReadResource(context.Background(), "/ztw/api/v1/ecgroup/123", &result)

		require.NoError(t, err)
		assert.Equal(t, float64(123), result["id"])
		assert.Equal(t, "EC Group 1", result["name"])
	})

	t.Run("CreateResource", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("POST", "/ztw/api/v1/ecgroup", common.SuccessResponseWithStatus(http.StatusCreated, map[string]interface{}{
			"ID":   456,
			"Name": "New EC Group",
		}))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		type ECGroup struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name"`
		}

		requestBody := ECGroup{Name: "New EC Group"}
		result, err := service.Client.CreateResource(context.Background(), "/ztw/api/v1/ecgroup", requestBody)

		require.NoError(t, err)
		require.NotNil(t, result)
	})

	t.Run("DeleteResource", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("DELETE", "/ztw/api/v1/ecgroup/789", common.SuccessResponseWithStatus(http.StatusNoContent, nil))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		err = service.Client.DeleteResource(context.Background(), "/ztw/api/v1/ecgroup/789")

		require.NoError(t, err)
	})

	t.Run("ReadTextResource", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/ztw/api/v1/ecgroup/status", common.RawResponse([]byte("active"), http.StatusOK, nil))

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		var result string
		err = service.Client.ReadTextResource(context.Background(), "/ztw/api/v1/ecgroup/status", &result)

		require.NoError(t, err)
		assert.Equal(t, "active", result)
	})
}

// =====================================================
// Client Helper Methods Tests
// =====================================================

func TestClientHelpers(t *testing.T) {
	t.Run("GetSandboxURL", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		// GetSandboxURL returns a URL (may or may not be empty based on config)
		url := service.Client.GetSandboxURL()
		// Just verify the function doesn't panic - the URL depends on cloud config
		assert.NotNil(t, url) // url is a string, so this checks it's callable
	})

	t.Run("GetSandboxToken", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		token := service.Client.GetSandboxToken()
		assert.Empty(t, token)
	})

	t.Run("GetLogger", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		service, err := common.CreateTestService(context.Background(), server, "123456")
		require.NoError(t, err)

		logger := service.Client.GetLogger()
		assert.NotNil(t, logger)
	})
}

