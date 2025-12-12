// Package unit provides unit tests for ZPA services
package unit

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

// =============================================================================
// SegmentGroup Structure Tests
// =============================================================================

func TestSegmentGroup_Structure(t *testing.T) {
	t.Parallel()

	t.Run("SegmentGroup JSON marshaling", func(t *testing.T) {
		sg := &segmentgroup.SegmentGroup{
			ID:                  "123456",
			Name:                "Test Segment Group",
			Description:         "Test description",
			Enabled:             true,
			ConfigSpace:         "DEFAULT",
			PolicyMigrated:      false,
			TcpKeepAliveEnabled: "1",
			Applications:        []segmentgroup.Application{},
		}

		data, err := json.Marshal(sg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"123456"`)
		assert.Contains(t, string(data), `"name":"Test Segment Group"`)
		assert.Contains(t, string(data), `"enabled":true`)
	})

	t.Run("SegmentGroup JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "789",
			"name": "Unmarshaled Group",
			"description": "From JSON",
			"enabled": true,
			"configSpace": "DEFAULT",
			"policyMigrated": false,
			"tcpKeepAliveEnabled": "0",
			"applications": []
		}`

		var sg segmentgroup.SegmentGroup
		err := json.Unmarshal([]byte(jsonData), &sg)
		require.NoError(t, err)

		assert.Equal(t, "789", sg.ID)
		assert.Equal(t, "Unmarshaled Group", sg.Name)
		assert.Equal(t, "From JSON", sg.Description)
		assert.True(t, sg.Enabled)
	})

	t.Run("Application structure", func(t *testing.T) {
		app := segmentgroup.Application{
			ID:                 "app-123",
			Name:               "Test App",
			Description:        "Test application",
			Enabled:            true,
			DomainNames:        []string{"example.com", "test.example.com"},
			DoubleEncrypt:      false,
			IPAnchored:         true,
			HealthCheckType:    "NONE",
			BypassType:         "NEVER",
			ConfigSpace:        "DEFAULT",
			DefaultIdleTimeout: "3600",
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"app-123"`)
		assert.Contains(t, string(data), `"name":"Test App"`)
		assert.Contains(t, string(data), `"enabled":true`)
		assert.Contains(t, string(data), `"ipAnchored":true`)
	})
}

// =============================================================================
// SegmentGroup Response Parsing Tests
// =============================================================================

func TestSegmentGroup_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse single segment group response", func(t *testing.T) {
		jsonResponse := `{
			"id": "12345",
			"name": "Production Segment Group",
			"description": "Production applications",
			"enabled": true,
			"configSpace": "DEFAULT",
			"creationTime": "2024-01-15T10:30:00Z",
			"modifiedBy": "admin@example.com",
			"modifiedTime": "2024-01-20T14:45:00Z",
			"policyMigrated": true,
			"tcpKeepAliveEnabled": "1",
			"microtenantId": "",
			"applications": [
				{
					"id": "app-001",
					"name": "Web App",
					"enabled": true,
					"domainNames": ["webapp.example.com"]
				}
			]
		}`

		var sg segmentgroup.SegmentGroup
		err := json.Unmarshal([]byte(jsonResponse), &sg)
		require.NoError(t, err)

		assert.Equal(t, "12345", sg.ID)
		assert.Equal(t, "Production Segment Group", sg.Name)
		assert.True(t, sg.Enabled)
		assert.True(t, sg.PolicyMigrated)
		assert.Len(t, sg.Applications, 1)
		assert.Equal(t, "Web App", sg.Applications[0].Name)
	})

	t.Run("Parse segment group list response", func(t *testing.T) {
		jsonResponse := `{
			"list": [
				{
					"id": "sg-001",
					"name": "Group 1",
					"enabled": true,
					"applications": []
				},
				{
					"id": "sg-002",
					"name": "Group 2",
					"enabled": false,
					"applications": []
				}
			],
			"totalPages": 1
		}`

		var response struct {
			List       []segmentgroup.SegmentGroup `json:"list"`
			TotalPages int                         `json:"totalPages"`
		}
		err := json.Unmarshal([]byte(jsonResponse), &response)
		require.NoError(t, err)

		assert.Len(t, response.List, 2)
		assert.Equal(t, "sg-001", response.List[0].ID)
		assert.Equal(t, "Group 1", response.List[0].Name)
		assert.True(t, response.List[0].Enabled)
		assert.Equal(t, "sg-002", response.List[1].ID)
		assert.False(t, response.List[1].Enabled)
	})
}

// =============================================================================
// Mock Server SegmentGroup Tests
// =============================================================================

func TestSegmentGroup_MockServerOperations(t *testing.T) {
	t.Run("GET segment group by ID", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		expectedResponse := common.MockSegmentGroupResponse("sg-123", "Test Group")
		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123", common.SuccessResponse(expectedResponse))

		// Simulate API call
		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var sg segmentgroup.SegmentGroup
		err = json.Unmarshal(body, &sg)
		require.NoError(t, err)

		assert.Equal(t, "sg-123", sg.ID)
		assert.Equal(t, "Test Group", sg.Name)
	})

	t.Run("GET all segment groups", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		groups := common.MockSegmentGroupListResponse(
			common.MockSegmentGroupResponse("sg-001", "Group 1"),
			common.MockSegmentGroupResponse("sg-002", "Group 2"),
		)
		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup", common.SuccessResponse(groups))

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response struct {
			List []segmentgroup.SegmentGroup `json:"list"`
		}
		err = json.Unmarshal(body, &response)
		require.NoError(t, err)

		assert.Len(t, response.List, 2)
	})

	t.Run("POST create segment group", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		createdGroup := common.MockSegmentGroupResponse("new-sg-123", "New Group")
		server.On("POST", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup", common.SuccessResponse(createdGroup))

		requestBody := `{"name": "New Group", "enabled": true}`
		resp, err := http.Post(
			server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup",
			"application/json",
			strings.NewReader(requestBody),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 200, resp.StatusCode)

		// Verify request was recorded
		lastReq := server.LastRequest()
		require.NotNil(t, lastReq)
		assert.Equal(t, "POST", lastReq.Method)
		assert.Contains(t, string(lastReq.Body), "New Group")
	})

	t.Run("PUT update segment group", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("PUT", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123", common.NoContentResponse())

		req, _ := http.NewRequest(
			"PUT",
			server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123",
			strings.NewReader(`{"name": "Updated Group", "enabled": false}`),
		)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 204, resp.StatusCode)
	})

	t.Run("DELETE segment group", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("DELETE", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123", common.NoContentResponse())

		req, _ := http.NewRequest(
			"DELETE",
			server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123",
			nil,
		)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 204, resp.StatusCode)
	})
}

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestSegmentGroup_ErrorHandling(t *testing.T) {
	t.Parallel()

	t.Run("404 Not Found", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/nonexistent", common.NotFoundResponse())

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/nonexistent")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("401 Unauthorized", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup", common.SessionInvalidResponse())

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 401, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "SESSION_NOT_VALID")
	})

	t.Run("409 Conflict - Edit Lock", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("PUT", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123", common.EditLockResponse())

		req, _ := http.NewRequest(
			"PUT",
			server.URL+"/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup/sg-123",
			strings.NewReader(`{"name": "Updated"}`),
		)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 409, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		assert.Contains(t, string(body), "EDIT_LOCK_NOT_AVAILABLE")
	})

	t.Run("429 Rate Limited", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup", common.TooManyRequestsResponseWithHeader("30"))

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 429, resp.StatusCode)
		assert.Equal(t, "30", resp.Header.Get("Retry-After"))
	})

	t.Run("500 Server Error", func(t *testing.T) {
		server := common.NewTestServer()
		defer server.Close()

		server.On("GET", "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup", common.ServerErrorResponse())

		resp, err := http.Get(server.URL + "/zpa/mgmtconfig/v1/admin/customers/123/segmentGroup")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, 500, resp.StatusCode)
	})
}

// =============================================================================
// MicroTenant Tests
// =============================================================================

func TestSegmentGroup_MicroTenant(t *testing.T) {
	t.Parallel()

	t.Run("SegmentGroup with MicroTenantID", func(t *testing.T) {
		sg := &segmentgroup.SegmentGroup{
			ID:            "sg-123",
			Name:          "Microtenant Group",
			Enabled:       true,
			MicroTenantID: "mt-456",
		}

		data, err := json.Marshal(sg)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"microtenantId":"mt-456"`)
	})

	t.Run("Parse response with MicroTenantName", func(t *testing.T) {
		jsonResponse := `{
			"id": "sg-789",
			"name": "MT Group",
			"enabled": true,
			"microtenantId": "mt-001",
			"microtenantName": "Test Microtenant",
			"applications": []
		}`

		var sg segmentgroup.SegmentGroup
		err := json.Unmarshal([]byte(jsonResponse), &sg)
		require.NoError(t, err)

		assert.Equal(t, "mt-001", sg.MicroTenantID)
		assert.Equal(t, "Test Microtenant", sg.MicroTenantName)
	})
}

