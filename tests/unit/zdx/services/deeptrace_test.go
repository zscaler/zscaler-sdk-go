// Package services provides unit tests for ZDX deeptrace service
package services

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/troubleshooting/deeptrace"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDeepTrace_GetDeepTraces_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/deeptraces"

	server.On("GET", path, common.SuccessResponse([]deeptrace.DeepTraceSession{
		{TraceID: "trace-001", Status: "completed"},
		{TraceID: "trace-002", Status: "running"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := deeptrace.GetDeepTraces(context.Background(), service, 12345)

	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "trace-001", result[0].TraceID)
}

func TestDeepTrace_GetDeepTraceSession_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/deeptraces/trace-001"

	server.On("GET", path, common.SuccessResponse(nil))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	resp, err := deeptrace.GetDeepTraceSession(context.Background(), service, 12345, "trace-001")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeepTrace_CreateDeepTraceSession_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/deeptraces"

	server.On("POST", path, common.SuccessResponse(deeptrace.DeepTraceSession{
		TraceID: "new-trace-123",
		Status:  "pending",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	payload := deeptrace.DeepTraceSessionPayload{
		SessionName:          "Debug Session",
		AppID:                100,
		WebProbeID:           1,
		SessionLengthMinutes: 30,
		ProbeDevice:          true,
	}

	result, _, err := deeptrace.CreateDeepTraceSession(context.Background(), service, 12345, payload)

	require.NoError(t, err)
	assert.Equal(t, "new-trace-123", result.TraceID)
}

func TestDeepTrace_DeleteDeepTraceSession_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/deeptraces/trace-001"

	server.On("DELETE", path, common.SuccessResponse(nil))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	resp, err := deeptrace.DeleteDeepTraceSession(context.Background(), service, 12345, "trace-001")

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDeepTrace_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeepTraceSession JSON marshaling", func(t *testing.T) {
		session := deeptrace.DeepTraceSession{
			TraceID: "trace-abc-123",
			Status:  "completed",
			TraceDetails: deeptrace.TraceDetails{
				SessionName:        "Performance Debug Session",
				AppID:              "12345",
				AppName:            "Microsoft 365",
				UserID:             "user-001",
				Username:           "john.doe@example.com",
				DeviceID:           "device-001",
				DeviceName:         "LAPTOP-ENG-001",
				WebProbeID:         "probe-001",
				WebProbeName:       "office365.com",
				CloudPathProbeID:   "cpprobe-001",
				CloudPathProbeName: "Azure Endpoint",
				SessionLength:      30,
				ProbeDevice:        true,
			},
			CreatedAt:           1699900000,
			StartedAt:           1699900100,
			EndedAt:             1699901900,
			ExpectedTimeMinutes: 30,
		}

		data, err := json.Marshal(session)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"trace_id":"trace-abc-123"`)
		assert.Contains(t, string(data), `"status":"completed"`)
		assert.Contains(t, string(data), `"session_name":"Performance Debug Session"`)
		assert.Contains(t, string(data), `"probe_device":true`)
	})

	t.Run("DeepTraceSession JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"trace_id": "trace-xyz-789",
			"status": "running",
			"trace_details": {
				"session_name": "Network Issue Debug",
				"app_id": "67890",
				"app_name": "Salesforce",
				"user_id": "user-002",
				"username": "jane.smith@example.com",
				"device_id": "device-002",
				"device_name": "DESKTOP-SALES-001",
				"session_length": 15,
				"probe_device": false
			},
			"created_at": 1699950000,
			"started_at": 1699950100,
			"expected_time_minutes": 15
		}`

		var session deeptrace.DeepTraceSession
		err := json.Unmarshal([]byte(jsonData), &session)
		require.NoError(t, err)

		assert.Equal(t, "trace-xyz-789", session.TraceID)
		assert.Equal(t, "running", session.Status)
		assert.Equal(t, "Network Issue Debug", session.TraceDetails.SessionName)
		assert.Equal(t, "67890", session.TraceDetails.AppID)
		assert.Equal(t, "Salesforce", session.TraceDetails.AppName)
		assert.Equal(t, 15, session.TraceDetails.SessionLength)
		assert.False(t, session.TraceDetails.ProbeDevice)
	})

	t.Run("TraceDetails JSON marshaling", func(t *testing.T) {
		details := deeptrace.TraceDetails{
			SessionName:        "Zoom Call Debug",
			AppID:              "11111",
			AppName:            "Zoom",
			UserID:             "user-100",
			Username:           "admin@example.com",
			DeviceID:           "device-100",
			DeviceName:         "MACBOOK-ADMIN",
			WebProbeID:         "wp-001",
			WebProbeName:       "zoom.us",
			CloudPathProbeID:   "cp-001",
			CloudPathProbeName: "Zoom CDN",
			SessionLength:      60,
			ProbeDevice:        true,
		}

		data, err := json.Marshal(details)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"session_name":"Zoom Call Debug"`)
		assert.Contains(t, string(data), `"app_name":"Zoom"`)
		assert.Contains(t, string(data), `"web_probe_name":"zoom.us"`)
		assert.Contains(t, string(data), `"session_length":60`)
	})

	t.Run("DeepTraceSessionPayload JSON marshaling", func(t *testing.T) {
		payload := deeptrace.DeepTraceSessionPayload{
			SessionName:          "New Debug Session",
			AppID:                12345,
			WebProbeID:           100,
			CloudPathProbeID:     200,
			SessionLengthMinutes: 30,
			ProbeDevice:          true,
		}

		data, err := json.Marshal(payload)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"session_name":"New Debug Session"`)
		assert.Contains(t, string(data), `"app_id":12345`)
		assert.Contains(t, string(data), `"web_probe_id":100`)
		assert.Contains(t, string(data), `"cloud_path_probe_id":200`)
		assert.Contains(t, string(data), `"session_length_minutes":30`)
		assert.Contains(t, string(data), `"probe_device":true`)
	})

	t.Run("DeepTraceSessionPayload JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"session_name": "Teams Debug",
			"app_id": 54321,
			"web_probe_id": 500,
			"cloud_path_probe_id": 600,
			"session_length_minutes": 45,
			"probe_device": false
		}`

		var payload deeptrace.DeepTraceSessionPayload
		err := json.Unmarshal([]byte(jsonData), &payload)
		require.NoError(t, err)

		assert.Equal(t, "Teams Debug", payload.SessionName)
		assert.Equal(t, 54321, payload.AppID)
		assert.Equal(t, 500, payload.WebProbeID)
		assert.Equal(t, 600, payload.CloudPathProbeID)
		assert.Equal(t, 45, payload.SessionLengthMinutes)
		assert.False(t, payload.ProbeDevice)
	})
}

func TestDeepTrace_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse deep traces list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"trace_id": "trace-001",
				"status": "completed",
				"trace_details": {
					"session_name": "Session 1",
					"app_id": "100",
					"app_name": "App 1"
				},
				"created_at": 1699900000,
				"ended_at": 1699901800
			},
			{
				"trace_id": "trace-002",
				"status": "running",
				"trace_details": {
					"session_name": "Session 2",
					"app_id": "200",
					"app_name": "App 2"
				},
				"created_at": 1699950000,
				"started_at": 1699950100,
				"expected_time_minutes": 30
			},
			{
				"trace_id": "trace-003",
				"status": "pending",
				"trace_details": {
					"session_name": "Session 3",
					"app_id": "300",
					"app_name": "App 3"
				},
				"created_at": 1699960000
			}
		]`

		var sessions []deeptrace.DeepTraceSession
		err := json.Unmarshal([]byte(jsonResponse), &sessions)
		require.NoError(t, err)

		assert.Len(t, sessions, 3)
		
		// Check first session
		assert.Equal(t, "trace-001", sessions[0].TraceID)
		assert.Equal(t, "completed", sessions[0].Status)
		assert.Equal(t, "Session 1", sessions[0].TraceDetails.SessionName)
		
		// Check second session
		assert.Equal(t, "trace-002", sessions[1].TraceID)
		assert.Equal(t, "running", sessions[1].Status)
		assert.Equal(t, 30, sessions[1].ExpectedTimeMinutes)
		
		// Check third session
		assert.Equal(t, "trace-003", sessions[2].TraceID)
		assert.Equal(t, "pending", sessions[2].Status)
	})

	t.Run("Parse single deep trace session", func(t *testing.T) {
		jsonResponse := `{
			"trace_id": "trace-detailed-001",
			"status": "completed",
			"trace_details": {
				"session_name": "Full Debug Session",
				"app_id": "12345",
				"app_name": "Microsoft Teams",
				"user_id": "user-admin",
				"username": "admin@company.com",
				"device_id": "dev-laptop-001",
				"device_name": "Admin Laptop",
				"web_probe_id": "wp-teams",
				"web_probe_name": "teams.microsoft.com",
				"cloudpath_probe_id": "cp-azure",
				"cloud_path_name": "Azure CDN",
				"session_length": 30,
				"probe_device": true
			},
			"created_at": 1699900000,
			"started_at": 1699900060,
			"ended_at": 1699901860,
			"expected_time_minutes": 30
		}`

		var session deeptrace.DeepTraceSession
		err := json.Unmarshal([]byte(jsonResponse), &session)
		require.NoError(t, err)

		assert.Equal(t, "trace-detailed-001", session.TraceID)
		assert.Equal(t, "completed", session.Status)
		assert.Equal(t, "Full Debug Session", session.TraceDetails.SessionName)
		assert.Equal(t, "Microsoft Teams", session.TraceDetails.AppName)
		assert.Equal(t, "admin@company.com", session.TraceDetails.Username)
		assert.Equal(t, "teams.microsoft.com", session.TraceDetails.WebProbeName)
		assert.True(t, session.TraceDetails.ProbeDevice)
		assert.Equal(t, 1699901860, session.EndedAt)
	})

	t.Run("Parse deep trace with different statuses", func(t *testing.T) {
		statuses := []string{"pending", "running", "completed", "failed", "cancelled"}
		
		for _, status := range statuses {
			jsonData := `{"trace_id": "trace-1", "status": "` + status + `"}`
			
			var session deeptrace.DeepTraceSession
			err := json.Unmarshal([]byte(jsonData), &session)
			require.NoError(t, err)
			
			assert.Equal(t, status, session.Status)
		}
	})
}

