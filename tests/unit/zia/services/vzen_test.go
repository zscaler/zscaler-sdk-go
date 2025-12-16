// Package services provides unit tests for ZIA services
package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/tests/unit/common"
	ziacommon "github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_clusters"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/vzen_nodes"
)

// =====================================================
// VZEN Clusters SDK Function Tests
// =====================================================

func TestVZENClusters_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	clusterID := 12345
	path := "/zia/api/v1/virtualZenClusters/12345"

	server.On("GET", path, common.SuccessResponse(vzen_clusters.VZENClusters{
		ID:             clusterID,
		Name:           "VZEN Cluster 1",
		Status:         "ENABLED",
		IpAddress:      "192.168.1.100",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "192.168.1.1",
		Type:           "CLOUD_CONNECTOR",
		IpSecEnabled:   true,
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_clusters.Get(context.Background(), service, clusterID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, clusterID, result.ID)
	assert.Equal(t, "VZEN Cluster 1", result.Name)
}

func TestVZENClusters_GetClusterByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	clusterName := "Production Cluster"
	path := "/zia/api/v1/virtualZenClusters"

	server.On("GET", path, common.SuccessResponse([]vzen_clusters.VZENClusters{
		{ID: 1, Name: "Other Cluster", Status: "ENABLED"},
		{ID: 2, Name: clusterName, Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_clusters.GetClusterByName(context.Background(), service, clusterName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, clusterName, result.Name)
}

func TestVZENClusters_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/virtualZenClusters"

	server.On("GET", path, common.SuccessResponse([]vzen_clusters.VZENClusters{
		{ID: 1, Name: "Cluster 1", Status: "ENABLED"},
		{ID: 2, Name: "Cluster 2", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_clusters.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestVZENClusters_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/virtualZenClusters"

	server.On("POST", path, common.SuccessResponse(vzen_clusters.VZENClusters{
		ID:             99999,
		Name:           "New Cluster",
		Status:         "ENABLED",
		IpAddress:      "10.0.0.100",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "10.0.0.1",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newCluster := &vzen_clusters.VZENClusters{
		Name:           "New Cluster",
		IpAddress:      "10.0.0.100",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "10.0.0.1",
	}

	result, _, err := vzen_clusters.Create(context.Background(), service, newCluster)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestVZENClusters_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	clusterID := 12345
	path := "/zia/api/v1/virtualZenClusters/12345"

	server.On("PUT", path, common.SuccessResponse(vzen_clusters.VZENClusters{
		ID:   clusterID,
		Name: "Updated Cluster",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateCluster := &vzen_clusters.VZENClusters{
		ID:   clusterID,
		Name: "Updated Cluster",
	}

	result, _, err := vzen_clusters.Update(context.Background(), service, clusterID, updateCluster)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Cluster", result.Name)
}

// =====================================================
// VZEN Nodes SDK Function Tests
// =====================================================

func TestVZENNodes_Get_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nodeID := 12345
	path := "/zia/api/v1/virtualZenNodes/12345"

	server.On("GET", path, common.SuccessResponse(vzen_nodes.VZENNodes{
		ID:             nodeID,
		Name:           "VZEN Node 1",
		Status:         "ENABLED",
		IPAddress:      "192.168.1.101",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "192.168.1.1",
		DeploymentMode: "CLUSTER",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_nodes.Get(context.Background(), service, nodeID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, nodeID, result.ID)
	assert.Equal(t, "VZEN Node 1", result.Name)
}

func TestVZENNodes_GetNodeByName_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nodeName := "Production Node"
	path := "/zia/api/v1/virtualZenNodes"

	server.On("GET", path, common.SuccessResponse([]vzen_nodes.VZENNodes{
		{ID: 1, Name: "Other Node", Status: "ENABLED"},
		{ID: 2, Name: nodeName, Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_nodes.GetNodeByName(context.Background(), service, nodeName)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, nodeName, result.Name)
}

func TestVZENNodes_GetAll_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/virtualZenNodes"

	server.On("GET", path, common.SuccessResponse([]vzen_nodes.VZENNodes{
		{ID: 1, Name: "Node 1", Status: "ENABLED"},
		{ID: 2, Name: "Node 2", Status: "ENABLED"},
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	result, err := vzen_nodes.GetAll(context.Background(), service)

	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestVZENNodes_Create_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	path := "/zia/api/v1/virtualZenNodes"

	server.On("POST", path, common.SuccessResponse(vzen_nodes.VZENNodes{
		ID:             99999,
		Name:           "New Node",
		Status:         "ENABLED",
		IPAddress:      "10.0.0.101",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "10.0.0.1",
		DeploymentMode: "STANDALONE",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	newNode := &vzen_nodes.VZENNodes{
		Name:           "New Node",
		IPAddress:      "10.0.0.101",
		SubnetMask:     "255.255.255.0",
		DefaultGateway: "10.0.0.1",
		DeploymentMode: "STANDALONE",
	}

	result, _, err := vzen_nodes.Create(context.Background(), service, newNode)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 99999, result.ID)
}

func TestVZENNodes_Update_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nodeID := 12345
	path := "/zia/api/v1/virtualZenNodes/12345"

	server.On("PUT", path, common.SuccessResponse(vzen_nodes.VZENNodes{
		ID:   nodeID,
		Name: "Updated Node",
	}))

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	updateNode := &vzen_nodes.VZENNodes{
		ID:   nodeID,
		Name: "Updated Node",
	}

	result, _, err := vzen_nodes.Update(context.Background(), service, nodeID, updateNode)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Updated Node", result.Name)
}

func TestVZENNodes_Delete_SDK(t *testing.T) {
	server := common.NewTestServer()
	defer server.Close()

	nodeID := 12345
	path := "/zia/api/v1/virtualZenNodes/12345"

	server.On("DELETE", path, common.NoContentResponse())

	service, err := common.CreateTestService(context.Background(), server, "123456")
	require.NoError(t, err)

	_, err = vzen_nodes.Delete(context.Background(), service, nodeID)

	require.NoError(t, err)
}

// =====================================================
// Structure Tests
// =====================================================

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
			VirtualZenNodes: []ziacommon.IDNameExternalID{
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

