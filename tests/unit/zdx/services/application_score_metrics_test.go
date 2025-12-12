// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
)

func TestApplicationScoreMetrics_Structure(t *testing.T) {
	t.Parallel()

	t.Run("Metric JSON marshaling", func(t *testing.T) {
		metric := common.Metric{
			Metric: "zdx_score",
			Unit:   "score",
			DataPoints: []common.DataPoint{
				{TimeStamp: 1699900000, Value: 85.5},
				{TimeStamp: 1699903600, Value: 88.2},
				{TimeStamp: 1699907200, Value: 91.0},
			},
		}

		data, err := json.Marshal(metric)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"metric":"zdx_score"`)
		assert.Contains(t, string(data), `"unit":"score"`)
		assert.Contains(t, string(data), `"datapoints"`)
	})

	t.Run("Metric JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"metric": "page_fetch_time",
			"unit": "ms",
			"datapoints": [
				{"timestamp": 1699900000, "value": 250.5},
				{"timestamp": 1699903600, "value": 245.2},
				{"timestamp": 1699907200, "value": 260.0}
			]
		}`

		var metric common.Metric
		err := json.Unmarshal([]byte(jsonData), &metric)
		require.NoError(t, err)

		assert.Equal(t, "page_fetch_time", metric.Metric)
		assert.Equal(t, "ms", metric.Unit)
		assert.Len(t, metric.DataPoints, 3)
		assert.Equal(t, 1699900000, metric.DataPoints[0].TimeStamp)
		assert.Equal(t, 250.5, metric.DataPoints[0].Value)
	})

	t.Run("DataPoint JSON marshaling", func(t *testing.T) {
		dataPoint := common.DataPoint{
			TimeStamp: 1699900000,
			Value:     95.75,
		}

		data, err := json.Marshal(dataPoint)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"timestamp":1699900000`)
		assert.Contains(t, string(data), `"value":95.75`)
	})

	t.Run("GetFromToFilters JSON marshaling", func(t *testing.T) {
		filters := common.GetFromToFilters{
			From:       1699900000,
			To:         1700000000,
			Loc:        []int{1, 2, 3},
			Dept:       []int{10, 20},
			Geo:        []string{"US-CA", "US-NY"},
			MetricName: "pft",
			Limit:      100,
			Offset:     "page2",
		}

		data, err := json.Marshal(filters)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"from":1699900000`)
		assert.Contains(t, string(data), `"to":1700000000`)
		assert.Contains(t, string(data), `"metric_name":"pft"`)
	})
}

func TestApplicationScoreMetrics_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse app score metrics response", func(t *testing.T) {
		jsonResponse := `[
			{
				"metric": "zdx_score",
				"unit": "score",
				"datapoints": [
					{"timestamp": 1699900000, "value": 85.0},
					{"timestamp": 1699903600, "value": 87.5},
					{"timestamp": 1699907200, "value": 90.0}
				]
			}
		]`

		var metrics []common.Metric
		err := json.Unmarshal([]byte(jsonResponse), &metrics)
		require.NoError(t, err)

		assert.Len(t, metrics, 1)
		assert.Equal(t, "zdx_score", metrics[0].Metric)
		assert.Len(t, metrics[0].DataPoints, 3)
	})

	t.Run("Parse app metrics response with multiple metrics", func(t *testing.T) {
		jsonResponse := `[
			{
				"metric": "page_fetch_time",
				"unit": "ms",
				"datapoints": [
					{"timestamp": 1699900000, "value": 250.0}
				]
			},
			{
				"metric": "server_response_time",
				"unit": "ms",
				"datapoints": [
					{"timestamp": 1699900000, "value": 150.0}
				]
			},
			{
				"metric": "dns_time",
				"unit": "ms",
				"datapoints": [
					{"timestamp": 1699900000, "value": 25.0}
				]
			}
		]`

		var metrics []common.Metric
		err := json.Unmarshal([]byte(jsonResponse), &metrics)
		require.NoError(t, err)

		assert.Len(t, metrics, 3)
		assert.Equal(t, "page_fetch_time", metrics[0].Metric)
		assert.Equal(t, "server_response_time", metrics[1].Metric)
		assert.Equal(t, "dns_time", metrics[2].Metric)
	})
}

