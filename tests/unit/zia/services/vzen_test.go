// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_clusters"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_nodes"
)

func TestVZENClusters_Structure(t *testing.T) {
	t.Parallel()

	t.Run("VZENClusters JSON marshaling", func(t *testing.T) {
		cluster := vzen_clusters.VZENClusters{
			ID:             12345,
			Name:           "VZEN Cluster 1",
			Status:         "ENABLED",
			IpAddress:      "192.168.1.100",
			SubnetMask:     "255.255.255.0",
			DefaultGateway: "192.168.1.1",
			Type:           "CLOUD_CONNECTOR",
			IpSecEnabled:   true,
			VirtualZenNodes: []common.IDNameExternalID{
				{ID: 1, Name: "Node 1"},
				{ID: 2, Name: "Node 2"},
			},
		}

		data, err := json.Marshal(cluster)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"ipAddress":"192.168.1.100"`)
		assert.Contains(t, string(data), `"ipSecEnabled":true`)
	})

	t.Run("VZENClusters JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Production VZEN Cluster",
			"status": "ENABLED",
			"ipAddress": "10.0.0.100",
			"subnetMask": "255.255.255.0",
			"defaultGateway": "10.0.0.1",
			"type": "BRANCH_CONNECTOR",
			"ipSecEnabled": false,
			"virtualZenNodes": [
				{"id": 10, "name": "Node A"},
				{"id": 20, "name": "Node B"},
				{"id": 30, "name": "Node C"}
			]
		}`

		var cluster vzen_clusters.VZENClusters
		err := json.Unmarshal([]byte(jsonData), &cluster)
		require.NoError(t, err)

		assert.Equal(t, 54321, cluster.ID)
		assert.False(t, cluster.IpSecEnabled)
		assert.Len(t, cluster.VirtualZenNodes, 3)
	})
}

func TestVZENNodes_Structure(t *testing.T) {
	t.Parallel()

	t.Run("VZENNodes JSON marshaling", func(t *testing.T) {
		node := vzen_nodes.VZENNodes{
			ID:                            12345,
			ZGatewayID:                    67890,
			Name:                          "VZEN Node 1",
			Status:                        "ENABLED",
			InProduction:                  true,
			IPAddress:                     "192.168.1.101",
			SubnetMask:                    "255.255.255.0",
			DefaultGateway:                "192.168.1.1",
			Type:                          "CLOUD_CONNECTOR",
			IPSecEnabled:                  true,
			OnDemandSupportTunnelEnabled:  true,
			EstablishSupportTunnelEnabled: false,
			DeploymentMode:                "CLUSTER",
			ClusterName:                   "Production Cluster",
			VzenSkuType:                   "LARGE",
		}

		data, err := json.Marshal(node)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"zgatewayId":67890`)
		assert.Contains(t, string(data), `"inProduction":true`)
		assert.Contains(t, string(data), `"vzenSkuType":"LARGE"`)
	})

	t.Run("VZENNodes JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"zgatewayId": 11111,
			"name": "Standalone VZEN Node",
			"status": "ENABLED",
			"inProduction": false,
			"ipAddress": "10.0.0.101",
			"subnetMask": "255.255.255.0",
			"defaultGateway": "10.0.0.1",
			"type": "BRANCH_CONNECTOR",
			"ipSecEnabled": false,
			"onDemandSupportTunnelEnabled": false,
			"establishSupportTunnelEnabled": true,
			"loadBalancerIpAddress": "10.0.0.200",
			"deploymentMode": "STANDALONE",
			"clusterName": "",
			"vzenSkuType": "MEDIUM"
		}`

		var node vzen_nodes.VZENNodes
		err := json.Unmarshal([]byte(jsonData), &node)
		require.NoError(t, err)

		assert.Equal(t, 54321, node.ID)
		assert.False(t, node.InProduction)
		assert.Equal(t, "STANDALONE", node.DeploymentMode)
		assert.Equal(t, "MEDIUM", node.VzenSkuType)
	})
}

func TestVZEN_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse VZEN clusters list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Cluster 1", "status": "ENABLED"},
			{"id": 2, "name": "Cluster 2", "status": "ENABLED"},
			{"id": 3, "name": "Cluster 3", "status": "DISABLED"}
		]`

		var clusters []vzen_clusters.VZENClusters
		err := json.Unmarshal([]byte(jsonResponse), &clusters)
		require.NoError(t, err)

		assert.Len(t, clusters, 3)
		assert.Equal(t, "DISABLED", clusters[2].Status)
	})

	t.Run("Parse VZEN nodes list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Node 1", "status": "ENABLED", "deploymentMode": "CLUSTER"},
			{"id": 2, "name": "Node 2", "status": "ENABLED", "deploymentMode": "STANDALONE"},
			{"id": 3, "name": "Node 3", "status": "DISABLED", "deploymentMode": "CLUSTER"}
		]`

		var nodes []vzen_nodes.VZENNodes
		err := json.Unmarshal([]byte(jsonResponse), &nodes)
		require.NoError(t, err)

		assert.Len(t, nodes, 3)
		assert.Equal(t, "STANDALONE", nodes[1].DeploymentMode)
	})
}

