// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/troubleshooting/analysis"
)

func TestAnalysis_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AnalysisRequest JSON marshaling", func(t *testing.T) {
		request := analysis.AnalysisRequest{
			DeviceID: 12345,
			AppID:    67890,
			T0:       1699900000,
			T1:       1700000000,
		}

		data, err := json.Marshal(request)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"device_id":12345`)
		assert.Contains(t, string(data), `"app_id":67890`)
		assert.Contains(t, string(data), `"t0":1699900000`)
		assert.Contains(t, string(data), `"t1":1700000000`)
	})

	t.Run("AnalysisRequest JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"device_id": 11111,
			"app_id": 22222,
			"t0": 1699800000,
			"t1": 1699900000
		}`

		var request analysis.AnalysisRequest
		err := json.Unmarshal([]byte(jsonData), &request)
		require.NoError(t, err)

		assert.Equal(t, 11111, request.DeviceID)
		assert.Equal(t, 22222, request.AppID)
		assert.Equal(t, 1699800000, request.T0)
		assert.Equal(t, 1699900000, request.T1)
	})

	t.Run("AnalysisResult JSON marshaling", func(t *testing.T) {
		result := analysis.AnalysisResult{
			ErrMsg: "",
			Result: analysis.Result{
				Issue:      "Network Latency",
				Confidence: 85,
				Message:    "High network latency detected between client and server",
				Times:      []int{1699900000, 1699903600, 1699907200},
			},
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"issue":"Network Latency"`)
		assert.Contains(t, string(data), `"confidence":85`)
		assert.Contains(t, string(data), `"message":"High network latency detected between client and server"`)
	})

	t.Run("AnalysisResult JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"err_msg": "",
			"result": {
				"issue": "CPU Throttling",
				"confidence": 92,
				"message": "CPU thermal throttling detected causing performance degradation",
				"times": [1699900000, 1699901000, 1699902000]
			}
		}`

		var result analysis.AnalysisResult
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Empty(t, result.ErrMsg)
		assert.Equal(t, "CPU Throttling", result.Result.Issue)
		assert.Equal(t, 92, result.Result.Confidence)
		assert.Len(t, result.Result.Times, 3)
	})

	t.Run("AnalysisResult with error message", func(t *testing.T) {
		jsonData := `{
			"err_msg": "Analysis failed: insufficient data",
			"result": {
				"issue": "",
				"confidence": 0,
				"message": "",
				"times": []
			}
		}`

		var result analysis.AnalysisResult
		err := json.Unmarshal([]byte(jsonData), &result)
		require.NoError(t, err)

		assert.Equal(t, "Analysis failed: insufficient data", result.ErrMsg)
		assert.Empty(t, result.Result.Issue)
		assert.Equal(t, 0, result.Result.Confidence)
	})

	t.Run("Result JSON marshaling", func(t *testing.T) {
		result := analysis.Result{
			Issue:      "Disk I/O Bottleneck",
			Confidence: 78,
			Message:    "High disk I/O wait times detected",
			Times:      []int{1699900000},
		}

		data, err := json.Marshal(result)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"issue":"Disk I/O Bottleneck"`)
		assert.Contains(t, string(data), `"confidence":78`)
		assert.Contains(t, string(data), `"times":[1699900000]`)
	})
}

func TestAnalysis_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse successful analysis response", func(t *testing.T) {
		jsonResponse := `{
			"err_msg": "",
			"result": {
				"issue": "WiFi Signal Strength",
				"confidence": 88,
				"message": "Weak WiFi signal detected causing packet loss and retransmissions",
				"times": [1699900000, 1699900300, 1699900600, 1699900900]
			}
		}`

		var result analysis.AnalysisResult
		err := json.Unmarshal([]byte(jsonResponse), &result)
		require.NoError(t, err)

		assert.Empty(t, result.ErrMsg)
		assert.Equal(t, "WiFi Signal Strength", result.Result.Issue)
		assert.Equal(t, 88, result.Result.Confidence)
		assert.Len(t, result.Result.Times, 4)
	})

	t.Run("Parse analysis with no issues found", func(t *testing.T) {
		jsonResponse := `{
			"err_msg": "",
			"result": {
				"issue": "No Issues",
				"confidence": 100,
				"message": "No performance issues detected during the analysis period",
				"times": []
			}
		}`

		var result analysis.AnalysisResult
		err := json.Unmarshal([]byte(jsonResponse), &result)
		require.NoError(t, err)

		assert.Equal(t, "No Issues", result.Result.Issue)
		assert.Equal(t, 100, result.Result.Confidence)
		assert.Empty(t, result.Result.Times)
	})

	t.Run("Parse multiple issue timestamps", func(t *testing.T) {
		jsonResponse := `{
			"err_msg": "",
			"result": {
				"issue": "Intermittent Connectivity",
				"confidence": 75,
				"message": "Multiple connection drops detected",
				"times": [1699900000, 1699910000, 1699920000, 1699930000, 1699940000, 1699950000]
			}
		}`

		var result analysis.AnalysisResult
		err := json.Unmarshal([]byte(jsonResponse), &result)
		require.NoError(t, err)

		assert.Equal(t, "Intermittent Connectivity", result.Result.Issue)
		assert.Len(t, result.Result.Times, 6)
		assert.Equal(t, 1699900000, result.Result.Times[0])
		assert.Equal(t, 1699950000, result.Result.Times[5])
	})
}

