// Package services provides unit tests for ZDX services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zdx/services/reports/devices"
)

func TestDeviceHealthMetrics_Structure(t *testing.T) {
	t.Parallel()

	t.Run("HealthMetrics JSON marshaling", func(t *testing.T) {
		healthMetrics := devices.HealthMetrics{
			Category: "CPU",
			Instances: []devices.Instances{
				{
					Name: "cpu_usage",
					Metrics: []common.Metric{
						{
							Metric: "cpu_percent",
							Unit:   "%",
							DataPoints: []common.DataPoint{
								{TimeStamp: 1699900000, Value: 45.5},
								{TimeStamp: 1699903600, Value: 52.3},
							},
						},
					},
				},
			},
		}

		data, err := json.Marshal(healthMetrics)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"category":"CPU"`)
		assert.Contains(t, string(data), `"instances"`)
	})

	t.Run("HealthMetrics JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"category": "Memory",
			"instances": [
				{
					"metric": "memory_usage",
					"metrics": [
						{
							"metric": "mem_percent",
							"unit": "%",
							"datapoints": [
								{"timestamp": 1699900000, "value": 65.0},
								{"timestamp": 1699903600, "value": 68.5}
							]
						}
					]
				}
			]
		}`

		var healthMetrics devices.HealthMetrics
		err := json.Unmarshal([]byte(jsonData), &healthMetrics)
		require.NoError(t, err)

		assert.Equal(t, "Memory", healthMetrics.Category)
		assert.Len(t, healthMetrics.Instances, 1)
		assert.Len(t, healthMetrics.Instances[0].Metrics, 1)
		assert.Equal(t, "mem_percent", healthMetrics.Instances[0].Metrics[0].Metric)
	})

	t.Run("Instances JSON marshaling", func(t *testing.T) {
		instance := devices.Instances{
			Name: "disk_io",
			Metrics: []common.Metric{
				{
					Metric: "read_bytes",
					Unit:   "MB/s",
					DataPoints: []common.DataPoint{
						{TimeStamp: 1699900000, Value: 125.5},
					},
				},
				{
					Metric: "write_bytes",
					Unit:   "MB/s",
					DataPoints: []common.DataPoint{
						{TimeStamp: 1699900000, Value: 85.2},
					},
				},
			},
		}

		data, err := json.Marshal(instance)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"metric":"disk_io"`)
		assert.Contains(t, string(data), `"metrics"`)
	})
}

func TestDeviceHealthMetrics_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse health metrics list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"category": "CPU",
				"instances": [
					{
						"metric": "cpu_usage",
						"metrics": [
							{
								"metric": "cpu_percent",
								"unit": "%",
								"datapoints": [
									{"timestamp": 1699900000, "value": 45.0}
								]
							}
						]
					}
				]
			},
			{
				"category": "Memory",
				"instances": [
					{
						"metric": "memory_usage",
						"metrics": [
							{
								"metric": "mem_percent",
								"unit": "%",
								"datapoints": [
									{"timestamp": 1699900000, "value": 70.0}
								]
							}
						]
					}
				]
			},
			{
				"category": "Disk",
				"instances": [
					{
						"metric": "disk_io",
						"metrics": [
							{
								"metric": "disk_read",
								"unit": "MB/s",
								"datapoints": [
									{"timestamp": 1699900000, "value": 50.0}
								]
							}
						]
					}
				]
			}
		]`

		var healthMetrics []devices.HealthMetrics
		err := json.Unmarshal([]byte(jsonResponse), &healthMetrics)
		require.NoError(t, err)

		assert.Len(t, healthMetrics, 3)
		assert.Equal(t, "CPU", healthMetrics[0].Category)
		assert.Equal(t, "Memory", healthMetrics[1].Category)
		assert.Equal(t, "Disk", healthMetrics[2].Category)
	})
}

