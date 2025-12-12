// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/cloudnss/nss_servers"
)

func TestNSSServers_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NSSServers JSON marshaling", func(t *testing.T) {
		server := nss_servers.NSSServers{
			ID:        12345,
			Name:      "NSS Server 1",
			Status:    "UP",
			State:     "ENABLED",
			Type:      "NSS_FOR_WEB",
			IcapSvrId: 100,
		}

		data, err := json.Marshal(server)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"NSS_FOR_WEB"`)
		assert.Contains(t, string(data), `"status":"UP"`)
	})

	t.Run("NSSServers JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "NSS Server 2",
			"status": "DOWN",
			"state": "DISABLED",
			"type": "NSS_FOR_FIREWALL",
			"icapSvrId": 200
		}`

		var server nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonData), &server)
		require.NoError(t, err)

		assert.Equal(t, 54321, server.ID)
		assert.Equal(t, "DOWN", server.Status)
		assert.Equal(t, "NSS_FOR_FIREWALL", server.Type)
	})

	t.Run("NSSServers with VZEN type", func(t *testing.T) {
		jsonData := `{
			"id": 67890,
			"name": "VZEN Server",
			"status": "UP",
			"state": "ENABLED",
			"type": "VZEN"
		}`

		var server nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonData), &server)
		require.NoError(t, err)

		assert.Equal(t, "VZEN", server.Type)
	})
}

func TestNSSServers_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse NSS servers list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "Server 1", "status": "UP", "type": "NSS_FOR_WEB"},
			{"id": 2, "name": "Server 2", "status": "UP", "type": "NSS_FOR_FIREWALL"},
			{"id": 3, "name": "Server 3", "status": "DOWN", "type": "VZEN"}
		]`

		var servers []nss_servers.NSSServers
		err := json.Unmarshal([]byte(jsonResponse), &servers)
		require.NoError(t, err)

		assert.Len(t, servers, 3)
		assert.Equal(t, "DOWN", servers[2].Status)
	})
}

