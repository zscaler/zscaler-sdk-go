// Package services provides unit tests for ZIA services
package services

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplicationgroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkapplications"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/firewallpolicies/networkservices"
)

func TestNetworkServices_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkServices JSON marshaling", func(t *testing.T) {
		svc := networkservices.NetworkServices{
			ID:            12345,
			Name:          "Custom HTTPS",
			Description:   "Custom HTTPS on alternate port",
			Type:          "CUSTOM",
			Tag:           "HTTPS_ALT",
			IsNameL10nTag: false,
		}

		data, err := json.Marshal(svc)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"type":"CUSTOM"`)
	})

	t.Run("NetworkServices JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "HTTP",
			"description": "Hypertext Transfer Protocol",
			"type": "PREDEFINED",
			"tag": "HTTP",
			"isNameL10nTag": true,
			"srcTcpPorts": [
				{"start": 1, "end": 65535}
			],
			"destTcpPorts": [
				{"start": 80, "end": 80}
			]
		}`

		var svc networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonData), &svc)
		require.NoError(t, err)

		assert.Equal(t, 54321, svc.ID)
		assert.True(t, svc.IsNameL10nTag)
	})

	t.Run("NetworkPorts JSON marshaling", func(t *testing.T) {
		ports := networkservices.NetworkPorts{
			Start: 443,
			End:   443,
		}

		data, err := json.Marshal(ports)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"start":443`)
		assert.Contains(t, string(data), `"end":443`)
	})
}

func TestNetworkServiceGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkServiceGroups JSON marshaling", func(t *testing.T) {
		group := networkservicegroups.NetworkServiceGroups{
			ID:          12345,
			Name:        "Web Services",
			Description: "HTTP and HTTPS services",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
		assert.Contains(t, string(data), `"name":"Web Services"`)
	})

	t.Run("NetworkServiceGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "Database Services",
			"description": "Database protocols",
			"services": [
				{"id": 100, "name": "MySQL"},
				{"id": 101, "name": "PostgreSQL"}
			]
		}`

		var group networkservicegroups.NetworkServiceGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Len(t, group.Services, 2)
	})
}

func TestNetworkApplications_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkApplications JSON marshaling", func(t *testing.T) {
		app := networkapplications.NetworkApplications{
			ID:          "SSH",
			Description: "Secure Shell Protocol",
			Deprecated:  false,
		}

		data, err := json.Marshal(app)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":"SSH"`)
		assert.Contains(t, string(data), `"deprecated":false`)
	})

	t.Run("NetworkApplications JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": "FTP",
			"description": "File Transfer Protocol",
			"deprecated": false,
			"parentCategory": "FILE_TRANSFER"
		}`

		var app networkapplications.NetworkApplications
		err := json.Unmarshal([]byte(jsonData), &app)
		require.NoError(t, err)

		assert.Equal(t, "FTP", app.ID)
		assert.Equal(t, "FILE_TRANSFER", app.ParentCategory)
	})
}

func TestNetworkApplicationGroups_Structure(t *testing.T) {
	t.Parallel()

	t.Run("NetworkApplicationGroups JSON marshaling", func(t *testing.T) {
		group := networkapplicationgroups.NetworkApplicationGroups{
			ID:          12345,
			Name:        "Remote Access",
			Description: "Remote access applications",
		}

		data, err := json.Marshal(group)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"id":12345`)
	})

	t.Run("NetworkApplicationGroups JSON unmarshaling", func(t *testing.T) {
		jsonData := `{
			"id": 54321,
			"name": "File Sharing",
			"description": "File sharing applications",
			"networkApplications": ["FTP", "SFTP", "SCP"]
		}`

		var group networkapplicationgroups.NetworkApplicationGroups
		err := json.Unmarshal([]byte(jsonData), &group)
		require.NoError(t, err)

		assert.Equal(t, 54321, group.ID)
		assert.Len(t, group.NetworkApplications, 3)
	})
}

func TestNetworkServices_ResponseParsing(t *testing.T) {
	t.Parallel()

	t.Run("Parse network services list", func(t *testing.T) {
		jsonResponse := `[
			{"id": 1, "name": "HTTP", "type": "PREDEFINED"},
			{"id": 2, "name": "HTTPS", "type": "PREDEFINED"},
			{"id": 3, "name": "Custom", "type": "CUSTOM"}
		]`

		var services []networkservices.NetworkServices
		err := json.Unmarshal([]byte(jsonResponse), &services)
		require.NoError(t, err)

		assert.Len(t, services, 3)
	})

	t.Run("Parse network applications list", func(t *testing.T) {
		jsonResponse := `[
			{"id": "SSH", "deprecated": false},
			{"id": "TELNET", "deprecated": true},
			{"id": "FTP", "deprecated": false}
		]`

		var apps []networkapplications.NetworkApplications
		err := json.Unmarshal([]byte(jsonResponse), &apps)
		require.NoError(t, err)

		assert.Len(t, apps, 3)
		assert.True(t, apps[1].Deprecated)
	})
}

