// Package services provides unit tests for ZDX device quality metrics service
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	zdxcommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

// =====================================================
// SDK Function Tests - Exercise actual SDK code paths
// =====================================================

func TestDeviceQualityMetrics_GetQualityMetrics_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zdx/v1/devices/12345/apps/100/call-quality-metrics"

	server.On("GET", path, common.SuccessResponse([]devices.CallQualityMetrics{
		{
			MeetID:        "meet-123",
			MeetSessionID: "session-456",
			MeetSubject:   "Team Standup",
			Metrics: []zdxcommon.Metric{
				{Metric: "jitter", Unit: "ms"},
				{Metric: "packet_loss", Unit: "%"},
			},
		},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, _, err := devices.GetQualityMetrics(context.Background(), service, 12345, 100, zdxcommon.GetFromToFilters{})

	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "meet-123", result[0].MeetID)
	assert.Len(t, result[0].Metrics, 2)
}

// =====================================================
// Structure Tests - JSON marshaling/unmarshaling
// =====================================================

func TestDeviceQualityMetrics_Structure(t *testing.T) {
	t.Parallel()

	t.Run("CallQualityMetrics JSON marshaling", func(t *testing.T) {
		cqm := devices.CallQualityMetrics{
			MeetID:        "meet-123-abc",
			MeetSessionID: "session-456",
			MeetSubject:   "Team Weekly Standup",
			Metrics: []zdxcommon.Metric{
				{
					Metric: "jitter",
					Unit:   "ms",
					DataPoints: []zdxcommon.DataPoint{
						{TimeStamp: 1699900000, Value: 5.5},
						{TimeStamp: 1699903600, Value: 7.2},
					},
				},
				{
					Metric: "packet_loss",
					Unit:   "%",
					DataPoints: []zdxcommon.DataPoint{
						{TimeStamp: 1699900000, Value: 0.5},
						{TimeStamp: 1699903600, Value: 1.2},
					},
				},
			},
		}

		data, err := json.Marshal(cqm)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"meet_id":"meet-123-abc"`)
		assert.Contains(t, string(data), `"meet_session_id":"session-456"`)
		assert.Contains(t, string(data), `"meet_subject":"Team Weekly Standup"`)
		assert.Contains(t, string(data), `"metrics"`)
	})

	t.Run("CallQualityMetrics JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"meet_id": "meet-789-xyz",
			"meet_session_id": "session-012",
			"meet_subject": "Customer Demo",
			"metrics": [
				{
					"metric": "audio_quality",
					"unit": "score",
					"datapoints": [
						{"timestamp": 1699900000, "value": 4.5}
					]
				},
				{
					"metric": "video_quality",
					"unit": "score",
					"datapoints": [
						{"timestamp": 1699900000, "value": 4.2}
					]
				}
			]
		}`

		var cqm devices.CallQualityMetrics
		err := json.Unmarshal([]byte(jsonData), &cqm)
		require.NoError(t, err)

		assert.Equal(t, "meet-789-xyz", cqm.MeetID)
		assert.Equal(t, "session-012", cqm.MeetSessionID)
		assert.Equal(t, "Customer Demo", cqm.MeetSubject)
		assert.Len(t, cqm.Metrics, 2)
		assert.Equal(t, "audio_quality", cqm.Metrics[0].Metric)
		assert.Equal(t, "video_quality", cqm.Metrics[1].Metric)
	})
}

func TestDeviceQualityMetrics_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse call quality metrics list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"meet_id": "meet-001",
				"meet_session_id": "sess-001",
				"meet_subject": "Morning Standup",
				"metrics": [
					{
						"metric": "jitter",
						"unit": "ms",
						"datapoints": [
							{"timestamp": 1699900000, "value": 3.5}
						]
					}
				]
			},
			{
				"meet_id": "meet-002",
				"meet_session_id": "sess-002",
				"meet_subject": "All Hands Meeting",
				"metrics": [
					{
						"metric": "packet_loss",
						"unit": "%",
						"datapoints": [
							{"timestamp": 1699950000, "value": 0.8}
						]
					}
				]
			}
		]`

		var metrics []devices.CallQualityMetrics
		err := json.Unmarshal([]byte(jsonResponse), &metrics)
		require.NoError(t, err)

		assert.Len(t, metrics, 2)
		assert.Equal(t, "meet-001", metrics[0].MeetID)
		assert.Equal(t, "Morning Standup", metrics[0].MeetSubject)
		assert.Equal(t, "meet-002", metrics[1].MeetID)
		assert.Equal(t, "All Hands Meeting", metrics[1].MeetSubject)
	})

	t.Run("Parse Zoom/Teams call quality", func(t *testing.T) {
		jsonResponse := `[
			{
				"meet_id": "zoom-abc123",
				"meet_session_id": "zoom-sess-1",
				"meet_subject": "Product Review",
				"metrics": [
					{
						"metric": "audio_jitter",
						"unit": "ms",
						"datapoints": [
							{"timestamp": 1699900000, "value": 2.5},
							{"timestamp": 1699903600, "value": 3.0},
							{"timestamp": 1699907200, "value": 2.8}
						]
					},
					{
						"metric": "video_latency",
						"unit": "ms",
						"datapoints": [
							{"timestamp": 1699900000, "value": 45.0},
							{"timestamp": 1699903600, "value": 52.0},
							{"timestamp": 1699907200, "value": 48.0}
						]
					},
					{
						"metric": "audio_packet_loss",
						"unit": "%",
						"datapoints": [
							{"timestamp": 1699900000, "value": 0.1},
							{"timestamp": 1699903600, "value": 0.2},
							{"timestamp": 1699907200, "value": 0.1}
						]
					}
				]
			}
		]`

		var metrics []devices.CallQualityMetrics
		err := json.Unmarshal([]byte(jsonResponse), &metrics)
		require.NoError(t, err)

		assert.Len(t, metrics, 1)
		assert.Equal(t, "zoom-abc123", metrics[0].MeetID)
		assert.Equal(t, "Product Review", metrics[0].MeetSubject)
		assert.Len(t, metrics[0].Metrics, 3)
		
		// Verify each metric type
		metricNames := make([]string, len(metrics[0].Metrics))
		for i, m := range metrics[0].Metrics {
			metricNames[i] = m.Metric
		}
		assert.Contains(t, metricNames, "audio_jitter")
		assert.Contains(t, metricNames, "video_latency")
		assert.Contains(t, metricNames, "audio_packet_loss")
	})
}

