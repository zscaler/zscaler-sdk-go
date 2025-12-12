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

func TestDeviceCloudPathProbes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("DeviceCloudPathProbe JSON marshaling", func(t *testing.T) {
		probe := devices.DeviceCloudPathProbe{
			ID:        12345,
			Name:      "azure-endpoint",
			NumProbes: 50,
			AverageLatency: []devices.AverageLatency{
				{LegSRC: "Client", LegDst: "Egress", Latency: 15.5},
				{LegSRC: "Egress", LegDst: "Application", Latency: 25.3},
			},
		}

		data, err := json.Marshal(probe)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"azure-endpoint"`)
		assert.Contains(t, string(data), `"num_probes":50`)
		assert.Contains(t, string(data), `"avg_latencies"`)
	})

	t.Run("DeviceCloudPathProbe JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "aws-endpoint",
			"num_probes": 100,
			"avg_latencies": [
				{"leg_src": "Client", "leg_dst": "ZIA", "latency": 10.5},
				{"leg_src": "ZIA", "leg_dst": "Application", "latency": 30.2}
			]
		}`

		var probe devices.DeviceCloudPathProbe
		err := json.Unmarshal([]byte(jsonData), &probe)
		require.NoError(t, err)

		assert.Equal(t, 67890, probe.ID)
		assert.Equal(t, "aws-endpoint", probe.Name)
		assert.Equal(t, 100, probe.NumProbes)
		assert.Len(t, probe.AverageLatency, 2)
		assert.Equal(t, float32(10.5), probe.AverageLatency[0].Latency)
	})

	t.Run("AverageLatency JSON marshaling", func(t *testing.T) {
		latency := devices.AverageLatency{
			LegSRC:  "Client",
			LegDst:  "Server",
			Latency: 45.7,
		}

		data, err := json.Marshal(latency)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"leg_src":"Client"`)
		assert.Contains(t, string(data), `"leg_dst":"Server"`)
		assert.Contains(t, string(data), `"latency":45.7`)
	})

	t.Run("NetworkStats JSON marshaling", func(t *testing.T) {
		stats := devices.NetworkStats{
			LegSRC: "Client",
			LegDst: "Egress",
			Stats: []common.Metric{
				{
					Metric: "latency",
					Unit:   "ms",
					DataPoints: []common.DataPoint{
						{TimeStamp: 1699900000, Value: 25.5},
					},
				},
			},
		}

		data, err := json.Marshal(stats)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"leg_src":"Client"`)
		assert.Contains(t, string(data), `"leg_dst":"Egress"`)
		assert.Contains(t, string(data), `"stats"`)
	})

	t.Run("CloudPathProbe JSON marshaling", func(t *testing.T) {
		probe := devices.CloudPathProbe{
			TimeStamp: 1699900000,
			CloudPath: []devices.CloudPath{
				{
					SRC:           "192.168.1.100",
					DST:           "40.97.100.1",
					NumHops:       5,
					Latency:       35.5,
					Loss:          0.5,
					NumUnrespHops: 1,
					TunnelType:    1,
					Hops: []devices.Hops{
						{IP: "192.168.1.1", LatencyAvg: 2, PktSent: 10, PktRcvd: 10},
						{IP: "10.0.0.1", LatencyAvg: 5, PktSent: 10, PktRcvd: 9},
					},
				},
			},
		}

		data, err := json.Marshal(probe)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"timestamp":1699900000`)
		assert.Contains(t, string(data), `"cloudpath"`)
		assert.Contains(t, string(data), `"num_hops":5`)
	})

	t.Run("CloudPath JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"src": "10.0.0.100",
			"dst": "52.96.0.1",
			"num_hops": 8,
			"latency": 45.5,
			"loss": 1.5,
			"num_unresp_hops": 2,
			"tunnel_type": 2,
			"hops": [
				{"ip": "10.0.0.1", "latency_avg": 1, "pkt_sent": 10, "pkt_rcvd": 10},
				{"ip": "10.0.1.1", "latency_avg": 3, "pkt_sent": 10, "pkt_rcvd": 10}
			]
		}`

		var cloudPath devices.CloudPath
		err := json.Unmarshal([]byte(jsonData), &cloudPath)
		require.NoError(t, err)

		assert.Equal(t, "10.0.0.100", cloudPath.SRC)
		assert.Equal(t, "52.96.0.1", cloudPath.DST)
		assert.Equal(t, 8, cloudPath.NumHops)
		assert.Equal(t, float32(45.5), cloudPath.Latency)
		assert.Len(t, cloudPath.Hops, 2)
	})

	t.Run("Hops JSON marshaling", func(t *testing.T) {
		hop := devices.Hops{
			IP:          "192.168.1.1",
			GWMac:       "00:11:22:33:44:55",
			GWMacVendor: "Cisco",
			PktSent:     100,
			PktRcvd:     98,
			LatencyMin:  1,
			LatencyMax:  10,
			LatencyAvg:  5,
			LatencyDiff: 9,
		}

		data, err := json.Marshal(hop)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"ip":"192.168.1.1"`)
		assert.Contains(t, string(data), `"gw_mac":"00:11:22:33:44:55"`)
		assert.Contains(t, string(data), `"pkt_sent":100`)
		assert.Contains(t, string(data), `"latency_avg":5`)
	})
}

func TestDeviceCloudPathProbes_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse cloudpath probes list response", func(t *testing.T) {
		jsonResponse := `[
			{
				"id": 1,
				"name": "probe-1",
				"num_probes": 50,
				"avg_latencies": [
					{"leg_src": "Client", "leg_dst": "Egress", "latency": 10.0}
				]
			},
			{
				"id": 2,
				"name": "probe-2",
				"num_probes": 40,
				"avg_latencies": [
					{"leg_src": "Client", "leg_dst": "Application", "latency": 50.0}
				]
			}
		]`

		var probes []devices.DeviceCloudPathProbe
		err := json.Unmarshal([]byte(jsonResponse), &probes)
		require.NoError(t, err)

		assert.Len(t, probes, 2)
		assert.Equal(t, "probe-1", probes[0].Name)
		assert.Equal(t, "probe-2", probes[1].Name)
	})

	t.Run("Parse cloudpath response", func(t *testing.T) {
		jsonResponse := `[
			{
				"timestamp": 1699900000,
				"cloudpath": [
					{
						"src": "192.168.1.100",
						"dst": "40.97.100.1",
						"num_hops": 6,
						"latency": 40.0,
						"loss": 0.0,
						"hops": [
							{"ip": "192.168.1.1", "latency_avg": 2},
							{"ip": "10.0.0.1", "latency_avg": 5},
							{"ip": "40.97.100.1", "latency_avg": 33}
						]
					}
				]
			}
		]`

		var cloudPaths []devices.CloudPathProbe
		err := json.Unmarshal([]byte(jsonResponse), &cloudPaths)
		require.NoError(t, err)

		assert.Len(t, cloudPaths, 1)
		assert.Equal(t, 1699900000, cloudPaths[0].TimeStamp)
		assert.Len(t, cloudPaths[0].CloudPath, 1)
		assert.Equal(t, 6, cloudPaths[0].CloudPath[0].NumHops)
		assert.Len(t, cloudPaths[0].CloudPath[0].Hops, 3)
	})
}

