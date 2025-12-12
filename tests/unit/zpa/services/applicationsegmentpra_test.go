// Package unit provides unit tests for ZPA Application Segment PRA service
package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AppSegmentPRA represents the PRA app segment for testing
type AppSegmentPRA struct {
	ID                        string              `json:"id,omitempty"`
	Name                      string              `json:"name,omitempty"`
	Description               string              `json:"description,omitempty"`
	DomainNames               []string            `json:"domainNames,omitempty"`
	Enabled                   bool                `json:"enabled"`
	PassiveHealthEnabled      bool                `json:"passiveHealthEnabled"`
	SelectConnectorCloseToApp bool                `json:"selectConnectorCloseToApp"`
	DoubleEncrypt             bool                `json:"doubleEncrypt"`
	BypassType                string              `json:"bypassType,omitempty"`
	MatchStyle                string              `json:"matchStyle,omitempty"`
	FQDNDnsCheck              bool                `json:"fqdnDnsCheck"`
	HealthCheckType           string              `json:"healthCheckType,omitempty"`
	IsCnameEnabled            bool                `json:"isCnameEnabled"`
	IpAnchored                bool                `json:"ipAnchored"`
	HealthReporting           string              `json:"healthReporting,omitempty"`
	IcmpAccessType            string              `json:"icmpAccessType,omitempty"`
	SegmentGroupID            string              `json:"segmentGroupId"`
	SegmentGroupName          string              `json:"segmentGroupName,omitempty"`
	TCPKeepAlive              string              `json:"tcpKeepAlive,omitempty"`
	TCPPortRanges             []string            `json:"tcpPortRanges,omitempty"`
	UDPPortRanges             []string            `json:"udpPortRanges,omitempty"`
	DefaultIdleTimeout        string              `json:"defaultIdleTimeout,omitempty"`
	DefaultMaxAge             string              `json:"defaultMaxAge,omitempty"`
	PRAApps                   []PRAApps           `json:"praApps"`
	CommonAppsDto             PRACommonAppsDto    `json:"commonAppsDto"`
	MicroTenantID             string              `json:"microtenantId,omitempty"`
	MicroTenantName           string              `json:"microtenantName,omitempty"`
	CreationTime              string              `json:"creationTime,omitempty"`
	ModifiedBy                string              `json:"modifiedBy,omitempty"`
	ModifiedTime              string              `json:"modifiedTime,omitempty"`
}

// PRAApps represents PRA app configuration
type PRAApps struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	AppID               string `json:"appId"`
	ApplicationPort     string `json:"applicationPort,omitempty"`
	ApplicationProtocol string `json:"applicationProtocol,omitempty"`
	CertificateID       string `json:"certificateId,omitempty"`
	CertificateName     string `json:"certificateName,omitempty"`
	ConnectionSecurity  string `json:"connectionSecurity,omitempty"`
	Hidden              bool   `json:"hidden"`
	Portal              bool   `json:"portal"`
	Description         string `json:"description,omitempty"`
	Domain              string `json:"domain,omitempty"`
	Enabled             bool   `json:"enabled"`
	MicroTenantID       string `json:"microtenantId,omitempty"`
}

// PRACommonAppsDto represents common apps configuration
type PRACommonAppsDto struct {
	AppsConfig     []PRAAppsConfig `json:"appsConfig,omitempty"`
	DeletedPraApps []string        `json:"deletedPraApps,omitempty"`
}

// PRAAppsConfig represents app config for PRA
type PRAAppsConfig struct {
	ID                  string   `json:"id,omitempty"`
	AppID               string   `json:"appId"`
	PRAAppID            string   `json:"praAppId"`
	Name                string   `json:"name,omitempty"`
	Description         string   `json:"description,omitempty"`
	Enabled             bool     `json:"enabled,omitempty"`
	AppTypes            []string `json:"appTypes,omitempty"`
	ApplicationPort     string   `json:"applicationPort,omitempty"`
	ApplicationProtocol string   `json:"applicationProtocol,omitempty"`
	ConnectionSecurity  string   `json:"connectionSecurity,omitempty"`
	Domain              string   `json:"domain,omitempty"`
	Hidden              bool     `json:"hidden,omitempty"`
	Portal              bool     `json:"portal,omitempty"`
}

// TestAppSegmentPRA_Structure tests the struct definitions
func TestAppSegmentPRA_Structure(t *testing.T) {
	t.Parallel()

	t.Run("AppSegmentPRA JSON marshaling", func(t *testing.T) {
		pra := AppSegmentPRA{
			ID:               "pra-123",
			Name:             "SSH Access",
			Description:      "SSH remote access",
			DomainNames:      []string{"ssh.example.com"},
			SegmentGroupID:   "sg-001",
			SegmentGroupName: "Remote Access",
			Enabled:          true,
			HealthReporting:  "ON_ACCESS",
			TCPPortRanges:    []string{"22", "22"},
			PRAApps: []PRAApps{
				{
					ID:                  "pa-001",
					Name:                "SSH App",
					ApplicationPort:     "22",
					ApplicationProtocol: "SSH",
					ConnectionSecurity:  "ANY",
					Enabled:             true,
				},
			},
		}

		data, err := json.Marshal(pra)
		require.NoError(t, err)

		var unmarshaled AppSegmentPRA
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, pra.ID, unmarshaled.ID)
		assert.Equal(t, pra.Name, unmarshaled.Name)
		assert.True(t, unmarshaled.Enabled)
		assert.Len(t, unmarshaled.PRAApps, 1)
		assert.Equal(t, "SSH", unmarshaled.PRAApps[0].ApplicationProtocol)
	})

	t.Run("AppSegmentPRA from API response", func(t *testing.T) {
		apiResponse := `{
			"id": "pra-456",
			"name": "RDP Access",
			"description": "Remote desktop access",
			"domainNames": ["rdp.example.com"],
			"segmentGroupId": "sg-002",
			"segmentGroupName": "Desktop Access",
			"enabled": true,
			"bypassType": "NEVER",
			"matchStyle": "EXCLUSIVE",
			"passiveHealthEnabled": true,
			"fqdnDnsCheck": true,
			"healthCheckType": "DEFAULT",
			"isCnameEnabled": true,
			"healthReporting": "ON_ACCESS",
			"tcpPortRanges": ["3389", "3389"],
			"defaultIdleTimeout": "0",
			"defaultMaxAge": "0",
			"praApps": [
				{
					"id": "pa-002",
					"appId": "pra-456",
					"name": "RDP App",
					"enabled": true,
					"applicationPort": "3389",
					"applicationProtocol": "RDP",
					"connectionSecurity": "ANY",
					"domain": "rdp.example.com",
					"hidden": false,
					"portal": false
				}
			],
			"creationTime": "1609459200000",
			"modifiedTime": "1612137600000"
		}`

		var pra AppSegmentPRA
		err := json.Unmarshal([]byte(apiResponse), &pra)
		require.NoError(t, err)

		assert.Equal(t, "pra-456", pra.ID)
		assert.Equal(t, "RDP Access", pra.Name)
		assert.True(t, pra.Enabled)
		assert.Len(t, pra.PRAApps, 1)
		assert.Equal(t, "RDP", pra.PRAApps[0].ApplicationProtocol)
		assert.Equal(t, "3389", pra.PRAApps[0].ApplicationPort)
	})

	t.Run("PRAAppsConfig structure", func(t *testing.T) {
		config := PRAAppsConfig{
			ID:                  "config-001",
			AppID:               "pra-001",
			PRAAppID:            "pa-001",
			Name:                "VNC App",
			Enabled:             true,
			AppTypes:            []string{"SECURE_REMOTE_ACCESS"},
			ApplicationPort:     "5900",
			ApplicationProtocol: "VNC",
			ConnectionSecurity:  "ANY",
			Portal:              false,
		}

		data, err := json.Marshal(config)
		require.NoError(t, err)

		var unmarshaled PRAAppsConfig
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Equal(t, config.PRAAppID, unmarshaled.PRAAppID)
		assert.Equal(t, "VNC", unmarshaled.ApplicationProtocol)
	})
}

// TestAppSegmentPRA_MockServerOperations tests CRUD operations
func TestAppSegmentPRA_MockServerOperations(t *testing.T) {
	t.Run("GET PRA app by ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Contains(t, r.URL.Path, "/application/")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "pra-123",
				"name": "Mock PRA App",
				"enabled": true,
				"praApps": [{"id": "pa-1", "applicationProtocol": "SSH"}]
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Get(server.URL + "/application/pra-123")
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("POST create PRA app", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "new-pra-456",
				"name": "New PRA App",
				"enabled": true
			}`
			w.Write([]byte(response))
		}))
		defer server.Close()

		resp, err := http.Post(server.URL+"/application", "application/json", nil)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DELETE PRA app with forceDelete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "DELETE", r.Method)
			assert.Equal(t, "true", r.URL.Query().Get("forceDelete"))
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		req, _ := http.NewRequest("DELETE", server.URL+"/application/pra-123?forceDelete=true", nil)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}

// TestAppSegmentPRA_SpecialCases tests edge cases
func TestAppSegmentPRA_SpecialCases(t *testing.T) {
	t.Parallel()

	t.Run("PRA protocols", func(t *testing.T) {
		protocols := []string{"SSH", "RDP", "VNC"}

		for _, protocol := range protocols {
			praApp := PRAApps{
				ID:                  "pa-" + protocol,
				Name:                protocol + " App",
				ApplicationProtocol: protocol,
			}

			data, err := json.Marshal(praApp)
			require.NoError(t, err)

			assert.Contains(t, string(data), protocol)
		}
	})

	t.Run("Connection security options", func(t *testing.T) {
		securityOptions := []string{"ANY", "NLA", "TLS", "RDP"}

		for _, security := range securityOptions {
			praApp := PRAApps{
				ID:                 "pa-" + security,
				Name:               security + " Security",
				ConnectionSecurity: security,
			}

			data, err := json.Marshal(praApp)
			require.NoError(t, err)

			assert.Contains(t, string(data), security)
		}
	})

	t.Run("CommonAppsDto with deleted apps", func(t *testing.T) {
		commonApps := PRACommonAppsDto{
			AppsConfig: []PRAAppsConfig{
				{ID: "config-1", Name: "Active App", Enabled: true},
			},
			DeletedPraApps: []string{"pa-deleted-1", "pa-deleted-2"},
		}

		data, err := json.Marshal(commonApps)
		require.NoError(t, err)

		var unmarshaled PRACommonAppsDto
		err = json.Unmarshal(data, &unmarshaled)
		require.NoError(t, err)

		assert.Len(t, unmarshaled.DeletedPraApps, 2)
	})

	t.Run("Portal PRA app", func(t *testing.T) {
		praApp := PRAApps{
			ID:     "pa-123",
			Name:   "Portal App",
			Portal: true,
			Hidden: false,
		}

		data, err := json.Marshal(praApp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"portal":true`)
	})
}

