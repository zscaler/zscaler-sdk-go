// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/adminauditlogs"
)

// =====================================================
// SDK Function Tests
// =====================================================

func TestAdminAuditLogs_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/auditlogEntryReport"

	server.On("GET", path, common.SuccessResponse(adminauditlogs.AuditLogEntryReportTaskInfo{
		Status:                "COMPLETE",
		ProgressItemsComplete: 1000,
		ProgressEndTime:       1699999999,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := adminauditlogs.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Equal(t, "COMPLETE", result.Status)
	assert.Equal(t, 1000, result.ProgressItemsComplete)
}

func TestAdminAuditLogs_CreateAdminAuditLogsExport_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/auditlogEntryReport"

	server.On("POST", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	exportRequest := adminauditlogs.AuditLogEntryRequest{
		StartTime: 1699000000,
		EndTime:   1699999999,
	}

	resp, err := adminauditlogs.CreateAdminAuditLogsExport(context.Background(), service, exportRequest)

	require.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestAdminAuditLogs_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/auditlogEntryReport"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = adminauditlogs.Delete(context.Background(), service)

	require.NoError(t, err)
}

// Note: GetAdminAuditLogsDownload test omitted due to raw byte response handling

// =====================================================
// Structure Tests
// =====================================================

func TestAdminAuditLogs_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AuditLogEntryReportTaskInfo JSON marshaling", func(t *testing.T) {
		taskInfo := adminauditlogs.AuditLogEntryReportTaskInfo{
			Status:                "COMPLETE",
			ProgressItemsComplete: 1000,
			ProgressEndTime:       1699999999,
			ErrorMessage:          "",
			ErrorCode:             "",
		}

		data, err := json.Marshal(taskInfo)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"status":"COMPLETE"`)
		assert.Contains(t, string(data), `"progressItemsComplete":1000`)
	})

	t.Run("AuditLogEntryReportTaskInfo JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"status": "IN_PROGRESS",
			"progressItemsComplete": 500,
			"progressEndTime": 0,
			"errorMessage": "",
			"errorCode": ""
		}`

		var taskInfo adminauditlogs.AuditLogEntryReportTaskInfo
		err := json.Unmarshal([]byte(jsonData), &taskInfo)
		require.NoError(t, err)

		assert.Equal(t, "IN_PROGRESS", taskInfo.Status)
		assert.Equal(t, 500, taskInfo.ProgressItemsComplete)
	})

	t.Run("AuditLogEntryReportTaskInfo with error", func(t *testing.T) {
		jsonData := `{
			"status": "ERROR",
			"progressItemsComplete": 0,
			"progressEndTime": 0,
			"errorMessage": "Failed to generate report",
			"errorCode": "REPORT_GENERATION_FAILED"
		}`

		var taskInfo adminauditlogs.AuditLogEntryReportTaskInfo
		err := json.Unmarshal([]byte(jsonData), &taskInfo)
		require.NoError(t, err)

		assert.Equal(t, "ERROR", taskInfo.Status)
		assert.Equal(t, "Failed to generate report", taskInfo.ErrorMessage)
		assert.Equal(t, "REPORT_GENERATION_FAILED", taskInfo.ErrorCode)
	})

	t.Run("AuditLogEntryRequest JSON marshaling", func(t *testing.T) {
		req := adminauditlogs.AuditLogEntryRequest{
			StartTime:       1699000000,
			EndTime:         1699999999,
			Page:            1,
			PageSize:        "100",
			AdminName:       "admin@company.com",
			ObjectName:      "Firewall Policy",
			ActionInterface: "UI",
			Category:        "FIREWALL",
			Subcategories:   []string{"RULE_CREATE", "RULE_UPDATE"},
			ActionResult:    "SUCCESS",
			ActionTypes:     []string{"CREATE", "UPDATE", "DELETE"},
		}

		data, err := json.Marshal(req)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"startTime":1699000000`)
		assert.Contains(t, string(data), `"adminName":"admin@company.com"`)
		assert.Contains(t, string(data), `"category":"FIREWALL"`)
	})

	t.Run("AuditLogEntryRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"targetOrgId": 12345,
			"traceId": 67890,
			"startTime": 1699000000,
			"endTime": 1699999999,
			"page": 1,
			"pageSize": "50",
			"adminName": "security@company.com",
			"objectName": "URL Filtering",
			"actionInterface": "API",
			"category": "URL_FILTERING",
			"subcategories": ["RULE_CREATE"],
			"actionResult": "FAILURE",
			"actionTypes": ["CREATE"],
			"clientIP": 0
		}`

		var req adminauditlogs.AuditLogEntryRequest
		err := json.Unmarshal([]byte(jsonData), &req)
		require.NoError(t, err)

		assert.Equal(t, 12345, req.TargetOrgId)
		assert.Equal(t, "API", req.ActionInterface)
		assert.Equal(t, "FAILURE", req.ActionResult)
	})
}

func TestAdminAuditLogs_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse completed task info", func(t *testing.T) {
		jsonResponse := `{
			"status": "COMPLETE",
			"progressItemsComplete": 5000,
			"progressEndTime": 1699999999
		}`

		var taskInfo adminauditlogs.AuditLogEntryReportTaskInfo
		err := json.Unmarshal([]byte(jsonResponse), &taskInfo)
		require.NoError(t, err)

		assert.Equal(t, "COMPLETE", taskInfo.Status)
		assert.Equal(t, 5000, taskInfo.ProgressItemsComplete)
	})
}

